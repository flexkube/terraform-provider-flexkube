package flexkube

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/flexkube/libflexkube/pkg/pki"
)

func kubernetesMarshal(e *pki.Kubernetes, sensitive bool) interface{} {
	if e == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"certificate":                         []interface{}{certificateMarshal(sensitive, &e.Certificate)},
			"ca":                                  []interface{}{certificateMarshal(sensitive, e.CA)},
			"front_proxy_ca":                      []interface{}{certificateMarshal(sensitive, e.FrontProxyCA)},
			"admin_certificate":                   []interface{}{certificateMarshal(sensitive, e.AdminCertificate)},
			"kube_controller_manager_certificate": []interface{}{certificateMarshal(sensitive, e.KubeControllerManagerCertificate)}, //nolint:lll
			"kube_scheduler_certificate":          []interface{}{certificateMarshal(sensitive, e.KubeSchedulerCertificate)},
			"service_account_certificate":         []interface{}{certificateMarshal(sensitive, e.ServiceAccountCertificate)},
			"kube_api_server":                     pkiKubeAPIServerMarshal(e.KubeAPIServer, sensitive),
		},
	}
}

func kubernetesUnmarshal(i interface{}) *pki.Kubernetes {
	j, ok := i.([]interface{})
	if !ok || len(j) != 1 {
		return nil
	}

	l, ok := j[0].(map[string]interface{})

	if !ok {
		return &pki.Kubernetes{}
	}

	e := &pki.Kubernetes{
		CA:                               certificateUnmarshal(l["ca"]),
		FrontProxyCA:                     certificateUnmarshal(l["front_proxy_ca"]),
		AdminCertificate:                 certificateUnmarshal(l["admin_certificate"]),
		KubeControllerManagerCertificate: certificateUnmarshal(l["kube_controller_manager_certificate"]),
		KubeSchedulerCertificate:         certificateUnmarshal(l["kube_scheduler_certificate"]),
		ServiceAccountCertificate:        certificateUnmarshal(l["service_account_certificate"]),
		KubeAPIServer:                    pkiKubeAPIServerUnmarshal(l["kube_api_server"]),
	}

	if c := certificateUnmarshal(l["certificate"]); c != nil {
		e.Certificate = *c
	}

	return e
}

func kubernetesSchema(computed bool) *schema.Schema {
	return optionalBlock(computed, false, func(computed bool) map[string]*schema.Schema {
		return map[string]*schema.Schema{
			"certificate":                         certificateBlockSchema(computed),
			"ca":                                  certificateBlockSchema(computed),
			"front_proxy_ca":                      certificateBlockSchema(computed),
			"admin_certificate":                   certificateBlockSchema(computed),
			"kube_controller_manager_certificate": certificateBlockSchema(computed),
			"kube_scheduler_certificate":          certificateBlockSchema(computed),
			"service_account_certificate":         certificateBlockSchema(computed),
			"kube_api_server":                     pkiKubeAPIServerSchema(computed),
		}
	})
}
