package game

import (
	"github.com/YuanJey/goutils2/pkg/utils"
	"github.com/gin-gonic/gin"
	"zf-server/internal/api/common"
	"zf-server/pkg/base_info"
	"zf-server/pkg/common/database"
)

// PullOrder 拉取订单
func PullOrder(c *gin.Context) {
	operationID := utils.OperationIDGenerator()
	req := base_info.PullOrderReq{}
	if err := c.BindJSON(&req); err != nil {
		common.ApiErr(c, operationID, 400, "参数错误"+err.Error())
		return
	}

	orderURL, err := database.DB.MysqlDB.GetOrderURLByAmountAndStatus(req.Amount)
	if err != nil {
		common.ApiErr(c, operationID, 500, "获取订单失败: "+err.Error())
		return
	}

	// 如果没有找到可用订单，返回空
	if orderURL.URL == "" {
		common.ApiSuccess(c, operationID, "暂无可用订单", nil)
		return
	}

	// 更新订单状态为已使用
	//orderURL.Status = 1 // 标记为已使用
	//err = database.DB.MysqlDB.InstallOrderURL(&orderURL)
	//if err != nil {
	//	common.ApiErr(c, operationID, 500, "更新订单状态失败: "+err.Error())
	//	return
	//}

	common.ApiSuccess(c, operationID, "获取订单成功", orderURL)
}
