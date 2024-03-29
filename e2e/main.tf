resource "flexkube_pki" "pki" {
  certificate {
    organization = "example"
  }

  etcd {
    peers   = zipmap(local.controller_names, local.controller_ips)
    servers = zipmap(local.controller_names, local.controller_ips)

    client_cns = [
      "root",
      "kube-apiserver",
      "prometheus",
    ]
  }

  kubernetes {
    kube_api_server {
      external_names = ["kube-apiserver.example.com"]
      server_ips     = concat(local.controller_ips, ["127.0.1.1", "11.0.0.1"])
    }
  }
}

resource "random_password" "bootstrap_token_id" {
  length  = 6
  upper   = false
  special = false
}

resource "random_password" "bootstrap_token_secret" {
  length  = 16
  upper   = false
  special = false
}

locals {
  api_port = 8443

  node_load_balancer_address = "127.0.0.1:7443"

  etcd_servers = formatlist("https://%s:2379", values(flexkube_pki.pki.etcd[0].peers))

  tls_bootstrapping_values = <<EOF
tokens:
- token-id: ${random_password.bootstrap_token_id.result}
  token-secret: ${random_password.bootstrap_token_secret.result}
EOF

  client_certificates_map = { for index, v in flexkube_pki.pki.state_sensitive[0].etcd[0].client_certificates : v.common_name => v.certificate[0] }

  kube_apiserver_values = templatefile("./templates/kube-apiserver-values.yaml.tmpl", {
    server_key                     = flexkube_pki.pki.state_sensitive[0].kubernetes[0].kube_api_server[0].server_certificate[0].private_key
    server_certificate             = flexkube_pki.pki.state_sensitive[0].kubernetes[0].kube_api_server[0].server_certificate[0].x509_certificate
    service_account_private_key    = flexkube_pki.pki.state_sensitive[0].kubernetes[0].service_account_certificate[0].private_key
    ca_certificate                 = flexkube_pki.pki.state_sensitive[0].kubernetes[0].ca[0].x509_certificate
    front_proxy_client_key         = flexkube_pki.pki.state_sensitive[0].kubernetes[0].kube_api_server[0].front_proxy_client_certificate[0].private_key
    front_proxy_client_certificate = flexkube_pki.pki.state_sensitive[0].kubernetes[0].kube_api_server[0].front_proxy_client_certificate[0].x509_certificate
    front_proxy_ca_certificate     = flexkube_pki.pki.state_sensitive[0].kubernetes[0].front_proxy_ca[0].x509_certificate
    kubelet_client_certificate     = flexkube_pki.pki.state_sensitive[0].kubernetes[0].kube_api_server[0].kubelet_certificate[0].x509_certificate
    kubelet_client_key             = flexkube_pki.pki.state_sensitive[0].kubernetes[0].kube_api_server[0].kubelet_certificate[0].private_key
    etcd_ca_certificate            = flexkube_pki.pki.state_sensitive[0].etcd[0].ca[0].x509_certificate
    etcd_client_certificate        = local.client_certificates_map["kube-apiserver"].x509_certificate
    etcd_client_key                = local.client_certificates_map["kube-apiserver"].private_key
    etcd_servers                   = local.etcd_servers
    replicas                       = var.controllers_count
  })

  api_servers = formatlist("%s:%d", local.controller_ips, local.api_port)

  kubernetes_values = templatefile("./templates/values.yaml.tmpl", {
    service_account_private_key = flexkube_pki.pki.state_sensitive[0].kubernetes[0].service_account_certificate[0].private_key
    kubernetes_ca_key           = flexkube_pki.pki.state_sensitive[0].kubernetes[0].ca[0].private_key
    root_ca_certificate         = flexkube_pki.pki.state_sensitive[0].root_ca[0].x509_certificate
    kubernetes_ca_certificate   = flexkube_pki.pki.state_sensitive[0].kubernetes[0].ca[0].x509_certificate
    api_servers                 = local.api_servers
    replicas                    = var.controllers_count
  })

  kube_proxy_values = <<EOF
apiServers:
%{for api_server in local.api_servers~}
- ${api_server}
%{endfor~}
podCIDR: ${var.pod_cidr}
EOF

  coredns_values = <<EOF
rbac:
  pspEnable: true
service:
  clusterIP: 11.0.0.10
nodeSelector:
  node-role.kubernetes.io/master: ""
tolerations:
  - key: node-role.kubernetes.io/master
    operator: Exists
    effect: NoSchedule
EOF

  calico_values = <<EOF
podCIDR: ${var.pod_cidr}
flexVolumePluginDir: /var/lib/kubelet/volumeplugins
EOF

  metrics_server_values = <<EOF
rbac:
  pspEnabled: true
args:
- --kubelet-preferred-address-types=InternalIP
podDisruptionBudget:
  enabled: true
  minAvailable: 1
tolerations:
- key: node-role.kubernetes.io/master
  operator: Exists
  effect: NoSchedule
resources:
  requests:
    memory: 20Mi
EOF

  kubeconfig_admin = templatefile("./templates/kubeconfig.tmpl", {
    name        = "admin"
    server      = "https://${local.first_controller_ip}:${local.api_port}"
    ca_cert     = base64encode(flexkube_pki.pki.state_sensitive[0].kubernetes[0].ca[0].x509_certificate)
    client_cert = base64encode(flexkube_pki.pki.state_sensitive[0].kubernetes[0].admin_certificate[0].x509_certificate)
    client_key  = base64encode(flexkube_pki.pki.state_sensitive[0].kubernetes[0].admin_certificate[0].private_key)
  })

  deploy_workers = var.workers_count > 0 ? 1 : 0

  ssh_private_key = file(var.ssh_private_key_path)
}

