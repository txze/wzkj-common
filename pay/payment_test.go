package pay

import (
	"testing"
	"time"

	"github.com/txze/wzkj-common/pay/allinpay"
)

const privateKey = `-----BEGIN PRIVATE KEY-----
MIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQDQp5A+YFBh5WTn
qcsCMlhiXZR0uo0cPImTeS7Jmr89V1F5aH/RKtI3gZvsiNjACgnSCXT2cUoZKkbJ
Gl38Dz6uFY2typqjLB1pdWlLpwHJKiOV7Ls/COkTU4lDzkjF8u+ITTWSqYu/Wvx2
P5sgPmyZvWrbNlxXDJMNeyIjEa3bqLjDE6Urh/ONxlQWRffMT+AxqyQ0cee8l63i
2kmUwuGZkEDTQ1nB0MtQMMu7CSGx5+Spl/O3xnn4FAjzBFUeCnU0UutDm+E7WdMK
+l1yDi6jwfIhQbQ2HEKN1Gd05MoRLWla7t08yzP1x3UVm5Z298FKmXxIcwbxoowr
Yo0rZ9sBAgMBAAECggEAEbvMFKEW89sNQms558vjmyic73bTe2zhvHj2MwhF7K65
K/pnsp1TFIidefL/iQLRZtqK6E8knxLqxTjKeBvLlfwa+IRZtDiRn17tPhLJohFE
yP8/wtG9DXlyFyM3KCvHk+wL+5URXYgcIOizBICJtl9U76ClJHjbHrAybIyaHCsJ
iS46GQUSFgfHNvT60FF4O6fxrTLrW9m8NDDGDKPKDLJbgydYYr8wUCQ5xhvtcGcF
UnurE2LQ8W+MD1PGV+F9ml91YegMRaGjXV4D6e5SDaoYPG+ZkZ+Tid7KFTGog0eO
JRuZ+k5NgRHksgFIgmDZz6QjJwEctPUVerUIOBDgYQKBgQDvnNOK23X5ofjT6/5I
+0BgKVSowPfos5rng31hpaISXhqRLp/MEAEwwTV0GEs/Y2iXQnOhNbVT5PgY74hn
gOitoEY2IuP90FE2y4Je+GvgOysOXykdskCvXGEUcNcZ+kcoI96N3U912Sf6FgCN
2FqnK+tiKzL09iboEy0SocOG7QKBgQDe7LgAPeh4E/p8Fjt3ycknvJY201Kg32OK
QnRtPxtavvreuqkXBTPuyfhnnVUF3WnONP/ZlWAX/8h/GCyGCN76DJno0QyKFes9
5AkrnF6PrFQRU51SWGsaKJxqRcn0npJbeVm0fYdVWcczfRQpc3CyCutgeSi8Twx9
Calwvlqt5QKBgG2z2nJnmfLpwlecY3acedPM+HKurpH+sPwwClaLk9Fe/kDcHNM7
vJ/KxaNagBEMfVVLWk9DnLpFSYV5HXVt4pmjmKGuhb2uA5DXyd+bUyB9VnAlB1kO
RGlFHTlTlFfTa4KoMXu4CGpHOvNX4XcPyCljhUgTySe4DwYPyYIPR8rdAoGAYbxN
K6X4zvSLZG3m4qzwaWCQRzc9SdTG8m4SV3dMieujV5Vk3vfj/fRE2UCsbybU5Zhs
97s65yq4f6hclOM8x0pRDDbjFYNooLjioGEtQDZgoTwUhG7Jfi2B7kHsujfvmPVK
NAy5Ed2LrXJQLaA0L4sECUb1aiIKKqPaythaL1UCgYBQcEiecxAKI7NyqTW8I52C
c8c/a8C947VDtLss+aTy//MFNiucASD8+5n3WZxwVvHpynX/OBRJIChpXP8ShVBu
fvRQU5fVlgd60pqOOJbpRerCbIi6vGLipJsXoUtkC5TRiIpcLV/5xuD3X+urCfXf
1ozFA1ltZIngnAt1fPgAKw==
-----END PRIVATE KEY-----`

const publicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0KeQPmBQYeVk56nLAjJY
Yl2UdLqNHDyJk3kuyZq/PVdReWh/0SrSN4Gb7IjYwAoJ0gl09nFKGSpGyRpd/A8+
rhWNrcqaoywdaXVpS6cBySojley7PwjpE1OJQ85IxfLviE01kqmLv1r8dj+bID5s
mb1q2zZcVwyTDXsiIxGt26i4wxOlK4fzjcZUFkX3zE/gMaskNHHnvJet4tpJlMLh
mZBA00NZwdDLUDDLuwkhsefkqZfzt8Z5+BQI8wRVHgp1NFLrQ5vhO1nTCvpdcg4u
o8HyIUG0NhxCjdRndOTKES1pWu7dPMsz9cd1FZuWdvfBSpl8SHMG8aKMK2KNK2fb
AQIDAQAB
-----END PUBLIC KEY-----`

func TestNewPayment(t *testing.T) {
	t.Run("payment", func(t *testing.T) {
		//初始化支付
		pay := NewPayment[allinpay.PayRequest, allinpay.Notify]().SetStrategy(allinpay.NewAllInPay(&allinpay.AllInPayConfig{
			AppId:      "123",
			CuSID:      "123",
			APIVersion: "123",
			PrivateKey: privateKey,
			PublicKey:  publicKey,
		}))
		rs, err := pay.Pay(&allinpay.PayRequest{
			TrxAmt:    100,
			Reqsn:     time.Now().String(),
			Validtime: time.Now().Add(time.Second * 3).Format(time.DateTime),
			NotifyUrl: "www.baidu.com",
			Body:      "你好吗？",
			Remark:    "我很好",
		})
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(rs)

		//isSign, err := pay.VerifySign(rs.Params, rs.Sign)
		//if err != nil {
		//	t.Error(err)
		//	return
		//}
		//t.Log(isSign)
	})
}
