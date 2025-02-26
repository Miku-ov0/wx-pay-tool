package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册所有路由
func RegisterRoutes(r *gin.Engine) {
	uploadHandler := NewUploadHandler()
	recordsHandler := NewRecordsHandler()
	transferHandler := NewTransferHandler()

	// 文件上传相关路由
	r.POST("/api/upload", uploadHandler.HandleFileUpload)

	// 转账记录相关路由
	r.GET("/api/records", recordsHandler.ListRecordFiles)
	r.GET("/api/records/:filename", recordsHandler.GetRecordFile)

	// 转账处理路由
	r.POST("/api/transfer", transferHandler.HandleTransfer)

	// 添加首页路由
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "upload.html", nil)
	})
}
