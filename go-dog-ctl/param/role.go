package param

//GetRoleAPIReq 获取角色api请求
type GetRoleAPIReq struct {
	RoleID uint `json:"roleID" description:"角色ID" type:"string"`
}

//GetRoleAPIRes 获取角色api响应
type GetRoleAPIRes struct {
	APIS []API `json:"apis" description:"API" type:"[]API"`
}

//GetAPIListReq 获取API
type GetAPIListReq struct {
}

//GetAPIListRes 获取API
type GetAPIListRes struct {
	APIS []API `json:"apis" description:"API" type:"[]API"`
}

//API API详情
type API struct {
	ID       uint   `gorm:"primary_key" json:"id" description:"ID" type:"uint"`
	Organize string `json:"organize" description:"组织" type:"string"`
	Describe string `json:"describe" description:"描述" type:"string"`
	API      string `json:"api" description:"api接口" type:"string"`
	Time     string `json:"time" description:"时间" type:"string"`
}

//DelAPIReq 删除API
type DelAPIReq struct {
	ID uint `json:"id" description:"ID" type:"uint"`
}

//DelAPIRes 删除API
type DelAPIRes struct {
	Success bool `json:"success" description:"结果" type:"bool"`
}

//BindRoleAPIReq 绑定api菜单
type BindRoleAPIReq struct {
	RoleID uint `json:"roleID" description:"角色ID" type:"string"`
	APIID  uint `json:"apiID" description:"api的ID" type:"string"`
}

//BindRoleAPIRes 绑定api菜单
type BindRoleAPIRes struct {
	Success bool `json:"success" description:"结果" type:"bool"`
}

//DelRoleAPIReq 删除角色API
type DelRoleAPIReq struct {
	RoleID uint `json:"roleID" description:"角色ID" type:"string"`
	APIID  uint `json:"apiID" description:"api的ID" type:"string"`
}

//DelRoleAPIRes 删除角色API
type DelRoleAPIRes struct {
	Success bool `json:"success" description:"结果" type:"bool"`
}

//GetRoleMenuReq 获取角色菜单
type GetRoleMenuReq struct {
	ID uint `json:"id" description:"ID" type:"uint"`
}

//GetRoleMenuRes 获取角色菜单
type GetRoleMenuRes struct {
	Menu []*RoleMenu `json:"menu" description:"菜单" type:"[]*RoleMenu"`
}

//RoleMenu 权限菜单
type RoleMenu struct {
	ID       uint        `json:"id" description:"ID" type:"uint"`
	ParentID uint        `json:"parentID" description:"父亲结点ID" type:"uint"`
	URL      string      `json:"url" description:"菜单路由" type:"string"`
	Describe string      `json:"describe" description:"菜单描述" type:"string"`
	Add      bool        `json:"add" description:"增加权限" type:"string"`
	Del      bool        `json:"del" description:"删除权限" type:"string"`
	Update   bool        `json:"update" description:"更新权限" type:"string"`
	Select   bool        `json:"select" description:"查询权限" type:"string"`
	Sort     uint        `json:"sort" description:"排序" type:"uint"`
	Time     int64       `json:"time" description:"时间" type:"int64"`
	Children []*RoleMenu `json:"children" json:"id" description:"子菜单" type:"[]*RoleMenu"`
}

//BindRoleMenuReq 绑定角色菜单
type BindRoleMenuReq struct {
	RoleID uint `json:"roleID" description:"角色id" type:"uint"`
	MenuID uint `json:"menuID" description:"菜单id" type:"uint"`
	Add    bool `json:"add" description:"增加权限" type:"string"`
	Del    bool `json:"del" description:"删除权限" type:"string"`
	Update bool `json:"update" description:"更新权限" type:"string"`
	Select bool `json:"select" description:"查询权限" type:"string"`
}

//BindRoleMenuRes 绑定角色菜单
type BindRoleMenuRes struct {
	Success bool `json:"success" description:"结果" type:"bool"`
}

//DelRoleMenuReq 删除角色菜单
type DelRoleMenuReq struct {
	RoleID uint `json:"roleID" description:"角色id" type:"uint"`
	MenuID uint `json:"menuID" description:"菜单id" type:"uint"`
}

//DelRoleMenuRes 删除角色菜单
type DelRoleMenuRes struct {
	Success bool `json:"success" description:"结果" type:"bool"`
}

//GetRoleReq 获取角色列表
type GetRoleListReq struct {
}

//GetRoleListRes 获取角色列表
type GetRoleListRes struct {
	Roles []Role `json:"role" description:"角色列表" type:"[]Role"`
}

//Role 角色
type Role struct {
	ID       uint   `json:"id" description:"ID" type:"uint"`
	Organize string `json:"organize" description:"组织" type:"string"`
	Name     string `json:"name" description:"角色名称" type:"string"`
	Describe string `json:"describe" description:"描述" type:"string"`
	IsSuper  bool   `json:"isSuper" description:"是否是超级管理员" type:"bool"`
	Time     string `json:"time" description:"时间" type:"string"`
}

//CreateRoleReq 创建角色
type CreateRoleReq struct {
	Name     string `json:"name" description:"角色名称" type:"string"`
	Describe string `json:"describe" description:"描述" type:"string"`
}

//CreateRoleRes 创建角色
type CreateRoleRes struct {
	ID uint `json:"id" description:"ID" type:"uint"`
}

//DelRoleReq 删除角色
type DelRoleReq struct {
	ID uint `json:"id" description:"ID" type:"uint"`
}

//DelRoleRes 删除角色
type DelRoleRes struct {
	Success bool `json:"success" description:"结果" type:"bool"`
}

//CreateMenuReq 创建菜单请求
type CreateMenuReq struct {
	Describe string `json:"describe" description:"描述" type:"string"`
	URL      string `json:"url" description:"菜单URL" type:"string"`
	ParentID uint   `json:"parentID" description:"父菜单ID" type:"uint"`
	Sort     uint   `json:"sort" description:"排序" type:"uint"`
}

//CreateMenuRes 创建菜单响应
type CreateMenuRes struct {
	Success bool `json:"success" description:"结果" type:"bool"`
}

//DelMenuReq 删除菜单请求
type DelMenuReq struct {
	MenuID uint `json:"menuID" description:"菜单ID" type:"uint"`
}

//DelMenuRes 删除菜单响应
type DelMenuRes struct {
	Success bool `json:"success" description:"结果" type:"bool"`
}

//GetMenuReq 获取菜单请求
type GetMenuReq struct {
}

//GetMenuRes 获取菜单响应
type GetMenuRes struct {
	Menu []*Menu `json:"menu" description:"菜单" type:"Role"`
}

//Menu 菜单配置
type Menu struct {
	ID       uint    `json:"id" description:"ID" type:"uint"`
	ParentID uint    `json:"parentID" description:"父亲结点ID" type:"uint"`
	Organize string  `json:"organize" description:"组织" type:"string"`
	Describe string  `json:"describe" description:"描述" type:"string"`
	URL      string  `json:"url" description:"菜单URL" type:"string"`
	Sort     uint    `json:"sort" description:"排序" type:"uint"`
	Children []*Menu `json:"children" description:"菜单" type:"[]*Menu"`
	Time     string  `json:"string" description:"时间" type:"int64"`
}
