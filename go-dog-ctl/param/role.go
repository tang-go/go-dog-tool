package param

//CreateMenuReq 创建菜单请求
type CreateMenuReq struct {
	Describe string `json:"describe" description:"描述" type:"string"`
	URL      string `json:"url" description:"菜单URL" type:"string"`
	ParentID uint   `json:"parentID" description:"父菜单ID" type:"string"`
	Sort     uint   `json:"sort" description:"排序" type:"uint"`
}

//CreateMenuRes 创建菜单响应
type CreateMenuRes struct {
	Success bool `json:"success" description:"结果" type:"bool"`
}
