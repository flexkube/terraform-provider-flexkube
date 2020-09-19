package flexkube

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/flexkube/libflexkube/pkg/pki"
)

func pkiMarshal(e *pki.PKI, sensitive bool) interface{} {
	if e == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"certificate": []interface{}{certificateMarshal(sensitive, &e.Certificate)},
			"root_ca":     []interface{}{certificateMarshal(sensitive, e.RootCA)},
			"etcd":        etcdMarshal(e.Etcd, sensitive),
			"kubernetes":  kubernetesMarshal(e.Kubernetes, sensitive),
		},
	}
}

func pkiSchema(computed bool, sensitive bool) *schema.Schema {
	return optionalBlock(computed, sensitive, func(computed bool) map[string]*schema.Schema {
		return map[string]*schema.Schema{
			"certificate": certificateBlockSchema(true),
			"root_ca":     certificateBlockSchema(true),
			"etcd":        etcdSchema(true),
			"kubernetes":  kubernetesSchema(true),
		}
	})
}
