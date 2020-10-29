package param

//GetBuildServiceReq 获取编译发布请求
type GetBuildServiceReq struct {
	PageNo   int `json:"pageNo" description:"当前请求页数" type:"int"`
	PageSize int `json:"pageSize" description:"每一页的总数" type:"int"`
}

//GetBuildServiceRes 获取编译发布响应
type GetBuildServiceRes struct {
	Data       []BuildService `json:"data" description:"编译发布记录" type:"[]BuildService"`
	PageSize   int            `json:"pageSize" description:"一页大小" type:"int"`
	PageNo     int            `json:"pageNo" description:"当前所处的页数" type:"int"`
	TotalPage  int            `json:"totalPage" description:"总页数" type:"int"`
	TotalCount int            `json:"totalCount" description:"总数量" type:"int"`
}

//BuildService 发布记录
type BuildService struct {
	//唯一主键
	ID int64 `json:"id" description:"唯一主键" type:"int64"`
	//编译发布的管理员
	AdminID int64 `json:"adminId" description:"管理员ID" type:"int64"`
	//发布镜像
	Image string `json:"image" description:"发布镜像" type:"string"`
	//状态
	Status bool `json:"status" description:"状态" type:"bool"`
	//业主ID
	OwnerID int64 `json:"ownerId" description:"业主ID" type:"int64"`
	//执行日志
	Log string `json:"log" description:"日志" type:"string"`
	//注册事件
	Time string `json:"time" description:"时间" type:"string"`
}
