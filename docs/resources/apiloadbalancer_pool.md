# API Load Balancer Pool Resource

This resource allows to create one or more Kubernetes API Load balancer containers on remote hosts over SSH using Docker container runtime.

## Example Usage

```hcl
resource "flexkube_apiloadbalancer_pool" "controllers" {
  name             = "api-loadbalancer-controllers"
  host_config_path = "/etc/haproxy/controllers.cfg"
  bind_address     = "0.0.0.0"
  servers          = ["10.0.0.10:6443"]

  ssh {
    port        = 22
    user        = "core"
  }

	api_load_balancer {
    host {
      ssh {
        address = "10.0.0.10"
      }
    }
  }
}
```

## Argument Reference

* `servers` - (Required) A list of Kubernetes API Server addresses with ports, which should be the target of the load balancer.

* `name` - (Optional) Unique identifier of the load balancers on the host. If you deploy more than one instance of load balancer on a single machine this must be defined to avoid the instances colliding.

* `host_config_path` - (Optional) Path where to store instance configuration file. If you deploy more than one instance of load balancer on a single machine, this must be defined to avoid configuration files collisions.

* `bind_address` - (Optional) Address on which the load balancer should bind for incoming requests.

* `api_load_balancer` - (Required) A `api_load_balancer` block as defined below. This block defines single API Load Balancer instance and can be specified multiple times.

* `image` - (Optional) Docker image with tag to be used to run HAProxy container. Defaults to `libflexkube` [default HAProxy Image](https://github.com/flexkube/libflexkube/blob/master/pkg/defaults/defaults.go#L12).

* `ssh` - (Optional) A `ssh` block as defined below. This block defines global SSH settings shared by all instances.

---

A `api_load_balancer` block supports the following:

* `host` - (Optional) A `host` block as defined below. This block defines where to connect for creating the container.

* `servers` - (Required) A list of Kubernetes API Server addresses with ports, which should be the target of the load balancer.

* `name` - (Optional) Unique identifier of the load balancers on the host. If you deploy more than one instance of load balance
r on a single machine this must be defined to avoid the instances colliding.

* `host_config_path` - (Optional) Path where to store instance configuration file. If you deploy more than one instance of load balancer on a single machine, this must be defined to avoid configuration files collisions.

* `bind_address` - (Optional) Address on which the load balancer should bind for incoming requests.

* `image` - (Optional) Docker image with tag to be used to run HAProxy container. Defaults to `libflexkube` [default HAProxy Image](https://github.com/flexkube/libflexkube/blob/master/pkg/defaults/defaults.go#L12).

---

A `host` block supports the following:

* `direct` - (Optional) A `direct` block as defined below. Mutually exclusive with all other fields in this block. If defined, container will be created on local machine.

* `ssh` - (Optional) A `ssh` block as defined below. Mutually exclusive with all other fields in this block. If defined, container will be created on a remote machine using SSH connection.

---

A `direct` block does not support any arguments.

---

A `ssh` block supports the following:

* `address` - (Required) An address where SSH client should connect to. Can be either hostname of IP address.

* `port` - (Optional) Port where to open SSH connection. Defaults to `22`.

* `user` - (Optional) Username to use when opening SSH connection. Defaults to `root`.

* `password` - (Optional) Password to use for SSH authentication. Can be combined with `private_key` and SSH agent authentication methods.

* `connection_timeout` - (Optional) Duration for how long to wait before connection attempts times out, expressed in [Go Duration format](https://golang.org/pkg/time/#ParseDuration). Defaults to `30s`.

* `retry_timeout` - (Optional) Duration for how long to wait before giving up on connection attempts, expressed in [Go Duration format](https://golang.org/pkg/time/#ParseDuration). Defaults to `60s`.

* `retry_interval` - (Optional) Duration for how long to wait between connection attempts, expressed in [Go Duration format](https://golang.org/pkg/time/#ParseDuration). Defaults to `1s`.

* `private_key` - (Optional) PEM encoded privat key to be used for authentication. Can be combined with `password` and SSH agent authentication methods.

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
