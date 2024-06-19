package types

import (
	"errors"
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
	if sf < 1 || sf > 4 {
		return errors.New("invalid value for type Predicate")
	}
	return mint.NewByteScalar(byte(sf)).Marshall(w)
}
func (sf *Predicate) Unmarshall(r io.Reader) (err error) {
	f := mint.NewByteScalar(byte(int32(0)))
	err = f.Unmarshall(r)
	if err != nil {
		return
	}
	*sf = Predicate(f.Value().(byte))
	if *sf < 1 || *sf > 4 {
		return errors.New("invalid value for type Predicate")
	}
	return
}
func (sf Predicate) Value() any {
	return sf
}
