package shadowsocks

import (
	"context"
	"github.com/v2fly/v2ray-core/v5/common"
	"github.com/v2fly/v2ray-core/v5/common/buf"
	"github.com/v2fly/v2ray-core/v5/transport/internet"
)

var (
	pluginLoader func(plugin string) SIP003Plugin
	plugins      map[string]func() SIP003Plugin
)

func init() {
	plugins = make(map[string]func() SIP003Plugin)
}

func SetPluginLoader(creator func(plugin string) SIP003Plugin) {
	pluginLoader = creator
}

func RegisterPlugin(name string, creator func() SIP003Plugin) {
	plugins[name] = creator
}

type SIP003Plugin interface {
	Init(ctx context.Context, localHost string, localPort string, remoteHost string, remotePort string, pluginOpts string, pluginArgs []string, account *MemoryAccount) error
	common.Closable
}

type StreamPlugin interface {
	StreamConn(conn internet.Connection) internet.Connection
}

type ProtocolConn struct {
	buf.Reader
	buf.Writer
	ProtocolReader buf.Reader
	ProtocolWriter buf.Writer
}

// 喵
type ProtocolPlugin interface {
	ProtocolConn(conn *ProtocolConn, iv []byte)
	EncodePacket(data []byte) ([]byte, error)
	DecodePacket(data []byte) ([]byte, error)
}
