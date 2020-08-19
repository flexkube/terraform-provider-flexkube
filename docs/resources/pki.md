# PKI Resource

This resource allows to create and manage all certificates required for creating Kubernetes
cluster by following [Kubernetes PKI certificates and requirements](https://kubernetes.io/docs/setup/best-practices/certificates/).

## Example Usage

```hcl
resource "flexkube_pki" "pki" {
  certificate {
    organization = "example"
  }

  etcd {
    peers   = {
			"foo" = "192.168.10.10"
		}
    servers = {
			"foo" = "192.168.10.10"
		}

    client_cns = [
      "root",
      "kube-apiserver",
      "prometheus",
    ]
  }

  kubernetes {
    kube_api_server {
      external_names = ["kube-apiserver.example.com"]
      server_ips     = ["192.168.10.10", "127.0.1.1", "11.0.0.1"]
    }
  }
}
```

## Argument Reference

* `certificate` - (Optional) A `certificate` block as defined below. Settings from this certificate will be populated to all certificates in the PKI unless more specific overrides are found. This certificate won't have `x509_certificate`, `public_key` and `private_key` attributes populated.

* `root_ca` - (Optional) A `certificate` block as defined below. This certificates represents PKI root certificate.

* `etcd` - (Optional) A `etcd` block as defined below. This block holds all etcd related certificates. If this block is set, etcd CA certificate will be generated.

* `kubernetes`- (Optional) A `kubernetes` block as defined below. This block holds all Kubernetes related certificates. If this block is set, Kubernetes CA certificate will be generated.

---

A `etcd` block supports the following:

* `certificate` - (Optional) A `certificate` block as defined below. Settings from this certificate will be populated to all etcd certificates unless more specific overrides are found. This certificate won't have `x509_certificate`, `public_key` and `private_key` attributes populated.

* `ca` - (Optional) A `certificate` block as defined below. This field contains settings for etcd CA certificate and the generated certificates.

* `client_cns` - (Optional) List of client common names to generate client certificates, which will be stored in `client_certificates` field.

* `peers` - (Optional) Map of name and IP address of cluster members for communicating between each other. This field will generate the certificates into `peer_certificates` field.

* `servers` - (Optional) Map of name and IP address of cluster members for external communication. This field will generate the certificates into `server_certificates` field.

* `peer_certificates` - (Optional) Map of `certificate` blocks as defined below. Stores peer certificate settings and generated certificates.

* `server_certificates` - (Optional) Map of `certificate` blocks as defined below. Stores server certificate settings and generated certificates.

* `client_certificates` - (Optional) Map of `certificate` block as defined below. Stores client certificate settings and generated certificates.

---

A `kubernetes` block supports the following:

* `certificate` - (Optional) A `certificate` block as defined below. Settings from this certificate will be populated to all etcd certificates unless more specific overrides are found. This certificate won't have `x509_certificate`, `public_key` and `private_key` attributes populated.

* `ca` - (Optional) A `certificate` block as defined below. This field contains settings for Kubernetes CA certificate and the generated certificates.

* `front_proxy_ca` - (Optional) A `certificate` block as defined below. This field contains settings for Kubernetes Front Proxy CA certificate and the generated certificates.

* `admin_certificate` - (Optional) A `certificate` block as defined below. This field contains settings for Kubernetes Admin/Root certificate and the generated certificates.

* `kube_controller_manager_certificate` - (Optional) A `certificate` block as defined below. This field contains settings for Kube Controller Manager certificate and the generated certificates.

* `kube_scheduler_certificate` - (Optional) A `certificate` block as defined below. This field contains settings for Kube Scheduler certificate and the generated certificates.

* `service_account_certificate` - (Optional) A `certificate` block as defined below. This field contains settings for certificates used for signing and validating service account tokens by kube-controller-manager and kube-apiserver.

* `kube_api_server` - (Optional) A `kube_api_server` block as defined below.

---

A `kube_api_server` block supports the following:

* `certificate` - (Optional) A `certificate` block as defined below. Settings from this certificate will be populated to all kube-apiserver certificates unless more specific overrides are found. This certificate won't have `x509_certificate`, `public_key` and `private_key` attributes populated.

* `external_names` - (Optional) List of DNS names where kube-apiserver will be served.

* `server_ips` - (Optional) List of IP addresses where kube-apiserver will be served.

* `server_certificate` - (Optional) A `certificate` block as defined below. This field contains settings for Kubernetes server certificate and the generated certificates.

* `kubelet_certificate` - (Optional) A `certificate` block as defined below. This field contains settings for kube-apiserver client certificate used for communicating with kubelets and the generated certificates.

* `front_proxy_client_certificate` - (Optional) A `certificate` block as defined below. This field contains settings for kube-apiserver client certificate used for communicating with extension API servers and the generated certificates.

---

A `certificate` block supports the following:

* `organization` - (Optional) Organization field in X.509 certificate. Defaults to `organization`.

* `rsa_bits` - (Optional) Number of bits to use for generating RSA key. Defaults to `2048`.

* `validity_duration` - (Optional) Duration for how long generated certificate should be valid, expressed in [Go Duration format](https://golang.org/pkg/time/#ParseDuration). Defaults to `8760h` (365 days).

* `renew_threshold` - (Optional) Duration for how early before end of certificate expiry the certificate should be renewed, expressed in [Go Duration format](https://golang.org/pkg/time/#ParseDuration). Defaults to `720h` (30 days). NOTE: This option currently has no effect, as certificate renewal is not supported.

* `common_name` - (Optional) Common Name field in X.509 certificate.

* `ca` - (Optional) If `true`, certificate will be marked as CA Certificate. Defaults to `false`.

* `key_usage` - (Optional) List of allowed key usages. It accepts the following keywords, combining the set of flags defined by
both [Key Usage](https://tools.ietf.org/html/rfc5280#section-4.2.1.3) and
[Extended Key Usage](https://tools.ietf.org/html/rfc5280#section-4.2.1.12) in
[RFC5280](https://tools.ietf.org/html/rfc5280):
	* `digital_signature`
	* `content_commitment`
	* `key_encipherment`
	* `data_encipherment`
	* `key_agreement`
	* `cert_signing`
	* `crl_signing`
	* `encipher_only`
	* `decipher_only`
	* `any_extended`
	* `server_auth`
	* `client_auth`
	* `code_signing`
	* `email_protection`
	* `ipsec_end_system`
	* `ipsec_tunnel`
	* `ipsec_user`
	* `timestamping`
	* `ocsp_signing`
	* `microsoft_server_gated_crypto`
	* `netscape_server_gated_crypto`

* `ip_addresses` - (Optional) List of allowed IP address for this certificate can be used for.

* `dns_names` - (Optional) List of allowed DNS names for this certificate to be used for.

* `x509_certificate` - (Optional) This field stores PEM encoded X.509 certificate. If the field is empty, it will store generated certificate after creation of the resource. If this field is not empty, defined certificate here might be used to sign other certificates.

* `public_key` - (Optional) This field stores PEM encoded RSA public key. If this field is empty, creating PKI resource will populate this field with generated certificate. If this field is not empty, defined certificate here might be used to sign other certificates.

* `private_key` - (Optional) This field stores PEM encoded RSA private key for X.509 certificate. If this field is empty, creating PKI resource will populate this field with generated privat key. If this field is not empty, defined certificate here might be used to sign other certificates.

## Attribute Reference

* `state_yaml` - PKI state in YAML format, which can be passed to `flexkube` CLI or to other resources like `flexkube_kubelet_pool`, which supports PKI integration.
