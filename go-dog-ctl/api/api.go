package api

import (
	"archive/tar"
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/tang-go/go-dog-tool/define"
	"github.com/tang-go/go-dog-tool/go-dog-ctl/cfg"
	"github.com/tang-go/go-dog-tool/go-dog-ctl/table"
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
	//获取编译发布记录
	pointer.service.GET("GetBuildServiceList", "v1", "get/build/service/list",
		3,
		true,
		"获取编译发布记录",
		pointer.GetBuildServiceList)
	//发布服务
	pointer.service.POST("BuildService", "v1", "build/service",
		3,
		true,
		"编译发布服务",
		pointer.BuildService)
	//获取docker运行服务
	pointer.service.GET("GetDockerList", "v1", "get/docker/list",
		3,
		true,
		"获取docker运行服务",
		pointer.GetDockerList)
	//docker启动服务
	pointer.service.POST("StartDocker", "v1", "strat/docker",
		3,
		true,
		"docker方式启动服务",
		pointer.StartDocker)
	//docker 关闭
	pointer.service.POST("CloseDocker", "v1", "clsoe/docker",
		3,
		true,
		"关闭docker服务",
		pointer.CloseDocker)
	//删除docker服务
	pointer.service.POST("DelDocker", "v1", "del/docker",
		3,
		true,
		"删除docker服务",
		pointer.DelDocker)
	//删除docker服务
	pointer.service.POST("RestartDocker", "v1", "restart/docker",
		3,
		true,
		"重启docker服务",
		pointer.RestartDocker)
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
	cfg       *cfg.Config
	docker    *client.Client
	mysql     *mysql.Mysql
	snowflake *snowflake.SnowFlake
	cache     *cache.Cache
	lock      sync.RWMutex
}

