package flexkube

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"sigs.k8s.io/yaml"

	"github.com/flexkube/libflexkube/pkg/container"
	"github.com/flexkube/libflexkube/pkg/container/resource"
	"github.com/flexkube/libflexkube/pkg/pki"
	"github.com/flexkube/libflexkube/pkg/types"
)

const (
	// If schema is Type: schema.TypeList and Elem: &schema.Resource,
	// MaxItems should be set to this value, to treat the property as a
	// standalone, singleton block.
	blockMaxItems           = 1
	stateSensitiveSchemaKey = "state_sensitive"
	stateSchemaKey          = "state"
	stateYAMLSchemaKey      = "state_yaml"
	configYAMLSchemaKey     = "config_yaml"
)

type getter interface {
	Get(key string) interface{}
	GetChange(key string) (interface{}, interface{})
	GetOk(key string) (interface{}, bool)
}

type unmarshalF func(getter, bool) types.ResourceConfig

func requiredString(computed bool) *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeString,
		Required: !computed,
		Computed: computed,
	}
}

func optionalString(computed bool) *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		Computed: computed,
	}
}

func sensitiveString(computed bool) *schema.Schema {
	return &schema.Schema{
		Type:      schema.TypeString,
		Optional:  true,
		Computed:  computed,
		Sensitive: true,
	}
}

func optionalStringList(computed bool) *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Computed: computed,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
}

func optionalBool(computed bool) *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Computed: computed,
	}
}

func optionalInt(computed bool) *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeInt,
		Optional: true,
		Computed: computed,
	}
}

func optionalMap(computed bool, elem func(bool) *schema.Resource) *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Computed: computed,
		Elem:     elem(computed),
	}
}

func optionalMapPrimitive(computed bool, elem func(bool) *schema.Schema) *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
		Computed: computed,
		Elem:     elem(computed),
	}
}

func requiredBlock(computed bool, elem func(bool) *schema.Resource) *schema.Schema {
	s := &schema.Schema{
		Type:     schema.TypeList,
		Required: !computed,
		Computed: computed,
		Elem:     elem(computed),
	}

	if !computed {
		s.MaxItems = blockMaxItems
	}

	return s
}

func optionalBlock(computed, sensitive bool, elem func(bool) map[string]*schema.Schema) *schema.Schema {
	s := &schema.Schema{
		Type:      schema.TypeList,
		Computed:  computed,
		Sensitive: sensitive,
		Elem: &schema.Resource{
			Schema: elem(computed),
		},
	}

	if !computed {
		s.MaxItems = blockMaxItems
		s.Optional = true
	}

	return s
}

func requiredList(computed, sensitive bool, elem func(bool) *schema.Resource) *schema.Schema {
	return &schema.Schema{
		Type:      schema.TypeList,
		Computed:  computed,
		Sensitive: sensitive,
		Required:  !computed,
		Elem:      elem(computed),
	}
}

func optionalList(computed bool, elem func(bool) *schema.Resource) *schema.Schema {
	s := &schema.Schema{
		Type:     schema.TypeList,
		Computed: computed,
		Elem:     elem(computed),
	}

	if !computed {
		s.Optional = true
	}

	return s
}

func sha256sum(data []byte) string {
	return fmt.Sprintf("%x", sha256.Sum256(data))
}

// withCommonFields adds common fields to the resource scheme. This should be used for resources,
// which implements types.Resource, so their state and configuration can be exposed to the user in
// standard way.
func withCommonFields(s map[string]*schema.Schema) map[string]*schema.Schema {
	// Existing state of the configuration, used for operations.
	// This filed is entirely marked as sensitive, to prevent leaking secrets when running
	// plan etc.
	s[stateSensitiveSchemaKey] = hostConfiguredContainerSchema(true, true)
	// This state has secrets stripped out and will be presented as a diff to the user.
	s[stateSchemaKey] = hostConfiguredContainerSchema(true, false)
	// Sensitive state in YAML format, which can be saved to disk and used with CLI tools.
	s[stateYAMLSchemaKey] = sensitiveString(true)
	// Sensitive user configuration in YAML format, which can be saved to disk and used with
	// CLI tools as well.
	s[configYAMLSchemaKey] = sensitiveString(true)

	return s
}

func getState(d getter) *container.ContainersState {
	ss, _ := d.GetChange(stateSensitiveSchemaKey)

	return containersStateUnmarshal(ss)
}

