package cmd

import "github.com/spf13/cobra"

var (
	host string
	port int
	RootCmd = &cobra.Command{
		Use:  "oteldemo",
		Long: "oteldemo cli",
	}
)

func init() {
	RootCmd.PersistentFlags().StringVar(&host, "host", "0.0.0.0", "host address to bind to")
	RootCmd.PersistentFlags().IntVar(&port, "port", port, "host port to bind to")
	RootCmd.AddCommand(IDGenCmd)
	RootCmd.AddCommand(NameGenCmd)
}
