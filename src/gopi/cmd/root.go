package cmd

import "github.com/spf13/cobra"

var (
	host string
	port int
	profile bool
	RootCmd = &cobra.Command{
		Use:  "gopi",
		Long: "kwlcat core api",
	}
)

func init() {
	RootCmd.PersistentFlags().StringVar(&host, "host", "0.0.0.0", "host address to bind to")
	RootCmd.PersistentFlags().IntVar(&port, "port", port, "host port to bind to")
	RootCmd.PersistentFlags().BoolVar(&profile, "profile", false, "enable profiler")
	RootCmd.MarkPersistentFlagRequired("port")
	RootCmd.AddCommand(IDGenCmd)
	RootCmd.AddCommand(NameGenCmd)
}
