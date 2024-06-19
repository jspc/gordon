package gordon

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/jspc/gordon/client"
	"github.com/jspc/gordon/types"
	"github.com/pion/dtls/v2"
	"github.com/pion/dtls/v2/pkg/crypto/selfsign"
)

type dummyHandler struct {
	err bool
}

func (h dummyHandler) Serve(*types.Request) (*types.Page, error) {
	if h.err {
		return nil, errors.New("An error")
	}

	return &types.Page{
		Title: "A Test Page",
		Meta: types.Metadata{
			Author: "A. N. Tester",
		},
		Sections: []types.Section{
			{
				Title: "Chapter 1 - I am born",
				Body:  "some text nobody will read",
			},
		},
		Status: types.StatusOK,
	}, nil
}

type nilNilHandler struct{}

func (nilNilHandler) Serve(*types.Request) (*types.Page, error) {
	return nil, nil
}

type emptyPageHandler struct{}

func (emptyPageHandler) Serve(*types.Request) (*types.Page, error) {
	return new(types.Page), nil
}

func TestNewListener(t *testing.T) {
	_, err := NewListener(new(dummyHandler), tls.Certificate{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestListener_ListenAndServe(t *testing.T) {
	cert, _ := selfsign.GenerateSelfSigned()

	l, _ := NewListener(new(dummyHandler), cert)
	defer func() {
		err := l.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()

	go l.ListenAndServe("localhost:4445")

	addr, _ := client.ParseAddress("//localhost:4445/")

	for _, test := range []struct {
		name        string
		handler     Handler
		expectPage  *types.Page
		expectError bool
	}{
		{"Handler returns nothing on errors", dummyHandler{err: true}, nil, true},
		{"Handler returns page when no error", dummyHandler{}, new(types.Page), false},
		{"Nil pages from handler errors appropriately", nilNilHandler{}, nil, true},
		{"Invalid pages from handler errors appropriately", emptyPageHandler{}, nil, true},
	} {
		t.Run(test.name, func(t *testing.T) {
			l.handler = test.handler

			rcvd, err := client.DoRequest(types.VerbRead, addr)
			if err != nil && !test.expectError {
				t.Errorf("unexpected error: %v", err)
			} else if err == nil && test.expectError {
				t.Error("expected error")
			}

			if (test.expectPage == nil) != (rcvd == nil) {
				t.Errorf("expected\n\t%#v\nreceived\n\t%#v",
					test.expectPage,
					rcvd,
				)
			}
		})
	}
}

func TestListener_ListenAndServe_DodgyListenAddress(t *testing.T) {
	cert, _ := selfsign.GenerateSelfSigned()

	l, _ := NewListener(new(dummyHandler), cert)

	err := l.ListenAndServe("kjshkas:sds::::2223")
	if err == nil {
		t.Fatal("expected error, received none")
	}
}

func TestListener_ListenAndServe_NonDTLSAreSilentlyDiscarded(t *testing.T) {
	cert, _ := selfsign.GenerateSelfSigned()

	l, _ := NewListener(new(dummyHandler), cert)
	defer func() {
		err := l.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()

	go l.ListenAndServe("localhost:4445")

	addr, _ := net.ResolveUDPAddr("udp", "localhost:4445")

	c, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		t.Fatal(err)
	}

	defer c.Close()

	_, err = c.Write([]byte("hello, world!"))
	if err != nil {
		t.Fatal(err)
	}

	c.SetDeadline(time.Now().Add(time.Millisecond * 100))

	buf := make([]byte, 1024)
	_, err = c.Read(buf)
	if err == nil {
		t.Errorf("expected error, received none with buffer %q", string(buf))
	}
}

func TestListener_ListenAndServe_NonMintEncodedDocumentFails(t *testing.T) {
	cert, _ := selfsign.GenerateSelfSigned()

	l, _ := NewListener(new(dummyHandler), cert)
	defer func() {
		l.Close()
	}()

	go l.ListenAndServe("localhost:4445")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	addr, _ := net.ResolveUDPAddr("udp", "localhost:4445")
	conn, err := dtls.DialWithContext(ctx, "udp", addr, &dtls.Config{
		Certificates:         []tls.Certificate{cert},
		InsecureSkipVerify:   true,
		ExtendedMasterSecret: dtls.RequireExtendedMasterSecret,
	})
	if err != nil {
		t.Error(err)
	}

	_, err = conn.Write([]byte("hello, world!"))
	if err != nil {
		t.Error(err)
	}

	buf := make([]byte, 1024)
	_, err = conn.Read(buf)
	if err == nil {
		t.Errorf("expected error, received none with buffer %q", string(buf))
	}
}

func TestVerbToString(t *testing.T) {
	for _, test := range []struct {
		name   string
		verb   types.Verb
		expect string
	}{
		{"Creates are recognised", types.VerbCreate, "create"},
		{"Reads are recognised", types.VerbRead, "read"},
		{"Updates are recognised", types.VerbUpdate, "update"},
		{"Deletes are recognised", types.VerbDelete, "delete"},

		{"0 is unknown", types.Verb(0), "unknown"},
		{"5 is unknown", types.Verb(5), "unknown"},
		{"8 is unknown", types.Verb(8), "unknown"},
		{"99 is unknown", types.Verb(99), "unknown"},
	} {
		t.Run(test.name, func(t *testing.T) {
			rcvd := verbToString(test.verb)
			if test.expect != rcvd {
				t.Errorf("expected %q, received %q", test.expect, rcvd)
			}
		})
	}
}

func TestIsValidVerb(t *testing.T) {
	for _, test := range []struct {
		name   string
		verb   types.Verb
		expect bool
	}{
		{"Create is valid", types.VerbCreate, true},
		{"Read is valid", types.VerbRead, true},
		{"Update is valid", types.VerbUpdate, true},
		{"Delete is valid", types.VerbDelete, true},

		{"0 is false", types.Verb(0), false},
		{"5 is false", types.Verb(5), false},
		{"8 is false", types.Verb(8), false},
		{"99 is false", types.Verb(99), false},
	} {
		t.Run(test.name, func(t *testing.T) {
			rcvd := isValidVerb(test.verb)
			if test.expect != rcvd {
				t.Errorf("expected %v, received %v", test.expect, rcvd)
			}
		})
	}

}
