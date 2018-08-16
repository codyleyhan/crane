package cmd

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.PersistentFlags().BoolP("unsecure", "u", false, "allows for accessing unsecure http repositories")
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "crane",
	Short: "A utility for working with private docker repositories",
	Long: `
	Crane is a CLI for making private docker repositories actually usable
	by providing intuitive and useful commands for doing useful things`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if cmd.Flag("unsecure").Value.String() == "true" {
			http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
