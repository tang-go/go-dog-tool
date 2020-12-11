package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/tang-go/go-dog-tool/define"
	authParam "github.com/tang-go/go-dog-tool/go-dog-auth/param"
	authRPC "github.com/tang-go/go-dog-tool/go-dog-auth/rpc"
	authTable "github.com/tang-go/go-dog-tool/go-dog-auth/table"
	"github.com/tang-go/go-dog-tool/go-dog-ctl/param"
	"github.com/tang-go/go-dog-tool/go-dog-ctl/table"
	customerror "github.com/tang-go/go-dog/error"
	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/plugins"
)

//MenuSort 菜单排序
type MenuSort []*param.Menu

//Len 长度
func (m MenuSort) Len() int { return len(m) }

//Swap 交换
func (m MenuSort) Swap(i, j int) { m[i], m[j] = m[j], m[i] }

//Less 比较
func (m MenuSort) Less(i, j int) bool {
	if m[i].Sort == m[j].Sort {
		return m[i].Time > m[j].Time
	}
	return m[i].Sort > m[j].Sort
}

//_AssembleMenu 组织菜单
func (s *Service) _AssembleMenu(m *param.Menu, source []authTable.SysMenu, parentID uint) {
	for _, menu := range source {
		if menu.ParentID == parentID {
			nm := param.Menu{
				ID:       menu.ID,
				ParentID: menu.ParentID,
				Organize: menu.Organize,
				URL:      menu.URL,
				Describe: menu.Describe,
				Sort:     menu.Sort,
				Time:     time.Unix(menu.Time, 0).Format("2006-01-02 15:04:05"),
			}
			s._AssembleMenu(&nm, source, menu.ID)
			m.Children = append(m.Children, &nm)
		}
	}
	sort.Sort(MenuSort(m.Children))
	return
}

//RoleMenuSort 菜单排序
type RoleMenuSort []*param.RoleMenu

//Len 长度
func (m RoleMenuSort) Len() int { return len(m) }

//Swap 交换
func (m RoleMenuSort) Swap(i, j int) { m[i], m[j] = m[j], m[i] }

//Less 比较
func (m RoleMenuSort) Less(i, j int) bool {
	if m[i].Sort == m[j].Sort {
		return m[i].Time > m[j].Time
	}
	return m[i].Sort > m[j].Sort
}

//_AssembleMenu 组织菜单
func (s *Service) _AssembleRoleMenu(m *param.RoleMenu, source []authParam.RoleMenu, parentID uint) {
	for _, menu := range source {
		if menu.ParentID == parentID {
			nm := param.RoleMenu{
				ID:       menu.ID,
				ParentID: menu.ParentID,
				URL:      menu.URL,
				Add:      menu.Add,
				Describe: menu.Describe,
				Del:      menu.Del,
				Update:   menu.Update,
				Select:   menu.Select,
				Sort:     menu.Sort,
				Time:     menu.Time,
			}
			s._AssembleRoleMenu(&nm, source, menu.ID)
			m.Children = append(m.Children, &nm)
		}
	}
	sort.Sort(RoleMenuSort(m.Children))
	return
}

//GetMenu 获取菜单
func (s *Service) GetMenu(ctx plugins.Context, request param.GetMenuReq) (response param.GetMenuRes, err error) {
	//获取菜单信息
	menus, e := authRPC.SelectMenu(ctx, define.Organize)
	if e != nil {
		log.Errorln(e.Error())
		err = customerror.EnCodeError(define.GetMenuErr, "管理员角色菜单失败")
		return
	}
	for _, menu := range menus {
		if menu.ParentID == 0 {
			nm := param.Menu{
				ID:       menu.ID,
				ParentID: menu.ParentID,
				Organize: menu.Organize,
				URL:      menu.URL,
				Describe: menu.Describe,
				Sort:     menu.Sort,
				Time:     time.Unix(menu.Time, 0).Format("2006-01-02 15:04:05"),
			}
			s._AssembleMenu(&nm, menus, menu.ID)
			response.Menu = append(response.Menu, &nm)
		}
	}
	sort.Sort(MenuSort(response.Menu))
	return
}

