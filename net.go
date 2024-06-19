package gordon

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
	defaultMaxWorkers int64 = 1024
	networkUDP              = "udp"
)

// A Handler responds to Gordon Requests with either a Page or an error
//
// This is the absolute minimum a Gordon Server implementation must provide
// to be able to serve documentation to users
type Handler interface {
	Serve(req *types.Request) (resp *types.Page, err error)
}

// A Listener wraps a Handler and a TLS Certificate and does all of the
// networky stuff
//
// This type should not be copied; it contains (amongst other things) a
// reference to a weighted semaphore- if you try to copy it or duplicate
// it then you'll end up with very strange behaviour
type Listener struct {
	handler        Handler
	listener       net.Listener
	listenerConfig *dtls.Config
	logger         *zap.Logger
	requestPool    *semaphore.Weighted

	MaxConnections int64
}

// NewListener accepts a Handler and a Certificate and configures a Listener
// ahead of `ListenAndServe` being called
//
// This call sets up a default MaxConnections size which may be overwritten
// afterwards, as per
//
//	l, _ := gordon.NewListener(someServer{}, someCert)
//	l.MaxConnections = 10
//
// The default size is 1024, which may be stupidly, ridiculously, overly
// large.
func NewListener(h Handler, cert tls.Certificate) (l Listener, err error) {
	l.MaxConnections = defaultMaxWorkers

	l.handler = h

	ctx := context.Background()
	l.listenerConfig = &dtls.Config{
		Certificates:         []tls.Certificate{cert},
		ExtendedMasterSecret: dtls.RequireExtendedMasterSecret,
		ConnectContextMaker: func() (context.Context, func()) {
			return context.WithTimeout(ctx, time.Second*5)
		},
	}

	l.logger, err = zap.NewProduction()

	return
}

// ListenAndServe creates a DTLS listener against the provided address,
// and handles Marshalling and Unmarshalling mint documents into Gordon
// types.
//
// This function will propagate errors creating a DTLS listener to the
// gordon implementation; any error in processing data, or any error returned
// from a Handler, is logged and moved on from.
func (l *Listener) ListenAndServe(address string) (err error) {
	addr, err := net.ResolveUDPAddr(networkUDP, address)
	if err != nil {
		return
	}

	l.listener, err = dtls.Listen(networkUDP, addr, l.listenerConfig)
	if err != nil {
		return
	}

	l.requestPool = semaphore.NewWeighted(l.MaxConnections)

	for {
		conn, err := l.listener.Accept()
		if err != nil {
			l.logger.Error(err.Error())

			return err
		}

		err = l.requestPool.Acquire(context.Background(), 1)
		if err != nil {
			l.connErr(conn, err)

			conn.Close()

			return err
		}

		go l.process(conn)
	}
}

// Close the underlying UDP listener
func (l Listener) Close() error {
	return l.listener.Close()
}

func (l *Listener) process(conn net.Conn) {
	defer l.requestPool.Release(1)
	defer conn.Close()

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

	if resp == nil {
		l.connErr(conn, new(NilPageError))

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

func isValidVerb(v types.Verb) bool {
	return v > types.VerbUnknown &&
		v <= types.VerbDelete
}
