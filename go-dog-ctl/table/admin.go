package table

const (
	//AdminDisable  管理员禁用
	AdminDisable = true
	//AdminAvailable 管理员可用
	AdminAvailable = false
)

//Admin 超级管理员
type Admin struct {
	AdminID     int64  `json:"adminID" description:"账号 唯一主键" type:"int64"`
	Name        string `json:"name" description:"名称" type:"string"`
	Phone       string `json:"phone" description:"电话" type:"string"`
	Pwd         string `json:"pwd" description:"密码" type:"string"`
	Salt        string `json:"salt" description:"盐值 md5使用" type:"string"`
	OwnerID     int64  `json:"ownerID" description:"所属业主" type:"int64"`
	IsDisable   bool   `json:"isDisable" description:"是否被禁用" type:"bool"`
	IsOnline    bool   `json:"isOnline" description:"是否在线" type:"bool"`
	GateAddress string `json:"gateAddress" description:"在线网关地址" type:"string"`
	RoleID      uint   `json:"roleID" description:"角色" type:"uint"`
	Time        int64  `json:"time" description:"角色" type:"int64"`
}
