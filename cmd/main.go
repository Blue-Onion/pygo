package main

import (
	"fmt"
	"os"

	"github.com/Blue-Onion/pygo/hanlder/object"
	"github.com/Blue-Onion/pygo/hanlder/repo"
)

func cmdInit(path string) (*repo.Gitrepo, error) {

	r, err := repo.RepoCreate(path)
	if err != nil {
		return nil, err
	}
	return r, nil
}
func CmdCatFile(path string, name string, typ string) {
	repo, err := repo.RepoFind(path, true)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	data, err := object.CatFile(repo, name, typ)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(string(data))
}
func main() {
	var path string
	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Println("Invalid command.")
		return
	}

	cmd := args[0]

	// Default path = current directory
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}
	path = cwd

	switch cmd {

	case "init":
		if len(args) >= 2 {
			path = args[1]
		}

		_, err := cmdInit(path)
		if err != nil {
			fmt.Println("Error initializing repo:", err)
			return
		}
		fmt.Println("Repository initialized successfully!")

	case "cat-file":
		// Usage: cat-file <type> <object> [path]
		if len(args) < 3 {
			fmt.Println("Usage: cat-file <type> <object> [path]")
			return
		}

		tag := args[1]
		name := args[2]

		if len(args) >= 4 {
			path = args[3]
		}
	
		CmdCatFile(path, name, tag)

	default:
		fmt.Println("Invalid command. Available commands: init, cat-file")
	}
}

