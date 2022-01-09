package shadowsocks

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/md5"
	"crypto/rc4"
	"crypto/sha1"
	"io"
	"strings"

	"github.com/aead/chacha20"
	"github.com/aead/chacha20/chacha"
	"github.com/dgryski/go-camellia"
	"github.com/dgryski/go-idea"
	"github.com/dgryski/go-rc2"
	"github.com/geeksbaek/seed"
	"github.com/kierdavis/cfb8"
	"golang.org/x/crypto/blowfish"
	"golang.org/x/crypto/cast5"
	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/hkdf"

	"github.com/v2fly/v2ray-core/v5/common"
	"github.com/v2fly/v2ray-core/v5/common/antireplay"
	"github.com/v2fly/v2ray-core/v5/common/buf"
	"github.com/v2fly/v2ray-core/v5/common/crypto"
	"github.com/v2fly/v2ray-core/v5/common/protocol"
)

// MemoryAccount is an account type converted from Account.
type MemoryAccount struct {
	Cipher Cipher
	Key    []byte

	replayFilter antireplay.GeneralizedReplayFilter
}

// Equals implements protocol.Account.Equals().
func (a *MemoryAccount) Equals(another protocol.Account) bool {
	if account, ok := another.(*MemoryAccount); ok {
		return bytes.Equal(a.Key, account.Key)
	}
	return false
}

func (a *MemoryAccount) CheckIV(iv []byte) error {
	if a.replayFilter == nil {
		return nil
	}
	if a.replayFilter.Check(iv) {
		return nil
	}
	return newError("IV is not unique")
}

func createAesGcm(key []byte) cipher.AEAD {
	block, err := aes.NewCipher(key)
	common.Must(err)
	gcm, err := cipher.NewGCM(block)
	common.Must(err)
	return gcm
}

func createChaCha20Poly1305(key []byte) cipher.AEAD {
	ChaChaPoly1305, err := chacha20poly1305.New(key)
	common.Must(err)
	return ChaChaPoly1305
}

func createXChaCha20Poly1305(key []byte) cipher.AEAD {
	XChaChaPoly1305, err := chacha20poly1305.NewX(key)
	common.Must(err)
	return XChaChaPoly1305
}

