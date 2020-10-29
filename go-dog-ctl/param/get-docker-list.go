package param

//GetDockerListReq 获取docker列表请求
type GetDockerListReq struct {
}

//GetDockerListRes 获取docker列表响应
type GetDockerListRes struct {
	Data       []Docker `json:"data" description:"编译发布记录" type:"[]BuildService"`
	PageSize   int      `json:"pageSize" description:"一页大小" type:"int"`
	PageNo     int      `json:"pageNo" description:"当前所处的页数" type:"int"`
	TotalPage  int      `json:"totalPage" description:"总页数" type:"int"`
	TotalCount int      `json:"totalCount" description:"总数量" type:"int"`
}

//Docker 发布记录
type Docker struct {
	//唯一主键
	ID string `json:"id" description:"id" type:"string"`
	//名称
	Name string `json:"name" description:"名称" type:"string"`
	//编译发布的管理员
	AdminID int64 `json:"adminId" description:"管理员ID" type:"int64"`
	//发布镜像
	Image string `json:"image" description:"发布镜像" type:"string"`
	//状态
	Status string `json:"status" description:"状态" type:"string"`
	//业主ID
	OwnerID int64 `json:"ownerId" description:"业主ID" type:"int64"`
	//执行命令
	Command string `json:"command" description:"执行命令" type:"string"`
	//端口
	Ports string `json:"ports" description:"端口" type:"string"`
	//注册事件
	Time string `json:"time" description:"时间" type:"string"`
}
