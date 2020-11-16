package param

import "omo-service/services/admin/table"

//CreateAdminTokenReq 创建管理员token请求
type CreateAdminTokenReq struct {
	SvcName string       `json:"svcName" description:"服务名称" type:"string"`
	Admin   *table.Admin `json:"value" description:"管理员数据" type:"table.Admin"`
}

//CreateAdminTokenRes 创建管理员token请求
type CreateAdminTokenRes struct {
	Token string `json:"token" description:"管理员Token" type:"string"`
}

//UpdateAdminTokenReq 更新管理员token请求
type UpdateAdminTokenReq struct {
	Token string       `json:"token" description:"管理员Token" type:"string"`
	Admin *table.Admin `json:"value" description:"管理员数据" type:"table.Admin"`
}

//UpdateAdminTokenRes 更新管理员token响应
type UpdateAdminTokenRes struct {
	Success bool `json:"success" description:"结果" type:"bool"`
}

//CheckAdminTokenReq 校验管理员token
type CheckAdminTokenReq struct {
	Token   string `json:"token" description:"管理员Token" type:"string"`
	Method  string `json:"method" description:"方法" type:"string"`
	SvcName string `json:"svcName" description:"服务名称" type:"string"`
}

//CheckAdminTokenRes 校验管理员token
type CheckAdminTokenRes struct {
	Admin *table.Admin `json:"value" description:"管理员数据" type:"table.Admin"`
}
