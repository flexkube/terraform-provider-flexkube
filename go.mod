module github.com/flexkube/terraform-provider-flexkube

go 1.16

require (
	github.com/flexkube/libflexkube v0.7.0
	github.com/google/go-cmp v0.5.6
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.7.0
	sigs.k8s.io/yaml v1.2.0
)

replace (
	// Borrowed from Helm.
	github.com/docker/distribution => github.com/docker/distribution v0.0.0-20191216044856-a8371794149d

	// Use forked version of Helm to workaround https://github.com/helm/helm/issues/9761.
	helm.sh/helm/v3 => github.com/flexkube/helm/v3 v3.1.0-rc.1.0.20210728081922-539dfe1e558a
)
