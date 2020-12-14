package service

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/tang-go/go-dog-tool/go-dog-auth/param"
	"github.com/tang-go/go-dog-tool/go-dog-auth/table"
	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/plugins"
)

//GetRoleMenu 获取角色菜单
func (s *Service) GetRoleMenu(ctx plugins.Context, request param.GetRoleMenuReq) (response param.GetRoleMenuRes, err error) {
	role := new(table.SysRole)
	err = s.mysql.GetReadEngine().Where("id = ? AND organize = ?", request.RoleID, request.Organize).First(role).Error
	if err != nil {
		log.Errorln(err.Error())
		return
	}

	if role.IsSuper {
		//超级管理员角色直接获取所有的
		var menus []table.SysMenu
		err = s.mysql.GetReadEngine().Where("organize = ?", request.Organize).Find(&menus).Error
		if err != nil {
			log.Errorln(err.Error())
			return
		}
		for _, menu := range menus {
			response.RoleMenus = append(response.RoleMenus, param.RoleMenu{
				ID:       menu.ID,
				ParentID: menu.ParentID,
				Describe: menu.Describe,
				URL:      menu.URL,
				Sort:     menu.Sort,
				Add:      true,
				Del:      true,
				Update:   true,
				Select:   true,
				Time:     menu.Time,
			})
		}
	} else {
		var roleMenus []table.SysRoleMenu
		err = s.mysql.GetReadEngine().Where("role_id = ?", request.RoleID).Find(&roleMenus).Error
		if err != nil {
			log.Errorln(err.Error())
			return
		}
		for _, menu := range roleMenus {
			sysMenu := new(table.SysMenu)
			if s.mysql.GetReadEngine().Where("id = ?", menu.MenuID).First(sysMenu).Error != nil {
				continue
			}
			response.RoleMenus = append(response.RoleMenus, param.RoleMenu{
				ID:       sysMenu.ID,
				ParentID: sysMenu.ParentID,
				Describe: sysMenu.Describe,
				URL:      sysMenu.URL,
				Sort:     sysMenu.Sort,
				Add:      menu.Add,
				Del:      menu.Del,
				Update:   menu.Update,
				Select:   menu.Select,
				Time:     sysMenu.Time,
			})
		}
	}
	return
}

//GetRoleAPI 获取角色API
func (s *Service) GetRoleAPI(ctx plugins.Context, request param.GetRoleAPIReq) (response param.GetRoleAPIRes, err error) {
	role := new(table.SysRole)
	err = s.mysql.GetReadEngine().Where("id = ? AND organize = ?", request.RoleID, request.Organize).First(role).Error
	if err != nil {
		log.Errorln(err.Error())
		return
	}

	if role.IsSuper {
		//超级管理员角色直接获取所有的
		err = s.mysql.GetReadEngine().Where("organize = ?", request.Organize).Find(&response.SysAPI).Error
		if err != nil {
			log.Errorln(err.Error())
			return
		}
	} else {
		var roleAPIs []table.SysRoleAPI
		err = s.mysql.GetReadEngine().Where("role_id = ?", request.RoleID).Find(&roleAPIs).Error
		if err != nil {
			log.Errorln(err.Error())
			return
		}
		var apis []uint
		for _, api := range roleAPIs {
			apis = append(apis, api.APIID)
		}
		err = s.mysql.GetReadEngine().Where("id IN (?)", apis).Find(&response.SysAPI).Error
		if err != nil {
			log.Errorln(err.Error())
			return
		}
	}
	return
}

//CreateRole 创建角色
func (s *Service) CreateRole(ctx plugins.Context, request param.CreateRoleReq) (response param.CreateRoleRes, err error) {
	sysRole := new(table.SysRole)
	if err = s.mysql.GetReadEngine().Where("organize = ? AND name = ?", request.Organize, request.Name).First(sysRole).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
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
		log.Errorln(err.Error())
		return

	}
	response.ID = sysRole.ID
	return
}

