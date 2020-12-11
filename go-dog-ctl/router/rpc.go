package router

import (
	"github.com/tang-go/go-dog-tool/go-dog-ctl/service"
	"github.com/tang-go/go-dog/plugins"
)

//RPCRouter rpc路由
func RPCRouter(router plugins.Service, s *service.Service) {
	router.RPC("AdminOnline", 3, false, "管理员上线", s.AdminOnline)
	router.RPC("AdminOffline", 3, false, "管理员下线", s.AdminOffline)
	router.RPC("AuthAdmin", 3, false, "验证管理员", s.AuthAdmin)
	router.RPC("StartListDockerLog", 3, false, "开始监听docker日志", s.StartListDockerLog)
	router.RPC("EndListDockerLog", 3, false, "结束监听docker日志", s.EndListDockerLog)
}
