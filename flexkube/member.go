package flexkube

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/flexkube/libflexkube/pkg/etcd"
)

func membersUnmarshal(i interface{}) map[string]etcd.Member {
	j := i.([]interface{})

	members := map[string]etcd.Member{}

	for _, k := range j {
		t := k.(map[string]interface{})

		m := etcd.Member{
			Name:              t["name"].(string),
			Image:             t["image"].(string),
			CACertificate:     t["ca_certificate"].(string),
			PeerCertificate:   t["peer_certificate"].(string),
			PeerKey:           t["peer_key"].(string),
			PeerAddress:       t["peer_address"].(string),
			InitialCluster:    t["initial_cluster"].(string),
			PeerCertAllowedCN: t["peer_cert_allowed_cn"].(string),
			ServerCertificate: t["server_certificate"].(string),
			ServerKey:         t["server_key"].(string),
			ServerAddress:     t["server_address"].(string),
		}

		if v, ok := t["host"]; ok && len(v.([]interface{})) == 1 {
			m.Host = hostUnmarshal(v.([]interface{})[0])
		}

		members[t["name"].(string)] = m
	}

	return members
}

func memberSchema() *schema.Schema {
	return requiredList(false, false, func(computed bool) *schema.Resource {
		return &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name":                 requiredString(false),
				"image":                optionalString(false),
				"host":                 hostSchema(false),
				"ca_certificate":       optionalString(false),
				"peer_certificate":     optionalString(false),
				"peer_key":             sensitiveString(true),
				"peer_address":         optionalString(false),
				"initial_cluster":      optionalString(false),
				"peer_cert_allowed_cn": optionalString(false),
				"server_certificate":   optionalString(false),
				"server_key":           sensitiveString(true),
				"server_address":       requiredString(false),
			},
		}
	})
}