func newResource(c types.ResourceConfig, refresh bool) (types.Resource, error) {
	// Validate the configuration.
	r, err := c.New()
	if err != nil {
		return nil, fmt.Errorf("creating resource: %w", err)
	}

	if !refresh {
		return r, nil
	}

	// Get current state of the containers.
	if err := r.CheckCurrentState(); err != nil {
		return nil, fmt.Errorf("checking current state: %w", err)
	}

	return r, nil
}

func initialize(d getter, uf unmarshalF, refresh bool) (types.Resource, error) {
	c := uf(d, true)

	r, err := newResource(c, refresh)
	if err != nil {
		return nil, fmt.Errorf("initializing resource: %w", err)
	}

	return r, nil
}

func resourceCreate(uf unmarshalF) func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		// Create Containers object.
		c, err := initialize(d, uf, true)
		if err != nil {
			return diagFromErr(fmt.Errorf("initializing configuration: %w", err))
		}

		// Deploy changes.
		deployErr := c.Deploy()

		// If there was at least one container created, set the ID to mark, that resource has been at least partially
		// created.
		// If the ID is already set, then also don't update it, as there is no need for that.
		if d.IsNewResource() && len(c.Containers().ToExported().PreviousState) != 0 {
			d.SetId(sha256sum([]byte(cmp.Diff(nil, uf(d, true)))))
		}

		return diagFromErr(saveState(d, c.Containers().ToExported().PreviousState, uf, deployErr))
	}
}

func resourceRead(uf unmarshalF) func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		c, err := initialize(d, uf, true)
		if err != nil {
			return diagFromErr(fmt.Errorf("initializing configuration: %w", err))
		}

		// If there is nothing in the current state, mark the resource as destroyed.
		if len(c.Containers().ToExported().PreviousState) == 0 {
			d.SetId("")
		}

		return diagFromErr(saveState(d, c.Containers().ToExported().PreviousState, uf, nil))
	}
}

func diagFromErr(err error) diag.Diagnostics {
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceDelete(uf unmarshalF, key string) schema.DeleteContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		// Reset user configuration to indicate, that we destroy everything.
		if err := d.Set(key, []interface{}{}); err != nil {
			return diagFromErr(fmt.Errorf("trigging a destroy: %w", err))
		}

		// Create Containers object.
		c, err := initialize(d, uf, true)
		if err != nil {
			return diagFromErr(fmt.Errorf("initializing configuration: %w", err))
		}

		// Deploy changes.
		deployErr := c.Deploy()

		// If deployment succeeded, we are done.
		if deployErr == nil {
			d.SetId("")

			return nil
		}

		return diagFromErr(saveState(d, c.Containers().ToExported().PreviousState, uf, deployErr))
	}
}

// prepareDiff generates all information, which needs to be written by resourceDiff.
func prepareDiff(d getter, uf unmarshalF) (string, types.Resource, map[string]interface{}, error) {
	cy, err := configYaml(d, uf)
	if err != nil {
		return "", nil, nil, fmt.Errorf("getting config in YAML format: %w", err)
	}

	// Initialize resource, but there is no need to refresh the state, as we will only write
	// desired states and configuration anyway.
	r, err := initialize(d, uf, false)
	if err != nil {
		return "", nil, nil, fmt.Errorf("initializing configuration: %w", err)
	}

	statesMap, err := states(r.Containers().DesiredState())
	if err != nil {
		return "", nil, nil, fmt.Errorf("getting serialized states: %w", err)
	}

	return cy, r, statesMap, nil
}

// resourceDiff customize resource diff for resources implementing types.Resource.
// It makes sure, that all fields are marked correctly and that diff will show valuable
// and secure output to the user.
func resourceDiff(uf unmarshalF) func(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
	return func(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
		if changes := d.GetChangedKeysPrefix(""); len(changes) > 0 {
			return setNewComputed(d)
		}

		cy, r, _, err := prepareDiff(d, uf)
		if err != nil {
			// If the configuration has not been fully populated yet, some of required fields might be empty, so validation
			// and initialization will fail, so just don't do anything if that happens.
			//
			//nolint:nilerr // It's deliberate to ignore the error here to not proceed with computation.
			return nil
		}

		// If there is some change to state, like container needs to be added or created, where we don't know exact
		// value which it will take, as we cannot know container ID in advance, then mark field as computed, so if
		// other resources takes this field as an input, they will get triggered.
		stateDiff := cmp.Diff(r.Containers().ToExported().PreviousState, r.Containers().DesiredState())

		// If fields, which builds the config are not known before the execution, e.g. when you include certificate
		// generated by Terraform, then mark 'config_yaml' field as new computed, to avoid producing inconsistent
		// state. If the config does not differ, then just write it as known value, as it should not produce any diff.
		// We still need to write it to mark the field, as it will get the value.
		configDiff := cmp.Diff(d.Get(configYAMLSchemaKey).(string), cy)

		if stateDiff != "" || configDiff != "" {
			return setNewComputed(d)
		}

		return nil
	}
}

