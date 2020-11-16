package service

import (
	"fmt"
	"time"

	"github.com/tang-go/go-dog-tool/define"
	"github.com/tang-go/go-dog-tool/go-dog-auth/param"
	"github.com/tang-go/go-dog/lib/md5"
	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/plugins"
)

//CreateAdminToken 创建管理员token
func (s *Service) CreateAdminToken(ctx plugins.Context, request param.CreateAdminTokenReq) (response param.CreateAdminTokenRes, err error) {
	//生成token缓存
	now := time.Now().Unix()
	//生成登录token
	token := md5.Md5(fmt.Sprintf("%s-%d-%d", request.SvcName, request.Admin.ID, now))
	if e := s.cache.GetCache().SetByTime(token, request.Admin, define.AdminTokenValidityTime); e != nil {
		log.Errorln(e.Error())
		return response, e
	}
	response.Token = token
	return response, nil
}

//UpdateAdminToken 更新管理员token
func (s *Service) UpdateAdminToken(ctx plugins.Context, request param.UpdateAdminTokenReq) (response param.UpdateAdminTokenRes, err error) {
	if e := s.cache.GetCache().SetByTime(request.Token, request.Admin, define.AdminTokenValidityTime); e != nil {
		log.Errorln(e.Error())
		return response, e
	}
	response.Success = true
	return response, nil
}

//CheckAdminToken 验证管理员token
func (s *Service) CheckAdminToken(ctx plugins.Context, request param.CheckAdminTokenReq) (response param.CheckAdminTokenRes, err error) {
	if e := s.cache.GetCache().Get(request.Token, &response.Admin); e != nil {
		log.Errorln(e.Error())
		return response, e
	}
	return response, nil
}
