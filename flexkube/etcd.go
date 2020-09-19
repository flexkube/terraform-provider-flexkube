package flexkube

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/flexkube/libflexkube/pkg/pki"
)

func etcdMarshal(e *pki.Etcd, sensitive bool) interface{} {
	if e == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"certificate":         []interface{}{certificateMarshal(sensitive, &e.Certificate)},
			"ca":                  []interface{}{certificateMarshal(sensitive, e.CA)},
			"client_cns":          stringSliceToInterfaceSlice(e.ClientCNs),
			"peers":               stringMapMarshal(e.Peers, false),
			"servers":             stringMapMarshal(e.Servers, false),
			"peer_certificates":   certificateMapMarshal(e.PeerCertificates, sensitive),
			"server_certificates": certificateMapMarshal(e.ServerCertificates, sensitive),
			"client_certificates": certificateMapMarshal(e.ClientCertificates, sensitive),
		},
	}
}

func etcdUnmarshal(i interface{}) *pki.Etcd {
	j, ok := i.([]interface{})
	if !ok || len(j) != 1 {
		return nil
	}

	k, ok := j[0].(map[string]interface{})

	if !ok {
		return &pki.Etcd{}
	}

	e := &pki.Etcd{
		CA:                 certificateUnmarshal(k["ca"]),
		ClientCNs:          stringListUnmarshal(k["client_cns"]),
		Peers:              stringMapUnmarshal(k["peers"]),
		Servers:            stringMapUnmarshal(k["servers"]),
		PeerCertificates:   certificateMapUnmarshal(k["peer_certificates"]),
		ServerCertificates: certificateMapUnmarshal(k["server_certificates"]),
		ClientCertificates: certificateMapUnmarshal(k["client_certificates"]),
	}

	if c := certificateUnmarshal(k["certificate"]); c != nil {
		e.Certificate = *c
	}

	return e
}

func etcdSchema(computed bool) *schema.Schema {
	return optionalBlock(computed, false, func(computed bool) map[string]*schema.Schema {
		return map[string]*schema.Schema{
			"certificate":         certificateBlockSchema(computed),
			"ca":                  certificateBlockSchema(computed),
			"client_cns":          optionalStringList(computed),
			"peers":               stringMapSchema(false, false),
			"servers":             stringMapSchema(false, false),
			"peer_certificates":   certificateMapSchema(computed),
			"server_certificates": certificateMapSchema(computed),
			"client_certificates": certificateMapSchema(computed),
		}
	})
}
