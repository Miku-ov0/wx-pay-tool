package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateNonceStr 生成随机字符串
func GenerateNonceStr() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// ConvertToWxAmount 将元转换为分（微信支付金额单位）
func ConvertToWxAmount(amount float64) int64 {
	return int64(amount * 100)
}
