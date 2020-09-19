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
		peer_certificates {}

		peer_certificates {
			common_name = ""
		}

		peer_certificates {
			certificate {}
		}
	}
}
`

	updatedConfig := `
resource "flexkube_pki" "pki" {
	etcd {
		certificate {
			organization = "bar"
		}

		peers = {
			"foo" = "1.1.1.1"
			"bar" = "2.2.2.2"
		}
	}

	kubernetes {
		certificate {
			organization = "foo"
		}

		kube_api_server {
			certificate {}
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
