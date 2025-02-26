package handler

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

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

// HandleBatchTransfer 处理批量转账
func (h *TransferHandler) HandleBatchTransfer(c *gin.Context) {
	utils.InfoLogger.Println("接收到批量转账请求")

	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		utils.ErrorLogger.Printf("获取上传文件失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请选择要上传的Excel文件",
		})
		return
	}

	// 验证文件类型
	ext := filepath.Ext(file.Filename)
	if ext != ".xlsx" && ext != ".xls" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "只支持.xlsx或.xls格式的Excel文件",
		})
		return
	}

	// 生成临时文件名
	timestamp := time.Now().Format("20060102150405")
	filename := fmt.Sprintf("upload_%s%s", timestamp, ext)
	tempPath := filepath.Join("data", "temp", filename)

	// 保存文件
	if err := c.SaveUploadedFile(file, tempPath); err != nil {
		utils.ErrorLogger.Printf("保存上传文件失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "文件保存失败",
		})
		return
	}

	// 解析Excel文件
	excelService := service.NewExcelService()
	records, err := excelService.ParseTransferRecords(tempPath)
	if err != nil {
		utils.ErrorLogger.Printf("解析Excel文件失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 获取配置信息
	var config model.TransferRequest
	config.MchID = c.PostForm("mchId")
	config.AppID = c.PostForm("appId")
	config.APIKey = c.PostForm("apiKey")
	config.SceneID = c.PostForm("sceneId")
	config.Remark = c.PostForm("remark")
	config.SceneInfo = c.PostForm("sceneInfo")
	// 验证必要字段
	if config.MchID == "" || config.AppID == "" || config.APIKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "转账配置信息不完整：需要提供商户号、应用ID和API密钥",
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

	// 返回处理结果
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("成功处理%d条转账记录", len(records)),
		"records": records,
	})
}
