package util

import (
	"database/sql/driver"
	"encoding/json"
	"math"
	"math/big"
)

type M map[string]interface{}

func (m M) Value() (driver.Value, error) {
	valueString, err := json.Marshal(m)
	return string(valueString), err
}

func (m *M) Scan(value interface{}) error {
	if err := json.Unmarshal(value.([]byte), &m); err != nil {
		return err
	}
	return nil
}

// PutIntoMapIf 根据条件放进map
func PutIntoMapIf(m map[string]interface{}, cond bool, key string, value interface{}) {
	if cond {
		m[key] = value
	}
}

// S2Json struct to json
func S2Json(v interface{}) string {
	buf, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(buf)
}

// Json2S trans json to object
func Json2S(src string, dest interface{}) error {
	return json.Unmarshal([]byte(src), dest)
}

func S2S(data, v interface{}) error {
	return Json2S(S2Json(data), v)
}

func Bytes2S(src []byte, desc interface{}) error {
	return json.Unmarshal(src, desc)
}

func S2Bytes(src interface{}) ([]byte, error) {
	return json.Marshal(src)
}

// BigIntToFloat64 大整数转为float64
func BigIntToFloat64(v *big.Int) float64 {
	a := &big.Float{}
	f, _ := a.SetInt(v).Float64()
	return f
}

// Float64ToBigInt float6¢ 转为大整数
func Float64ToBigInt(v float64) *big.Int {
	f := big.NewFloat(v)
	r, _ := f.Int(nil)
	return r
}

// StringToU64 string to uint64
func StringToU64(s string) uint64 {
	bi, ok := (&big.Int{}).SetString(s, 10)
	if ok {
		return bi.Uint64()
	}
	return 0
}

// Gas 计算以太坊gas， gasPrice 单位是gwei
func Gas(gasLimit uint64, gasPrice uint64) *big.Int {
	limit := (&big.Int{}).SetUint64(gasLimit)
	price := (&big.Int{}).Mul((&big.Int{}).SetUint64(gasPrice), big.NewInt(1e9))

	return (&big.Int{}).Mul(limit, price)
}

// LeftPadString 字符串前面填充到指定到个数
func LeftPadString(s, padding string, total int) string {
	if len(s) >= total {
		return s
	}

	count := total - len(s)
	for i := 0; i < count; i++ {
		s = padding + s
	}

	return s
}

// AddStrInt string int add
func AddStrInt(a, b string) *big.Int {
	ai, _ := (&big.Int{}).SetString(a, 10)
	bi, _ := (&big.Int{}).SetString(b, 10)

	sum := (big.NewInt(0)).Add(ai, bi)
	return sum
}

// SubStrInt return a-b
func SubStrInt(a, b string) *big.Int {
	ai, _ := (&big.Int{}).SetString(a, 10)
	bi, _ := (&big.Int{}).SetString(b, 10)

	diff := (big.NewInt(0)).Sub(ai, bi)
	return diff
}

// F64ToU64 float64 转 uint64 并且乘以给定到小数点位数
func F64ToU64(value float64, decimal int) uint64 {
	return uint64(value * math.Pow10(decimal))
}

// FormatUsdt 格式化usdt
func FormatUsdt(v float64) float64 {
	f := new(big.Float).Quo(big.NewFloat(v), big.NewFloat(1e6))
	res, _ := f.Float64()
	return res
}

// FormatEther 格式化以太币
func FormatEther(v string) string {
	i, _ := new(big.Int).SetString(v, 10)
	i = i.Div(i, big.NewInt(1e9))

	f := new(big.Float).SetUint64(i.Uint64())

	f = f.Quo(f, big.NewFloat(1e9))
	return f.Text('f', 9)
}
