# Flexkube Provider

The Flexkube provider allows to create and manage Kubernetes cluster components using [libflexkube](https://github.com/flexkube/libflexkube). With this provider, you can create containers like etcd or kubelet on remote machines over SSH using Docker container runtime.

This provider also provides `flexkube_helm_release` resource, so you can use it to manage cluster-essential workloads like CNI plugins or [CoreDNS](https://coredns.io/).

## Example Usage

```hcl
terraform {
  required_providers {
    sshcommand = {
      source  = "flexkube/flexkube"
      version = "0.3.2"
    }
  }
}

variable "ip" {}

variable "name" {
  default = "member01"
}

resource "flexkube_pki" "pki" {
  etcd {
    peers = {
      var.name = var.ip
    }

    servers = {
      var.name = var.ip
    }

    client_cns = ["root"]
  }
}

resource "flexkube_etcd_cluster" "etcd" {
  pki_yaml = flexkube_pki.pki.state_yaml

  member {
    name           = var.name
    peer_address   = var.ip
    server_address = var.ip
  }
}
```

## Argument Reference

This provider currently takes no arguments.
