package param

import "github.com/tang-go/go-dog-tool/go-dog-ctl/table"

//CreateDocsReq 创建文档
type CreateDocsReq struct {
	Name string `json:"name" description:"名称" type:"string"`
	URL  string `json:"url" description:"文档地址" type:"string"`
}

//CreateDocsRes 创建文档
type CreateDocsRes struct {
	Success bool `json:"success" description:"结果" type:"bool"`
}

//DelDocsReq 删除文档
type DelDocsReq struct {
	ID int64 `json:"id" description:"文档ID" type:"int64"`
}

//DelDocsRes 删除文档
type DelDocsRes struct {
	Success bool `json:"success" description:"结果" type:"bool"`
}

//GetDocsReq 获取文档
type GetDocsReq struct {
}

//GetDocsRes 获取文档
type GetDocsRes struct {
	Docs []table.Docs `json:"docs" description:"文档ID" type:"[]table.Docs"`
}
