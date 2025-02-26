package handler

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"wx-mch-trans/internal/service"
	"wx-mch-trans/internal/utils"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	excelService   *service.ExcelService
	storageService *service.StorageService
}

func NewUploadHandler() *UploadHandler {
	return &UploadHandler{
		excelService:   service.NewExcelService(),
		storageService: service.NewStorageService(),
	}
}

// HandleFileUpload 处理文件上传
func (h *UploadHandler) HandleFileUpload(c *gin.Context) {
	utils.InfoLogger.Println("接收到文件上传请求")

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
	records, err := h.excelService.ParseTransferRecords(tempPath)
	if err != nil {
		utils.ErrorLogger.Printf("解析Excel文件失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 保存转账记录
	if err := h.storageService.SaveTransferRecords(records); err != nil {
		utils.ErrorLogger.Printf("保存转账记录失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "保存转账记录失败",
		})
		return
	}

	// 返回解析结果
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("成功解析并保存%d条转账记录", len(records)),
		"records": records,
	})
}
