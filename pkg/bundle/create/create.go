package create

import (
	"fmt"
	"os"

	bundle "github.com/RedHatGov/bundle/pkg/bundle"
	archive "github.com/RedHatGov/bundle/pkg/bundle/archive"
	"github.com/sirupsen/logrus"
)

// CreateFull performs all tasks in creating full imagesets
func CreateFull(ext string, rootDir string, segSize int64) error {

	err := bundle.MakeCreateDirs(rootDir)
	if err != nil {
		logrus.Error(err)
		return err
	}
	// Open Metadata
	metadata, err := bundle.ReadMeta(rootDir)
	if err != nil {
		logrus.Error(err)
		return err
	}
	lastRun := metadata.Imagesets[len(metadata.Imagesets)-1]
	logrus.Info(lastRun)

	// Read the bundle-config.yaml
	config, err := bundle.ReadBundleConfig(rootDir)
	if err != nil {
		logrus.Error(err)
		return err
	}
	logrus.Info(config)

	if config.Mirror.Ocp.Channels != nil {
		bundle.GetReleases(&lastRun, config, rootDir)
	}
	/*if config.Mirror.Operators != nil {
	//GetOperators(*config, rootDir)
	//}
	//if config.Mirror.Samples != nil {
	//GetSamples(*config, rootDir)
	//}*/

	// User defined image download

	if config.Mirror.AdditionalImages != nil {
		if err := bundle.GetAdditional(config, rootDir); err != nil {
			return fmt.Errorf("error downloading additional images: %v", err)
		}
	}

	// Get current working directory
	cwd, err := os.Getwd()

	if err != nil {
		return err
	}

	// Create archiver
	arc, err := archive.NewArchiver(ext)

	if err != nil {
		return fmt.Errorf("failed to create archiver: %v", err)
	}

	os.Chdir(rootDir)

	logrus.Info("Creating split archive")
	// Create tar archive
	if err := archive.CreateSplitArchive(arc, cwd, "bundle", segSize, "."); err != nil {
		return fmt.Errorf("failed to create archive: %v", err)
	}

	return nil
}

// CreateDiff performs all tasks in creating differential imagesets
//func CreateDiff(rootDir string) error {
//    return err
//}

//func downloadObjects() {
//
//}
