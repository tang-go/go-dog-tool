package api

import (
	"github.com/tang-go/go-dog-tool/define"
	"github.com/tang-go/go-dog-tool/go-dog-auth/param"
	"github.com/tang-go/go-dog-tool/go-dog-auth/table"
	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/plugins"
)

//GetRoleMenu 获取角色菜单
func GetRoleMenu(ctx plugins.Context, organize string, roleID uint) (sysMenus []param.SysMenu, err error) {
	info := param.GetRoleMenuRes{}
	if err := ctx.GetClient().Call(ctx, plugins.RandomMode, define.SvcAuth, "GetRoleMenu", param.GetRoleMenuReq{
		Organize: organize,
		RoleID:   roleID,
	}, &info); err != nil {
		log.Errorln(err.Error())
		return info.SysMenu, err
	}
	return info.SysMenu, nil
}

//GetRoleApi 获取角色Api
func GetRoleApi(ctx plugins.Context, organize string, roleID uint) (sysApis []table.SysApi, err error) {
	info := param.GetRoleApiRes{}
	if err := ctx.GetClient().Call(ctx, plugins.RandomMode, define.SvcAuth, "GetRoleApi", param.GetRoleApiReq{
		Organize: organize,
		RoleID:   roleID,
	}, &info); err != nil {
		log.Errorln(err.Error())
		return info.SysApi, err
	}
	return info.SysApi, nil
}

//CreateRole 创建角色
func CreateRole(ctx plugins.Context, organize, name, describe string, isSuper bool) (id uint, err error) {
	info := param.CreateRoleRes{}
	if err := ctx.GetClient().Call(ctx, plugins.RandomMode, define.SvcAuth, "CreateRole", param.CreateRoleReq{
		Name:     name,
		Describe: describe,
		IsSuper:  isSuper,
		Organize: organize,
	}, &info); err != nil {
		log.Errorln(err.Error())
		return info.ID, err
	}
	return info.ID, nil
}

//SelectRoleByOrganize 查询角色
func SelectRoleByOrganize(ctx plugins.Context, organize string) (sysRoles []table.SysRole, err error) {
	info := param.SelectRoleByOrganizeRes{}
	if err := ctx.GetClient().Call(ctx, plugins.RandomMode, define.SvcAuth, "SelectRoleByOrganize", param.SelectRoleByOrganizeReq{
		Organize: organize,
	}, &info); err != nil {
		log.Errorln(err.Error())
		return info.SysRoles, err
	}
	return info.SysRoles, nil
}

//SelectRoleByID 查询角色
func SelectRoleByID(ctx plugins.Context, organize string, roleID uint) (sysRole table.SysRole, err error) {
	info := param.SelectRoleByIDRes{}
	if err := ctx.GetClient().Call(ctx, plugins.RandomMode, define.SvcAuth, "SelectRoleByID", param.SelectRoleByIDReq{
		Organize: organize,
		RoleID:   roleID,
	}, &info); err != nil {
		log.Errorln(err.Error())
		return info.SysRole, err
	}
	return info.SysRole, nil
}

//CreateMenu 创建菜单
func CreateMenu(ctx plugins.Context, organize string, describe string, url string, parentID uint, sort uint) (id uint, err error) {
	info := param.CreateMenuRes{}
	if err := ctx.GetClient().Call(ctx, plugins.RandomMode, define.SvcAuth, "CreateMenu", param.CreateMenuReq{
		Organize: organize,
		Describe: describe,
		URL:      url,
		ParentID: parentID,
		Sort:     sort,
	}, &info); err != nil {
		log.Errorln(err.Error())
		return info.ID, err
	}
	return info.ID, nil
}

//CreateApi 创建api
func CreateApi(ctx plugins.Context, organize string, describe string, api string) (id uint, err error) {
	info := param.CreateApiRes{}
	if err := ctx.GetClient().Call(ctx, plugins.RandomMode, define.SvcAuth, "CreateApi", param.CreateApiReq{
		Organize: organize,
		Describe: describe,
		API:      api,
	}, &info); err != nil {
		log.Errorln(err.Error())
		return info.ID, err
	}
	return info.ID, nil
}

//BindRoleApi 绑定角色api
func BindRoleApi(ctx plugins.Context, roleID uint, apiID uint) (success bool, err error) {
	info := param.BindRoleApiRes{}
	if err := ctx.GetClient().Call(ctx, plugins.RandomMode, define.SvcAuth, "BindRoleApi", param.BindRoleApiReq{
		RoleID: roleID,
		ApiID:  apiID,
	}, &info); err != nil {
		log.Errorln(err.Error())
		return info.Success, err
	}
	return info.Success, nil
}

//BindRoleMenu 绑定角色菜单
func BindRoleMenu(ctx plugins.Context, menuID uint, roleID uint, add bool, del bool, update bool, sel bool) (success bool, err error) {
	info := param.BindRoleMenuRes{}
	if err := ctx.GetClient().Call(ctx, plugins.RandomMode, define.SvcAuth, "BindRoleMenu", param.BindRoleMenuReq{
		MenuID: menuID,
		RoleID: roleID,
		Add:    add,
		Del:    del,
		Update: update,
		Select: sel,
	}, &info); err != nil {
		log.Errorln(err.Error())
		return info.Success, err
	}
	return info.Success, nil
}