resource "local_file" "kubeconfig" {
  sensitive_content = local.kubeconfig_admin
  filename          = "./kubeconfig"
}

resource "flexkube_etcd_cluster" "etcd" {
  ssh {
    user        = "core"
    port        = var.node_ssh_port
    private_key = local.ssh_private_key
  }

  pki_yaml = flexkube_pki.pki.state_yaml

  dynamic "member" {
    for_each = flexkube_pki.pki.etcd[0].peers

    content {
      name           = member.key
      peer_address   = member.value
      server_address = member.value

      host {
        ssh {
          address = member.value
        }
      }
    }
  }
}

resource "flexkube_apiloadbalancer_pool" "controllers" {
  name             = "api-loadbalancer-controllers"
  host_config_path = "/etc/haproxy/controllers.cfg"
  bind_address     = local.node_load_balancer_address
  servers          = formatlist("%s:%d", local.controller_ips, local.api_port)

  ssh {
    private_key = local.ssh_private_key
    port        = var.node_ssh_port
    user        = "core"
  }

  dynamic "api_load_balancer" {
    for_each = local.controller_ips

    content {
      host {
        ssh {
          address = api_load_balancer.value
        }
      }
    }
  }
}

resource "flexkube_controlplane" "bootstrap" {
  pki_yaml = flexkube_pki.pki.state_yaml

  kube_apiserver {
    service_cidr      = "11.0.0.0/24"
    etcd_servers      = local.etcd_servers
    bind_address      = local.first_controller_ip
    advertise_address = local.first_controller_ip
    secure_port       = local.api_port
  }

  kube_controller_manager {
    flex_volume_plugin_dir = "/var/lib/kubelet/volumeplugins"
  }

  api_server_address = local.first_controller_ip
  api_server_port    = local.api_port

  ssh {
    user        = "core"
    address     = local.first_controller_ip
    port        = var.node_ssh_port
    private_key = local.ssh_private_key
  }

  depends_on = [
    flexkube_etcd_cluster.etcd,
  ]
}

