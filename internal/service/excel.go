package service

import (
	"fmt"
	"strconv"

	"wx-mch-trans/internal/model"
	"wx-mch-trans/internal/utils"

	"github.com/xuri/excelize/v2"
)

type ExcelService struct{}

// NewExcelService 创建Excel服务
func NewExcelService() *ExcelService {
	return &ExcelService{}
}

// ParseTransferRecords 解析Excel中的转账记录
func (s *ExcelService) ParseTransferRecords(filePath string) ([]model.TransferRecord, error) {
	utils.InfoLogger.Printf("开始解析Excel文件: %s", filePath)

	// 打开Excel文件
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("打开Excel文件失败: %v", err)
	}
	defer f.Close()

	// 获取第一个工作表
	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("读取工作表失败: %v", err)
	}

	if len(rows) < 2 { // 至少需要标题行和一行数据
		return nil, fmt.Errorf("Excel文件为空或格式不正确")
	}

	// 验证表头
	expectedHeaders := []string{"商户单号", "收款用户OpenID", "转账金额"}
	headers := rows[0]
	if err := validateHeaders(headers, expectedHeaders); err != nil {
		return nil, err
	}

	// 解析数据行
	var records []model.TransferRecord
	for i, row := range rows[1:] { // 跳过表头
		if len(row) < 3 {
			return nil, fmt.Errorf("第%d行数据不完整", i+2)
		}

		// 解析金额
		amount, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			return nil, fmt.Errorf("第%d行转账金额格式不正确: %s", i+2, row[2])
		}

		// 验证金额
		if amount <= 0 {
			return nil, fmt.Errorf("第%d行转账金额必须大于0: %f", i+2, amount)
		}

		// 验证商户单号
		if row[0] == "" {
			return nil, fmt.Errorf("第%d行商户单号不能为空", i+2)
		}

		// 验证OpenID
		if row[1] == "" {
			return nil, fmt.Errorf("第%d行收款用户OpenID不能为空", i+2)
		}

		record := model.TransferRecord{
			OutBatchNo: row[0],
			OpenID:     row[1],
			Amount:     amount,
			Status:     "pending", // 初始状态为待处理
		}

		records = append(records, record)
	}

	utils.InfoLogger.Printf("成功解析%d条转账记录", len(records))
	return records, nil
}

// validateHeaders 验证表头是否符合要求
func validateHeaders(actual, expected []string) error {
	if len(actual) < len(expected) {
		return fmt.Errorf("表头列数不足，需要%d列，实际%d列", len(expected), len(actual))
	}

	for i, header := range expected {
		if actual[i] != header {
			return fmt.Errorf("表头格式不正确，第%d列应为'%s'，实际为'%s'", i+1, header, actual[i])
		}
	}

	return nil
}
