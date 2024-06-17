package types

import (
	mint "github.com/vinyl-linux/mint"
	"io"
)

type Verb byte

const (
	VerbUnknown Verb = iota
	VerbCreate
	VerbRead
	VerbUpdate
	VerbDelete
)

func (sf Verb) Marshall(w io.Writer) (err error) {
	return mint.NewByteScalar(byte(sf)).Marshall(w)
}
func (sf *Verb) Unmarshall(r io.Reader) (err error) {
	f := mint.NewByteScalar(byte(int32(0)))
	err = f.Unmarshall(r)
	if err != nil {
		return
	}
	*sf = Verb(f.Value().(byte))
	return
}
func (sf Verb) Value() any {
	return sf
}
