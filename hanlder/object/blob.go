package object
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