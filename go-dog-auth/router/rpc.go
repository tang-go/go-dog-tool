package router

import (
	"github.com/tang-go/go-dog-tool/go-dog-auth/service"
)

//RPCRouter rpc路由
func RPCRouter(s *service.Service) {
	s.RPC("CreateUserToken", 4, false, "创建用户Token", s.CreateUserToken)
	s.RPC("CheckUserToken", 4, false, "校验用户token", s.CheckUserToken)
	s.RPC("UpdateUserToken", 4, false, "更新用户Token", s.UpdateUserToken)

	s.RPC("CreateAdminToken", 4, false, "创建管理员Token", s.CreateAdminToken)
	s.RPC("CheckAdminToken", 4, false, "校验管理员token", s.CheckAdminToken)
	s.RPC("UpdateAdminToken", 4, false, "更管理员Token", s.UpdateAdminToken)

	s.RPC("GetRoleMenu", 4, false, "获取角色拥有的菜单", s.GetRoleMenu)
	s.RPC("GetRoleApi", 4, false, "获取角色拥有的api", s.GetRoleApi)
	s.RPC("CreateRole", 4, false, "创建角色", s.CreateRole)
	s.RPC("SelectRole", 4, false, "更具组织查询角色", s.SelectRole)
	s.RPC("CreateMenu", 4, false, "创建菜单", s.CreateMenu)
	s.RPC("CreateApi", 4, false, "创建api", s.CreateApi)
	s.RPC("BindRoleApi", 4, false, "绑定角色拥有的api", s.BindRoleApi)
	s.RPC("BindRoleMenu", 4, false, "绑定角色拥有的菜单", s.BindRoleMenu)

}
