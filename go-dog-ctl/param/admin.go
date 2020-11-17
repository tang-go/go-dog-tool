package param

//GetAdminInfoReq 获取管理员信息
type GetAdminInfoReq struct {
}

//GetAdminInfoRes 获取用户信息返回
type GetAdminInfoRes struct {
	ID       string      `json:"id" description:"用户ID" type:"string"`
	Name     string      `json:"name" description:"名称" type:"string"`
	Avatar   string      `json:"avatar" description:"头像" type:"string"`
	Phone    string      `json:"phone" description:"电话" type:"string"`
	RoleID   uint        `json:"roleId" description:"权限的ID" type:"uint"`
	RoleName string      `json:"roleName" description:"权限名称" type:"string"`
	Menu     []*RoleMenu `json:"menu" description:"菜单" type:"[]*RoleMenu"`
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

//AdminLoginReq 管理员登录
type AdminLoginReq struct {
	Phone string `json:"phone" description:"电话" type:"string"`
	Pwd   string `json:"pwd" description:"密码" type:"string"`
}

//AdminLoginRes 管理员登录返回
type AdminLoginRes struct {
	Name    string `json:"name" description:"名称" type:"string"`
	OwnerID string `json:"ownerId" description:"业主ID" type:"string"`
	Token   string `json:"token" description:"注册用户的token" type:"string"`
}

//AdminOnlineReq 管理员上线
type AdminOnlineReq struct {
	Address string `json:"address" description:"用户上线" type:"string"`
}

//AdminOnlineRes 管理员上线
type AdminOnlineRes struct {
	Success bool `success:"address" description:"结果" type:"bool"`
}

//AdminOfflineReq 管理员下线
type AdminOfflineReq struct {
	Address string `json:"address" description:"用户上线" type:"string"`
}

//AdminOfflineRes 管理员下线
type AdminOfflineRes struct {
	Success bool `success:"address" description:"结果" type:"bool"`
}
