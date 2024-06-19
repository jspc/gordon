package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"log"
	"net"
	"time"

	"github.com/jspc/gordon/types"
	"github.com/pion/dtls/v2"
	"golang.org/x/sync/semaphore"
)

const (
	maxWorkers int64 = 1024
	networkUDP       = "udp"
)

type Handler interface {
	Serve(req *types.Request) (resp types.Page, err error)
}

type Listener struct {
	handler        Handler
	requestPool    *semaphore.Weighted
	listenerConfig *dtls.Config
}

func NewListener(h Handler, cert tls.Certificate) (l Listener, err error) {
	ctx := context.Background()

	l.handler = h
	l.listenerConfig = &dtls.Config{
		Certificates:         []tls.Certificate{cert},
		ExtendedMasterSecret: dtls.RequireExtendedMasterSecret,
		ConnectContextMaker: func() (context.Context, func()) {
			return context.WithTimeout(ctx, time.Second*5)
		},
	}

	l.requestPool = semaphore.NewWeighted(maxWorkers)

	return
}

func (l *Listener) ListenAndServe(address string) (err error) {
	addr, err := net.ResolveUDPAddr(networkUDP, address)
	if err != nil {
		return
	}

	s, err := dtls.Listen(networkUDP, addr, l.listenerConfig)
	if err != nil {
		return
	}

	for {
		conn, err := s.Accept()
		if err != nil {
			return err
		}

		err = l.requestPool.Acquire(context.Background(), 1)
		if err != nil {
			connErr(conn, err)

			conn.Close()
		}

		go l.process(conn)
	}
}

func (l *Listener) process(conn net.Conn) {
	defer l.requestPool.Release(1)

	start := time.Now()

	data := make([]byte, 1024)
	_, err := conn.Read(data)
	if err != nil {
		connErr(conn, err)

		return
	}

	buf := bytes.NewBuffer(data)
	req := new(types.Request)

	err = req.Unmarshall(buf)
	if err != nil {
		connErr(conn, err)

		return
	}

	resp, err := l.handler.Serve(req)
	if err != nil {
		connErr(conn, err)

		return
	}

	buf = new(bytes.Buffer)
	err = resp.Marshall(buf)
	if err != nil {
		connErr(conn, err)

		return
	}

	_, err = conn.Write(buf.Bytes())
	if err != nil {
		connErr(conn, err)

		return
	}

	duration := time.Now().Sub(start)
	log.Printf("%q %v %q %q",
		conn.RemoteAddr().String(),
		req.Verb,
		req.ID,
		duration.String(),
	)
}

func connErr(conn net.Conn, err error) {
	log.Printf("%s %s", conn.RemoteAddr().String(), err.Error())
}
