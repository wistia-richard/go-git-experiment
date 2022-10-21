package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/google/go-github/v48/github"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

func main() {

	message := flag.String("m", "WIP", "# commit message")
	flag.Parse()
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

	// commit the staged changes
	commit, err := w.Commit(*message, opts_commit)
	if err != nil {
		log.Fatal("Error unable to commit")
	}

	fmt.Println(r.CommitObject(commit))

	// prep for pushing changes to remote
	// headref, err := r.Head()
	// if err != nil {
	// 	log.Fatal("Error unable get the head ref of the main branch")
	// }

	// list remotes

	// list, err := r.Remotes()
	// if err != nil {
	// 	log.Fatalf("Error unable to list remotes %s", err)
	// }
	// fmt.Println(list)

	// creating a new branch

	token := os.Getenv("token")
	fmt.Println(token)

	// comment to prevent push and creating a new branch
	// if err := r.Push(&git.PushOptions{RemoteName: "origin", Auth: &http.BasicAuth{Username: "wistia-richard", Password: token}}); err != nil {
	// 	log.Fatalf("Error unable push the commit to origin %s", err)
	// }

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

	// create a pull request
	// curl \
	// -X POST \
	// -H "Accept: application/vnd.github+json" \
	// -H "Authorization: Bearer <YOUR-TOKEN>" \
	// https://api.github.com/repos/OWNER/REPO/pulls \
	// -d '{"title":"Amazing new feature","body":"Please pull these awesome changes in!","head":"octocat:new-feature","base":"master"}'

	context := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(context, ts)

	client := github.NewClient(tc)

	repos, _, _ := client.Repositories.List(context, "", &github.RepositoryListOptions{Affiliation: "owner"})
	fmt.Println(repos)
	// fmt.Println(rsp)
}
