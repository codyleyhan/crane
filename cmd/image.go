package cmd

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/codyleyhan/crane/docker"
)

func init() {
	rootCmd.AddCommand(imageCmd)

	imageCmd.Flags().BoolP("all", "a", false, "includes all information for image")
	imageCmd.Flags().String("tag", "", "a specific tag for an image")
}

// imageCmd represents the image command
var imageCmd = &cobra.Command{
	Use:   "image [image]",
	Short: "Info about a docker image in repo",
	Long: `image determines information around docker images such as
	tags, manifests and deleting docker images`,
	ArgAliases: []string{"repo", "image"},
	Args:       cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if config.Repo == nil {
			fmt.Fprint(os.Stderr, "You must supply --repo=REPO or a profile")
			return
		}
		image := args[0]

		if !strings.HasPrefix(*config.Repo, "http") {
			*config.Repo = "https://" + *config.Repo
		}

		client := http.Client{Timeout: 10 * time.Second}

		if cmd.Flag("tag").Value.String() != "" {
			manifest, err := docker.GetImageManifest(*config.Repo, image, cmd.Flag("tag").Value.String(), &client, config.Auth)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Print(manifest)
			return
		}

		i, err := docker.GetImage(*config.Repo, image, &client, config.Auth)
		if err != nil {
			fmt.Println(err)
			return
		}

		table := tablewriter.NewWriter(os.Stdout)

		if cmd.Flag("all").Value.String() == "true" {
			table.SetHeader([]string{"tag", "digest", "size", "media type"})

			for _, tag := range i.Tags {
				manifest, err := docker.GetImageManifest(*config.Repo, image, tag, &client, config.Auth)
				if err != nil {
					fmt.Println(err)
					return
				}

				table.Append([]string{
					tag,
					manifest.Config.Digest,
					strconv.FormatInt(int64(manifest.Config.Size), 10),
					manifest.Config.MediaType},
				)
			}
		} else {
			table.SetHeader([]string{"tag"})
			for _, tag := range i.Tags {
				table.Append([]string{tag})
			}
		}

		table.Render()
	},
}
