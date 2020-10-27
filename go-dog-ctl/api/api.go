package api

import (
	"bufio"
	"fmt"
	"math/big"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/tang-go/go-dog-tool/define"
	"github.com/tang-go/go-dog-tool/go-dog-ctl/table"
	"github.com/tang-go/go-dog-tool/go-dog-gw/param"
	"github.com/tang-go/go-dog/cache"
	"github.com/tang-go/go-dog/lib/md5"
	"github.com/tang-go/go-dog/lib/rand"
	"github.com/tang-go/go-dog/lib/snowflake"
	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/mysql"
	"github.com/tang-go/go-dog/pkg/service"
	"github.com/tang-go/go-dog/plugins"
	"github.com/tang-go/go-dog/serviceinfo"
)

//Router 注册路由
func (pointer *API) Router() {
	//获取图片验证码
	pointer.service.GET("GetCode", "v1", "get/code",
		3,
		false,
		"获取图片验证码",
		pointer.GetCode)
	//管理员登录
	pointer.service.POST("AdminLogin", "v1", "admin/login",
		3,
		false,
		"管理员登录",
		pointer.AdminLogin)
	//获取管理员信息
	pointer.service.GET("GetAdminInfo", "v1", "get/admin/info",
		3,
		true,
		"获取管理员信息",
		pointer.GetAdminInfo)
	//获取角色列表
	pointer.service.GET("GetRoleList", "v1", "get/role/list",
		3,
		true,
		"获取角色列表",
		pointer.GetRoleList)
	//发布服务
	pointer.service.POST("BuildService", "v1", "build/service",
		3,
		true,
		"编译发布服务",
		pointer.BuildService)
	//docker启动服务
	pointer.service.POST("StartDocker", "v1", "strat/docker",
		3,
		true,
		"docker方式启动服务",
		pointer.StartDocker)
	//获取服务列表
	pointer.service.GET("GetServiceList", "v1", "get/service/list",
		3,
		true,
		"获取服务列表",
		pointer.GetServiceList)
}

//APIService API服务
type _APIService struct {
	method  *serviceinfo.API
	name    string
	explain string
	count   int32
}

//API 控制服务
type API struct {
	service   plugins.Service
	mysql     *mysql.Mysql
	snowflake *snowflake.SnowFlake
	cache     *cache.Cache
	lock      sync.RWMutex
}

//NewService 初始化服务
func NewService() *API {
	ctl := new(API)
	//初始化rpc服务端
	ctl.service = service.CreateService(define.SvcController)
	//验证函数
	ctl.service.Auth(ctl.Auth)
	//设置服务端最大访问量
	ctl.service.GetLimit().SetLimit(define.MaxServiceRequestCount)
	//设置客户端最大访问量
	ctl.service.GetClient().GetLimit().SetLimit(define.MaxClientRequestCount)
	//初始化数据库
	ctl.mysql = mysql.NewMysql(ctl.service.GetCfg())
	//初始化数据库表
	ctl.mysql.GetWriteEngine().AutoMigrate(
		table.Admin{},
		table.Owner{},
		table.OwnerRole{},
		table.Permission{},
		table.RolePermission{},
		table.Log{},
	)
	//初始化缓存
	ctl.cache = cache.NewCache(ctl.service.GetCfg())
	//初始化雪花算法
	ret := big.NewInt(0)
	ret.SetBytes(net.ParseIP(ctl.service.GetCfg().GetHost()).To4())
	id, err := strconv.ParseInt(fmt.Sprintf("%d%d", ret.Int64(), ctl.service.GetCfg().GetPort()), 10, 64)
	if err != nil {
		panic(err)
	}
	ctl.snowflake = snowflake.NewSnowFlake(id)
	//初始化API
	ctl.Router()
	//初始化数据库数据
	ctl._InitMysql("13688460148", "admin")
	return ctl
}

//Run 启动
func (pointer *API) Run() error {
	return pointer.service.Run()
}

