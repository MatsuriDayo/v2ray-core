package shadowsocks

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"hash/crc32"
	"io"
	gonet "net"

	"github.com/v2fly/v2ray-core/v5/common"
	"github.com/v2fly/v2ray-core/v5/common/buf"
	"github.com/v2fly/v2ray-core/v5/common/drain"
	"github.com/v2fly/v2ray-core/v5/common/net"
	"github.com/v2fly/v2ray-core/v5/common/protocol"
)

const (
	Version = 1
)

var addrParser = protocol.NewAddressParser(
	protocol.AddressFamilyByte(0x01, net.AddressFamilyIPv4),
	protocol.AddressFamilyByte(0x04, net.AddressFamilyIPv6),
	protocol.AddressFamilyByte(0x03, net.AddressFamilyDomain),
	protocol.WithAddressTypeParser(func(b byte) byte {
		return b & 0x0F
	}),
)

// ReadTCPSession reads a Shadowsocks TCP session from the given reader, returns its header and remaining parts.
func ReadTCPSession(user *protocol.MemoryUser, reader io.Reader, conn *ProtocolConn) (*protocol.RequestHeader, buf.Reader, error) {
	account := user.Account.(*MemoryAccount)

	hashkdf := hmac.New(sha256.New, []byte("SSBSKDF"))
	hashkdf.Write(account.Key)

	behaviorSeed := crc32.ChecksumIEEE(hashkdf.Sum(nil))

	drainer, err := drain.NewBehaviorSeedLimitedDrainer(int64(behaviorSeed), 16+38, 3266, 64)
	if err != nil {
		return nil, nil, newError("failed to initialize drainer").Base(err)
	}

	buffer := buf.New()
	defer buffer.Release()

	ivLen := account.Cipher.IVSize()
	var iv []byte
	if ivLen > 0 {
		if _, err := buffer.ReadFullFrom(reader, ivLen); err != nil {
			drainer.AcknowledgeReceive(int(buffer.Len()))
			return nil, nil, drain.WithError(drainer, reader, newError("failed to read IV").Base(err))
		}

		iv = append([]byte(nil), buffer.BytesTo(ivLen)...)
	}

	r, err := account.Cipher.NewDecryptionReader(account.Key, iv, reader)
	if err != nil {
		drainer.AcknowledgeReceive(int(buffer.Len()))
		return nil, nil, drain.WithError(drainer, reader, newError("failed to initialize decoding stream").Base(err).AtError())
	}

	if conn != nil {
		conn.Reader = r
		r = conn.ProtocolReader
	}

	br := &buf.BufferedReader{Reader: r}

	request := &protocol.RequestHeader{
		Version: Version,
		User:    user,
		Command: protocol.RequestCommandTCP,
	}

	drainer.AcknowledgeReceive(int(buffer.Len()))
	buffer.Clear()

	addr, port, err := addrParser.ReadAddressPort(buffer, br)
	if err != nil {
		drainer.AcknowledgeReceive(int(buffer.Len()))
		return nil, nil, drain.WithError(drainer, reader, newError("failed to read address").Base(err))
	}

	request.Address = addr
	request.Port = port

	if request.Address == nil {
		drainer.AcknowledgeReceive(int(buffer.Len()))
		return nil, nil, drain.WithError(drainer, reader, newError("invalid remote address."))
	}

	if ivError := account.CheckIV(iv); ivError != nil {
		drainer.AcknowledgeReceive(int(buffer.Len()))
		return nil, nil, drain.WithError(drainer, reader, newError("failed iv check").Base(ivError))
	}

	return request, br, nil
}

