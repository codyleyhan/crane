package cmd

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/codyleyhan/crane/docker"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(imageCmd)

	imageCmd.Flags().BoolP("all", "a", false, "includes all information for image")
	imageCmd.Flags().StringP("tag", "t", "", "a specific tag for an image")
}

// imageCmd represents the image command
var imageCmd = &cobra.Command{
	Use:   "image [repo] [image]",
	Short: "Info about a docker image in repo",
	Long: `image determines information around docker images such as
	tags, manifests and deleting docker images`,
	ArgAliases: []string{"repo", "image"},
	Args:       cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		repo := args[0]
		image := args[1]

		if !strings.HasPrefix(repo, "http") {
			repo = "https://" + repo
		}

		client := http.Client{Timeout: 10 * time.Second}

		if cmd.Flag("tag").Value.String() != "" {
			manifest, err := docker.GetImageManifest(repo, image, cmd.Flag("tag").Value.String(), &client)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Print(manifest)
			return
		}

		i, err := docker.GetImage(repo, image, &client)
		if err != nil {
			fmt.Println(err)
			return
		}

		if cmd.Flag("all").Value.String() == "true" {
			for _, tag := range i.Tags {
				color.Cyan(tag)
				manifest, err := docker.GetImageManifest(repo, image, tag, &client)
				if err != nil {
					fmt.Println(err)
					return
				}

				fmt.Print(manifest)
			}

		} else {
			fmt.Print(i)
		}
	},
}
