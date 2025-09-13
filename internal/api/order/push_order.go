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
	"net/url"
	"time"
	"zf-server/pkg/base_info"

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
	req := base_info.PushOrderReq{}
	if err := c.BindJSON(&req); err != nil {
		common.ApiErr(c, operationID, 400, "参数错误"+err.Error())
		return
	}
	//https://qr.alipay.com/_d?_b=peerpay&enableWK=YES&biz_no=2025091304200364641064636490_4adb46c73eb419ba5c4ce3ac222b7a48&app_name=tb&sc=qr_code&v=20250920&sign=c28a82&__webview_options__=pd%3dNO&channel=qr_code
	userId := middleware.GetUserId(c)
	urlParams, err2 := parseAlipayURLParams(req.URL)
	if err2 != nil {
		common.ApiErr(c, operationID, 401, "URL错误"+err2.Error())
		return
	}
	if _, ok := urlParams["biz_no"]; !ok {
		common.ApiErr(c, operationID, 402, "URL错误")
		return
	}
	// 创建订单URL记录
	orderURL := mysql_model_struct.OrderURL{
		UserID:     userId,
		Amount:     req.Amount, // Changed from amount to req.Amount
		URL:        req.URL,
		BizNo:      urlParams["biz_no"],
		Status:     0,                                          // 0表示未使用
		ExpireTime: time.Now().Add(72 * 60 * 60 * time.Second), // 3天后过期
	}

	// 保存到数据库
	err := database.DB.MysqlDB.InstallOrderURL(&orderURL)
	if err != nil {
		common.ApiErr(c, operationID, http.StatusInternalServerError, "保存订单失败: "+err.Error())
		return
	}

	common.ApiSuccess(c, operationID, "订单推送成功", req.URL)
}
func parseAlipayURLParams(urlString string) (map[string]string, error) {
	// 解析URL
	parsedURL, err := url.Parse(urlString)
	if err != nil {
		return nil, fmt.Errorf("解析URL失败: %v", err)
	}

	// 解析查询参数
	params := parsedURL.Query()

	// 转换为map[string]string
	result := make(map[string]string)
	for key, values := range params {
		if len(values) > 0 {
			result[key] = values[0] // 取第一个值
		}
	}

	return result, nil
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
