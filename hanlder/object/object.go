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
	case "commit":
		obj = &Commit{}
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
				fmt.Printf("%s %s\t%s\n",
					v.Mode,
					fmt.Sprintf("%x", v.Sha), // <-- hex string
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
