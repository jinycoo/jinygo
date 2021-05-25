/**------------------------------------------------------------**
 * @filename auth/jwt.go
 * @author   jinycoo
 * @version  1.0.0
 * @date     2019-07-24 14:15
 * @desc     auth - jwt token
 **------------------------------------------------------------**/
package jiny

import (
	"time"

	"jinycoo.com/jinygo/auth/jwt"
	"jinycoo.com/jinygo/ctime"
	"jinycoo.com/jinygo/errors"
	"jinycoo.com/jinygo/log"
	"jinycoo.com/jinygo/net/http/jiny/server"
)

const (
	Authorization = "Authorization"
	UCKey         = "78abffcf33c17bf9f3aef9be2a17284e6d2a1909"
	UCID          = "97d1474d94350970168b6ccb02bda90a585b69d0"
	SigningKey    = "aa9f16c7e942d86bd0738cb7cf61924a"
)

type JWT struct {
	SigningKey []byte
}

type Claims struct {
	ID       int64  `json:"id"`
	MID      int64  `json:"mid"`
	Username string `json:"username"`
	jwt.StandardClaims
}

type AccInfo struct {
	ID       int64
	Email    string
	MID      int64
	Username string
	LoginAt  time.Time
	Issuer   string
}

func NewJWT() (j *JWT) {
	sk := conf.SigningKey
	if len(sk) == 0 {
		sk = SigningKey
	}
	return &JWT{[]byte(sk)}
}

func AccessToken(acc *AccInfo, redura *ctime.Duration) (string, error) {
	if acc == nil {
		acc = new(AccInfo)
	}
	if conf.Expiry == 0 {
		conf.Expiry = ctime.Duration(8 * time.Hour)
	}
	exp := time.Duration(conf.Expiry)
	if redura != nil {
		exp = time.Duration(*redura)
	}
	now := time.Now()
	if acc.LoginAt.Unix() == -62135596800 || acc.LoginAt.Add(exp).Before(now) {
		acc.LoginAt = now
	}
	claims := Claims{
		ID:       acc.ID,
		MID:      acc.MID,
		Username: acc.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(acc.LoginAt.Add(exp)),
			Issuer:    acc.Issuer,
		},
	}
	return NewJWT().Gen(claims)
}

func (j *JWT) Gen(claims Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	token.Header["kid"] = j.SigningKey
	return token.SignedString(j.SigningKey)
}

func (j *JWT) Parse(c *server.Context) (*Claims, error) {
	token, err := jwt.ParseWithClaims(c.GetHeader(Authorization), &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if token != nil {
		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			return claims, nil
		}
	}
	log.Warnf("jwt token err(%v)", err.Error())
	return nil, err
}

func (j *JWT) Refresh(tokenStr string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		claims.StandardClaims.ExpiresAt = jwt.At(time.Now().Add(time.Duration(conf.Expiry)))
		return j.Gen(*claims)
	}
	return "", errors.TokenInvalid
}

func JwtAuth() server.HandlerFn {
	return func(c *server.Context) {
		claims, err := NewJWT().Parse(c)
		if err != nil {
			c.JSON(nil, errors.AuthTokenErr)
			c.Abort()
			return
		}
		c.Set(UCKey, claims.MID)
		c.Set(UCID, claims.ID)
	}
}
