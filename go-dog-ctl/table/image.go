package table

//Image image账号密码
type Image struct {
	ID      uint   `gorm:"primary_key" json:"id" description:"ID" type:"uint"`
	Address string `json:"address" description:"镜像仓库地址" type:"string"`
	Account string `json:"account" description:"账号" type:"string"`
	Pwd     string `json:"pwd" description:"密码" type:"string"`
	OwnerID int64  `json:"ownerID" description:"业主ID" type:"int64"`
	Time    int64  `json:"time" description:"时间" type:"int64"`
}
