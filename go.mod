module github.com/flexkube/terraform-provider-flexkube

go 1.15

require (
	cloud.google.com/go v0.66.0 // indirect
	cloud.google.com/go/storage v1.11.0 // indirect
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/aws/aws-sdk-go v1.34.24 // indirect
	github.com/flexkube/libflexkube v0.5.0
	github.com/google/go-cmp v0.5.4
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-hclog v0.15.0 // indirect
	github.com/hashicorp/go-multierror v1.1.0 // indirect
	github.com/hashicorp/go-plugin v1.4.0 // indirect
	github.com/hashicorp/go-uuid v1.0.2 // indirect
	github.com/hashicorp/hcl/v2 v2.8.2 // indirect
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.4.3
	github.com/hashicorp/yamux v0.0.0-20200609203250-aecfd211c9ce // indirect
	github.com/mitchellh/go-testing-interface v1.14.1 // indirect
	github.com/mitchellh/mapstructure v1.4.1 // indirect
	github.com/oklog/run v1.1.0 // indirect
	github.com/vmihailenco/msgpack v4.0.4+incompatible // indirect
	github.com/zclconf/go-cty v1.7.1 // indirect
	go.etcd.io/etcd v3.3.25+incompatible // indirect
	google.golang.org/api v0.32.0 // indirect
	sigs.k8s.io/yaml v1.2.0
)

replace (
	// Borrowed from Helm.
	github.com/docker/distribution => github.com/docker/distribution v0.0.0-20191216044856-a8371794149d

	// Force updating docker/docker to v19.03.15.
	github.com/docker/docker => github.com/moby/moby v17.12.0-ce-rc1.0.20200618181300-9dc6525e6118+incompatible

	// With v0.2.0 package has been renames, so until all dependencies are updated to use new import name,
	// we need to use older version.
	//
	// See: https://github.com/moby/spdystream/releases/tag/v0.2.0
	github.com/docker/spdystream => github.com/moby/spdystream v0.1.0

	// For testing.
	github.com/flexkube/libflexkube => github.com/invidian/libflexkube v0.1.1-0.20210218101112-095c8a8f774c

	// sigs.k8s.io/kustomize@v2.0.3+incompatible pulled by
	// k8s.io/cli-runtime pulled by helm.sh/helm/v3
	// is not compatible with spec v0.19.9.
	github.com/go-openapi/spec => github.com/go-openapi/spec v0.19.8

	// Newer version requires new version of GRCP, which is not compatible with used etcd version.
	github.com/hashicorp/terraform-plugin-sdk/v2 => github.com/hashicorp/terraform-plugin-sdk/v2 v2.1.0

	// With https://github.com/moby/term/pull/14 merged to fix windows builds.
	github.com/moby/term => github.com/flexkube/term v0.0.0-20200813115617-6a39c2cf564e

	// k8s.io/kubectl is not compatible with never version.
	github.com/russross/blackfriday => github.com/russross/blackfriday v1.5.2

	// Force updating etcd to most recent version.
	go.etcd.io/etcd => go.etcd.io/etcd v0.5.0-alpha.5.0.20200824191128-ae9734ed278b

	// Most recent etcd version is not compatible with grpc v1.13.x.
	google.golang.org/grpc => google.golang.org/grpc v1.29.1

	// Force updating client-go to most recent version.
	k8s.io/client-go => k8s.io/client-go v0.20.2
)
