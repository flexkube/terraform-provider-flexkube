package flexkube

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/flexkube/libflexkube/pkg/kubernetes/client"
	"github.com/flexkube/libflexkube/pkg/types"
)

func clientUnmarshal(i interface{}) *client.Config {
	// If block is not defined, don't return anything.
	if i == nil {
		return nil
	}

	j, ok := i.(map[string]interface{})
	if !ok || len(j) == 0 {
		return nil
	}

	c := &client.Config{}

	if h, ok := j["server"]; ok {
		c.Server = h.(string)
	}

	if h, ok := j["ca_certificate"]; ok {
		c.CACertificate = types.Certificate(h.(string))
	}

	if h, ok := j["client_certificate"]; ok {
		c.ClientCertificate = types.Certificate(h.(string))
	}

	if h, ok := j["client_key"]; ok {
		c.ClientKey = types.PrivateKey(h.(string))
	}

	if h, ok := j["token"]; ok {
		c.Token = h.(string)
	}

	return c
}

//nolint:unparam // False positive.
func clientSchema(computed bool) *schema.Schema {
	return optionalBlock(computed, false, func(computed bool) map[string]*schema.Schema {
		return map[string]*schema.Schema{
			"server":             optionalString(computed),
			"ca_certificate":     optionalString(computed),
			"client_certificate": optionalString(computed),
			"client_key":         sensitiveString(computed),
			"token":              sensitiveString(computed),
		}
	})
}
