package poetryrun_test

import (
	"os"
	"path/filepath"
	"testing"

	poetryrun "github.com/paketo-buildpacks/poetry-run"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testPyProjectConfigParser(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect     = NewWithT(t).Expect
		parser     poetryrun.PyProjectConfigParser
		workingDir string
	)

	it.Before(func() {
		var err error
		workingDir, err = os.MkdirTemp("", "working-dir")
		Expect(err).NotTo(HaveOccurred())

		contents := `
[tool.poetry.scripts]
my-script = "my_module:main"
`

		Expect(os.WriteFile(filepath.Join(workingDir, "pyproject.toml"), []byte(contents), 0644)).To(Succeed())

		parser = poetryrun.NewPyProjectConfigParser()
	})

	context("parsing", func() {
		it("returns the key of the only provided script", func() {
			result, err := parser.Parse(filepath.Join(workingDir, "pyproject.toml"))
			Expect(err).NotTo(HaveOccurred())

			Expect(result).To(Equal("my-script"))

		})

		context("when there is no pyproject.toml file", func() {
			it.Before(func() {
				Expect(os.Remove(filepath.Join(workingDir, "pyproject.toml"))).To(Succeed())
			})
			it("returns an empty string without error", func() {
				script, err := parser.Parse(filepath.Join(workingDir, "pyproject.toml"))
				Expect(err).NotTo(HaveOccurred())

				Expect(script).To(BeEmpty())
			})
		})

		context("when there are no poetry scripts in the pyproject.toml file", func() {
			it.Before(func() {
				Expect(os.Remove(filepath.Join(workingDir, "pyproject.toml"))).To(Succeed())
				contents := `
[some.other.valid.toml]
a-key = "a value"`
				Expect(os.WriteFile(filepath.Join(workingDir, "pyproject.toml"), []byte(contents), 0644)).To(Succeed())
			})

			it("returns an empty string without error", func() {
				script, err := parser.Parse(filepath.Join(workingDir, "pyproject.toml"))
				Expect(err).NotTo(HaveOccurred())

				Expect(script).To(BeEmpty())
			})
		})

		context("when there are multiple poetry scripts in the pyproject.toml file", func() {
			it.Before(func() {
				Expect(os.Remove(filepath.Join(workingDir, "pyproject.toml"))).To(Succeed())
				contents := `
[tool.poetry.scripts]
my-script = "my_module:main"
my-other-script = "my_other_module:main"
`
				Expect(os.WriteFile(filepath.Join(workingDir, "pyproject.toml"), []byte(contents), 0644)).To(Succeed())
			})

			it("returns an empty string without error", func() {
				script, err := parser.Parse(filepath.Join(workingDir, "pyproject.toml"))
				Expect(err).NotTo(HaveOccurred())

				Expect(script).To(BeEmpty())
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
					_, err := parser.Parse(filepath.Join(workingDir, "pyproject.toml"))
					Expect(err).To(MatchError(ContainSubstring("permission denied")))
				})
			})

			context("when the pyproject.toml does not contain the expected TOML structure", func() {
				it.Before(func() {
					Expect(os.Remove(filepath.Join(workingDir, "pyproject.toml"))).To(Succeed())
					contents := `
[tool.poetry.scripts]
a-key = [ "a value", "another value"]`

					Expect(os.WriteFile(filepath.Join(workingDir, "pyproject.toml"), []byte(contents), 0644)).To(Succeed())
				})

				it("returns an error", func() {
					_, err := parser.Parse(filepath.Join(workingDir, "pyproject.toml"))
					Expect(err).To(MatchError(ContainSubstring("incompatible types: TOML")))
				})
			})
		})

	})
}
