package nekoutils

import (
	"sync"
	"sync/atomic"
	"time"
)

type ManagedV2rayConn struct {
	id      uint32
	lock    sync.Mutex
	corePtr uintptr

	CloseFunc func() error

	Dest       string
	RouteDest  string
	InboundTag string
	InboundUid uint32
	Tag        string

	StartTime int64
	EndTime   int64
}

func (c *ManagedV2rayConn) Close() error {
	if c.CloseFunc != nil {
		return c.CloseFunc()
	}
	return nil
}

func (c *ManagedV2rayConn) RemoteAddress() string {
	return c.Dest
}

func (c *ManagedV2rayConn) ID() uint32 {
	return c.id
}

func (c *ManagedV2rayConn) Instance() uintptr {
	return c.corePtr
}

// 在此添加连接

func (c *ManagedV2rayConn) ConnectionStart(core uintptr) {
	c.StartTime = time.Now().Unix()
	c.id = atomic.AddUint32(&ConnectionPool_V2Ray.cnt, 1)
	c.corePtr = core
	ConnectionPool_V2Ray.AddConnection(c)
}

func (c *ManagedV2rayConn) ConnectionEnd() {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.EndTime > 0 {
		return
	}
	c.EndTime = time.Now().Unix()

	// Move to log
	ConnectionPool_V2Ray.RemoveConnection(c)
	ConnectionLog_V2Ray.AddConnection(c)
}
