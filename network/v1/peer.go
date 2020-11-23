package network

import "net"

type Peer interface {
	//service.Service
	FlushStop()

	RemoteIP() net.IP     // remote IP of the connection
	RemoteAddr() net.Addr // remote address of the connection

	IsOutbound() bool   // did we dial the peer
	IsPersistent() bool // do we redial this peer when we disconnect

	CloseConn() error // close original connection

	NodeInfo() NodeInfo // peer's info
	//Status() ConnectionStatus
	SocketAddr() *NetAddress // actual address of the socket

	Send(byte, []byte) bool
	TrySend(byte, []byte) bool

	Set(string, interface{})
	Get(string) interface{}
}

// NetAddress defines information about a peer on the network
// including its ID, IP address, and port.
type NetAddress struct {
	IP   net.IP `json:"ip"`
	Port uint16 `json:"port"`

	str string
}
