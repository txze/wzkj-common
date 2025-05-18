package wechat

import (
	"github.com/silenceper/wechat/oauth"
)

type Oauth struct {
	*oauth.Oauth
}

func (o *Oauth) UserInfo(code string) (result oauth.UserInfo, err error) {
	accessToken, err := o.GetUserAccessToken(code)
	if err != nil {
		return
	}

	result, err = o.GetUserInfo(accessToken.AccessToken, accessToken.OpenID)
	if err != nil {
		return
	}
	return
}
