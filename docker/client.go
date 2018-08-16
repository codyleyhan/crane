package docker

import (
	"encoding/json"
	"fmt"
	"net/http"

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
)

func (i Image) String() string {
	str := i.Name + "\n"
	for _, tag := range i.Tags {
		str += "\t- " + tag + "\n"
	}

	return str
}

func (m Manifest) String() string {
	str := ""
	str += fmt.Sprintf("- Schema Version: %d\n", m.SchemaVersion)
	str += fmt.Sprintln("- Digest:", m.Config.Digest)
	str += fmt.Sprintln("- Media Type:", m.Config.MediaType)
	str += fmt.Sprintf("- Size: %d\n", m.Config.Size)

	str += fmt.Sprintln("- Layers:")
	for _, layer := range m.Layers {
		str += fmt.Sprintln("\t- Digest:", layer.Digest)
		str += fmt.Sprintln("\t- Media Type:", layer.MediaType)
		str += fmt.Sprintf("\t- Size: %d\n", layer.Size)
	}

	return str
}

func GetCatalog(repoURL string, client *http.Client) (*Catalog, error) {
	var catalog Catalog

	catalogBody, err := client.Get(repoURL + catalogURL)
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
func GetAllImages(repoURL string, client *http.Client) ([]*Image, error) {
	catalog, err := GetCatalog(repoURL, client)
	if err != nil {
		return nil, err
	}

	images := make([]*Image, 0, len(catalog.Repositories))

	for _, repo := range catalog.Repositories {
		image, err := GetImage(repoURL, repo, client)
		if err != nil {
			return nil, err
		}

		images = append(images, image)
	}

	return images, nil
}

func GetImage(repoURL, name string, client *http.Client) (*Image, error) {
	var image Image
	url := fmt.Sprintf(imageURL, name)

	resp, err := client.Get(repoURL + url)
	if err != nil {
		return nil, errors.Wrap(err, "problem getting tags for "+name)
	}

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&image); err != nil {
		return nil, errors.Wrap(err, "problem decoding "+name)
	}

	return &image, nil
}

func GetImageManifest(repoURL, image, tag string, client *http.Client) (*Manifest, error) {
	var manifest Manifest

	req, err := http.NewRequest("GET", fmt.Sprintf(repoURL+manifestURL, image, tag), nil)
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
