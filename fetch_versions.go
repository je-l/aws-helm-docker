package main

import (
	"context"
	"fmt"
	"sort"

	"github.com/Masterminds/semver"
	"github.com/google/go-github/github"
	"github.com/heroku/docker-registry-client/registry"
)

func fetchAwsCliVersion() *semver.Version {
	docker, err := registry.New("https://registry-1.docker.io", "", "")

	if err != nil {
		panic(err)
	}

	tags, err := docker.Tags("amazon/aws-cli")

	if err != nil {
		panic(err)
	}

	var parsedTags semver.Collection = make([]*semver.Version, 0)

	for _, tag := range tags {
		if tag == "amd64" || tag == "arm64" || tag == "latest" {
			continue
		}

		parsedTags = append(parsedTags, semver.MustParse(tag))
	}

	sort.Sort(parsedTags)

	return parsedTags[len(parsedTags)-1]
}

func fetchKubectlVersion(client *github.Client) *semver.Version {
	releases, _, error := client.Repositories.ListReleases(
		context.Background(),
		"kubernetes",
		"kubernetes",
		nil,
	)

	if error != nil {
		panic(error)
	}

	var newReleases semver.Collection = make([]*semver.Version, 0)

	for _, release := range releases {
		version := semver.MustParse(release.GetName())

		// pin kubectl version to 17 because of AWS EKS suppored
		// kubernetes versions
		if version.Prerelease() == "" && version.Minor() <= 17 {
			newReleases = append(newReleases, version)
		}
	}

	sort.Sort(newReleases)

	return newReleases[len(newReleases)-1]
}

func fetchHelmVersion(client *github.Client) *semver.Version {
	releases, _, error := client.Repositories.ListReleases(
		context.Background(),
		"helm",
		"helm",
		nil,
	)

	if error != nil {
		panic(error)
	}

	var newReleases semver.Collection = make([]*semver.Version, 0)

	for _, release := range releases {
		version := semver.MustParse(*release.TagName)

		if version.Prerelease() == "" {
			newReleases = append(newReleases, version)
		}
	}

	sort.Sort(newReleases)

	return newReleases[len(newReleases)-1]
}

func main() {
	awsCliVersion := fetchAwsCliVersion()

	githubClient := github.NewClient(nil)
	kubectlVersion := fetchKubectlVersion(githubClient)
	helmVersion := fetchHelmVersion(githubClient)

	fmt.Printf("AWS_CLI_VERSION=%s HELM_VERSION=%s KUBECTL_VERSION=%s", awsCliVersion, helmVersion, kubectlVersion)
}
