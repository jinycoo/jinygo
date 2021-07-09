package jwt

import (
	"bytes"
	"encoding/json"
	"github.com/jinycoo/jinygo/errors"
	"strings"
)

// Parser is the type used to parse and validate a JWT token from string
type Parser struct {
	validMethods         []string          // If populated, only these methods will be considered valid
	useJSONNumber        bool              // Use JSON Number format in JSON decoder
	skipClaimsValidation bool              // Skip claims validation during token parsing
	unmarshaller         TokenUnmarshaller // Use this instead of encoding/json
	*ValidationHelper
}

// NewParser returns a new Parser with the specified options
func NewParser(options ...ParserOption) *Parser {
	p := &Parser{
		ValidationHelper: new(ValidationHelper),
	}
	for _, option := range options {
		option(p)
	}
	return p
}

// Parse will parse, validate, and return a token.
// keyFunc will receive the parsed token and should return the key for validating.
// If everything is kosher, err will be nil
func (p *Parser) Parse(tokenString string, keyFunc Keyfunc) (*Token, error) {
	return p.ParseWithClaims(tokenString, MapClaims{}, keyFunc)
}

// ParseWithClaims is just like parse, but with the claims type specified
func (p *Parser) ParseWithClaims(tokenString string, claims Claims, keyFunc Keyfunc) (*Token, error) {
	token, parts, err := p.ParseUnverified(tokenString, claims)
	if err != nil {
		return token, err
	}

	// Verify signing method is in the required set
	if p.validMethods != nil {
		var signingMethodValid = false
		var alg = token.Method.Alg()
		for _, m := range p.validMethods {
			if m == alg {
				signingMethodValid = true
				break
			}
		}
		if !signingMethodValid {
			// signing method is not in the listed set
			return token, errors.AuthJwtAlgInvalid
		}
	}

	// Lookup key
	var key interface{}
	if keyFunc == nil {
		// keyFunc was not provided.  short circuiting validation
		return token, &UnverfiableTokenError{Message: "no Keyfunc was provided."}
	}
	if key, err = keyFunc(token); err != nil {
		// keyFunc returned an error
		return token, wrapError(&UnverfiableTokenError{Message: "Keyfunc returned an error"}, err)
	}

	var vErr error

	// Perform validation
	token.Signature = parts[2]
	if err = token.Method.Verify(strings.Join(parts[0:2], "."), token.Signature, key); err != nil {
		vErr = errors.AuthTokenInvalid
	}

	// Validate Claims
	if !p.skipClaimsValidation && vErr == nil {
		if err := token.Claims.Valid(p.ValidationHelper); err != nil {
			vErr = wrapError(err, vErr)
		}
	}

	if vErr == nil {
		token.Valid = true
	}

	return token, vErr
}

// ParseUnverified is used to inspect a token without validating it
// WARNING: Don't use this method unless you know what you're doing
//
// This method parses the token but doesn't validate the signature. It's only
// ever useful in cases where you know the signature is valid (because it has
// been checked previously in the stack) and you want to extract values from
// it. Or for debuggery.
func (p *Parser) ParseUnverified(tokenString string, claims Claims) (token *Token, parts []string, err error) {
	parts = strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, parts, errors.AuthJwtSegmentsNumInvalid
	}

	token = &Token{Raw: tokenString}

	// choose unmarshaller
	var unmarshaller = p.unmarshaller
	if unmarshaller == nil {
		unmarshaller = p.defaultUnmarshaller
	}

	// parse Header
	var headerBytes []byte
	if headerBytes, err = DecodeSegment(parts[0]); err != nil {
		if strings.HasPrefix(strings.ToLower(tokenString), "bearer ") {
			return token, parts, errors.AuthJwtBearerErr
		}
		return token, parts, errors.AuthJwtHeaderInvalid
	}
	if err = unmarshaller(CodingContext{HeaderFieldDescriptor, nil}, headerBytes, &token.Header); err != nil {
		return token, parts, errors.AuthJwtHeaderUnmarshalErr
	}

	// parse Claims
	var claimBytes []byte
	token.Claims = claims

	if claimBytes, err = DecodeSegment(parts[1]); err != nil {
		return token, parts, errors.AuthJwtClaimsErr
	}
	// JSON Decode.  Special case for map type to avoid weird pointer behavior
	ctx := CodingContext{ClaimsFieldDescriptor, token.Header}
	if c, ok := token.Claims.(MapClaims); ok {
		err = unmarshaller(ctx, claimBytes, &c)
	} else {
		err = unmarshaller(ctx, claimBytes, &claims)
	}
	// Handle decode error
	if err != nil {
		return token, parts, errors.AuthJwtClaimsUnmarshalErr
	}

	// Lookup signature method
	if method, ok := token.Header["alg"].(string); ok {
		if token.Method = GetSigningMethod(method); token.Method == nil {
			return token, parts, errors.AuthJwtAlgUnverfiable
		}
	} else {
		return token, parts, errors.AuthJwtAlgUnspecified
	}

	return token, parts, nil
}

func (p *Parser) defaultUnmarshaller(ctx CodingContext, data []byte, v interface{}) error {
	// If we don't need a special parser, use Unmarshal
	// We never use a special encoder for the header
	if !p.useJSONNumber || ctx.FieldDescriptor == HeaderFieldDescriptor {
		return json.Unmarshal(data, v)
	}

	// To enable the JSONNumber mode, we must use Decoder instead of Unmarshal
	dec := json.NewDecoder(bytes.NewBuffer(data))
	dec.UseNumber()
	return dec.Decode(v)
}
