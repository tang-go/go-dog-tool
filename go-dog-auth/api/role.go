package api

import (
	"github.com/tang-go/go-dog-tool/define"
	"github.com/tang-go/go-dog-tool/go-dog-auth/param"
	"github.com/tang-go/go-dog-tool/go-dog-auth/table"
	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/plugins"
)

//GetRoleMenu 获取角色菜单
func GetRoleMenu(ctx plugins.Context, organize string, roleID uint) (sysMenus []param.RoleMenu, err error) {
	info := param.GetRoleMenuRes{}
	if err := ctx.GetClient().Call(ctx, plugins.RandomMode, define.SvcAuth, "GetRoleMenu", param.GetRoleMenuReq{
		Organize: organize,
		RoleID:   roleID,
	}, &info); err != nil {
		log.Errorln(err.Error())
		return info.RoleMenus, err
	}
	return info.RoleMenus, nil
}

//GetRoleAPI 获取角色API
func GetRoleAPI(ctx plugins.Context, organize string, roleID uint) (sysAPIs []table.SysAPI, err error) {
	info := param.GetRoleAPIRes{}
	if err := ctx.GetClient().Call(ctx, plugins.RandomMode, define.SvcAuth, "GetRoleAPI", param.GetRoleAPIReq{
		Organize: organize,
		RoleID:   roleID,
	}, &info); err != nil {
		log.Errorln(err.Error())
		return info.SysAPI, err
	}
	return info.SysAPI, nil
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

//SelectMenu 查询菜单
func SelectMenu(ctx plugins.Context, organize string) (sysMenu []table.SysMenu, err error) {
	info := param.SelectMenuRes{}
	if err := ctx.GetClient().Call(ctx, plugins.RandomMode, define.SvcAuth, "SelectMenu", param.SelectMenuReq{
		Organize: organize,
	}, &info); err != nil {
		log.Errorln(err.Error())
		return info.Menus, err
	}
	return info.Menus, nil
}

//CreateAPI 创建api
func CreateAPI(ctx plugins.Context, organize string, describe string, api string) (id uint, err error) {
	info := param.CreateAPIRes{}
	if err := ctx.GetClient().Call(ctx, plugins.RandomMode, define.SvcAuth, "CreateAPI", param.CreateAPIReq{
		Organize: organize,
		Describe: describe,
		API:      api,
	}, &info); err != nil {
		log.Errorln(err.Error())
		return info.ID, err
	}
	return info.ID, nil
}

//SelectAPI 查询API
func SelectAPI(ctx plugins.Context, organize string) (sysAPI []table.SysAPI, err error) {
	info := param.SelectAPIRes{}
	if err := ctx.GetClient().Call(ctx, plugins.RandomMode, define.SvcAuth, "SelectAPI", param.SelectAPIReq{
		Organize: organize,
	}, &info); err != nil {
		log.Errorln(err.Error())
		return info.APIS, err
	}
	return info.APIS, nil
}

//BindRoleAPI 绑定角色api
func BindRoleAPI(ctx plugins.Context, roleID uint, apiID uint) (success bool, err error) {
	info := param.BindRoleAPIRes{}
	if err := ctx.GetClient().Call(ctx, plugins.RandomMode, define.SvcAuth, "BindRoleAPI", param.BindRoleAPIReq{
		RoleID: roleID,
		APIID:  apiID,
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
