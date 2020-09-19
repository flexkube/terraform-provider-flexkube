package flexkube

import "testing"

func TestCertificateMarshalNil(t *testing.T) {
	if c := certificateMarshal(false, nil); c != nil {
		t.Fatalf("from nil certificate, no data should be returned")
	}
}

func TestCertificateUnMarshalNil(t *testing.T) {
	if c := certificateUnmarshal(nil); c != nil {
		t.Fatalf("from nil certificate, no data should be returned")
	}
}
