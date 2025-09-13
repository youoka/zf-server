package user

import (
	"github.com/YuanJey/go-log/pkg/log"
	"github.com/YuanJey/goutils2/pkg/utils"
	"github.com/gin-gonic/gin"
	"zf-server/internal/api/common"
	"zf-server/pkg/base_info"
	"zf-server/pkg/common/database"
)

func Recharge(c *gin.Context) {
	operationID := utils.OperationIDGenerator()
	req := base_info.UserRechargeReq{}
	if err := c.BindJSON(&req); err != nil {
		common.ApiErr(c, operationID, 400, "参数错误"+err.Error())
		return
	}
	if req.Amount <= 0 {
		common.ApiErr(c, operationID, 400, "充值金额必须大于0")
		return
	}
	account, err2 := database.DB.MysqlDB.GetUserByAccount(req.UserID)
	if err2 != nil {
		common.ApiErr(c, operationID, 500, "查询用户失败"+err2.Error())
		return
	}
	err := database.DB.MysqlDB.UpdateBalance(req.UserID, account.Balance+req.Amount)
	if err != nil {
		common.ApiErr(c, operationID, 500, "充值失败"+err.Error())
		return
	}
	err = database.DB.MysqlDB.InstallBalanceLog(req.UserID, req.Amount)
	if err != nil {
		account2, err3 := database.DB.MysqlDB.GetUserByAccount(req.UserID)
		if err3 != nil {
			common.ApiErr(c, operationID, 500, "充值失败"+err.Error())
			return
		}
		if account2.Balance != account.Balance {
			err := database.DB.MysqlDB.UpdateBalance(req.UserID, account.Balance)
			if err != nil {
				log.Error(operationID, "修复余额失败 ", err.Error())
			}
		}
		log.Error(operationID, "充值失败 ", err.Error())
		common.ApiErr(c, operationID, 500, "充值失败"+err.Error())
		return
	}
	common.ApiSuccess(c, operationID, "充值成功", nil)
}
