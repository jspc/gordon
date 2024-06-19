package types

import (
	"errors"
	mint "github.com/vinyl-linux/mint"
	"io"
)

type Status byte

const (
	StatusUnknown Status = iota
	StatusOK
	StatusError
)

func (sf Status) Marshall(w io.Writer) (err error) {
	if sf < 1 || sf > 2 {
		return errors.New("invalid value for type Status")
	}
	return mint.NewByteScalar(byte(sf)).Marshall(w)
}
func (sf *Status) Unmarshall(r io.Reader) (err error) {
	f := mint.NewByteScalar(byte(int32(0)))
	err = f.Unmarshall(r)
	if err != nil {
		return
	}
	*sf = Status(f.Value().(byte))
	if *sf < 1 || *sf > 2 {
		return errors.New("invalid value for type Status")
	}
	return
}
func (sf Status) Value() any {
	return sf
}
