package flexkube_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/flexkube/terraform-provider-flexkube/flexkube"
)

func TestPKIPropagateEtcdPeersToServers(t *testing.T) {
	config := `
resource "flexkube_pki" "pki" {
	etcd {
		peers = {
			"foo" = "1.1.1.1"
		}
	}
}
`

	updatedConfig := `
resource "flexkube_pki" "pki" {
	etcd {
		peers = {
			"foo" = "1.1.1.1"
			"bar" = "2.2.2.2"
		}
	}
}
`

	resource.UnitTest(t, resource.TestCase{
		Providers: map[string]*schema.Provider{
			"flexkube": flexkube.Provider(),
		},
		Steps: []resource.TestStep{
			{
				Config: config,
			},
			{
				Config:             config,
				ExpectNonEmptyPlan: false,
			},
			{
				Config: updatedConfig,
			},
		},
	})
}
