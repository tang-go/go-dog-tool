package param

import "github.com/tang-go/go-dog-tool/go-dog-ctl/table"

//AuthReq 权限验证响应
type AuthReq struct {
	Token string `json:"token" description:"玩家token" type:"string"`
	URL   string `json:"url" description:"请求路径" type:"string"`
}

//AuthRes 权限验证返回
type AuthRes struct {
	Admin table.Admin `json:"admin" description:"管理员信息" type:"*table.Admin"`
}
