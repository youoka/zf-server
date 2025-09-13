package base_info

type UserRegisterReq struct {
	UserID      string `json:"userID" binding:"required"`
	Password    string `json:"password" binding:"required"`
	PhoneNumber string `json:"phoneNumber" binding:"required"`
	Email       string `json:"email" binding:"required"`
	ZFBAccount  string `json:"zfbAccount" binding:"required"`
}
type UserRechargeReq struct {
	UserID string  `json:"userID" binding:"required"`
	Amount float64 `json:"amount" binding:"required"`
}

type UserLoginReq struct {
	UserID   string `json:"userID" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserLoginResp struct {
	UserInfo any    `json:"userInfo"`
	Token    string `json:"token"`
	ExpireAt int64  `json:"expireAt"`
}
