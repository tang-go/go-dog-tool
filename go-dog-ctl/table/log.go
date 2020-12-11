package table

const (
	//LoginType 登录类型
	LoginType int32 = iota
	//BuildServiceType 编译发布服务类型
	BuildServiceType
	//CloseDockerType 关闭docker服务
	CloseDockerType
	//StartDockerType 启动docker服务
	StartDockerType
	//RestartDockerType 重启docker服务
	RestartDockerType
	//CreateMenuType 创建菜单
	CreateMenuType
	//DelMenuType 删除菜单
	DelMenuType
	//CreateRoleType 创建角色
	CreateRoleType
	//DelRoleType 删除角色
	DelRoleType
	//BindRoleMenuType 绑定角色菜单
	BindRoleMenuType
	//DelRoleMenuType 删除角色菜单
	DelRoleMenuType
	//DelAPIType 删除API
	DelAPIType
	//BindRoleAPIType 绑定角色api
	BindRoleAPIType
	//DelRoleAPIType 删除角色api
	DelRoleAPIType
	//CreateAdminType 创建管理员
	CreateAdminType
	//DisableAdminType 禁用管理员
	DisableAdminType
	//OpenAdminType 开启管理员
	OpenAdminType
	//DelAdminType 删除管理员
	DelAdminType
	//CreateImageType 创建镜像账号密码
	CreateImageType
	//DelImageType 删除镜像账号密码
	DelImageType
	//CreateGitType 创建git账号密码
	CreateGitType
	//DelGitType 删除git账号密码
	DelGitType
	//CreateDocsType 创建文档
	CreateDocsType
	//DelDocsType 删除文档
	DelDocsType
)

//Log 日志
type Log struct {
	//日志ID
	LogID int64
	//类型
	Type int32
	//操作人
	AdminID int64
	//名称
	AdminName string
	//操作方法
	Method string
	//描述
	Description string
	//业主ID
	OwnerID int64
	//操作IP
	IP string
	//操作URL
	URL string
	//操作时间
	Time int64
}
