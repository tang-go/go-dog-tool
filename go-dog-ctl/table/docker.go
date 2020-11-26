package table

//Docker 记录
type Docker struct {
	//唯一主键
	ID int64
	//名称
	Name string
	//编译发布的管理员
	AdminID int64
	//发布镜像
	Image string
	//账号
	Account string
	//密码
	Pwd string
	//端口
	Ports string
	//业主ID
	OwnerID int64
	//注册事件
	Time int64
}

//Image image账号密码
type Image struct {
	ID      uint   `gorm:"primary_key" json:"id" description:"ID" type:"uint"`
	Image   string `json:"image" description:"镜像" type:"string"`
	Account string `json:"content" description:"账号" type:"string"`
	Pwd     string `json:"pwd" description:"密码" type:"string"`
	OwnerID int64  `json:"ownerID" description:"业主ID" type:"string"`
	Time    int64  `json:"time" description:"时间" type:"int64"`
}
