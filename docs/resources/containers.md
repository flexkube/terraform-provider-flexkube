# Containers Resource

This resource allows to manager an arbitrary containers. This resource might be helpful, if your cluster requires some side-containers to run.

## Example Usage

```hcl
resource "flexkube_containers" "foo" {
  host_configured_container {
    name = "bar"

    container {
      config {
        name  = "bazhh"
        image = "nginx"
      }
    }
  }
}
```

## Argument Reference

* `host_configured_container` - (Required) A list of `host_configured_container` blocks as defined below. One block represents one managed container.

---

A `host_configured_container` block supports the following:

* `name` - (Required) Name of the container.

* `container` - (Required) A `container` block as defined below. Defines container configuration.

* `config_files` - (Optional) A map of configuration files to be created for the container on the host. The key is the path on the host to the configuration file and the value is a configuration file content.

* `host` - (Optional) A `host` block as defined above. This block allows configuring on which host the container should be created.

---

A `container` block supports the following:

* `config` - (Required) A `config` block as defined below. Defines container parameters.

* `runtime` - (Optional) A `runtime` block as defined below. Allows configuring container runtime to be used by the container.

---

A `config` block supports the following:

* `name` - (Required) Name of the managed container.

* `image` - (Required) Container image used.

* `privileged` - (Optional) If `true`, the container runs as a privileged process on the host. Defaults to `false`.

* `args` - (Optional) Arguments used for a container.

* `entrypoint` - (Optional) Binary name which runs in the container.

* `port` - (Optional) A `port` block as defined below. Contains ports, which are exposed by the container.

* `mount` - (Optional) A list of `mount` blocks as defined above. Contains information which host paths are mounted into the container.

* `network_mode` - (Optional) Defines what network mode container use. Valid value may depend on used container runtime.

* `pid_mode` - (Optional) Defines in which PID mode container runs. Valid value may depend on used container runtime.

* `ipc_mode` - (Optional) Defines in which IPC mode container runs. Valid value may depend on used container runtime.

* `user` - (Optional) Name of the user or UID used by the container.

* `group` - (Optional) Name of the group or GID used by the container.

---

A `port` block supports the following:

* `port` - (Required) Exposed port number.

* `protocol` - (Required) Exposed protocol.

* `ip` - (Optional) IP on which the port is exposed.

---

A `runtime` block supports the following:

* `docker` - (Optional) A `docker` block as defined below. Container Docker runtime configuration attributes.

---

A `docker` block supports the following:

* `host` - (Optional) URL used to talk to Docker runtime. Defaults to `unix:///var/run/docker.sock`.

## Attribute Reference

* `state` - A list of `host_configured_container` blocks as defined below. This attribute represents generated configuration of the managed containers. Sensitive values like configuration files content, environment variables or SSH password are replaced with SHA256 of the values. To get the actual value, use `state_sensitive` block.

* `state_sensitive` - A list of `host_configured_container` blocks as defined below. This attribute represents generated configuration of the managed containers. This attribute is marked entirely as sensitive, so it won't show up detailed in a plan. To see specific changes in generated configuration, use `state` attribute.

* `state_yaml` - State of created containers in YAML format. Can be dumped to a `state.yaml` file and used together with `flexkube kubelet-pool` command.

* `config_yaml` - Generated configuration in YAML format, which can be used by the `flexkube kubelet-pool` command.

---

A `host_configured_container` block supports the following:

* `name` - Name of the container.

* `container` - A `container` block as defined below. Contains container configuration parameters.

* `config_files` - A map of configuration files which are created for the container on the host. The key is the path on the host to the configuration file and the value is either a configuration file content or it's SHA256, depending if read from `state` or from `state_sensitive` attribute.

* `host` - A `host` block as defined above. Describes on which host the container is created.

---

A `container` block supports the following:

* `config` - A `config` block as defined below. Contains container configuration attributes.

* `runtime` - A `runtime` block as defined below. Includes container runtime configuration.

* `status` - A `status` block as defined below. Contains container status information.

---

A `config` block supports the following:

* `name` - Name of the managed container.

* `image` - Container image used.

* `privileged` - If `true`, the container runs as a privileged process on the host.

* `args` - Arguments used for a container.

* `entrypoint` - Binary name which runs in the container.

* `port` - A `port` block as defined below. Contains ports, which are exposed by the container.

* `mount` - A list of `mount` blocks as defined above. Contains information which host paths are mounted into the container.

* `network_mode` - Defines what network mode container use. Actual value may depend on used container runtime.

* `pid_mode` - Defines in which PID mode container runs. Actual value may depend on used container runtime.

* `ipc_mode` - Defines in which IPC mode container runs. Actual value may depend on used container runtime.

* `user` - Name of the user or UID used by the container.

* `group` - Name of the group or GID used by the container.

---

A `port` block supports the following:

* `ip` - IP on which the port is exposed.

* `port` - Exposed port number.

* `protocol` - Exposed protocol.

---

A `runtime` block supports the following:

* `docker` - A `docker` block as defined below. Container Docker runtime configuration attributes.

---

A `docker` block supports the following:

* `host` - URL used to talk to Docker runtime. Defaults to `unix:///var/run/docker.sock`.

---

A `status` block supports the following:

* `id` - ID of the created container given by used container runtime.

* `status` - Text status of the container. If field is empty, it means that the container does not exist.
