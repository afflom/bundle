package catalog

import (
	"github.com/RedHatGov/bundle/pkg/bundle"
	"github.com/blang/semver"
	"github.com/sirupsen/logrus"
)

const (
	redhatoperators    = "registry.redhat.io/redhat/redhat-operator-index"
	communityoperators = "quay.io/openshift-community-operators/catalog"
	certifiedoperators = ""
)

func PullChannelHeads() {}

func GetChannelImageName() {}

func GetOperators(i *bundle.Imageset, c *bundle.BundleSpec, rootDir string) error {
	// Check for metadata
	if i != nil {
		// For each Catalog in the config file
		for _, r := range c.Mirror.Operators {
			// Check for specific version declarations
			if r.Packages != nil {
				// for each package, determin if latest or specific version
				for _, rn := range r.Packages {
					// Check for declared versions
					if rn.Version != "" {
						// Download each version
						for _, rv := range rn.Version {
							// Convert the string to a semver
							logrus.Infof("rn is: %v", rv)
							rs, err := semver.Parse(rn.Version)
							if err != nil {
								logrus.Errorln(err)
								return err
							}
							// This dumps the available upgrades from the last downloaded version
							requested, _, err := calculateUpgradePath(r, rs)
							if err != nil {
								logrus.Errorln("Failed get upgrade graph")
								logrus.Error(err)
								return err
							}

							logrus.Infof("requested: %v", requested.Version)
							err = downloadMirror(requested.Image, rootDir)
							if err != nil {
								logrus.Errorln(err)
							}
							logrus.Infof("Channel Latest version %v", requested.Version)

							/* Select the requested version from the available versions
							for _, d := range neededVersions {
								logrus.Infof("Available Release Version: %v \n Requested Version: %o", d.Version, rs)
								if d.Version.Equals(rs) {
									logrus.Infof("Image to download: %v", d.Image)
									err := downloadMirror(d.Image)
									if err != nil {
										logrus.Errorln(err)
									}
									logrus.Infof("Image to download: %v", d.Image)
									break
								}
							} */

							// download the selected version

							logrus.Infof("Current Object: %v", rn)
							logrus.Infoln("")
							logrus.Infoln("")
							//logrus.Infof("Next-Versions: %v", neededVersions.)
							//nv = append(nv, neededVersions)
						}
					} else {
						// Download each declared operator head

					}
				}
			} else {
				// Download all heads within the catalog
				latest, err := GetLatestVersion(r)
				if err != nil {
					logrus.Errorln(err)
					return err
				}
				logrus.Infof("Image to download: %v", latest.Image)
				// Download the release
				err = downloadMirror(latest.Image, rootDir)
				if err != nil {
					logrus.Errorln(err)
				}
				logrus.Infof("Channel Latest version %v", latest.Version)
			}
		}
		// If there was no metadata
	} else {
		// For each Catalog in the config file
		for _, r := range c.Mirror.Operators {
			// Check for specific version declarations
			if r.Packages != nil {
				// for each package, determin if latest or specific version
				for _, rn := range r.Packages {
					// Check for declared versions
					if rn.Version != nil {
						// Download each version
						for _, rv := range rn.Version {
							// Convert the string to a semver
							logrus.Infof("rn is: %v", rn)
							rs, err := semver.Parse(rn)
							if err != nil {
								logrus.Errorln(err)
								return err
							}
							// This dumps the available upgrades from the last downloaded version
							requested, _, err := calculateUpgradePath(r, rs)
							if err != nil {
								logrus.Errorln("Failed get upgrade graph")
								logrus.Error(err)
								return err
							}

							logrus.Infof("requested: %v", requested.Version)
							err = downloadMirror(requested.Image, rootDir)
							if err != nil {
								logrus.Errorln(err)
							}
							logrus.Infof("Channel Latest version %v", requested.Version)

							/* Select the requested version from the available versions
							for _, d := range neededVersions {
								logrus.Infof("Available Release Version: %v \n Requested Version: %o", d.Version, rs)
								if d.Version.Equals(rs) {
									logrus.Infof("Image to download: %v", d.Image)
									err := downloadMirror(d.Image)
									if err != nil {
										logrus.Errorln(err)
									}
									logrus.Infof("Image to download: %v", d.Image)
									break
								}
							} */

							// download the selected version

							logrus.Infof("Current Object: %v", rn)
							logrus.Infoln("")
							logrus.Infoln("")
							//logrus.Infof("Next-Versions: %v", neededVersions.)
							//nv = append(nv, neededVersions)
						}
					} else {
						// Download each declared operator head

					}
				}
			} else {
				// Download all heads within the catalog
				latest, err := GetLatestVersion(r)
				if err != nil {
					logrus.Errorln(err)
					return err
				}
				logrus.Infof("Image to download: %v", latest.Image)
				// Download the release
				err = downloadMirror(latest.Image, rootDir)
				if err != nil {
					logrus.Errorln(err)
				}
				logrus.Infof("Channel Latest version %v", latest.Version)
			}
		}
	}
	return nil
}

// if operators are populated, start processing

// If entire catalog, download heads of each operator

// if