//NewService 初始化服务
func NewService() *API {
	ctl := new(API)
	ctl.cfg = cfg.NewConfig()
	//cli, err := client.NewClient("tcp://127.0.0.1:3375", "v1.39", nil, nil)
	cli, err := client.NewClient("unix:///var/run/docker.sock", "v1.39", nil, nil)
	if err != nil {
		panic(err.Error())
	}
	ctl.docker = cli
	//初始化rpc服务端
	ctl.service = service.CreateService(define.SvcController, ctl.cfg)
	//验证函数
	ctl.service.Auth(ctl.Auth)
	//设置服务端最大访问量
	ctl.service.GetLimit().SetLimit(define.MaxServiceRequestCount)
	//设置客户端最大访问量
	ctl.service.GetClient().GetLimit().SetLimit(define.MaxClientRequestCount)
	//初始化数据库
	ctl.mysql = mysql.NewMysql(ctl.cfg)
	//初始化数据库表
	ctl.mysql.GetWriteEngine().AutoMigrate(
		table.Admin{},
		table.Owner{},
		table.OwnerRole{},
		table.Permission{},
		table.RolePermission{},
		table.BuildService{},
		table.Docker{},
		table.Log{},
	)
	//初始化缓存
	ctl.cache = cache.NewCache(ctl.cfg)
	//初始化雪花算法
	ret := big.NewInt(0)
	ret.SetBytes(net.ParseIP(ctl.cfg.GetHost()).To4())
	id, err := strconv.ParseInt(fmt.Sprintf("%d%d", ret.Int64(), ctl.cfg.GetPort()), 10, 64)
	if err != nil {
		panic(err)
	}
	ctl.snowflake = snowflake.NewSnowFlake(id)
	//初始化API
	ctl.Router()
	//初始化数据库数据
	ctl._InitMysql(ctl.cfg.GetPhone(), ctl.cfg.GetPwd())
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

func (pointer *API) _RunInLinux(cmd string, success func(string), fail func(string)) error {
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
			if success != nil {
				success(logScan.Text())
			}
		}
	}()
	//错误
	errScan := bufio.NewScanner(stderr)
	go func() {
		for errScan.Scan() {
			if fail != nil {
				fail(errScan.Text())
			}
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

//BuildImage 编译镜像
func (pointer *API) _BuildImage(tarFile, project, imageName string, success func(string)) error {
	dockerBuildContext, err := os.Open(tarFile)
	if err != nil {
		return err
	}
	defer dockerBuildContext.Close()

	buildOptions := types.ImageBuildOptions{
		Dockerfile: "Dockerfile", // optional, is the default
		Tags:       []string{imageName},
		Labels: map[string]string{
			project: "project",
		},
	}
	output, err := pointer.docker.ImageBuild(context.Background(), dockerBuildContext, buildOptions)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(output.Body)
	for scanner.Scan() {
		success(scanner.Text())
	}
	return nil
}

//_CloseDocker 关闭镜像
func (pointer *API) _CloseDocker(id string) error {
	err := pointer.docker.ContainerStop(context.TODO(), id, nil)
	if err != nil {
		return err
	}
	err = pointer.docker.ContainerRemove(context.TODO(), id, types.ContainerRemoveOptions{})
	if err != nil {
		return err
	}
	return nil
}

//PushImage 推送镜像
func (pointer *API) _PushImage(registryUser, registryPassword, image string, success func(string)) error {
	config := types.ImagePushOptions{}
	if len(registryUser) > 0 || len(registryPassword) > 0 {
		authConfig := types.AuthConfig{
			Username: registryUser,
			Password: registryPassword,
		}
		encodedJSON, err := json.Marshal(authConfig)
		if err != nil {
			return err
		}
		authStr := base64.URLEncoding.EncodeToString(encodedJSON)
		config.RegistryAuth = authStr
	}
	out, err := pointer.docker.ImagePush(context.TODO(), image, config)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(out)
	for scanner.Scan() {
		success(scanner.Text())
	}
	return nil
}

//PushImage 推送镜像
func (pointer *API) _PullImage(registryUser, registryPassword, image string, success func(string)) error {
	config := types.ImagePullOptions{}
	if registryUser != "" || registryPassword != "" {
		authConfig := types.AuthConfig{
			Username: registryUser,
			Password: registryPassword,
		}
		encodedJSON, err := json.Marshal(authConfig)
		if err != nil {
			return err
		}
		authStr := base64.URLEncoding.EncodeToString(encodedJSON)
		config.RegistryAuth = authStr
	}
	out, err := pointer.docker.ImagePull(context.TODO(), image, config)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(out)
	for scanner.Scan() {
		success(scanner.Text())
	}
	return nil
}

// 打包tar
func (pointer *API) _CreateTar(filesource, filetarget string, deleteIfExist bool) error {
	tarfile, err := os.Create(filetarget)
	if err != nil {
		if err == os.ErrExist {
			if err := os.Remove(filetarget); err != nil {
				log.Errorln(err.Error())
				return err
			}
		} else {
			log.Errorln(err.Error())
			return err
		}
	}
	defer tarfile.Close()
	tarwriter := tar.NewWriter(tarfile)

	files, _ := ioutil.ReadDir(filesource)
	for _, f := range files {
		path := filesource + "/" + f.Name()
		sfileInfo, err := os.Stat(path)
		if err != nil {
			log.Errorln(err.Error())
			return err
		}
		if !sfileInfo.IsDir() {
			if err := pointer._TarFile(f.Name(), path, sfileInfo, tarwriter); err != nil {
				log.Errorln(err.Error())
				return err
			}
		} else {
			if err := pointer._TarFolder(path, tarwriter); err != nil {
				log.Errorln(err.Error())
				return err
			}
		}
	}
	return nil
}

func (pointer *API) _TarFile(directory string, filesource string, sfileInfo os.FileInfo, tarwriter *tar.Writer) error {
	sfile, err := os.Open(filesource)
	if err != nil {
		log.Errorln(err.Error())
		return err
	}
	defer sfile.Close()
	header, err := tar.FileInfoHeader(sfileInfo, "")
	if err != nil {
		log.Errorln(err.Error())
		return err
	}
	header.Name = directory
	err = tarwriter.WriteHeader(header)
	if err != nil {
		log.Errorln(err.Error())
		return err
	}
	if _, err = io.Copy(tarwriter, sfile); err != nil {
		log.Errorln(err.Error())
		return err
	}
	return nil
}

func (pointer *API) _TarFolder(directory string, tarwriter *tar.Writer) error {
	baseFolder := filepath.Base(directory)
	return filepath.Walk(directory, func(targetpath string, file os.FileInfo, err error) error {
		if file == nil {
			log.Errorln(err.Error())
			return err
		}
		if file.IsDir() {
			header, err := tar.FileInfoHeader(file, "")
			if err != nil {
				log.Errorln(err.Error())
				return err
			}
			header.Name = filepath.Join(baseFolder, strings.TrimPrefix(targetpath, directory))
			if err = tarwriter.WriteHeader(header); err != nil {
				log.Errorln(err.Error())
				return err
			}
			os.Mkdir(strings.TrimPrefix(baseFolder, file.Name()), os.ModeDir)
			return nil
		}
		fileFolder := baseFolder + "/" + file.Name()
		return pointer._TarFile(fileFolder, targetpath, file, tarwriter)
	})
}
