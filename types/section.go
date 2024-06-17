package types

import (
	mint "github.com/vinyl-linux/mint"
	"io"
)

type Section struct {
	Title string
	Body  string
}

func (sf Section) Validate() error {
	errors := make([]error, 0)
	for _, err := range []error{} {
		if err != nil {
			errors = append(errors, err)
		}
	}
	return mint.ValidationErrors("Section", errors)
}
func (sf *Section) Transform() (err error) {
	return
}
func (sf Section) Value() any {
	return sf
}
func (sf *Section) unmarshallTitle(r io.Reader) (err error) {
	f := mint.NewStringScalar("")
	err = f.Unmarshall(r)
	if err != nil {
		return
	}
	sf.Title = f.Value().(string)
	return
}
func (sf *Section) unmarshallBody(r io.Reader) (err error) {
	f := mint.NewStringScalar("")
	err = f.Unmarshall(r)
	if err != nil {
		return
	}
	sf.Body = f.Value().(string)
	return
}
func (sf *Section) Unmarshall(r io.Reader) (err error) {
	if err = sf.unmarshallTitle(r); err != nil {
		return
	}
	if err = sf.unmarshallBody(r); err != nil {
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
func (sf Section) Marshall(w io.Writer) (err error) {
	if err = sf.Transform(); err != nil {
		return
	}
	if err = sf.Validate(); err != nil {
		return
	}
	if err = mint.NewStringScalar(sf.Title).Marshall(w); err != nil {
		return
	}
	if err = mint.NewStringScalar(sf.Body).Marshall(w); err != nil {
		return
	}
	return
}