func setNewComputed(d *schema.ResourceDiff) error {
	setNewComputed := []string{
		stateSensitiveSchemaKey,
		stateYAMLSchemaKey,
		stateSchemaKey,
		configYAMLSchemaKey,
	}

	// Now apply selected fields.
	for _, k := range setNewComputed {
		if err := d.SetNewComputed(k); err != nil {
			return fmt.Errorf("setting key %q as new computed: %w", k, err)
		}
	}

	return nil
}

func states(s container.ContainersState) (map[string]interface{}, error) {
	states := map[string]interface{}{
		stateSensitiveSchemaKey: stateSensitiveMarshal(s),
		stateSchemaKey:          stateMarshal(s),
	}

	sy, err := stateYaml(s)
	if err != nil {
		return nil, fmt.Errorf("converting state to YAML: %w", err)
	}

	states[stateYAMLSchemaKey] = sy

	return states, nil
}

func saveState(d *schema.ResourceData, s container.ContainersState, uf unmarshalF, origErr error) error {
	states, err := states(s)
	if err != nil {
		return fmt.Errorf("getting serialized states: %w", err)
	}

	// If config is build on values passed from other resources, we won't know the exact content during
	// planning, so we need to make sure, that after creating the right content is written to the field.
	cy, err := configYaml(d, uf)
	if err != nil {
		return fmt.Errorf("getting config in YAML format: %w", err)
	}

	states[configYAMLSchemaKey] = cy

	for k, v := range states {
		if err := d.Set(k, v); err != nil {
			return fmt.Errorf("persisting key %q to state: %w", k, err)
		}
	}

	return origErr
}

// stateYaml returns data in format compatible for writing to 'state_yaml' field.
func stateYaml(s container.ContainersState) (interface{}, error) {
	cc := &resource.Containers{
		State: s,
	}

	ccy, err := yaml.Marshal(cc)
	if err != nil {
		return "", fmt.Errorf("serializing state: %w", err)
	}

	return string(ccy), nil
}

// stateSensitiveMarshal returns data in format compatible for writing to 'state_sensitive' field.
func stateSensitiveMarshal(s container.ContainersState) interface{} {
	return containersStateMarshal(s, false)
}

// stateMarshal returns data in format compatible for writing to 'state' field.
func stateMarshal(s container.ContainersState) interface{} {
	return containersStateMarshal(s, true)
}

// configYaml returns data in format compatible for writing to 'config_yaml' field.
func configYaml(d getter, uf unmarshalF) (string, error) {
	rc := uf(d, false)

	b, err := yaml.Marshal(rc)
	if err != nil {
		return "", fmt.Errorf("serializing generated configuration: %w", err)
	}

	return string(b), nil
}

func stringMapUnmarshal(i interface{}) map[string]string {
	r := map[string]string{}

	if i == nil {
		return nil
	}

	j := i.(map[string]interface{})

	if len(j) == 0 {
		return nil
	}

	for k, v := range j {
		r[k] = v.(string)
	}

	return r
}

func stringListUnmarshal(i interface{}) []string {
	r := []string{}

	if i == nil {
		return nil
	}

	j, ok := i.([]interface{})

	if !ok || len(j) == 0 {
		return nil
	}

	for _, v := range j {
		r = append(r, v.(string))
	}

	return r
}

func stringSliceToInterfaceSlice(i []string) []interface{} {
	//nolint:prealloc // Don't preallocate, as then the diff shows diff between nil and empty slice.
	var o []interface{}

	for _, v := range i {
		o = append(o, v)
	}

	return o
}

func stringMapSchema(computed, sensitive bool) *schema.Schema {
	s := optionalMapPrimitive(computed, func(computed bool) *schema.Schema {
		return &schema.Schema{
			Type: schema.TypeString,
		}
	})

	s.Sensitive = sensitive

	return s
}

func stringMapMarshal(c map[string]string, sensitive bool) interface{} {
	i := map[string]interface{}{}

	for k, v := range c {
		if sensitive && v != "" {
			i[k] = sha256sum([]byte(v))

			continue
		}

		i[k] = v
	}

	return i
}

func unmarshalPKI(d getter) *pki.PKI {
	pki := &pki.PKI{}

	v, ok := d.GetOk("pki_yaml")
	if !ok || v.(string) == "" {
		return pki
	}

	_ = yaml.Unmarshal([]byte(v.(string)), pki)

	return pki
}