//CreateMenu 创建菜单
func (s *Service) CreateMenu(ctx plugins.Context, request param.CreateMenuReq) (response param.CreateMenuRes, err error) {
	admin, e := s.GetAdmin(ctx)
	if e != nil {
		log.Errorln(e.Error())
		err = customerror.EnCodeError(define.GetAdminInfoErr, e.Error())
		return
	}
	if _, e := authRPC.CreateMenu(ctx, define.Organize, request.Describe, request.URL, request.ParentID, request.Sort); e != nil {
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
		URL:         ctx.GetURL(),
		Time:        time.Now().Unix(),
	}
	if e := s.mysql.GetWriteEngine().Create(&tbLog).Error; e != nil {
		log.Errorln(e.Error())
	}
	response.Success = true
	return
}

//DelMenu 删除菜单
func (s *Service) DelMenu(ctx plugins.Context, request param.DelMenuReq) (response param.DelMenuRes, err error) {
	admin, e := s.GetAdmin(ctx)
	if e != nil {
		log.Errorln(e.Error())
		err = customerror.EnCodeError(define.GetAdminInfoErr, e.Error())
		return
	}
	if _, e := authRPC.DelMenu(ctx, define.Organize, request.MenuID); e != nil {
		log.Errorln(e.Error())
		err = customerror.EnCodeError(define.DelMenuErr, "删除菜单失败")
		return
	}
	//生成登录记录
	tbLog := &table.Log{
		LogID:       s.snowflake.GetID(),
		Type:        table.DelMenuType,
		AdminID:     admin.AdminID,
		AdminName:   admin.Name,
		Method:      "DelMenu",
		Description: "删除菜单",
		OwnerID:     admin.OwnerID,
		IP:          ctx.GetAddress(),
		URL:         ctx.GetURL(),
		Time:        time.Now().Unix(),
	}
	if e := s.mysql.GetWriteEngine().Create(&tbLog).Error; e != nil {
		log.Errorln(e.Error())
	}
	response.Success = true
	return
}

