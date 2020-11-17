package router

import (
	"github.com/tang-go/go-dog-tool/go-dog-auth/service"
)

//RPCRouter rpc路由
func RPCRouter(s *service.Service) {
	s.RPC("GetRoleMenu", 4, false, "获取角色拥有的菜单", s.GetRoleMenu)
	s.RPC("GetRoleApi", 4, false, "获取角色拥有的api", s.GetRoleApi)
	s.RPC("CreateRole", 4, false, "创建角色", s.CreateRole)
	s.RPC("SelectRoleByOrganize", 4, false, "更具组织查询角色", s.SelectRoleByOrganize)
	s.RPC("SelectRoleByID", 4, false, "更具ID查询角色", s.SelectRoleByID)
	s.RPC("CreateMenu", 4, false, "创建菜单", s.CreateMenu)
	s.RPC("CreateApi", 4, false, "创建api", s.CreateApi)
	s.RPC("BindRoleApi", 4, false, "绑定角色拥有的api", s.BindRoleApi)
	s.RPC("BindRoleMenu", 4, false, "绑定角色拥有的菜单", s.BindRoleMenu)

}
