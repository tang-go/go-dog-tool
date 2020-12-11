package table

//Docs 文档
type Docs struct {
	ID      int64  `json:"id" description:"文档ID" type:"int64"`
	Name    string `json:"name" description:"名称" type:"string"`
	URL     string `json:"url" description:"文档地址" type:"string"`
	OwnerID int64  `json:"ownerID" description:"所属业主" type:"int64"`
	Time    int64  `json:"time" description:"角色" type:"int64"`
}
