package shadowpredecl

import (
	"testing"

	"github.com/theQRL/qrysm/v4/build/bazel"
	"golang.org/x/tools/go/analysis/analysistest"
)

func init() {
	if bazel.BuiltWithBazel() {
		bazel.SetGoEnv()
	}
}

func TestAnalyzer(t *testing.T) {
	testdata := bazel.TestDataPath(t)
	analysistest.TestData = func() string { return testdata }
	analysistest.Run(t, testdata, Analyzer)
}
