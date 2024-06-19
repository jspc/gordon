package client

import (
	"errors"
	"net"
	"net/url"
	"strings"

	"github.com/gofrs/uuid/v5"
	"github.com/jspc/gordon/types"
)

type Address struct {
	orig string

	addr  *net.UDPAddr
	docID uuid.UUID
}

func (a Address) String() string {
	return a.orig
}

func (a Address) Server() string {
	return a.addr.AddrPort().String()
}

func (a Address) Page() string {
	return a.docID.String()
}

func ParseAddress(s string) (a Address, err error) {
	a.orig = s

	u, err := url.Parse(s)
	if err != nil {
		return
	}

	// Set default port to 4444
	if !strings.Contains(u.Host, ":") {
		u.Host += ":4444"
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
