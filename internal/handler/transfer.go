package handler

import (
	"fmt"
	"net/http"

	"wx-mch-trans/internal/model"
	"wx-mch-trans/internal/service"
	"wx-mch-trans/internal/utils"

	"github.com/gin-gonic/gin"
)

type TransferHandler struct {
	transferService *service.TransferService
	storageService  *service.StorageService
}

func NewTransferHandler() *TransferHandler {
	return &TransferHandler{
		storageService: service.NewStorageService(),
	}
}

// HandleTransfer 处理转账请求
func (h *TransferHandler) HandleTransfer(c *gin.Context) {
	// 获取文件名
	filename := c.PostForm("filename")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "未指定转账记录文件",
		})
		return
	}

	// 获取配置信息
	var config model.TransferRequest
	if err := c.ShouldBind(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "转账配置信息不完整",
		})
		return
	}

	// 加载转账记录
	records, err := h.storageService.LoadTransferRecords(filename)
	if err != nil {
		utils.ErrorLogger.Printf("加载转账记录失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "加载转账记录失败",
		})
		return
	}

	// 创建转账服务
	wxConfig := model.WxPayConfig{
		MchID:      config.MchID,
		AppID:      config.AppID,
		APIv3Key:   config.APIKey,
		PrivateKey: "", // 需要从配置或其他安全位置获取
		SerialNo:   "", // 需要从配置或其他安全位置获取
	}
	transferService := service.NewTransferService(wxConfig)

	// 处理转账
	if err := transferService.ProcessTransfer(records, config); err != nil {
		utils.ErrorLogger.Printf("处理转账失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("处理转账失败: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "转账处理完成",
		"records": records,
	})
}