//DelRole 删除角色
func (s *Service) DelRole(ctx plugins.Context, request param.DelRoleReq) (response param.DelRoleRes, err error) {
	role := new(table.SysRole)
	err = s.mysql.GetReadEngine().Where("id = ? AND organize = ?", request.ID, request.Organize).First(role).Error
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	tx := s.mysql.GetWriteEngine().Begin()
	if err = tx.Where("role_id = ?", role.ID).Delete(&table.SysRoleMenu{}).Error; err != nil {
		log.Errorln(err.Error())
		tx.Rollback()
		return
	}
	if err = tx.Where("role_id = ?", role.ID).Delete(&table.SysRoleAPI{}).Error; err != nil {
		log.Errorln(err.Error())
		tx.Rollback()
		return
	}
	if err = tx.Where("id = ?", role.ID).Delete(&table.SysRole{}).Error; err != nil {
		log.Errorln(err.Error())
		tx.Rollback()
		return
	}
	tx.Commit()
	response.Success = true
	return
}

//SelectRoleByOrganize 查询角色
func (s *Service) SelectRoleByOrganize(ctx plugins.Context, request param.SelectRoleByOrganizeReq) (response param.SelectRoleByOrganizeRes, err error) {
	err = s.mysql.GetReadEngine().Where("organize = ?", request.Organize).Find(&response.SysRoles).Error
	if err != nil {
		log.Errorln(err.Error())
	}
	return
}

//SelectRoleByID 查询角色
func (s *Service) SelectRoleByID(ctx plugins.Context, request param.SelectRoleByIDReq) (response param.SelectRoleByIDRes, err error) {
	err = s.mysql.GetReadEngine().Where("organize = ? AND id = ?", request.Organize, request.RoleID).Find(&response.SysRole).Error
	if err != nil {
		log.Errorln(err.Error())
	}
	return
}