func (a *Account) getCipher() (Cipher, error) {
	switch a.CipherType {
	case CipherType_AES_128_GCM:
		return &AEADCipher{
			KeyBytes:        16,
			IVBytes:         16,
			AEADAuthCreator: createAesGcm,
		}, nil
	case CipherType_AES_192_GCM:
		return &AEADCipher{
			KeyBytes:        24,
			IVBytes:         24,
			AEADAuthCreator: createAesGcm,
		}, nil
	case CipherType_AES_256_GCM:
		return &AEADCipher{
			KeyBytes:        32,
			IVBytes:         32,
			AEADAuthCreator: createAesGcm,
		}, nil
	case CipherType_CHACHA20_IETF_POLY1305:
		return &AEADCipher{
			KeyBytes:        32,
			IVBytes:         32,
			AEADAuthCreator: createChaCha20Poly1305,
		}, nil
	case CipherType_XCHACHA20_IETF_POLY1305:
		return &AEADCipher{
			KeyBytes:        32,
			IVBytes:         32,
			AEADAuthCreator: createXChaCha20Poly1305,
		}, nil
	case CipherType_NONE:
		return &NoneCipher{}, nil

	case CipherType_AES_128_CTR:
		return &StreamCipher{
			KeyBytes:       16,
			IVBytes:        aes.BlockSize,
			EncryptCreator: blockStream(aes.NewCipher, cipher.NewCTR),
			DecryptCreator: blockStream(aes.NewCipher, cipher.NewCTR),
		}, nil
	case CipherType_AES_192_CTR:
		return &StreamCipher{
			KeyBytes:       24,
			IVBytes:        aes.BlockSize,
			EncryptCreator: blockStream(aes.NewCipher, cipher.NewCTR),
			DecryptCreator: blockStream(aes.NewCipher, cipher.NewCTR),
		}, nil
	case CipherType_AES_256_CTR:
		return &StreamCipher{
			KeyBytes:       32,
			IVBytes:        aes.BlockSize,
			EncryptCreator: blockStream(aes.NewCipher, cipher.NewCTR),
			DecryptCreator: blockStream(aes.NewCipher, cipher.NewCTR),
		}, nil
	case CipherType_AES_128_CFB:
		return &StreamCipher{
			KeyBytes:       16,
			IVBytes:        aes.BlockSize,
			EncryptCreator: blockStream(aes.NewCipher, cipher.NewCFBEncrypter),
			DecryptCreator: blockStream(aes.NewCipher, cipher.NewCFBDecrypter),
		}, nil
	case CipherType_AES_192_CFB:
		return &StreamCipher{
			KeyBytes:       24,
			IVBytes:        aes.BlockSize,
			EncryptCreator: blockStream(aes.NewCipher, cipher.NewCFBEncrypter),
			DecryptCreator: blockStream(aes.NewCipher, cipher.NewCFBDecrypter),
		}, nil
	case CipherType_AES_256_CFB:
		return &StreamCipher{
			KeyBytes:       32,
			IVBytes:        aes.BlockSize,
			EncryptCreator: blockStream(aes.NewCipher, cipher.NewCFBEncrypter),
			DecryptCreator: blockStream(aes.NewCipher, cipher.NewCFBDecrypter),
		}, nil
	case CipherType_AES_128_CFB8:
		return &StreamCipher{
			KeyBytes:       16,
			IVBytes:        aes.BlockSize,
			EncryptCreator: blockStream(aes.NewCipher, cfb8.NewEncrypter),
			DecryptCreator: blockStream(aes.NewCipher, cfb8.NewDecrypter),
		}, nil
	case CipherType_AES_192_CFB8:
		return &StreamCipher{
			KeyBytes:       24,
			IVBytes:        aes.BlockSize,
			EncryptCreator: blockStream(aes.NewCipher, cfb8.NewEncrypter),
			DecryptCreator: blockStream(aes.NewCipher, cfb8.NewDecrypter),
		}, nil
	case CipherType_AES_256_CFB8:
		return &StreamCipher{
			KeyBytes:       32,
			IVBytes:        aes.BlockSize,
			EncryptCreator: blockStream(aes.NewCipher, cfb8.NewEncrypter),
			DecryptCreator: blockStream(aes.NewCipher, cfb8.NewDecrypter),
		}, nil
	case CipherType_AES_128_OFB:
		return &StreamCipher{
			KeyBytes:       16,
			IVBytes:        aes.BlockSize,
			EncryptCreator: blockStream(aes.NewCipher, cfb8.NewEncrypter),
			DecryptCreator: blockStream(aes.NewCipher, cfb8.NewDecrypter),
		}, nil
	case CipherType_AES_192_OFB:
		return &StreamCipher{
			KeyBytes:       24,
			IVBytes:        aes.BlockSize,
			EncryptCreator: blockStream(aes.NewCipher, cfb8.NewEncrypter),
			DecryptCreator: blockStream(aes.NewCipher, cfb8.NewDecrypter),
		}, nil
	case CipherType_AES_256_OFB:
		return &StreamCipher{
			KeyBytes:       32,
			IVBytes:        aes.BlockSize,
			EncryptCreator: blockStream(aes.NewCipher, cfb8.NewEncrypter),
			DecryptCreator: blockStream(aes.NewCipher, cfb8.NewDecrypter),
		}, nil
	case CipherType_RC4:
		return &StreamCipher{
			KeyBytes: 16,
			IVBytes:  16,
			EncryptCreator: func(key []byte, iv []byte) (cipher.Stream, error) {
				return rc4.NewCipher(key)
			},
			DecryptCreator: func(key []byte, iv []byte) (cipher.Stream, error) {
				return rc4.NewCipher(key)
			},
		}, nil
	case CipherType_RC4_MD5:
		return &StreamCipher{
			KeyBytes: 16,
			IVBytes:  16,
			EncryptCreator: func(key []byte, iv []byte) (cipher.Stream, error) {
				h := md5.New()
				h.Write(key)
				h.Write(iv)
				return rc4.NewCipher(h.Sum(nil))
			},
			DecryptCreator: func(key []byte, iv []byte) (cipher.Stream, error) {
				h := md5.New()
				h.Write(key)
				h.Write(iv)
				return rc4.NewCipher(h.Sum(nil))
			},
		}, nil
	case CipherType_BF_CFB:
		return &StreamCipher{
			KeyBytes:       16,
			IVBytes:        blowfish.BlockSize,
			EncryptCreator: blockStream(func(key []byte) (cipher.Block, error) { return blowfish.NewCipher(key) }, cipher.NewCFBEncrypter),
			DecryptCreator: blockStream(func(key []byte) (cipher.Block, error) { return blowfish.NewCipher(key) }, cipher.NewCFBDecrypter),
		}, nil
	case CipherType_CAST5_CFB:
		return &StreamCipher{
			KeyBytes:       16,
			IVBytes:        cast5.BlockSize,
			EncryptCreator: blockStream(func(key []byte) (cipher.Block, error) { return cast5.NewCipher(key) }, cipher.NewCFBEncrypter),
			DecryptCreator: blockStream(func(key []byte) (cipher.Block, error) { return cast5.NewCipher(key) }, cipher.NewCFBDecrypter),
		}, nil
	case CipherType_DES_CFB:
		return &StreamCipher{
			KeyBytes:       8,
			IVBytes:        des.BlockSize,
			EncryptCreator: blockStream(des.NewCipher, cipher.NewCFBEncrypter),
			DecryptCreator: blockStream(des.NewCipher, cipher.NewCFBDecrypter),
		}, nil
	case CipherType_IDEA_CFB:
		return &StreamCipher{
			KeyBytes:       16,
			IVBytes:        8,
			EncryptCreator: blockStream(idea.NewCipher, cipher.NewCFBEncrypter),
			DecryptCreator: blockStream(idea.NewCipher, cipher.NewCFBDecrypter),
		}, nil
	case CipherType_RC2_CFB:
		return &StreamCipher{
			KeyBytes:       16,
			IVBytes:        rc2.BlockSize,
			EncryptCreator: blockStream(func(key []byte) (cipher.Block, error) { return rc2.New(key, 16) }, cipher.NewCFBEncrypter),
			DecryptCreator: blockStream(func(key []byte) (cipher.Block, error) { return rc2.New(key, 16) }, cipher.NewCFBDecrypter),
		}, nil
	case CipherType_SEED_CFB:
		return &StreamCipher{
			KeyBytes:       16,
			IVBytes:        seed.BlockSize,
			EncryptCreator: blockStream(seed.NewCipher, cipher.NewCFBEncrypter),
			DecryptCreator: blockStream(seed.NewCipher, cipher.NewCFBDecrypter),
		}, nil
	case CipherType_CAMELLIA_128_CFB:
		return &StreamCipher{
			KeyBytes:       16,
			IVBytes:        camellia.BlockSize,
			EncryptCreator: blockStream(camellia.New, cipher.NewCFBEncrypter),
			DecryptCreator: blockStream(camellia.New, cipher.NewCFBDecrypter),
		}, nil
	case CipherType_CAMELLIA_192_CFB:
		return &StreamCipher{
			KeyBytes:       24,
			IVBytes:        camellia.BlockSize,
			EncryptCreator: blockStream(camellia.New, cipher.NewCFBEncrypter),
			DecryptCreator: blockStream(camellia.New, cipher.NewCFBDecrypter),
		}, nil
	case CipherType_CAMELLIA_256_CFB:
		return &StreamCipher{
			KeyBytes:       32,
			IVBytes:        camellia.BlockSize,
			EncryptCreator: blockStream(camellia.New, cipher.NewCFBEncrypter),
			DecryptCreator: blockStream(camellia.New, cipher.NewCFBDecrypter),
		}, nil
	case CipherType_CAMELLIA_128_CFB8:
		return &StreamCipher{
			KeyBytes:       16,
			IVBytes:        camellia.BlockSize,
			EncryptCreator: blockStream(camellia.New, cfb8.NewEncrypter),
			DecryptCreator: blockStream(camellia.New, cfb8.NewDecrypter),
		}, nil
	case CipherType_CAMELLIA_192_CFB8:
		return &StreamCipher{
			KeyBytes:       24,
			IVBytes:        camellia.BlockSize,
			EncryptCreator: blockStream(camellia.New, cfb8.NewEncrypter),
			DecryptCreator: blockStream(camellia.New, cfb8.NewDecrypter),
		}, nil
	case CipherType_CAMELLIA_256_CFB8:
		return &StreamCipher{
			KeyBytes:       32,
			IVBytes:        camellia.BlockSize,
			EncryptCreator: blockStream(camellia.New, cfb8.NewEncrypter),
			DecryptCreator: blockStream(camellia.New, cfb8.NewDecrypter),
		}, nil
	case CipherType_SALSA20:
		return &StreamCipher{
			KeyBytes:       32,
			IVBytes:        8,
			EncryptCreator: crypto.NewSalsa20,
			DecryptCreator: crypto.NewSalsa20,
		}, nil
	case CipherType_CHACHA20:
		return &StreamCipher{
			KeyBytes: chacha.KeySize,
			IVBytes:  chacha.NonceSize,
			EncryptCreator: func(key []byte, iv []byte) (cipher.Stream, error) {
				return chacha20.NewCipher(iv, key)
			},
			DecryptCreator: func(key []byte, iv []byte) (cipher.Stream, error) {
				return chacha20.NewCipher(iv, key)
			},
		}, nil
	case CipherType_CHACHA20_IETF:
		return &StreamCipher{
			KeyBytes: chacha.KeySize,
			IVBytes:  chacha.INonceSize,
			EncryptCreator: func(key []byte, iv []byte) (cipher.Stream, error) {
				return chacha20.NewCipher(iv, key)
			},
			DecryptCreator: func(key []byte, iv []byte) (cipher.Stream, error) {
				return chacha20.NewCipher(iv, key)
			},
		}, nil
	case CipherType_XCHACHA20:
		return &StreamCipher{
			KeyBytes: chacha.KeySize,
			IVBytes:  chacha.XNonceSize,
			EncryptCreator: func(key []byte, iv []byte) (cipher.Stream, error) {
				return chacha20.NewCipher(iv, key)
			},
			DecryptCreator: func(key []byte, iv []byte) (cipher.Stream, error) {
				return chacha20.NewCipher(iv, key)
			},
		}, nil
	default:
		return nil, newError("Unsupported cipher.")
	}
}

