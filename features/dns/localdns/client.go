package localdns

import (
	"context"

	"github.com/v2fly/v2ray-core/v5/common/net"
	"github.com/v2fly/v2ray-core/v5/features/dns"
)

var (
	Instance   = &Client{}
	LookupFunc func(network string, host string) ([]net.IP, error)
)

func init() {
	SetLookupFunc(nil)
}

func SetLookupFunc(lookupFunc func(network, host string) ([]net.IP, error)) {
	if lookupFunc == nil {
		resolver := &net.Resolver{PreferGo: false}
		LookupFunc = func(network string, host string) ([]net.IP, error) {
			return resolver.LookupIP(context.Background(), network, host)
		}
	} else {
		LookupFunc = lookupFunc
	}
}

// Client is an implementation of dns.Client, which queries localhost for DNS.
type Client struct {
	resolver *net.Resolver
}

// Type implements common.HasType.
func (*Client) Type() interface{} {
	return dns.ClientType()
}

// Start implements common.Runnable.
func (*Client) Start() error {
	return nil
}

// Close implements common.Closable.
func (*Client) Close() error { return nil }

// LookupIP implements Client.
func (c *Client) LookupIP(host dns.MatsuriDomainString) ([]net.IP, error) {
	return LookupFunc("ip", matsuriHookGetDomain(host))
}

// LookupIPv4 implements IPv4Lookup.
func (c *Client) LookupIPv4(host dns.MatsuriDomainString) ([]net.IP, error) {
	return LookupFunc("ip4", matsuriHookGetDomain(host))
}

// LookupIPv6 implements IPv6Lookup.
func (c *Client) LookupIPv6(host dns.MatsuriDomainString) ([]net.IP, error) {
	return LookupFunc("ip6", matsuriHookGetDomain(host))
}

// New create a new dns.Client that queries localhost for DNS.
func New() *Client {
	return Instance
}

// Masuri: hook
func matsuriHookGetDomain(_domain dns.MatsuriDomainString) string {
	var domain string
	if a, ok := _domain.(*dns.MatsuriDomainStringEx); ok {
		domain = a.Domain
	} else if a, ok := _domain.(string); ok {
		domain = a
	}
	return domain
}
