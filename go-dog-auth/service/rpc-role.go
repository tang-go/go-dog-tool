package service

import (
	"errors"
	"time"

	"github.com/tang-go/go-dog-tool/go-dog-auth/param"
	"github.com/tang-go/go-dog-tool/go-dog-auth/table"
	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/plugins"
)

//GetRoleMenu 获取角色菜单
func (s *Service) GetRoleMenu(ctx plugins.Context, request param.GetRoleMenuReq) (response param.GetRoleMenuRes, err error) {
	role := new(table.SysRole)
	err = s.mysql.GetReadEngine().Where("id = ?", request.RoleID).First(role).Error
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	var roleMenus []table.SysRoleMenu
	if role.IsSuper {
		//超级管理员角色直接获取所有的
		err = s.mysql.GetReadEngine().Find(&roleMenus).Error
		if err != nil {
			log.Errorln(err.Error())
			return
		}
	} else {
		err = s.mysql.GetReadEngine().Where("role_id = ?", request.RoleID).Find(&roleMenus).Error
		if err != nil {
			log.Errorln(err.Error())
			return
		}
	}
	for _, menu := range roleMenus {
		sysMenu := new(table.SysMenu)
		if s.mysql.GetReadEngine().Where("id = ?", menu.MenuID).First(sysMenu).Error != nil {
			continue
		}
		response.SysMenu = append(response.SysMenu, param.SysMenu{
			ID:       sysMenu.ID,
			ParentID: sysMenu.ParentID,
			Describe: sysMenu.Describe,
			URL:      sysMenu.Describe,
			Sort:     sysMenu.Sort,
			Add:      menu.Add,
			Del:      menu.Del,
			Update:   menu.Update,
			Select:   menu.Select,
			Time:     sysMenu.Time,
		})
	}
	return
}

//GetRoleApi 获取角色Api
func (s *Service) GetRoleApi(ctx plugins.Context, request param.GetRoleApiReq) (response param.GetRoleApiRes, err error) {
	role := new(table.SysRole)
	err = s.mysql.GetReadEngine().Where("id = ?", request.RoleID).First(role).Error
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	var roleApis []table.SysRoleApi
	if role.IsSuper {
		//超级管理员角色直接获取所有的
		err = s.mysql.GetReadEngine().Find(&roleApis).Error
		if err != nil {
			log.Errorln(err.Error())
			return
		}
	} else {
		err = s.mysql.GetReadEngine().Where("role_id = ?", request.RoleID).Find(&roleApis).Error
		if err != nil {
			log.Errorln(err.Error())
			return
		}
	}
	var apis []uint
	for _, api := range roleApis {
		apis = append(apis, api.ApiID)
	}
	err = s.mysql.GetReadEngine().Where("id IN (?)", apis).Find(&response.SysApi).Error
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	return
}

//CreateRole 创建角色
func (s *Service) CreateRole(ctx plugins.Context, request param.CreateRoleReq) (response param.CreateRoleRes, err error) {
	if s.mysql.GetReadEngine().Where("organize = ? AND name = ?", request.Organize, request.Name).First(&table.SysRole{}).RecordNotFound() == true {
		role := table.SysRole{
			Organize: request.Organize,
			Name:     request.Name,
			Describe: request.Describe,
			IsSuper:  request.IsSuper,
			Time:     time.Now().Unix(),
		}
		if err = s.mysql.GetWriteEngine().Create(&role).Error; err != nil {
			log.Errorln(err.Error())
			return
		}
		response.ID = role.ID
		return
	}
	err = errors.New("权限名称已经存在")
	return
}

//SelectRole 查询角色
func (s *Service) SelectRole(ctx plugins.Context, request param.SelectRoleReq) (response param.SelectRoleRes, err error) {
	err = s.mysql.GetReadEngine().Where("organize = ?", request.Organize).Find(&response.SysRoles).Error
	if err != nil {
		log.Errorln(err.Error())
	}
	return
}

