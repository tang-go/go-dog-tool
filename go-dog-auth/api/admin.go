package api

import (
	"omo-service/define"
	adminTable "omo-service/services/admin/table"
	"omo-service/services/auth/param"

	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/plugins"
)

//CreateAdminToken 创建管理员token
func CreateAdminToken(ctx plugins.Context, svcName string, admin *adminTable.Admin) (string, error) {
	info := param.CreateAdminTokenRes{}
	if err := ctx.GetClient().Call(ctx, plugins.RandomMode, define.SvcAuth, "CreateAdminToken", param.CreateAdminTokenReq{
		SvcName: svcName,
		Admin:   admin,
	}, &info); err != nil {
		log.Errorln(err.Error())
		return info.Token, err
	}
	return info.Token, nil
}

//UpdateAdminToken 更新管理员token
func UpdateAdminToken(ctx plugins.Context, token string, admin *adminTable.Admin) (bool, error) {
	info := param.UpdateAdminTokenRes{}
	if err := ctx.GetClient().Call(ctx, plugins.RandomMode, define.SvcAuth, "UpdateAdminToken", param.UpdateAdminTokenReq{
		Token: token,
		Admin: admin,
	}, &info); err != nil {
		log.Errorln(err.Error())
		return info.Success, err
	}
	return info.Success, nil
}

//CheckAdminToken 验证管理员token
func CheckAdminToken(ctx plugins.Context, token, svcName, method string) (*adminTable.Admin, error) {
	info := param.CheckAdminTokenRes{}
	if err := ctx.GetClient().Call(ctx, plugins.RandomMode, define.SvcAuth, "CheckAdminToken", &param.CheckAdminTokenReq{Token: token, SvcName: svcName, Method: method}, &info); err != nil {
		log.Errorln(err.Error())
		return info.Admin, err
	}
	return info.Admin, nil
}
