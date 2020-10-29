package param

//StartDockerReq 启动Docker请求
type StartDockerReq struct {
	Ports  []*Ports `json:"ports" description:"服务暴露端口" type:"[]*Ports"`
	Name   string   `json:"name" description:"服务名称" type:"string"`
	Images string   `json:"images" description:"镜像地址" type:"string"`
}

//Ports 暴露端口
type Ports struct {
	InsidePort   int `json:"insidePort" description:"内部端口" type:"int"`
	ExternalPort int `json:"externalPort" description:"外部暴露端口" type:"int"`
}

//StartDockerRes 启动Docker响应
type StartDockerRes struct {
	Success bool `json:"success" description:"结果" type:"bool"`
}
