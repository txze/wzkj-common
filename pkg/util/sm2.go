package util

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"sort"
	"strings"

	"github.com/tjfoc/gmsm/x509"
)

// Hex 排序并组装签名明文串
func Hex(m map[string]interface{}) string {
	// 移除signature和url参数
	delete(m, "signature")
	delete(m, "url")

	// 获取所有key并排序
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 组装签名串
	var sb strings.Builder
	for i, k := range keys {
		if i > 0 {
			sb.WriteString("&")
		}
		sb.WriteString(fmt.Sprintf("%s=%v", k, m[k]))
	}

	return sb.String()
}

// SM2Sign SM2签名
func SM2Sign(privateKey string, data string) (string, error) {
	// 解码私钥
	privateKeyBytes, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to decode private key: %v", err)
	}

	privKey, err := x509.ParsePKCS8UnecryptedPrivateKey(privateKeyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %v", err)
	}

	// 生成签名
	signature, err := privKey.Sign(rand.Reader, []byte(data), nil)
	if err != nil {
		return "", fmt.Errorf("failed to sign data: %v", err)
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

// SM2Verify SM2验签
func SM2Verify(publicKey string, data string, signature string) (bool, error) {
	pubbytes, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return false, fmt.Errorf("failed to decode public key: %v", err)
	}

	pubKey, err := x509.ParseSm2PublicKey(pubbytes)
	if err != nil {
		return false, fmt.Errorf("failed to parse public key: %v", err)
	}

	// 解码签名
	signBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false, fmt.Errorf("failed to decode signature: %v", err)
	}

	// 验证签名
	return pubKey.Verify([]byte(data), signBytes), nil
}
