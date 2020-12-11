package router

import (
	"github.com/tang-go/go-dog-tool/go-dog-ctl/service"
	"github.com/tang-go/go-dog/plugins"
)

//HTTPRouter http路由
func HTTPRouter(router plugins.Service, s *service.Service) {
	noauth := router.HTTP().APINoAuth()
	{
		noauth.APIGroup("管理员登陆").APIVersion("v1").POST("AdminLogin", "admin/login", "管理员登录", s.AdminLogin)
	}
	auth := router.HTTP().APIAuth()
	{
		authAdminV1 := auth.APIGroup("管理员相关").APIVersion("v1")
		{
			authAdminV1.APILevel(4).GET("GetAdminInfo", "get/admin/info", "获取管理员信息", s.GetAdminInfo)
			authAdminV1.APILevel(4).GET("GetAdminList", "get/admin/list", "获取管理员列表", s.GetAdminList)

			authAdminV1.APILevel(4).POST("CreateAdmin", "create/admin", "创建管理员", s.CreateAdmin)
			authAdminV1.APILevel(4).POST("DelAdmin", "del/admin", "删除管理员", s.DelAdmin)
			authAdminV1.APILevel(4).POST("DisableAdmin", "disable/admin", "禁用管理员", s.DisableAdmin)
			authAdminV1.APILevel(4).POST("OpenAdmin", "open/admin", "开启管理员", s.OpenAdmin)
		}
		authRoleV1 := auth.APIGroup("角色相关").APIVersion("v1")
		{
			authRoleV1.APILevel(4).GET("GetMenu", "get/menu", "获取菜单", s.GetMenu)
			authRoleV1.APILevel(4).GET("GetAPIList", "get/api/list", "获取API列表", s.GetAPIList)
			authRoleV1.APILevel(4).GET("GetRoleList", "get/role/list", "获取角色列表", s.GetRoleList)
			authRoleV1.APILevel(4).GET("GetRoleMenu", "get/role/menu/list", "获取角色菜单列表", s.GetRoleMenu)
			authRoleV1.APILevel(4).GET("GetRoleAPI", "get/role/api/list", "获取校色api列表", s.GetRoleAPI)

			authRoleV1.APILevel(4).POST("CreateMenu", "create/menu", "创建菜单", s.CreateMenu)
			authRoleV1.APILevel(4).POST("DelMenu", "del/menu", "删除菜单", s.DelMenu)
			authRoleV1.APILevel(4).POST("DelAPI", "del/api", "删除API", s.DelAPI)
			authRoleV1.APILevel(4).POST("CreateRole", "create/role", "创建角色", s.CreateRole)
			authRoleV1.APILevel(4).POST("DelRole", "del/role", "删除角色", s.DelRole)
			authRoleV1.APILevel(4).POST("BindRoleMenu", "bind/role/menu", "绑定角色菜单", s.BindRoleMenu)
			authRoleV1.APILevel(4).POST("DelRoleMenu", "del/role/menu", "删除角色菜单", s.DelRoleMenu)
			authRoleV1.APILevel(4).POST("BindRoleAPI", "bind/role/api", "绑定角色API", s.BindRoleAPI)
			authRoleV1.APILevel(4).POST("DelRoleAPI", "del/role/api", "删除角色API", s.DelRoleAPI)

		}
		authServiceV1 := auth.APIGroup("在线服务").APIVersion("v1")
		{
			authServiceV1.APILevel(4).GET("GetServiceList", "get/service/list", "获取服务列表", s.GetServiceList)
		}
		authBuildV1 := auth.APIGroup("编译发布").APIVersion("v1")
		{
			authBuildV1.APILevel(4).POST("BuildService", "build/service", "编译发布服务", s.BuildService)
			authBuildV1.APILevel(4).GET("GetBuildServiceList", "get/build/service/list", "获取编译发布记录", s.GetBuildServiceList)
		}
		authDockerV1 := auth.APIGroup("docker相关").APIVersion("v1")
		{
			authDockerV1.APILevel(4).GET("GetDockerList", "get/docker/list", "获取docker运行服务", s.GetDockerList)

			authDockerV1.APILevel(4).POST("StartDocker", "strat/docker", "docker方式启动服务", s.StartDocker)
			authDockerV1.APILevel(4).POST("CloseDocker", "clsoe/docker", "关闭docker服务", s.CloseDocker)
			authDockerV1.APILevel(4).POST("DelDocker", "del/docker", "删除docker服务", s.DelDocker)
			authDockerV1.APILevel(4).POST("RestartDocker", "restart/docker", "重启docker服务", s.RestartDocker)
		}
		authImageV1 := auth.APIGroup("镜像相关").APIVersion("v1")
		{
			authImageV1.APILevel(4).GET("GetImageList", "get/image/list", "获取镜像仓库列表", s.GetImageList)

			authImageV1.APILevel(4).POST("CreateImage", "create/image", "创建镜像仓库", s.CreateImage)
			authImageV1.APILevel(4).POST("DelImage", "del/image", "删除镜像仓库", s.DelImage)
		}
		authGitV1 := auth.APIGroup("git相关").APIVersion("v1")
		{
			authGitV1.APILevel(4).GET("GetGitList", "get/git/list", "获取GIT仓库列表", s.GetGitList)

			authGitV1.APILevel(4).POST("CreateGit", "create/git", "创建GIT仓库", s.CreateGit)
			authGitV1.APILevel(4).POST("DelGit", "del/git", "删除GIT仓库", s.DelGit)
		}
	}
	//k8s相关
	//s.GET("GetKubernetesNameSpace",  "get/kubernetes/namespace",  "获取k8s的namespace", s.GetKubernetesNameSpace)
	//s.GET("GetKubernetesDeployments",  "get/kubernetes/deployments",  "获取kubernetes的Deployments部署", s.GetKubernetesDeployments)
	//s.GET("GetKubernetesDeploymentInfoByName",  "get/kubernetes/deployment/info/by/name",  "根据Name获取kubernetes的Deployments部署的详情", s.GetKubernetesDeploymentInfoByName)
	//s.GET("GetKubernetesPodLog",  "get/kubernetes/pod/log",  "获取kubernetes的pod日志", s.GetKubernetesPodLog)
}