// AsAccount implements protocol.AsAccount.
func (a *Account) AsAccount() (protocol.Account, error) {
	Cipher, err := a.getCipher()
	if err != nil {
		return nil, newError("failed to get cipher").Base(err)
	}
	return &MemoryAccount{
		Cipher: Cipher,
		Key:    passwordToCipherKey([]byte(a.Password), Cipher.KeySize()),
		replayFilter: func() antireplay.GeneralizedReplayFilter {
			if a.IvCheck {
				return antireplay.NewBloomRing()
			}
			return nil
		}(),
	}, nil
}

// Cipher is an interface for all Shadowsocks ciphers.
type Cipher interface {
	KeySize() int32
	IVSize() int32
	NewEncryptionWriter(key []byte, iv []byte, writer io.Writer) (buf.Writer, error)
	NewDecryptionReader(key []byte, iv []byte, reader io.Reader) (buf.Reader, error)
	IsAEAD() bool
	EncodePacket(key []byte, b *buf.Buffer) error
	DecodePacket(key []byte, b *buf.Buffer) error
}

var _ Cipher = (*AEADCipher)(nil)

type AEADCipher struct {
	KeyBytes        int32
	IVBytes         int32
	AEADAuthCreator func(key []byte) cipher.AEAD
}

func (*AEADCipher) IsAEAD() bool {
	return true
}

