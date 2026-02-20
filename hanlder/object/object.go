package object

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/Blue-Onion/pygo/hanlder/repo"
)

type GitObject interface {
	Serialize() ([]byte, error)
	Deserialize([]byte) error
	Type() string
}
type Blob struct {
	Data []byte
	Fmt  []byte
}

func (b *Blob) Serialize() ([]byte, error) {
	return b.Data, nil
}

func (b *Blob) Deserialize(raw []byte) error {
	b.Data = raw
	return nil
}

type Tree struct {
	Data []byte
	Fmt  []byte
}

func (b *Tree) Serialize() ([]byte, error) {
	return b.Data, nil
}

func (b *Tree) Deserialize(raw []byte) error {
	b.Data = raw
	return nil
}

func (b *Tree) Type() string {
	return "tree"
}
func (b *Blob) Type() string {
	return "blob"
}
func lengthAndContent(raw []byte) (int, []byte, error) {
	parts := bytes.Split(raw, []byte(" "))
	if len(parts) != 2 {
		return 0, []byte(""), fmt.Errorf("Malformed Content")
	}
	length, err := strconv.Atoi(string(parts[1]))
	if err != nil {
		return 0, []byte(""), fmt.Errorf("Malformed Content-length")
	}
	return length, parts[0], nil

}
func ObjectRead(repo *repo.Gitrepo, name string) (GitObject, error) {
	file := name[2:]
	dir := name[:2]
	path := repo.Gitdir + "/objects/" + dir + "/" + file
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	r, err := zlib.NewReader(bytes.NewReader(raw))
	if err != nil {
		return nil, err
	}
	defer r.Close()
	rawdata, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	i := bytes.IndexByte(rawdata, 0)
	headers := rawdata[:i]
	content := rawdata[i+1:]
	_, typ, err := lengthAndContent(headers)
	if err != nil {
		return nil, err
	}
	var obj GitObject
	switch string(typ) {

	case "blob":
		obj = &Blob{
			Data: content,
			Fmt:  []byte("blob"),
		}
	case "tree":
		obj = &Tree{
			Data: content,
			Fmt:  []byte("tree"),
		}
	default:
		return nil, fmt.Errorf("Type not found")
	}
	return obj, nil
}

func CatFile(repo *repo.Gitrepo, name string, tag string) ([]byte, error) {
	obj, err := ObjectRead(repo, name)
	if err != nil {
		return []byte(""), err
	}
	var data []byte
	switch tag{
	case "-p":
		data, err = obj.Serialize()
	case "-t":
		data = []byte(obj.Type())

	}
	if err != nil {
		return []byte(""), err
	}
	return data, nil
}
