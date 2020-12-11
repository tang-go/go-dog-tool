package rpc

import (
	"github.com/tang-go/go-dog-tool/define"
	"github.com/tang-go/go-dog-tool/go-dog-ctl/param"
	"github.com/tang-go/go-dog-tool/go-dog-ctl/table"
	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/plugins"
)

//Auth 验证
func Auth(client plugins.Client, ctx plugins.Context, token, url string) (table.Admin, error) {
	res := new(param.AuthRes)
	if err := client.Call(ctx, plugins.RandomMode, define.SvcController, "Auth", &param.AuthReq{Token: token, URL: url}, res); err != nil {
		log.Errorln(err.Error())
		return res.Admin, err
	}
	return res.Admin, nil
}
