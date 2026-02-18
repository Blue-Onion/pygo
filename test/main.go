package main

import (
	"fmt"
	"log"

	"github.com/Blue-Onion/pygo/hanlder/object"
	"github.com/Blue-Onion/pygo/hanlder/repo"
)

func main() {
	r, err := repo.RepoFind(".", true)
	if err != nil {
		log.Fatal(err)
	}

	// get a real blob SHA from your repo:
	// git hash-object somefile
	sha := "0becd8c218e287eed9d25e98f4b302de8113c06d"

	obj, err := object.ObjectRead(r, sha, "blob")
	if err != nil {
		log.Fatal(err)
	}

	blob := obj.(*object.Blob)

	fmt.Println("Blob content:")
	fmt.Println(string(blob.Data))
}