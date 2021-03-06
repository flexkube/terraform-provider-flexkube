package flexkube

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/flexkube/libflexkube/pkg/apiloadbalancer"
	"github.com/flexkube/libflexkube/pkg/types"
)

func resourceAPILoadBalancerPool() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreate(apiLoadBalancersUnmarshal),
		ReadContext:   resourceRead(apiLoadBalancersUnmarshal),
		DeleteContext: resourceDelete(apiLoadBalancersUnmarshal, "api_load_balancer"),
		UpdateContext: resourceCreate(apiLoadBalancersUnmarshal),
		CustomizeDiff: resourceDiff(apiLoadBalancersUnmarshal),
		Schema: withCommonFields(map[string]*schema.Schema{
			"image":             optionalString(false),
			"ssh":               sshSchema(false),
			"servers":           optionalStringList(false),
			"api_load_balancer": apiLoadBalancerSchema(),
			"name":              optionalString(false),
			"host_config_path":  optionalString(false),
			"bind_address":      optionalString(false),
		}),
	}
}

func apiLoadBalancersUnmarshal(d getter, includeState bool) types.ResourceConfig {
	servers := []string{}

	if i, ok := d.GetOk("servers"); ok {
		s := i.([]interface{})

		for _, v := range s {
			servers = append(servers, v.(string))
		}
	}

	cc := &apiloadbalancer.APILoadBalancers{
		Image:            d.Get("image").(string),
		Servers:          servers,
		APILoadBalancers: apiLoadBalancerUnmarshal(d.Get("api_load_balancer")),
		Name:             d.Get("name").(string),
		HostConfigPath:   d.Get("host_config_path").(string),
		BindAddress:      d.Get("bind_address").(string),
	}

	if s := getState(d); includeState && s != nil {
		cc.State = *s
	}

	if d, ok := d.GetOk("ssh"); ok && len(d.([]interface{})) == 1 {
		cc.SSH = sshUnmarshal(d.([]interface{})[0])
	}

	return cc
}
