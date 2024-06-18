package main

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/jspc/gordon/types"
)

type Server struct {
	pages map[uuid.UUID]types.Page
}

func (s Server) Serve(req *types.Request) (resp types.Page, err error) {
	switch req.Verb {
	case types.VerbRead:
		return s.serveRead(req)

	default:
		return s.verbNotSupported(req.ID), nil
	}
}

func (s Server) serveRead(req *types.Request) (resp types.Page, err error) {
	if req.ID.IsNil() {
		return s.indexPage(), nil
	}

	resp, ok := s.pages[req.ID]
	if ok {
		return
	}

	return s.pageNotFound(req.ID), nil
}

func (s Server) indexPage() (p types.Page) {
	p = types.Page{
		Title:  "Page Index",
		Status: types.StatusOK,
		Meta: types.Metadata{
			Author:    "Gordon",
			Published: time.Now(),
		},
		Sections: make([]types.Section, 0),
		Links:    make([]types.PageRef, 0),
		Tags:     []string{"index"},
	}

	for id, page := range s.pages {
		p.Sections = append(p.Sections, types.Section{
			Title: id.String(),
			Body:  pageSummary(page),
		})

		p.Links = append(p.Links, types.PageRef{
			Page: id,
		})
	}

	return
}

func (s Server) pageNotFound(id uuid.UUID) types.Page {
	return s.error(id, "Page Not Found")
}

func (s Server) verbNotSupported(id uuid.UUID) types.Page {
	return s.error(id, "Verb Not Supported")
}

func (s Server) error(id uuid.UUID, msg string) types.Page {
	return types.Page{
		Title:  msg,
		Status: types.StatusError,
		Meta: types.Metadata{
			ID:        id,
			Author:    "Gordon",
			Published: time.Now(),
		},
	}

}

func pageSummary(p types.Page) string {
	return fmt.Sprintf("%s\n\nPublished Last by %s (%s)\n", p.Preamble, p.Meta.Author, p.Meta.Published)
}
