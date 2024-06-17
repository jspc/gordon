package types

import (
	mint "github.com/vinyl-linux/mint"
	"io"
)

type Page struct {
	// Metadata contains metadata about this page, including useful stuff like author details, and revisions, and dates and all that stuff
	Meta Metadata
	// History is a slice of Metadatas that point to the previous versions of this page, such as they may be. There are no limits to this field, and indeed no requirement for it to even contain anything
	History []Metadata
	// Title represents the title of a page, funnily enough, and must be between 1 and 512 characters
	Title string
	// Preamble is used on Page Indexes and other list operations
	Preamble string
	// Sections contains a list of sections in a document. A section is how we divvy up data
	Sections []Section
	// Tags are an arbitrary list of strings to attach to a page for searching and organising
	Tags []string
	// Labels are Much like tags; handy for searching and organising data
	Labels map[string]string
	// Links are references to other pages; within the body of a section, they are referenced by their index
	Links []PageRef
	// Relationships are used to link pages
	Relationships []Relationship
	// Status reflects whether this page is to be treated as an error page or not
	Status Status
}

func (sf Page) Validate() error {
	errors := make([]error, 0)
	for _, err := range []error{mint.StringNotEmpty("Title", sf.Title), sf.TitleNotTooLong("Title", sf.Title)} {
		if err != nil {
			errors = append(errors, err)
		}
	}
	return mint.ValidationErrors("Page", errors)
}
func (sf *Page) Transform() (err error) {
	return
}
func (sf Page) Value() any {
	return sf
}
func (sf *Page) unmarshallMeta(r io.Reader) (err error) {
	f := new(Metadata)
	err = f.Unmarshall(r)
	if err != nil {
		return
	}
	sf.Meta = f.Value().(Metadata)
	return
}
func (sf *Page) unmarshallHistory(r io.Reader) (err error) {
	f := mint.NewSliceCollection(nil, false)
	err = f.ReadSize(r)
	if err != nil {
		return
	}
	if f.Len() == 0 {
		sf.History = nil
		return
	}
	f.V = make([]mint.MarshallerUnmarshallerValuer, f.Len())
	for i := range f.V {
		f.V[i] = new(Metadata)
	}
	err = f.Unmarshall(r)
	if err != nil {
		return
	}
	sf.History = make([]Metadata, f.Len())
	for i, v := range f.Value().([]mint.MarshallerUnmarshallerValuer) {
		sf.History[i] = v.Value().(Metadata)
	}
	return
}
func (sf *Page) unmarshallTitle(r io.Reader) (err error) {
	f := mint.NewStringScalar("")
	err = f.Unmarshall(r)
	if err != nil {
		return
	}
	sf.Title = f.Value().(string)
	return
}
func (sf *Page) unmarshallPreamble(r io.Reader) (err error) {
	f := mint.NewStringScalar("")
	err = f.Unmarshall(r)
	if err != nil {
		return
	}
	sf.Preamble = f.Value().(string)
	return
}
func (sf *Page) unmarshallSections(r io.Reader) (err error) {
	f := mint.NewSliceCollection(nil, false)
	err = f.ReadSize(r)
	if err != nil {
		return
	}
	if f.Len() == 0 {
		sf.Sections = nil
		return
	}
	f.V = make([]mint.MarshallerUnmarshallerValuer, f.Len())
	for i := range f.V {
		f.V[i] = new(Section)
	}
	err = f.Unmarshall(r)
	if err != nil {
		return
	}
	sf.Sections = make([]Section, f.Len())
	for i, v := range f.Value().([]mint.MarshallerUnmarshallerValuer) {
		sf.Sections[i] = v.Value().(Section)
	}
	return
}
func (sf *Page) unmarshallTags(r io.Reader) (err error) {
	f := mint.NewSliceCollection(nil, false)
	err = f.ReadSize(r)
	if err != nil {
		return
	}
	if f.Len() == 0 {
		sf.Tags = nil
		return
	}
	f.V = make([]mint.MarshallerUnmarshallerValuer, f.Len())
	for i := range f.V {
		f.V[i] = mint.NewStringScalar("")
	}
	err = f.Unmarshall(r)
	if err != nil {
		return
	}
	sf.Tags = make([]string, f.Len())
	for i, v := range f.Value().([]mint.MarshallerUnmarshallerValuer) {
		sf.Tags[i] = v.Value().(string)
	}
	return
}
func (sf *Page) unmarshallLabels(r io.Reader) (err error) {
	f := mint.NewMapCollection(map[mint.MarshallerUnmarshallerValuer]mint.MarshallerUnmarshallerValuer{})
	err = f.ReadSize(r)
	if err != nil {
		return
	}
	if f.Len() == 0 {
		sf.Labels = nil
		return
	}
	for i := 0; i < f.Len(); i++ {
		f.V[mint.NewStringScalar("")] = mint.NewStringScalar("")
	}
	err = f.Unmarshall(r)
	if err != nil {
		return
	}
	sf.Labels = make(map[string]string)
	for k, v := range f.Value().(map[mint.MarshallerUnmarshallerValuer]mint.MarshallerUnmarshallerValuer) {
		sf.Labels[k.Value().(string)] = v.Value().(string)
	}
	return
}
func (sf *Page) unmarshallLinks(r io.Reader) (err error) {
	f := mint.NewSliceCollection(nil, false)
	err = f.ReadSize(r)
	if err != nil {
		return
	}
	if f.Len() == 0 {
		sf.Links = nil
		return
	}
	f.V = make([]mint.MarshallerUnmarshallerValuer, f.Len())
	for i := range f.V {
		f.V[i] = new(PageRef)
	}
	err = f.Unmarshall(r)
	if err != nil {
		return
	}
	sf.Links = make([]PageRef, f.Len())
	for i, v := range f.Value().([]mint.MarshallerUnmarshallerValuer) {
		sf.Links[i] = v.Value().(PageRef)
	}
	return
}
func (sf *Page) unmarshallRelationships(r io.Reader) (err error) {
	f := mint.NewSliceCollection(nil, false)
	err = f.ReadSize(r)
	if err != nil {
		return
	}
	if f.Len() == 0 {
		sf.Relationships = nil
		return
	}
	f.V = make([]mint.MarshallerUnmarshallerValuer, f.Len())
	for i := range f.V {
		f.V[i] = new(Relationship)
	}
	err = f.Unmarshall(r)
	if err != nil {
		return
	}
	sf.Relationships = make([]Relationship, f.Len())
	for i, v := range f.Value().([]mint.MarshallerUnmarshallerValuer) {
		sf.Relationships[i] = v.Value().(Relationship)
	}
	return
}
func (sf *Page) unmarshallStatus(r io.Reader) (err error) {
	f := new(Status)
	err = f.Unmarshall(r)
	if err != nil {
		return
	}
	sf.Status = f.Value().(Status)
	return
}
func (sf *Page) Unmarshall(r io.Reader) (err error) {
	if err = sf.unmarshallMeta(r); err != nil {
		return
	}
	if err = sf.unmarshallHistory(r); err != nil {
		return
	}
	if err = sf.unmarshallTitle(r); err != nil {
		return
	}
	if err = sf.unmarshallPreamble(r); err != nil {
		return
	}
	if err = sf.unmarshallSections(r); err != nil {
		return
	}
	if err = sf.unmarshallTags(r); err != nil {
		return
	}
	if err = sf.unmarshallLabels(r); err != nil {
		return
	}
	if err = sf.unmarshallLinks(r); err != nil {
		return
	}
	if err = sf.unmarshallRelationships(r); err != nil {
		return
	}
	if err = sf.unmarshallStatus(r); err != nil {
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
func (sf Page) marshallHistory(w io.Writer) (err error) {
	f := make([]mint.MarshallerUnmarshallerValuer, len(sf.History))
	for i := range f {
		f[i] = &(sf.History[i])
	}
	return mint.NewSliceCollection(f, false).Marshall(w)
}
func (sf Page) marshallSections(w io.Writer) (err error) {
	f := make([]mint.MarshallerUnmarshallerValuer, len(sf.Sections))
	for i := range f {
		f[i] = &(sf.Sections[i])
	}
	return mint.NewSliceCollection(f, false).Marshall(w)
}
func (sf Page) marshallTags(w io.Writer) (err error) {
	f := make([]mint.MarshallerUnmarshallerValuer, len(sf.Tags))
	for i := range f {
		f[i] = mint.NewStringScalar(sf.Tags[i])
	}
	return mint.NewSliceCollection(f, false).Marshall(w)
}
func (sf Page) marshallLabels(w io.Writer) (err error) {
	f := make(map[mint.MarshallerUnmarshallerValuer]mint.MarshallerUnmarshallerValuer)
	for k, v := range sf.Labels {
		f[mint.NewStringScalar(k)] = mint.NewStringScalar(v)
	}
	return mint.NewMapCollection(f).Marshall(w)
}
func (sf Page) marshallLinks(w io.Writer) (err error) {
	f := make([]mint.MarshallerUnmarshallerValuer, len(sf.Links))
	for i := range f {
		f[i] = &(sf.Links[i])
	}
	return mint.NewSliceCollection(f, false).Marshall(w)
}
func (sf Page) marshallRelationships(w io.Writer) (err error) {
	f := make([]mint.MarshallerUnmarshallerValuer, len(sf.Relationships))
	for i := range f {
		f[i] = &(sf.Relationships[i])
	}
	return mint.NewSliceCollection(f, false).Marshall(w)
}
func (sf Page) Marshall(w io.Writer) (err error) {
	if err = sf.Transform(); err != nil {
		return
	}
	if err = sf.Validate(); err != nil {
		return
	}
	if err = sf.Meta.Marshall(w); err != nil {
		return
	}
	if err = sf.marshallHistory(w); err != nil {
		return
	}
	if err = mint.NewStringScalar(sf.Title).Marshall(w); err != nil {
		return
	}
	if err = mint.NewStringScalar(sf.Preamble).Marshall(w); err != nil {
		return
	}
	if err = sf.marshallSections(w); err != nil {
		return
	}
	if err = sf.marshallTags(w); err != nil {
		return
	}
	if err = sf.marshallLabels(w); err != nil {
		return
	}
	if err = sf.marshallLinks(w); err != nil {
		return
	}
	if err = sf.marshallRelationships(w); err != nil {
		return
	}
	if err = sf.Status.Marshall(w); err != nil {
		return
	}
	return
}
