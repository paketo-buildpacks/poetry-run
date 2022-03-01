package poetryrun

import (
	"errors"
	"os"

	"github.com/BurntSushi/toml"
)

type PyProjectConfig struct {
	Tool struct {
		Poetry struct {
			Scripts map[string]string `toml:"scripts"`
		} `toml:"poetry"`
	} `toml:"tool"`
}

type PyProjectConfigParser struct {
}

func NewPyProjectConfigParser() PyProjectConfigParser {
	return PyProjectConfigParser{}
}

// Parse returns the name of the script for Poetry to execute
// If there is no file, no script to run, or multiple scripts to run,
// Parse returns an empty string and a nil error
// If there is an error reading the file, Parse returns an error
func (p PyProjectConfigParser) Parse(filepath string) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", nil
		}

		return "", err
	}

	var pyProjectConfig PyProjectConfig

	_, err = toml.NewDecoder(file).Decode(&pyProjectConfig)
	if err != nil {
		return "", err
	}

	if len(pyProjectConfig.Tool.Poetry.Scripts) != 1 {
		return "", nil
	}

	for key := range pyProjectConfig.Tool.Poetry.Scripts {
		return key, nil
	}

	panic("should not be able to get here")
}
