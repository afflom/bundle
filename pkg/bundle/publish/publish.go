package publish

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/RedHatGov/bundle/pkg/archive"
	"github.com/RedHatGov/bundle/pkg/config"
	"github.com/RedHatGov/bundle/pkg/config/v1alpha1"
	"github.com/RedHatGov/bundle/pkg/metadata/storage"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var (
	PublishOpts struct {
		FromBundle string
		ToMirror   string
	}
)

type UuidError struct {
	InUuid     uuid.UUID
	CurentUuid uuid.UUID
}

func (u *UuidError) Error() string {
	return fmt.Sprintf("Mismatched UUIDs. Want %v, got %v", u.CurentUuid, u.InUuid)
}

func isMultiple(file fs.FileInfo, a archive.Archiver, imageset, target, dest string) {
	if file.IsDir() {

		// find first file and load metadata from that
		logrus.Infoln("Detected multiple incoming archive files")

	} else {
		// Pass file to ExtractFile Function
		if err := archive.ExtractFile(a, imageset, target, dest); err != nil {
			logrus.Error(err)
		}
	}
}

func Publish(rootDir, imageset, fqdn string, dryRun, insecure bool) error {

	ctx := context.Background()

	// Create Publish dirs
	// if err := bundle.MakePublishDirs(rootDir); err != nil {
	// 	logrus.Error(err)
	// }

	// Create backend
	backend, err := storage.NewLocalBackend(rootDir)
	if err != nil {
		return fmt.Errorf("error opening local backend: %v", err)
	}

	// Copy metadata to publish dir
	a := archive.NewArchiver()

	file, err := os.Stat(imageset)

	currentMeta := v1alpha1.Metadata{}
	incomingMeta := v1alpha1.Metadata{}
	target := filepath.Join(config.PublishDir, config.MetadataFile)

	tempMeta := config.MetadataFile + ".incoming"

	if err != nil {
		logrus.Error(err)
	}

	if _, err := os.Stat(target); os.IsNotExist(err) {
		isMultiple(file, a, imageset, target, target)

		// find first file and load metadata from that
		logrus.Infof("No existing metadata found. Setting up new workspace: %v ", err)
		backend.ReadMetadata(ctx, &incomingMeta, tempMeta)

	} else {
		// Pass file to ExtractFile Function
		isMultiple(file, a, imageset, target, tempMeta)

	}
	// extract metadata
	if err := backend.ReadMetadata(ctx, &incomingMeta, tempMeta); err != nil {
		logrus.Error(err)
	}
	backend.ReadMetadata(ctx, &currentMeta, tempMeta)

	if incomingMeta.MetadataSpec.Uid != currentMeta.MetadataSpec.Uid {
		logrus.Fatalln("Refereneced imageset has mismatched UUID for the target workspace. UUID must match")
		return &UuidError{incomingMeta.MetadataSpec.Uid, currentMeta.MetadataSpec.Uid}
	}

	// check sequence/uid
	// If there was previously imported metadata:
	// it must have the same uuid and be in sequence with the incoming metadata
	// If new workspace, dont validate uuid or sequence number

	// if multiple files, combine

	// import containers

	// import imagecontentsourcepolicy

	// import catalogsource

	// mirror to registry

	// install imagecontentsourcepolicy

	// install catalogsource

	return nil
}
