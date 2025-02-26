package handler

import (
	"net/http"

	"wx-mch-trans/internal/service"
	"wx-mch-trans/internal/utils"

	"github.com/gin-gonic/gin"
)

type RecordsHandler struct {
	storageService *service.StorageService
}

func NewRecordsHandler() *RecordsHandler {
	return &RecordsHandler{
		storageService: service.NewStorageService(),
	}
}

// ListRecordFiles 列出所有转账记录文件
func (h *RecordsHandler) ListRecordFiles(c *gin.Context) {
	files, err := h.storageService.ListTransferFiles()
	if err != nil {
		utils.ErrorLogger.Printf("列出转账记录文件失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取转账记录列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"files": files,
	})
}

// GetRecordFile 获取指定的转账记录文件内容
func (h *RecordsHandler) GetRecordFile(c *gin.Context) {
	filename := c.Param("filename")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "未指定文件名",
		})
		return
	}

	records, err := h.storageService.LoadTransferRecords(filename)
	if err != nil {
		utils.ErrorLogger.Printf("读取转账记录文件失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "读取转账记录失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"records": records,
	})
}
