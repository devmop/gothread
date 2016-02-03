package shutdown

import (
	"errors"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var shutdown error

func init() {
	shutdown = errors.New("Shutting down...")
}

type stoppable struct {
	ch chan struct{}
	ListenerWithDeadline
}

type ListenerWithDeadline interface {
	SetDeadline(time.Time) error
	net.Listener
}

func Wrap(l net.Listener) net.Listener {
	lwd, ok := l.(ListenerWithDeadline)

	if !ok {
		return l
	}

	ch := make(chan struct{})

	go awaitQuit(ch)

	return stoppable{ch, lwd}
}

func WasUserShutdown(e error) bool {
	return shutdown == e
}

func (s stoppable) Accept() (net.Conn, error) {
	for {

		select {
		case _ = <-s.ch:
			return nil, shutdown
		default:
		}

		s.ListenerWithDeadline.SetDeadline(time.Now().Add(500 * time.Millisecond))

		conn, err := s.ListenerWithDeadline.Accept()

		if isTimeout(err) {
			continue
		}

		return conn, err
	}
}

func isTimeout(err error) bool {
	if err == nil {
		return false
	}

	neterr, ok := err.(net.Error)

	return ok && neterr.Timeout()
}

func awaitQuit(ch chan struct{}) {
	signals := make(chan os.Signal, 100)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)

	<-signals
	ch <- struct{}{}
}
