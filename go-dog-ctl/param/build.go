package param

//BuildServiceReq 编译发布服务请求
type BuildServiceReq struct {
	Git      string `json:"git" description:"代码git地址" type:"string"`
	Path     string `json:"path" description:"编译项目的目录" type:"string"`
	Name     string `json:"name" description:"镜像名称" type:"string"`
	Harbor   string `json:"harbor" description:"镜像仓库" type:"string"`
	Accouunt string `json:"account" description:"Hrabor账号" type:"string"`
	Pwd      string `json:"pwd" description:"Harbor密码" type:"string"`
	Version  string `json:"version" description:"版本号" type:"string"`
}

//BuildServiceRes 编译发布服务返回
type BuildServiceRes struct {
	Success bool `json:"success" description:"结果" type:"bool"`
}

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
	ID string `json:"id" description:"唯一主键" type:"string"`
	//发布镜像
	Image string `json:"image" description:"发布镜像" type:"string"`
	//状态
	Status bool `json:"status" description:"状态" type:"bool"`
	//执行日志
	Log string `json:"log" description:"日志" type:"string"`
	//注册事件
	Time string `json:"time" description:"时间" type:"string"`
}