resource "flexkube_helm_release" "kube-apiserver" {
  kubeconfig = local.kubeconfig_admin
  namespace  = "kube-system"
  chart      = var.kube_apiserver_helm_chart_source
  version    = var.kube_apiserver_helm_chart_version
  name       = "kube-apiserver"
  values     = local.kube_apiserver_values
  wait       = true

  depends_on = [
    flexkube_controlplane.bootstrap,
    flexkube_helm_release.tls-bootstrapping,
  ]
}

resource "flexkube_helm_release" "kubernetes" {
  kubeconfig = local.kubeconfig_admin
  namespace  = "kube-system"
  chart      = var.kubernetes_helm_chart_source
  version    = var.kubernetes_helm_chart_version
  name       = "kubernetes"
  values     = local.kubernetes_values
  wait       = true

  depends_on = [
    flexkube_controlplane.bootstrap,
    flexkube_helm_release.calico,
  ]
}

resource "flexkube_helm_release" "kube-proxy" {
  kubeconfig = local.kubeconfig_admin
  namespace  = "kube-system"
  chart      = var.kube_proxy_helm_chart_source
  version    = var.kube_proxy_helm_chart_version
  name       = "kube-proxy"
  values     = local.kube_proxy_values
  wait       = true

  depends_on = [
    flexkube_controlplane.bootstrap,
    flexkube_helm_release.tls-bootstrapping,
  ]
}

resource "flexkube_helm_release" "tls-bootstrapping" {
  kubeconfig = local.kubeconfig_admin
  namespace  = "kube-system"
  chart      = var.tls_bootstrapping_helm_chart_source
  version    = var.tls_bootstrapping_helm_chart_version
  name       = "tls-bootstrapping"
  values     = local.tls_bootstrapping_values
  wait       = true

  depends_on = [
    flexkube_controlplane.bootstrap,
  ]
}

resource "flexkube_helm_release" "coredns" {
  kubeconfig = local.kubeconfig_admin
  namespace  = "kube-system"
  chart      = var.coredns_helm_chart_source
  version    = var.coredns_helm_chart_version
  name       = "coredns"
  values     = local.coredns_values
  wait       = true

  depends_on = [
    flexkube_controlplane.bootstrap,
    flexkube_helm_release.calico,
  ]
}

resource "flexkube_helm_release" "metrics-server" {
  kubeconfig = local.kubeconfig_admin
  namespace  = "kube-system"
  chart      = var.metrics_server_helm_chart_source
  version    = var.metrics_server_helm_chart_version
  name       = "metrics-server"
  values     = local.metrics_server_values
  wait       = true

  depends_on = [
    flexkube_controlplane.bootstrap,
    flexkube_helm_release.calico,
  ]
}

resource "flexkube_helm_release" "kubelet-rubber-stamp" {
  kubeconfig = local.kubeconfig_admin
  namespace  = "kube-system"
  chart      = var.kubelet_rubber_stamp_helm_chart_source
  version    = var.kubelet_rubber_stamp_helm_chart_version
  name       = "kubelet-rubber-stamp"
  wait       = true

  depends_on = [
    flexkube_controlplane.bootstrap,
    flexkube_helm_release.calico,
  ]
}

resource "flexkube_helm_release" "calico" {
  kubeconfig = local.kubeconfig_admin
  namespace  = "kube-system"
  chart      = var.calico_helm_chart_source
  version    = var.calico_helm_chart_version
  name       = "calico"
  values     = local.calico_values
  wait       = true

  depends_on = [
    flexkube_controlplane.bootstrap,
    flexkube_helm_release.kube-proxy,
  ]
}

