package jwt_test

import (
	"testing"

	"github.com/txze/wzkj-common/pkg/jwt"

	"github.com/stretchr/testify/assert"
)

func TestJWT(t *testing.T) {
	var token string
	var err error
	var secret = "4bf7d08fb0c146a3b89de9dd0768f167"
	var payload = &jwt.Payload{
		UserId: "xxxx",
	}
	token, err = jwt.GenerateToken(payload, secret)
	assert.NoError(t, err)
	t.Log(token)

	_, err = jwt.ParseToken(token, secret)
	assert.NoError(t, err)
}
