package param

//CreateMenuReq 创建菜单请求
type CreateMenuReq struct {
	Describe string `json:"describe" description:"描述" type:"string"`
	URL      string `json:"url" description:"菜单URL" type:"string"`
	ParentID uint   `json:"parentID" description:"父菜单ID" type:"uint"`
	Sort     uint   `json:"sort" description:"排序" type:"uint"`
}

//CreateMenuRes 创建菜单响应
type CreateMenuRes struct {
	Success bool `json:"success" description:"结果" type:"bool"`
}

//GetMenuReq 获取菜单请求
type GetMenuReq struct {
}

//GetMenuRes 获取菜单响应
type GetMenuRes struct {
	Menu []*Menu `json:"menu" description:"菜单" type:"Role"`
}

//Menu 菜单配置
type Menu struct {
	ID       uint    `json:"id" description:"ID" type:"uint"`
	ParentID uint    `json:"parentID" description:"父亲结点ID" type:"uint"`
	Organize string  `json:"organize" description:"组织" type:"string"`
	Describe string  `json:"describe" description:"描述" type:"string"`
	URL      string  `json:"url" description:"菜单URL" type:"string"`
	Sort     uint    `json:"sort" description:"排序" type:"uint"`
	Children []*Menu `json:"children" description:"菜单" type:"[]*Menu"`
	Time     string  `json:"string" description:"时间" type:"int64"`
}
