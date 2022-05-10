package poetryrun

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/scribe"
)

// Build will return a packit.BuildFunc that will be invoked during the build
// phase of the buildpack lifecycle.
//
// Build assigns the image a launch process of 'poetry run <target>' where <target>
// is the key of a poetry script or system executable. This can be set via `BP_POETRY_RUN_TARGET`
// or inferred from pyproject.toml when there is exactly one script.
func Build(pyProjectParser PyProjectParser, logger scribe.Emitter) packit.BuildFunc {
	return func(context packit.BuildContext) (packit.BuildResult, error) {
		logger.Title("%s %s", context.BuildpackInfo.Name, context.BuildpackInfo.Version)

		var command string

		logger.Debug.Process("Finding the poetry run target")
		if runTarget, ok := os.LookupEnv("BP_POETRY_RUN_TARGET"); ok {
			command = fmt.Sprintf("poetry run %s", runTarget)
			logger.Debug.Subprocess("Found BP_POETRY_RUN_TARGET=%s", runTarget)
		} else {
			scriptKey, err := pyProjectParser.Parse(filepath.Join(context.WorkingDir, "pyproject.toml"))
			if err != nil {
				return packit.BuildResult{}, err
			}
			command = fmt.Sprintf("poetry run %s", scriptKey)
			logger.Debug.Subprocess("Found pyproject.toml script=%s", scriptKey)
		}

		logger.Process("Assigning launch process")
		logger.Subprocess("web: %s", command)

		return packit.BuildResult{
			Launch: packit.LaunchMetadata{
				Processes: []packit.Process{
					{
						Type:    "web",
						Command: command,
						Default: true,
					},
				},
			},
		}, nil
	}
}
