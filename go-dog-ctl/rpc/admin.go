package rpc

import (
	"github.com/tang-go/go-dog-tool/define"
	"github.com/tang-go/go-dog-tool/go-dog-ctl/param"
	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/plugins"
)

//AdminOnline 管理员上线
func AdminOnline(client plugins.Client, ctx plugins.Context, address string) (bool, error) {
	res := new(param.AdminOnlineRes)
	if err := client.Call(ctx, plugins.RandomMode, define.SvcController, "AdminOnline", &param.AdminOnlineReq{Address: address}, res); err != nil {
		log.Errorln(err.Error())
		return res.Success, err
	}
	return res.Success, nil
}

//AdminOffline 管理员下线
func AdminOffline(client plugins.Client, ctx plugins.Context, address string) (bool, error) {
	res := new(param.AdminOfflineRes)
	if err := client.Call(ctx, plugins.RandomMode, define.SvcController, "AdminOffline", &param.AdminOfflineReq{Address: address}, res); err != nil {
		log.Errorln(err.Error())
		return res.Success, err
	}
	return res.Success, nil
}

//AuthAdmin 验证管理员
func AuthAdmin(client plugins.Client, ctx plugins.Context, token string) (bool, error) {
	res := new(param.AuthAdminRes)
	if err := client.Call(ctx, plugins.RandomMode, define.SvcController, "AuthAdmin", &param.AuthAdminReq{Token: token}, res); err != nil {
		log.Errorln(err.Error())
		return res.Success, err
	}
	return res.Success, nil
}
