package service

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"
	"time"
	"wx-mch-trans/internal/model"
	"wx-mch-trans/internal/utils"
)

type WxPayService struct {
	config model.WxPayConfig
}

// NewWxPayService 创建微信支付服务
func NewWxPayService(config model.WxPayConfig) *WxPayService {
	return &WxPayService{
		config: config,
	}
}

// 生成签名
func (s *WxPayService) sign(method, path, body string, nonceStr string, timestamp int64) (string, error) {
	message := fmt.Sprintf("%s\n%s\n%d\n%s\n%s\n",
		method,
		path,
		timestamp,
		nonceStr,
		body,
	)

	// 解析私钥
	block, _ := pem.Decode([]byte(s.config.PrivateKey))
	if block == nil {
		return "", fmt.Errorf("failed to parse private key")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %v", err)
	}

	h := sha256.New()
	h.Write([]byte(message))
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey.(*rsa.PrivateKey), crypto.SHA256, h.Sum(nil))
	if err != nil {
		return "", fmt.Errorf("failed to sign message: %v", err)
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

// 发起商家转账
func (s *WxPayService) TransferBatches(req model.TransferBatchesRequest) (*http.Response, error) {
	url := "https://api.mch.weixin.qq.com/v3/transfer/batches"
	method := "POST"

	// 生成请求体
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request body failed: %v", err)
	}

	// 生成随机字符串
	nonceStr := utils.GenerateNonceStr()
	timestamp := time.Now().Unix()

	// 生成签名
	signature, err := s.sign(method, "/v3/transfer/batches", string(body), nonceStr, timestamp)
	if err != nil {
		return nil, fmt.Errorf("generate signature failed: %v", err)
	}

	// 构造认证信息
	token := fmt.Sprintf("WECHATPAY2-SHA256-RSA2048 mchid=\"%s\",nonce_str=\"%s\",signature=\"%s\",timestamp=\"%d\",serial_no=\"%s\"",
		s.config.MchID,
		nonceStr,
		signature,
		timestamp,
		s.config.SerialNo,
	)

	// 创建HTTP请求
	httpReq, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("create request failed: %v", err)
	}

	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Authorization", token)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send request failed: %v", err)
	}

	utils.InfoLogger.Printf("发起转账请求: %s\n响应状态码: %d", url, resp.StatusCode)
	return resp, nil
}
