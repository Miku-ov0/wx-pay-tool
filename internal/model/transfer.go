package model

// TransferRequest 转账请求的配置信息
type TransferRequest struct {
	// 商户配置
	MchID  string `json:"mchId" binding:"required"`
	AppID  string `json:"appId" binding:"required"`
	APIKey string `json:"apiKey" binding:"required"`

	// 转账配置
	SceneID   string `json:"sceneId" binding:"required"`
	Remark    string `json:"remark" binding:"required"`
	SceneInfo string `json:"sceneInfo" binding:"required"`
}

// TransferRecord Excel中的转账记录
type TransferRecord struct {
	// Excel中的字段
	OutBatchNo string  `json:"out_batch_no"` // 商户单号
	OpenID     string  `json:"open_id"`      // 收款用户OpenID
	Amount     float64 `json:"amount"`       // 转账金额

	// 处理后添加的字段
	TransferID   string `json:"transfer_id"`   // 微信转账单号
	SceneID      string `json:"scene_id"`      // 转账场景ID
	Remark       string `json:"remark"`        // 转账备注
	SceneInfo    string `json:"scene_info"`    // 转账场景报备信息
	Status       string `json:"status"`        // 转账状态
	RequestData  string `json:"request_data"`  // HTTP请求内容
	ResponseData string `json:"response_data"` // HTTP响应内容
}
