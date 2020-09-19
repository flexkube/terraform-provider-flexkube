package flexkube

import (
	"testing"
)

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("Provider internal validation should succeed, got: %v", err)
	}
}
