package peer

import (
	"bytes"
	"fmt"
	"github.com/dariusigna/bitcoin-handshake/encoding"
	"github.com/dariusigna/bitcoin-handshake/protocol"
	"github.com/sirupsen/logrus"
	"io"
	"math/rand"
	"time"
)

func (p *Peer) negotiateOutboundProtocol() error {
	if err := p.writeLocalVersionMsg(); err != nil {
		return err
	}

	if err := p.readRemoteVersionMsg(); err != nil {
		return err
	}

	if err := p.writeVerackMsg(); err != nil {
		return err
	}

	return p.waitToFinishNegotiation()
}

func (p *Peer) readMessage() (*protocol.MessageHeader, error) {
	tmp := make([]byte, protocol.MsgHeaderLength)
	byteCount, err := p.conn.Read(tmp)
	if err != nil {
		return nil, err
	}

	logrus.Debugf("received header: %x", tmp[:byteCount])
	var msgHeader protocol.MessageHeader
	if err = encoding.NewDecoder(bytes.NewReader(tmp[:byteCount])).Decode(&msgHeader); err != nil {
		return nil, fmt.Errorf("invalid header: %s", err)
	}

	if err = msgHeader.Validate(); err != nil {
		return nil, fmt.Errorf("invalid header: %s", err)
	}

	return &msgHeader, nil
}

func (p *Peer) readRemoteVersionMsg() error {
	msg, err := p.readMessage()
	if err != nil {
		return err
	}

	if msg.CommandString() != "version" {
		return fmt.Errorf("unexpected command %s -> expected 'version", msg.CommandString())
	}

	var version protocol.MsgVersion
	lr := io.LimitReader(p.conn, int64(msg.Length))
	if err = encoding.NewDecoder(lr).Decode(&version); err != nil {
		return err
	}

	p.flagsMtx.Lock()
	p.protocolVersion = uint32(version.Version)
	p.versionKnown = true
	p.flagsMtx.Unlock()

	logrus.Infof("VERSION: %+v", version.Version)
	return nil
}

func (p *Peer) writeVerackMsg() error {
	verack, err := protocol.NewVerackMsg(p.network)
	if err != nil {
		return err
	}

	msg, err := encoding.Marshal(verack)
	if err != nil {
		return err
	}

	_, err = p.conn.Write(msg)
	if err != nil {
		return err
	}

	return nil
}

func (p *Peer) writeLocalVersionMsg() error {
	peerAddr, err := ParseNodeAddr(p.address)
	if err != nil {
		return err
	}

	version := protocol.MsgVersion{
		Version:   protocol.Version,
		Services:  protocol.SrvNodeNetwork,
		Timestamp: time.Now().UTC().Unix(),
		AddrRecv: protocol.VersionNetAddr{
			Services: protocol.SrvNodeNetwork,
			IP:       peerAddr.IP,
			Port:     peerAddr.Port,
		},
		AddrFrom: protocol.VersionNetAddr{
			Services: protocol.SrvNodeNetwork,
			IP:       protocol.NewIPv4(127, 0, 0, 1),
			Port:     9334,
		},
		Nonce:       nonce(),
		UserAgent:   protocol.NewUserAgent("Darius Awesome Node"),
		StartHeight: -1,
		Relay:       true,
	}

	msg, err := protocol.NewMessage("version", p.network, version)
	if err != nil {
		logrus.Fatalln(err)
	}

	msgSerialized, err := encoding.Marshal(msg)
	if err != nil {
		logrus.Fatalln(err)
	}

	_, err = p.conn.Write(msgSerialized)
	if err != nil {
		logrus.Fatalln(err)
	}

	return nil
}

func (p *Peer) processRemoteVerAckMsg() {
	p.flagsMtx.Lock()
	p.verAckReceived = true
	p.flagsMtx.Unlock()
}

func (p *Peer) waitToFinishNegotiation() error {
	for {
		msg, err := p.readMessage()
		if err != nil {
			return err
		}

		switch msg.CommandString() {
		case "verack":
			p.processRemoteVerAckMsg()
			return nil
		default:
			return fmt.Errorf("unexpected command %s", msg.CommandString())
		}
	}
}

func nonce() uint64 {
	return rand.Uint64()
}
