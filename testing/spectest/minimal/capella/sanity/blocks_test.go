package sanity

import (
	"testing"

	"github.com/theQRL/qrysm/v4/testing/spectest/shared/capella/sanity"
)

func TestMinimal_Capella_Sanity_Blocks(t *testing.T) {
	sanity.RunBlockProcessingTest(t, "minimal", "sanity/blocks/pyspec_tests")
}
