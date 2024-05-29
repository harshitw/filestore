package p2p

import (
	"fmt"
	"net"
	"sync"
)

// TCPPeer represents the remote node over TCP established
// connection.
type TCPPeer struct {
	// conn is the underlying connection of the peer
	conn net.Conn

	// If we dial and retrieve a conn =>  outbound == true
	// If we accept and retrieve a conn =>  outbound == false
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

type TCPTransport struct {
	listenAddress string
	listener      net.Listener
	handshake     Handshaker

	// put mutex over the resource you want to protect
	mu    sync.RWMutex
	peers map[net.Addr]Peer
}

func NewTCPTransport(listenAddr string) *TCPTransport {
	return &TCPTransport{
		listenAddress: listenAddr,
	}
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error
	t.listener, err = net.Listen("tcp", t.listenAddress)
	if err != nil {
		return err
	}
	go t.startAcceptLoop()
	fmt.Printf("listen and accept")
	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("Tcp accept error %s\n", err)
		}
		fmt.Printf("start accept loop")
		go t.handleConn(conn)
	}
}

func (t *TCPTransport) handleConn(conn net.Conn) {
	peer := NewTCPPeer(conn, true)
	fmt.Printf("new incomming connection %+v\n", peer)
}

// func Tset() {
// 	t := NewTCPTransport(":4343").(*TCPTransport)
// 	t.listener.Accept()
// }
