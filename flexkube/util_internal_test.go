package flexkube

import (
	"context"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/flexkube/libflexkube/pkg/container"
	"github.com/flexkube/libflexkube/pkg/container/resource"
	"github.com/flexkube/libflexkube/pkg/container/runtime/docker"
	"github.com/flexkube/libflexkube/pkg/container/types"
	"github.com/flexkube/libflexkube/pkg/host"
	"github.com/flexkube/libflexkube/pkg/host/transport/direct"
)

const (
	tfACC = "TF_ACC"
)

// saveState() tests.
//
//nolint:paralleltest // This function modifies environment variables, which are global.
func TestSaveStateBadScheme(t *testing.T) {
	r := resourceContainers()
	delete(r.Schema, stateYAMLSchemaKey)

	t.Setenv(tfACC, "")

	d := r.Data(&terraform.InstanceState{})

	if err := saveState(d, container.ContainersState{}, containersUnmarshal, nil); err == nil {
		t.Fatalf("save state should fail when called on bad scheme")
	}
}

// resourceDelete() tests.
func TestResourceDeleteRuntimeFail(t *testing.T) {
	t.Parallel()

	// Get the resource object we will work on.
	r := resourceContainers()
	r.DeleteContext = resourceDelete(containersUnmarshal, stateSensitiveSchemaKey)

	// Prepare some fake state.
	s := container.ContainersState{
		"foo": &container.HostConfiguredContainer{
			Host: host.Host{
				DirectConfig: &direct.Config{},
			},
			Container: container.Container{
				Runtime: container.RuntimeConfig{
					Docker: &docker.Config{
						Host: "unix:///nonexistent",
					},
				},
				Config: types.ContainerConfig{
					Name:  "foo",
					Image: "busybox:latest",
				},
				Status: &types.ContainerStatus{
					ID:     "foo",
					Status: "running",
				},
			},
		},
	}

	// Create raw configuration to create ResourceData object.
	raw := map[string]interface{}{
		stateSensitiveSchemaKey: containersStateMarshal(s, false),
	}

	// Create ResourceData object.
	d := schema.TestResourceDataRaw(t, r.Schema, raw)

	// Mark newly created object as created, so it's state is persisted.
	d.SetId("foo")

	// Create new ResourceData from the state, so it's persisted and there is no diff included.
	dn := r.Data(d.State())

	// Finally, try to call Delete.
	if err := r.DeleteContext(context.TODO(), dn, nil); err != nil {
		t.Fatalf("destroying should work with unreachable runtime: %v", err)
	}
}

func TestResourceDeleteEmpty(t *testing.T) {
	t.Parallel()

	r := resourceContainers()
	r.DeleteContext = resourceDelete(containersUnmarshal, stateSensitiveSchemaKey)

	s := container.ContainersState{
		"foo": &container.HostConfiguredContainer{
			Host: host.Host{
				DirectConfig: &direct.Config{},
			},
			Container: container.Container{
				Runtime: container.RuntimeConfig{
					Docker: &docker.Config{
						Host: "unix:///nonexistent",
					},
				},
				Config: types.ContainerConfig{
					Name:  "foo",
					Image: "busybox:latest",
				},
			},
		},
	}

	d := r.Data(&terraform.InstanceState{})
	if err := d.Set(stateSensitiveSchemaKey, containersStateMarshal(container.ContainersState{}, false)); err != nil {
		t.Fatalf("Failed writing: %v", err)
	}

	if err := d.Set("host_configured_container", containersStateMarshal(s, false)); err != nil {
		t.Fatalf("writing containers configuration to state failed: %v", err)
	}

	if err := r.DeleteContext(context.TODO(), d, nil); !strings.Contains(err[0].Summary, "Is the docker daemon running") {
		t.Fatalf("destroying should fail for unreachable runtime")
	}
}

func TestResourceDeleteEmptyState(t *testing.T) {
	t.Parallel()

	r := resourceContainers()

	if err := r.DeleteContext(context.TODO(), r.Data(&terraform.InstanceState{}), nil); err == nil {
		t.Fatalf("initializing from empty state should fail")
	}
}

//nolint:paralleltest // This function modifies environment variables, which are global.
func TestResourceDeleteBadKey(t *testing.T) {
	r := resourceContainers()
	r.DeleteContext = resourceDelete(containersUnmarshal, "foo")

	t.Setenv(tfACC, "")

	if err := r.DeleteContext(context.TODO(), r.Data(&terraform.InstanceState{}), nil); err == nil {
		t.Fatalf("emptying key not existing in scheme should fail")
	}
}

// newResource() tests.
func TestNewResourceFailRefresh(t *testing.T) {
	t.Parallel()

	cc := &resource.Containers{
		State: container.ContainersState{
			"foo": &container.HostConfiguredContainer{
				Host: host.Host{
					DirectConfig: &direct.Config{},
				},
				Container: container.Container{
					Runtime: container.RuntimeConfig{
						Docker: &docker.Config{
							Host: "unix:///nonexistent",
						},
					},
					Config: types.ContainerConfig{
						Name:  "foo",
						Image: "busybox:latest",
					},
					Status: &types.ContainerStatus{
						ID:     "foo",
						Status: "running",
					},
				},
			},
		},
	}

	if _, err := newResource(cc, true); err != nil {
		t.Fatalf("should not return any errors: %v", err)
	}
}

