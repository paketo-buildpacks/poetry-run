package poetryrun

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/paketo-buildpacks/packit"
)

// BuildPlanMetadata is the buildpack specific data included in build plan
// requirements.
type BuildPlanMetadata struct {
	// Build denotes the dependency is needed at build-time.
	Build bool `toml:"build"`
}

type PyProjectConfig struct {
	Tool struct {
		Poetry struct {
			Scripts map[string]string
		}
	}
}

// Detect will return a packit.DetectFunc that will be invoked during the
// detect phase of the buildpack lifecycle.
//
// Detection will contribute a Build Plan that provides site-packages,
// and requires cpython and pip at build.
//
// Detection is contingent on there being one or more scripts to run
// defined in the pyproject.toml under [tool.poetry.scripts]
func Detect() packit.DetectFunc {
	return func(context packit.DetectContext) (packit.DetectResult, error) {
		file, err := os.Open(filepath.Join(context.WorkingDir, "pyproject.toml"))
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return packit.DetectResult{}, packit.Fail
			}

			return packit.DetectResult{}, err
		}

		var pyProjectConfig PyProjectConfig

		_, err = toml.NewDecoder(file).Decode(&pyProjectConfig)
		if err != nil {
			return packit.DetectResult{}, err
		}

		if len(pyProjectConfig.Tool.Poetry.Scripts) == 0 {
			return packit.DetectResult{}, packit.Fail
		}

		return packit.DetectResult{
			Plan: packit.BuildPlan{
				Provides: []packit.BuildPlanProvision{},
				Requires: []packit.BuildPlanRequirement{
					{
						Name: CPython,
						Metadata: BuildPlanMetadata{
							Build: true,
						},
					},
					{
						Name: Poetry,
						Metadata: BuildPlanMetadata{
							Build: true,
						},
					},
					{
						Name: PoetryVenv,
						Metadata: BuildPlanMetadata{
							Build: true,
						},
					},
				},
			},
		}, nil
	}
}
