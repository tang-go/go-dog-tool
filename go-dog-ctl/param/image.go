package param

//CreateImageReq 创建镜像
type CreateImageReq struct {
	Address string `json:"address" description:"镜像仓库地址" type:"string"`
	Account string `json:"account" description:"账号" type:"string"`
	Pwd     string `json:"pwd" description:"密码" type:"string"`
}

//CreateImageRes 创建镜像
type CreateImageRes struct {
	Success bool `json:"success" description:"结果" type:"bool"`
}

//DelImageReq 删除镜像
type DelImageReq struct {
	ID uint `json:"id" description:"镜像仓库ID" type:"uint"`
}

//DelImageRes 删除镜像
type DelImageRes struct {
	Success bool `json:"success" description:"结果" type:"bool"`
}

//GetImageListReq 获取镜像仓库列表
type GetImageListReq struct {
}

//GetImageListRes 获取镜像仓库列表
type GetImageListRes struct {
	Images []Image `json:"images" description:"列表" type:"[]Image"`
}

//Image image账号密码
type Image struct {
	ID      uint   `json:"id" description:"ID" type:"uint"`
	Address string `json:"address" description:"镜像仓库地址" type:"string"`
	Account string `json:"account" description:"账号" type:"string"`
	Pwd     string `json:"pwd" description:"密码" type:"string"`
	Time    string `json:"time" description:"时间" type:"string"`
}
