package user

import (
	"github.com/YuanJey/goutils2/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"zf-server/internal/api/common"
	"zf-server/internal/api/middleware"
	"zf-server/pkg/base_info"
	"zf-server/pkg/common/auth"
	"zf-server/pkg/common/database"
	"zf-server/pkg/common/database/mysql_model_struct"
)

func Register(c *gin.Context) {
	operationID := utils.OperationIDGenerator()
	req := base_info.UserRegisterReq{}
	if err := c.BindJSON(&req); err != nil {
		common.ApiErr(c, operationID, 400, "参数错误"+err.Error())
		return
	}
	count, _ := database.DB.MysqlDB.GetUserCount(req.UserID)
	if count != 0 {
		common.ApiErr(c, operationID, 400, "用户已存在")
		return
	}
	user := mysql_model_struct.User{}
	utils.CopyStructFields(&user, &req)
	err := database.DB.MysqlDB.UserRegister(user)
	if err != nil {
		common.ApiErr(c, operationID, 500, "注册失败"+err.Error())
		return
	}
	common.ApiSuccess(c, operationID, "注册成功", nil)
}

// Login 用户登录
func Login(c *gin.Context) {
	operationID := utils.OperationIDGenerator()
	req := base_info.UserLoginReq{}

	if err := c.BindJSON(&req); err != nil {
		common.ApiErr(c, operationID, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	// 查询用户
	user, err := database.DB.MysqlDB.GetUserByAccount(req.UserID)
	if err != nil {
		common.ApiErr(c, operationID, http.StatusUnauthorized, "用户名或密码错误")
		return
	}

	// 验证密码（这里应该使用加密后的密码比较）
	if user.Password != req.Password {
		common.ApiErr(c, operationID, http.StatusUnauthorized, "用户名或密码错误")
		return
	}

	// 生成JWT token
	token, err := auth.GenerateToken(req.UserID)
	if err != nil {
		common.ApiErr(c, operationID, http.StatusInternalServerError, "生成token失败: "+err.Error())
		return
	}
	account, err := database.DB.MysqlDB.GetUserByAccount(req.UserID)
	if err != nil {
		common.ApiErr(c, operationID, 500, "查询用户失败"+err.Error())
		return
	}
	resp := base_info.UserLoginResp{
		UserInfo: account,
		Token:    token,
		ExpireAt: 24 * 60 * 60, // 24小时，单位秒
	}

	common.ApiSuccess(c, operationID, "登录成功", resp)
}

// UserInfo 获取用户信息
func UserInfo(c *gin.Context) {
	operationID := utils.OperationIDGenerator()

	// 从上下文中获取用户ID
	userID := middleware.GetUserId(c)
	if userID == "" {
		common.ApiErr(c, operationID, http.StatusUnauthorized, "用户未登录")
		return
	}

	// 查询用户信息
	user, err := database.DB.MysqlDB.GetUserByAccount(userID)
	if err != nil {
		common.ApiErr(c, operationID, http.StatusInternalServerError, "查询用户信息失败: "+err.Error())
		return
	}

	common.ApiSuccess(c, operationID, "获取用户信息成功", user)
}
