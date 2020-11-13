package api

import (
	"archive/tar"
	"bufio"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/tang-go/go-dog-tool/define"
	"github.com/tang-go/go-dog-tool/go-dog-ctl/cfg"
	"github.com/tang-go/go-dog-tool/go-dog-ctl/param"
	"github.com/tang-go/go-dog-tool/go-dog-ctl/table"
	gateParam "github.com/tang-go/go-dog-tool/go-dog-gw/param"
	"github.com/tang-go/go-dog/cache"
	"github.com/tang-go/go-dog/etcd"
	"github.com/tang-go/go-dog/lib/md5"
	"github.com/tang-go/go-dog/lib/rand"
	"github.com/tang-go/go-dog/lib/snowflake"
	"github.com/tang-go/go-dog/lib/uuid"
	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/mysql"
	"github.com/tang-go/go-dog/pkg/context"
	"github.com/tang-go/go-dog/pkg/service"
	"github.com/tang-go/go-dog/plugins"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

//Router 注册路由
func (pointer *API) Router() {
	pointer.service.RPC("AdminOnline", 3, true, "管理员上线", pointer.AdminOnline)
	pointer.service.RPC("AdminOffline", 3, true, "管理员下线", pointer.AdminOffline)
	pointer.service.GET("GetCode", "v1", "get/code", 3, false, "获取图片验证码", pointer.GetCode)
	pointer.service.POST("AdminLogin", "v1", "admin/login", 3, false, "管理员登录", pointer.AdminLogin)
	pointer.service.GET("GetAdminInfo", "v1", "get/admin/info", 3, true, "获取管理员信息", pointer.GetAdminInfo)
	pointer.service.GET("GetRoleList", "v1", "get/role/list", 3, true, "获取角色列表", pointer.GetRoleList)
	pointer.service.GET("GetBuildServiceList", "v1", "get/build/service/list", 3, true, "获取编译发布记录", pointer.GetBuildServiceList)
	pointer.service.POST("BuildService", "v1", "build/service", 3, true, "编译发布服务", pointer.BuildService)
	pointer.service.GET("GetDockerList", "v1", "get/docker/list", 3, true, "获取docker运行服务", pointer.GetDockerList)
	pointer.service.POST("StartDocker", "v1", "strat/docker", 3, true, "docker方式启动服务", pointer.StartDocker)
	pointer.service.POST("CloseDocker", "v1", "clsoe/docker", 3, true, "关闭docker服务", pointer.CloseDocker)
	pointer.service.POST("DelDocker", "v1", "del/docker", 3, true, "删除docker服务", pointer.DelDocker)
	pointer.service.POST("RestartDocker", "v1", "restart/docker", 3, true, "重启docker服务", pointer.RestartDocker)
	pointer.service.GET("GetServiceList", "v1", "get/service/list", 3, true, "获取服务列表", pointer.GetServiceList)
	pointer.service.GET("GetKubernetesNameSpace", "v1", "get/kubernetes/namespace", 3, true, "获取k8s的namespace", pointer.GetKubernetesNameSpace)
	pointer.service.GET("GetKubernetesDeployments", "v1", "get/kubernetes/deployments", 3, true, "获取kubernetes的Deployments部署", pointer.GetKubernetesDeployments)
	pointer.service.GET("GetKubernetesDeploymentInfoByName", "v1", "get/kubernetes/deployment/info/by/name", 3, true, "根据Name获取kubernetes的Deployments部署的详情", pointer.GetKubernetesDeploymentInfoByName)
	pointer.service.POST("CreateKubernetesDeployment", "v1", "create/kubernetes/deployment", 3, true, "创建一个kubernetes部署", pointer.CreateKubernetesDeployment)
	pointer.service.POST("DeleteKubernetesDeployment", "v1", "delete/kubernetes/deployment", 3, true, "删除一个kubernetes部署", pointer.DeleteKubernetesDeployment)
	pointer.service.GET("GetKubernetesPodLog", "v1", "get/kubernetes/pod/log", 3, true, "获取kubernetes的pod日志", pointer.GetKubernetesPodLog)
}

//write 实现
type write struct {
	read func([]byte)
}

func newWrite(read func([]byte)) *write {
	return &write{read: read}
}

//Write 写实现
func (w *write) Write(p []byte) (int, error) {
	if w.read != nil {
		w.read(p)
	}
	return len(p), nil
}

//buildEvent 编译事件
type buildEvent struct {
	ctx     plugins.Context
	request *param.BuildServiceReq
	buildID int64
}

//API 控制服务
type API struct {
	service    plugins.Service
	cfg        *cfg.Config
	docker     *client.Client
	mysql      *mysql.Mysql
	etcd       *etcd.Etcd
	snowflake  *snowflake.SnowFlake
	cache      *cache.Cache
	clientSet  *kubernetes.Clientset
	buildEvent chan *buildEvent
	closeEvent chan bool
	wait       sync.WaitGroup
}

//NewService 初始化服务
func NewService() *API {
	ctl := new(API)
	//初始化日志
	ctl.cfg = cfg.NewConfig()
	//初始化k8s
	config, err := clientcmd.BuildConfigFromFlags("", "./config/admin.conf")
	if err != nil {
		panic(err.Error())
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	ctl.clientSet = clientSet
	//初始化docker
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
	//初始化etcd
	ctl.etcd = etcd.NewEtcd(ctl.cfg)
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
	//启动事件执行器
	ctl.buildEvent = make(chan *buildEvent)
	ctl.closeEvent = make(chan bool)
	//启动
	go ctl._EventExecution()
	return ctl
}

//Run 启动
func (pointer *API) Run() error {
	err := pointer.service.Run()
	pointer.etcd.Close()
	pointer.wait.Wait()
	pointer.closeEvent <- true
	return err
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
	err := pointer.docker.ContainerStop(context.Background(), id, nil)
	if err != nil {
		return err
	}
	err = pointer.docker.ContainerRemove(context.Background(), id, types.ContainerRemoveOptions{})
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
	out, err := pointer.docker.ImagePush(context.Background(), image, config)
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
	out, err := pointer.docker.ImagePull(context.Background(), image, config)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(out)
	for scanner.Scan() {
		success(scanner.Text())
	}
	return nil
}

//_CreateTar 打包tar
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

//_TarFile 打包文件
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

//_TarFolder 打包文件架
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

//_SendEvent 发送事件
func (pointer *API) _SendEvent(id int64, ctx plugins.Context, request *param.BuildServiceReq) {
	pointer.wait.Add(1)
	pointer._PuseMsgToAdmin(ctx.GetToken(), define.BuildServiceTopic, "正在排队进行编译，请稍等")
	pointer.buildEvent <- &buildEvent{
		buildID: id,
		ctx:     ctx,
		request: request,
	}
}

//_EventExecution 事件执行队列
func (pointer *API) _EventExecution() {
	for {
		select {
		case event := <-pointer.buildEvent:
			request := event.request
			ctx := event.ctx
			paths := strings.Split(request.Git, "/")
			name := strings.Replace(paths[len(paths)-1], ".git", "", -1)
			image := request.Harbor + `/` + request.Name + `:` + request.Version
			logTxt := ""

			pointer._PuseMsgToAdmin(ctx.GetToken(), define.BuildServiceTopic, "开始编译 "+request.Git)
			logTxt = logTxt + `开始编译<p/>`
			if len(paths) <= 0 {
				pointer._PuseMsgToAdmin(ctx.GetToken(), define.BuildServiceTopic, "编译路径不正确")
				logTxt = logTxt + `路径不正确<p/>`
			} else {
				if _, e := git.PlainClone(name, false, &git.CloneOptions{
					URL: request.Git,
					Progress: newWrite(func(b []byte) {
						pointer._PuseMsgToAdmin(ctx.GetToken(), define.BuildServiceTopic, string(b))
						logTxt = logTxt + string(b) + `<p/>`
					}),
					Depth: 1,
					Auth: &http.BasicAuth{
						Username: "tangjiework@outlook.com",
						Password: "tangjie520@",
					},
				}); e == nil {
					system := runtime.GOOS
					build := ""
					log.Traceln("当前系统环境", system)
					switch system {
					case "darwin":
						build = "CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o " + request.Name
					case "linxu":
						build = "go build -o " + request.Name
					default:
						build = "CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o " + request.Name
					}
					//开始编译代码
					shell := `
					cd ` + name + `
					go mod tidy
					cd ` + request.Path + `
					` + build
					pointer._RunInLinux(shell, func(success string) {
						pointer._PuseMsgToAdmin(ctx.GetToken(), define.BuildServiceTopic, success)
						logTxt = logTxt + success + `<p/>`
					}, func(err string) {
						pointer._PuseMsgToAdmin(ctx.GetToken(), define.BuildServiceTopic, err)
						logTxt = logTxt + err + `<p/>`
					})
					path := fmt.Sprintf("./%s/%s", name, request.Path)
					tarName := uuid.GetToken() + ".tar"
					//打包
					if e := pointer._CreateTar(path, tarName, false); e != nil {
						pointer._PuseMsgToAdmin(ctx.GetToken(), define.BuildServiceTopic, e.Error())
						logTxt = logTxt + e.Error() + `<p/>`
					} else {
						//编译镜像
						if err := pointer._BuildImage("./"+tarName, "", image, func(res string) {
							pointer._PuseMsgToAdmin(ctx.GetToken(), define.BuildServiceTopic, res)
							logTxt = logTxt + res + `<p/>`
						}); err != nil {
							pointer._PuseMsgToAdmin(ctx.GetToken(), define.BuildServiceTopic, err.Error())
							logTxt = logTxt + err.Error() + `<p/>`
						}
						//执行push
						if e := pointer._PushImage(request.Accouunt, request.Pwd, image, func(res string) {
							pointer._PuseMsgToAdmin(ctx.GetToken(), define.BuildServiceTopic, res)
							logTxt = logTxt + res + `<p/>`
						}); e != nil {
							pointer._PuseMsgToAdmin(ctx.GetToken(), define.BuildServiceTopic, e.Error())
							logTxt = logTxt + e.Error() + `<p/>`
						}
					}
					//删除执行文件夹
					os.RemoveAll(name)
					os.RemoveAll(tarName)
					pointer._PuseMsgToAdmin(ctx.GetToken(), define.BuildServiceTopic, "执行完成")
					logTxt = logTxt + "执行完成" + `<p/>`
				} else {
					log.Errorln(e.Error())
					pointer._PuseMsgToAdmin(ctx.GetToken(), define.BuildServiceTopic, e.Error())
					logTxt = logTxt + e.Error() + `<p/>`
				}
			}
			//完成
			err := pointer.mysql.GetWriteEngine().Model(&table.BuildService{}).Where("id = ?", event.buildID).Update(
				map[string]interface{}{
					"Log":    logTxt,
					"Status": true,
				}).Error
			if err != nil {
				log.Errorln(err.Error())
			}
			pointer.wait.Done()
		case <-pointer.closeEvent:
			return
		}
	}
}

//_PuseMsgToAdmin 给admin推送消息
func (pointer *API) _PuseMsgToAdmin(token, topic, msg string) error {
	admin := new(table.Admin)
	if e := pointer.cache.GetCache().Get(token, admin); e != nil {
		//用户下线了,不推送任何消息
		return errors.New("token fail")
	}
	ctx := context.Background()
	ctx.SetIsTest(false)
	ctx.SetTraceID(uuid.GetToken())
	ctx.SetToken(token)
	if e := pointer.service.GetClient().CallByAddress(context.WithTimeout(ctx, int64(time.Second*time.Duration(6))), admin.GateAddress, define.SvcGateWay, "Push", &gateParam.PushReq{
		Token: token,
		Topic: topic,
		Msg:   msg,
	}, &gateParam.PushRes{}); e != nil {
		log.Warnln(e.Error())
		return e
	}
	return nil
}
