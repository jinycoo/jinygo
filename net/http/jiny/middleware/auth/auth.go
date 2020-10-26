package auth

import "jinycoo.com/jinygo/ctime"

const (
	ReqAuthKey = "Authorization"
	ResAuthKey = "WWW-Authorization"
	Account    = "account"
)

type AccInfo struct {
	MID      int64
	Username string
	Password string
	Avatar   string
	LoginAt  int64
	Expiry   ctime.Duration
	Issuer   string
}

type Accounts map[string]string

type authPair struct {
	value string
	user  string
}

type authPairs []authPair
