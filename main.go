package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
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

	if _, err := w.Commit(*message, opts_commit); err != nil {
		log.Fatal("Error unable to commit")
	}

	
}
