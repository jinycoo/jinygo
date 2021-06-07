package signature

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"

	"jinycoo.com/jinygo/errors"
)

// Sign 加签函数
type Sign func(content, privateKey string) (sign string, err error)

// Verify 验签函数
type Verify func(content, sign, pubKey string) (err error)

// NewSigner 初始化Signer，默认RSA2签名
func NewSigner(s Sign) *Signer {
	if s == nil {
		s = rsa2Sign
	}
	return &Signer{
		S: s,
	}
}

var ErrPemDecode = errors.New("pem.Decode failed") // pem解析失败

func rsa2Sign(content, privateKey string) (sign string, err error) {
	// 1、将密钥解析成密钥实例
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		err = ErrPemDecode
		return
	}
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return
	}

	// 2、生成签名
	hash := sha256.New()
	_, err = hash.Write([]byte(content))
	if err != nil {
		return
	}
	signature, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, hash.Sum(nil))
	if err != nil {
		return
	}

	// 3、签名base64编码
	sign = base64.StdEncoding.EncodeToString(signature)
	return
}
