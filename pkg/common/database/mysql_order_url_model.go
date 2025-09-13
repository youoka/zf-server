package database

import (
	"time"
	"zf-server/pkg/common/database/mysql_model_struct"
)

func (m *mysqlDB) InstallOrderURL(OrderURL *mysql_model_struct.OrderURL) error {
	return m.DefaultGormDB().Table(OrderURL.TableName()).Create(OrderURL).Error
}

// GetOrderURLByAmountAndStatus 根据金额获取一个有效的订单URL（未过期且状态为0的记录）
func (m *mysqlDB) GetOrderURLByAmountAndStatus(amount float64) (mysql_model_struct.OrderURL, error) {
	var orderURL mysql_model_struct.OrderURL
	// 查询金额匹配、状态为0（未使用）、且过期时间大于当前时间的记录，只取第一条
	err := m.DefaultGormDB().Table(mysql_model_struct.OrderURL{}.TableName()).
		Where("amount = ? AND status = ? AND expire_time > ?", amount, 0, time.Now()).
		First(&orderURL).Error
	return orderURL, err
}
