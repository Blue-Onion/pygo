package main

import (
	"fmt"
	"os"

	"github.com/Blue-Onion/pygo/hanlder/repo"
)

func cmdInit(path string) (*repo.Gitrepo, error) {
	
	r, err := repo.RepoCreate(path)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func main() {
	var path string
	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Println("Invalid command. Usage: init [path]")
		return
	}

	cmd := args[0]

	// Default path = current directory if not given
	if len(args) >= 2 {
		path = args[1]
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Error getting current directory:", err)
			return
		}
		path = cwd
	}

	switch cmd {
	case "init":
		_, err := cmdInit(path)
		if err != nil {
			fmt.Println("Error initializing repo:", err)
			return
		}
		fmt.Println("Repository initialized successfully!")
	default:
		fmt.Println("Invalid command. Available commands: init")
	}
}

