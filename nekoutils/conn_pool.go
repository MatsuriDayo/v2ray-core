package nekoutils

import (
	"net"
	"sync"
	"sync/atomic"

	"github.com/v2fly/v2ray-core/v5/common"
)

type connectionPool struct {
	sync.Map
	cnt uint32
}

var ConnectionPool_System = &connectionPool{sync.Map{}, 0}
var ConnectionPool_V2Ray = &connectionPool{sync.Map{}, 0}
var ConnectionLog_V2Ray = &connectionPool{sync.Map{}, 0}
var Connection_V2Ray_Enabled = false

// For one conn

func (p *connectionPool) AddConnection(c ManagedConn) {
	p.Store(c.ID(), c)
}

func (p *connectionPool) RemoveConnection(c ManagedConn) {
	p.Delete(c.ID())
}

// For all conn

func (p *connectionPool) ResetConnections(reset bool) {
	p.Range(func(key interface{}, value interface{}) bool {
		p.Delete(key)
		if reset {
			common.Close(value)
		}
		return true
	})
	p.cnt = 0
}

// conn

type ManagedConn interface {
	ID() uint32
	Close() error
	RemoteAddress() string
}

type mangedNetConn struct {
	net.Conn //wtf type

	id     uint32
	Closed int32
	Pool   *connectionPool
}

func (c *mangedNetConn) Close() error {
	cnt := atomic.AddInt32(&c.Closed, 1)
	if cnt > 1 { // already closed
		return nil
	}
	c.Pool.RemoveConnection(c)
	return c.Conn.Close()
}

func (c *mangedNetConn) RemoteAddress() string {
	return c.RemoteAddr().Network() + ":" + c.RemoteAddr().String()
}

func (c *mangedNetConn) ID() uint32 {
	return c.id
}

// packet conn?

var _ FusedConn = (*mangedFusedConn)(nil)

type mangedFusedConn struct {
	mangedNetConn
	c2 net.PacketConn
}

func (c *mangedFusedConn) WriteTo(p []byte, d net.Addr) (int, error) {
	return c.c2.WriteTo(p, d)
}

func (c *mangedFusedConn) ReadFrom(p []byte) (int, net.Addr, error) {
	return c.c2.ReadFrom(p)
}

// 在此添加连接

func (p *connectionPool) StartNetConn(c net.Conn) net.Conn {
	mc := mangedNetConn{
		id:   atomic.AddUint32(&ConnectionPool_System.cnt, 1),
		Conn: c,
		Pool: p,
	}

	// PacketConn -> FusedConn
	// Conn -> Conn
	// buf.Copy call ReadFrom if have, and fails if error returned, so do a check here

	if c2, ok := c.(net.PacketConn); ok {
		mfc := mangedFusedConn{
			mangedNetConn: mc,
			c2:            c2,
		}
		p.AddConnection(&mfc)
		return &mfc
	} else {
		p.AddConnection(&mc)
		return &mc
	}
}
