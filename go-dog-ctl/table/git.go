package table

//Git git账号密码
type Git struct {
	ID      uint   `gorm:"primary_key" json:"id" description:"ID" type:"uint"`
	Account string `json:"content" description:"账号" type:"string"`
	Pwd     string `json:"pwd" description:"密码" type:"string"`
	Time    int64  `json:"time" description:"时间" type:"int64"`
}
