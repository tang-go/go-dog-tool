package param

import "omo-service/services/auth/table"

//GetRoleMenuReq 获取角色菜单请求
type GetRoleMenuReq struct {
	RoleID uint `json:"roleID" description:"角色ID" type:"string"`
}

//GetRoleMenuRes 获取角色菜单响应
type GetRoleMenuRes struct {
	SysMenu []SysMenu `json:"SysMenu" description:"菜单列表" type:"[]SysMenu"`
}

//SysMenu 系统菜单表
type SysMenu struct {
	ID       uint   `gorm:"primary_key" json:"id" description:"ID" type:"uint"`
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

//GetRoleApiReq 获取角色api请求
type GetRoleApiReq struct {
	RoleID uint `json:"roleID" description:"角色ID" type:"string"`
}

//GetRoleApiRes 获取角色api响应
type GetRoleApiRes struct {
	SysApi []table.SysApi `json:"sysApi" description:"api列表" type:"[]table.SysApi"`
}

//CreateRoleReq 创建普通角色请求
type CreateRoleReq struct {
	Organize string `json:"organize" description:"组织" type:"string"`
	Name     string `json:"name" description:"角色名称" type:"string"`
	Describe string `json:"describe" description:"描述" type:"string"`
	IsSuper  bool   `json:"isSuper" description:"是否是超级管理员" type:"isSuper"`
}

//CreateRoleRes 创建普通角色响应
type CreateRoleRes struct {
	ID uint `json:"id" description:"ID" type:"uint"`
}

//SelectRoleReq 通过组织查询角色
type SelectRoleReq struct {
	Organize string `json:"organize" description:"组织" type:"string"`
}

//CreateRoleRes 通过组织查询角色
type SelectRoleRes struct {
	SysRoles []table.SysRole `json:"id" description:"角色列表" type:" []table.SysRole "`
}

//CreateMenuReq 创建菜单请求
type CreateMenuReq struct {
	Describe string `json:"describe" description:"描述" type:"string"`
	URL      string `json:"url" description:"菜单URL" type:"string"`
	ParentID uint   `json:"parentID" description:"父菜单ID" type:"string"`
	Sort     uint   `json:"sort" description:"排序" type:"uint"`
}

//CreateMenuRes 创建菜单请求
type CreateMenuRes struct {
	ID uint `json:"id" description:"ID" type:"uint"`
}

//CreateApiReq 创建api请求
type CreateApiReq struct {
	Describe string `json:"describe" description:"描述" type:"string"`
	API      string `json:"api" description:"api接口" type:"string"`
}

//CreateApiRes 创建api请求
type CreateApiRes struct {
	ID uint `json:"id" description:"ID" type:"uint"`
}

//BindRoleMenuReq 绑定角色菜单
type BindRoleMenuReq struct {
	MenuID uint `json:"menuID" description:"系统菜单ID" type:"string"`
	RoleID uint `json:"roleID" description:"角色ID" type:"string"`
	Add    bool `json:"add" description:"增加权限" type:"string"`
	Del    bool `json:"del" description:"删除权限" type:"string"`
	Update bool `json:"update" description:"更新权限" type:"string"`
	Select bool `json:"select" description:"查询权限" type:"string"`
}

//BindRoleMenuRes 绑定角色菜单
type BindRoleMenuRes struct {
	Success bool `json:"success" description:"结果" type:"bool"`
}

//BindRoleApiReq 绑定api菜单
type BindRoleApiReq struct {
	RoleID uint `json:"roleID" description:"角色ID" type:"string"`
	ApiID  uint `json:"apiID" description:"api的ID" type:"string"`
}

//BindRoleApiRes 绑定api菜单
type BindRoleApiRes struct {
	Success bool `json:"success" description:"结果" type:"bool"`
}
