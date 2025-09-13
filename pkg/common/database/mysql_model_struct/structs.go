package mysql_model_struct

import "time"

type User struct {
	UserID      string    `gorm:"column:user_id;primary_key;type:varchar(64)" json:"userID"`
	Password    string    `gorm:"column:password;type:varchar(255)" json:"password"`
	ZFBAccount  string    `gorm:"column:zfb_account;type:varchar(64);uniqueIndex" json:"zfbAccount"`
	PhoneNumber string    `gorm:"column:phone_number;type:varchar(32);uniqueIndex" json:"phoneNumber"`
	CreateTime  time.Time `gorm:"column:create_time;autoCreateTime" json:"createTime"`
	UpdateTime  time.Time `gorm:"column:update_time;autoUpdateTime" json:"updateTime"`
	Balance     float64   `gorm:"column:balance;type:REAL;default:0" json:"balance"`
	Status      int       `gorm:"column:status;type:int;default:0" json:"status"` // 0 normal, 1 disabled
}

func (u User) TableName() string {
	return "users"
}

type BalanceLog struct {
	UserID     string    `gorm:"column:user_id;primary_key;type:varchar(64)" json:"userID"`
	Amount     float64   `gorm:"column:amount;type:REAL" json:"amount"`
	CreateTime time.Time `gorm:"column:create_time;primary_key;autoCreateTime" json:"createTime"`
}

func (b BalanceLog) TableName() string {
	return "balance_logs"
}

type OrderURL struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     string    `gorm:"column:user_id;type:varchar(64)" json:"userID"`
	BizNo      string    `gorm:"column:biz_no;type:varchar(255);uniqueIndex" json:"biz_no"`
	Amount     float64   `gorm:"column:amount;type:REAL" json:"amount"`
	URL        string    `gorm:"column:url;type:varchar(2000)" json:"url"` //uniqueIndex
	Status     int       `gorm:"column:status;type:int;default:0" json:"status"`
	ExpireTime time.Time `gorm:"column:expire_time;type:datetime" json:"expireTime"`
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime" json:"createTime"`
	UpdateTime time.Time `gorm:"column:update_time;autoUpdateTime" json:"updateTime"`
}

func (o OrderURL) TableName() string {
	return "order_urls"
}
