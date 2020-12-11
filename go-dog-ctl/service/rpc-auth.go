package service

import (
	"fmt"

	"github.com/tang-go/go-dog-tool/define"
	authRPC "github.com/tang-go/go-dog-tool/go-dog-auth/rpc"
	"github.com/tang-go/go-dog-tool/go-dog-ctl/param"
	customerror "github.com/tang-go/go-dog/error"
	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/plugins"
)

//Auth 插件
func (s *Service) Auth(ctx plugins.Context, request param.AuthReq) (response param.AuthRes, err error) {
	token := request.Token
	url := request.URL
	if e := s.cache.GetCache().Get(token, &response.Admin); e != nil {
		log.Errorln(token, e.Error())
		err = customerror.EnCodeError(define.AdminTokenErr, "token失效或者不正确")
		return
	}
	//获取权限
	adminapis := make(map[string]string)
	if e := s.cache.GetCache().Get(fmt.Sprintf("%s-%d", define.Organize, response.Admin.RoleID), &adminapis); e != nil {
		//验证接口权限
		apis, e := authRPC.GetRoleAPI(ctx, define.Organize, response.Admin.RoleID)
		if e != nil {
			log.Errorln(e.Error())
			err = customerror.EnCodeError(define.APIPowerErr, "api权限不正确")
			return
		}
		for _, api := range apis {
			adminapis[api.API] = api.API
		}
		if e := s.cache.GetCache().SetByTime(fmt.Sprintf("%s-%d", define.Organize, response.Admin.RoleID), adminapis, 60*10); e != nil {
			log.Errorln(e.Error())
			err = customerror.EnCodeError(define.APIPowerErr, "api权限不正确")
			return
		}
	}
	if _, ok := adminapis[url[1:]]; !ok {
		log.Warnln("请求url：", url[1:])
		err = customerror.EnCodeError(define.APIPowerErr, "api权限不正确")
		return
	}
	return
}
