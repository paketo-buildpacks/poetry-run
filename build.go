package poetryrun

import (
	"fmt"
	"path/filepath"

	"github.com/paketo-buildpacks/packit"
	"github.com/paketo-buildpacks/packit/scribe"
)

// Build will return a packit.BuildFunc that will be invoked during the build
// phase of the buildpack lifecycle.
//
// Build assigns the image a launch process of 'poetry run <script>' where <script>
// is the key of the only script present for Poetry to run.
func Build(pyProjectParser PyProjectParser, logger scribe.Emitter) packit.BuildFunc {
	return func(context packit.BuildContext) (packit.BuildResult, error) {
		logger.Title("%s %s", context.BuildpackInfo.Name, context.BuildpackInfo.Version)

		scriptKey, err := pyProjectParser.Parse(filepath.Join(context.WorkingDir, "pyproject.toml"))
		if err != nil {
			return packit.BuildResult{}, err
		}

		logger.Process("Assigning launch process")
		command := fmt.Sprintf("poetry run %s", scriptKey)
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
