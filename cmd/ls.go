package cmd

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/codyleyhan/crane/docker"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(lsCmd)

	lsCmd.Flags().BoolP("all", "a", false, "includes tag information in the list")
}

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls [repo]",
	Short: "lists all images in the docker repo",
	Long:  `Lists all docker images that are currently in the repo and associated information`,
	Run: func(cmd *cobra.Command, args []string) {
		if config.Repo == nil {
			fmt.Fprint(os.Stderr, "You must supply --repo=REPO or a profile")
			return
		}

		if !strings.HasPrefix(*config.Repo, "http") {
			*config.Repo = "https://" + *config.Repo
		}

		client := http.Client{Timeout: 10 * time.Second}

		if cmd.Flag("all").Value.String() == "true" {
			images, err := docker.GetAllImages(*config.Repo, &client, config.Auth)
			if err != nil {
				fmt.Println(err)
				return
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Name", "Tag"})

			for _, image := range images {
				for _, tag := range image.Tags {
					table.Append([]string{image.Name, tag})
				}
			}

			table.Render()
		} else {
			catalog, err := docker.GetCatalog(*config.Repo, &client, config.Auth)
			if err != nil {
				fmt.Println(err)
				return
			}

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Name"})

			for _, image := range catalog.Repositories {
				table.Append([]string{image})
			}

			table.Render()
		}
	},
}
