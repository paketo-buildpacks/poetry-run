package poetryrun_test

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/paketo-buildpacks/packit/v2"
	poetryrun "github.com/paketo-buildpacks/poetry-run"
	"github.com/paketo-buildpacks/poetry-run/fakes"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testDetect(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect
		detect packit.DetectFunc

		pyProjectParser *fakes.PyProjectParser

		workingDir string
	)

	it.Before(func() {
		var err error
		workingDir, err = ioutil.TempDir("", "working-dir")
		Expect(err).NotTo(HaveOccurred())

		Expect(ioutil.WriteFile(filepath.Join(workingDir, "pyproject.toml"), []byte("some contents"), 0644)).To(Succeed())

		pyProjectParser = &fakes.PyProjectParser{}
		pyProjectParser.ParseCall.Returns.String = "some-script"

		detect = poetryrun.Detect(pyProjectParser)
	})

	context("detection", func() {
		it("returns a build plan", func() {
			result, err := detect(packit.DetectContext{
				WorkingDir: workingDir,
			})
			Expect(err).NotTo(HaveOccurred())

			Expect(result.Plan).To(Equal(packit.BuildPlan{
				Provides: []packit.BuildPlanProvision{},
				Requires: []packit.BuildPlanRequirement{
					{
						Name: poetryrun.CPython,
						Metadata: poetryrun.BuildPlanMetadata{
							Launch: true,
						},
					},
					{
						Name: poetryrun.Poetry,
						Metadata: poetryrun.BuildPlanMetadata{
							Launch: true,
						},
					},
					{
						Name: poetryrun.PoetryVenv,
						Metadata: poetryrun.BuildPlanMetadata{
							Launch: true,
						},
					},
				},
			}))

			Expect(pyProjectParser.ParseCall.Receives.String).To(Equal(filepath.Join(workingDir, "pyproject.toml")))
		})

		context("when there is no script returned by the paser", func() {
			it.Before(func() {
				pyProjectParser.ParseCall.Returns.String = ""
			})

			it("fails detection", func() {
				_, err := detect(packit.DetectContext{
					WorkingDir: workingDir,
				})

				Expect(err).To(MatchError(packit.Fail))
			})
		})

		context("failure cases", func() {
			context("when the pyproject.toml parser returns an error", func() {
				it.Before(func() {
					pyProjectParser.ParseCall.Returns.Error = fmt.Errorf("some error")
				})

				it("returns the error", func() {
					_, err := detect(packit.DetectContext{
						WorkingDir: workingDir,
					})
					Expect(err).To(MatchError(ContainSubstring("some error")))
				})
			})
		})

	})
}
