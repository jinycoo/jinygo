/**------------------------------------------------------------**
 * @filename auth/jwt.go
 * @author   jinycoo
 * @version  1.0.0
 * @date     2019-07-24 14:15
 * @desc     auth - jwt token
 **------------------------------------------------------------**/
package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"

	"jinycoo.com/jinygo/errors"
	"jinycoo.com/jinygo/net/http/jiny"
)

const UCKey = "UkID"

var SigningKey = "api.jinycoo.com"

type JWT struct {
	SigningKey []byte
}

type Claims struct {
	MID      int64  `json:"mid"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	LoginAt  int64  `json:"login_at"`
	jwt.StandardClaims
}

func NewJWT() (j *JWT) {
	return &JWT{[]byte(GetSigningKey())}
}

func GetSigningKey() string {
	return SigningKey
}

func SetSigningKey(key string) {
	SigningKey = key
}

func BuildAccToken(acc *AccInfo) (string, error) {
	if acc == nil {
		acc = new(AccInfo)
	}
	exp := int64(time.Duration(acc.Expiry))
	if acc.LoginAt == 0 || acc.LoginAt+exp < time.Now().Unix() {
		acc.LoginAt = time.Now().Unix()
	}
	claims := Claims{
		MID:      acc.MID,
		Username: acc.Username,
		Avatar:   acc.Avatar,
		LoginAt:  acc.LoginAt,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: acc.LoginAt + exp,
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

func (j *JWT) Parse(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errors.TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, errors.TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, errors.TokenNotValidYet
			} else {
				return nil, errors.TokenInvalid
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			return claims, nil
		}
	}
	return nil, errors.AccessKeyErr
}

func (j *JWT) Refresh(tokenStr string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return j.Gen(*claims)
	}
	return "", errors.TokenInvalid
}

func JwtAuth() jiny.HandlerFn {
	return func(c *jiny.Context) {
		tokenStr := c.Request.Header.Get(ReqAuthKey)
		if len(tokenStr) == 0 {
			c.JSON(nil, errors.Unauthorized)
			c.Abort()
			return
		}
		j := NewJWT()
		claims, err := j.Parse(tokenStr)
		if err != nil {
			c.JSON(nil, err)
			c.Abort()
			return
		}
		c.Set(UCKey, claims.MID)
	}
}
