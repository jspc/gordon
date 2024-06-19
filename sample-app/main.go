package main

import (
	"fmt"

	"github.com/gofrs/uuid/v5"
	"github.com/jspc/gordon"
	"github.com/jspc/gordon/types"
	"github.com/pion/dtls/v2/pkg/crypto/selfsign"
)

func main() {
	fmt.Println("gordon")

	s := Server{
		pages: map[uuid.UUID]*types.Page{
			rootDocID: rootDoc,
			mintDocID: mintDoc,
		},
	}

	certificate, err := selfsign.GenerateSelfSigned()
	if err != nil {
		panic(err)
	}

	l, err := gordon.NewListener(s, certificate)
	if err != nil {
		panic(err)
	}

	panic(l.ListenAndServe("0.0.0.0:4444"))
}
