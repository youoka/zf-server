package order

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/YuanJey/goutils2/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/liyue201/goqr"
	"github.com/tuotoo/qrcode"
	"zf-server/internal/api/common"
	"zf-server/internal/api/middleware"
	"zf-server/pkg/common/database"
	"zf-server/pkg/common/database/mysql_model_struct"
)

// PushOrder 处理推送订单请求
func PushOrder(c *gin.Context) {
	operationID := utils.OperationIDGenerator()

	// 从表单中获取金额参数
	amountStr := c.PostForm("amount")
	if amountStr == "" {
		common.ApiErr(c, operationID, http.StatusBadRequest, "缺少金额参数")
		return
	}

	// 解析金额
	var amount float64
	_, err := fmt.Sscanf(amountStr, "%f", &amount)
	if err != nil || amount <= 0 {
		common.ApiErr(c, operationID, http.StatusBadRequest, "金额参数无效")
		return
	}

	// 从表单中获取二维码图片文件
	file, err := c.FormFile("qrcode")
	if err != nil {
		common.ApiErr(c, operationID, http.StatusBadRequest, "缺少二维码图片文件: "+err.Error())
		return
	}

	// 解析二维码图片
	qrContent, err := parseQRCode(file)
	if err != nil {
		common.ApiErr(c, operationID, http.StatusBadRequest, "二维码解析失败: "+err.Error())
		return
	}

	// 获取当前用户ID
	userID := middleware.GetUserId(c)

	// 创建订单URL记录
	orderURL := mysql_model_struct.OrderURL{
		UserID:     userID,
		Amount:     amount,
		URL:        qrContent,
		Status:     0,                                          // 0表示未使用
		ExpireTime: time.Now().Add(72 * 60 * 60 * time.Second), // 3天后过期
	}

	// 保存到数据库
	err = database.DB.MysqlDB.InstallOrderURL(&orderURL)
	if err != nil {
		common.ApiErr(c, operationID, http.StatusInternalServerError, "保存订单失败: "+err.Error())
		return
	}

	common.ApiSuccess(c, operationID, "订单推送成功", qrContent)
}

// parseQRCode 解析二维码图片，使用多种方法提高成功率
func parseQRCode(file *multipart.FileHeader) (string, error) {
	// 打开文件
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// 方法1: 使用tuotoo/qrcode解析
	qr, err := qrcode.Decode(src)
	if err == nil {
		return qr.Content, nil
	}

	// 重新打开文件用于其他方法
	src2, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("重新打开文件失败: %v", err)
	}
	defer src2.Close()

	// 解码图片
	img, _, err := image.Decode(src2)
	if err != nil {
		return "", fmt.Errorf("图片解码失败: %v", err)
	}

	// 预处理图片以提高识别率
	processedImg := preprocessImage(img)

	// 方法2: 使用goqr库解析预处理后的图片
	qrCodes, err := goqr.Recognize(processedImg)
	if err == nil && len(qrCodes) > 0 {
		return string(qrCodes[0].Payload), nil
	}

	// 如果所有方法都失败，返回错误
	return "", fmt.Errorf("无法识别二维码内容，可能是因为二维码中间有图标导致识别失败，请使用无图标的二维码")
}

// preprocessImage 预处理图片以提高二维码识别率
func preprocessImage(img image.Image) image.Image {
	// 转换为灰度图像
	bounds := img.Bounds()
	gray := image.NewGray(bounds)
	draw.Draw(gray, bounds, img, bounds.Min, draw.Src)

	// 应用简单的对比度增强
	enhanced := enhanceContrast(gray)
	return enhanced
}

// enhanceContrast 简单的对比度增强
func enhanceContrast(img *image.Gray) image.Image {
	bounds := img.Bounds()
	enhanced := image.NewRGBA(bounds)

	// 计算平均亮度
	var total uint32
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			total += uint32(img.GrayAt(x, y).Y)
		}
	}
	avg := uint8(total / uint32(bounds.Dx()*bounds.Dy()))

	// 应用对比度增强
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			pixel := img.GrayAt(x, y)
			var newPixel uint8
			if pixel.Y > avg {
				// 增强亮部
				if pixel.Y+50 > 255 {
					newPixel = 255
				} else {
					newPixel = pixel.Y + 50
				}
			} else {
				// 增强暗部
				if pixel.Y < 50 {
					newPixel = 0
				} else {
					newPixel = pixel.Y - 50
				}
			}
			enhanced.SetRGBA(x, y, color.RGBA{newPixel, newPixel, newPixel, 255})
		}
	}

	return enhanced
}
