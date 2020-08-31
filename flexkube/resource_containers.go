package flexkube

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/flexkube/libflexkube/pkg/container/resource"
	"github.com/flexkube/libflexkube/pkg/types"
)

func resourceContainers() *schema.Resource {
	return &schema.Resource{
		Create: resourceCreate(containersUnmarshal),
		// Update should be exactly the same operation as Create.
		Update:        resourceCreate(containersUnmarshal),
		Read:          resourceRead(containersUnmarshal),
		Delete:        resourceDelete(containersUnmarshal, "host_configured_container"),
		CustomizeDiff: resourceDiff(containersUnmarshal),
		Schema: withCommonFields(map[string]*schema.Schema{
			// Configuration specified by the user.
			"host_configured_container": hostConfiguredContainerSchema(false, false),
		}),
	}
}

func containersUnmarshal(d getter, includeState bool) types.ResourceConfig {
	c := &resource.Containers{}

	if cs := containersStateUnmarshal(d.Get("host_configured_container")); cs != nil {
		c.Containers = *cs
	}

	if s := getState(d); includeState && s != nil {
		c.State = *s
	}

	return c
}
