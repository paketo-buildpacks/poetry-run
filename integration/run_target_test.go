package integration_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/paketo-buildpacks/occam"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
	. "github.com/paketo-buildpacks/occam/matchers"
)

func testRunTargets(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect     = NewWithT(t).Expect
		Eventually = NewWithT(t).Eventually

		pack   occam.Pack
		docker occam.Docker
	)

	it.Before(func() {
		pack = occam.NewPack()
		docker = occam.NewDocker()
	})

	context("when BP_POETRY_RUN_TARGET is set", func() {
		var (
			image     occam.Image
			container occam.Container
			name      string
			source    string
		)

		it.Before(func() {
			var err error
			name, err = occam.RandomName()
			Expect(err).NotTo(HaveOccurred())

			source, err = occam.Source(filepath.Join("testdata", "run_target_app"))
			Expect(err).NotTo(HaveOccurred())
		})

		it.After(func() {
			Expect(docker.Container.Remove.Execute(container.ID)).To(Succeed())
			Expect(docker.Image.Remove.Execute(image.ID)).To(Succeed())
			Expect(docker.Volume.Remove.Execute(occam.CacheVolumeNames(name))).To(Succeed())
			Expect(os.RemoveAll(source)).To(Succeed())
		})

		context("when BP_POETRY_RUN_TARGET is set to an executable command", func() {
			it("builds and runs successfully", func() {
				var err error
				var logs fmt.Stringer

				image, logs, err = pack.WithNoColor().Build.
					WithPullPolicy("never").
					WithBuildpacks(
						settings.Buildpacks.CPython.Online,
						settings.Buildpacks.Pip.Online,
						settings.Buildpacks.Poetry.Online,
						settings.Buildpacks.PoetryInstall.Online,
						settings.Buildpacks.PoetryRun.Online,
						settings.Buildpacks.BuildPlan.Online,
					).
					WithEnv(map[string]string{
						"BP_POETRY_RUN_TARGET": "python -V",
					}).
					Execute(name, source)
				Expect(err).ToNot(HaveOccurred(), logs.String)

				Expect(logs).To(ContainLines(
					MatchRegexp(fmt.Sprintf(`%s \d+\.\d+\.\d+`, buildpackInfo.Buildpack.Name)),
					"  Assigning launch processes:",
					"    web (default): poetry run python -V",
				))

				container, err = docker.Container.Run.
					Execute(image.ID)
				Expect(err).ToNot(HaveOccurred())

				Eventually(func() string {
					cLogs, err := docker.Container.Logs.Execute(container.ID)
					Expect(err).NotTo(HaveOccurred())
					return cLogs.String()
				}).Should(MatchRegexp(`Python 3\.\d+\.\d+`))
			})
		})

		context("when BP_POETRY_RUN_TARGET is set to a script key", func() {
			it("builds and runs successfully", func() {
				var err error
				var logs fmt.Stringer

				image, logs, err = pack.WithNoColor().Build.
					WithPullPolicy("never").
					WithBuildpacks(
						settings.Buildpacks.CPython.Online,
						settings.Buildpacks.Pip.Online,
						settings.Buildpacks.Poetry.Online,
						settings.Buildpacks.PoetryInstall.Online,
						settings.Buildpacks.PoetryRun.Online,
						settings.Buildpacks.BuildPlan.Online,
					).
					WithEnv(map[string]string{
						"BP_POETRY_RUN_TARGET": "working-script-key",
					}).
					Execute(name, source)
				Expect(err).ToNot(HaveOccurred(), logs.String)

				Expect(logs).To(ContainLines(
					MatchRegexp(fmt.Sprintf(`%s \d+\.\d+\.\d+`, buildpackInfo.Buildpack.Name)),
					"  Assigning launch processes:",
					"    web (default): poetry run working-script-key",
				))

				container, err = docker.Container.Run.
					WithEnv(map[string]string{"PORT": "8080"}).
					WithPublish("8080").
					Execute(image.ID)
				Expect(err).ToNot(HaveOccurred())

				Eventually(container).Should(BeAvailable())
				Eventually(container).Should(Serve(ContainSubstring("Hello, World!")).OnPort(8080))
			})
		})
	})
}
