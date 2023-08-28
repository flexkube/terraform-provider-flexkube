package flexkube

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/flexkube/libflexkube/pkg/host/transport/direct"
)

func directMarshal(direct.Config) interface{} {
	return []interface{}{map[string]interface{}{}}
}

func directUnmarshal(interface{}) *direct.Config {
	return &direct.Config{}
}

func directSchema(computed bool) *schema.Schema {
	return optionalBlock(computed, false, func(computed bool) map[string]*schema.Schema {
		return map[string]*schema.Schema{}
	})
}
