package flexkube

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/flexkube/libflexkube/pkg/container"
	"github.com/flexkube/libflexkube/pkg/container/runtime/docker"
	"github.com/flexkube/libflexkube/pkg/container/types"
	"github.com/flexkube/libflexkube/pkg/host"
	"github.com/flexkube/libflexkube/pkg/host/transport/direct"
)

func TestContainersStateMarshal(t *testing.T) { //nolint:funlen
	c := container.ContainersState{
		"foo": &container.HostConfiguredContainer{
			Container: container.Container{
				Config: types.ContainerConfig{},
				Runtime: container.RuntimeConfig{
					Docker: docker.DefaultConfig(),
				},
			},
			ConfigFiles: map[string]string{
				"/foo": "bar",
			},
			Host: host.Host{
				DirectConfig: &direct.Config{},
			},
		},
	}

	var s []interface{}

	e := []interface{}{
		map[string]interface{}{
			"name": "foo",
			"config_files": map[string]interface{}{
				"/foo": "bar",
			},
			"host": []interface{}{
				map[string]interface{}{
					"direct": []interface{}{
						map[string]interface{}{},
					},
				},
			},
			"container": []interface{}{
				map[string]interface{}{
					"config": []interface{}{
						map[string]interface{}{
							"name":         "",
							"image":        "",
							"privileged":   false,
							"args":         s,
							"entrypoint":   s,
							"port":         []interface{}{},
							"mount":        []interface{}{},
							"network_mode": "",
							"pid_mode":     "",
							"ipc_mode":     "",
							"user":         "",
							"group":        "",
						},
					},
					"runtime": []interface{}{
						map[string]interface{}{
							"docker": []interface{}{
								map[string]interface{}{
									"host": "unix:///var/run/docker.sock",
								},
							},
						},
					},
				},
			},
		},
	}

	if diff := cmp.Diff(containersStateMarshal(c, false), e); diff != "" {
		t.Errorf("Unexpected diff:\n%s", diff)
	}
}

func TestContainersStateUnmarshal(t *testing.T) {
	c := &container.ContainersState{
		"foo": &container.HostConfiguredContainer{
			Container: container.Container{
				Config: types.ContainerConfig{},
				Runtime: container.RuntimeConfig{
					Docker: docker.DefaultConfig(),
				},
			},
			ConfigFiles: map[string]string{
				"/foo": "bar",
			},
			Host: host.Host{
				DirectConfig: &direct.Config{},
			},
		},
	}

	e := hostConfiguredContainerMarshaled()

	if diff := cmp.Diff(containersStateUnmarshal(e), c); diff != "" {
		t.Errorf("Unexpected diff:\n%s", diff)
	}
}
