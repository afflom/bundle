package bundle

import "strings"

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
