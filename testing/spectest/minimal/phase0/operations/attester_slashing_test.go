package operations

import (
	"testing"

	"github.com/theQRL/qrysm/v4/testing/spectest/shared/phase0/operations"
)

func TestMinimal_Phase0_Operations_AttesterSlashing(t *testing.T) {
	operations.RunAttesterSlashingTest(t, "minimal")
}
