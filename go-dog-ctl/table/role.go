package table

const (
	//IsAdmin 是超级管理员
	IsAdmin = true
	//NoAdmin 不是管理员
	NoAdmin = false
)

//OwnerRole 业主角色表
type OwnerRole struct {
	ID          uint   `gorm:"primary_key" json:"id" description:"ID" type:"uint"`
	RoleID      uint   `json:"roleID" description:"角色ID" type:"int64"`
	Name        string `json:"name" description:"角色名称" type:"string"`
	Description string `json:"description" description:"角色描述" type:"string"`
	OwnerID     int64  `json:"ownerID" description:"业主ID" type:"int64"`
	Time        int64  `json:"time" description:"角色创建时间" type:"int64"`
}
