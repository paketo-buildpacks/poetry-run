package poetryrun

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/paketo-buildpacks/libreload-packit"
	"github.com/paketo-buildpacks/packit/v2"
	"github.com/paketo-buildpacks/packit/v2/scribe"
)

// Build will return a packit.BuildFunc that will be invoked during the build
// phase of the buildpack lifecycle.
//
// Build assigns the image a launch process of 'poetry run <target>' where <target>
// is the key of a poetry script or system executable. This can be set via `BP_POETRY_RUN_TARGET`
// or inferred from pyproject.toml when there is exactly one script.
func Build(pyProjectParser PyProjectParser, logger scribe.Emitter, reloader libreload.Reloader) packit.BuildFunc {
	return func(context packit.BuildContext) (packit.BuildResult, error) {
		logger.Title("%s %s", context.BuildpackInfo.Name, context.BuildpackInfo.Version)

		args := []string{"run"}

		logger.Debug.Process("Finding the poetry run target")
		if runTarget, ok := os.LookupEnv("BP_POETRY_RUN_TARGET"); ok {
			args = append(args, strings.Split(runTarget, " ")...)
			logger.Debug.Subprocess("Found BP_POETRY_RUN_TARGET=%s", runTarget)
		} else {
			scriptKey, err := pyProjectParser.Parse(filepath.Join(context.WorkingDir, "pyproject.toml"))
			if err != nil {
				return packit.BuildResult{}, err
			}
			args = append(args, scriptKey)
			logger.Debug.Subprocess("Found pyproject.toml script=%s", scriptKey)
		}

		processes := make([]packit.Process, 0)

		poetryProcess := packit.Process{
			Type:    "web",
			Command: "poetry",
			Args:    args,
			Default: true,
			Direct:  true,
		}

		if shouldEnableReload, err := reloader.ShouldEnableLiveReload(); err != nil {
			return packit.BuildResult{}, err
		} else if shouldEnableReload {
			spec := libreload.ReloadableProcessSpec{
				WatchPaths: []string{context.WorkingDir},
			}
			nonReloadableProcess, reloadableProcess := reloader.TransformReloadableProcesses(poetryProcess, spec)
			processes = append(processes, reloadableProcess, nonReloadableProcess)
		} else {
			processes = append(processes, poetryProcess)
		}

		logger.LaunchProcesses(processes)

		return packit.BuildResult{
			Launch: packit.LaunchMetadata{
				Processes: processes,
			},
		}, nil
	}
}
