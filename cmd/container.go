package cmd

import (
	"github.com/puppetlabs/lumogon/analytics"
	"github.com/puppetlabs/lumogon/scheduler"
	"github.com/puppetlabs/lumogon/types"
	"github.com/spf13/cobra"
	"os"
)

var opts = types.ClientOptions{}

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan one or more containers and print the collected information",
	Long:  `Creates and attaches a container to the specified containers, inspect the container and then output that information as JSON to STDOUT`,
	PreRun: func(cmd *cobra.Command, args []string) {
		analytics.ScreenView("scan")
	},
	Run: func(cmd *cobra.Command, args []string) {
		scheduler := scheduler.New(args, opts)
		scheduler.Run()
	},
}

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Scan one or more containers and send the collected information to the Lumogon service",
	Long:  `Creates and attaches a container to the specified containers, inspect the container and then output that information as JSON over HTTP `,
	PreRun: func(cmd *cobra.Command, args []string) {
		analytics.ScreenView("report")
	},
	Run: func(cmd *cobra.Command, args []string) {
		if opts.ConsumerURL == "" {
			opts.ConsumerURL = os.Getenv("LUMOGON_ENDPOINT")
		}
		scheduler := scheduler.New(args, opts)
		scheduler.Run()
	},
}

func init() {
	RootCmd.AddCommand(scanCmd)
	RootCmd.AddCommand(reportCmd)
	RootCmd.PersistentFlags().BoolVarP(&opts.KeepHarvesters, "keep-harvesters", "k", false, "Keeps harvester containers instead of automatically deleting")
	reportFlags := reportCmd.Flags()

	reportFlags.StringVar(&opts.ConsumerURL, "endpoint", "", "Use a custom HTTP endpoint for sending the results of the scan")
}
