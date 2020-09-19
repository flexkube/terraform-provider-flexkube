module github.com/flexkube/terraform-provider-flexkube

go 1.15

require (
	cloud.google.com/go v0.66.0 // indirect
	cloud.google.com/go/storage v1.11.0 // indirect
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/asaskevich/govalidator v0.0.0-20200907205600-7a23bdc65eef // indirect
	github.com/aws/aws-sdk-go v1.34.24 // indirect
	github.com/containerd/containerd v1.4.1 // indirect
	github.com/emicklei/go-restful v2.14.2+incompatible // indirect
	github.com/flexkube/libflexkube v0.4.3
	github.com/go-logr/logr v0.2.1 // indirect
	github.com/google/go-cmp v0.5.2
	github.com/hashicorp/go-hclog v0.14.1 // indirect
	github.com/hashicorp/go-multierror v1.1.0 // indirect
	github.com/hashicorp/go-uuid v1.0.2 // indirect
	github.com/hashicorp/hcl/v2 v2.6.0 // indirect
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.0.3
	github.com/hashicorp/yamux v0.0.0-20200609203250-aecfd211c9ce // indirect
	github.com/mitchellh/go-testing-interface v1.14.1 // indirect
	github.com/mitchellh/mapstructure v1.3.3 // indirect
	github.com/moby/term v0.0.0-20200915141129-7f0af18e79f2 // indirect
	github.com/oklog/run v1.1.0 // indirect
	github.com/ulikunitz/xz v0.5.8 // indirect
	github.com/zclconf/go-cty v1.6.1 // indirect
	go.etcd.io/etcd v3.3.25+incompatible // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.16.0 // indirect
	golang.org/x/sys v0.0.0-20200916084744-dbad9cb7cb7a // indirect
	golang.org/x/tools v0.0.0-20200915201639-f4cefd1cb5ba // indirect
	google.golang.org/api v0.32.0 // indirect
	google.golang.org/genproto v0.0.0-20200915202801-9f80d0600517 // indirect
	google.golang.org/grpc v1.32.0 // indirect
	helm.sh/helm/v3 v3.3.1 // indirect
	k8s.io/apiextensions-apiserver v0.19.1 // indirect
	k8s.io/kube-openapi v0.0.0-20200831175022-64514a1d5d59 // indirect
	k8s.io/kubectl v0.19.1 // indirect
	k8s.io/kubelet v0.19.1 // indirect
	k8s.io/utils v0.0.0-20200912215256-4140de9c8800 // indirect
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
