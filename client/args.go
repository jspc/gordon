package main

import (
	"errors"
	"net"
	"net/url"
	"strings"

	"github.com/gofrs/uuid/v5"
	"github.com/jspc/gordon/types"
)

type Address struct {
	addr  *net.UDPAddr
	docID uuid.UUID
}

func ParseAddress(s string) (a Address, err error) {
	u, err := url.Parse(s)
	if err != nil {
		return
	}

	a.addr, err = net.ResolveUDPAddr("udp", u.Host)
	if err != nil {
		return
	}

	trimmedPath := strings.TrimPrefix(u.Path, "/")
	if len(trimmedPath) == 0 {
		return
	}

	a.docID, err = uuid.FromString(trimmedPath)

	return
}

func ParseVerb(s string) (types.Verb, error) {
	switch strings.ToLower(s) {
	case "create":
		return types.VerbCreate, nil

	case "read":
		return types.VerbRead, nil

	case "update":
		return types.VerbUpdate, nil

	case "delete":
		return types.VerbDelete, nil
	}

	return types.VerbUnknown, errors.New("Unknown or invalid verb " + s)
}
