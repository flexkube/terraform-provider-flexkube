resource "null_resource" "controllers" {
  count = var.controllers_count

  triggers = {
    name = format("controller%02d", count.index + 1)
    ip   = cidrhost(var.nodes_cidr, count.index + var.cidr_ips_offset)
    cidr = cidrsubnet(var.pod_cidr, 8, count.index + 2)
  }
}

locals {
  controller_ips   = null_resource.controllers.*.triggers.ip
  controller_names = null_resource.controllers.*.triggers.name
  controller_cidrs = null_resource.controllers.*.triggers.cidr

  first_controller_ip = local.controller_ips[0]

  worker_ips   = null_resource.workers.*.triggers.ip
  worker_cidrs = null_resource.workers.*.triggers.cidr
  worker_names = null_resource.workers.*.triggers.name

  kubelet_extra_args = var.container_runtime == "containerd" ? concat(var.kubelet_extra_args, ["--container-runtime=remote", "--container-runtime-endpoint=unix:///run/containerd/containerd.sock"]) : var.kubelet_extra_args
  kubelet_extra_mounts = var.container_runtime == "containerd" ? [
    {
      source = "/run/containerd/",
      target = "/run/containerd",
    },
    {
      source = "/var/lib/containerd/",
      target = "/var/lib/containerd",
    },
  ] : []
}

resource "null_resource" "workers" {
  count = var.workers_count

  triggers = {
    name = format("worker%02d", count.index + 1)
    ip   = cidrhost(var.nodes_cidr, count.index + var.cidr_ips_offset + 10)
    cidr = cidrsubnet(var.pod_cidr, 8, count.index + 2 + 10)
  }
}
