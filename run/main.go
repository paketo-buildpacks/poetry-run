package main

import (
	"os"

	"github.com/paketo-buildpacks/libreload-packit/watchexec"
	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/scribe"
	poetryrun "github.com/paketo-buildpacks/poetry-run"
)

func main() {
	logger := scribe.NewEmitter(os.Stdout).WithLevel(os.Getenv("BP_LOG_LEVEL"))
	pyProjectParser := poetryrun.NewPyProjectConfigParser()

	reload := watchexec.NewWatchexecReloader()
	packit.Run(
		poetryrun.Detect(pyProjectParser, reload),
		poetryrun.Build(
			pyProjectParser,
			logger,
			reload,
		),
	)
}
