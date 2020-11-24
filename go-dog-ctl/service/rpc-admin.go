package service

import (
	"github.com/tang-go/go-dog-tool/define"
	"github.com/tang-go/go-dog-tool/go-dog-ctl/param"
	"github.com/tang-go/go-dog-tool/go-dog-ctl/table"
	customerror "github.com/tang-go/go-dog/error"
	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/plugins"
)

//AdminOnline 管理员上线
func (s *Service) AdminOnline(ctx plugins.Context, request param.AdminOnlineReq) (response param.AdminOnlineRes, err error) {
	admin := new(table.Admin)
	if e := s.cache.GetCache().Get(ctx.GetToken(), admin); e != nil {
		log.Errorln(e.Error())
		err = customerror.EnCodeError(define.GetAdminInfoErr, "管理员信息失败")
		return
	}
	//如果还在想
	if admin.IsOnline {
		if admin.GateAddress != request.Address {
			err = customerror.EnCodeError(define.AdminOnlineErr, "请勿重复登陆")
			return
		}
	}
	admin.GateAddress = request.Address
	admin.IsOnline = true
	//生成token缓存
	s.cache.GetCache().SetByTime(ctx.GetToken(), admin, define.AdminTokenValidityTime)
	response.Success = true
	return
}

//AdminOffline 管理员下线
func (s *Service) AdminOffline(ctx plugins.Context, request param.AdminOfflineReq) (response param.AdminOfflineRes, err error) {
	admin := new(table.Admin)
	if e := s.cache.GetCache().Get(ctx.GetToken(), admin); e != nil {
		log.Errorln(e.Error())
		err = customerror.EnCodeError(define.GetAdminInfoErr, "管理员信息失败")
		return
	}
	if admin.GateAddress != "" || admin.IsOnline {
		//设置不在线
		admin.GateAddress = ""
		admin.IsOnline = false
		//生成token缓存
		s.cache.GetCache().Del(ctx.GetToken())
	}
	response.Success = true
	return
}

//AuthAdmin 验证管理员
func (s *Service) AuthAdmin(ctx plugins.Context, request param.AuthAdminReq) (response param.AuthAdminRes, err error) {
	admin := new(table.Admin)
	if e := s.cache.GetCache().Get(request.Token, admin); e != nil {
		log.Errorln(e.Error())
		err = customerror.EnCodeError(define.AdminTokenErr, "token失效或者不正确")
		return
	}
	response.Success = true
	return
}
