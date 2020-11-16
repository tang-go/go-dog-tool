package service

import (
	"omo-service/define"
	"omo-service/services/auth/table"

	"github.com/tang-go/go-dog/cache"
	"github.com/tang-go/go-dog/mysql"
	"github.com/tang-go/go-dog/pkg/service"
	"github.com/tang-go/go-dog/plugins"
)

//Service 用户服务
type Service struct {
	service plugins.Service
	cache   *cache.Cache
	mysql   *mysql.Mysql
}

//NewService 初始化服务
func NewService(routers ...func(*Service)) *Service {
	s := new(Service)
	//初始化rpc服务端
	s.service = service.CreateService(define.SvcAuth)
	//设置服务端最大访问量
	s.service.GetLimit().SetLimit(define.MaxServiceRequestCount)
	//设置客户端最大访问量
	s.service.GetClient().GetLimit().SetLimit(define.MaxClientRequestCount)
	//初始化数据库
	s.mysql = mysql.NewMysql(s.service.GetClient().GetCfg())
	//初始化数据库表
	s.mysql.GetWriteEngine().AutoMigrate(
		//系统角色
		table.SysRole{},
		//系统角色菜单映射表
		table.SysRoleMenu{},
		//系统角色菜单表
		table.SysMenu{},
		//系统角色api映射表
		table.SysRoleApi{},
		//系统角色api表
		table.SysApi{},
	)
	//初始化缓存
	s.cache = cache.NewCache(s.service.GetClient().GetCfg())
	//初始化路由
	for _, router := range routers {
		router(s)
	}
	return s
}

//RegisterRPC 	注册RPC方法
//name			方法名称
//level			方法等级
//isAuth		是否需要鉴权
//explain		方法说明
//fn			注册的方法
func (s *Service) RPC(name string, level int8, isAuth bool, explain string, fn interface{}) {
	s.service.RPC(name, level, isAuth, explain, fn)
}

//Run 启动
func (s *Service) Run() error {
	return s.service.Run()
}