//_InitMysql 第一次加载初始化数据库数据
func (pointer *API) _InitMysql(phone, pwd string) {
	//读取是否有业主了
	owner := new(table.Owner)
	if pointer.mysql.GetReadEngine().Where("phone = ?", phone).First(owner).RecordNotFound() == false {
		return
	}
	//如果没有业主则新增默认业主
	owner.OwnerID = pointer.snowflake.GetID()
	owner.Name = "超级业主"
	owner.Phone = phone
	owner.Level = 1
	owner.IsDisable = table.OwnerAvailable
	owner.IsAdminOwner = table.IsAdminOwner
	owner.Time = time.Now().Unix()
	//超级管理员
	ownerRole := &table.OwnerRole{
		RoleID: pointer.snowflake.GetID(),
		//角色名称
		Name: "admin",
		//角色描述
		Description: "系统自带的超级管理员",
		//是否为超级管理员
		IsAdmin: table.IsAdmin,
		//业主ID
		OwnerID: owner.OwnerID,
		//角色创建时间
		Time: owner.Time,
	}
	//管理员
	admin := &table.Admin{
		//账号 唯一主键
		AdminID: pointer.snowflake.GetID(),
		//名称
		Name: "admin",
		//电话
		Phone: phone,
		//盐值 md5使用
		Salt: rand.StringRand(6),
		//等级
		Level: owner.Level,
		//所属业主
		OwnerID: owner.OwnerID,
		//是否被禁用
		IsDisable: table.AdminAvailable,
		//类型
		RoleID: ownerRole.RoleID,
		//注册事件
		Time: owner.Time,
	}
	//生成密码
	admin.Pwd = md5.Md5(md5.Md5(pwd) + admin.Salt)
	//开启数据库操作
	tx := pointer.mysql.GetWriteEngine().Begin()
	if err := tx.Create(owner).Error; err != nil {
		tx.Rollback()
		panic(err)
	}
	if err := tx.Create(ownerRole).Error; err != nil {
		tx.Rollback()
		panic(err)
	}
	if err := tx.Create(admin).Error; err != nil {
		tx.Rollback()
		panic(err)
	}
	tx.Commit()
}

// Set 设置验证码ID
func (pointer *API) Set(id string, value string) {
	if err := pointer.cache.GetCache().SetByTime(id, value, define.CodeValidityTime); err != nil {
		log.Errorln(err.Error())
	}
}

// Get 更具验证ID获取验证码
func (pointer *API) Get(id string, clear bool) (vali string) {
	err := pointer.cache.GetCache().Get(id, &vali)
	if err != nil {
		log.Errorln(err.Error())
	}
	if clear {
		pointer.cache.GetCache().Del(id)
	}
	return
}

//Verify 验证验证码
func (pointer *API) Verify(id, answer string, clear bool) bool {
	vali := pointer.Get(id, clear)
	if strings.ToLower(vali) != strings.ToLower(answer) {
		return false
	}
	return true
}

func (pointer *API) _RunInLinux(ctx plugins.Context, token string, topic string, cmd string) error {
	c := exec.Command("sh", "-c", cmd)
	stdout, err := c.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stdin, "error=>", err.Error())
		return err
	}
	stderr, err := c.StderrPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error=>", err.Error())
		return err
	}
	c.Start()
	// 正常日志
	logScan := bufio.NewScanner(stdout)
	go func() {
		for logScan.Scan() {
			res := new(param.PushRes)
			ctx.GetClient().Broadcast(
				ctx,
				define.SvcGateWay,
				"Push",
				&param.PushReq{
					Token: token,
					Topic: topic,
					Msg:   logScan.Text(),
				},
				res)
		}
	}()
	//错误
	errScan := bufio.NewScanner(stderr)
	go func() {
		for errScan.Scan() {
			res := new(param.PushRes)
			ctx.GetClient().Broadcast(
				ctx,
				define.SvcGateWay,
				"Push",
				&param.PushReq{
					Token: token,
					Topic: topic,
					Msg:   errScan.Text(),
				},
				res)
		}
	}()
	c.Wait()
	return nil
}

func (pointer *API) _PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
