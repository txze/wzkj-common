package util

import (
	"math/rand"
)

var defaultLetters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

var UpperLetters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

var DigitLetters = []rune("0123456789")

// RandomString returns a random string with a fixed length
func RandomString(n int, allowedChars ...[]rune) string {
	var letters []rune

	if len(allowedChars) == 0 {
		letters = defaultLetters
	} else {
		letters = allowedChars[0]
	}

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

// RandomStr 随机生成字符串
func RandomStr(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func RandInt64(min, max int64) int64 {
	s := rand.NewSource(Now().UnixNano())
	r := rand.New(s)
	return min + r.Int63n(max-min+1)
}

// RandInt64L 左闭右开[min, max)
func RandInt64L(min, max int64) int64 {
	s := rand.NewSource(Now().UnixNano())
	r := rand.New(s)
	return min + r.Int63n(max-min)
}
func RandFloat64(min, max float64) float64 {
	s := rand.NewSource(Now().UnixNano())
	r := rand.New(s)

	return min + r.Float64()*(max-min)
}

func RandomInt(n int) int {
	var num = 0
	for i := 0; i < n; i++ {
		v := rand.Intn(10)
		if i == 0 && v == 0 {
			i = -1
			continue
		}
		num = num*10 + v
	}
	return num
}
