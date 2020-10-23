package param

//BuildServiceReq 编译发布服务请求
type BuildServiceReq struct {
	Git     string `json:"git" description:"代码git地址" type:"string"`
	Path    string `json:"path" description:"编译项目的目录" type:"string"`
	Name    string `json:"name" description:"镜像名称" type:"string"`
	Harbor  string `json:"harbor" description:"镜像仓库" type:"string"`
	Version string `json:"version" description:"版本号" type:"string"`
}

//BuildServiceRes 编译发布服务返回
type BuildServiceRes struct {
	Result string `json:"result" description:"结果" type:"string"`
}
