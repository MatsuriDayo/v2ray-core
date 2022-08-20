package simplified

import (
	"context"

	"github.com/v2fly/v2ray-core/v5/common"
	"github.com/v2fly/v2ray-core/v5/common/protocol"
	"github.com/v2fly/v2ray-core/v5/common/serial"
	"github.com/v2fly/v2ray-core/v5/proxy/shadowsocks"
)

func init() {
	common.Must(common.RegisterConfig((*ServerConfig)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		simplifiedServer := config.(*ServerConfig)
		fullServer := &shadowsocks.ServerConfig{
			User: &protocol.User{
				Account: serial.ToTypedMessage(&shadowsocks.Account{
					Password:   simplifiedServer.Password,
					CipherType: simplifiedServer.Method,
				}),
			},
			Network:        simplifiedServer.Networks.Network,
			PacketEncoding: simplifiedServer.PacketEncoding,
			Plugin:     simplifiedServer.Plugin,
			PluginOpts: simplifiedServer.PluginOpts,
			PluginArgs: simplifiedServer.PluginArgs,
		}

		return common.CreateObject(ctx, fullServer)
	}))

	common.Must(common.RegisterConfig((*ClientConfig)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		simplifiedClient := config.(*ClientConfig)
		fullClient := &shadowsocks.ClientConfig{
			Server: []*protocol.ServerEndpoint{
				{
					Address: simplifiedClient.Address,
					Port:    simplifiedClient.Port,
					User: []*protocol.User{
						{
							Account: serial.ToTypedMessage(&shadowsocks.Account{
								Password:                       simplifiedClient.Password,
								CipherType:                     simplifiedClient.Method,
								ExperimentReducedIvHeadEntropy: simplifiedClient.ExperimentReducedIvHeadEntropy,
							}),
						},
					},
				},
			},
			Plugin:     simplifiedClient.Plugin,
			PluginOpts: simplifiedClient.PluginOpts,
			PluginArgs: simplifiedClient.PluginArgs,
		}

		return common.CreateObject(ctx, fullClient)
	}))
}
