package param

//PushK8sLogMsg 推送K8s日志的消息
type PushK8sLogMsg struct {
	PodName string `json:"podName" description:"pod名称" type:"string"`
	PodLog  string `json:"podLog" description:"日志" type:"string"`
}