//GetAPIList 获取API列表
func (s *Service) GetAPIList(ctx plugins.Context, request param.GetAPIListReq) (response param.GetAPIListRes, err error) {
	apis, e := authRPC.SelectAPI(ctx, define.Organize)
	if e != nil {
		log.Errorln(e.Error())
		err = customerror.EnCodeError(define.GetAPIListErr, "获取API列表失败")
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
	return
}

//DelAPI 删除API
func (s *Service) DelAPI(ctx plugins.Context, request param.DelAPIReq) (response param.DelAPIRes, err error) {
	admin, e := s.GetAdmin(ctx)
	if e != nil {
		log.Errorln(e.Error())
		err = customerror.EnCodeError(define.GetAdminInfoErr, e.Error())
		return
	}
	_, e = authRPC.DelAPI(ctx, define.Organize, request.ID)
	if e != nil {
		log.Errorln(e.Error())
		err = customerror.EnCodeError(define.DelAPIErr, "删除API失败")
		return
	}
	//生成记录
	tbLog := &table.Log{
		LogID:       s.snowflake.GetID(),
		Type:        table.DelAPIType,
		AdminID:     admin.AdminID,
		AdminName:   admin.Name,
		Method:      "DelAPI",
		Description: "删除API",
		OwnerID:     admin.OwnerID,
		IP:          ctx.GetAddress(),
		URL:         ctx.GetURL(),
		Time:        time.Now().Unix(),
	}
	if e := s.mysql.GetWriteEngine().Create(&tbLog).Error; e != nil {
		log.Errorln(e.Error())
	}
	response.Success = true
	return
}

//GetRoleList 获取角色列表
func (s *Service) GetRoleList(ctx plugins.Context, request param.GetRoleListReq) (response param.GetRoleListRes, err error) {
	roles, e := authRPC.SelectRoleByOrganize(ctx, define.Organize)
	if e != nil {
		log.Errorln(e.Error())
		err = customerror.EnCodeError(define.GetRoleListErr, "管理员角色菜单失败")
		return
	}
	for _, role := range roles {
		response.Roles = append(response.Roles, param.Role{
			ID:       role.ID,
			Organize: role.Organize,
			Name:     role.Name,
			Describe: role.Describe,
			IsSuper:  role.IsSuper,
			Time:     time.Unix(role.Time, 0).Format("2006-01-02 15:04:05"),
		})
	}
	return
}

//CreateRole 创建角色
func (s *Service) CreateRole(ctx plugins.Context, request param.CreateRoleReq) (response param.CreateRoleRes, err error) {
	admin, e := s.GetAdmin(ctx)
	if e != nil {
		log.Errorln(e.Error())
		err = customerror.EnCodeError(define.GetAdminInfoErr, e.Error())
		return
	}
	id, e := authRPC.CreateRole(ctx, define.Organize, request.Name, request.Describe, false)
	if e != nil {
		log.Errorln(e.Error())
		err = customerror.EnCodeError(define.CreateRoleErr, "创建角色失败")
		return
	}
	//生成记录
	tbLog := &table.Log{
		LogID:       s.snowflake.GetID(),
		Type:        table.CreateRoleType,
		AdminID:     admin.AdminID,
		AdminName:   admin.Name,
		Method:      "CreateRole",
		Description: "创建角色",
		OwnerID:     admin.OwnerID,
		IP:          ctx.GetAddress(),
		URL:         ctx.GetURL(),
		Time:        time.Now().Unix(),
	}
	if e := s.mysql.GetWriteEngine().Create(&tbLog).Error; e != nil {
		log.Errorln(e.Error())
	}
	response.ID = id
	return
}

//DelRole 删除角色
func (s *Service) DelRole(ctx plugins.Context, request param.DelRoleReq) (response param.DelRoleRes, err error) {
	admin, e := s.GetAdmin(ctx)
	if e != nil {
		log.Errorln(e.Error())
		err = customerror.EnCodeError(define.GetAdminInfoErr, e.Error())
		return
	}
	if _, e := authRPC.DelRole(ctx, define.Organize, request.ID); e != nil {
		log.Errorln(e.Error())
		err = customerror.EnCodeError(define.DelRoleErr, "删除角色失败")
		return
	}
	//生成登录记录
	tbLog := &table.Log{
		LogID:       s.snowflake.GetID(),
		Type:        table.DelRoleType,
		AdminID:     admin.AdminID,
		AdminName:   admin.Name,
		Method:      "DelRole",
		Description: "删除角色",
		OwnerID:     admin.OwnerID,
		IP:          ctx.GetAddress(),
		URL:         ctx.GetURL(),
		Time:        time.Now().Unix(),
	}
	if e := s.mysql.GetWriteEngine().Create(&tbLog).Error; e != nil {
		log.Errorln(e.Error())
	}
	response.Success = true
	return
}

//GetRoleMenu 获取角色菜单列表
func (s *Service) GetRoleMenu(ctx plugins.Context, request param.GetRoleMenuReq) (response param.GetRoleMenuRes, err error) {
	//获取菜单信息
	menus, e := authRPC.GetRoleMenu(ctx, define.Organize, request.RoleID)
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
	return
}

//BindRoleMenu 绑定角色菜单
func (s *Service) BindRoleMenu(ctx plugins.Context, request param.BindRoleMenuReq) (response param.BindRoleMenuRes, err error) {
	admin, e := s.GetAdmin(ctx)
	if e != nil {
		log.Errorln(e.Error())
		err = customerror.EnCodeError(define.GetAdminInfoErr, e.Error())
		return
	}
	_, e = authRPC.BindRoleMenu(ctx, define.Organize, request.MenuID, request.RoleID, request.Add, request.Del, request.Update, request.Select)
	if e != nil {
		log.Errorln(e.Error())
		err = customerror.EnCodeError(define.BindRoleMenuErr, "绑定角色菜单失败")
		return
	}
	//生成记录
	tbLog := &table.Log{
		LogID:       s.snowflake.GetID(),
		Type:        table.CreateRoleType,
		AdminID:     admin.AdminID,
		AdminName:   admin.Name,
		Method:      "BindRoleMenu",
		Description: "绑定角色菜单",
		OwnerID:     admin.OwnerID,
		IP:          ctx.GetAddress(),
		URL:         ctx.GetURL(),
		Time:        time.Now().Unix(),
	}
	if e := s.mysql.GetWriteEngine().Create(&tbLog).Error; e != nil {
		log.Errorln(e.Error())
	}
	response.Success = true
	return
}

//DelRoleMenu 删除角色菜单
func (s *Service) DelRoleMenu(ctx plugins.Context, request param.DelRoleMenuReq) (response param.DelRoleMenuRes, err error) {
	admin, e := s.GetAdmin(ctx)
	if e != nil {
		log.Errorln(e.Error())
		err = customerror.EnCodeError(define.GetAdminInfoErr, e.Error())
		return
	}
	if _, e := authRPC.DelRoleMenu(ctx, define.Organize, request.RoleID, request.MenuID); e != nil {
		log.Errorln(e.Error())
		err = customerror.EnCodeError(define.DelRoleErr, "删除角色失败")
		return
	}
	//生成登录记录
	tbLog := &table.Log{
		LogID:       s.snowflake.GetID(),
		Type:        table.DelRoleType,
		AdminID:     admin.AdminID,
		AdminName:   admin.Name,
		Method:      "DelRole",
		Description: "删除角色菜单",
		OwnerID:     admin.OwnerID,
		IP:          ctx.GetAddress(),
		URL:         ctx.GetURL(),
		Time:        time.Now().Unix(),
	}
	if e := s.mysql.GetWriteEngine().Create(&tbLog).Error; e != nil {
		log.Errorln(e.Error())
	}
	response.Success = true
	return
}

//GetRoleAPI 获取角色API列表
func (s *Service) GetRoleAPI(ctx plugins.Context, request param.GetRoleAPIReq) (response param.GetRoleAPIRes, err error) {
	apis, e := authRPC.GetRoleAPI(ctx, define.Organize, request.RoleID)
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
	return
}

//BindRoleAPI 绑定角色API
func (s *Service) BindRoleAPI(ctx plugins.Context, request param.BindRoleAPIReq) (response param.BindRoleAPIRes, err error) {
	admin, e := s.GetAdmin(ctx)
	if e != nil {
		log.Errorln(e.Error())
		err = customerror.EnCodeError(define.GetAdminInfoErr, e.Error())
		return
	}
	_, e = authRPC.BindRoleAPI(ctx, define.Organize, request.RoleID, request.APIID)
	if e != nil {
		log.Errorln(e.Error())
		err = customerror.EnCodeError(define.BindRoleAPIErr, "绑定角色API失败")
		return
	}
	//生成记录
	tbLog := &table.Log{
		LogID:       s.snowflake.GetID(),
		Type:        table.BindRoleAPIType,
		AdminID:     admin.AdminID,
		AdminName:   admin.Name,
		Method:      "BindRoleAPI",
		Description: "绑定角色API",
		OwnerID:     admin.OwnerID,
		IP:          ctx.GetAddress(),
		URL:         ctx.GetURL(),
		Time:        time.Now().Unix(),
	}
	if e := s.mysql.GetWriteEngine().Create(&tbLog).Error; e != nil {
		log.Errorln(e.Error())
	}
	s.cache.GetCache().Del(fmt.Sprintf("%s-%d", define.Organize, request.RoleID))
	response.Success = true
	return
}

//DelRoleAPI 删除角色API
func (s *Service) DelRoleAPI(ctx plugins.Context, request param.DelRoleAPIReq) (response param.DelRoleAPIRes, err error) {
	admin, e := s.GetAdmin(ctx)
	if e != nil {
		log.Errorln(e.Error())
		err = customerror.EnCodeError(define.GetAdminInfoErr, e.Error())
		return
	}
	if _, e := authRPC.DelRoleAPI(ctx, define.Organize, request.RoleID, request.APIID); e != nil {
		log.Errorln(e.Error())
		err = customerror.EnCodeError(define.DelRoleAPIErr, "删除角色API失败")
		return
	}
	//生成登录记录
	tbLog := &table.Log{
		LogID:       s.snowflake.GetID(),
		Type:        table.DelRoleAPIType,
		AdminID:     admin.AdminID,
		AdminName:   admin.Name,
		Method:      "DelRoleAPI",
		Description: "删除角色API",
		OwnerID:     admin.OwnerID,
		IP:          ctx.GetAddress(),
		URL:         ctx.GetURL(),
		Time:        time.Now().Unix(),
	}
	if e := s.mysql.GetWriteEngine().Create(&tbLog).Error; e != nil {
		log.Errorln(e.Error())
	}
	s.cache.GetCache().Del(fmt.Sprintf("%s-%d", define.Organize, request.RoleID))
	response.Success = true
	return
}
