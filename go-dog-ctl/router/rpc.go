package router

import (
	"github.com/tang-go/go-dog-tool/go-dog-ctl/service"
)

//RPCRouter rpc路由
func RPCRouter(s *service.Service) {
	s.RPC("AdminOnline", 3, false, "管理员上线", s.AdminOnline)
	s.RPC("AdminOffline", 3, false, "管理员下线", s.AdminOffline)
	s.RPC("AuthAdmin", 3, false, "验证管理员", s.AuthAdmin)
	s.RPC("StartListDockerLog", 3, false, "开始监听docker日志", s.StartListDockerLog)
	s.RPC("EndListDockerLog", 3, false, "结束监听docker日志", s.EndListDockerLog)
}
