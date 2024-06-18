package main

import (
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/jspc/gordon/types"
)

var (
	rootDocID = uuid.FromStringOrNil("208b43d9-a95d-476d-ba3b-3b64fda2507b")
	mintDocID = uuid.FromStringOrNil("996b046f-11d2-41c9-8b45-9294c7215e38")
)

var rootDoc = types.Page{
	Meta: types.Metadata{
		ID:        rootDocID,
		Author:    "jspc",
		Published: time.Now(),
	},
	History:  make([]types.Metadata, 0),
	Title:    "The Gordon Documentation Protocol",
	Preamble: "The Canonical Gordon Documentation",
	Sections: []types.Section{
		{
			Title: "Introduction",
			Body: `Welcome to the Gordon Documentation Protocol, umm...., documentation.
This piece of documentation is a living document, and reflects the absolute bleeding edge.
`,
		},
		{
			Title: "The Data Format",
			Body: `Gordon pages are defined in the mint format.

Mint is a binary data format that produces small payloads and supports uuids, datetimes, and stuff natively. It has a DDL and code generators, too and may be found at https://github.com/vinyl-linux/mint.

The definition for Gordon documents can be found in the git repo, or at [l:0].
`,
		},
		{
			Title: "The Network Format",
			Body:  `Gordon pages are served via DTLS and uses a default max payload size of 256kb. Documents larger that this will work perfectly fine, they might just be a little slower to load.`,
		},
	},
	Tags: []string{"gordon", "protocol"},
	Labels: map[string]string{
		"status": "draft",
	},
	Links: []types.PageRef{
		{
			Page: mintDocID,
		},
	},
	Relationships: []types.Relationship{
		{
			Subject: types.PageRef{
				Page: mintDocID,
			},
			Predicate: types.PredicateSupplements,
			Object: types.PageRef{
				Page: rootDocID,
			},
		},
	},
	Status: types.StatusOK,
}

var mintDoc = types.Page{
	Meta: types.Metadata{
		ID:        mintDocID,
		Author:    "jspc",
		Published: time.Now(),
	},
	History:  make([]types.Metadata, 0),
	Title:    "Mint DDL for Gordon",
	Preamble: "Mint DDL for Gordon",
	Sections: []types.Section{
		{
			Title: "page.mint",
			Body: `type Page {
     +mint:doc:"Metadata contains metadata about this page, including"
     +mint:doc:"useful stuff like author details, and revisions, and dates"
     +mint:doc:"and all that stuff"
     Metadata Meta = 0;

     +mint:doc:"History is a slice of Metadatas that point to the previous"
     +mint:doc:"versions of this page, such as they may be."
     +mint:doc:"There are no limits to this field, and indeed no requirement"
     +mint:doc:"for it to even contain anything"
     []Metadata History = 1;

     +mint:doc:"Title represents the title of a page, funnily enough, and"
     +mint:doc:"must be between 1 and 512 characters"
     +mint:validate:string_not_empty
     +custom:validate:title_not_too_long
     string Title = 2;

     +mint:doc:"Sections contains a list of sections in a document. A section"
     +mint:doc:"is how we divvy up data"
     []Section Sections = 3;

     +mint:doc:"Tags are an arbitrary list of strings to attach to a page for"
     +mint:doc:"searching and organising"
     []string Tags = 4;

     +mint:doc:"Labels are Much like tags; handy for searching and organising"
     +mint:doc:"data"
     map<string,string> Labels = 5;

     +mint:doc:"Links are references to other pages; within the body of a section,"
     +mint:doc:"they are referenced by their index"
     []Link Links = 6;

     +mint:doc:"Status reflects whether this page is to be treated as an error page"
     +mint:doc:"or not"
     Status Status = 7;

     +mint:doc:"Preamble is used on Page Indexes and other list operations"
     string Preamble = 8;
}

type Metadata {
     +mint:doc:"ID of this resource"
     uuid ID = 0;

     +mint:doc:"Author of the resource"
     string Author = 1;

     +mint:doc:"Published date of the resource"
     +mint:transform:date_in_utc
     datetime Published = 2;
}

type Section {
     string Title = 0;
     string Body = 1;
}

type Link {
     uuid Page = 0;
     string Section = 1;
     string DisplayText = 2;
}

enum Status {
     OK
     Error
}`,
		},
		{
			Title: "request.mint",
			Body: `enum Verb {
     Create
     Read
     Update
     Delete
}

type Request {
     +mint:doc:"Verb is analogous to an HTTP request type and defines the"
     +mint:doc:"expectations the client has for the type of response from"
     +mint:doc:"the server"
     Verb Verb = 0;

     +mint:doc:"ID is the ID of the Page being operated on"
     uuid ID = 1;

     +mint:doc:"Args are arbitrary arguments that a server may or may not"
     +mint:doc:"respond to."
     +mint:doc:""
     +mint:doc:"There are some arguments that will always be expected for verbs"
     +mint:doc:"such as Section for an Update, or Body for both Create and Update"
     map<string,string> Args = 2;
}
`,
		},
	},
	Tags: []string{"gordon", "protocol", "mint"},
	Labels: map[string]string{
		"status": "draft",
	},
	Relationships: []types.Relationship{
		{
			Subject: types.PageRef{
				Page: mintDocID,
			},
			Predicate: types.PredicateSupplements,
			Object: types.PageRef{
				Page: rootDocID,
			},
		},
	},
	Status: types.StatusOK,
}
