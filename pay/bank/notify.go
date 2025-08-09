package bank

import (
	"crypto"
	"crypto/md5"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"

	"github.com/pkg/errors"
)

func (p *Pay) Notify(data string) (bool, error) {
	var bodyData map[string]string
	err := json.Unmarshal([]byte(data), &bodyData)
	if err != nil {
		return false, errors.Wrapf(err, "Error parsing body data")
	}

	signStr, ok := bodyData["signature"]
	if !ok {
		return false, errors.New("Error extracting signature from body data")
	}

	// Get the content to be signed
	signContent := p.getContent(bodyData)
	isValid, err := p.verify(signContent, signStr)
	if err != nil {
		return false, errors.Wrapf(err, "Error verifying signature")
	}
	return isValid, nil
}

// verify verifies the signature using the provided public key
func (p *Pay) verify(data, sign string) (bool, error) {
	block, _ := pem.Decode([]byte(p.config.PublicKey))
	if block == nil || block.Type != "PUBLIC KEY" {
		return false, errors.New("failed to decode PEM block containing public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return false, err
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return false, errors.New("not a valid RSA public key")
	}

	decodedSign, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return false, err
	}

	hashed := md5.Sum([]byte(data))
	err = rsa.VerifyPKCS1v15(rsaPub, crypto.MD5, hashed[:], decodedSign)
	if err != nil {
		return false, err
	}
	return true, nil
}
