package types

import (
	v5 "github.com/gofrs/uuid/v5"
	mint "github.com/vinyl-linux/mint"
	"io"
	"time"
)

type Metadata struct {
	// ID of this resource
	ID v5.UUID
	// Author of the resource
	Author string
	// Published date of the resource
	Published time.Time
}

func (sf Metadata) Validate() error {
	errors := make([]error, 0)
	for _, err := range []error{} {
		if err != nil {
			errors = append(errors, err)
		}
	}
	return mint.ValidationErrors("Metadata", errors)
}
func (sf *Metadata) Transform() (err error) {
	sf.Published, err = mint.DateInUtc(sf.Published)
	if err != nil {
		return
	}
	return
}
func (sf Metadata) Value() any {
	return sf
}
func (sf *Metadata) unmarshallID(r io.Reader) (err error) {
	f := mint.NewUuidScalar(v5.UUID{})
	err = f.Unmarshall(r)
	if err != nil {
		return
	}
	sf.ID = f.Value().(v5.UUID)
	return
}
func (sf *Metadata) unmarshallAuthor(r io.Reader) (err error) {
	f := mint.NewStringScalar("")
	err = f.Unmarshall(r)
	if err != nil {
		return
	}
	sf.Author = f.Value().(string)
	return
}
func (sf *Metadata) unmarshallPublished(r io.Reader) (err error) {
	f := mint.NewDatetimeScalar(time.Time{})
	err = f.Unmarshall(r)
	if err != nil {
		return
	}
	sf.Published = f.Value().(time.Time)
	return
}
func (sf *Metadata) Unmarshall(r io.Reader) (err error) {
	if err = sf.unmarshallID(r); err != nil {
		return
	}
	if err = sf.unmarshallAuthor(r); err != nil {
		return
	}
	if err = sf.unmarshallPublished(r); err != nil {
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
func (sf Metadata) Marshall(w io.Writer) (err error) {
	if err = sf.Transform(); err != nil {
		return
	}
	if err = sf.Validate(); err != nil {
		return
	}
	if err = mint.NewUuidScalar(sf.ID).Marshall(w); err != nil {
		return
	}
	if err = mint.NewStringScalar(sf.Author).Marshall(w); err != nil {
		return
	}
	if err = mint.NewDatetimeScalar(sf.Published).Marshall(w); err != nil {
		return
	}
	return
}
