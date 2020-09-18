package flexkube

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"sigs.k8s.io/yaml"

	"github.com/flexkube/libflexkube/pkg/pki"
)

func resourcePKI() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePKICreate,
		ReadContext:   resourcePKIRead,
		DeleteContext: resourcePKIDelete,
		UpdateContext: resourcePKICreate,
		Schema: map[string]*schema.Schema{
			"certificate": certificateBlockSchema(false),
			"root_ca":     certificateBlockSchema(true),
			"etcd":        etcdSchema(true),
			"kubernetes":  kubernetesSchema(true),
			"state_yaml":  sensitiveString(true),
		},
	}
}

func getPKI(d *schema.ResourceData) *pki.PKI {
	return &pki.PKI{
		Certificate: *certificateUnmarshal(d.Get("certificate")),
		RootCA:      certificateUnmarshal(d.Get("root_ca")),
		Etcd:        etcdUnmarshal(d.Get("etcd")),
		Kubernetes:  kubernetesUnmarshal(d.Get("kubernetes")),
	}
}

func savePKI(d *schema.ResourceData, p *pki.PKI) error {
	b, err := yaml.Marshal(p)
	if err != nil {
		return fmt.Errorf("failed converting PKI to YAML: %w", err)
	}

	props := map[string]interface{}{
		"state_yaml": string(b),
		"root_ca":    []interface{}{certificateMarshal(p.RootCA)},
		"etcd":       etcdMarshal(p.Etcd),
		"kubernetes": kubernetesMarshal(p.Kubernetes),
	}

	for k, v := range props {
		if err := d.Set(k, v); err != nil {
			return fmt.Errorf("failed setting property %q: %w", k, err)
		}
	}

	return nil
}

func resourcePKICreate(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	p := getPKI(d)

	if err := p.Generate(); err != nil {
		return diagFromErr(err)
	}

	b, err := yaml.Marshal(p)
	if err != nil {
		return diagFromErr(err)
	}

	d.SetId(sha256sum(b))

	return diagFromErr(savePKI(d, p))
}

func resourcePKIDelete(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.SetId("")

	return nil
}

func resourcePKIRead(c context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}
