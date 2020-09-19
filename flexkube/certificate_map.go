package flexkube

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/flexkube/libflexkube/pkg/pki"
)

func certificateMapMarshal(cm map[string]*pki.Certificate, sensitive bool) interface{} {
	if len(cm) == 0 {
		return nil
	}

	r := []interface{}{}

	for k, v := range cm {
		r = append(r, map[string]interface{}{
			"common_name": k,
			"certificate": []interface{}{certificateMarshal(sensitive, v)},
		})
	}

	return r
}

func certificateMapUnmarshal(i interface{}) map[string]*pki.Certificate {
	if i == nil {
		return nil
	}

	j, ok := i.([]interface{})

	if !ok || len(j) == 0 {
		return nil
	}

	r := map[string]*pki.Certificate{}

	for _, v := range j {
		if v == nil {
			continue
		}

		vv, ok := v.(map[string]interface{})
		if !ok || len(vv) == 0 {
			continue
		}

		cn := vv["common_name"].(string)
		if cn == "" {
			continue
		}

		if c := certificateUnmarshal(vv["certificate"]); c != nil {
			r[cn] = c
		}
	}

	return r
}

func certificateMapSchema(computed bool) *schema.Schema {
	return optionalList(computed, func(computed bool) *schema.Resource {
		return &schema.Resource{
			Schema: map[string]*schema.Schema{
				"common_name": optionalString(computed),
				"certificate": certificateBlockSchema(computed),
			},
		}
	})
}
