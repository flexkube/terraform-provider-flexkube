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

variable "network_plugin" {
  default = "calico"
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
  default = "0.3.4"
}

variable "kubernetes_helm_chart_version" {
  default = "0.4.4"
}

variable "kube_proxy_helm_chart_version" {
  default = "0.3.4"
}

variable "tls_bootstrapping_helm_chart_version" {
  default = "0.1.1"
}

variable "coredns_helm_chart_version" {
  default = "2.0.0"
}

variable "metrics_server_helm_chart_version" {
  default = "3.0.2"
}

variable "kubelet_rubber_stamp_helm_chart_version" {
  default = "0.1.6"
}

variable "calico_helm_chart_version" {
  default = "0.4.4"
}

variable "flatcar_channel" {
  default = "stable"
}
