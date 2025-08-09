package bank

import (
	"crypto"
	"crypto/md5"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"sort"
	"strings"
)

func (p *Pay) sign(params map[string]string) (string, error) {
	content := p.getContent(params)
	block, _ := pem.Decode([]byte(p.config.PrivateKey))
	if block == nil {
		return "", errors.New("failed to parse PEM block")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", errors.New("failed to parse private key: " + err.Error())
	}

	// 确保私钥是 RSA 类型
	rsaPrivateKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return "", errors.New("not an RSA private key")
	}

	// 使用 SHA-256 哈希算法和 PKCS#1 v1.5 填充方案进行签名
	// 注意：这里使用了 nil 作为随机数生成器，通常应该使用 crypto/rand.Reader
	// 但为了与你的 PHP 代码尽可能一致（尽管 PHP 代码中的做法是不安全的），这里保留了 nil
	// 在生产环境中，请务必使用 crypto/rand.Reader
	hashed := md5.Sum([]byte(content))
	signature, err := rsa.SignPKCS1v15(nil, rsaPrivateKey, crypto.MD5, hashed[:])
	if err != nil {
		return "", errors.New("failed to sign data err:" + err.Error())
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

func (p *Pay) checkEmpty(value string) bool {
	return strings.TrimSpace(value) == ""
}

func (p *Pay) getContent(params map[string]string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		if k == "signature" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var builder strings.Builder
	for i, k := range keys {
		v := params[k]
		if !p.checkEmpty(v) && !strings.HasPrefix(v, "@") {
			if i == 0 {
				builder.WriteString(fmt.Sprintf("%s=%s", k, v))
			} else {
				builder.WriteString(fmt.Sprintf("&%s=%s", k, v))
			}
		}
	}
	return builder.String()
}
