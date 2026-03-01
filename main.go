package main

import (
	"bytes"
	"fmt"
)

type Commit struct {
	Fmt     []byte
	Header  map[string][]string
	Message []byte
}

func kvlmParse(raw []byte, start int, kvlm map[string][]string) (map[string][]string, []byte, error) {
	spc := bytes.Index(raw[start:], []byte(" "))
	nl := bytes.Index(raw[start:], []byte("\n"))
	if spc == -1 || nl == -1 || spc > nl {
		return kvlm, raw[start+1:], nil
	}
	spc += start
	nl += start

	key := raw[start:spc]
	end := nl

	for end+1 < len(raw) && raw[end+1] == ' ' {
		nextNl := bytes.Index(raw[end+1:], []byte("\n"))
		if nextNl == -1 {
			end = len(raw)
			break
		}
		end += nextNl + 1
	}
	value := bytes.ReplaceAll(raw[spc+1:end], []byte("\n "), []byte("\n"))
	v, ok := kvlm[string(key)]
	if ok {
		kvlm[string(key)] = append(v, string(value))
	} else {
		kvlm[string(key)] = []string{string(value)}
	}

	return kvlmParse(raw, end+1, kvlm)
}

func main() {

	raw := []byte(`tree abc123
parent def456
parent fedcba
author Blue Onion
 <blue@onion.com>
committer Blue Onion <blue@onion.com>

This is the commit message
With multiple lines
And even more lines
`)

	kvlm := make(map[string][]string)

	headers, message, err := kvlmParse(raw, 0, kvlm)
	if err != nil {
		panic(err)
	}

	commit := Commit{
		Fmt:     []byte("commit"),
		Header:  headers,
		Message: message,
	}

	fmt.Println("========== COMMIT OBJECT ==========")
	fmt.Println("Format:", string(commit.Fmt))

	fmt.Println("\n---- HEADERS ----")
	for k, v := range commit.Header {
		fmt.Printf("%s:\n", k)
		for _, val := range v {
			fmt.Printf("  %s\n", val)
		}
	}

	fmt.Println("\n---- MESSAGE ----")
	fmt.Println(string(commit.Message))
}
