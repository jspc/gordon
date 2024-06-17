enum Verb {
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
     +mint:doc:"such as `Section` for an Update, or `Body` for both Create and Update"
     map<string,string> Args = 2;
}
