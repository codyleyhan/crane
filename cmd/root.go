package cmd

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"

	"github.com/codyleyhan/crane/docker"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	username string
	password string
	token    string

	auth docker.Auth
)

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&username, "username", "n", "", "username for docker repo")
	rootCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "password for docker repo")
	rootCmd.PersistentFlags().StringVarP(&token, "token", "t", "", "token for docker repo")
	rootCmd.PersistentFlags().BoolP("unsecure", "u", false, "allows for accessing unsecure http repositories")
}

func initConfig() {
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	viper.AddConfigPath(home)
	viper.SetConfigName(".crane")

	if username != "" {
		auth.Username = &username
	}
	if password != "" {
		auth.Password = &password
	}
	if token != "" {
		auth.Token = &token
	}

	viper.ReadInConfig()
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
