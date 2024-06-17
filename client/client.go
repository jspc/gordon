package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"time"

	"github.com/jspc/gordon/types"
	"github.com/pion/dtls/v2"
	"github.com/pion/dtls/v2/pkg/crypto/selfsign"
)

func DoRequest(verb types.Verb, addr Address) (page *types.Page, err error) {
	req := types.Request{
		Verb: verb,
		ID:   addr.docID,
	}

	buf := new(bytes.Buffer)

	err = req.Marshall(buf)
	if err != nil {
		return
	}

	certificate, err := selfsign.GenerateSelfSigned()
	if err != nil {
		panic(err)
	}

	config := &dtls.Config{
		Certificates:         []tls.Certificate{certificate},
		InsecureSkipVerify:   true,
		ExtendedMasterSecret: dtls.RequireExtendedMasterSecret,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	conn, err := dtls.DialWithContext(ctx, "udp", addr.addr, config)
	if err != nil {
		return
	}

	_, err = conn.Write(buf.Bytes())
	if err != nil {
		return
	}

	payload := make([]byte, 0)
	for {
		// Work in buffers of 256kb
		data := make([]byte, 256000)

		read, err := conn.Read(data)
		if err != nil {
			return nil, err
		}

		payload = append(payload, data...)
		if read < 256000 {
			break
		}
	}

	buf = bytes.NewBuffer(payload)
	page = new(types.Page)

	err = page.Unmarshall(buf)

	return
}
