package param

import "github.com/tang-go/go-dog-tool/go-dog-auth/table"

//GetRoleMenuReq 获取角色菜单请求
type GetRoleMenuReq struct {
	Organize string `json:"organize" description:"组织" type:"string"`
	RoleID   uint   `json:"roleID" description:"角色ID" type:"string"`
}

//GetRoleMenuRes 获取角色菜单响应
type GetRoleMenuRes struct {
	RoleMenus []RoleMenu `json:"roleMenus" description:"菜单列表" type:"[]RoleMenu"`
}

//RoleMenu 权限菜单表
type RoleMenu struct {
	ID       uint   `json:"id" description:"菜单ID" type:"uint"`
	Describe string `json:"describe" description:"描述" type:"string"`
	URL      string `json:"url" description:"菜单URL" type:"string"`
	ParentID uint   `json:"parentID" description:"父菜单ID" type:"string"`
	Add      bool   `json:"add" description:"增加权限" type:"string"`
	Del      bool   `json:"del" description:"删除权限" type:"string"`
	Update   bool   `json:"update" description:"更新权限" type:"string"`
	Select   bool   `json:"select" description:"查询权限" type:"string"`
	Sort     uint   `json:"sort" description:"排序" type:"uint"`
	Time     int64  `json:"time" description:"时间" type:"int64"`
}

//GetRoleAPIReq 获取角色api请求
type GetRoleAPIReq struct {
	Organize string `json:"organize" description:"组织" type:"string"`
	RoleID   uint   `json:"roleID" description:"角色ID" type:"string"`
}

//GetRoleAPIRes 获取角色api响应
type GetRoleAPIRes struct {
	SysAPI []table.SysAPI `json:"sysAPI" description:"api列表" type:"[]table.SysAPI"`
}

//CreateRoleReq 创建普通角色请求
type CreateRoleReq struct {
	Organize string `json:"organize" description:"组织" type:"string"`
	Name     string `json:"name" description:"角色名称" type:"string"`
	Describe string `json:"describe" description:"描述" type:"string"`
	IsSuper  bool   `json:"isSuper" description:"是否是超级管理员" type:"bool"`
}

//CreateRoleRes 创建普通角色响应
type CreateRoleRes struct {
	ID uint `json:"id" description:"ID" type:"uint"`
}

//DelRoleReq 删除角色
type DelRoleReq struct {
	Organize string `json:"organize" description:"组织" type:"string"`
	ID       uint   `json:"id" description:"ID" type:"uint"`
}

//DelRoleRes 删除角色
type DelRoleRes struct {
	Success bool `json:"success" description:"结果" type:"bool"`
}

//SelectRoleByOrganizeReq 通过组织查询角色
type SelectRoleByOrganizeReq struct {
	Organize string `json:"organize" description:"组织" type:"string"`
}

//SelectRoleByOrganizeRes 通过组织查询角色
type SelectRoleByOrganizeRes struct {
	SysRoles []table.SysRole `json:"sysRoles" description:"角色列表" type:" []table.SysRole "`
}

//SelectRoleByIDReq 通过ID查询角色
type SelectRoleByIDReq struct {
	Organize string `json:"organize" description:"组织" type:"string"`
	RoleID   uint   `json:"RoleID" description:"角色ID" type:"uint"`
}

//SelectRoleByIDRes 通过ID查询角色
type SelectRoleByIDRes struct {
	SysRole table.SysRole `json:"sysRole" description:"系统角色" type:" table.SysRole"`
}

//SelectMenuReq 获取菜单
type SelectMenuReq struct {
	Organize string `json:"organize" description:"组织" type:"string"`
}

//SelectMenuRes 获取菜单
type SelectMenuRes struct {
	Menus []table.SysMenu `json:"menus" description:"菜单" type:"[]table.SysMenu"`
}

