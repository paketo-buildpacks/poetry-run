package integration_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/paketo-buildpacks/occam"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	. "github.com/onsi/gomega"
)

var buildpackInfo struct {
	Buildpack struct {
		ID   string
		Name string
	}
	Metadata struct {
		Dependencies []struct {
			Version string
		}
	}
}

var settings struct {
	Buildpacks struct {
		CPython struct {
			Online string
		}
		Pip struct {
			Online string
		}
		Poetry struct {
			Online string
		}
		PoetryInstall struct {
			Online string
		}
		PoetryRun struct {
			Online string
		}
		BuildPlan struct {
			Online string
		}
	}

	Config struct {
		CPython       string `json:"cpython"`
		Pip           string `json:"pip"`
		Poetry        string `json:"poetry"`
		PoetryInstall string `json:"poetry-install"`
		BuildPlan     string `json:"build-plan"`
	}
}

func TestIntegration(t *testing.T) {
	Expect := NewWithT(t).Expect

	file, err := os.Open("../integration.json")
	Expect(err).NotTo(HaveOccurred())

	Expect(json.NewDecoder(file).Decode(&settings.Config)).To(Succeed())
	Expect(file.Close()).To(Succeed())

	file, err = os.Open("../buildpack.toml")
	Expect(err).NotTo(HaveOccurred())

	_, err = toml.NewDecoder(file).Decode(&buildpackInfo)
	Expect(err).NotTo(HaveOccurred())

	root, err := filepath.Abs("./..")
	Expect(err).ToNot(HaveOccurred())

	buildpackStore := occam.NewBuildpackStore()

	settings.Buildpacks.PoetryRun.Online, err = buildpackStore.Get.
		WithVersion("1.2.3").
		Execute(root)
	Expect(err).NotTo(HaveOccurred())

	settings.Buildpacks.PoetryInstall.Online, err = buildpackStore.Get.
		Execute(settings.Config.PoetryInstall)
	Expect(err).NotTo(HaveOccurred())

	settings.Buildpacks.Poetry.Online, err = buildpackStore.Get.
		Execute(settings.Config.Poetry)
	Expect(err).NotTo(HaveOccurred())

	settings.Buildpacks.Pip.Online, err = buildpackStore.Get.
		Execute(settings.Config.Pip)
	Expect(err).NotTo(HaveOccurred())

	settings.Buildpacks.CPython.Online, err = buildpackStore.Get.
		Execute(settings.Config.CPython)
	Expect(err).NotTo(HaveOccurred())

	settings.Buildpacks.BuildPlan.Online, err = buildpackStore.Get.
		Execute(settings.Config.BuildPlan)
	Expect(err).NotTo(HaveOccurred())

	SetDefaultEventuallyTimeout(5 * time.Second)

	suite := spec.New("Integration", spec.Report(report.Terminal{}))
	suite("Default", testDefault, spec.Parallel())
	suite.Run(t)
}
