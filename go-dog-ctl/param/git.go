package param

//CreateGitReq 创建Git
type CreateGitReq struct {
	Address string `json:"address" description:"git仓库地址" type:"string"`
	Account string `json:"account" description:"账号" type:"string"`
	Pwd     string `json:"pwd" description:"密码" type:"string"`
}

//CreateGitRes 创建Git
type CreateGitRes struct {
	Success bool `json:"success" description:"结果" type:"bool"`
}

//DelGitReq 删除Git
type DelGitReq struct {
	ID uint `json:"id" description:"Git仓库ID" type:"uint"`
}

//DelGitRes 删除Git
type DelGitRes struct {
	Success bool `json:"success" description:"结果" type:"bool"`
}

//GetGitListReq 获取Git仓库列表
type GetGitListReq struct {
}

//GetGitListRes 获取Git仓库列表
type GetGitListRes struct {
	Gits []Git `json:"gits" description:"列表" type:"[]Git"`
}

//Git image账号密码
type Git struct {
	ID      uint   `json:"id" description:"ID" type:"uint"`
	Address string `json:"address" description:"Git" type:"string"`
	Account string `json:"account" description:"账号" type:"string"`
	Pwd     string `json:"pwd" description:"密码" type:"string"`
	Time    string `json:"time" description:"时间" type:"string"`
}
