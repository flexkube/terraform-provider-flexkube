package flexkube

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/flexkube/libflexkube/pkg/kubelet"
	"github.com/flexkube/libflexkube/pkg/types"
)

func resourceKubeletPool() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreate(kubeletPoolUnmarshal),
		ReadContext:   resourceRead(kubeletPoolUnmarshal),
		DeleteContext: resourceDelete(kubeletPoolUnmarshal, "kubelet"),
		UpdateContext: resourceCreate(kubeletPoolUnmarshal),
		CustomizeDiff: resourceDiff(kubeletPoolUnmarshal),
		Schema: withCommonFields(map[string]*schema.Schema{
			"image":                     optionalString(false),
			"ssh":                       sshSchema(false),
			"bootstrap_config":          clientSchema(false),
			"kubernetes_ca_certificate": optionalString(false),
			"cluster_dns_ips":           optionalStringList(false),
			"taints": optionalMapPrimitive(false, func(computed bool) *schema.Schema {
				return &schema.Schema{
					Type: schema.TypeString,
				}
			}),
			"labels": optionalMapPrimitive(false, func(computed bool) *schema.Schema {
				return &schema.Schema{
					Type: schema.TypeString,
				}
			}),
			"privileged_labels": optionalMapPrimitive(false, func(computed bool) *schema.Schema {
				return &schema.Schema{
					Type: schema.TypeString,
				}
			}),
			"admin_config":  clientSchema(false),
			"cgroup_driver": optionalString(false),
			"system_reserved": optionalMapPrimitive(false, func(computed bool) *schema.Schema {
				return &schema.Schema{
					Type: schema.TypeString,
				}
			}),
			"kube_reserved": optionalMapPrimitive(false, func(computed bool) *schema.Schema {
				return &schema.Schema{
					Type: schema.TypeString,
				}
			}),
			"hairpin_mode":        optionalString(false),
			"volume_plugin_dir":   optionalString(false),
			"kubelet":             kubeletSchema(),
			"extra_mount":         mountsSchema(false),
			"pki_yaml":            sensitiveString(false),
			"wait_for_node_ready": optionalBool(false),
			"extra_args":          optionalStringList(false),
		}),
	}
}

func kubeletPoolUnmarshal(d getter, includeState bool) types.ResourceConfig {
	c := &kubelet.Pool{
		Image:                   d.Get("image").(string),
		KubernetesCACertificate: types.Certificate(d.Get("kubernetes_ca_certificate").(string)),
		CgroupDriver:            d.Get("cgroup_driver").(string),
		HairpinMode:             d.Get("hairpin_mode").(string),
		VolumePluginDir:         d.Get("volume_plugin_dir").(string),
		Kubelets:                kubeletsUnmarshal(d.Get("kubelet")),
		ClusterDNSIPs:           stringListUnmarshal(d.Get("cluster_dns_ips")),
		Taints:                  stringMapUnmarshal(d.Get("taints")),
		Labels:                  stringMapUnmarshal(d.Get("labels")),
		PrivilegedLabels:        stringMapUnmarshal(d.Get("privileged_labels")),
		SystemReserved:          stringMapUnmarshal(d.Get("system_reserved")),
		KubeReserved:            stringMapUnmarshal(d.Get("kube_reserved")),
		ExtraMounts:             mountsUnmarshal(d.Get("extra_mount")),
		PKI:                     unmarshalPKI(d),
		WaitForNodeReady:        d.Get("wait_for_node_ready").(bool),
		ExtraArgs:               stringListUnmarshal(d.Get("extra_args")),
	}

	if s := getState(d); includeState && s != nil {
		c.State = *s
	}

	if d, ok := d.GetOk("ssh"); ok && len(d.([]interface{})) == 1 {
		c.SSH = sshUnmarshal(d.([]interface{})[0])
	}

	if v, ok := d.GetOk("bootstrap_config"); ok && len(v.([]interface{})) == 1 {
		c.BootstrapConfig = clientUnmarshal(v.([]interface{})[0])
	}

	if v, ok := d.GetOk("admin_config"); ok && len(v.([]interface{})) == 1 {
		c.AdminConfig = clientUnmarshal(v.([]interface{})[0])
	}

	return c
}
