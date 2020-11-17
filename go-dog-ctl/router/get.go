package router

import (
	"github.com/tang-go/go-dog-tool/go-dog-ctl/service"
)

//GETRouter get路由
func GETRouter(s *service.Service) {
	s.GET("GetCode", "v1", "get/code", 3, false, "获取图片验证码", s.GetCode)
	s.GET("GetAdminInfo", "v1", "get/admin/info", 3, true, "获取管理员信息", s.GetAdminInfo)
	s.GET("GetMenu", "v1", "get/menu", 3, true, "获取菜单", s.GetMenu)
	s.GET("GetBuildServiceList", "v1", "get/build/service/list", 3, true, "获取编译发布记录", s.GetBuildServiceList)
	s.GET("GetDockerList", "v1", "get/docker/list", 3, true, "获取docker运行服务", s.GetDockerList)
	s.GET("GetServiceList", "v1", "get/service/list", 3, true, "获取服务列表", s.GetServiceList)
	s.GET("GetKubernetesNameSpace", "v1", "get/kubernetes/namespace", 3, true, "获取k8s的namespace", s.GetKubernetesNameSpace)
	s.GET("GetKubernetesDeployments", "v1", "get/kubernetes/deployments", 3, true, "获取kubernetes的Deployments部署", s.GetKubernetesDeployments)
	s.GET("GetKubernetesDeploymentInfoByName", "v1", "get/kubernetes/deployment/info/by/name", 3, true, "根据Name获取kubernetes的Deployments部署的详情", s.GetKubernetesDeploymentInfoByName)
	s.GET("GetKubernetesPodLog", "v1", "get/kubernetes/pod/log", 3, true, "获取kubernetes的pod日志", s.GetKubernetesPodLog)
}
