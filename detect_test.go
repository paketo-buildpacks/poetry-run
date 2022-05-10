package poetryrun_test

import (
	"fmt"
	"os"
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
	)

	it.Before(func() {
		pyProjectParser = &fakes.PyProjectParser{}
		pyProjectParser.ParseCall.Returns.String = "some-script"

		detect = poetryrun.Detect(pyProjectParser)
	})

	context("detection", func() {
		it("returns a build plan", func() {
			result, err := detect(packit.DetectContext{
				WorkingDir: "a-working-dir",
			})
			Expect(err).NotTo(HaveOccurred())

			Expect(result.Plan).To(Equal(packit.BuildPlan{
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

			Expect(pyProjectParser.ParseCall.Receives.String).To(Equal(filepath.Join("a-working-dir", "pyproject.toml")))
		})

		context("when BP_LIVE_RELOAD_ENABLED=true", func() {
			it.Before(func() {
				Expect(os.Setenv("BP_LIVE_RELOAD_ENABLED", "true")).To(Succeed())
			})

			it.After(func() {
				Expect(os.Unsetenv("BP_LIVE_RELOAD_ENABLED")).To(Succeed())
			})

			it("requires watchexec at launch", func() {
				result, err := detect(packit.DetectContext{})
				Expect(err).NotTo(HaveOccurred())

				Expect(result.Plan).To(Equal(packit.BuildPlan{
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
						{
							Name: poetryrun.Watchexec,
							Metadata: poetryrun.BuildPlanMetadata{
								Launch: true,
							},
						},
					},
				}))
			})
		})

		context("when there is no script returned by the parser", func() {
			it.Before(func() {
				pyProjectParser.ParseCall.Returns.String = ""
			})

			it("fails detection", func() {
				_, err := detect(packit.DetectContext{})

				Expect(err).To(MatchError(packit.Fail.WithMessage("Expects one and exactly one script defined in pyproject.toml")))
			})
		})

		context("failure cases", func() {
			context("when the pyproject.toml parser returns an error", func() {
				it.Before(func() {
					pyProjectParser.ParseCall.Returns.Error = fmt.Errorf("some error")
				})

				it("returns the error", func() {
					_, err := detect(packit.DetectContext{})
					Expect(err).To(MatchError(ContainSubstring("some error")))
				})
			})

			context("when BP_LIVE_RELOAD_ENABLED is set to an invalid value", func() {
				it.Before(func() {
					Expect(os.Setenv("BP_LIVE_RELOAD_ENABLED", "not-a-bool")).To(Succeed())
				})

				it.After(func() {
					Expect(os.Unsetenv("BP_LIVE_RELOAD_ENABLED")).To(Succeed())
				})

				it("returns an error", func() {
					_, err := detect(packit.DetectContext{})
					Expect(err).To(MatchError(ContainSubstring("failed to parse BP_LIVE_RELOAD_ENABLED value not-a-bool")))
				})
			})
		})
	})
}
