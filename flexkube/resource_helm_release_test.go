package flexkube_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/flexkube/terraform-provider-flexkube/flexkube"
)

const helmReleasePlanOnly = `
resource "flexkube_helm_release" "metrics-server" {
  kubeconfig = <<EOF
apiVersion: v1
kind: Config
clusters:
- name: admin-cluster
  cluster: {}
users:
- name: admin-user
  user:
current-context: admin-context
contexts:
- name: admin-context
  context:
    cluster: admin-cluster
    namespace: kube-system
    user: admin-use
EOF
  namespace  = "kube-system"
  chart      = "foo/bar"
  version    = "1.2.3"
  name       = "metrics-server"
  wait       = true
  values     = <<EOF
foo: bar
EOF

	create_namespace = true
}

`

func TestHelmRelease(t *testing.T) {
	t.Parallel()

	resource.UnitTest(t, resource.TestCase{
		Providers: map[string]*schema.Provider{
			"flexkube": flexkube.Provider(),
		},
		Steps: []resource.TestStep{
			{
				Config:             helmReleasePlanOnly,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config:      helmReleasePlanOnly,
				ExpectError: regexp.MustCompile(`failed to create kubernetes client`),
			},
		},
	})
}

const helmReleaseBadValues = `
resource "flexkube_helm_release" "metrics-server" {
  kubeconfig = <<EOF
apiVersion: v1
kind: Config
clusters:
- name: admin-cluster
  cluster:
		server: https://foo:6443
		certificate-authority-data: Zm9vCg==
users:
- name: admin-user
  user:
		token: bar
current-context: admin-context
contexts:
- name: admin-context
  context:
    cluster: admin-cluster
    namespace: kube-system
    user: admin-user
EOF
  namespace  = "kube-system"
  chart      = "foo/bar"
  version    = "1.2.3"
  name       = "metrics-server"
  wait       = true
  values     = <<EOF
	foo: bar
EOF

  create_namespace = true
}

`

func TestHelmReleaseBadValues(t *testing.T) {
	t.Parallel()

	resource.UnitTest(t, resource.TestCase{
		Providers: map[string]*schema.Provider{
			"flexkube": flexkube.Provider(),
		},
		Steps: []resource.TestStep{
			{
				Config:      helmReleaseBadValues,
				ExpectError: regexp.MustCompile(`parsing field as YAM`),
			},
		},
	})
}

const helmReleaseBadKubeconfig = `
resource "flexkube_helm_release" "metrics-server" {
  kubeconfig = <<EOF
apiVersion: v1
kind: Config
clusters:
- name: admin-cluster
	cluster:
		server: https://foo:6443
		certificate-authority-data: Zm9vCg==
users:
- name: admin-user
  user:
		token: bar
current-context: admin-context
contexts:
- name: admin-context
  context:
    cluster: admin-cluster
    namespace: kube-system
    user: admin-user
EOF
  namespace  = "kube-system"
  chart      = "foo/bar"
  version    = "1.2.3"
  name       = "metrics-server"
  wait       = true

  create_namespace = true
}

`

func TestHelmReleaseBadKubeconfig(t *testing.T) {
	t.Parallel()

	resource.UnitTest(t, resource.TestCase{
		Providers: map[string]*schema.Provider{
			"flexkube": flexkube.Provider(),
		},
		Steps: []resource.TestStep{
			{
				Config:      helmReleaseBadKubeconfig,
				ExpectError: regexp.MustCompile(`parsing field as YAM`),
			},
		},
	})
}
