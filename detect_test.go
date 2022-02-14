package poetryrun_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/paketo-buildpacks/packit"
	poetryrun "github.com/paketo-buildpacks/poetry-run"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testDetect(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect     = NewWithT(t).Expect
		detect     packit.DetectFunc
		workingDir string
	)

	it.Before(func() {
		var err error
		workingDir, err = ioutil.TempDir("", "working-dir")
		Expect(err).NotTo(HaveOccurred())

		contents := `
[tool.poetry.scripts]
my-script = "my_module:main"
`

		Expect(ioutil.WriteFile(filepath.Join(workingDir, "pyproject.toml"), []byte(contents), 0644)).To(Succeed())

		detect = poetryrun.Detect()
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
							Build: true,
						},
					},
					{
						Name: poetryrun.Poetry,
						Metadata: poetryrun.BuildPlanMetadata{
							Build: true,
						},
					},
					{
						Name: poetryrun.PoetryVenv,
						Metadata: poetryrun.BuildPlanMetadata{
							Build: true,
						},
					},
				},
			}))
		})

		context("when there is no pyproject.toml file", func() {
			it.Before(func() {
				Expect(os.Remove(filepath.Join(workingDir, "pyproject.toml"))).To(Succeed())
			})

			it("fails detection", func() {
				_, err := detect(packit.DetectContext{
					WorkingDir: workingDir,
				})
				Expect(err).To(MatchError(packit.Fail))
			})
		})

		context("when there are no poetry scripts in the pyproject.toml file", func() {
			it.Before(func() {
				Expect(os.Remove(filepath.Join(workingDir, "pyproject.toml"))).To(Succeed())
				contents := `
[some.other.valid.toml]
a-key = "a value"`
				Expect(ioutil.WriteFile(filepath.Join(workingDir, "pyproject.toml"), []byte(contents), 0644)).To(Succeed())
			})

			it("fails detection", func() {
				_, err := detect(packit.DetectContext{
					WorkingDir: workingDir,
				})
				Expect(err).To(MatchError(packit.Fail))
			})
		})

		context("failure cases", func() {
			context("when the pyproject.toml cannot be read", func() {
				it.Before(func() {
					Expect(os.Chmod(workingDir, 0000)).To(Succeed())
				})

				it.After(func() {
					Expect(os.Chmod(workingDir, os.ModePerm)).To(Succeed())
				})

				it("returns an error", func() {
					_, err := detect(packit.DetectContext{
						WorkingDir: workingDir,
					})
					Expect(err).To(MatchError(ContainSubstring("permission denied")))
				})
			})

			context("when the pyproject.toml does not contain the expected TOML structure", func() {
				it.Before(func() {
					Expect(os.Remove(filepath.Join(workingDir, "pyproject.toml"))).To(Succeed())
					contents := `
[tool.poetry.scripts]
a-key = [ "a value", "another value"]`

					Expect(ioutil.WriteFile(filepath.Join(workingDir, "pyproject.toml"), []byte(contents), 0644)).To(Succeed())
				})

				it("returns an error", func() {
					_, err := detect(packit.DetectContext{
						WorkingDir: workingDir,
					})
					Expect(err).To(MatchError(ContainSubstring("incompatible types: TOML")))
				})
			})
		})

	})
}
