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
func cmdCatFile(path string, name string, typ string) {
	repo, err := repo.RepoFind(path, true)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	object.CatFile(repo, name, typ)
	


}
func cmdHashObject(path string, args []string) {
	repo, err := repo.RepoFind(path, true)
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(args) < 1 {
		fmt.Println("Usage: hash-object [-t type] <file>")
		return
	}

	typ := "blob" // default like real git
	file := ""

	for i := 0; i < len(args); i++ {

		switch args[i] {
		case "-t":
			if i+1 >= len(args) {
				fmt.Println("Missing type after -t")
				return
			}
			typ = args[i+1]
			i++
		default:
			file = args[i]
		}
	}

	if file == "" {
		fmt.Println("No file specified")
		return
	}

	sha, err := object.ObjectHash(file, typ, repo)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(sha)
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

		if len(args) < 3 {
			fmt.Println("Usage: cat-file <type> <object> [path]")
			return
		}

		tag := args[1]
		name := args[2]

		if len(args) >= 4 {
			path = args[3]
		}
		cmdCatFile(path, name, tag)
	case "hash-object":
		// Usage: hash-object [-t type] [-w] <file> [path]
		if len(args) < 2 {
			fmt.Println("Usage: hash-object [-t type] [-w] <file> [path]")
			return
		}
	
		// Optional repo path as last argument
		if len(args) >= 3 {
			last := args[len(args)-1]
			if info, err := os.Stat(last); err == nil && info.IsDir() {
				path = last
				args = args[1 : len(args)-1]
			} else {
				args = args[1:]
			}
		} else {
			args = args[1:]
		}
	
		cmdHashObject(path, args)
	default:
		fmt.Println("Invalid command. Available commands: init, cat-file")
	}
}

func getConcatenation(nums []int) []int {
    res:=make([]int,len(nums)*2)
	for i,v:=range nums{
		res[i]=v
		res[i+len(nums)]=v

	}
	return res
}
