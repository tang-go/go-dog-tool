package param

//StartListDockerLogReq 开始监听docker日志
type StartListDockerLogReq struct {
	ID      string `json:"id" description:"镜像ID" type:"string"`
	UID     string `json:"uid" description:"Uid" type:"string"`
	Address string `json:"address" description:"网关地址" type:"string"`
}

//StartListDockerLogRes 开始监听docker日志
type StartListDockerLogRes struct {
	Success bool `json:"success" description:"结果" type:"bool"`
}

//EndListDockerLogReq 结束监听docker日志
type EndListDockerLogReq struct {
	UID string `json:"uid" description:"Uid" type:"string"`
}

//EndListDockerLogRes 结束监听docker日志
type EndListDockerLogRes struct {
	Success bool `json:"success" description:"结果" type:"bool"`
}

//CreateImageReq 创建镜像
type CreateImageReq struct {
	Image   string `json:"image" description:"镜像" type:"string"`
	Account string `json:"content" description:"账号" type:"string"`
	Pwd     string `json:"pwd" description:"密码" type:"string"`
}

//CreateImageRes 创建镜像
type CreateImageRes struct {
	Success bool `json:"success" description:"结果" type:"bool"`
}

//DelImageReq 删除镜像
type DelImageReq struct {
	ID string `json:"id" description:"镜像仓库ID" type:"string"`
}

//DelImageRes 删除镜像
type DelImageRes struct {
	Success bool `json:"success" description:"结果" type:"bool"`
}

//GetImageListReq 获取镜像仓库列表
type GetImageListReq struct {
	ID string `json:"id" description:"镜像仓库ID" type:"string"`
}

//GetImageListRes 获取镜像仓库列表
type GetImageListRes struct {
	Images []Image `json:"images" description:"列表" type:"[]Image"`
}

//Image image账号密码
type Image struct {
	ID      uint   `json:"id" description:"ID" type:"uint"`
	Image   string `json:"image" description:"镜像" type:"string"`
	Account string `json:"content" description:"账号" type:"string"`
	Pwd     string `json:"pwd" description:"密码" type:"string"`
	Time    string `json:"time" description:"时间" type:"string"`
}

//DelDockerReq 删除镜像请求
type DelDockerReq struct {
	DockerID string `json:"dockerId" description:"ID" type:"string"`
}

//DelDockerRes 删除镜像响应
type DelDockerRes struct {
	Success bool `json:"success" description:"结果" type:"bool"`
}

//RestartDockerReq 重启镜像请求
type RestartDockerReq struct {
	DockerID string `json:"dockerId" description:"ID" type:"string"`
}

//RestartDockerRes 重启镜像响应
type RestartDockerRes struct {
	Success bool `json:"success" description:"结果" type:"bool"`
}

//CloseDockerReq 关闭docker镜像
type CloseDockerReq struct {
	ID string `json:"id" description:"镜像ID" type:"string"`
}

//CloseDockerRes 关闭docker容器返回
type CloseDockerRes struct {
	Success bool `json:"success" description:"结果" type:"bool"`
}

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
	//容器ID
	ID string `json:"id" description:"id" type:"string"`
	//ID
	DockerID string `json:"dockerId" description:"id" type:"string"`
	//名称
	Name string `json:"name" description:"名称" type:"string"`
	//发布镜像
	Image string `json:"image" description:"发布镜像" type:"string"`
	//状态
	RunStatus bool `json:"runStatus" description:"允许状态" type:"bool"`
	//状态
	Status string `json:"status" description:"状态" type:"string"`
	//执行命令
	Command string `json:"command" description:"执行命令" type:"string"`
	//端口
	Ports string `json:"ports" description:"端口" type:"string"`
	//注册事件
	Time string `json:"time" description:"时间" type:"string"`
}
