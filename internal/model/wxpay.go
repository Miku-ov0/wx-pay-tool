package model

// WxPayConfig 微信支付配置
type WxPayConfig struct {
	MchID      string `json:"mch_id"`      // 商户号
	AppID      string `json:"app_id"`      // 应用ID
	APIv3Key   string `json:"apiv3_key"`   // API v3密钥
	PrivateKey string `json:"private_key"` // 商户私钥
	SerialNo   string `json:"serial_no"`   // 证书序列号
}

// TransferBatchesRequest 商家转账API请求参数
type TransferBatchesRequest struct {
	AppID         string               `json:"appid"`                // 应用ID
	OutBatchNo    string               `json:"out_batch_no"`         // 商户批次单号
	BatchName     string               `json:"batch_name"`           // 批次名称
	BatchRemark   string               `json:"batch_remark"`         // 批次备注
	TotalAmount   int64                `json:"total_amount"`         // 转账总金额(分)
	TotalNum      int                  `json:"total_num"`            // 转账总笔数
	TransferList  []TransferDetailList `json:"transfer_detail_list"` // 转账明细列表
	TransferScene TransferScene        `json:"transfer_scene"`       // 转账场景
}

// TransferDetailList 转账明细
type TransferDetailList struct {
	OutDetailNo    string `json:"out_detail_no"`   // 商户明细单号
	TransferAmount int64  `json:"transfer_amount"` // 转账金额(分)
	OpenID         string `json:"openid"`          // 用户openid
	UserName       string `json:"user_name"`       // 收款用户姓名
	Remark         string `json:"remark"`          // 转账备注
}

// TransferScene 转账场景
type TransferScene struct {
	SceneID   string `json:"scene_id"`   // 场景ID
	SceneInfo string `json:"scene_info"` // 场景信息
}
