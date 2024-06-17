package types

import (
	v5 "github.com/gofrs/uuid/v5"
	mint "github.com/vinyl-linux/mint"
	"io"
)

type Request struct {
	// Verb is analogous to an HTTP request type and defines the expectations the client has for the type of response from the server
	Verb Verb
	// ID is the ID of the Page being operated on
	ID v5.UUID
	// Args are arbitrary arguments that a server may or may not respond to.  There are some arguments that will always be expected for verbs such as `Section` for an Update, or `Body` for both Create and Update
	Args map[string]string
}

func (sf Request) Validate() error {
	errors := make([]error, 0)
	for _, err := range []error{} {
		if err != nil {
			errors = append(errors, err)
		}
	}
	return mint.ValidationErrors("Request", errors)
}
func (sf *Request) Transform() (err error) {
	return
}
func (sf Request) Value() any {
	return sf
}
func (sf *Request) unmarshallVerb(r io.Reader) (err error) {
	f := new(Verb)
	err = f.Unmarshall(r)
	if err != nil {
		return
	}
	sf.Verb = f.Value().(Verb)
	return
}
func (sf *Request) unmarshallID(r io.Reader) (err error) {
	f := mint.NewUuidScalar(v5.UUID{})
	err = f.Unmarshall(r)
	if err != nil {
		return
	}
	sf.ID = f.Value().(v5.UUID)
	return
}
func (sf *Request) unmarshallArgs(r io.Reader) (err error) {
	f := mint.NewMapCollection(map[mint.MarshallerUnmarshallerValuer]mint.MarshallerUnmarshallerValuer{})
	err = f.ReadSize(r)
	if err != nil {
		return
	}
	if f.Len() == 0 {
		sf.Args = nil
		return
	}
	for i := 0; i < f.Len(); i++ {
		f.V[mint.NewStringScalar("")] = mint.NewStringScalar("")
	}
	err = f.Unmarshall(r)
	if err != nil {
		return
	}
	sf.Args = make(map[string]string)
	for k, v := range f.Value().(map[mint.MarshallerUnmarshallerValuer]mint.MarshallerUnmarshallerValuer) {
		sf.Args[k.Value().(string)] = v.Value().(string)
	}
	return
}
func (sf *Request) Unmarshall(r io.Reader) (err error) {
	if err = sf.unmarshallVerb(r); err != nil {
		return
	}
	if err = sf.unmarshallID(r); err != nil {
		return
	}
	if err = sf.unmarshallArgs(r); err != nil {
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
func (sf Request) marshallArgs(w io.Writer) (err error) {
	f := make(map[mint.MarshallerUnmarshallerValuer]mint.MarshallerUnmarshallerValuer)
	for k, v := range sf.Args {
		f[mint.NewStringScalar(k)] = mint.NewStringScalar(v)
	}
	return mint.NewMapCollection(f).Marshall(w)
}
func (sf Request) Marshall(w io.Writer) (err error) {
	if err = sf.Transform(); err != nil {
		return
	}
	if err = sf.Validate(); err != nil {
		return
	}
	if err = sf.Verb.Marshall(w); err != nil {
		return
	}
	if err = mint.NewUuidScalar(sf.ID).Marshall(w); err != nil {
		return
	}
	if err = sf.marshallArgs(w); err != nil {
		return
	}
	return
}
