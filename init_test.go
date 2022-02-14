package poetryrun_test

import (
	"testing"

	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestUnitPoetryRun(t *testing.T) {
	suite := spec.New("poetryrun", spec.Report(report.Terminal{}))
	suite("Detect", testDetect)
	suite("Build", testBuild)
	suite.Run(t)
}
