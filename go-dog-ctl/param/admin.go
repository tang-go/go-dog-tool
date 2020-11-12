package param

//GetAdminInfoReq 获取管理员信息
type GetAdminInfoReq struct {
}

//GetAdminInfoRes 获取用户信息返回
type GetAdminInfoRes struct {
	ID     string `json:"id" description:"用户ID" type:"string"`
	Name   string `json:"name" description:"名称" type:"string"`
	Avatar string `json:"avatar" description:"头像" type:"string"`
	Phone  string `json:"phone" description:"电话" type:"string"`
	RoleID string `json:"roleId" description:"权限名称" type:"string"`
	Role   Role   `json:"role" description:"权限内容" type:"Role"`
}

//Role 权限配置
type Role struct {
	ID          string         `json:"id" description:"权限名称" type:"string"`
	Permissions []*Permissions `json:"permissions" json:"id" description:"权限内容" type:"[]*Permissions "`
}

//Permissions 权限
type Permissions struct {
	RoleID          string             `json:"roleId" description:"权限名称" type:"string"`
	PermissionID    string             `json:"permissionId" description:"操作名称" type:"string"`
	PermissionName  string             `json:"permissionName" description:"操作描述" type:"string"`
	ActionEntitySet []*ActionEntitySet `json:"actionEntitySet" description:"操作动作" type:"[]*ActionEntitySet "`
}

//ActionEntitySet 动作
type ActionEntitySet struct {
	Action       string `json:"action" description:"动作名称" type:"string"`
	Describe     string `json:"describe" description:"动作描述" type:"string"`
	DefaultCheck bool   `json:"defaultCheck" description:"是否可以操作" type:"bool"`
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
