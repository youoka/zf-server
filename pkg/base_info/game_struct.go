package base_info

type PullOrderReq struct {
	Amount float64 `json:"amount" binding:"required"`
}
