package operations

import (
	"testing"

	"github.com/theQRL/qrysm/v4/testing/spectest/shared/deneb/operations"
)

func TestMinimal_Deneb_Operations_Attestation(t *testing.T) {
	operations.RunAttestationTest(t, "minimal")
}
