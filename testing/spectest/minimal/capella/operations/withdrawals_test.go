package operations

import (
	"testing"

	"github.com/theQRL/qrysm/v4/testing/spectest/shared/capella/operations"
)

func TestMinimal_Capella_Operations_Withdrawals(t *testing.T) {
	operations.RunWithdrawalsTest(t, "minimal")
}
