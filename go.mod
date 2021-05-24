module github.com/flexkube/terraform-provider-flexkube

go 1.16

require (
	cloud.google.com/go v0.66.0 // indirect
	cloud.google.com/go/storage v1.11.0 // indirect
	github.com/MakeNowJust/heredoc v1.0.0 // indirect
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/asaskevich/govalidator v0.0.0-20200907205600-7a23bdc65eef // indirect
	github.com/aws/aws-sdk-go v1.34.24 // indirect
	github.com/containerd/cgroups v0.0.0-20210114181951-8a68de567b68 // indirect
	github.com/containerd/continuity v0.0.0-20210208174643-50096c924a4e // indirect
	github.com/docker/go-metrics v0.0.1 // indirect
	github.com/emicklei/go-restful v2.15.0+incompatible // indirect
	github.com/exponent-io/jsonpath v0.0.0-20201116121440-e84ac1befdf8 // indirect
	github.com/fatih/color v1.10.0 // indirect
	github.com/flexkube/libflexkube v0.6.0
	github.com/google/go-cmp v0.5.5
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/googleapis/gnostic v0.5.4 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.0 // indirect
	github.com/hashicorp/go-uuid v1.0.2 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/hashicorp/hcl/v2 v2.8.2 // indirect
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.6.1
	github.com/hashicorp/yamux v0.0.0-20200609203250-aecfd211c9ce // indirect
	github.com/jmoiron/sqlx v1.3.1 // indirect
	github.com/magefile/mage v1.11.0 // indirect
	github.com/mattn/go-runewidth v0.0.10 // indirect
	github.com/mitchellh/go-testing-interface v1.14.1 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/mitchellh/mapstructure v1.4.1 // indirect
	github.com/oklog/run v1.1.0 // indirect
	github.com/prometheus/client_golang v1.9.0 // indirect
	github.com/prometheus/common v0.16.0 // indirect
	github.com/prometheus/procfs v0.6.0 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/rubenv/sql-migrate v0.0.0-20210215143335-f84234893558 // indirect
	github.com/russross/blackfriday v2.0.0+incompatible // indirect
	github.com/sirupsen/logrus v1.8.0 // indirect
	github.com/spf13/cobra v1.1.3 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	go.opencensus.io v0.22.6 // indirect
	golang.org/x/oauth2 v0.0.0-20210216194517-16ff1888fd2e // indirect
	google.golang.org/api v0.32.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	k8s.io/apimachinery v0.21.1 // indirect
	k8s.io/utils v0.0.0-20210111153108-fddb29f9d009 // indirect
	sigs.k8s.io/yaml v1.2.0
)

replace (
	// Borrowed from Helm.
	github.com/docker/distribution => github.com/docker/distribution v0.0.0-20191216044856-a8371794149d

	// Force updating docker/docker to v19.03.15.
	github.com/docker/docker => github.com/moby/moby v17.12.0-ce-rc1.0.20200618181300-9dc6525e6118+incompatible

	// sigs.k8s.io/kustomize@v2.0.3+incompatible pulled by
	// k8s.io/cli-runtime pulled by helm.sh/helm/v3
	// is not compatible with spec v0.19.9.
	github.com/go-openapi/spec => github.com/go-openapi/spec v0.19.8

	// k8s.io/kubectl is not compatible with never version.
	github.com/russross/blackfriday => github.com/russross/blackfriday v1.5.2
)
