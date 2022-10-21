package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
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

	if err := w.AddWithOptions(&git.AddOptions{All: true}); err != nil {
		fmt.Println(err)
	}

	status, err = w.Status()
	if err != nil {
		log.Fatal("Error getting the git status")
	}
	fmt.Println(status)

	opts_commit := &git.CommitOptions{
		All: true,
		Author: &object.Signature{
			Name:  "wistia-richard",
			Email: "rdonbosco@wistia.com",
			When:  time.Now(),
		},
	}

	commit, err := w.Commit(*message, opts_commit)
	if err != nil {
		log.Fatal("Error unable to commit")
	}

	fmt.Println(r.CommitObject(commit))

	headref, err := r.Head()
	if err != nil {
		log.Fatal("Error unable to commit")
	}

	// list remotes

	list, err := r.Remotes()
	if err != nil {
		log.Fatalf("Error unable to list remotes %s", err)
	}
	fmt.Println(list)

	// push commit

	token := os.Getenv("token")
	// password := os.Getenv("git_password")
	// fmt.Println(password)
	// if err := r.Push(&git.PushOptions{RemoteName: "origin", Auth: &http.BasicAuth{Username: "wistia-richard", Password: password}}); err != nil {
	// 	log.Fatalf("Error unable push the commit to origin %s", err)
	// }
	fmt.Println(token)
	if err := r.Push(&git.PushOptions{RemoteName: "origin", Auth: &http.TokenAuth{Token: token}}); err != nil {
		log.Fatalf("Error unable push the commit to origin %s", err)
	}
	// token := os.Getenv("token")
	// fmt.Println(token)
	// if err := r.Push(&git.PushOptions{RemoteName: "origin", Auth: &http.TokenAuth{Token: token}}); err != nil {
	// 	log.Fatalf("Error unable push the commit to origin %s", err)
	// }

	branchRef := plumbing.NewHashReference("refs/heads/RD/test", headref.Hash())

	if err := r.Storer.SetReference(branchRef); err != nil {
		log.Fatal("Error unable to store new ref")
	}

	branchConfig, err := r.Branch("RD/test")
	if err != nil {
		log.Fatal("Error unable to find the branch ")
	}

	fmt.Println(branchConfig.Name)

	if err := w.Checkout(&git.CheckoutOptions{Branch: "refs/heads/RD/test"}); err != nil {
		log.Fatalf("Error unable checkout branch %s", branchConfig.Name)
	}

}
