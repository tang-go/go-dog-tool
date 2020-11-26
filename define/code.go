package define

//返回码定义
const (
	//SuccessCode 成功返回码
	SuccessCode = 10000
	//GetCodeErr 获取验证码失败
	GetCodeErr = 10001
	//VerfiyCodeErr 验证码失败
	VerfiyCodeErr = 10002
	//AdminLoginErr 管理员登录失败
	AdminLoginErr = 10003
	//AdminTokenErr 管理员Token验证失败
	AdminTokenErr = 10004
	//GetAdminInfoErr 获取管理员信息失败
	GetAdminInfoErr = 10005
	//APIPowerErr API权限不够
	APIPowerErr = 10006
	//BuildServiceErr 编译发布服务失败
	BuildServiceErr = 10007
	//GetBuildServiceListErr 获取编译发布服务失败
	GetBuildServiceListErr = 10008
	//StartDockerErr 启动docker失败
	StartDockerErr = 10009
	//CloseDockerErr 关闭docker失败
	CloseDockerErr = 10010
	//DelDockerErr 删除失败
	DelDockerErr = 10011
	//RestartDockerErr 重启失败
	RestartDockerErr = 10012
	//GetKubernetesNameSpaceErr 获取k8s的namespqce失败
	GetKubernetesNameSpaceErr = 10013
	//GetKubernetesPodsByNamsespaceErr 获取k8s的pods失败
	GetKubernetesPodsByNamsespace = 10014
	//GetKubernetesDeploymentInfoByNameErr Deployments部署的详情失败
	GetKubernetesDeploymentInfoByNameErr = 10015
	//CreateKubernetesDeploymentErr 创建部署失败
	CreateKubernetesDeploymentErr = 10016
	//DeleteKubernetesDeploymentErr 删除失败
	DeleteKubernetesDeploymentErr = 10017
	//AdminOnlineErr 管理员上线错误
	AdminOnlineErr = 10018
	//AdminOfflineErr 管理员下线错误
	AdminOfflineErr = 10019
	//GetKubernetesPodLogErr 获取pod日志失败
	GetKubernetesPodLogErr = 10020
	//CreateMenuErr 创建菜单
	CreateMenuErr = 10021
	//DelMenuErr 删除菜单
	DelMenuErr = 10022
	//GetRoleListErr 获取角色列表失败
	GetRoleListErr = 10023
	//CreateRoleErr 创建角色失败
	CreateRoleErr = 10024
	//DelRoleErr 删除角色失败
	DelRoleErr = 10025
	//GetRoleMenuErr 获取角色菜单失败
	GetRoleMenuErr = 10026
	//BindRoleMenuErr 绑定角色菜单失败
	BindRoleMenuErr = 10027
	//DelRoleMenuErr 删除角色菜单失败
	DelRoleMenuErr = 10028
	//GetAPIListErr 获取API列表
	GetAPIListErr = 10029
	//GetMenuErr 获取菜单失败
	GetMenuErr = 10030
	//DelAPIErr 删除api
	DelAPIErr = 10031
	//GetRoleAPIErr 获取角色api失败
	GetRoleAPIErr = 10032
	//BindRoleAPIErr 绑定角色api失败
	BindRoleAPIErr = 10033
	//DelRoleAPIErr 删除角色api失败
	DelRoleAPIErr = 10034
	//GetAdminListErr 获取管理员列表失败
	GetAdminListErr = 10035
	//CreateAdminErr 创建管理员失败
	CreateAdminErr = 10036
	//DisableAdminErr 禁用管理员失败
	DisableAdminErr = 10037
	//OpenAdminErr 开启管理员
	OpenAdminErr = 10038
	//DelAdminErr 删除管理员
	DelAdminErr = 10039
	//StartListDockerLogErr 开始监听docker日志
	StartListDockerLogErr = 10040
	//EndListDockerLogErr 结束监听docker日志
	EndListDockerLogErr = 10041
)
