//Package define 定义框架常量
package define

//组织
const (
	Organize = "go-dog"
)

//服务定义
const (
	// SvcPrefix 统一前缀
	SvcPrefix = "go-dog/"
)

//定义服务
const (
	// SvcAdmin 控制中心
	SvcController = SvcPrefix + "controller"
	// SvcGateWay 网管服务
	SvcGateWay = SvcPrefix + "gateway"
	// SvcAuth 权限服务
	SvcAuth = SvcPrefix + "auth"
)

//Aes key加密
// const (
// 	AdminAes = "aes.go-dog-dog.code"
// )

//定义事件
const (
	//RunDockerEvent 启动docker
	RunDockerTopic = "run-docker-topic"
	//BuildServiceEvent 编译服务
	BuildServiceTopic = "build-service-topic"
	//K8sLogTopic k8s日志话题
	K8sLogTopic = "k8s-log-topic"
)
const (
	//AdminTokenValidityTime 管理员token有效时间 1小时
	AdminTokenValidityTime = 10
	//CodeValidityTime 验证码有效时间 5分钟
	CodeValidityTime = 60 * 5
)

const (
	//MaxClientRequestCount 客户端最大访问量
	MaxClientRequestCount = 1000000
	//MaxServiceRequestCount 服务端最大访问量
	MaxServiceRequestCount = 1000000
)
