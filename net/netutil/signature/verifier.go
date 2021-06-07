package signature

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"strings"

	"jinycoo.com/jinygo/errors"
	"jinycoo.com/jinygo/log"
)

const (
	// 私钥 PEMBEGIN 开头
	PEMBEGIN = "-----BEGIN RSA PRIVATE KEY-----\n"
	// 私钥 PEMEND 结尾
	PEMEND = "\n-----END RSA PRIVATE KEY-----"
	// 公钥 PEMBEGIN 开头
	PUBPEMBEGIN = "-----BEGIN PUBLIC KEY-----\n"
	// 公钥 PEMEND 结尾
	PUBPEMEND = "\n-----END PUBLIC KEY-----"
)

type Verifier struct {
	V Verify
}

func NewVerifier(v Verify) *Verifier {
	if v == nil {
		v = rsa2Verify
	}
	return &Verifier{
		V: v,
	}
}

func NewRSA2Verifier() *Verifier {
	return &Verifier{
		V: rsa2Verify,
	}
}

func (s *Verifier) Verify(content, sign, pubKey string) (err error) {
	return s.V(content, sign, pubKey)
}

func rsa2Verify(content, sign, pubKey string) (err error) {
	return Rsa2PubSign(content, sign, pubKey, crypto.SHA256)
}

// Rsa2PubSign RSA2公钥验证签名
func Rsa2PubSign(signContent, sign, publicKey string, hash crypto.Hash) (err error) {
	hashed := sha256.Sum256([]byte(signContent))
	pubKey, err := ParsePublicKey(publicKey)
	if err != nil {
		log.Errorf("rsa2 public key check failed. %v", err)
		return err
	}
	sig, _ := base64.StdEncoding.DecodeString(sign)
	err = rsa.VerifyPKCS1v15(pubKey, hash, hashed[:], sig)
	if err != nil {
		log.Errorf("rsa2 public check sign failed. %v", err)
	}
	return err
}

// ParsePublicKey 公钥验证
func ParsePublicKey(publicKey string) (*rsa.PublicKey, error) {
	publicKey = FormatPublicKey(publicKey)
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		return nil, errors.New("公钥信息错误！")
	}
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pubKey.(*rsa.PublicKey), nil
}

// FormatPublicKey 组装公钥
func FormatPublicKey(publicKey string) string {
	if !strings.HasPrefix(publicKey, PUBPEMBEGIN) {
		publicKey = PUBPEMBEGIN + publicKey
	}
	if !strings.HasSuffix(publicKey, PUBPEMEND) {
		publicKey = publicKey + PUBPEMEND
	}
	return publicKey
}
