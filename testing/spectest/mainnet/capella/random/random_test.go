package random

import (
	"testing"

	"github.com/theQRL/qrysm/v4/testing/spectest/shared/capella/sanity"
)

func TestMainnet_Capella_Random(t *testing.T) {
	sanity.RunBlockProcessingTest(t, "mainnet", "random/random/pyspec_tests")
}
