package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
)

func init() {
	configCmd.AddCommand(addCmd)
}

// configCmd represents the config command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a profile",
	Long: `Add a profile so that you can move betwen multiple
	docker repositiories with ease without having to remember all of 
	the config for each one
	`,
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)

		profileMap := make(map[string]interface{})
		name := ""

		state := 0
		finished := false
		for !finished {
			switch state {
			case 0:
				{
					fmt.Print("What is the profile name?\n\xF0\x9F\x9A\x80 ")
					input, err := reader.ReadString('\n')
					if err != nil {
						fmt.Fprintln(os.Stderr, err)
						return
					}
					input = strings.TrimSpace(input)
					if input == "" {
						break
					}
					name = input
					state++
					break
				}
			case 1:
				{
					fmt.Print("What is the url for the Repo?\n\xF0\x9F\x9A\x80 ")
					input, err := reader.ReadString('\n')
					if err != nil {
						fmt.Fprintln(os.Stderr, err)
						return
					}
					input = strings.TrimSpace(input)
					if input == "" {
						break
					}
					profileMap["repo"] = input
					state++
					break
				}
			case 2:
				{
					fmt.Print("Is this repo unsecure ie it does not have a verified https certificate? 'yes' or 'no'\n\xF0\x9F\x9A\x80 ")
					input, err := reader.ReadString('\n')
					if err != nil {
						fmt.Fprintln(os.Stderr, err)
						return
					}
					input = strings.TrimSpace(input)
					input = strings.ToLower(input)
					if input != "yes" {
						state++
					}
					profileMap["unsecure"] = true
					state++
					break
				}
			case 3:
				{
					fmt.Print("Is that a username for this repo? Please press return or enter 'no' if none\n\xF0\x9F\x9A\x80 ")
					input, err := reader.ReadString('\n')
					if err != nil {
						fmt.Fprintln(os.Stderr, err)
						return
					}
					input = strings.TrimSpace(input)
					if input == "" || input == "no" {
						state += 2
						break
					}
					profileMap["username"] = input
					state++
					break
				}
			case 4:
				{
					fmt.Print("What is the password for this username?\n\xF0\x9F\x9A\x80 ")
					inputBytes, err := terminal.ReadPassword(0)
					if err != nil {
						fmt.Fprintln(os.Stderr, err)
						return
					}
					input := string(inputBytes)
					if input == "" {
						break
					}
					profileMap["password"] = input
					state += 2

					break
				}
			case 5:
				{
					fmt.Print("Do you have a token for this repo? Please press return or enter 'no' if none\n\xF0\x9F\x9A\x80 ")
					input, err := reader.ReadString('\n')
					input = strings.TrimSpace(input)
					if err != nil {
						fmt.Fprintln(os.Stderr, err)
						return
					}
					if input == "" || input == "no" {
						state++
						break
					}
					profileMap["token"] = input
					state++
					break
				}
			case 6:
				{
					fmt.Println("Saving profile to config file")
					profiles := viper.GetStringMap("profiles")
					if profiles == nil {
						profiles = make(map[string]interface{})
					}

					profiles[name] = profileMap
					viper.Set("profiles", profiles)
					if err := viper.WriteConfig(); err != nil {
						fmt.Print("Could not write profile because %+v\n", err)
						finished = true
					}
					fmt.Println("Successfully wrote profile to config")
					finished = true
				}
			}
		}
	},
}
