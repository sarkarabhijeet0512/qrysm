package random

import (
	"testing"

	"github.com/theQRL/qrysm/v4/testing/spectest/shared/altair/sanity"
)

func TestMinimal_Altair_Random(t *testing.T) {
	sanity.RunBlockProcessingTest(t, "minimal", "random/random/pyspec_tests")
}
