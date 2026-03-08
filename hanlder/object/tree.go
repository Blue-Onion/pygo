package object

import (
	"bytes"
	"fmt"
)
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
	fmt.Println(string(raw))
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