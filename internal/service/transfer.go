package service

import (
	"fmt"
	"time"

	"wx-mch-trans/internal/model"
	"wx-mch-trans/internal/utils"
)

type TransferService struct {
	wxPayService   *WxPayService
	storageService *StorageService
}

// NewTransferService 创建转账服务
func NewTransferService(config model.WxPayConfig) *TransferService {
	return &TransferService{
		wxPayService:   NewWxPayService(config),
		storageService: NewStorageService(),
	}
}

// ProcessTransfer 处理转账
func (s *TransferService) ProcessTransfer(records []model.TransferRecord, config model.TransferRequest) error {
	utils.InfoLogger.Printf("开始处理%d条转账记录", len(records))

	// 遍历处理每条记录
	for i := range records {
		// 生成批次号
		timestamp := time.Now().Format("20060102150405")
		batchNo := fmt.Sprintf("B%s%03d", timestamp, i+1)

		// 构建转账请求
		transferReq := model.TransferBatchesRequest{
			AppID:       config.AppID,
			OutBatchNo:  batchNo,
			BatchName:   fmt.Sprintf("批量转账-%s", batchNo),
			BatchRemark: config.Remark,
			TotalAmount: utils.ConvertToWxAmount(records[i].Amount),
			TotalNum:    1,
			TransferList: []model.TransferDetailList{
				{
					OutDetailNo:    records[i].OutBatchNo,
					TransferAmount: utils.ConvertToWxAmount(records[i].Amount),
					OpenID:         records[i].OpenID,
					UserName:       "", // 微信不再要求上送用户姓名
					Remark:         config.Remark,
				},
			},
			TransferScene: model.TransferScene{
				SceneID:   config.SceneID,
				SceneInfo: config.SceneInfo,
			},
		}

		// 调用微信转账接口
		resp, err := s.wxPayService.TransferBatches(transferReq)
		if err != nil {
			utils.ErrorLogger.Printf("转账失败[%s]: %v", records[i].OutBatchNo, err)
			records[i].Status = "failed"
			records[i].ResponseData = err.Error()
			continue
		}
		defer resp.Body.Close()

		// 更新记录状态
		records[i].Status = fmt.Sprintf("success_%d", resp.StatusCode)
		records[i].TransferID = batchNo
		records[i].SceneID = config.SceneID
		records[i].Remark = config.Remark
		records[i].SceneInfo = config.SceneInfo

		utils.InfoLogger.Printf("转账成功[%s]: 批次号=%s", records[i].OutBatchNo, batchNo)
	}

	// 保存更新后的记录
	if err := s.storageService.SaveTransferRecords(records); err != nil {
		return fmt.Errorf("保存转账记录失败: %v", err)
	}

	return nil
}
