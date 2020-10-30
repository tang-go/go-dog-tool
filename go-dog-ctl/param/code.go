package param

//GetCodeReq 获取验证码请求
type GetCodeReq struct {
}

//GetCodeRes 获取验证码响应
type GetCodeRes struct {
	ID  string `json:"id" description:"验证码ID" type:"string"`
	Img string `json:"img" description:"验证码图片" type:"string"`
}
