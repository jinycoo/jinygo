package jwt

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/jinycoo/jinygo/errors"
	"github.com/jinycoo/jinygo/log"
)

type Parser struct {
	ValidMethods         []string // If populated, only these methods will be considered valid
	UseJSONNumber        bool     // Use JSON Number format in JSON decoder
	SkipClaimsValidation bool     // Skip claims validation during token parsing
}

// Parse, validate, and return a token.
// keyFunc will receive the parsed token and should return the key for validating.
// If everything is kosher, err will be nil
func (p *Parser) Parse(tokenString string, keyFunc Keyfunc) (*Token, error) {
	return p.ParseWithClaims(tokenString, MapClaims{}, keyFunc)
}

func (p *Parser) ParseWithClaims(tokenString string, claims Claims, keyFunc Keyfunc) (*Token, error) {
	token, parts, err := p.ParseUnverified(tokenString, claims)
	if err != nil {
		return token, err
	}

	// Verify signing method is in the required set
	if p.ValidMethods != nil {
		var signingMethodValid = false
		var alg = token.Method.Alg()
		for _, m := range p.ValidMethods {
			if m == alg {
				signingMethodValid = true
				break
			}
		}
		if !signingMethodValid {
			log.Errorf("Signature validation failed - signing method %v is invalid", alg)
			// signing method is not in the listed set
			return token, errors.TokenSignatureInvalid
		}
	}

	// Lookup key
	var key interface{}
	if keyFunc == nil {
		log.ZError("Token could not be verified because of signing problems - no Keyfunc was provided.")
		// keyFunc was not provided.  short circuiting validation
		return token, errors.TokenSigningErr
	}
	if key, err = keyFunc(token); err != nil {
		// keyFunc returned an error
		//if ve, ok := err.(*ValidationError); ok {
		//	return token, ve
		//}
		return token, errors.TokenSigningErr
	}

	// Validate Claims
	if !p.SkipClaimsValidation {
		if err := token.Claims.Valid(); err != nil {
			return token, err
		}
	}

	// Perform validation
	token.Signature = parts[2]
	if err = token.Method.Verify(strings.Join(parts[0:2], "."), token.Signature, key); err != nil {
		return token, err
	}

	return token, nil
}

// WARNING: Don't use this method unless you know what you're doing
//
// This method parses the token but doesn't validate the signature. It's only
// ever useful in cases where you know the signature is valid (because it has
// been checked previously in the stack) and you want to extract values from
// it.
func (p *Parser) ParseUnverified(tokenString string, claims Claims) (token *Token, parts []string, err error) {
	parts = strings.Split(tokenString, ".")
	if len(parts) != 3 {
		log.ZError("Token is malformed - token contains an invalid number of segments")
		return nil, parts, errors.TokenMalformed
	}

	token = &Token{Raw: tokenString}

	// parse Header
	var headerBytes []byte
	if headerBytes, err = DecodeSegment(parts[0]); err != nil {
		if strings.HasPrefix(strings.ToLower(tokenString), "bearer ") {
			log.ZError("Token is malformed - tokenstring should not contain 'bearer '")
			return token, parts, errors.TokenMalformed
		}
		log.Errorf("Token is malformed - %v", err)
		return token, parts, errors.TokenMalformed
	}
	if err = json.Unmarshal(headerBytes, &token.Header); err != nil {
		log.Errorf("Token is malformed json unmarshal - %v", err)
		return token, parts, errors.TokenMalformed
	}

	// parse Claims
	var claimBytes []byte
	token.Claims = claims

	if claimBytes, err = DecodeSegment(parts[1]); err != nil {
		log.Errorf("Token is malformed - %v", err)
		return token, parts, errors.TokenMalformed
	}
	dec := json.NewDecoder(bytes.NewBuffer(claimBytes))
	if p.UseJSONNumber {
		dec.UseNumber()
	}
	// JSON Decode.  Special case for map type to avoid weird pointer behavior
	if c, ok := token.Claims.(MapClaims); ok {
		err = dec.Decode(&c)
	} else {
		err = dec.Decode(&claims)
	}
	// Handle decode error
	if err != nil {
		log.Errorf("Token is malformed - %v", err)
		return token, parts, errors.TokenMalformed
	}

	// Lookup signature method
	if method, ok := token.Header["alg"].(string); ok {
		if token.Method = GetSigningMethod(method); token.Method == nil {
			log.ZError("Token could not be verified because of signing problems - signing method (alg) is unavailable.")
			return token, parts, errors.TokenSigningErr
		}
	} else {
		log.ZError("Token could not be verified because of signing problems - signing method (alg) is unspecified.")
		return token, parts, errors.TokenSigningErr
	}

	return token, parts, nil
}
