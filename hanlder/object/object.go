package object

import (
	"bytes"
	"compress/zlib"
	"io"
	"os"

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
func ObjectRead(repo *repo.Gitrepo, name string, typ string) (GitObject, error) {
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
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	i := bytes.IndexByte(data, 0)
	data = data[i+1:]

	obj := &Blob{Data: data, Fmt: []byte("blob")}
	return obj, nil
}
func ObjectWrite(repo *repo.Gitrepo) {

}
