package object

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
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
func (b *Blob) Type() string {
	return "blob"
}

type Tree struct {
	Data []TreeData
	Fmt  []byte
}
type TreeData struct {
	Mode []byte
	Name []byte
	Sha  []byte
}

func (t *Tree) Serialize() ([]byte, error) {
	var out bytes.Buffer

	for _, entry := range t.Data {

		if len(entry.Sha) != 20 {
			return nil, fmt.Errorf("invalid sha length: expected 20 bytes")
		}
		out.Write(entry.Mode)
		out.WriteByte(' ')
		out.Write(entry.Name)
		out.WriteByte(0)
		out.Write(entry.Sha)
	}

	return out.Bytes(), nil
}

func (t *Tree) Deserialize(raw []byte) error {
	t.Data = nil
	n := 0

	for n < len(raw) {
		spaceI := bytes.IndexByte(raw[n:], ' ')
		if spaceI == -1 {
			return fmt.Errorf("invalid tree: no space found")
		}
		spaceI += n
		mode := raw[n:spaceI]
		nullI := bytes.IndexByte(raw[spaceI+1:], 0)
		if nullI == -1 {
			return fmt.Errorf("invalid tree: no null found")
		}
		nullI += spaceI + 1
		name := raw[spaceI+1 : nullI]
		shaStart := nullI + 1
		shaEnd := shaStart + 20
		if shaEnd > len(raw) {
			return fmt.Errorf("invalid tree: sha overflow")
		}
		sha := raw[shaStart:shaEnd]
		entry := TreeData{
			Mode: mode,
			Name: name,
			Sha:  sha,
		}
		t.Data = append(t.Data, entry)
		n = shaEnd
	}

	return nil
}

func (t *Tree) Type() string {
	return "tree"
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
		obj = &Blob{}
		obj.Deserialize(content)
	case "tree":
		obj = &Tree{}
		obj.Deserialize(content)
	default:
		return nil, fmt.Errorf("Type not found")
	}
	return obj, nil
}

func HashString(objType string, data []byte) string {

	header := objType + " " + strconv.Itoa(len(data)) + "\x00"

	store := append([]byte(header), data...)

	hash := sha1.Sum(store)

	return fmt.Sprintf("%x", hash)
}
func ObjectHash(path string, typ string, repo *repo.Gitrepo) (string, error) {
	var obj GitObject
	switch typ {
	case "blob":
		obj = &Blob{}
	case "tree":
		obj = &Tree{}
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	err = obj.Deserialize(data)
	if err != nil {
		return "", err
	}

	return ObjectWrite(repo, obj)

}
func ObjectWrite(Gitrepo *repo.Gitrepo, obj GitObject) (string, error) {

	data, err := obj.Serialize()
	if err != nil {
		return "", err
	}

	header := obj.Type() + " " + strconv.Itoa(len(data)) + "\x00"
	store := append([]byte(header), data...)

	hash := sha1.Sum(store)
	sha := fmt.Sprintf("%x", hash)

	path, err := repo.RepoFile(Gitrepo, true, "objects", sha[:2], sha[2:])
	if err != nil {
		return "", err
	}

	exist, _ := repo.PathExist(path)
	if exist {
		return sha, nil
	}

	var buf bytes.Buffer
	w := zlib.NewWriter(&buf)
	if _, err := w.Write(store); err != nil {
		return "", err
	}
	if err := w.Close(); err != nil {
		return "", err
	}

	err = os.WriteFile(path, buf.Bytes(), 0644)
	if err != nil {
		return "", err
	}

	return sha, nil
}
func CatFile(repo *repo.Gitrepo, name string, tag string) {
	obj, err := ObjectRead(repo, name)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	switch tag {
	case "-p":
		switch o := obj.(type) {
		case *Tree:
			for _, v := range o.Data {
				fmt.Printf("%s %x\t%s\n",
					v.Mode,
					v.Sha,
					v.Name,
				)
			}
		case *Blob:
			fmt.Print(string(o.Data))
		}
	case "-t":
		fmt.Println(obj.Type())
	}
}
