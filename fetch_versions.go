// Fetch latest kubectl, helm and AWS cli versions and print them to stdout
// Example:
// $ ./fetchversions
// KUBECTL_VERSION=v1.17.12 HELM_VERSION=v3.3.3 AWS_CLI_VERSION=2.0.50
package main

import (
	"context"
	"fmt"
	"sort"

	"github.com/Masterminds/semver"
	"github.com/google/go-github/github"
	"github.com/heroku/docker-registry-client/registry"
)

func fetchAwsCliVersion(versions chan string) {
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

	versions <- "AWS_CLI_VERSION=" + parsedTags[len(parsedTags)-1].Original()
}

func fetchKubeVersion(client *github.Client, versions chan string) {
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
		if *release.Prerelease {
			continue
		}

		version := semver.MustParse(*release.TagName)

		// pin kubectl version to 21 because of AWS EKS suppored
		// kubernetes versions
		// https://docs.aws.amazon.com/eks/latest/userguide/kubernetes-versions.html
		if version.Major() == 1 && version.Minor() <= 21 {
			newReleases = append(newReleases, version)
		}
	}

	sort.Sort(newReleases)

	versions <- "KUBECTL_VERSION=" + newReleases[len(newReleases)-1].Original()
}

func fetchHelmVersion(client *github.Client, versions chan string) {
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
	versions <- "HELM_VERSION=" + newReleases[len(newReleases)-1].Original()
}

func main() {
	versions := make(chan string)

	go fetchAwsCliVersion(versions)

	githubClient := github.NewClient(nil)
	go fetchKubeVersion(githubClient, versions)
	go fetchHelmVersion(githubClient, versions)

	for i := 0; i < 3; i++ {
		version := <-versions
		fmt.Printf("%s ", version)
	}

	fmt.Println()
}
