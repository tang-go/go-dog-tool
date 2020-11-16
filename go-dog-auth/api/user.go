package api

import (
	"omo-service/define"
	"omo-service/services/auth/param"
	userTable "omo-service/services/user/dao/table"

	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/plugins"
)

//CreateUserToken 创建用户token
func CreateUserToken(ctx plugins.Context, svcName string, user *userTable.User) (string, error) {
	info := param.CreateUserTokenRes{}
	if err := ctx.GetClient().Call(ctx, plugins.RandomMode, define.SvcAuth, "CreateUserToken", param.CreateUserTokenReq{
		SvcName: svcName,
		User:    user,
	}, &info); err != nil {
		log.Errorln(err.Error())
		return info.Token, err
	}
	return info.Token, nil
}

//UpdateUserToken 更新用户token
func UpdateUserToken(ctx plugins.Context, token string, user *userTable.User) (bool, error) {
	info := param.UpdateUserTokenRes{}
	if err := ctx.GetClient().Call(ctx, plugins.RandomMode, define.SvcAuth, "UpdateUserToken", param.UpdateUserTokenReq{
		Token: token,
		User:  user,
	}, &info); err != nil {
		log.Errorln(err.Error())
		return info.Success, err
	}
	return info.Success, nil
}

//CheckUserToken 验证用户token
func CheckUserToken(ctx plugins.Context, token, svcName, method string) (*userTable.User, error) {
	info := param.CheckUserTokenRes{}
	if err := ctx.GetClient().Call(ctx, plugins.RandomMode, define.SvcAuth, "CheckUserToken", &param.CheckUserTokenReq{Token: token, SvcName: svcName, Method: method}, &info); err != nil {
		log.Errorln(err.Error())
		return info.User, err
	}
	return info.User, nil
}
