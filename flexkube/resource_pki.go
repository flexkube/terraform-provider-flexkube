package flexkube

import (
	"context"
	"fmt"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"sigs.k8s.io/yaml"

	"github.com/flexkube/libflexkube/pkg/pki"
)

func resourcePKI() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePKICreate,
		ReadContext:   resourcePKIRead,
		DeleteContext: resourcePKIDelete,
		UpdateContext: resourcePKICreate,
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
			p, err := getPKI(d)
			if err != nil {
				return fmt.Errorf("getting PKI: %w", err)
			}

			// Generate PKI to find pending changes.
			if err := p.Generate(); err != nil {
				return fmt.Errorf("generating PKI: %w", err)
			}

			pkiState := &pki.PKI{}
			if err := yaml.Unmarshal([]byte(d.Get("state_yaml").(string)), pkiState); err != nil {
				return fmt.Errorf("malformed %q field, unable to parse as PKI: %w", "state_yaml", err)
			}

			if diff := cmp.Diff(pkiState, p); diff != "" {
				if err := d.SetNewComputed("state"); err != nil {
					return fmt.Errorf("failed setting key %q as new computed: %w", "state", err)
				}

				if err := d.SetNewComputed("state_sensitive"); err != nil {
					return fmt.Errorf("failed setting key %q as new computed: %w", "state_sensitive", err)
				}

				if err := d.SetNewComputed("state_yaml"); err != nil {
					return fmt.Errorf("failed setting key %q as new computed: %w", "state_yaml", err)
				}
			}

			return nil
		},
		Schema: map[string]*schema.Schema{
			"certificate":     certificateBlockSchema(false),
			"root_ca":         certificateBlockSchema(false),
			"etcd":            etcdSchema(false),
			"kubernetes":      kubernetesSchema(false),
			"state_yaml":      sensitiveString(true),
			"state":           pkiSchema(true, false),
			"state_sensitive": pkiSchema(true, true),
		},
	}
}

func getPKI(d getter) (*pki.PKI, error) {
	p := &pki.PKI{
		RootCA:     certificateUnmarshal(d.Get("root_ca")),
		Etcd:       etcdUnmarshal(d.Get("etcd")),
		Kubernetes: kubernetesUnmarshal(d.Get("kubernetes")),
	}

	if c := certificateUnmarshal(d.Get("certificate")); c != nil {
		p.Certificate = *c
	}

	pkiState := &pki.PKI{}
	if err := yaml.Unmarshal([]byte(d.Get("state_yaml").(string)), pkiState); err != nil {
		return nil, fmt.Errorf("malformed %q field, unable to parse as PKI: %w", "state_yaml", err)
	}

	b, err := yaml.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("failed converting PKI to YAML: %w", err)
	}

	if err := yaml.Unmarshal(b, &pkiState); err != nil {
		return nil, fmt.Errorf("failed unmarshaling PKI: %w", err)
	}

	return pkiState, nil
}

func savePKI(d *schema.ResourceData, p *pki.PKI) error {
	b, err := yaml.Marshal(p)
	if err != nil {
		return fmt.Errorf("failed converting PKI to YAML: %w", err)
	}

	props := map[string]interface{}{
		"state_yaml":      string(b),
		"state":           pkiMarshal(p, true),
		"state_sensitive": pkiMarshal(p, false),
	}

	for k, v := range props {
		if err := d.Set(k, v); err != nil {
			return fmt.Errorf("failed setting property %q: %w", k, err)
		}
	}

	return nil
}

func resourcePKICreate(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	p, err := getPKI(d)
	if err != nil {
		return diagFromErr(fmt.Errorf("getting PKI: %w", err))
	}

	if err := p.Generate(); err != nil {
		return diagFromErr(err)
	}

	b, err := yaml.Marshal(p)
	if err != nil {
		return diagFromErr(err)
	}

	if d.IsNewResource() {
		d.SetId(sha256sum(b))
	}

	return diagFromErr(savePKI(d, p))
}

func resourcePKIDelete(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.SetId("")

	return nil
}

func resourcePKIRead(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}
