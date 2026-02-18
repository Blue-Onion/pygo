package object

import (
	"github.com/Blue-Onion/pygo/hanlder/repo"
)

type GitObject interface {
	Serialize() ([]byte, error)
	Deserialize([]byte) error
	Type() string
}
type Blob struct {
	Data []byte
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
func ObjectRead(repo *repo.Gitrepo)  {
	
}
func ObjectWrite(repo *repo.Gitrepo){
	
}