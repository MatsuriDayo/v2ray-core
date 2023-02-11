package nekoutils

import (
	"encoding/json"
	"net"
	"sort"
	"sync"
	"sync/atomic"
)

type connectionPool struct {
	sync.Map
	cnt uint32
}

var ConnectionPool_System = &connectionPool{sync.Map{}, 0}
var ConnectionPool_V2Ray = &connectionPool{sync.Map{}, 0}
var ConnectionLog_V2Ray = &connectionPool{sync.Map{}, 0}
var coreUseConnectionPool sync.Map

func GetConnectionPoolV2RayEnabled(core uintptr) bool {
	if v, ok := coreUseConnectionPool.Load(core); ok {
		return v.(bool)
	}
	return false
}

func SetConnectionPoolV2RayEnabled(core uintptr, enable bool) {
	if enable {
		coreUseConnectionPool.Store(core, true)
	} else {
		coreUseConnectionPool.Delete(core)
	}
}

// For one conn

func (p *connectionPool) AddConnection(c ManagedConn) {
	p.Store(c.ID(), c)
}

func (p *connectionPool) RemoveConnection(c ManagedConn) {
	p.Delete(c.ID())
}

// For all conn

func (p *connectionPool) ResetConnections(corePtr uintptr) {
	p.Range(func(key interface{}, value interface{}) bool {
		c, ok := value.(ManagedConn)

		if !ok {
			return true
		}

		if corePtr == 0 || c.Instance() == corePtr {
			p.Delete(key)
			c.Close()
		}
		return true
	})
}

// conn

type ManagedConn interface {
	ID() uint32
	Instance() uintptr
	Close() error
	RemoteAddress() string
}

type mangedNetConn struct {
	net.Conn //wtf type

	id      uint32
	corePtr uintptr

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

func (c *mangedNetConn) Instance() uintptr {
	return c.corePtr
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

func (p *connectionPool) StartNetConn(c net.Conn, core uintptr) net.Conn {
	mc := mangedNetConn{
		id:      atomic.AddUint32(&ConnectionPool_System.cnt, 1),
		Conn:    c,
		Pool:    p,
		corePtr: core,
	}

	// PacketConn -> FusedConn
	// Conn -> Conn
	// buf.Copy call ReadFrom, and fails if error returned, so do a check here

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

//

var ListConnections_MaxLineCount = 100
var ListConnections_IgnoreTags = func(tag string) bool {
	if tag == "dns-out" || tag == "direct" {
		return true
	}
	return false
}

func ListConnections(corePtr uintptr) string {
	list2 := make([]interface{}, 0)

	rangeMap := func(m *sync.Map) []interface{} {
		vs := make(map[uint32]interface{}, 0)
		ks := make([]uint32, 0)

		m.Range(func(key interface{}, value interface{}) bool {
			if k, ok := key.(uint32); ok {
				vs[k] = value
				ks = append(ks, k)
			}
			return true
		})

		sort.Slice(ks, func(i, j int) bool { return ks[i] > ks[j] })

		ret := make([]interface{}, 0)
		for _, id := range ks {
			ret = append(ret, vs[id])
		}
		return ret
	}

	addToList := func(list interface{}) {
		for i, c := range list.([]interface{}) {
			if i >= ListConnections_MaxLineCount {
				return
			}
			if c2, ok := c.(*ManagedV2rayConn); ok {
				if ListConnections_IgnoreTags(c2.Tag) {
					continue
				}
				if corePtr != 0 && corePtr != c2.corePtr {
					continue
				}
				item := &struct {
					ID    uint32
					Dest  string
					RDest string
					Uid   uint32
					Start int64
					End   int64
					Tag   string
				}{
					c2.ID(),
					c2.Dest,
					c2.RouteDest,
					c2.InboundUid,
					c2.StartTime,
					c2.EndTime,
					c2.Tag,
				}
				list2 = append(list2, item)
			}
		}
	}

	addToList(rangeMap(&ConnectionPool_V2Ray.Map))
	addToList(rangeMap(&ConnectionLog_V2Ray.Map))

	b, _ := json.Marshal(&list2)
	return string(b)
}

func ResetAllConnections(system bool) {
	if system {
		ConnectionPool_System.ResetConnections(0)
	} else {
		ConnectionPool_V2Ray.ResetConnections(0)
		ConnectionLog_V2Ray.ResetConnections(0)
	}
}

func ResetConnections(corePtr uintptr) {
	ConnectionLog_V2Ray.ResetConnections(corePtr)
	ConnectionPool_V2Ray.ResetConnections(corePtr)
	ConnectionPool_System.ResetConnections(corePtr)
}