// WriteTCPRequest writes Shadowsocks request into the given writer, and returns a writer for body.
func WriteTCPRequest(request *protocol.RequestHeader, writer io.Writer, iv []byte, conn *ProtocolConn) (buf.Writer, error) {
	user := request.User
	account := user.Account.(*MemoryAccount)

	if len(iv) > 0 {
		if err := buf.WriteAllBytes(writer, iv); err != nil {
			return nil, newError("failed to write IV")
		}
	}

	w, err := account.Cipher.NewEncryptionWriter(account.Key, iv, writer)
	if err != nil {
		return nil, newError("failed to create encoding stream").Base(err).AtError()
	}

	if conn != nil {
		conn.Writer = w
		w = conn.ProtocolWriter
	}

	header := buf.New()

	if err := addrParser.WriteAddressPort(header, request.Address, request.Port); err != nil {
		return nil, newError("failed to write address").Base(err)
	}

	if err := w.WriteMultiBuffer(buf.MultiBuffer{header}); err != nil {
		return nil, newError("failed to write header").Base(err)
	}

	return w, nil
}

func ReadTCPResponse(user *protocol.MemoryUser, reader io.Reader, conn *ProtocolConn) (buf.Reader, error) {
	account := user.Account.(*MemoryAccount)

	hashkdf := hmac.New(sha256.New, []byte("SSBSKDF"))
	hashkdf.Write(account.Key)

	behaviorSeed := crc32.ChecksumIEEE(hashkdf.Sum(nil))

	drainer, err := drain.NewBehaviorSeedLimitedDrainer(int64(behaviorSeed), 16+38, 3266, 64)
	if err != nil {
		return nil, newError("failed to initialize drainer").Base(err)
	}

	var iv []byte
	if account.Cipher.IVSize() > 0 {
		iv = make([]byte, account.Cipher.IVSize())
		if n, err := io.ReadFull(reader, iv); err != nil {
			return nil, newError("failed to read IV").Base(err)
		} else { // nolint: revive
			drainer.AcknowledgeReceive(n)
		}
	}

	if ivError := account.CheckIV(iv); ivError != nil {
		return nil, drain.WithError(drainer, reader, newError("failed iv check").Base(ivError))
	}

	r, err := account.Cipher.NewDecryptionReader(account.Key, iv, reader)

	if conn != nil {
		conn.Reader = r
		r = conn.ProtocolReader
	}

	return r, err
}

func WriteTCPResponse(request *protocol.RequestHeader, writer io.Writer, iv []byte, conn *ProtocolConn) (buf.Writer, error) {
	user := request.User
	account := user.Account.(*MemoryAccount)

	if len(iv) > 0 {
		if err := buf.WriteAllBytes(writer, iv); err != nil {
			return nil, newError("failed to write IV.").Base(err)
		}
	}

	w, err := account.Cipher.NewEncryptionWriter(account.Key, iv, writer)

	if err == nil && conn != nil {
		conn.Writer = w
		w = conn.ProtocolWriter
	}

	return w, err
}

func EncodeUDPPacket(request *protocol.RequestHeader, payload []byte) (*buf.Buffer, error) {
	user := request.User
	account := user.Account.(*MemoryAccount)

	buffer := buf.New()
	ivLen := account.Cipher.IVSize()
	if ivLen > 0 {
		common.Must2(buffer.ReadFullFrom(rand.Reader, ivLen))
	}

	if err := addrParser.WriteAddressPort(buffer, request.Address, request.Port); err != nil {
		return nil, newError("failed to write address").Base(err)
	}

	buffer.Write(payload)

	if err := account.Cipher.EncodePacket(account.Key, buffer); err != nil {
		return nil, newError("failed to encrypt UDP payload").Base(err)
	}

	return buffer, nil
}

func DecodeUDPPacket(user *protocol.MemoryUser, payload *buf.Buffer) (*protocol.RequestHeader, *buf.Buffer, error) {
	account := user.Account.(*MemoryAccount)

	var iv []byte
	if !account.Cipher.IsAEAD() && account.Cipher.IVSize() > 0 {
		// Keep track of IV as it gets removed from payload in DecodePacket.
		iv = make([]byte, account.Cipher.IVSize())
		copy(iv, payload.BytesTo(account.Cipher.IVSize()))
	}

	if err := account.Cipher.DecodePacket(account.Key, payload); err != nil {
		return nil, nil, newError("failed to decrypt UDP payload").Base(err)
	}

	request := &protocol.RequestHeader{
		Version: Version,
		User:    user,
		Command: protocol.RequestCommandUDP,
	}

	payload.SetByte(0, payload.Byte(0)&0x0F)

	addr, port, err := addrParser.ReadAddressPort(nil, payload)
	if err != nil {
		return nil, nil, newError("failed to parse address").Base(err)
	}

	request.Address = addr
	request.Port = port

	return request, payload, nil
}