//CreateMenuReq 创建菜单请求
type CreateMenuReq struct {
	Organize string `json:"organize" description:"组织" type:"string"`
	Describe string `json:"describe" description:"描述" type:"string"`
	URL      string `json:"url" description:"菜单URL" type:"string"`
	ParentID uint   `json:"parentID" description:"父菜单ID" type:"string"`
	Sort     uint   `json:"sort" description:"排序" type:"uint"`
}

//CreateMenuRes 创建菜单
type CreateMenuRes struct {
	ID uint `json:"id" description:"ID" type:"uint"`
}

//DelMenuReq 删除菜单
type DelMenuReq struct {
	Organize string `json:"organize" description:"组织" type:"string"`
	ID       uint   `json:"id" description:"ID" type:"uint"`
}

//DelMenuRes 删除菜单
type DelMenuRes struct {
	Success bool `json:"success" description:"结果" type:"bool"`
}

//SelectAPIReq 获取API
type SelectAPIReq struct {
	Organize string `json:"organize" description:"组织" type:"string"`
}

//SelectAPIRes 获取API
type SelectAPIRes struct {
	APIS []table.SysAPI `json:"apis" description:"菜单" type:"[]table.SysAPI"`
}

//DelAPIReq 删除API
type DelAPIReq struct {
	Organize string `json:"organize" description:"组织" type:"string"`
	ID       uint   `json:"id" description:"ID" type:"uint"`
}

//DelAPIRes 删除API
type DelAPIRes struct {
	Success bool `json:"success" description:"结果" type:"bool"`
}

//CreateAPIReq 创建api请求
type CreateAPIReq struct {
	Organize string `json:"organize" description:"组织" type:"string"`
	Describe string `json:"describe" description:"描述" type:"string"`
	API      string `json:"api" description:"api接口" type:"string"`
}

//CreateAPIRes 创建api请求
type CreateAPIRes struct {
	ID uint `json:"id" description:"ID" type:"uint"`
}

//BindRoleMenuReq 绑定角色菜单
type BindRoleMenuReq struct {
	Organize string `json:"organize" description:"组织" type:"string"`
	MenuID   uint   `json:"menuID" description:"系统菜单ID" type:"string"`
	RoleID   uint   `json:"roleID" description:"角色ID" type:"string"`
	Add      bool   `json:"add" description:"增加权限" type:"string"`
	Del      bool   `json:"del" description:"删除权限" type:"string"`
	Update   bool   `json:"update" description:"更新权限" type:"string"`
	Select   bool   `json:"select" description:"查询权限" type:"string"`
}

//BindRoleMenuRes 绑定角色菜单
type BindRoleMenuRes struct {
	Success bool `json:"success" description:"结果" type:"bool"`
}

//DelRoleMenuReq 删除角色菜单
type DelRoleMenuReq struct {
	Organize string `json:"organize" description:"组织" type:"string"`
	MenuID   uint   `json:"menuID" description:"系统菜单ID" type:"string"`
	RoleID   uint   `json:"roleID" description:"角色ID" type:"string"`
}

//DelRoleRes 删除角色菜单
type DelRoleMenuRes struct {
	Success bool `json:"success" description:"结果" type:"bool"`
}

//BindRoleAPIReq 绑定api菜单
type BindRoleAPIReq struct {
	Organize string `json:"organize" description:"组织" type:"string"`
	RoleID   uint   `json:"roleID" description:"角色ID" type:"string"`
	APIID    uint   `json:"apiID" description:"api的ID" type:"string"`
}

//BindRoleAPIRes 绑定api菜单
type BindRoleAPIRes struct {
	Success bool `json:"success" description:"结果" type:"bool"`
}

//DelRoleAPIReq 删除角色API
type DelRoleAPIReq struct {
	Organize string `json:"organize" description:"组织" type:"string"`
	RoleID   uint   `json:"roleID" description:"角色ID" type:"string"`
	APIID    uint   `json:"apiID" description:"api的ID" type:"string"`
}

//DelRoleAPIRes 删除角色API
type DelRoleAPIRes struct {
	Success bool `json:"success" description:"结果" type:"bool"`
}
