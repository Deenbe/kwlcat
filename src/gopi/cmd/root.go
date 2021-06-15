package cmd

import (
	"github.com/spf13/cobra"
	"log"
)

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

	err := RootCmd.MarkPersistentFlagRequired("port")
	if err != nil {
		log.Fatal(err)
	}

	RootCmd.AddCommand(IDGenCmd)
	RootCmd.AddCommand(NameGenCmd)
}
