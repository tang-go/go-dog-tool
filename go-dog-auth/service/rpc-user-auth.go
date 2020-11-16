package service

import (
	"fmt"
	"omo-service/define"
	"
	param"
	"time"

	"github.com/tang-go/go-dog/lib/md5"
	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/plugins"
)

//CreateUserToken 创建用户token
func (s *Service) CreateUserToken(ctx plugins.Context, request param.CreateUserTokenReq) (response param.CreateUserTokenRes, err error) {
	//生成token缓存
	now := time.Now().Unix()
	//生成登录token
	token := md5.Md5(fmt.Sprintf("%s-%d-%d", request.SvcName, request.User.ID, now))
	if e := s.cache.GetCache().SetByTime(token, request.User, define.UserTokenValidityTime); e != nil {
		log.Errorln(e.Error())
		return response, e
	}
	response.Token = token
	return response, nil
}

//UpdateUserToken 更新用户token
func (s *Service) UpdateUserToken(ctx plugins.Context, request param.UpdateUserTokenReq) (response param.UpdateUserTokenRes, err error) {
	if e := s.cache.GetCache().SetByTime(request.Token, request.User, define.UserTokenValidityTime); e != nil {
		log.Errorln(e.Error())
		return response, e
	}
	response.Success = true
	return response, nil
}

//CheckUserToken 验证用户token
func (s *Service) CheckUserToken(ctx plugins.Context, request param.CheckUserTokenReq) (response param.CheckUserTokenRes, err error) {
	if e := s.cache.GetCache().Get(request.Token, &response.User); e != nil {
		log.Errorln(e.Error())
		return response, e
	}
	return response, nil
}
