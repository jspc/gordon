package main

import (
	"flag"
	"os"

	"github.com/kr/pretty"
)

var (
	verb = flag.String("X", "READ", "Verb to use during rquest, case insensitive")
)

func main() {
	flag.Parse()

	if len(flag.Args()) != 1 {
		panic("invalid command; call " + os.Args[0] + " [-X READ] //example.com/some-document")
	}

	addr, err := ParseAddress(flag.Arg(0))
	if err != nil {
		panic(err)
	}

	v, err := ParseVerb(*verb)
	if err != nil {
		panic(err)
	}

	page, err := DoRequest(v, addr)
	if err != nil {
		panic(err)
	}

	pretty.Print(page)
}
