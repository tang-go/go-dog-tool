package service

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/tang-go/go-dog-tool/define"
	authAPI "github.com/tang-go/go-dog-tool/go-dog-auth/api"
	"github.com/tang-go/go-dog-tool/go-dog-ctl/param"
	"github.com/tang-go/go-dog-tool/go-dog-ctl/table"
	customerror "github.com/tang-go/go-dog/error"
	"github.com/tang-go/go-dog/lib/md5"
	"github.com/tang-go/go-dog/lib/rand"
	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/plugins"
)

//CreateAdmin 创建管理员
func (s *Service) CreateAdmin(ctx plugins.Context, request param.CreateAdminReq) (response param.CreateAdminRes, err error) {
	admin, ok := ctx.GetShareByKey("Admin").(*table.Admin)
	if ok == false {
		err = customerror.EnCodeError(define.GetAdminInfoErr, "管理员信息失败")
		return
	}
	admininfo := new(table.Admin)
	if e := s.mysql.GetReadEngine().Where("owner_id = ? AND phone = ?", admin.OwnerID, request.Phone).First(&admininfo).Error; e != nil {
		if e.Error() == gorm.ErrRecordNotFound.Error() {
			admininfo.AdminID = s.snowflake.GetID()
			admininfo.Name = request.Name
			admininfo.Phone = request.Phone
			admininfo.Salt = rand.StringRand(6)
			admininfo.OwnerID = admin.OwnerID
			admininfo.IsDisable = table.AdminAvailable
			admininfo.RoleID = request.RoleID
			admininfo.Time = time.Now().Unix()
			admininfo.Pwd = md5.Md5(md5.Md5(request.Pwd) + admininfo.Salt)
			tbLog := &table.Log{
				LogID:       s.snowflake.GetID(),
				Type:        table.CreateAdminType,
				AdminID:     admin.AdminID,
				AdminName:   admin.Name,
				Method:      "CreateAdmin",
				Description: "创建管理员",
				OwnerID:     admin.OwnerID,
				IP:          ctx.GetAddress(),
				URL:         ctx.GetDataByKey("URL").(string),
				Time:        time.Now().Unix(),
			}
			tx := s.mysql.GetWriteEngine().Begin()
			if e := tx.Create(admininfo).Error; e != nil {
				log.Errorln(e.Error())
				tx.Rollback()
				err = customerror.EnCodeError(define.GetAdminInfoErr, "创建管理员失败")
				return
			}
			if e := tx.Create(tbLog).Error; e != nil {
				log.Errorln(e.Error())
				tx.Rollback()
				err = customerror.EnCodeError(define.GetAdminInfoErr, "创建管理员失败")
				return
			}
			tx.Commit()
			response.Success = true
			return
		}
		log.Errorln(e.Error())
		err = customerror.EnCodeError(define.CreateAdminErr, "创建管理员失败")
		return
	}
	err = customerror.EnCodeError(define.CreateAdminErr, "此电话号码已经存在")
	return
}

//DisableAdmin 禁用管理员
func (s *Service) DisableAdmin(ctx plugins.Context, request param.DisableAdminReq) (response param.DisableAdminRes, err error) {
	admin, ok := ctx.GetShareByKey("Admin").(*table.Admin)
	if ok == false {
		err = customerror.EnCodeError(define.GetAdminInfoErr, "管理员信息失败")
		return
	}
	admininfo := new(table.Admin)
	if e := s.mysql.GetReadEngine().Where("owner_id = ? AND admin_id = ?", admin.OwnerID, request.AdminID).First(&admininfo).Error; e != nil {
		log.Errorln(e.Error())
		err = customerror.EnCodeError(define.DisableAdminErr, "管理员ID错误")
		return
	}
	tbLog := &table.Log{
		LogID:       s.snowflake.GetID(),
		Type:        table.DisableAdminType,
		AdminID:     admin.AdminID,
		AdminName:   admin.Name,
		Method:      "DisableAdmin",
		Description: "禁用管理员",
		OwnerID:     admin.OwnerID,
		IP:          ctx.GetAddress(),
		URL:         ctx.GetDataByKey("URL").(string),
		Time:        time.Now().Unix(),
	}
	tx := s.mysql.GetWriteEngine().Begin()
	if e := tx.Model(&table.Admin{}).Where("admin_id = ?", request.AdminID).Updates(map[string]interface{}{"IsDisable": table.AdminDisable}).Error; e != nil {
		log.Errorln(e.Error())
		tx.Rollback()
		err = customerror.EnCodeError(define.DisableAdminErr, "禁用管理员失败")
		return
	}
	if e := tx.Create(tbLog).Error; e != nil {
		log.Errorln(e.Error())
		tx.Rollback()
		err = customerror.EnCodeError(define.DisableAdminErr, "禁用管理员失败")
		return
	}
	tx.Commit()
	response.Success = true
	return
}