//CreateMenu 创建菜单
func (s *Service) CreateMenu(ctx plugins.Context, request param.CreateMenuReq) (response param.CreateMenuRes, err error) {
	sysMenu := new(table.SysMenu)
	if s.mysql.GetReadEngine().Where("organize = ? AND url = ?", request.Organize, request.URL).First(sysMenu).RecordNotFound() == true {
		menu := table.SysMenu{
			Organize: request.Organize,
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
	response.ID = sysMenu.ID
	return
}

//SelectMenu 查询角色
func (s *Service) SelectMenu(ctx plugins.Context, request param.SelectMenuReq) (response param.SelectMenuRes, err error) {
	err = s.mysql.GetReadEngine().Where("organize = ? ", request.Organize).Find(&response.Menus).Error
	if err != nil {
		log.Errorln(err.Error())
	}
	return
}

//DelMenu 删除菜单
func (s *Service) DelMenu(ctx plugins.Context, request param.DelMenuReq) (response param.DelMenuRes, err error) {
	if err = s.mysql.GetWriteEngine().Where("organize = ? AND id = ?", request.Organize, request.ID).Delete(&table.SysMenu{}).Error; err != nil {
		log.Errorln(err.Error())
		return
	}
	response.Success = true
	return
}

//CreateAPI 创建api
func (s *Service) CreateAPI(ctx plugins.Context, request param.CreateAPIReq) (response param.CreateAPIRes, err error) {
	//request.API = strings.ToLower(request.API)
	sysAPI := new(table.SysAPI)
	if s.mysql.GetReadEngine().Where("organize = ? AND api = ?", request.Organize, request.API).First(&sysAPI).RecordNotFound() == true {
		api := table.SysAPI{
			Organize: request.Organize,
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
	} else {
		//更新
		if request.Describe != sysAPI.Describe {
			s.mysql.GetWriteEngine().Model(&table.SysAPI{}).Where("id = ?", sysAPI.ID).Updates(map[string]interface{}{
				"Describe": request.Describe,
			})
		}
	}
	response.ID = sysAPI.ID
	return
}

//SelectAPI 查询API
func (s *Service) SelectAPI(ctx plugins.Context, request param.SelectAPIReq) (response param.SelectAPIRes, err error) {
	err = s.mysql.GetReadEngine().Where("organize = ? ", request.Organize).Find(&response.APIS).Error
	if err != nil {
		log.Errorln(err.Error())
	}
	return
}

//DelAPI 删除API
func (s *Service) DelAPI(ctx plugins.Context, request param.DelAPIReq) (response param.DelAPIRes, err error) {
	if err = s.mysql.GetWriteEngine().Where("organize = ? AND id = ?", request.Organize, request.ID).Delete(&table.SysAPI{}).Error; err != nil {
		log.Errorln(err.Error())
		return
	}
	response.Success = true
	return
}

//BindRoleAPI 绑定角色api
func (s *Service) BindRoleAPI(ctx plugins.Context, request param.BindRoleAPIReq) (response param.BindRoleAPIRes, err error) {
	role := new(table.SysRole)
	err = s.mysql.GetReadEngine().Where("id = ? AND organize = ?", request.RoleID, request.Organize).First(role).Error
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	api := new(table.SysAPI)
	err = s.mysql.GetReadEngine().Where("id = ? AND organize = ?", request.APIID, request.Organize).First(api).Error
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	if role.Organize != api.Organize {
		err = errors.New("不是跨组织绑定")
		return
	}
	if s.mysql.GetReadEngine().Where("role_id = ? AND api_id = ?", request.RoleID, request.APIID).First(&table.SysRoleAPI{}).RecordNotFound() == true {
		if err = s.mysql.GetReadEngine().Where("id = ?", request.RoleID).First(&table.SysRole{}).Error; err != nil {
			log.Errorln(err.Error())
			return
		}
		if err = s.mysql.GetReadEngine().Where("id = ?", request.APIID).First(&table.SysAPI{}).Error; err != nil {
			log.Errorln(err.Error())
			return
		}
		roleAPI := table.SysRoleAPI{
			RoleID: request.RoleID,
			APIID:  request.APIID,
			Time:   time.Now().Unix(),
		}
		if err = s.mysql.GetWriteEngine().Create(&roleAPI).Error; err != nil {
			log.Errorln(err.Error())
			return
		}
		response.Success = true
		return
	}
	err = errors.New("请勿重复绑定")
	return
}

//DelRoleAPI 删除角色PAI
func (s *Service) DelRoleAPI(ctx plugins.Context, request param.DelRoleAPIReq) (response param.DelRoleAPIRes, err error) {
	role := new(table.SysRole)
	err = s.mysql.GetReadEngine().Where("id = ? AND organize = ?", request.RoleID, request.Organize).First(role).Error
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	api := new(table.SysAPI)
	err = s.mysql.GetReadEngine().Where("id = ? AND organize = ?", request.APIID, request.Organize).First(api).Error
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	if role.Organize != api.Organize {
		err = errors.New("不是跨组织删除")
		return
	}
	if err = s.mysql.GetWriteEngine().Where("role_id = ? AND api_id = ?", role.ID, api.ID).Delete(&table.SysRoleAPI{}).Error; err != nil {
		log.Errorln(err.Error())
		return
	}
	response.Success = true
	return
}

//DelRoleMenu 删除角色菜单
func (s *Service) DelRoleMenu(ctx plugins.Context, request param.DelRoleMenuReq) (response param.DelRoleMenuRes, err error) {
	role := new(table.SysRole)
	err = s.mysql.GetReadEngine().Where("id = ? AND organize = ?", request.RoleID, request.Organize).First(role).Error
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	menu := new(table.SysMenu)
	err = s.mysql.GetReadEngine().Where("id = ? AND organize = ?", request.MenuID, request.Organize).First(menu).Error
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	if role.Organize != menu.Organize {
		err = errors.New("不是跨组织删除")
		return
	}
	if err = s.mysql.GetWriteEngine().Where("role_id = ? AND menu_id = ?", role.ID, menu.ID).Delete(&table.SysRoleMenu{}).Error; err != nil {
		log.Errorln(err.Error())
		return
	}
	response.Success = true
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
	if s.mysql.GetReadEngine().Where("role_id = ? AND menu_id = ?", request.RoleID, request.MenuID).First(&table.SysRoleMenu{}).RecordNotFound() == true {
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
