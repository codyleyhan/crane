# Crane

```
	Crane is a CLI for making private docker repositories actually usable
	by providing intuitive and useful commands for doing useful things like listing
    images and their tags

Usage:
  crane [command]

Available Commands:
  help        Help about any command
  image       Info about a docker image in repo
  ls          lists all images in the docker repo

Flags:
  -h, --help       help for crane
  -u, --unsecure   allows for accessing unsecure http repositories

Use "crane [command] --help" for more information about a command.
```

## How to install

```
go install github.com/codyleyhan/crane
```

## Example usage

Showing all docker images in a repository
```
crane ls repoURL
```

Showing all docker images with tags in a repository

```
crane ls -a repoURL
```

Showing metadata for a given image 

```
crane image -a repoURL image
```