package param

//PushReq 推送请求
type PushReq struct {
	Token string `json:"token" description:"推送目标token" type:"string"`
	Topic string `json:"topic" description:"推送消息话题" type:"string"`
	Msg   string `json:"msg" description:"推送消息" type:"string"`
}

//PushRes 推送响应
type PushRes struct {
	Success bool `json:"success" description:"结果" type:"bool"`
}

//XtermPustReq xterm推送日志
type XtermPushReq struct {
	Uid string `json:"uid" description:"推送目标uid" type:"string"`
	Msg string `json:"msg" description:"内容" type:"string"`
}

//XtermPustReq xterm推送日志
type XtermPushRes struct {
	Success bool `json:"success" description:"结果" type:"bool"`
}
