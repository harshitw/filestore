package p2p

import (
	"fmt"
	"net"
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

// close implements the peer interface
func (p *TCPPeer) Close() error {
	return p.conn.Close()
}

type TCPTransportOpts struct {
	ListenAddr    string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder
	OnPeer        func(Peer) error // If this function returns error we
	// will drop this, peer list will be maintained on server
}

type TCPTransport struct {
	listener net.Listener
	TCPTransportOpts
	rpcCh chan RPC

	// put mutex over the resource you want to protect
	// mu    sync.RWMutex
	// peers map[net.Addr]Peer
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
		rpcCh:            make(chan RPC),
	}
}

// can only read from the channels, cannot send to the channel
func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcCh
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error
	t.listener, err = net.Listen("tcp", t.ListenAddr)
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
			fmt.Printf("TCP accept error %s\n", err)
		}
		fmt.Printf("new incomming connection %s\n", conn)
		go t.handleConn(conn)
	}
}

func (t *TCPTransport) handleConn(conn net.Conn) {

	var err error
	defer func() {
		fmt.Printf("dropping peer connection : %s", err)
		conn.Close()
	}()

	peer := NewTCPPeer(conn, true)
	if err := t.HandshakeFunc(peer); err != nil {
		fmt.Printf("TCP handshake error %s\n", err)
		conn.Close()
		return
	}

	if t.OnPeer != nil {
		if err := t.OnPeer(peer); err != nil {
			return
		}
	}

	// Read loop
	rpc := RPC{}

	for {
		// Not decoding of payload but RPC between two peers
		if err := t.Decoder.Decode(conn, &rpc); err != nil {
			fmt.Printf("TCP error %s\n", err)
			continue
		}
		rpc.From = conn.RemoteAddr()

		fmt.Printf("message : +%v\n", rpc)
		t.rpcCh <- rpc
	}
}
