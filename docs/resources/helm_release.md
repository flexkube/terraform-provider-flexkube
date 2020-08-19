# Helm Release Resource

This resource allows to create a Helm release on a Kubernetes cluster.

## Example Usage

```hcl
resource "flexkube_helm_release" "cert_manager" {
  kubeconfig = ""
  namespace  = "cert-manager"
  chart      = "jetstack/cert-manager"
  name       = "cert-manager"
  values     = <<EOF
global:
  podSecurityPolicy:
    enabled: true
    useAppArmor: false
  installCRDs: true
prometheus:
  servicemonitor:
    enabled: true
EOF
}
```

## Argument Reference

* `kubeconfig` - (Required) Content of `kubeconfig` file, which should be used to create a release.

* `namespace` - (Required) Namespace in which the release should be created.

* `chart` - (Required) Chart which should be used when creating the release. This can be either a string in `repository/chart` format or path to chart on a local filesystem.

* `name` - (Required) Name of the release to create.

* `values` - (Optional) Values to configure the chart in YAML format.

* `version` - (Optional) Version of the chart to install. Defaults to `>0.0.0-0`, which will install latest version of the chart.

* `create_namespace` - (Optional) If `true`, namespace for the release will be automatically craeted. Defaults to `false`.

## Attribute Reference

This resource supports no attributes.
