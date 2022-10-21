package main

import (
	"context"
	"fmt"

	"github.com/google/go-github/v48/github"
)

func createPullRequest(client *github.Client, context context.Context, prInformation *pullRequestData) {
	newPR := prInformation.prMetadata

	pr, _, err := client.PullRequests.Create(context, prInformation.owner, prInformation.repository, &newPR)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("PR created: %s\n", pr.GetHTMLURL())
}
