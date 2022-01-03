package dns

import (
	"github.com/v2fly/v2ray-core/v5/common/errors"
	"github.com/v2fly/v2ray-core/v5/common/net"
	"github.com/v2fly/v2ray-core/v5/common/serial"
	"github.com/v2fly/v2ray-core/v5/common/session"
	"github.com/v2fly/v2ray-core/v5/features"
)

// IPOption is an object for IP query options.
type IPOption struct {
	IPv4Enable bool
	IPv6Enable bool
	FakeEnable bool
}

// Matsuri: it's a string before.
type MatsuriDomainString interface {
}

// Matsuri: private options for Lookup
type MatsuriDomainStringEx struct {
	Domain     string
	OptInbound *session.Inbound
}

// Client is a V2Ray feature for querying DNS information.
//
// v2ray:api:stable
type Client interface {
	features.Feature

	// LookupIP returns IP address for the given domain. IPs may contain IPv4 and/or IPv6 addresses.
	LookupIP(domain MatsuriDomainString) ([]net.IP, error)
}

// IPv4Lookup is an optional feature for querying IPv4 addresses only.
//
// v2ray:api:beta
type IPv4Lookup interface {
	LookupIPv4(domain MatsuriDomainString) ([]net.IP, error)
}

// IPv6Lookup is an optional feature for querying IPv6 addresses only.
//
// v2ray:api:beta
type IPv6Lookup interface {
	LookupIPv6(domain MatsuriDomainString) ([]net.IP, error)
}

// ClientWithIPOption is an optional feature for querying DNS information.
//
// v2ray:api:beta
type ClientWithIPOption interface {
	// GetIPOption returns IPOption for the DNS client.
	GetIPOption() *IPOption

	// SetQueryOption sets IPv4Enable and IPv6Enable for the DNS client.
	SetQueryOption(isIPv4Enable, isIPv6Enable bool)

	// SetFakeDNSOption sets FakeEnable option for DNS client.
	SetFakeDNSOption(isFakeEnable bool)
}

// ClientType returns the type of Client interface. Can be used for implementing common.HasType.
//
// v2ray:api:beta
func ClientType() interface{} {
	return (*Client)(nil)
}

// ErrEmptyResponse indicates that DNS query succeeded but no answer was returned.
var ErrEmptyResponse = errors.New("empty response")

type RCodeError uint16

func (e RCodeError) Error() string {
	return serial.Concat("rcode: ", uint16(e))
}

func RCodeFromError(err error) uint16 {
	if err == nil {
		return 0
	}
	cause := errors.Cause(err)
	if r, ok := cause.(RCodeError); ok {
		return uint16(r)
	}
	return 0
}
