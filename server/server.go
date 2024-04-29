package server

import (
	"fmt"
	"github.com/dariusigna/bitcoin-handshake/peer"
	"github.com/dariusigna/bitcoin-handshake/protocol"
	"net"
	"sync"
)

// Server implements a light Bitcoin node that creates an outbound network handshake with a target node.
type Server struct {
	network      string
	networkMagic protocol.Magic
	userAgent    string
	outboundPeer *peer.Peer
	wg           sync.WaitGroup
	targetAddr   string
}

// New returns a new Server.
func New(network, userAgent, targetAddr string) (*Server, error) {
	networkMagic, ok := protocol.Networks[network]
	if !ok {
		return nil, fmt.Errorf("unsupported network %s", network)
	}

	return &Server{
		network:      network,
		networkMagic: networkMagic,
		userAgent:    userAgent,
		targetAddr:   targetAddr,
	}, nil
}

func (s *Server) Start() {
	s.wg.Add(1)
	go s.peerHandler()
}

func (s *Server) Stop() {
	if s.outboundPeer != nil {
		s.outboundPeer.Disconnect()
	}
}

func (s *Server) WaitForShutdown() {
	s.wg.Wait()
}

func (s *Server) peerHandler() {
	err := s.createOutBoundConnectionPeer(s.targetAddr)
	if err != nil {
		fmt.Printf("Cannot connect to peer: %v", err)
	}

	s.wg.Done()
}

func (s *Server) peerDoneHandler(p *peer.Peer) {
	p.WaitForDisconnect()
}

func (s *Server) createOutBoundConnectionPeer(targetAddr string) error {
	conn, err := net.Dial("tcp", targetAddr)
	if err != nil {
		return err
	}

	s.outboundPeer = peer.New(conn, false, targetAddr, s.network)
	s.outboundPeer.SetConnection(conn)
	go s.peerDoneHandler(s.outboundPeer)
	return nil
}
