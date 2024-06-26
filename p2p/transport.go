package p2p

// Peer is an interface the represents the remote node.
type Peer interface {
	Close() error
}

// Transport is an interface that handles the communication
// between nodes in the network. Can be of the form (TCP,
// UDP, websockets..)
type Transport interface {
	ListenAndAccept() error
	Consume() <-chan RPC
}
