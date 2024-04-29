package peer

import (
	"errors"
	"github.com/sirupsen/logrus"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

const (
	// negotiateTimeout is the duration of inactivity before we timeout a
	// peer that hasn't completed the initial version negotiation.
	negotiateTimeout = 30 * time.Second
)

type Peer struct {
	conn            net.Conn
	outbound        bool
	address         string
	quit            chan struct{}
	disconnect      int32
	connected       int32
	network         string
	protocolVersion uint32
	versionKnown    bool
	verAckReceived  bool
	flagsMtx        sync.RWMutex
}

func New(conn net.Conn, outbound bool, address, network string) *Peer {
	return &Peer{
		conn:     conn,
		outbound: outbound,
		address:  address,
		network:  network,
		quit:     make(chan struct{}),
	}
}

func (p *Peer) SetConnection(conn net.Conn) {
	p.conn = conn
	if !atomic.CompareAndSwapInt32(&p.connected, 0, 1) {
		return
	}

	go func() {
		if err := p.start(); err != nil {
			logrus.Debugf("Cannot start peer: %v", err)
			p.Disconnect()
		}
	}()
}

func (p *Peer) start() error {
	logrus.Info("Starting peer")
	negotiateErr := make(chan error, 1)
	go func() {
		if p.outbound {
			negotiateErr <- p.negotiateOutboundProtocol()
		}
	}()

	// Negotiate the protocol within the specified negotiateTimeout.
	select {
	case err := <-negotiateErr:
		if err != nil {
			p.Disconnect()
			return err
		}
	case <-time.After(negotiateTimeout):
		p.Disconnect()
		return errors.New("protocol negotiation timeout")
	}

	logrus.Debugf("Connected to %s", p.Addr())
	// Keep the connection alive for 5 more seconds.
	time.Sleep(5 * time.Second)
	p.Disconnect()
	return nil
}

func (p *Peer) Addr() string {
	return p.address
}

func (p *Peer) WaitForDisconnect() {
	<-p.quit
}

func (p *Peer) Disconnect() {
	if atomic.AddInt32(&p.disconnect, 1) != 1 {
		return
	}

	logrus.Debugf("Disconnecting %s", p.address)
	if atomic.LoadInt32(&p.connected) != 0 {
		p.conn.Close()
	}

	close(p.quit)
}
