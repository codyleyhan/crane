package cmd

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"

	"github.com/codyleyhan/crane/docker"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	username string
	password string
	token    string
	profile  string
	repo     string

	config struct {
		docker.Auth `mapstructure:",squash"`
		Profile     *string
		Repo        *string
		Unsecure    bool
	}
)

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().BoolVarP(&config.Unsecure, "unsecure", "u", false, "allows for accessing unsecure http repositories")
	rootCmd.PersistentFlags().StringVarP(&username, "username", "n", "", "username for docker repo")
	rootCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "password for docker repo")
	rootCmd.PersistentFlags().StringVarP(&token, "token", "t", "", "token for docker repo")
	rootCmd.PersistentFlags().StringVarP(&repo, "repo", "r", "", "private docker repo")
	rootCmd.PersistentFlags().StringVar(&profile, "profile", "", "profile in the config file for the docker repo")
}

func filledOrNil(src string, dest **string) {
	if src != "" {
		*dest = &src
	}
}

func initConfig() {
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	viper.SetConfigType("json")
	viper.AddConfigPath(home + "/.crane")
	viper.SetConfigName("config")
	viper.ReadInConfig()

	filledOrNil(username, &config.Username)
	filledOrNil(password, &config.Password)
	filledOrNil(token, &config.Token)
	filledOrNil(profile, &config.Profile)
	filledOrNil(repo, &config.Repo)

	if config.Profile != nil {
		profiles := viper.GetStringMap("profiles")
		if profiles == nil {
			fmt.Fprintf(os.Stderr, "There is no profile named %s\n", *config.Profile)
		}

		data, ok := profiles[*config.Profile]
		if !ok {
			fmt.Fprintf(os.Stderr, "There is no profile named %s\n", *config.Profile)
		}
		if err := mapstructure.Decode(data, &config); err != nil {
			fmt.Fprintf(os.Stderr, "There is no profile named %s: %+v\n", *config.Profile, err)
		}
	}
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "crane",
	Short: "A utility for working with private docker repositories",
	Long: `
	Crane is a CLI for making private docker repositories actually usable
	by providing intuitive and useful commands for doing useful things`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if config.Unsecure {
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
