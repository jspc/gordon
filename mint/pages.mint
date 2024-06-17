type Page {
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

     +mint:doc:"Preamble is used on Page Indexes and other list operations"
     string Preamble = 3;

     +mint:doc:"Sections contains a list of sections in a document. A section"
     +mint:doc:"is how we divvy up data"
     []Section Sections = 4;

     +mint:doc:"Tags are an arbitrary list of strings to attach to a page for"
     +mint:doc:"searching and organising"
     []string Tags = 5;

     +mint:doc:"Labels are Much like tags; handy for searching and organising"
     +mint:doc:"data"
     map<string,string> Labels = 6;

     +mint:doc:"Links are references to other pages; within the body of a section,"
     +mint:doc:"they are referenced by their index"
     []PageRef Links = 7;

     +mint:doc:"Relationships are used to link pages"
     []Relationship Relationships = 8;

     +mint:doc:"Status reflects whether this page is to be treated as an error page"
     +mint:doc:"or not"
     Status Status = 9;
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

type PageRef {
     uuid Page = 0;
     string Section = 1;
     string Server = 2;
}

enum Status {
     OK
     Error
}

type Relationship {
     PageRef Subject = 0;
     Predicate Predicate = 1;
     PageRef Object = 2;
}

enum Predicate {
     HasChild
     Extends
     Supercedes
     Supplements
}
