package param

import userTable "omo-service/services/user/dao/table"

//CreateUserTokenReq 创建用户token请求
type CreateUserTokenReq struct {
	SvcName string          `json:"svcName" description:"服务名称" type:"string"`
	User    *userTable.User `json:"value" description:"用户数据" type:"table.User"`
}

//CreateUserTokenRes 创建用户token请求
type CreateUserTokenRes struct {
	Token string `json:"token" description:"用户Token" type:"string"`
}

//UpdateUserTokenReq 更新用户token请求
type UpdateUserTokenReq struct {
	Token string          `json:"token" description:"用户Token" type:"string"`
	User  *userTable.User `json:"value" description:"用户数据" type:"table.User"`
}

//UpdateUserTokenRes 更新用户token响应
type UpdateUserTokenRes struct {
	Success bool `json:"success" description:"结果" type:"bool"`
}

//CheckUserTokenReq 校验用户token
type CheckUserTokenReq struct {
	Token   string `json:"token" description:"用户Token" type:"string"`
	Method  string `json:"method" description:"方法" type:"string"`
	SvcName string `json:"svcName" description:"服务名称" type:"string"`
}

//CheckUserTokenRes 校验用户token
type CheckUserTokenRes struct {
	User *userTable.User `json:"value" description:"用户数据" type:"table.User"`
}
