package flexkube

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/flexkube/libflexkube/pkg/pki"
)

func pkiKubeAPIServerMarshal(e *pki.KubeAPIServer, sensitive bool) interface{} {
	return []interface{}{
		map[string]interface{}{
			"certificate":                    []interface{}{certificateMarshal(sensitive, &e.Certificate)},
			"external_names":                 stringSliceToInterfaceSlice(e.ExternalNames),
			"server_ips":                     stringSliceToInterfaceSlice(e.ServerIPs),
			"server_certificate":             []interface{}{certificateMarshal(sensitive, e.ServerCertificate)},
			"kubelet_certificate":            []interface{}{certificateMarshal(sensitive, e.KubeletCertificate)},
			"front_proxy_client_certificate": []interface{}{certificateMarshal(sensitive, e.FrontProxyClientCertificate)},
		},
	}
}

func pkiKubeAPIServerUnmarshal(i interface{}) *pki.KubeAPIServer {
	a := &pki.KubeAPIServer{}

	if i == nil {
		return a
	}

	j, ok := i.([]interface{})
	if !ok || len(j) != 1 {
		return a
	}

	k, ok := j[0].(map[string]interface{})

	if !ok || len(j) == 0 {
		return a
	}

	ka := &pki.KubeAPIServer{
		ExternalNames:               stringListUnmarshal(k["external_names"]),
		ServerIPs:                   stringListUnmarshal(k["server_ips"]),
		ServerCertificate:           certificateUnmarshal(k["server_certificate"]),
		KubeletCertificate:          certificateUnmarshal(k["kubelet_certificate"]),
		FrontProxyClientCertificate: certificateUnmarshal(k["front_proxy_client_certificate"]),
	}

	if c := certificateUnmarshal(k["certificate"]); c != nil {
		ka.Certificate = *c
	}

	return ka
}

func pkiKubeAPIServerSchema(computed bool) *schema.Schema {
	return optionalBlock(computed, false, func(computed bool) map[string]*schema.Schema {
		return map[string]*schema.Schema{
			"certificate":                    certificateBlockSchema(computed),
			"external_names":                 optionalStringList(computed),
			"server_ips":                     optionalStringList(computed),
			"server_certificate":             certificateBlockSchema(computed),
			"kubelet_certificate":            certificateBlockSchema(computed),
			"front_proxy_client_certificate": certificateBlockSchema(computed),
		}
	})
}
