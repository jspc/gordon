package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/jspc/gordon/types"
	"github.com/pion/dtls/v2"
	"github.com/pion/dtls/v2/pkg/crypto/selfsign"
)

func main() {
	fmt.Println("gordon")

	gordon := Server{
		pages: map[uuid.UUID]types.Page{
			rootDocID: rootDoc,
			mintDocID: mintDoc,
		},
	}

	addr := &net.UDPAddr{IP: net.ParseIP("0.0.0.0"), Port: 4444}
	certificate, err := selfsign.GenerateSelfSigned()
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config := &dtls.Config{
		Certificates:         []tls.Certificate{certificate},
		ExtendedMasterSecret: dtls.RequireExtendedMasterSecret,
		ConnectContextMaker: func() (context.Context, func()) {
			return context.WithTimeout(ctx, time.Second*5)
		},
	}

	server, err := dtls.Listen("udp", addr, config)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := server.Accept()
		if err != nil {
			panic(err)
		}

		go func(c net.Conn) {
			data := make([]byte, 1024)
			_, err := c.Read(data)
			if err != nil {
				panic(err)
			}

			buf := bytes.NewBuffer(data)
			req := new(types.Request)

			err = req.Unmarshall(buf)
			if err != nil {
				panic(err)
			}

			resp, err := gordon.Serve(req)
			if err != nil {
				panic(err)
			}

			buf = new(bytes.Buffer)
			err = resp.Marshall(buf)
			if err != nil {
				panic(err)
			}

			_, err = c.Write(buf.Bytes())
			if err != nil {
				panic(err)
			}
		}(conn)
	}
}
