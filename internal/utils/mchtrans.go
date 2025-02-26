package main

import (
	"wx-mch-trans/internal/utils"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	wxutils "github.com/wechatpay-apiv3/wechatpay-go/utils"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
)

// TransferSceneInfo 转账场景信息
type TransferSceneInfo struct {
	InfoType    string `json:"info_type"`
	InfoContent string `json:"info_content"`
}

// MchTransferRequest 商家转账请求参数
type MchTransferRequest struct {
	Appid                   string              `json:"appid"`
	OutBillNo               string              `json:"out_bill_no"`
	TransferSceneID         string              `json:"transfer_scene_id"`
	OpenID                  string              `json:"openid"`
	UserName                string              `json:"user_name"`
	TransferAmount          int64               `json:"transfer_amount"`
	TransferRemark          string              `json:"transfer_remark"`
	NotifyURL               string              `json:"notify_url"`
	UserRecvPerception      string              `json:"user_recv_perception"`
	TransferSceneReportInfos []TransferSceneInfo `json:"transfer_scene_report_infos"`
}

type TransferResp struct {
	OutBillNo       string `json:"out_bill_no"`
	TransferBillNo  string `json:"transfer_bill_no"`
	CreateTime      string `json:"create_time"`
	State           string `json:"state"`
	FailReason      string `json:"fail_reason"`
	PackageInfo     string `json:"package_info"`
}


func MchTransfer(req MchTransferRequest) (resp *TransferResp, result, err error) {
	
	var (
		appid string = "Appid"
		mchID string = "MchID"
		mchCertificateSerialNumber string = "MchCertificateSerialNumber"
		mchPrivateKey string = "MchPrivateKey"
		wechatpayPublicKeyID       string = "PUB_KEY_ID_00000000000000000000000000000000"          // 微信支付公钥ID，注意微信支付公钥ID带有PUB_KEY_ID_前缀
	)
	wechatpayPublicKey, err := wxutils.LoadPublicKeyWithPath("./config/pub_key.pem")
	if err != nil {
		panic(fmt.Errorf("load wechatpay public key err:%s", err.Error()))
	}
		
	// 初始化 Client
	opts := []core.ClientOption{
		option.WithWechatPayPublicKeyAuthCipher(
			mchID, //商户号
			mchCertificateSerialNumber, mchPrivateKey, 
			wechatpayPublicKeyID, wechatpayPublicKey),
	}
	client, err := core.NewClient(ctx, opts...)
	if err != nil {
		utils.InfoLogger.Printf("new wechat pay client err:%s", err)
	}
	req.Appid = appid

	utils.InfoLogger.Printf("商家转账请求:%s", req)
	result, err = client.Post(ctx, 
		"https://api.mch.weixin.qq.com/v3/fund-app/mch-transfer/transfer-bills",
		req)
		
	if err != nil {
		utils.InfoLogger.Printf("商家转账请求失败: %v", err)
	}

	utils.InfoLogger.Printf("商家转账请求成功: %v", result)

	resp := new(TransferResp)
	err = core.UnMarshalResponse(result.Response, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func main () {


	req := &MchTransferRequest{
		Appid:            "wxf636efh567hg4356",
		OutBillNo:        "plfk2020042013",
		TransferSceneID:  "1000",
		OpenID:           "o-MYE42l80oelYMDE34nYD456Xoy",
		TransferAmount:   1,
		TransferRemark:   "测试发红包",
		TransferSceneReportInfos: []TransferSceneInfo{
			{
				InfoType:    "测试发红包",
				InfoContent: "测试发红包",
			},
		},
	}

	MchTransfer(req)
}
