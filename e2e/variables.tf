variable "ssh_private_key_path" {
  default = "/root/.ssh/id_rsa"
}

variable "controllers_count" {
  default = 1
}

variable "workers_count" {
  default = 0
}

variable "nodes_cidr" {
  default = "192.168.50.0/24"
}

variable "pod_cidr" {
  default = "10.1.0.0/16"
}

variable "node_ssh_port" {
  default = 22
}

variable "kube_apiserver_helm_chart_source" {
  default = "flexkube/kube-apiserver"
}

variable "kubernetes_helm_chart_source" {
  default = "flexkube/kubernetes"
}

variable "kube_proxy_helm_chart_source" {
  default = "flexkube/kube-proxy"
}

variable "tls_bootstrapping_helm_chart_source" {
  default = "flexkube/tls-bootstrapping"
}

variable "kubelet_rubber_stamp_helm_chart_source" {
  default = "flexkube/kubelet-rubber-stamp"
}

variable "calico_helm_chart_source" {
  default = "flexkube/calico"
}

variable "coredns_helm_chart_source" {
  default = "flexkube/coredns"
}

variable "metrics_server_helm_chart_source" {
  default = "flexkube/metrics-server"
}

variable "kube_apiserver_helm_chart_version" {
  default = "0.5.0"
}

variable "kubernetes_helm_chart_version" {
  default = "0.6.0"
}

variable "kube_proxy_helm_chart_version" {
  default = "0.5.0"
}

variable "tls_bootstrapping_helm_chart_version" {
  default = "0.1.1"
}

variable "coredns_helm_chart_version" {
  default = "2.0.6"
}

variable "metrics_server_helm_chart_version" {
  default = "3.0.9"
}

variable "kubelet_rubber_stamp_helm_chart_version" {
  default = "0.1.8"
}

variable "calico_helm_chart_version" {
  default = "0.4.22"
}

variable "flatcar_channel" {
  default = "stable"
}

variable "cidr_ips_offset" {
  default = 2
}

variable "kubelet_extra_args" {
  default = [
    "--container-runtime-endpoint=unix:///run/containerd/containerd.sock",
  ]
}

variable "cgroup_driver" {
  default = "systemd"
}