func (c *AEADCipher) KeySize() int32 {
	return c.KeyBytes
}

func (c *AEADCipher) IVSize() int32 {
	return c.IVBytes
}

func (c *AEADCipher) createAuthenticator(key []byte, iv []byte) *crypto.AEADAuthenticator {
	subkey := make([]byte, c.KeyBytes)
	hkdfSHA1(key, iv, subkey)
	aead := c.AEADAuthCreator(subkey)
	nonce := crypto.GenerateAEADNonceWithSize(aead.NonceSize())
	return &crypto.AEADAuthenticator{
		AEAD:           aead,
		NonceGenerator: nonce,
	}
}

func (c *AEADCipher) NewEncryptionWriter(key []byte, iv []byte, writer io.Writer) (buf.Writer, error) {
	auth := c.createAuthenticator(key, iv)
	return crypto.NewAuthenticationWriter(auth, &crypto.AEADChunkSizeParser{
		Auth: auth,
	}, writer, protocol.TransferTypeStream, nil), nil
}

func (c *AEADCipher) NewDecryptionReader(key []byte, iv []byte, reader io.Reader) (buf.Reader, error) {
	auth := c.createAuthenticator(key, iv)
	return crypto.NewAuthenticationReader(auth, &crypto.AEADChunkSizeParser{
		Auth: auth,
	}, reader, protocol.TransferTypeStream, nil), nil
}

func (c *AEADCipher) EncodePacket(key []byte, b *buf.Buffer) error {
	ivLen := c.IVSize()
	payloadLen := b.Len()
	auth := c.createAuthenticator(key, b.BytesTo(ivLen))

	b.Extend(int32(auth.Overhead()))
	_, err := auth.Seal(b.BytesTo(ivLen), b.BytesRange(ivLen, payloadLen))
	return err
}

func (c *AEADCipher) DecodePacket(key []byte, b *buf.Buffer) error {
	if b.Len() <= c.IVSize() {
		return newError("insufficient data: ", b.Len())
	}
	ivLen := c.IVSize()
	payloadLen := b.Len()
	auth := c.createAuthenticator(key, b.BytesTo(ivLen))

	bbb, err := auth.Open(b.BytesTo(ivLen), b.BytesRange(ivLen, payloadLen))
	if err != nil {
		return err
	}
	b.Resize(ivLen, int32(len(bbb)))
	return nil
}

