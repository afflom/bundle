package create

import (
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
	//}
	//if config.Mirror.AdditionalImages != nil {
	//	getAdditional(*config, rootDir)
	//}
	*/

	// Create archiver
	arc, err := archive.NewArchiver(ext)

	if err != nil {
		logrus.Errorf("failed to create archiver: %v", err)
		return err
	}

	cwd, err := os.Getwd()

	if err != nil {
		return err
	}

	os.Chdir(rootDir)

	logrus.Info("Creating split archive")
	// Create tar archive
	if err := archive.CreateSplitArchive(arc, cwd, "bundle", segSize, "."); err != nil {
		logrus.Errorf("failed to create archive: %v", err)
		return err
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
