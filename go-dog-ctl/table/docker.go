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
	//业主ID
	OwnerID int64
	//注册事件
	Time int64
}