//OpenAdmin 开启管理员
func (s *Service) OpenAdmin(ctx plugins.Context, request param.OpenAdminReq) (response param.OpenAdminRes, err error) {
	admin, ok := ctx.GetShareByKey("Admin").(*table.Admin)
	if ok == false {
		err = customerror.EnCodeError(define.GetAdminInfoErr, "管理员信息失败")
		return
	}
	admininfo := new(table.Admin)
	if e := s.mysql.GetReadEngine().Where("owner_id = ? AND admin_id = ?", admin.OwnerID, request.AdminID).First(&admininfo).Error; e != nil {
		log.Errorln(e.Error())
		err = customerror.EnCodeError(define.OpenAdminErr, "管理员ID错误")
		return
	}
	tbLog := &table.Log{
		LogID:       s.snowflake.GetID(),
		Type:        table.OpenAdminType,
		AdminID:     admin.AdminID,
		AdminName:   admin.Name,
		Method:      "OpenAdmin",
		Description: "开启管理员",
		OwnerID:     admin.OwnerID,
		IP:          ctx.GetAddress(),
		URL:         ctx.GetDataByKey("URL").(string),
		Time:        time.Now().Unix(),
	}
	tx := s.mysql.GetWriteEngine().Begin()
	if e := tx.Model(&table.Admin{}).Where("admin_id = ?", request.AdminID).Updates(map[string]interface{}{"IsDisable": table.AdminAvailable}).Error; e != nil {
		log.Errorln(e.Error())
		tx.Rollback()
		err = customerror.EnCodeError(define.OpenAdminErr, "开启管理员失败")
		return
	}
	if e := tx.Create(tbLog).Error; e != nil {
		log.Errorln(e.Error())
		tx.Rollback()
		err = customerror.EnCodeError(define.OpenAdminErr, "开启管理员失败")
		return
	}
	tx.Commit()
	response.Success = true
	return
}

//DelAdmin 删除管理员
func (s *Service) DelAdmin(ctx plugins.Context, request param.DelAdminReq) (response param.DelAdminRes, err error) {
	admin, ok := ctx.GetShareByKey("Admin").(*table.Admin)
	if ok == false {
		err = customerror.EnCodeError(define.GetAdminInfoErr, "管理员信息失败")
		return
	}
	admininfo := new(table.Admin)
	if e := s.mysql.GetReadEngine().Where("owner_id = ? AND admin_id = ?", admin.OwnerID, request.AdminID).First(&admininfo).Error; e != nil {
		log.Errorln(e.Error())
		err = customerror.EnCodeError(define.DelAdminErr, "管理员ID错误")
		return
	}
	tbLog := &table.Log{
		LogID:       s.snowflake.GetID(),
		Type:        table.DelAdminType,
		AdminID:     admin.AdminID,
		AdminName:   admin.Name,
		Method:      "DelAdmin",
		Description: "删除管理员",
		OwnerID:     admin.OwnerID,
		IP:          ctx.GetAddress(),
		URL:         ctx.GetDataByKey("URL").(string),
		Time:        time.Now().Unix(),
	}
	tx := s.mysql.GetWriteEngine().Begin()
	if e := tx.Where("admin_id = ?", request.AdminID).Delete(&table.Admin{}).Error; e != nil {
		log.Errorln(e.Error())
		tx.Rollback()
		err = customerror.EnCodeError(define.DelAdminErr, "删除管理员失败")
		return
	}
	if e := tx.Create(tbLog).Error; e != nil {
		log.Errorln(e.Error())
		tx.Rollback()
		err = customerror.EnCodeError(define.DelAdminErr, "删除管理员失败")
		return
	}
	tx.Commit()
	response.Success = true
	return
}

