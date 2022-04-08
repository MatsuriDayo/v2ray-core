package nekoutils

import (
	"net"
)

// Compact: udp & conn_pool

type FusedConn interface {
	net.PacketConn
	net.Conn
}

// Wrap a net.PacketConn to a net.Conn or FusedConn
// Different from which in conncetion_adaptor.go

var _ FusedConn = (*PacketConnWrapper)(nil)

type PacketConnWrapper struct {
	net.PacketConn
	Dest net.Addr
}

func (c *PacketConnWrapper) RemoteAddr() net.Addr {
	return c.Dest
}

func (c *PacketConnWrapper) Write(p []byte) (int, error) {
	return c.PacketConn.WriteTo(p, c.Dest)
}

func (c *PacketConnWrapper) Read(p []byte) (int, error) {
	n, _, err := c.PacketConn.ReadFrom(p)
	return n, err
}
