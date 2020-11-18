package service

import (
	"sort"
	"strconv"
	"time"

	"github.com/tang-go/go-dog-tool/define"
	authAPI "github.com/tang-go/go-dog-tool/go-dog-auth/api"
	"github.com/tang-go/go-dog-tool/go-dog-ctl/param"
	"github.com/tang-go/go-dog-tool/go-dog-ctl/table"
	customerror "github.com/tang-go/go-dog/error"
	"github.com/tang-go/go-dog/lib/md5"
	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/plugins"
)

//GetAdminInfo 获取管理员信息
func (s *Service) GetAdminInfo(ctx plugins.Context, request param.GetAdminInfoReq) (response param.GetAdminInfoRes, err error) {
	admin, ok := ctx.GetShareByKey("Admin").(*table.Admin)
	if ok == false {
		err = customerror.EnCodeError(define.GetAdminInfoErr, "管理员信息失败")
		return
	}
	response.ID = strconv.FormatInt(admin.AdminID, 10)
	response.Name = admin.Name
	response.Avatar = "/avatar2.jpg"
	response.Phone = admin.Phone
	//获取权限信息
	role, e := authAPI.SelectRoleByID(ctx, define.Organize, admin.RoleID)
	if e != nil {
		log.Errorln(e.Error())
		err = customerror.EnCodeError(define.GetAdminInfoErr, "管理员角色信息失败")
		return
	}
	//获取菜单信息
	menus, e := authAPI.GetRoleMenu(ctx, define.Organize, role.ID)
	if e != nil {
		log.Errorln(e.Error())
		err = customerror.EnCodeError(define.GetAdminInfoErr, "管理员角色菜单失败")
		return
	}
	for _, menu := range menus {
		if menu.ParentID == 0 {
			nm := param.RoleMenu{
				ID:       menu.ID,
				ParentID: menu.ParentID,
				URL:      menu.URL,
				Describe: menu.Describe,
				Add:      menu.Add,
				Del:      menu.Del,
				Update:   menu.Update,
				Select:   menu.Select,
				Sort:     menu.Sort,
				Time:     menu.Time,
			}
			s._AssembleRoleMenu(&nm, menus, menu.ID)
			response.Menu = append(response.Menu, &nm)
		}
	}
	sort.Sort(RoleMenuSort(response.Menu))
	response.RoleID = role.ID
	response.RoleName = role.Name
	return
}

//AdminLogin 管理员登录
func (s *Service) AdminLogin(ctx plugins.Context, request param.AdminLoginReq) (response param.AdminLoginRes, err error) {
	//查询是否拥有此用户
	admin := new(table.Admin)
	if s.mysql.GetReadEngine().Where("phone = ?", request.Phone).First(admin).RecordNotFound() == true {
		err = customerror.EnCodeError(define.AdminLoginErr, "管理员登录失败")
		return
	}
	//密码对比
	if md5.Md5(md5.Md5(request.Pwd)+admin.Salt) != admin.Pwd {
		err = customerror.EnCodeError(define.AdminLoginErr, "管理员登录失败")
		return
	}
	//生成登录记录
	mysqllog := &table.Log{
		//日志ID
		LogID: s.snowflake.GetID(),
		//类型
		Type: table.LoginType,
		//操作人
		AdminID: admin.AdminID,
		//名称
		AdminName: admin.Name,
		//操作方法
		Method: "AdminLogin",
		//描述
		Description: "管理员登录",
		//业主ID
		OwnerID: admin.OwnerID,
		//操作IP
		IP: ctx.GetAddress(),
		//操作URL
		URL: ctx.GetDataByKey("URL").(string),
		//操作时间
		Time: time.Now().Unix(),
	}
	if e := s.mysql.GetWriteEngine().Create(mysqllog).Error; e != nil {
		err = customerror.EnCodeError(define.AdminLoginErr, e.Error())
		return
	}
	//生成token
	token := md5.Md5(admin.AdminID)
	//生成token缓存
	s.cache.GetCache().SetByTime(token, admin, define.AdminTokenValidityTime)
	//登录成功返回
	response.Name = admin.Name
	response.OwnerID = strconv.FormatInt(admin.OwnerID, 10)
	response.Token = token
	return
}
