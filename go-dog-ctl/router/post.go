package router

import (
	"github.com/tang-go/go-dog-tool/go-dog-ctl/service"
)

//POSTRouter post路由
func POSTRouter(s *service.Service) {
	//管理员相关
	s.POST("AdminLogin", "v1", "admin/login", 3, false, "管理员登录", s.AdminLogin)
	s.POST("CreateAdmin", "v1", "create/admin", 3, true, "创建管理员", s.CreateAdmin)
	s.POST("DelAdmin", "v1", "del/admin", 3, true, "删除管理员", s.DelAdmin)
	s.POST("DisableAdmin", "v1", "disable/admin", 3, true, "禁用管理员", s.DisableAdmin)
	s.POST("OpenAdmin", "v1", "open/admin", 3, true, "开启管理员", s.OpenAdmin)
	//角色相关
	s.POST("CreateMenu", "v1", "create/menu", 3, true, "创建菜单", s.CreateMenu)
	s.POST("DelMenu", "v1", "del/menu", 3, true, "删除菜单", s.DelMenu)
	s.POST("DelAPI", "v1", "del/api", 3, true, "删除API", s.DelAPI)
	s.POST("CreateRole", "v1", "create/role", 3, true, "创建角色", s.CreateRole)
	s.POST("DelRole", "v1", "del/role", 3, true, "删除角色", s.DelRole)
	s.POST("BindRoleMenu", "v1", "bind/role/menu", 3, true, "绑定角色菜单", s.BindRoleMenu)
	s.POST("DelRoleMenu", "v1", "del/role/menu", 3, true, "删除角色菜单", s.DelRoleMenu)
	s.POST("BindRoleAPI", "v1", "bind/role/api", 3, true, "绑定角色API", s.BindRoleAPI)
	s.POST("DelRoleAPI", "v1", "del/role/api", 3, true, "删除角色API", s.DelRoleAPI)
	//docker相关
	s.POST("BuildService", "v1", "build/service", 3, true, "编译发布服务", s.BuildService)
	s.POST("StartDocker", "v1", "strat/docker", 3, true, "docker方式启动服务", s.StartDocker)
	s.POST("CloseDocker", "v1", "clsoe/docker", 3, true, "关闭docker服务", s.CloseDocker)
	s.POST("DelDocker", "v1", "del/docker", 3, true, "删除docker服务", s.DelDocker)
	s.POST("RestartDocker", "v1", "restart/docker", 3, true, "重启docker服务", s.RestartDocker)
	//镜像仓库相关
	s.POST("CreateImage", "v1", "create/image", 3, true, "创建镜像仓库", s.CreateImage)
	s.POST("DelImage", "v1", "del/image", 3, true, "删除镜像仓库", s.DelImage)
	//GIT仓库相关
	s.POST("CreateGit", "v1", "create/git", 3, true, "创建GIT仓库", s.CreateGit)
	s.POST("DelGit", "v1", "del/git", 3, true, "删除GIT仓库", s.DelGit)
	//k8s相关
	//s.POST("CreateKubernetesDeployment", "v1", "create/kubernetes/deployment", 3, true, "创建一个kubernetes部署", s.CreateKubernetesDeployment)
	//s.POST("DeleteKubernetesDeployment", "v1", "delete/kubernetes/deployment", 3, true, "删除一个kubernetes部署", s.DeleteKubernetesDeployment)
}
