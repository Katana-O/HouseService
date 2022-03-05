package handler

import (
	"bj38web/service/getCaptcha/proto/getCaptcha"
	"context"
	"encoding/json"
	"fmt"
	"github.com/afocus/captcha"
	"github.com/micro/go-micro/util/log"
	"image/color"
	"bj38web/service/getCaptcha/model"
)

type GetCaptcha struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *GetCaptcha) Call(ctx context.Context, req *getCaptcha.Request, rsp *getCaptcha.Response) error {
	fmt.Println("MicroService Call request received")
	log.Log("Received GetCaptcha.Call request")
	cap := captcha.New()
	// 可以设置多个字体 或使用cap.AddFont("xx.ttf")追加
	cap.SetFont("./conf/comic.ttf")
	// 设置验证码大小
	cap.SetSize(128, 64)
	// 设置干扰强度
	cap.SetDisturbance(captcha.MEDIUM)
	// 设置前景色 可以多个 随机替换文字颜色 默认黑色
	cap.SetFrontColor(color.RGBA{255, 255, 255, 255})
	// 设置背景色 可以多个 随机替换背景色 默认白色
	cap.SetBkgColor(color.RGBA{255, 0, 0, 255}, color.RGBA{0, 0, 255, 255}, color.RGBA{0, 153, 0, 255})
	// Create Image
	img, imgCode := cap.Create(4,captcha.NUM)
	// save uuid & imgCode to redis DB.
	err := model.SaveImgCode(imgCode, req.Uuid)
	if err != nil {
		return err
	}
	imgBuf, _ := json.Marshal(img)
	rsp.Img = imgBuf
	return nil
}
