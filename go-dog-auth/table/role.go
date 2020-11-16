package table

//SysRole 系统角色表
type SysRole struct {
	ID       uint   `gorm:"primary_key" json:"id" description:"ID" type:"uint"`
	Organize string `json:"organize" description:"组织" type:"string"`
	Name     string `json:"name" description:"角色名称" type:"string"`
	Describe string `json:"describe" description:"描述" type:"string"`
	IsSuper  bool   `json:"isSuper" description:"是否是超级管理员" type:"isSuper"`
	Time     int64  `json:"time" description:"时间" type:"int64"`
}

//SysRoleMenu 系统角色菜单映射表
type SysRoleMenu struct {
	ID     uint  `gorm:"primary_key" json:"id" description:"ID" type:"uint"`
	MenuID uint  `json:"menuID" description:"系统菜单ID" type:"string"`
	RoleID uint  `json:"roleID" description:"角色ID" type:"string"`
	Add    bool  `json:"add" description:"增加权限" type:"string"`
	Del    bool  `json:"del" description:"删除权限" type:"string"`
	Update bool  `json:"update" description:"更新权限" type:"string"`
	Select bool  `json:"select" description:"查询权限" type:"string"`
	Time   int64 `json:"time" description:"时间" type:"int64"`
}

//SysRoleApi 系统角色API映射表
type SysRoleApi struct {
	ID     uint  `gorm:"primary_key" json:"id" description:"ID" type:"uint"`
	RoleID uint  `json:"roleID" description:"角色ID" type:"string"`
	ApiID  uint  `json:"apiID" description:"api的ID" type:"string"`
	Time   int64 `json:"time" description:"时间" type:"int64"`
}

//SysMenu 系统菜单表
type SysMenu struct {
	ID       uint   `gorm:"primary_key" json:"id" description:"ID" type:"uint"`
	Organize string `json:"organize" description:"组织" type:"string"`
	Describe string `json:"describe" description:"描述" type:"string"`
	URL      string `json:"url" description:"菜单URL" type:"string"`
	ParentID uint   `json:"parentID" description:"父菜单ID" type:"string"`
	Sort     uint   `json:"sort" description:"排序" type:"uint"`
	Time     int64  `json:"time" description:"时间" type:"int64"`
}

//SysApi 系统API访问权限
type SysApi struct {
	ID       uint   `gorm:"primary_key" json:"id" description:"ID" type:"uint"`
	Organize string `json:"organize" description:"组织" type:"string"`
	Describe string `json:"describe" description:"描述" type:"string"`
	API      string `json:"api" description:"api接口" type:"string"`
	Time     int64  `json:"time" description:"时间" type:"int64"`
}
