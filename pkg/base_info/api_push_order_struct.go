package base_info

type PushOrderReq struct {
	Amount float64 `json:"amount" binding:"required"`
	URL    string  `json:"url" binding:"required"`
}
