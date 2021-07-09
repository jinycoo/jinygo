package request

import (
	"github.com/jinycoo/jinygo/auth/jwt"
	"net/url"
)

var requestTestData = []struct {
	name      string
	claims    jwt.MapClaims
	extractor Extractor
	headers   map[string]string
	query     url.Values
	valid     bool
}{
	{
		"authorization bearer token",
		jwt.MapClaims{"foo": "bar"},
		AuthorizationHeaderExtractor,
		map[string]string{"Authorization": "Bearer %v"},
		url.Values{},
		true,
	},
	{
		"oauth bearer token - header",
		jwt.MapClaims{"foo": "bar"},
		OAuth2Extractor,
		map[string]string{"Authorization": "Bearer %v"},
		url.Values{},
		true,
	},
	{
		"oauth bearer token - url",
		jwt.MapClaims{"foo": "bar"},
		OAuth2Extractor,
		map[string]string{},
		url.Values{"access_token": {"%v"}},
		true,
	},
	{
		"url token",
		jwt.MapClaims{"foo": "bar"},
		ArgumentExtractor{"token"},
		map[string]string{},
		url.Values{"token": {"%v"}},
		true,
	},
}
