package transport

import (
	"github.com/v2fly/v2ray-core/v5/common"
	"github.com/v2fly/v2ray-core/v5/common/buf"
)

// Link is a utility for connecting between an inbound and an outbound proxy handler.
type Link struct {
	Reader buf.Reader
	Writer buf.Writer
}

func LinkWithCloseHook(link *Link, hook func() bool) *Link {
	newLink := &Link{
		Reader: &readerWithCloseHook{link.Reader, hook},
		Writer: &writerWithCloseHook{link.Writer, hook},
	}
	return newLink
}

type readerWithCloseHook struct {
	buf.Reader
	hook func() bool
}

func (r *readerWithCloseHook) Interrupt() {
	if r.hook() {
		common.Interrupt(r.Reader)
	}
}

func (r *readerWithCloseHook) Close() (err error) {
	if r.hook() {
		return common.Close(r.Reader)
	}
	return nil
}

type writerWithCloseHook struct {
	buf.Writer
	hook func() bool
}

func (r *writerWithCloseHook) Interrupt() {
	if r.hook() {
		common.Interrupt(r.Writer)
	}
}

func (r *writerWithCloseHook) Close() (err error) {
	if r.hook() {
		return common.Close(r.Writer)
	}
	return nil
}