//GetAdminList 获取管理员列表请求
func (s *Service) GetAdminList(ctx plugins.Context, request param.GetAdminListReq) (response param.GetAdminListRes, err error) {
	admin, ok := ctx.GetShareByKey("Admin").(*table.Admin)
	if ok == false {
		err = customerror.EnCodeError(define.GetAdminInfoErr, "管理员信息失败")
		return
	}
	var admins []table.Admin
	if e := s.mysql.GetReadEngine().Where("owner_id = ?", admin.OwnerID).Find(&admins).Error; e != nil {
		log.Errorln(e.Error())
		err = customerror.EnCodeError(define.GetAdminListErr, "获取管理员列表失败")
		return
	}
	for _, item := range admins {
		admininfo := param.AdminInfo{
			AdminID:   fmt.Sprintf("%d", item.AdminID),
			Name:      item.Name,
			Phone:     item.Phone,
			IsDisable: item.IsDisable,
			RoleID:    item.RoleID,
			Time:      time.Unix(item.Time, 0).Format("2006-01-02 15:04:05"),
		}
		//获取权限信息
		role, e := authAPI.SelectRoleByID(ctx, define.Organize, item.RoleID)
		if e != nil {
			log.Errorln(e.Error())
			err = customerror.EnCodeError(define.GetAdminInfoErr, "管理员角色信息失败")
			return
		}
		//获取菜单信息
		menus, e := authAPI.GetRoleMenu(ctx, define.Organize, item.RoleID)
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
				admininfo.Menu = append(admininfo.Menu, &nm)
			}
		}
		sort.Sort(RoleMenuSort(admininfo.Menu))
		//获取api
		apis, e := authAPI.GetRoleAPI(ctx, define.Organize, item.RoleID)
		if e != nil {
			log.Errorln(e.Error())
			err = customerror.EnCodeError(define.GetAdminInfoErr, "管理员角色菜单失败")
			return
		}
		for _, api := range apis {
			admininfo.APIS = append(admininfo.APIS, param.API{
				ID:       api.ID,
				Organize: api.Organize,
				API:      api.API,
				Describe: api.Describe,
				Time:     time.Unix(api.Time, 0).Format("2006-01-02 15:04:05"),
			})
		}
		admininfo.RoleName = role.Name
		response.AdminInfos = append(response.AdminInfos, admininfo)
	}
	return
}

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
	//获取api
	apis, e := authAPI.GetRoleAPI(ctx, define.Organize, admin.RoleID)
	if e != nil {
		log.Errorln(e.Error())
		err = customerror.EnCodeError(define.GetAdminInfoErr, "管理员角色菜单失败")
		return
	}
	for _, api := range apis {
		response.APIS = append(response.APIS, param.API{
			ID:       api.ID,
			Organize: api.Organize,
			API:      api.API,
			Describe: api.Describe,
			Time:     time.Unix(api.Time, 0).Format("2006-01-02 15:04:05"),
		})
	}
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
		log.Warnln("密码不正确", request.Pwd, md5.Md5(md5.Md5(request.Pwd)+admin.Salt))
		err = customerror.EnCodeError(define.AdminLoginErr, "管理员登录失败")
		return
	}
	//管理员不可用
	if admin.IsDisable == table.AdminDisable {
		err = customerror.EnCodeError(define.AdminLoginErr, "此管理员已经禁用")
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
	token := md5.Md5(fmt.Sprintf("%d-%d", admin.AdminID, time.Now().Unix()))
	//生成token缓存
	s.cache.GetCache().SetByTime(token, admin, define.AdminTokenValidityTime)
	//登录成功返回
	response.Name = admin.Name
	response.OwnerID = strconv.FormatInt(admin.OwnerID, 10)
	response.Token = token
	return
}
