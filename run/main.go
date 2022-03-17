package main

import (
	"os"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/scribe"
	poetryrun "github.com/paketo-buildpacks/poetry-run"
)

func main() {
	logger := scribe.NewEmitter(os.Stdout)
	pyProjectParser := poetryrun.NewPyProjectConfigParser()

	packit.Run(
		poetryrun.Detect(pyProjectParser),
		poetryrun.Build(
			pyProjectParser,
			logger,
		),
	)
}
