package poetryrun_test

import (
	"errors"
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
		reloader        *fakes.Reloader
	)

	it.Before(func() {
		pyProjectParser = &fakes.PyProjectParser{}
		reloader = &fakes.Reloader{}

		detect = poetryrun.Detect(pyProjectParser, reloader)
	})

	context("with BP_POETRY_RUN_TARGET not set", func() {
		it.Before(func() {
			Expect(os.Unsetenv("BP_POETRY_RUN_TARGET")).To(Succeed())
		})

		context("when pyproject.toml parser returns a valid script", func() {
			it.Before(func() {
				pyProjectParser.ParseCall.Returns.String = "some-script"
			})

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

			context("when live reload is enabled", func() {
				it.Before(func() {
					reloader.ShouldEnableLiveReloadCall.Returns.Bool = true
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
		})

		context("when the pyproject.toml parser cannot find a script", func() {
			it.Before(func() {
				pyProjectParser.ParseCall.Returns.String = ""
			})

			it("fails detection", func() {
				_, err := detect(packit.DetectContext{})

				Expect(err).To(MatchError(packit.Fail.WithMessage("Expects one and exactly one script defined in pyproject.toml")))
			})
		})
	})

	context("with BP_POETRY_RUN_TARGET set", func() {
		it.Before(func() {
			Expect(os.Setenv("BP_POETRY_RUN_TARGET", "a custom command")).To(Succeed())
		})

		it.After(func() {
			Expect(os.Unsetenv("BP_POETRY_RUN_TARGET")).To(Succeed())
		})

		it("returns a build plan", func() {
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
				},
			}))

			Expect(pyProjectParser.ParseCall.CallCount).To(Equal(0))
		})

		context("when live reload is enabled", func() {
			it.Before(func() {
				reloader.ShouldEnableLiveReloadCall.Returns.Bool = true
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
	})

	context("failure cases", func() {
		context("when BP_POETRY_RUN_TARGET is not set", func() {
			it.Before(func() {
				Expect(os.Unsetenv("BP_POETRY_RUN_TARGET")).To(Succeed())
			})

			context("when the pyproject.toml parser returns an error", func() {
				it.Before(func() {
					pyProjectParser.ParseCall.Returns.Error = fmt.Errorf("some error")
				})

				it("returns the error", func() {
					_, err := detect(packit.DetectContext{})
					Expect(err).To(MatchError(ContainSubstring("some error")))
				})
			})
		})

		context("when the reloader returns an error", func() {
			it.Before(func() {
				reloader.ShouldEnableLiveReloadCall.Returns.Error = errors.New("failed to parse")
			})

			context("when pyproject.toml parser returns a valid script", func() {
				it.Before(func() {
					pyProjectParser.ParseCall.Returns.String = "some-script"
				})

				it("returns an error", func() {
					_, err := detect(packit.DetectContext{})
					Expect(err).To(MatchError(ContainSubstring("failed to parse")))
				})
			})
		})
	})
}
