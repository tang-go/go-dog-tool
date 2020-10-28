package table

//BuildService 发布记录
type BuildService struct {
	//唯一主键
	ID int64
	//编译发布的管理员
	AdminID int64
	//发布镜像
	Image string
	//状态
	Status bool
	//业主ID
	OwnerID int64
	//执行日志
	Log string
	//注册事件
	Time int64
}
