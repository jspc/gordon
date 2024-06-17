package types

import (
	v5 "github.com/gofrs/uuid/v5"
	mint "github.com/vinyl-linux/mint"
	"io"
)

type PageRef struct {
	Page    v5.UUID
	Section string
	Server  string
}

func (sf PageRef) Validate() error {
	errors := make([]error, 0)
	for _, err := range []error{} {
		if err != nil {
			errors = append(errors, err)
		}
	}
	return mint.ValidationErrors("PageRef", errors)
}
func (sf *PageRef) Transform() (err error) {
	return
}
func (sf PageRef) Value() any {
	return sf
}
func (sf *PageRef) unmarshallPage(r io.Reader) (err error) {
	f := mint.NewUuidScalar(v5.UUID{})
	err = f.Unmarshall(r)
	if err != nil {
		return
	}
	sf.Page = f.Value().(v5.UUID)
	return
}
func (sf *PageRef) unmarshallSection(r io.Reader) (err error) {
	f := mint.NewStringScalar("")
	err = f.Unmarshall(r)
	if err != nil {
		return
	}
	sf.Section = f.Value().(string)
	return
}
func (sf *PageRef) unmarshallServer(r io.Reader) (err error) {
	f := mint.NewStringScalar("")
	err = f.Unmarshall(r)
	if err != nil {
		return
	}
	sf.Server = f.Value().(string)
	return
}
func (sf *PageRef) Unmarshall(r io.Reader) (err error) {
	if err = sf.unmarshallPage(r); err != nil {
		return
	}
	if err = sf.unmarshallSection(r); err != nil {
		return
	}
	if err = sf.unmarshallServer(r); err != nil {
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
func (sf PageRef) Marshall(w io.Writer) (err error) {
	if err = sf.Transform(); err != nil {
		return
	}
	if err = sf.Validate(); err != nil {
		return
	}
	if err = mint.NewUuidScalar(sf.Page).Marshall(w); err != nil {
		return
	}
	if err = mint.NewStringScalar(sf.Section).Marshall(w); err != nil {
		return
	}
	if err = mint.NewStringScalar(sf.Server).Marshall(w); err != nil {
		return
	}
	return
}
