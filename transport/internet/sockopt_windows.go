package internet

import (
	"encoding/binary"
	"log"
	"strings"
	"syscall"
	"unsafe"

	"github.com/v2fly/v2ray-core/v5/nekoutils"
	"golang.org/x/sys/windows"
)

const (
	TCP_FASTOPEN = 15 // nolint: revive,stylecheck
)

func setTFO(fd syscall.Handle, settings SocketConfig_TCPFastOpenState) error {
	switch settings {
	case SocketConfig_Enable:
		if err := syscall.SetsockoptInt(fd, syscall.IPPROTO_TCP, TCP_FASTOPEN, 1); err != nil {
			return err
		}
	case SocketConfig_Disable:
		if err := syscall.SetsockoptInt(fd, syscall.IPPROTO_TCP, TCP_FASTOPEN, 0); err != nil {
			return err
		}
	}
	return nil
}

func applyOutboundSocketOptions(network string, address string, fd uintptr, config *SocketConfig) error {
	//
	if nekoutils.Windows_Protect_BindInterfaceIndex != nil {
		BindInterfaceIndex := nekoutils.Windows_Protect_BindInterfaceIndex()
		if BindInterfaceIndex != 0 {
			var v4, v6 bool
			if strings.HasSuffix(network, "6") {
				v4 = false
				v6 = true
			} else {
				v4 = true
				v6 = false
			}
			if err := bindInterface(fd, BindInterfaceIndex, v4, v6); err != nil {
				log.Println("bind outbound interface", err)
				return err
			}
		}
	}
	if config == nil {
		return nil
	}
	//

	if isTCPSocket(network) {
		if err := setTFO(syscall.Handle(fd), config.Tfo); err != nil {
			return err
		}
		if config.TcpKeepAliveIdle > 0 {
			if err := syscall.SetsockoptInt(syscall.Handle(fd), syscall.SOL_SOCKET, syscall.SO_KEEPALIVE, 1); err != nil {
				return newError("failed to set SO_KEEPALIVE", err)
			}
		}
	}

	if config.TxBufSize != 0 {
		if err := windows.SetsockoptInt(windows.Handle(fd), windows.SOL_SOCKET, windows.SO_SNDBUF, int(config.TxBufSize)); err != nil {
			return newError("failed to set SO_SNDBUF").Base(err)
		}
	}

	if config.RxBufSize != 0 {
		if err := windows.SetsockoptInt(windows.Handle(fd), windows.SOL_SOCKET, windows.SO_RCVBUF, int(config.TxBufSize)); err != nil {
			return newError("failed to set SO_RCVBUF").Base(err)
		}
	}

	return nil
}

func applyInboundSocketOptions(network string, fd uintptr, config *SocketConfig) error {
	//
	if nekoutils.Windows_Protect_BindInterfaceIndex != nil {
		BindInterfaceIndex := nekoutils.Windows_Protect_BindInterfaceIndex()
		if BindInterfaceIndex != 0 {
			if err := bindInterface(fd, BindInterfaceIndex, true, true); err != nil {
				log.Println("bind inbound interface", err)
				return err
			}
		}
	}
	if config == nil {
		return nil
	}
	//

	if isTCPSocket(network) {
		if err := setTFO(syscall.Handle(fd), config.Tfo); err != nil {
			return err
		}
		if config.TcpKeepAliveIdle > 0 {
			if err := syscall.SetsockoptInt(syscall.Handle(fd), syscall.SOL_SOCKET, syscall.SO_KEEPALIVE, 1); err != nil {
				return newError("failed to set SO_KEEPALIVE", err)
			}
		}
	}

	if config.TxBufSize != 0 {
		if err := windows.SetsockoptInt(windows.Handle(fd), windows.SOL_SOCKET, windows.SO_SNDBUF, int(config.TxBufSize)); err != nil {
			return newError("failed to set SO_SNDBUF").Base(err)
		}
	}

	if config.RxBufSize != 0 {
		if err := windows.SetsockoptInt(windows.Handle(fd), windows.SOL_SOCKET, windows.SO_RCVBUF, int(config.TxBufSize)); err != nil {
			return newError("failed to set SO_RCVBUF").Base(err)
		}
	}

	return nil
}

func bindAddr(fd uintptr, ip []byte, port uint32) error {
	return nil
}

func setReuseAddr(fd uintptr) error {
	return nil
}

func setReusePort(fd uintptr) error {
	return nil
}

const (
	IP_UNICAST_IF   = 31 // nolint: golint,stylecheck
	IPV6_UNICAST_IF = 31 // nolint: golint,stylecheck
)

func bindInterface(fd uintptr, interfaceIndex uint32, v4, v6 bool) error {
	if v4 {
		/* MSDN says for IPv4 this needs to be in net byte order, so that it's like an IP address with leading zeros. */
		bytes := make([]byte, 4)
		binary.BigEndian.PutUint32(bytes, interfaceIndex)
		interfaceIndex_v4 := *(*uint32)(unsafe.Pointer(&bytes[0]))

		if err := windows.SetsockoptInt(windows.Handle(fd), windows.IPPROTO_IP, IP_UNICAST_IF, int(interfaceIndex_v4)); err != nil {
			return err
		}
	}

	if v6 {
		if err := windows.SetsockoptInt(windows.Handle(fd), windows.IPPROTO_IPV6, IPV6_UNICAST_IF, int(interfaceIndex)); err != nil {
			return err
		}
	}

	return nil
}
