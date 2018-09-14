package docker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/hashicorp/go-version"
	"github.com/pkg/errors"
)

const (
	catalogURL   = "/v2/_catalog"
	imageURL     = "/v2/%s/tags/list"
	manifestURL  = "/v2/%s/manifests/%s"
	dockerHeader = "application/vnd.docker.distribution.manifest.v2+json"
)

type (
	Catalog struct {
		Repositories []string `json:"repositories"`
	}

	Image struct {
		Name string   `json:"name"`
		Tags []string `json:"tags"`
	}

	Digest struct {
		MediaType string
		Size      int
		Digest    string
	}

	Manifest struct {
		SchemaVersion int
		Config        Digest
		Layers        []Digest
	}

	Auth struct {
		Username *string
		Password *string
		Token    *string
	}
)

func (i Image) String() string {
	str := strings.Builder{}
	str.WriteString(color.CyanString(i.Name + "\n"))
	for _, tag := range i.Tags {
		str.WriteString(color.WhiteString("- " + tag + "\n"))
	}

	return str.String()
}

func (m Manifest) String() string {
	str := strings.Builder{}
	str.WriteString(color.WhiteString("- Schema Version: %d\n-", m.SchemaVersion))
	str.WriteString(color.RedString(" Digest: "))
	str.WriteString(color.WhiteString("%s\n", m.Config.Digest))
	str.WriteString(color.WhiteString("- Media Type: %s\n", m.Config.MediaType))
	str.WriteString(color.WhiteString("- Size: %d\n", m.Config.Size))

	str.WriteString(color.WhiteString("- Layers:\n"))
	for _, layer := range m.Layers {
		str.WriteString(color.WhiteString("  -"))
		str.WriteString(color.RedString(" Digest: "))
		str.WriteString(color.WhiteString("%s\n", m.Config.Digest))
		str.WriteString(color.WhiteString("  - Media Type: %s\n", layer.MediaType))
		str.WriteString(color.WhiteString("  - Size: %d\n", layer.Size))
	}

	return str.String()
}

func createGetRequest(url string, auth Auth) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	if auth.Username != nil && auth.Password != nil {
		req.SetBasicAuth(*auth.Username, *auth.Password)
	}
	if auth.Token != nil {
		req.Header.Add("Authorization", "Basic "+*auth.Token)
	}

	return req, nil
}

func GetCatalog(repoURL string, client *http.Client, auth Auth) (*Catalog, error) {
	var catalog Catalog

	req, err := createGetRequest(repoURL+catalogURL, auth)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create request to get image manifeast")
	}

	catalogBody, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "could not get list of images in repo")
	}

	decoder := json.NewDecoder(catalogBody.Body)
	if err := decoder.Decode(&catalog); err != nil {
		return nil, errors.Wrap(err, "there was a problem decoding the catalog")
	}

	return &catalog, nil
}

// GetAllImages returns all images in the repo
func GetAllImages(repoURL string, client *http.Client, auth Auth) ([]*Image, error) {
	catalog, err := GetCatalog(repoURL, client, auth)
	if err != nil {
		return nil, err
	}

	results := make(chan *Image, len(catalog.Repositories))
	wg := sync.WaitGroup{}
	wg.Add(len(catalog.Repositories))

	for _, repo := range catalog.Repositories {
		repo := repo
		newClient := *client
		go func() {
			image, err := GetImage(repoURL, repo, &newClient, auth)
			if err != nil {
				fmt.Println(err)
				panic(err)
			}
			results <- image
			wg.Done()
		}()
	}

	wg.Wait()
	close(results)

	images := make([]*Image, 0, len(catalog.Repositories))
	for result := range results {
		images = append(images, result)
	}

	return images, nil
}

func GetImage(repoURL, name string, client *http.Client, auth Auth) (*Image, error) {
	var image Image
	url := fmt.Sprintf(imageURL, name)

	req, err := createGetRequest(repoURL+url, auth)
	if err != nil {
		return nil, errors.Wrap(err, "problem getting tags for "+name)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "problem getting tags for "+name)
	}

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&image); err != nil {
		return nil, errors.Wrap(err, "problem decoding "+name)
	}

	// Sort by semantic versioning
	versions := []*version.Version{}
	notSemVer := []string{}
	for _, tag := range image.Tags {
		v, err := version.NewVersion(tag)
		if err != nil {
			notSemVer = append(notSemVer, tag)
			continue
		}
		versions = append(versions, v)
	}
	sort.Sort(version.Collection(versions))
	tags := make([]string, 0, len(image.Tags))
	for _, tag := range versions {
		tags = append(tags, tag.String())
	}
	image.Tags = append(tags, notSemVer...)

	return &image, nil
}

func GetImageManifest(repoURL, image, tag string, client *http.Client, auth Auth) (*Manifest, error) {
	var manifest Manifest

	req, err := createGetRequest(fmt.Sprintf(repoURL+manifestURL, image, tag), auth)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create request to get image manifeast")
	}

	req.Header.Add("Accept", dockerHeader)

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "could not get manifest of image")
	}

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&manifest); err != nil {
		return nil, errors.Wrap(err, "problem decoding "+image)
	}

	return &manifest, nil
}
