package types

import (
	mint "github.com/vinyl-linux/mint"
	"io"
)

type Predicate byte

const (
	PredicateUnknown Predicate = iota
	PredicateHasChild
	PredicateExtends
	PredicateSupercedes
	PredicateSupplements
)

func (sf Predicate) Marshall(w io.Writer) (err error) {
	return mint.NewByteScalar(byte(sf)).Marshall(w)
}
func (sf *Predicate) Unmarshall(r io.Reader) (err error) {
	f := mint.NewByteScalar(byte(int32(0)))
	err = f.Unmarshall(r)
	if err != nil {
		return
	}
	*sf = Predicate(f.Value().(byte))
	return
}
func (sf Predicate) Value() any {
	return sf
}
