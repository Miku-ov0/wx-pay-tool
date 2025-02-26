package service

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"wx-mch-trans/internal/model"
	"wx-mch-trans/internal/utils"
)

type StorageService struct {
	baseDir string // 基础存储目录
}

// NewStorageService 创建存储服务
func NewStorageService() *StorageService {
	return &StorageService{
		baseDir: "data/trans",
	}
}

// SaveTransferRecords 保存转账记录到JSON文件
func (s *StorageService) SaveTransferRecords(records []model.TransferRecord) error {
	// 确保存储目录存在
	if err := os.MkdirAll(s.baseDir, 0755); err != nil {
		return fmt.Errorf("创建存储目录失败: %v", err)
	}

	// 生成文件名
	timestamp := time.Now().Format("20060102150405")
	filename := fmt.Sprintf("trans-%s.json", timestamp)
	filepath := filepath.Join(s.baseDir, filename)

	// 转换为JSON
	data, err := json.MarshalIndent(records, "", "    ")
	if err != nil {
		return fmt.Errorf("转换JSON失败: %v", err)
	}

	// 写入文件
	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return fmt.Errorf("写入文件失败: %v", err)
	}

	utils.InfoLogger.Printf("成功保存%d条转账记录到文件: %s", len(records), filepath)
	return nil
}

// LoadTransferRecords 从JSON文件加载转账记录
func (s *StorageService) LoadTransferRecords(filename string) ([]model.TransferRecord, error) {
	filepath := filepath.Join(s.baseDir, filename)

	// 读取文件
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %v", err)
	}

	// 解析JSON
	var records []model.TransferRecord
	if err := json.Unmarshal(data, &records); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %v", err)
	}

	utils.InfoLogger.Printf("成功从文件加载%d条转账记录: %s", len(records), filepath)
	return records, nil
}

// ListTransferFiles 列出所有转账记录文件
func (s *StorageService) ListTransferFiles() ([]string, error) {
	pattern := filepath.Join(s.baseDir, "trans-*.json")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("列出文件失败: %v", err)
	}

	// 只返回文件名，不包含路径
	var filenames []string
	for _, file := range files {
		filenames = append(filenames, filepath.Base(file))
	}

	return filenames, nil
}
