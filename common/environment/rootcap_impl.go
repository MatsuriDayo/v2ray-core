package environment

import (
	"context"

	"github.com/v2fly/v2ray-core/v5/common/platform/filesystem/fsifce"
	"github.com/v2fly/v2ray-core/v5/features/extension/storage"
	"github.com/v2fly/v2ray-core/v5/transport/internet"
	"github.com/v2fly/v2ray-core/v5/transport/internet/tagged"
)

func NewRootEnvImpl(ctx context.Context, transientStorage storage.ScopedTransientStorage) RootEnvironment {
	return &rootEnvImpl{transientStorage: transientStorage, ctx: ctx}
}

type rootEnvImpl struct {
	transientStorage storage.ScopedTransientStorage

	ctx context.Context
}

func (r *rootEnvImpl) doNotImpl() {
	panic("placeholder doNotImpl")
}

func (r *rootEnvImpl) AppEnvironment(tag string) AppEnvironment {
	transientStorage, err := r.transientStorage.NarrowScope(r.ctx, tag)
	if err != nil {
		return nil
	}
	return &appEnvImpl{
		transientStorage: transientStorage,
		ctx:              r.ctx,
	}
}

func (r *rootEnvImpl) ProxyEnvironment(tag string) ProxyEnvironment {
	transientStorage, err := r.transientStorage.NarrowScope(r.ctx, tag)
	if err != nil {
		return nil
	}
	return &proxyEnvImpl{
		transientStorage: transientStorage,
		ctx:              r.ctx,
	}
}

type appEnvImpl struct {
	transientStorage storage.ScopedTransientStorage

	ctx context.Context
}

func (a *appEnvImpl) RequireFeatures() interface{} {
	panic("implement me")
}

func (a *appEnvImpl) RecordLog() interface{} {
	panic("implement me")
}

func (a *appEnvImpl) Dialer() internet.SystemDialer {
	panic("implement me")
}

func (a *appEnvImpl) Listener() internet.SystemListener {
	panic("implement me")
}

func (a *appEnvImpl) OutboundDialer() tagged.DialFunc {
	panic("implement me")
}

func (a *appEnvImpl) OpenFileForReadSeek() fsifce.FileSeekerFunc {
	panic("implement me")
}

func (a *appEnvImpl) OpenFileForRead() fsifce.FileReaderFunc {
	panic("implement me")
}

func (a *appEnvImpl) OpenFileForWrite() fsifce.FileWriterFunc {
	panic("implement me")
}

func (a *appEnvImpl) PersistentStorage() storage.ScopedPersistentStorage {
	panic("implement me")
}

func (a *appEnvImpl) TransientStorage() storage.ScopedTransientStorage {
	return a.transientStorage
}

func (a *appEnvImpl) NarrowScope(key string) (AppEnvironment, error) {
	transientStorage, err := a.transientStorage.NarrowScope(a.ctx, key)
	if err != nil {
		return nil, err
	}
	return &appEnvImpl{
		transientStorage: transientStorage,
		ctx:              a.ctx,
	}, nil
}

func (a *appEnvImpl) doNotImpl() {
	panic("placeholder doNotImpl")
}

type proxyEnvImpl struct {
	transientStorage storage.ScopedTransientStorage

	ctx context.Context
}

func (p *proxyEnvImpl) RequireFeatures() interface{} {
	panic("implement me")
}

func (p *proxyEnvImpl) RecordLog() interface{} {
	panic("implement me")
}

func (p *proxyEnvImpl) OutboundDialer() tagged.DialFunc {
	panic("implement me")
}

func (p *proxyEnvImpl) TransientStorage() storage.ScopedTransientStorage {
	return p.transientStorage
}

func (p *proxyEnvImpl) NarrowScope(key string) (ProxyEnvironment, error) {
	transientStorage, err := p.transientStorage.NarrowScope(p.ctx, key)
	if err != nil {
		return nil, err
	}
	return &proxyEnvImpl{
		transientStorage: transientStorage,
		ctx:              p.ctx,
	}, nil
}

func (p *proxyEnvImpl) NarrowScopeToTransport(key string) (TransportEnvironment, error) {
	transientStorage, err := p.transientStorage.NarrowScope(p.ctx, key)
	if err != nil {
		return nil, err
	}
	return &transportEnvImpl{
		ctx:              p.ctx,
		transientStorage: transientStorage,
	}, nil
}

func (p *proxyEnvImpl) doNotImpl() {
	panic("placeholder doNotImpl")
}

type transportEnvImpl struct {
	transientStorage storage.ScopedTransientStorage

	ctx context.Context
}

func (t *transportEnvImpl) RequireFeatures() interface{} {
	panic("implement me")
}

func (t *transportEnvImpl) RecordLog() interface{} {
	panic("implement me")
}

func (t *transportEnvImpl) Dialer() internet.SystemDialer {
	panic("implement me")
}

func (t *transportEnvImpl) Listener() internet.SystemListener {
	panic("implement me")
}

func (t *transportEnvImpl) OutboundDialer() tagged.DialFunc {
	panic("implement me")
}

func (t *transportEnvImpl) TransientStorage() storage.ScopedTransientStorage {
	return t.transientStorage
}

func (t *transportEnvImpl) NarrowScope(key string) (TransportEnvironment, error) {
	transientStorage, err := t.transientStorage.NarrowScope(t.ctx, key)
	if err != nil {
		return nil, err
	}
	return &transportEnvImpl{
		ctx:              t.ctx,
		transientStorage: transientStorage,
	}, nil
}

func (t *transportEnvImpl) doNotImpl() {
	panic("implement me")
}