// resourceCreate() tests.
func TestResourceCreate(t *testing.T) {
	t.Parallel()

	r := resourceContainers()

	s := container.ContainersState{
		"foo": &container.HostConfiguredContainer{
			Host: host.Host{
				DirectConfig: &direct.Config{},
			},
			Container: container.Container{
				Runtime: container.RuntimeConfig{
					Docker: &docker.Config{
						Host: "unix:///nonexistent",
					},
				},
				Config: types.ContainerConfig{
					Name:  "foo",
					Image: "busybox:latest",
				},
			},
		},
	}

	d := r.Data(&terraform.InstanceState{})
	if err := d.Set("host_configured_container", containersStateMarshal(s, false)); err != nil {
		t.Fatalf("writing containers configuration to state failed: %v", err)
	}

	if err := r.CreateContext(context.TODO(), d, nil); !strings.Contains(err[0].Summary, "Is the docker daemon running") {
		t.Fatalf("creating should fail for unreachable runtime, got: %v", err)
	}
}

func TestResourceCreateFailInitialize(t *testing.T) {
	t.Parallel()

	r := resourceContainers()

	s := container.ContainersState{
		"foo": &container.HostConfiguredContainer{
			Host: host.Host{
				DirectConfig: &direct.Config{},
			},
			Container: container.Container{
				Runtime: container.RuntimeConfig{
					Docker: &docker.Config{
						Host: "unix:///nonexistent",
					},
				},
				Config: types.ContainerConfig{
					Name:  "",
					Image: "busybox:latest",
				},
			},
		},
	}

	d := r.Data(&terraform.InstanceState{})
	if err := d.Set("host_configured_container", containersStateMarshal(s, false)); err != nil {
		t.Fatalf("writing containers configuration to state failed: %v", err)
	}

	if err := r.CreateContext(context.TODO(), d, nil); !strings.Contains(err[0].Summary, "name must be set") {
		t.Fatalf("creating should fail for unreachable runtime, got: %v", err)
	}
}

// resourceRead() tests.
func TestResourceRead(t *testing.T) {
	t.Parallel()

	r := resourceContainers()

	s := container.ContainersState{
		"foo": &container.HostConfiguredContainer{
			Host: host.Host{
				DirectConfig: &direct.Config{},
			},
			Container: container.Container{
				Runtime: container.RuntimeConfig{
					Docker: &docker.Config{
						Host: "unix:///nonexistent",
					},
				},
				Config: types.ContainerConfig{
					Name:  "foo",
					Image: "busybox:latest",
				},
			},
		},
	}

	d := r.Data(&terraform.InstanceState{})
	if err := d.Set("host_configured_container", containersStateMarshal(s, false)); err != nil {
		t.Fatalf("writing containers configuration to state failed: %v", err)
	}

	if err := r.ReadContext(context.TODO(), d, nil); err != nil {
		t.Fatalf("reading with no previous state should succeed, got: %v", err)
	}
}

func TestResourceReadFailInitialize(t *testing.T) {
	t.Parallel()

	r := resourceContainers()

	s := container.ContainersState{
		"foo": &container.HostConfiguredContainer{
			Host: host.Host{
				DirectConfig: &direct.Config{},
			},
			Container: container.Container{
				Runtime: container.RuntimeConfig{
					Docker: &docker.Config{
						Host: "unix:///nonexistent",
					},
				},
				Config: types.ContainerConfig{
					Name:  "",
					Image: "busybox:latest",
				},
			},
		},
	}

	d := r.Data(&terraform.InstanceState{})
	if err := d.Set("host_configured_container", containersStateMarshal(s, false)); err != nil {
		t.Fatalf("writing containers configuration to state failed: %v", err)
	}

	if err := r.ReadContext(context.TODO(), d, nil); err == nil {
		t.Fatalf("read should check for initialize errors and fail")
	}
}

func TestStringMapMarshal(t *testing.T) {
	t.Parallel()

	f := map[string]string{
		"/foo": "bar",
	}

	e := map[string]interface{}{
		"/foo": "bar",
	}

	if diff := cmp.Diff(stringMapMarshal(f, false), e); diff != "" {
		t.Fatalf("Unexpected diff: %s", diff)
	}
}

func TestStringMapMarshalSensitive(t *testing.T) {
	t.Parallel()

	f := map[string]string{
		"/foo": "bar",
	}

	e := map[string]interface{}{
		"/foo": sha256sum([]byte("bar")),
	}

	if diff := cmp.Diff(stringMapMarshal(f, true), e); diff != "" {
		t.Fatalf("Unexpected diff: %s", diff)
	}
}

func TestStringMapMarshalDontChecksumEmpty(t *testing.T) {
	t.Parallel()

	f := map[string]string{
		"/foo": "",
	}

	e := map[string]interface{}{
		"/foo": "",
	}

	if diff := cmp.Diff(stringMapMarshal(f, true), e); diff != "" {
		t.Fatalf("Unexpected diff: %s", diff)
	}
}

func TestStringMapUnmarshal(t *testing.T) {
	t.Parallel()

	i := map[string]interface{}{
		"/foo": "bar",
	}

	e := map[string]string{
		"/foo": "bar",
	}

	if diff := cmp.Diff(stringMapUnmarshal(i), e); diff != "" {
		t.Fatalf("Unexpected diff: %s", diff)
	}
}

func TestStringMapUnmarshalEmpty(t *testing.T) {
	t.Parallel()

	var e map[string]string

	if diff := cmp.Diff(stringMapUnmarshal(nil), e); diff != "" {
		t.Fatalf("Unexpected diff: %s", diff)
	}
}
