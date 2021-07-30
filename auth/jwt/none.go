package jwt

import (
	"github.com/jinycoo/jinygo/errors"
	"github.com/jinycoo/jinygo/log"
)

// Implements the none signing method.  This is required by the spec
// but you probably should never use it.
var SigningMethodNone *signingMethodNone

const UnsafeAllowNoneSignatureType unsafeNoneMagicConstant = "none signing method allowed"

type signingMethodNone struct{}
type unsafeNoneMagicConstant string

func init() {
	SigningMethodNone = &signingMethodNone{}

	RegisterSigningMethod(SigningMethodNone.Alg(), func() SigningMethod {
		return SigningMethodNone
	})
}

func (m *signingMethodNone) Alg() string {
	return "none"
}

// Only allow 'none' alg type if UnsafeAllowNoneSignatureType is specified as the key
func (m *signingMethodNone) Verify(signingString, signature string, key interface{}) (err error) {
	// Key must be UnsafeAllowNoneSignatureType to prevent accidentally
	// accepting 'none' signing method
	if _, ok := key.(unsafeNoneMagicConstant); !ok {
		log.ZError("Signature validation failed - 'none' signature type is not allowed")
		return errors.TokenSignatureInvalid
	}
	// If signing method is none, signature must be an empty string
	if signature != "" {
		log.ZError("Signature validation failed - 'none' signing method with non-empty signature")
		return errors.TokenSignatureInvalid
	}

	// Accept 'none' signing method.
	return nil
}

// Only allow 'none' signing if UnsafeAllowNoneSignatureType is specified as the key
func (m *signingMethodNone) Sign(signingString string, key interface{}) (string, error) {
	if _, ok := key.(unsafeNoneMagicConstant); ok {
		return "", nil
	}
	log.ZError("Signature validation failed - 'none' signature type is not allowed")
	return "", errors.TokenSignatureInvalid
}
