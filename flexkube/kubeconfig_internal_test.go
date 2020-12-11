package flexkube

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/flexkube/libflexkube/pkg/kubernetes/client"
	"github.com/flexkube/libflexkube/pkg/types"
)

func TestKubeconfigUnmarshal(t *testing.T) {
	t.Parallel()

	c := client.Config{
		Server:            "127.0.0.1",
		CACertificate:     types.Certificate("foo"),
		ClientCertificate: "bar",
		ClientKey:         "1s",
	}

	e := []interface{}{
		map[string]interface{}{
			"server":             "127.0.0.1",
			"ca_certificate":     "foo",
			"client_certificate": "bar",
			"client_key":         "1s",
		},
	}

	if diff := cmp.Diff(kubeconfigUnmarshal(e[0]), c); diff != "" {
		t.Errorf("Unexpected diff:\n%s", diff)
	}
}

func TestKubeconfigUnmarshalEmpty(t *testing.T) {
	t.Parallel()

	var c client.Config

	if diff := cmp.Diff(kubeconfigUnmarshal(nil), c); diff != "" {
		t.Errorf("Unexpected diff:\n%s", diff)
	}
}

func TestKubeconfigUnmarshalEmptyBock(t *testing.T) {
	t.Parallel()

	var c client.Config

	if diff := cmp.Diff(kubeconfigUnmarshal(map[string]interface{}{}), c); diff != "" {
		t.Errorf("Unexpected diff:\n%s", diff)
	}
}
