package jwt

import (
	"fmt"
)

// Error constants
var (
	ErrHashUnavailable = new(HashUnavailableError)
)

// Embeds b within a, if a is a valid wrapper. returns a
// If a is not a valid wrapper, b is dropped
// If one of the errors is nil, the other is returned
func wrapError(a, b error) error {
	if b == nil {
		return a
	}
	if a == nil {
		return b
	}

	type iErrorWrapper interface {
		Wrap(error)
	}
	if w, ok := a.(iErrorWrapper); ok {
		w.Wrap(b)
	}
	return a
}

// ErrorWrapper provides a simple, concrete helper for implementing nestable errors
type ErrorWrapper struct{ err error }

// Unwrap implements xerrors.Wrapper
func (w ErrorWrapper) Unwrap() error {
	return w.err
}

// Wrap stores the provided error value and returns it when Unwrap is called
func (w ErrorWrapper) Wrap(err error) {
	w.err = err
}

// InvalidKeyError is returned if the key is unusable for some reason other than type
type InvalidKeyError struct {
	Message string
	ErrorWrapper
}

func (e *InvalidKeyError) Error() string {
	return fmt.Sprintf("key is invalid: %v", e.Message)
}

// InvalidKeyTypeError is returned if the key is unusable because it is of an incompatible type
type InvalidKeyTypeError struct {
	Expected, Received string // String descriptions of expected and received types
	ErrorWrapper
}

func (e *InvalidKeyTypeError) Error() string {
	if e.Expected == "" && e.Received == "" {
		return "key is of invalid type"
	}
	return fmt.Sprintf("key is of invalid type: expected %v, received %v", e.Expected, e.Received)
}

// NewInvalidKeyTypeError creates an InvalidKeyTypeError, automatically capturing the type
// of received
func NewInvalidKeyTypeError(expected string, received interface{}) error {
	return &InvalidKeyTypeError{Expected: expected, Received: fmt.Sprintf("%T", received)}
}

// MalformedTokenError means the token failed to parse or exhibits some other
// non-standard property that prevents it being processed by this library
type MalformedTokenError struct {
	Message string
	ErrorWrapper
}

func (e *MalformedTokenError) Error() string {
	if e.Message == "" {
		return "token is malformed"
	}
	return fmt.Sprintf("token is malformed: %v", e.Message)
}

// UnverfiableTokenError means there's something wrong with the signature that prevents
// this library from verifying it.
type UnverfiableTokenError struct {
	Message string
	ErrorWrapper
}

func (e *UnverfiableTokenError) Error() string {
	if e.Message == "" {
		return "token is unverifiable"
	}
	return fmt.Sprintf("token is unverifiable: %v", e.Message)
}

// InvalidSignatureError means the signature on the token is invalid
type InvalidSignatureError struct {
	Message string
	ErrorWrapper
}

func (e *InvalidSignatureError) Error() string {
	if e.Message == "" {
		return "token signature is invalid"
	}
	return fmt.Sprintf("token signature is invalid: %v", e.Message)
}

// InvalidAudienceError means the token failed the audience check
// per the spec, if an 'aud' claim is present, the value must be verified
// See: WithAudience and WithoutAudienceValidation
type InvalidAudienceError struct {
	Message string
	ErrorWrapper
}

func (e *InvalidAudienceError) Error() string {
	if e.Message == "" {
		return "token audience is invalid"
	}
	return fmt.Sprintf("token audience is invalid: %v", e.Message)
}

// InvalidIssuerError means the token failed issuer validation
// Issuer validation is only run, by default, if the WithIssuer option is provided
type InvalidIssuerError struct {
	Message string
	ErrorWrapper
}

func (e *InvalidIssuerError) Error() string {
	if e.Message == "" {
		return "token issuer is invalid"
	}
	return fmt.Sprintf("token issuer is invalid: %v", e.Message)
}

type HashUnavailableError struct {
	ErrorWrapper
}

func (e *HashUnavailableError) Error() string {
	return "the requested hash function is unavailable"
}
