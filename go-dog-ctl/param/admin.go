package param

//CreateAdminReq 创建管理员
type CreateAdminReq struct {
	Name   string `json:"name" description:"名称" type:"string"`
	Phone  string `json:"phone" description:"电话" type:"string"`
	Pwd    string `json:"pwd" description:"密码" type:"string"`
	RoleID uint   `json:"roleID" description:"角色ID" type:"uint"`
}

//CreateAdminRes 创建管理员
type CreateAdminRes struct {
	Success bool `success:"address" description:"结果" type:"bool"`
}

//DelAdminReq 删除管理员
type DelAdminReq struct {
	AdminID string `json:"adminID" description:"账号 唯一主键" type:"string"`
}

//DelAdminRes 删除管理员
type DelAdminRes struct {
	Success bool `success:"address" description:"结果" type:"bool"`
}

//DisableAdminReq 禁用管理员
type DisableAdminReq struct {
	AdminID string `json:"adminID" description:"账号 唯一主键" type:"string"`
}

//DisableAdminRes 禁用管理员
type DisableAdminRes struct {
	Success bool `success:"address" description:"结果" type:"bool"`
}

//OpenAdminReq 开启管理员
type OpenAdminReq struct {
	AdminID string `json:"adminID" description:"账号 唯一主键" type:"string"`
}

//OpenAdminRes 开启管理员
type OpenAdminRes struct {
	Success bool `success:"address" description:"结果" type:"bool"`
}

//GetAdminListReq 获取管理员列表
type GetAdminListReq struct {
}

// GetAdminListRes 获取管理员列表
type GetAdminListRes struct {
	AdminInfos []AdminInfo `json:"adminInfos" description:"管理员列表" type:"[]AdminInfo"`
}

//AdminInfo 管理员信息
type AdminInfo struct {
	AdminID   string      `json:"adminID" description:"账号 唯一主键" type:"string"`
	Name      string      `json:"name" description:"名称" type:"string"`
	Phone     string      `json:"phone" description:"电话" type:"string"`
	IsDisable bool        `json:"isDisable" description:"是否被禁用" type:"bool"`
	RoleID    uint        `json:"roleID" description:"名称" type:"uint"`
	RoleName  string      `json:"roleName" description:"权限名称" type:"string"`
	Menu      []*RoleMenu `json:"menu" description:"菜单" type:"[]*RoleMenu"`
	APIS      []API       `json:"apis" description:"API" type:"[]API"`
	Time      string      `json:"time" description:"时间" type:"string"`
}

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
	APIS     []API       `json:"apis" description:"API" type:"[]API"`
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
