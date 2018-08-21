package cmd

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/codyleyhan/crane/docker"

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
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo := args[0]

		if !strings.HasPrefix(repo, "http") {
			repo = "https://" + repo
		}

		client := http.Client{Timeout: 10 * time.Second}

		if cmd.Flag("all").Value.String() == "true" {
			images, err := docker.GetAllImages(repo, &client)
			if err != nil {
				fmt.Println(err)
				return
			}

			for _, image := range images {
				fmt.Print(image)
			}
		} else {
			catalog, err := docker.GetCatalog(repo, &client)
			if err != nil {
				fmt.Println(err)
				return
			}

			for _, image := range catalog.Repositories {
				fmt.Println(image)
			}
		}
	},
}
