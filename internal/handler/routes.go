package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册所有路由
func RegisterRoutes(r *gin.Engine) {
	// uploadHandler := NewUploadHandler()
	recordsHandler := NewRecordsHandler()
	transferHandler := NewTransferHandler()

	// 转账处理路由
	r.POST("/api/transfer/batch", transferHandler.HandleBatchTransfer)

	// 转账记录相关路由
	r.GET("/api/records", recordsHandler.ListRecordFiles)
	r.GET("/api/records/:filename", recordsHandler.GetRecordFile)

	// 添加首页路由
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "upload.html", nil)
	})
}