type UDPReader struct {
	Reader io.Reader
	User   *protocol.MemoryUser
	Plugin ProtocolPlugin
}

func (v *UDPReader) ReadMultiBuffer() (buf.MultiBuffer, error) {
	buffer := buf.New()
	_, err := buffer.ReadFrom(v.Reader)
	if err != nil {
		buffer.Release()
		return nil, err
	}
	if v.Plugin != nil {
		newBuffer, err := v.Plugin.DecodePacket(buffer)
		if err != nil {
			return nil, err
		}
		buffer = newBuffer
	}

	header, payload, err := DecodeUDPPacket(v.User, buffer)
	if err != nil {
		buffer.Release()
		return nil, err
	}
	endpoint := header.Destination()
	payload.Endpoint = &endpoint
	return buf.MultiBuffer{payload}, nil
}

func (v *UDPReader) ReadFrom(p []byte) (n int, addr gonet.Addr, err error) {
	buffer := buf.New()
	_, err = buffer.ReadFrom(v.Reader)
	if err != nil {
		buffer.Release()
		return 0, nil, err
	}
	vaddr, payload, err := DecodeUDPPacket(v.User, buffer)
	if err != nil {
		buffer.Release()
		return 0, nil, err
	}
	n = copy(p, payload.Bytes())
	payload.Release()
	return n, &gonet.UDPAddr{IP: vaddr.Address.IP(), Port: int(vaddr.Port)}, nil
}

type UDPWriter struct {
	Writer  io.Writer
	Request *protocol.RequestHeader
	Plugin  ProtocolPlugin
}

func (w *UDPWriter) WriteMultiBuffer(mb buf.MultiBuffer) error {
	for _, buffer := range mb {
		if buffer == nil {
			continue
		}
		request := w.Request
		if buffer.Endpoint != nil {
			request = &protocol.RequestHeader{
				User:    w.Request.User,
				Address: buffer.Endpoint.Address,
				Port:    buffer.Endpoint.Port,
			}
		}
		packet, err := EncodeUDPPacket(request, buffer.Bytes())
		buffer.Release()
		if err != nil {
			buf.ReleaseMulti(mb)
			return err
		}
		if w.Plugin != nil {
			newBuffer, err := w.Plugin.EncodePacket(buffer)
			if err != nil {
				return err
			}
			buffer = newBuffer
		}
		_, err = w.Writer.Write(packet.Bytes())
		packet.Release()
		if err != nil {
			buf.ReleaseMulti(mb)
			return err
		}
	}
	return nil
}

// Write implements io.Writer.
func (w *UDPWriter) Write(payload []byte) (int, error) {
	packet, err := EncodeUDPPacket(w.Request, payload)
	if err != nil {
		return 0, err
	}
	if w.Plugin != nil {
		newBuffer, err := w.Plugin.EncodePacket(packet)
		if err != nil {
			return 0, err
		}
		packet = newBuffer
	}
	_, err = w.Writer.Write(packet.Bytes())
	packet.Release()
	return len(payload), err
}

func (w *UDPWriter) WriteTo(payload []byte, addr gonet.Addr) (n int, err error) {
	request := *w.Request
	udpAddr := addr.(*gonet.UDPAddr)
	request.Command = protocol.RequestCommandUDP
	request.Address = net.IPAddress(udpAddr.IP)
	request.Port = net.Port(udpAddr.Port)
	packet, err := EncodeUDPPacket(&request, payload)
	if err != nil {
		return 0, err
	}
	_, err = w.Writer.Write(packet.Bytes())
	packet.Release()
	return len(payload), err
}
