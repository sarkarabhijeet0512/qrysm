package operations

import (
	"testing"

	"github.com/theQRL/qrysm/v4/testing/spectest/shared/altair/operations"
)

func TestMinimal_Altair_Operations_SyncCommittee(t *testing.T) {
	operations.RunProposerSlashingTest(t, "minimal")
}