resource "flexkube_kubelet_pool" "controller" {
  bootstrap_config {
    server = local.node_load_balancer_address
    token  = "${random_password.bootstrap_token_id.result}.${random_password.bootstrap_token_secret.result}"
  }

  wait_for_node_ready = true

  pki_yaml          = flexkube_pki.pki.state_yaml
  cgroup_driver     = local.cgroup_driver
  hairpin_mode      = "hairpin-veth"
  volume_plugin_dir = "/var/lib/kubelet/volumeplugins"
  cluster_dns_ips = [
    "11.0.0.10"
  ]

  system_reserved = {
    "cpu"    = "100m"
    "memory" = "500Mi"
  }

  kube_reserved = {
    // 100MB for kubelet and 200MB for etcd.
    "memory" = "300Mi"
    "cpu"    = "100m"
  }

  privileged_labels = {
    "node-role.kubernetes.io/master" = ""
  }

  admin_config {
    server = "${local.first_controller_ip}:${local.api_port}"
  }

  taints = var.workers_count > 0 ? {
    "node-role.kubernetes.io/master" = "NoSchedule"
  } : {}

  dynamic "extra_mount" {
    for_each = local.kubelet_extra_mounts

    content {
      source = extra_mount.value.source
      target = extra_mount.value.target
    }
  }

  extra_args = local.kubelet_extra_args

  ssh {
    user        = "core"
    port        = var.node_ssh_port
    private_key = local.ssh_private_key
  }

  dynamic "kubelet" {
    for_each = local.controller_ips

    content {
      name    = local.controller_names[kubelet.key]
      address = local.controller_ips[kubelet.key]

      host {
        ssh {
          address = kubelet.value
        }
      }
    }
  }

  depends_on = [
    flexkube_apiloadbalancer_pool.controllers,
    flexkube_helm_release.tls-bootstrapping,
  ]
}

resource "flexkube_apiloadbalancer_pool" "workers" {
  count = local.deploy_workers

  name             = "api-loadbalancer-workers"
  host_config_path = "/etc/haproxy/workers.cfg"
  bind_address     = local.node_load_balancer_address
  servers          = formatlist("%s:%d", local.controller_ips, local.api_port)

  ssh {
    private_key = local.ssh_private_key
    port        = var.node_ssh_port
    user        = "core"
  }

  dynamic "api_load_balancer" {
    for_each = local.worker_ips

    content {
      host {
        ssh {
          address = api_load_balancer.value
        }
      }
    }
  }
}

resource "flexkube_kubelet_pool" "workers" {
  count = local.deploy_workers

  bootstrap_config {
    server = local.node_load_balancer_address
    token  = "${random_password.bootstrap_token_id.result}.${random_password.bootstrap_token_secret.result}"
  }

  admin_config {
    server = "${local.first_controller_ip}:${local.api_port}"
  }

  wait_for_node_ready = true

  pki_yaml = flexkube_pki.pki.state_yaml

  cgroup_driver     = local.cgroup_driver
  hairpin_mode      = "hairpin-veth"
  volume_plugin_dir = "/var/lib/kubelet/volumeplugins"
  cluster_dns_ips = [
    "11.0.0.10"
  ]

  system_reserved = {
    "cpu"    = "100m"
    "memory" = "500Mi"
  }

  kube_reserved = {
    "memory" = "100Mi"
    "cpu"    = "100m"
  }

  dynamic "extra_mount" {
    for_each = local.kubelet_extra_mounts

    content {
      source = extra_mount.value.source
      target = extra_mount.value.target
    }
  }

  extra_args = local.kubelet_extra_args

  ssh {
    user        = "core"
    port        = var.node_ssh_port
    private_key = local.ssh_private_key
  }

  dynamic "kubelet" {
    for_each = local.worker_ips

    content {
      name    = local.worker_names[kubelet.key]
      address = local.worker_ips[kubelet.key]

      host {
        ssh {
          address = kubelet.value
        }
      }
    }
  }

  depends_on = [
    flexkube_apiloadbalancer_pool.workers,
    flexkube_helm_release.tls-bootstrapping,
  ]
}
