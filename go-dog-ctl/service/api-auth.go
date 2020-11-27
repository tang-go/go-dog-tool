package service

import (
	"fmt"

	"github.com/tang-go/go-dog-tool/define"
	authRPC "github.com/tang-go/go-dog-tool/go-dog-auth/rpc"
	"github.com/tang-go/go-dog-tool/go-dog-ctl/table"
	customerror "github.com/tang-go/go-dog/error"
	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/plugins"
)

//Auth 插件
func (s *Service) Auth(ctx plugins.Context, method, token string) error {
	admin := new(table.Admin)
	if e := s.cache.GetCache().Get(token, admin); e != nil {
		return customerror.EnCodeError(define.AdminTokenErr, "token失效或者不正确")
	}
	//获取权限
	adminapis := make(map[string]string)
	if e := s.cache.GetCache().Get(fmt.Sprintf("%s-%d", define.Organize, admin.RoleID), &adminapis); e != nil {
		//验证接口权限
		apis, e := authRPC.GetRoleAPI(ctx, define.Organize, admin.RoleID)
		if e != nil {
			log.Errorln(e.Error())
			return customerror.EnCodeError(define.APIPowerErr, "api权限不正确")
		}
		for _, api := range apis {
			adminapis[api.API] = api.API
		}
		if e := s.cache.GetCache().SetByTime(fmt.Sprintf("%s-%d", define.Organize, admin.RoleID), adminapis, 60*10); e != nil {
			log.Errorln(e.Error())
			return customerror.EnCodeError(define.APIPowerErr, "api权限不正确")
		}
	}
	url := ctx.GetURL()
	if _, ok := adminapis[url[1:]]; !ok {
		log.Warnln("请求url：", url[1:])
		return customerror.EnCodeError(define.APIPowerErr, "api权限不正确")
	}
	ctx.SetShare("Admin", admin)
	return nil
}
