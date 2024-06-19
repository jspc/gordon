package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"net"
	"time"

	"github.com/jspc/gordon/types"
	"github.com/pion/dtls/v2"
	"go.uber.org/zap"
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
	logger         *zap.Logger
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
	l.logger, err = zap.NewProduction()

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
			l.logger.Error(err.Error())

			return err
		}

		err = l.requestPool.Acquire(context.Background(), 1)
		if err != nil {
			l.connErr(conn, err)

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
		l.connErr(conn, err)

		return
	}

	buf := bytes.NewBuffer(data)
	req := new(types.Request)

	err = req.Unmarshall(buf)
	if err != nil {
		l.connErr(conn, err)

		return
	}

	resp, err := l.handler.Serve(req)
	if err != nil {
		l.connErr(conn, err)

		return
	}

	buf = new(bytes.Buffer)
	err = resp.Marshall(buf)
	if err != nil {
		l.connErr(conn, err)

		return
	}

	_, err = conn.Write(buf.Bytes())
	if err != nil {
		l.connErr(conn, err)

		return
	}

	duration := time.Now().Sub(start)
	l.logger.Info("Request",
		zap.String("verb", verbToString(req.Verb)),
		zap.String("document", req.ID.String()),
		zap.String("remote_address", conn.RemoteAddr().String()),
		zap.Duration("duration", duration),
		zap.Bool("is_error", resp.Status == types.StatusError),
		zap.Int("size", buf.Len()),
	)
}

func (l Listener) connErr(conn net.Conn, err error) {
	l.logger.Error(err.Error(),
		zap.Error(err),
		zap.String("RemoteAddress", conn.RemoteAddr().String()),
	)
}

func verbToString(v types.Verb) string {
	switch v {
	case types.VerbCreate:
		return "create"
	case types.VerbRead:
		return "read"
	case types.VerbUpdate:
		return "update"
	case types.VerbDelete:
		return "delete"
	}

	return "unknown"
}
