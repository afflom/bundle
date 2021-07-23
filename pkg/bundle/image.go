package bundle

import (
	"context"
	"fmt"
	"strings"

	"github.com/containers/image/v5/transports/alltransports"
	"github.com/containers/image/v5/types"
)

// IsBlocked will return a boolean value on whether an image
// is specified as blocked in the BundleSpec
func IsBlocked(config *BundleSpec, imgRef string) bool {

	// TODO: parse reference and break it down by actual image name
	// to ensure we don't have false positives
	for _, block := range config.BlockedImages {
		if strings.Contains(imgRef, block) {
			return true
		}
	}
	return false
}

// InspectImages will inspect an image to determine the base image
func InspectImages(ctx context.Context, name string) ([]byte, error) {
	ref, err := alltransports.ParseImageName(name)
	if err != nil {
		return nil, err
	}
	sys := &types.SystemContext{}
	src, err := ref.NewImageSource(ctx, sys)

	if err != nil {
		return nil, fmt.Errorf("error creating image source: %v", err)
	}

	manifest, _, err := src.GetManifest(ctx, nil)

	return manifest, err
}
