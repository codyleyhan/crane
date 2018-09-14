package cmd

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(configCmd)
}

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "The saved config",
	Long: `The config used for this tool, save at ~/.crane. 
	Look at the documentation for all the settings that can be set
	`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%+v\n", viper.AllSettings())
		fmt.Printf("%+v\n", config)
	},
}
