// Package flexkube implements Terraform provider for libflexkube.
package flexkube

import (
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider returns Terraform Flexkube provider instance.
func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"flexkube_etcd_cluster":         resourceEtcdCluster(),
			"flexkube_kubelet_pool":         resourceKubeletPool(),
			"flexkube_controlplane":         resourceControlplane(),
			"flexkube_apiloadbalancer_pool": resourceAPILoadBalancerPool(),
			"flexkube_helm_release":         resourceHelmRelease(),
			"flexkube_containers":           resourceContainers(),
			"flexkube_pki":                  resourcePKI(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(*schema.ResourceData) (interface{}, error) {
	return &meta{
		helmClientLock: sync.Mutex{},
	}, nil
}

// Meta is the meta information structure for the provider.
type meta struct {
	// Mutex to create only one helm client as a time, as it is not thread-safe.
	helmClientLock sync.Mutex
}
