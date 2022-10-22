package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-git/go-git/v5/plumbing/transport/http"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/google/go-github/v48/github"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

type pullRequestData struct {
	prMetadata github.NewPullRequest
	owner      string
	repository string
}

func main() {

	// url := "https://github.com/wistia-richard/go-git-experiment"
	path, err := os.Getwd()
	if err != nil {
		log.Fatal("Error getting the current working directory")
	}

	fmt.Printf("The current working directory is: %s \n", path)

	r, err := git.PlainOpen(path)
	if err != nil {
		log.Fatal("Error opening the git repo")
	}

	w, err := r.Worktree()
	if err != nil {
		log.Fatal("Error getting the current working directory")
	}
	status, _ := w.Status()
	fmt.Println(status)

	// add the files
	if err := w.AddWithOptions(&git.AddOptions{All: true}); err != nil {
		fmt.Println(err)
	}

	// get git status
	status, err = w.Status()
	if err != nil {
		log.Fatal("Error getting the git status")
	}
	fmt.Println(status)

	// prep for commit
	opts_commit := &git.CommitOptions{
		All: true,
		Author: &object.Signature{
			Name:  "wistia-richard",
			Email: "rdonbosco@wistia.com",
			When:  time.Now(),
		},
	}

	message := flag.String("m", "WIP", "# commit message")
	flag.Parse()


	// commit the staged changes
	commit, err := w.Commit(*message, opts_commit)
	if err != nil {
		log.Fatal("Error unable to commit")
	}

	fmt.Println(r.CommitObject(commit))

	// // prep for pushing changes to remote
	// headref, err := r.Head()
	// if err != nil {
	// 	log.Fatal("Error unable get the head ref of the main branch")
	// }

	// list remotes

	list, err := r.Remotes()
	if err != nil {
		log.Fatalf("Error unable to list remotes %s", err)
	}
	fmt.Println(list)

	// creating a new branch

	token := os.Getenv("token")
	fmt.Println(token)

	if err := r.Push(&git.PushOptions{RemoteName: "origin", Auth: &http.BasicAuth{Username: "wistia-richard", Password: token}}); err != nil {
		log.Fatalf("Error unable push the commit to origin %s", err)
	}

	// creating a new branch
	// branchRef := plumbing.NewHashReference("refs/heads/RD/test", headref.Hash())

	// if err := r.Storer.SetReference(branchRef); err != nil {
	// 	log.Fatal("Error unable to store new ref")
	// }

	// branchConfig, err := r.Branch("RD/test")
	// if err != nil {
	// 	log.Fatal("Error unable to find the branch ")
	// }

	// fmt.Println(branchConfig.Name)

	// checkout the new branch
	// if err := w.Checkout(&git.CheckoutOptions{Branch: "refs/heads/RD/test"}); err != nil {
	// 	log.Fatalf("Error unable checkout branch %s", branchConfig.Name)
	// }

	context := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(context, ts)

	client := github.NewClient(tc)

	repos, _, _ := client.Repositories.List(context, "", &github.RepositoryListOptions{Affiliation: "owner"})

	fmt.Println("The list of repositories are:")
	for _, repo := range repos {
		fmt.Println(*repo.HTMLURL)
	}

	owner := "wistia-richard"
	repositoryName := "go-git-experiment"
	// create a PR
	prInformation := pullRequestData{
		prMetadata: github.NewPullRequest{
			Title:               github.String("Creating pr using github-go"),
			Head:                github.String("RD/test"),
			Base:                github.String("main"),
			Body:                github.String("This is the description of the PR created with the package `github.com/google/go-github/github`"),
			MaintainerCanModify: github.Bool(true),
		},
		owner:      owner,
		repository: repositoryName,
	}

	createPullRequest(client, context, &prInformation)

	/// list the PR

	prOptions := &github.PullRequestListOptions{
		State: "open",
	}

	prList, _, err := client.PullRequests.List(context, owner, repositoryName, prOptions)

	fmt.Println("The list of PRs:")
	for _, pr := range prList {
		fmt.Printf("URL:%s \t; state:%s\n", *pr.URL, *pr.State)
	}

	// close the PR

	// change the state to close
	for _, pr := range prList {
		*pr.State = "close"
	}

	for _, pr := range prList {
		fmt.Printf("URL:%s edited_state:%s", *pr.URL, *pr.State)

	}

	fmt.Println(len(prList))
	for i := len(prList) - 1; i > -1; i-- {
		fmt.Printf("Closing pr %d \n", *prList[i].Number)
		prClosed, _, _ := client.PullRequests.Edit(context, owner, repositoryName, *prList[i].Number, prList[i])
		if prClosed == nil {
			fmt.Printf("The pr %d was closed\n", *prList[i].Number)
		}
	}

	prList, _, err = client.PullRequests.List(context, owner, repositoryName, &github.PullRequestListOptions{State: "close"})

	fmt.Println("The list of PRs:")
	for _, pr := range prList {
		fmt.Printf("URL:%s \t; state:%s\n", *pr.URL, *pr.State)
	}

}
