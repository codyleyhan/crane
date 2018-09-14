# Crane

```
    Crane is a CLI for making private docker repositories actually usable
    by providing intuitive and useful commands for doing useful things

Usage:
  crane [command]

Available Commands:
  config      The saved config
  help        Help about any command
  image       Info about a docker image in repo
  ls          lists all images in the docker repo

Flags:
  -h, --help              help for crane
  -p, --password string   password for docker repo
      --profile string    profile in the config file for the docker repo
  -r, --repo string       private docker repo
  -t, --token string      token for docker repo
  -u, --unsecure          allows for accessing unsecure http repositories
  -n, --username string   username for docker repo

Use "crane [command] --help" for more information about a command.
```

## How to install

```
go get github.com/codyleyhan/crane && go install github.com/codyleyhan/crane
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

## Adding a profile

The following will walk you through creating or updating a profile
```
crane config add
``