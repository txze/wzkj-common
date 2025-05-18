package util

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
)

func MD5Hex(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func Sha1(str string) (hex string) {
	if str == "" {
		return
	}

	hex = fmt.Sprintf("%x", sha1.Sum([]byte(str)))
	return
}

func Sha512(str string) string {
	if str == "" {
		return ""
	}
	sha_h := sha512.New()
	sha_h.Write([]byte(str))
	return fmt.Sprintf("%x", sha_h.Sum(nil))
}

func Sha256(str string) string {
	if str == "" {
		return ""
	}
	sha_h := sha256.New()
	sha_h.Write([]byte(str))
	return fmt.Sprintf("%x", sha_h.Sum(nil))
}