type StreamCipher struct {
	KeyBytes       int32
	IVBytes        int32
	EncryptCreator func(key []byte, iv []byte) (cipher.Stream, error)
	DecryptCreator func(key []byte, iv []byte) (cipher.Stream, error)
}

func blockStream(blockCreator func(key []byte) (cipher.Block, error), streamCreator func(block cipher.Block, iv []byte) cipher.Stream) func([]byte, []byte) (cipher.Stream, error) {
	return func(key []byte, iv []byte) (cipher.Stream, error) {
		block, err := blockCreator(key)
		if err != nil {
			return nil, err
		}
		return streamCreator(block, iv), err
	}
}

func (*StreamCipher) IsAEAD() bool {
	return false
}

func (v *StreamCipher) KeySize() int32 {
	return v.KeyBytes
}

func (v *StreamCipher) IVSize() int32 {
	return v.IVBytes
}

func (v *StreamCipher) NewEncryptionWriter(key []byte, iv []byte, writer io.Writer) (buf.Writer, error) {
	stream, err := v.EncryptCreator(key, iv)
	if err != nil {
		return nil, err
	}
	return &buf.SequentialWriter{Writer: crypto.NewCryptionWriter(stream, writer)}, nil
}

func (v *StreamCipher) NewDecryptionReader(key []byte, iv []byte, reader io.Reader) (buf.Reader, error) {
	stream, err := v.DecryptCreator(key, iv)
	if err != nil {
		return nil, err
	}
	return &buf.SingleReader{Reader: crypto.NewCryptionReader(stream, reader)}, nil
}

func (v *StreamCipher) EncodePacket(key []byte, b *buf.Buffer) error {
	iv := b.BytesTo(v.IVSize())
	stream, err := v.EncryptCreator(key, iv)
	if err != nil {
		return err
	}
	stream.XORKeyStream(b.BytesFrom(v.IVSize()), b.BytesFrom(v.IVSize()))
	return nil
}

func (v *StreamCipher) DecodePacket(key []byte, b *buf.Buffer) error {
	if b.Len() <= v.IVSize() {
		return newError("insufficient data: ", b.Len())
	}
	iv := b.BytesTo(v.IVSize())
	stream, err := v.DecryptCreator(key, iv)
	if err != nil {
		return err
	}
	stream.XORKeyStream(b.BytesFrom(v.IVSize()), b.BytesFrom(v.IVSize()))
	b.Advance(v.IVSize())
	return nil
}

var _ Cipher = (*NoneCipher)(nil)

type NoneCipher struct{}

func (*NoneCipher) KeySize() int32 { return 16 }
func (*NoneCipher) IVSize() int32  { return 0 }
func (*NoneCipher) IsAEAD() bool {
	return false
}

func (*NoneCipher) NewDecryptionReader(key []byte, iv []byte, reader io.Reader) (buf.Reader, error) {
	return buf.NewReader(reader), nil
}

func (*NoneCipher) NewEncryptionWriter(key []byte, iv []byte, writer io.Writer) (buf.Writer, error) {
	return buf.NewWriter(writer), nil
}

func (*NoneCipher) EncodePacket(key []byte, b *buf.Buffer) error {
	return nil
}

func (*NoneCipher) DecodePacket(key []byte, b *buf.Buffer) error {
	return nil
}

func CipherFromString(c string) CipherType {
	c = strings.ToUpper(c)
	c = strings.ReplaceAll(c, "-", "_")
	if c == "PLAIN" {
		c = "NONE"
	}
	if c == "CHACHA20_POLY1305" {
		c = "CHACHA20_IETF_POLY1305"
	}
	return CipherType(CipherType_value[c])
}

func passwordToCipherKey(password []byte, keySize int32) []byte {
	const md5Len = 16

	cnt := (int(keySize)-1)/md5Len + 1
	m := make([]byte, cnt*md5Len)
	copy(m, md5sum(password))

	// Repeatedly call md5 until bytes generated is enough.
	// Each call to md5 uses data: prev md5 sum + password.
	d := make([]byte, md5Len+len(password))
	start := 0
	for i := 1; i < cnt; i++ {
		start += md5Len
		copy(d, m[start-md5Len:start])
		copy(d[md5Len:], password)
		copy(m[start:], md5sum(d))
	}
	return m[:keySize]
}

func md5sum(d []byte) []byte {
	h := md5.New()
	h.Write(d)
	return h.Sum(nil)
}

func hkdfSHA1(secret, salt, outKey []byte) {
	r := hkdf.New(sha1.New, secret, salt, []byte("ss-subkey"))
	common.Must2(io.ReadFull(r, outKey))
}
