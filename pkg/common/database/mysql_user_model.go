package database

import (
	"zf-server/pkg/common/database/mysql_model_struct"
)

func (m *mysqlDB) GetUserCount(userId string) (int64, error) {
	var count int64
	return count, m.DefaultGormDB().Table(mysql_model_struct.User{}.TableName()).Where("user_id = ?", userId).Count(&count).Error
}
func (m *mysqlDB) GetUserByAccount(userId string) (mysql_model_struct.User, error) {
	var user mysql_model_struct.User
	return user, m.DefaultGormDB().Table(user.TableName()).Where("user_id = ?", userId).First(&user).Error
}
func (m *mysqlDB) UserRegister(user mysql_model_struct.User) error {
	err := m.DefaultGormDB().Table(user.TableName()).Create(&user).Error
	if err != nil {
		return err
	}
	return nil
}
func (m *mysqlDB) UpdateUser(user mysql_model_struct.User) error {
	err := m.DefaultGormDB().Table(user.TableName()).Where("user_id = ?", user.UserID).Updates(&user).Error
	if err != nil {
		return err
	}
	return nil
}

func (m *mysqlDB) UpdateBalance(userID string, amount float64) error {
	err := m.DefaultGormDB().Table(mysql_model_struct.User{}.TableName()).Where("user_id = ?", userID).Update("balance", amount).Error
	if err != nil {
		return err
	}
	return nil
}
