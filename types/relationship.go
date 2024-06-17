package types

import (
	mint "github.com/vinyl-linux/mint"
	"io"
)

type Relationship struct {
	Subject   PageRef
	Predicate Predicate
	Object    PageRef
}

func (sf Relationship) Validate() error {
	errors := make([]error, 0)
	for _, err := range []error{} {
		if err != nil {
			errors = append(errors, err)
		}
	}
	return mint.ValidationErrors("Relationship", errors)
}
func (sf *Relationship) Transform() (err error) {
	return
}
func (sf Relationship) Value() any {
	return sf
}
func (sf *Relationship) unmarshallSubject(r io.Reader) (err error) {
	f := new(PageRef)
	err = f.Unmarshall(r)
	if err != nil {
		return
	}
	sf.Subject = f.Value().(PageRef)
	return
}
func (sf *Relationship) unmarshallPredicate(r io.Reader) (err error) {
	f := new(Predicate)
	err = f.Unmarshall(r)
	if err != nil {
		return
	}
	sf.Predicate = f.Value().(Predicate)
	return
}
func (sf *Relationship) unmarshallObject(r io.Reader) (err error) {
	f := new(PageRef)
	err = f.Unmarshall(r)
	if err != nil {
		return
	}
	sf.Object = f.Value().(PageRef)
	return
}
func (sf *Relationship) Unmarshall(r io.Reader) (err error) {
	if err = sf.unmarshallSubject(r); err != nil {
		return
	}
	if err = sf.unmarshallPredicate(r); err != nil {
		return
	}
	if err = sf.unmarshallObject(r); err != nil {
		return
	}
	if err = sf.Transform(); err != nil {
		return
	}
	if err = sf.Validate(); err != nil {
		return
	}
	return
}
func (sf Relationship) Marshall(w io.Writer) (err error) {
	if err = sf.Transform(); err != nil {
		return
	}
	if err = sf.Validate(); err != nil {
		return
	}
	if err = sf.Subject.Marshall(w); err != nil {
		return
	}
	if err = sf.Predicate.Marshall(w); err != nil {
		return
	}
	if err = sf.Object.Marshall(w); err != nil {
		return
	}
	return
}
