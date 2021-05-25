package jiny

import "jinycoo.com/jinygo/ctime"

type Config struct {
	Port          string
	AllowHosts    []string
	AllowPatterns []string
	SignPaths     []string
	Headers       map[string]string

	Prefix     string
	SigningKey string
	SignActive bool
	Sign       *Sign
	Expiry     ctime.Duration
	ReExpiry   ctime.Duration
	CasServer  string
}

type Sign struct {
	AppID   string
	PubKeys []string
}
