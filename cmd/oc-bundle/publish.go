package main

import (
	"github.com/RedHatGov/bundle/pkg/bundle/publish"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type publishOpts struct {
	imagesetLocation string
	targetRegistry   string
}

func newPublishCmd() *cobra.Command {

	opts := publishOpts{}

	return &cobra.Command{
		Use:   "publish",
		Short: "Publish OCP related content to an internet-disconnected environment",
		Args:  cobra.ExactArgs(0),
		Run: func(_ *cobra.Command, _ []string) {
			cleanup := setupFileHook(rootOpts.dir)
			defer cleanup()
			err := publish.Publish(rootOpts.dir, opts.imagesetLocation, opts.targetRegistry, rootOpts.dryRun, rootOpts.skipTLS)
			logrus.Infoln("Publish Was called")
			if err != nil {
				logrus.Fatal(err)
			}
		},
	}
}
