package flexkube

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/flexkube/libflexkube/pkg/helm/release"
)

func resourceHelmRelease() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceHelmReleaseCreate,
		ReadContext:   resourceHelmReleaseRead,
		DeleteContext: resourceHelmReleaseDelete,
		UpdateContext: resourceHelmReleaseCreate,
		Schema: map[string]*schema.Schema{
			"kubeconfig": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"chart": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"values": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  ">0.0.0-0",
			},
			"create_namespace": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"wait": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func getRelease(d *schema.ResourceData, m interface{}) (release.Release, error) {
	r := release.Config{
		Kubeconfig: d.Get("kubeconfig").(string),
		Namespace:  d.Get("namespace").(string),
		Name:       d.Get("name").(string),
		Chart:      d.Get("chart").(string),
		Values:     d.Get("values").(string),
		Version:    d.Get("version").(string),
	}

	if v, ok := d.GetOk("wait"); ok {
		r.Wait = v.(bool)
	}

	if v, ok := d.GetOk("create_namespace"); ok {
		r.CreateNamespace = v.(bool)
	}

	l := m.(*meta)
	l.helmClientLock.Lock()
	defer l.helmClientLock.Unlock()

	return r.New()
}

func getReleaseID(d *schema.ResourceData) string {
	chart := d.Get("chart").(string)
	name := d.Get("name").(string)
	namespace := d.Get("namespace").(string)
	kubeconfig := d.Get("kubeconfig").(string)

	return sha256sum([]byte(chart + name + namespace + kubeconfig))
}

func resourceHelmReleaseCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	release, err := getRelease(d, m)
	if err != nil {
		return diagFromErr(err)
	}

	if err := release.InstallOrUpgrade(); err != nil {
		return diagFromErr(err)
	}

	d.SetId(getReleaseID(d))

	return nil
}

func resourceHelmReleaseRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	release, err := getRelease(d, m)
	if err != nil {
		return diagFromErr(err)
	}

	e, err := release.Exists()
	if err != nil {
		return diagFromErr(err)
	}

	if e {
		d.SetId(getReleaseID(d))
	} else {
		d.SetId("")
	}

	return nil
}

func resourceHelmReleaseDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	release, err := getRelease(d, m)
	if err != nil {
		return diagFromErr(err)
	}

	if err := release.Uninstall(); err != nil {
		return diagFromErr(err)
	}

	d.SetId("")

	return nil
}
