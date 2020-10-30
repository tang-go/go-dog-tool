package param

//StartDockerReq 启动Docker请求
type StartDockerReq struct {
	Name    string   `json:"name" description:"服务名称" type:"string"`
	Images  string   `json:"images" description:"镜像地址" type:"string"`
	Account string   `json:"account" description:"账号" type:"string"`
	Pwd     string   `json:"pwd" description:"密码" type:"string"`
	Ports   []*Ports `json:"ports" description:"服务暴露端口" type:"[]*Ports"`
}

//Ports 暴露端口
type Ports struct {
	InsidePort   string `json:"insidePort" description:"内部端口" type:"string"`
	ExternalPort string `json:"externalPort" description:"外部暴露端口" type:"string"`
}

//StartDockerRes 启动Docker响应
type StartDockerRes struct {
	Success bool `json:"success" description:"结果" type:"bool"`
}

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
