module github.com/flexkube/terraform-provider-flexkube

go 1.15

require (
	github.com/flexkube/libflexkube v0.4.0
	github.com/google/go-cmp v0.5.2
	github.com/hashicorp/terraform-plugin-sdk v1.15.0
	go.etcd.io/etcd v3.3.25+incompatible // indirect
	sigs.k8s.io/yaml v1.2.0
)

replace (
	// sigs.k8s.io/kustomize@v2.0.3+incompatible pulled by
	// k8s.io/cli-runtime pulled by helm.sh/helm/v3
	// is not compatible with spec v0.19.9.
	github.com/go-openapi/spec => github.com/go-openapi/spec v0.19.8

	// Force updating docker/docker to most recent version.
	github.com/moby/moby => github.com/moby/moby v17.12.0-ce-rc1.0.20200618181300-9dc6525e6118+incompatible

	// With https://github.com/moby/term/pull/14 merged to fix windows builds.
	github.com/moby/term => github.com/flexkube/term v0.0.0-20200813115617-6a39c2cf564e

	// k8s.io/kubectl is not compatible with never version.
	github.com/russross/blackfriday => github.com/russross/blackfriday v1.5.2

	// Force updating etcd to most recent version.
	go.etcd.io/etcd => go.etcd.io/etcd v0.5.0-alpha.5.0.20200824191128-ae9734ed278b

	// To fix building github.com/moby/term. See
	// https://github.com/moby/term/issues/15 for more details.
	golang.org/x/sys => golang.org/x/sys v0.0.0-20200826173525-f9321e4c35a6

	// Most recent etcd version is not compatible with grpc v1.13.x.
	google.golang.org/grpc => google.golang.org/grpc v1.29.1

	// Use flexkube fork of Helm, which has K8s dependencies updated to v0.19.0.
	helm.sh/helm/v3 => github.com/flexkube/helm/v3 v3.3.1

	// Force updating client-go to most recent version.
	k8s.io/client-go => k8s.io/client-go v0.19.0
)
