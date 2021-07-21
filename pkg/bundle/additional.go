package bundle

import (
	"fmt"
	"os"

	"github.com/openshift/oc/pkg/cli/image/imagesource"
	"github.com/openshift/oc/pkg/cli/image/mirror"
	"github.com/sirupsen/logrus"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

// GetAdditional downloads specified images in the imageset-config.yaml under mirror.additonalImages
func GetAdditional(i *BundleSpec, rootDir string) error {

	var mappings []mirror.Mapping

	stream := genericclioptions.IOStreams{
		In:     os.Stdin,
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}

	// TODO: Add credential options
	opts := mirror.NewMirrorImageOptions(stream)
	opts.FileDir = rootDir + "/src/"

	logrus.Infof("Downloading %d image(s) to %s", len(i.Mirror.AdditionalImages), opts.FileDir)

	for _, img := range i.Mirror.AdditionalImages {

		// Get source image information
		srcRef, err := imagesource.ParseReference(img)

		if err != nil {
			return fmt.Errorf("error parsing source image %s: %v", img, err)
		}

		// Set destination image information
		path := "file://" + img

		dstRef, err := imagesource.ParseReference(path)

		if err != nil {
			return fmt.Errorf("error parsing destination reference %s: %v", path, err)
		}

		// Create mapping from source and destination images
		mappings = append(mappings, mirror.Mapping{
			Source:      srcRef,
			Destination: dstRef,
			Name:        srcRef.Ref.Name,
		})

	}

	opts.Mappings = mappings

	err := opts.Run()

	return err
}
