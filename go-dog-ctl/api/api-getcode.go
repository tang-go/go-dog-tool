package api

import (
	"github.com/mojocn/base64Captcha"
	"github.com/tang-go/go-dog-tool/define"
	"github.com/tang-go/go-dog-tool/go-dog-ctl/param"
	customerror "github.com/tang-go/go-dog/error"
	"github.com/tang-go/go-dog/lib/rand"
	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/plugins"
)

//GetCode 验证图片验证码
func (pointer *API) GetCode(ctx plugins.Context, request param.GetCodeReq) (response param.GetCodeRes, err error) {
	//查看是否是测试接口
	if ctx.GetIsTest() {
		response.ID = request.Code
		response.Img = "123456"
		pointer.Set(response.ID, response.Img)
		return
	}
	number := rand.StringRand(6)
	d := base64Captcha.NewDriverString(80, 240, 80, base64Captcha.OptionShowHollowLine, 6, number, nil, []string{})
	driver := d.ConvertFonts()
	code := base64Captcha.NewCaptcha(driver, pointer)
	id, b64s, err := code.Generate()
	if err != nil {
		log.Errorln(err.Error())
		err = customerror.EnCodeError(define.GetCodeErr, err.Error())
		return
	}
	response.ID = id
	response.Img = b64s
	return
}
