package flexkube

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestContainersPlanOnly(t *testing.T) {
	t.Parallel()

	config := `
resource "flexkube_containers" "foo" {
  host_configured_container {
    name = "bar"

    container {
      config {
        name  = "bazhh"
        image = "nginx"

        env = {
          FOO = "bar"
        }
      }
    }
  }
}
`

	resource.UnitTest(t, resource.TestCase{
		Providers: map[string]*schema.Provider{
			"flexkube": Provider(),
		},
		Steps: []resource.TestStep{
			{
				Config:             config,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestContainersCreateRuntimeError(t *testing.T) {
	t.Parallel()

	config := `
resource "flexkube_containers" "foo" {
  host_configured_container {
    name = "bar"

    container {
      runtime {
        docker {
          host = "unix:///nonexistent"
        }
      }

      config {
        name  = "bazhh"
        image = "nginx"
      }
    }
  }
}
`

	resource.UnitTest(t, resource.TestCase{
		Providers: map[string]*schema.Provider{
			"flexkube": Provider(),
		},
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(`Cannot connect to the Docker daemon`),
			},
		},
	})
}

func TestContainersValidateFail(t *testing.T) {
	t.Parallel()

	config := `
resource "flexkube_containers" "foo" {
  host_configured_container {
    name = "bar"

    container {
      config {
        name = ""
        image = "nginx"
      }
    }
  }
}
`

	resource.UnitTest(t, resource.TestCase{
		Providers: map[string]*schema.Provider{
			"flexkube": Provider(),
		},
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(`name must be set`),
			},
		},
	})
}