//CreateMenu 创建菜单
func (s *Service) CreateMenu(ctx plugins.Context, request param.CreateMenuReq) (response param.CreateMenuRes, err error) {
	if s.mysql.GetReadEngine().Where("url = ?", request.URL).First(&table.SysMenu{}).RecordNotFound() == true {
		menu := table.SysMenu{
			URL:      request.URL,
			Describe: request.Describe,
			ParentID: request.ParentID,
			Sort:     request.Sort,
			Time:     time.Now().Unix(),
		}
		if err = s.mysql.GetWriteEngine().Create(&menu).Error; err != nil {
			log.Errorln(err.Error())
			return
		}
		response.ID = menu.ID
		return
	}
	err = errors.New("菜单已经存在")
	return
}

//CreateApi 创建api
func (s *Service) CreateApi(ctx plugins.Context, request param.CreateApiReq) (response param.CreateApiRes, err error) {
	if s.mysql.GetReadEngine().Where("api = ?", request.API).First(&table.SysApi{}).RecordNotFound() == true {
		api := table.SysApi{
			API:      request.API,
			Describe: request.Describe,
			Time:     time.Now().Unix(),
		}
		if err = s.mysql.GetWriteEngine().Create(&api).Error; err != nil {
			log.Errorln(err.Error())
			return
		}
		response.ID = api.ID
		return
	}
	err = errors.New("api已经存在")
	return
}

//BindRoleApi 绑定角色api
func (s *Service) BindRoleApi(ctx plugins.Context, request param.BindRoleApiReq) (response param.BindRoleApiRes, err error) {
	role := new(table.SysRole)
	err = s.mysql.GetReadEngine().Where("id = ?", request.RoleID).First(role).Error
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	api := new(table.SysApi)
	err = s.mysql.GetReadEngine().Where("id = ?", request.ApiID).First(api).Error
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	if role.Organize != api.Organize {
		err = errors.New("不是跨组织绑定")
		return
	}
	if s.mysql.GetReadEngine().Where("role_id = ? AND api_id = ?", request.RoleID, request.ApiID).First(&table.SysRoleApi{}).RecordNotFound() == true {
		if err = s.mysql.GetReadEngine().Where("id = ?", request.RoleID).First(&table.SysRole{}).Error; err != nil {
			log.Errorln(err.Error())
			return
		}
		if err = s.mysql.GetReadEngine().Where("id = ?", request.ApiID).First(&table.SysApi{}).Error; err != nil {
			log.Errorln(err.Error())
			return
		}
		roleApi := table.SysRoleApi{
			RoleID: request.RoleID,
			ApiID:  request.ApiID,
			Time:   time.Now().Unix(),
		}
		if err = s.mysql.GetWriteEngine().Create(&roleApi).Error; err != nil {
			log.Errorln(err.Error())
			return
		}
		response.Success = true
		return
	}
	err = errors.New("请勿重复绑定")
	return
}

//BindRoleMenu 绑定角色菜单
func (s *Service) BindRoleMenu(ctx plugins.Context, request param.BindRoleMenuReq) (response param.BindRoleMenuRes, err error) {
	role := new(table.SysRole)
	err = s.mysql.GetReadEngine().Where("id = ?", request.RoleID).First(role).Error
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	menu := new(table.SysMenu)
	err = s.mysql.GetReadEngine().Where("id = ?", request.MenuID).First(menu).Error
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	if role.Organize != menu.Organize {
		err = errors.New("不是跨组织绑定")
		return
	}
	if s.mysql.GetReadEngine().Where("role_id = ? AND menu_id = ?", request.RoleID, request.MenuID).First(&table.SysRoleApi{}).RecordNotFound() == true {
		if err = s.mysql.GetReadEngine().Where("id = ?", request.RoleID).First(&table.SysRole{}).Error; err != nil {
			log.Errorln(err.Error())
			return
		}
		if err = s.mysql.GetReadEngine().Where("id = ?", request.MenuID).First(&table.SysMenu{}).Error; err != nil {
			log.Errorln(err.Error())
			return
		}
		roleMenu := table.SysRoleMenu{
			RoleID: request.RoleID,
			MenuID: request.MenuID,
			Add:    request.Add,
			Del:    request.Del,
			Update: request.Update,
			Select: request.Select,
			Time:   time.Now().Unix(),
		}
		if err = s.mysql.GetWriteEngine().Create(&roleMenu).Error; err != nil {
			log.Errorln(err.Error())
			return
		}
		response.Success = true
		return
	}
	err = errors.New("请勿重复绑定")
	return
}
