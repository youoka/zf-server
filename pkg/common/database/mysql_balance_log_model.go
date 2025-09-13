package database

import "zf-server/pkg/common/database/mysql_model_struct"

func (m *mysqlDB) InstallBalanceLog(userId string, amount float64) error {
	balanceLog := mysql_model_struct.BalanceLog{
		UserID: userId,
		Amount: amount,
	}
	err := m.DefaultGormDB().Table(balanceLog.TableName()).Create(&balanceLog).Error
	if err != nil {
		return err
	}
	return nil
}
