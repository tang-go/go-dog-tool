package service

import (
	"time"

	"github.com/tang-go/go-dog-tool/define"
	authAPI "github.com/tang-go/go-dog-tool/go-dog-auth/api"
	"github.com/tang-go/go-dog-tool/go-dog-ctl/param"
	"github.com/tang-go/go-dog-tool/go-dog-ctl/table"
	customerror "github.com/tang-go/go-dog/error"
	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/plugins"
)

//CreateMenu 创建菜单
func (s *Service) CreateMenu(ctx plugins.Context, request param.CreateMenuReq) (response param.CreateMenuRes, err error) {
	admin, ok := ctx.GetShareByKey("Admin").(*table.Admin)
	if ok == false {
		err = customerror.EnCodeError(define.GetBuildServiceListErr, "管理员信息失败")
		return
	}
	if _, e := authAPI.CreateMenu(ctx, define.Organize, request.Describe, request.URL, request.ParentID, request.Sort); e != nil {
		log.Errorln(e.Error())
		err = customerror.EnCodeError(define.CreateMenuErr, "创建菜单失败")
		return
	}
	//生成登录记录
	tbLog := &table.Log{
		LogID:       s.snowflake.GetID(),
		Type:        table.CreateMenuType,
		AdminID:     admin.AdminID,
		AdminName:   admin.Name,
		Method:      "CreateMenu",
		Description: "创建菜单",
		OwnerID:     admin.OwnerID,
		IP:          ctx.GetAddress(),
		URL:         ctx.GetDataByKey("URL").(string),
		Time:        time.Now().Unix(),
	}
	if e := s.mysql.GetWriteEngine().Create(&tbLog).Error; e != nil {
		log.Errorln(e.Error())
	}
	response.Success = true
	return
}
