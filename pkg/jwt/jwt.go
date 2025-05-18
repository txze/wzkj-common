package jwt

import (
	"errors"
	"wzkj-common/pkg/util"

	"github.com/dgrijalva/jwt-go"
)

type Payload struct {
	UserId string
	Name   string
	Expire int64
	Role   string
}

func (p *Payload) Valid() error {
	// 预设，如果过期时间小于等于0，代表token不过期
	if p.Expire <= 0 {
		return nil
	}
	//目前这里想到的是验证超时
	if p.Expire < util.Now().Unix() {
		return errors.New("token expire")
	}
	return nil
}
func GenerateToken(pl *Payload, secretKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, pl)
	return token.SignedString([]byte(secretKey))
}

func ParseToken(tokenStr string, secretKey string) (*Payload, error) {
	var pl Payload
	token, err := jwt.ParseWithClaims(tokenStr, &pl, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Payload); ok {
		return claims, nil
	}
	return nil, errors.New("not type")
}
