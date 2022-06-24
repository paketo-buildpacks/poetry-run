package poetryrun_test

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/scribe"
	poetryrun "github.com/paketo-buildpacks/poetry-run"
	"github.com/paketo-buildpacks/poetry-run/fakes"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
	. "github.com/paketo-buildpacks/occam/matchers"
)

func testBuild(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		layersDir  string
		workingDir string
		cnbDir     string
		buffer     *bytes.Buffer

		pyProjectParser *fakes.PyProjectParser

		build        packit.BuildFunc
		buildContext packit.BuildContext
	)

	it.Before(func() {
		var err error
		layersDir, err = os.MkdirTemp("", "layers")
		Expect(err).NotTo(HaveOccurred())

		cnbDir, err = os.MkdirTemp("", "cnb")
		Expect(err).NotTo(HaveOccurred())

		workingDir, err = os.MkdirTemp("", "working-dir")
		Expect(err).NotTo(HaveOccurred())

		buffer = bytes.NewBuffer(nil)
		logger := scribe.NewEmitter(buffer).WithLevel("DEBUG")

		pyProjectParser = &fakes.PyProjectParser{}
		pyProjectParser.ParseCall.Returns.String = "some-script"

		build = poetryrun.Build(pyProjectParser, logger)
		buildContext = packit.BuildContext{
			WorkingDir: workingDir,
			CNBPath:    cnbDir,
			Stack:      "some-stack",
			BuildpackInfo: packit.BuildpackInfo{
				Name:    "Some Buildpack",
				Version: "some-version",
			},
			Plan: packit.BuildpackPlan{
				Entries: []packit.BuildpackPlanEntry{},
			},
			Layers: packit.Layers{Path: layersDir},
		}
	})

	it.After(func() {
		Expect(os.RemoveAll(layersDir)).To(Succeed())
		Expect(os.RemoveAll(cnbDir)).To(Succeed())
		Expect(os.RemoveAll(workingDir)).To(Succeed())
	})

	context("with BP_POETRY_RUN_TARGET not set", func() {
		it("returns a result that sets the 'poetry run' launch command", func() {
			result, err := build(buildContext)
			Expect(err).NotTo(HaveOccurred())

			Expect(result).To(Equal(packit.BuildResult{
				Plan: packit.BuildpackPlan{
					Entries: nil,
				},
				Layers: nil,
				Launch: packit.LaunchMetadata{
					Processes: []packit.Process{
						{
							Type:    "web",
							Command: "poetry",
							Args: []string{
								"run",
								"some-script",
							},
							Default: true,
							Direct:  true,
						},
					},
				},
			}))

			Expect(buffer.String()).To(ContainLines(
				ContainSubstring("Finding the poetry run target"),
				ContainSubstring("Found pyproject.toml script=some-script"),
				ContainSubstring("Assigning launch processes:"),
				ContainSubstring("web (default): poetry run some-script"),
			))
		})

		context("when BP_LIVE_RELOAD_ENABLED=true in the build environment", func() {
			it.Before(func() {
				Expect(os.Setenv("BP_LIVE_RELOAD_ENABLED", "true")).To(Succeed())
			})

			it.After(func() {
				Expect(os.Unsetenv("BP_LIVE_RELOAD_ENABLED")).To(Succeed())
			})

			it("adds a reloadable start command and makes it the default", func() {
				result, err := build(buildContext)
				Expect(err).NotTo(HaveOccurred())

				Expect(result).To(Equal(packit.BuildResult{
					Plan: packit.BuildpackPlan{
						Entries: nil,
					},
					Layers: nil,
					Launch: packit.LaunchMetadata{
						Processes: []packit.Process{
							{
								Type:    "web",
								Command: "watchexec",
								Args: []string{
									"--restart",
									"--watch", workingDir,
									"--shell", "none",
									"--",
									"poetry",
									"run",
									"some-script",
								},
								Default: true,
								Direct:  true,
							},
							{
								Type:    "no-reload",
								Command: "poetry",
								Args:    []string{"run", "some-script"},
								Direct:  true,
							},
						},
					},
				}))
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

		it("will use the value of BP_POETRY_RUN_TARGET and not use the pyproject.toml parser", func() {
			result, err := build(buildContext)
			Expect(err).NotTo(HaveOccurred())

			Expect(result).To(Equal(packit.BuildResult{
				Plan: packit.BuildpackPlan{
					Entries: nil,
				},
				Layers: nil,
				Launch: packit.LaunchMetadata{
					Processes: []packit.Process{
						{
							Type:    "web",
							Command: "poetry",
							Args: []string{
								"run",
								"a",
								"custom",
								"command",
							},
							Default: true,
							Direct:  true,
						},
					},
				},
			}))

			Expect(buffer.String()).To(ContainLines(
				ContainSubstring("Finding the poetry run target"),
				ContainSubstring("Found BP_POETRY_RUN_TARGET=a custom command"),
				ContainSubstring("Assigning launch processes:"),
				ContainSubstring("web (default): poetry run a custom command"),
			))
			Expect(pyProjectParser.ParseCall.CallCount).To(Equal(0))
		})

		context("when BP_LIVE_RELOAD_ENABLED=true in the build environment", func() {
			it.Before(func() {
				Expect(os.Setenv("BP_LIVE_RELOAD_ENABLED", "true")).To(Succeed())
			})

			it.After(func() {
				Expect(os.Unsetenv("BP_LIVE_RELOAD_ENABLED")).To(Succeed())
			})

			it("adds a reloadable start command and makes it the default", func() {
				result, err := build(buildContext)
				Expect(err).NotTo(HaveOccurred())

				Expect(result).To(Equal(packit.BuildResult{
					Plan: packit.BuildpackPlan{
						Entries: nil,
					},
					Layers: nil,
					Launch: packit.LaunchMetadata{
						Processes: []packit.Process{
							{
								Type:    "web",
								Command: "watchexec",
								Args: []string{
									"--restart",
									"--watch", workingDir,
									"--shell", "none",
									"--",
									"poetry",
									"run",
									"a",
									"custom",
									"command",
								},
								Default: true,
								Direct:  true,
							},
							{
								Type:    "no-reload",
								Command: "poetry",
								Args: []string{
									"run",
									"a",
									"custom",
									"command",
								},
								Direct: true,
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

			context(" and the pyproject.toml parser returns an error", func() {
				it.Before(func() {
					pyProjectParser.ParseCall.Returns.Error = fmt.Errorf("some error")
				})

				it("returns the error", func() {
					_, err := build(buildContext)
					Expect(err).To(MatchError(ContainSubstring("some error")))
				})
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
				_, err := build(buildContext)
				Expect(err).To(MatchError(ContainSubstring("failed to parse BP_LIVE_RELOAD_ENABLED value not-a-bool")))
			})
		})
	})
}
