package service

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
	authAPI "github.com/tang-go/go-dog-tool/go-dog-auth/api"
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
func (s *Service) Router() {

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

//Service 控制服务
type Service struct {
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
func NewService(routers ...func(*Service)) *Service {
	ctl := new(Service)
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
		table.BuildService{},
		table.Docker{},
		table.Log{},
	)
	//初始化缓存
	ctl.cache = cache.NewCache(ctl.cfg)
	//初始化etcd
	//ctl.etcd = etcd.NewEtcd(ctl.cfg)
	//初始化雪花算法
	ret := big.NewInt(0)
	ret.SetBytes(net.ParseIP(ctl.cfg.GetHost()).To4())
	id, err := strconv.ParseInt(fmt.Sprintf("%d%d", ret.Int64(), ctl.cfg.GetPort()), 10, 64)
	if err != nil {
		panic(err)
	}
	ctl.snowflake = snowflake.NewSnowFlake(id)
	//初始化Service
	//初始化路由
	for _, router := range routers {
		router(ctl)
	}
	//初始化数据库数据
	ctl._InitMysql(ctl.cfg.GetPhone(), ctl.cfg.GetPwd())
	//启动事件执行器
	ctl.buildEvent = make(chan *buildEvent)
	ctl.closeEvent = make(chan bool)
	//启动
	go ctl._EventExecution()
	return ctl
}

//RPC 	注册RPC方法
//name			方法名称
//level			方法等级
//isAuth		是否需要鉴权
//explain		方法说明
//fn			注册的方法
func (s *Service) RPC(name string, level int8, isAuth bool, explain string, fn interface{}) {
	s.service.RPC(name, level, isAuth, explain, fn)
}

//POST 			注册POST方法
//methodname 	Service方法名称
//version 		Service方法版本
//path 			Service路由
//level 		Service等级
//isAuth 		是否需要鉴权
//explain		方法描述
//fn 			注册的方法
func (s *Service) POST(methodname, version, path string, level int8, isAuth bool, explain string, fn interface{}) {
	ctx := context.WithTimeout(context.Background(), int64(time.Second*time.Duration(6)))
	ctx.SetClient(s.service.GetClient())
	if _, err := authAPI.CreateAPI(
		ctx,
		define.Organize,
		explain,
		fmt.Sprintf("api/%s/%s/%s", define.SvcController, version, path)); err != nil {
		panic(err.Error())
	}
	s.service.POST(methodname, version, path, level, isAuth, explain, fn)
}

//GET GET方法
//methodname 	Service方法名称
//version 		Service方法版本
//path 			Service路由
//level 		Service等级
//isAuth 		是否需要鉴权
//explain		方法描述
//fn 			注册的方法
func (s *Service) GET(methodname, version, path string, level int8, isAuth bool, explain string, fn interface{}) {
	ctx := context.WithTimeout(context.Background(), int64(time.Second*time.Duration(6)))
	ctx.SetClient(s.service.GetClient())
	if _, err := authAPI.CreateAPI(
		ctx,
		define.Organize,
		explain,
		fmt.Sprintf("api/%s/%s/%s", define.SvcController, version, path)); err != nil {
		panic(err.Error())
	}
	s.service.GET(methodname, version, path, level, isAuth, explain, fn)
}

//Run 启动
func (s *Service) Run() error {
	err := s.service.Run()
	s.etcd.Close()
	s.wait.Wait()
	s.closeEvent <- true
	return err
}

//_InitMysql 第一次加载初始化数据库数据
func (s *Service) _InitMysql(phone, pwd string) {
	//添加系统菜单
	ctx := context.WithTimeout(context.Background(), int64(time.Second*time.Duration(10)))
	ctx.SetClient(s.service.GetClient())
	_, err := authAPI.CreateMenu(ctx, define.Organize, "首页", "/index", 0, 100000)
	if err != nil {
		panic(err.Error())
	}
	powerID, err := authAPI.CreateMenu(ctx, define.Organize, "权限管理", "/power", 0, 0)
	if err != nil {
		panic(err.Error())
	}
	_, err = authAPI.CreateMenu(ctx, define.Organize, "菜单管理", "/power/menu", powerID, 0)
	if err != nil {
		panic(err.Error())
	}
	//读取是否有业主了
	owner := new(table.Owner)
	if s.mysql.GetReadEngine().Where("phone = ?", phone).First(owner).RecordNotFound() == false {
		return
	}
	//调用权限服务创建一个角色
	roleID, err := authAPI.CreateRole(ctx, define.Organize, "超级管理员", "超级管理员", true)
	if err != nil {
		panic(err.Error())
	}
	//如果没有业主则新增默认业主
	owner.OwnerID = s.snowflake.GetID()
	owner.Name = "系统业主"
	owner.Phone = phone
	owner.Level = 1
	owner.IsDisable = table.OwnerAvailable
	owner.IsAdminOwner = table.IsAdminOwner
	owner.RoleID = roleID
	owner.Time = time.Now().Unix()
	//管理员
	admin := &table.Admin{
		//账号 唯一主键
		AdminID: s.snowflake.GetID(),
		//名称
		Name: "admin",
		//电话
		Phone: phone,
		//盐值 md5使用
		Salt: rand.StringRand(6),
		//所属业主
		OwnerID: owner.OwnerID,
		//是否被禁用
		IsDisable: table.AdminAvailable,
		//管理员绑定角色
		RoleID: roleID,
		//注册事件
		Time: owner.Time,
	}
	//生成密码
	admin.Pwd = md5.Md5(md5.Md5(pwd) + admin.Salt)
	//开启数据库操作
	tx := s.mysql.GetWriteEngine().Begin()
	if err := tx.Create(owner).Error; err != nil {
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
func (s *Service) Set(id string, value string) {
	if err := s.cache.GetCache().SetByTime(id, value, define.CodeValidityTime); err != nil {
		log.Errorln(err.Error())
	}
}

// Get 更具验证ID获取验证码
func (s *Service) Get(id string, clear bool) (vali string) {
	err := s.cache.GetCache().Get(id, &vali)
	if err != nil {
		log.Errorln(err.Error())
	}
	if clear {
		s.cache.GetCache().Del(id)
	}
	return
}

//Verify 验证验证码
func (s *Service) Verify(id, answer string, clear bool) bool {
	vali := s.Get(id, clear)
	if strings.ToLower(vali) != strings.ToLower(answer) {
		return false
	}
	return true
}

func (s *Service) _RunInLinux(cmd string, success func(string), fail func(string)) error {
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

func (s *Service) _PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	log.Errorln(err.Error())
	return false
}

//BuildImage 编译镜像
func (s *Service) _BuildImage(tarFile, project, imageName string, success func(string)) error {
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
	output, err := s.docker.ImageBuild(context.Background(), dockerBuildContext, buildOptions)
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
func (s *Service) _CloseDocker(id string) error {
	err := s.docker.ContainerStop(context.Background(), id, nil)
	if err != nil {
		return err
	}
	err = s.docker.ContainerRemove(context.Background(), id, types.ContainerRemoveOptions{})
	if err != nil {
		return err
	}
	return nil
}

//PushImage 推送镜像
func (s *Service) _PushImage(registryUser, registryPassword, image string, success func(string)) error {
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
	out, err := s.docker.ImagePush(context.Background(), image, config)
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
func (s *Service) _PullImage(registryUser, registryPassword, image string, success func(string)) error {
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
	out, err := s.docker.ImagePull(context.Background(), image, config)
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
func (s *Service) _CreateTar(filesource, filetarget string, deleteIfExist bool) error {
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
			if err := s._TarFile(f.Name(), path, sfileInfo, tarwriter); err != nil {
				log.Errorln(err.Error())
				return err
			}
		} else {
			if err := s._TarFolder(path, tarwriter); err != nil {
				log.Errorln(err.Error())
				return err
			}
		}
	}
	return nil
}

//_TarFile 打包文件
func (s *Service) _TarFile(directory string, filesource string, sfileInfo os.FileInfo, tarwriter *tar.Writer) error {
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
func (s *Service) _TarFolder(directory string, tarwriter *tar.Writer) error {
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
		return s._TarFile(fileFolder, targetpath, file, tarwriter)
	})
}

//_SendEvent 发送事件
func (s *Service) _SendEvent(id int64, ctx plugins.Context, request *param.BuildServiceReq) {
	s.wait.Add(1)
	s._PuseMsgToAdmin(ctx.GetToken(), define.BuildServiceTopic, "正在排队进行编译，请稍等")
	s.buildEvent <- &buildEvent{
		buildID: id,
		ctx:     ctx,
		request: request,
	}
}

//_EventExecution 事件执行队列
func (s *Service) _EventExecution() {
	for {
		select {
		case event := <-s.buildEvent:
			request := event.request
			ctx := event.ctx
			paths := strings.Split(request.Git, "/")
			name := strings.Replace(paths[len(paths)-1], ".git", "", -1)
			image := request.Harbor + `/` + request.Name + `:` + request.Version
			logTxt := ""

			s._PuseMsgToAdmin(ctx.GetToken(), define.BuildServiceTopic, "开始编译 "+request.Git)
			logTxt = logTxt + `开始编译<p/>`
			if len(paths) <= 0 {
				s._PuseMsgToAdmin(ctx.GetToken(), define.BuildServiceTopic, "编译路径不正确")
				logTxt = logTxt + `路径不正确<p/>`
			} else {
				if s._PathExists(name) {
					s._PuseMsgToAdmin(ctx.GetToken(), define.BuildServiceTopic, "开始git pull")
					logTxt = logTxt + "开始git pull" + `<p/>`
					log.Traceln("开始git pull")
					if r, err := git.PlainOpen(name); err != nil {
						s._PuseMsgToAdmin(ctx.GetToken(), define.BuildServiceTopic, err.Error())
						logTxt = logTxt + err.Error() + `<p/>`
						log.Traceln(err.Error())
					} else {
						if w, err := r.Worktree(); err != nil {
							s._PuseMsgToAdmin(ctx.GetToken(), define.BuildServiceTopic, err.Error())
							logTxt = logTxt + err.Error() + `<p/>`
							log.Traceln(err.Error())
						} else {
							if err = w.Pull(&git.PullOptions{
								RemoteName:        "main",
								RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
								Depth:             1,
								Progress: newWrite(func(b []byte) {
									s._PuseMsgToAdmin(ctx.GetToken(), define.BuildServiceTopic, string(b))
									logTxt = logTxt + string(b) + `<p/>`
									log.Traceln(string(b))
								}),
								Auth: &http.BasicAuth{
									Username: request.GitAccount,
									Password: request.GitPwd,
								},
							}); err != nil && err.Error() != git.NoErrAlreadyUpToDate.Error() {
								s._PuseMsgToAdmin(ctx.GetToken(), define.BuildServiceTopic, err.Error())
								logTxt = logTxt + err.Error() + `<p/>`
								log.Traceln(err.Error())
							} else {
								// Print the latest commit that was just pulled
								ref, err := r.Head()
								if err != nil {
									log.Errorln(err.Error())
								}
								commit, err := r.CommitObject(ref.Hash())
								if err != nil {
									log.Errorln(err.Error())
								}
								log.Traceln("pull commit %v", commit)
								system := runtime.GOOS
								build := ""
								s._PuseMsgToAdmin(ctx.GetToken(), define.BuildServiceTopic, "开始 go build 系统"+system)
								logTxt = logTxt + "开始 go build 系统" + system + `<p/>`
								log.Traceln("开始 go build 系统" + system)
								switch system {
								case "darwin":
									build = "CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o " + request.Name
								case "linxu":
									build = "CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o " + request.Name
								default:
									build = "CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o " + request.Name
								}
								//开始编译代码
								shell := `
								cd ` + name + `
								go mod tidy
								cd ` + request.Path + `
								` + build
								s._RunInLinux(shell, func(success string) {
									s._PuseMsgToAdmin(ctx.GetToken(), define.BuildServiceTopic, success)
									logTxt = logTxt + success + `<p/>`
									log.Traceln(success)
								}, func(err string) {
									s._PuseMsgToAdmin(ctx.GetToken(), define.BuildServiceTopic, err)
									logTxt = logTxt + err + `<p/>`
									log.Traceln(err)
								})
								path := fmt.Sprintf("./%s/%s", name, request.Path)
								tarName := uuid.GetToken() + ".tar"

								s._PuseMsgToAdmin(ctx.GetToken(), define.BuildServiceTopic, "打包 tar")
								logTxt = logTxt + "打包 tar" + `<p/>`
								log.Traceln("打包 tar")
								if e := s._CreateTar(path, tarName, false); e != nil {
									s._PuseMsgToAdmin(ctx.GetToken(), define.BuildServiceTopic, e.Error())
									logTxt = logTxt + e.Error() + `<p/>`
									log.Traceln(e.Error())
								} else {
									s._PuseMsgToAdmin(ctx.GetToken(), define.BuildServiceTopic, "开始编译镜像")
									logTxt = logTxt + "开始编译镜像" + `<p/>`
									log.Traceln("开始编译镜像")
									if err := s._BuildImage("./"+tarName, "", image, func(res string) {
										s._PuseMsgToAdmin(ctx.GetToken(), define.BuildServiceTopic, res)
										logTxt = logTxt + res + `<p/>`
										log.Traceln(res)
									}); err != nil {
										s._PuseMsgToAdmin(ctx.GetToken(), define.BuildServiceTopic, err.Error())
										logTxt = logTxt + err.Error() + `<p/>`
										log.Traceln(err.Error())
									}
									s._PuseMsgToAdmin(ctx.GetToken(), define.BuildServiceTopic, "开始推送镜像")
									logTxt = logTxt + "开始推送镜像" + `<p/>`
									log.Traceln("开始推送镜像")
									if e := s._PushImage(request.Accouunt, request.Pwd, image, func(res string) {
										s._PuseMsgToAdmin(ctx.GetToken(), define.BuildServiceTopic, res)
										logTxt = logTxt + res + `<p/>`
										log.Traceln(res)
									}); e != nil {
										s._PuseMsgToAdmin(ctx.GetToken(), define.BuildServiceTopic, e.Error())
										logTxt = logTxt + e.Error() + `<p/>`
										log.Traceln(e.Error())
									}
								}
								s._PuseMsgToAdmin(ctx.GetToken(), define.BuildServiceTopic, "删除执行文件夹")
								logTxt = logTxt + "删除执行文件夹" + `<p/>`
								log.Traceln("删除执行文件夹")
								os.RemoveAll(fmt.Sprintf("%s/%s/%s", name, request.Path, request.Name))
								os.RemoveAll(tarName)
							}
						}
					}
				} else {
					s._PuseMsgToAdmin(ctx.GetToken(), define.BuildServiceTopic, "开始git clone")
					logTxt = logTxt + "开始git clone" + `<p/>`
					log.Traceln("开始git clone")
					if r, e := git.PlainClone(name, false, &git.CloneOptions{
						URL: request.Git,
						Progress: newWrite(func(b []byte) {
							s._PuseMsgToAdmin(ctx.GetToken(), define.BuildServiceTopic, string(b))
							logTxt = logTxt + string(b) + `<p/>`
							log.Traceln(string(b))
						}),
						RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
						Depth:             1,
						Auth: &http.BasicAuth{
							Username: request.GitAccount,
							Password: request.GitPwd,
						},
					}); e == nil {
						ref, err := r.Head()
						if err != nil {
							log.Errorln(err.Error())
						}
						commit, err := r.CommitObject(ref.Hash())
						if err != nil {
							log.Errorln(err.Error())
						}
						log.Traceln("commit %v", commit)
						system := runtime.GOOS
						build := ""
						s._PuseMsgToAdmin(ctx.GetToken(), define.BuildServiceTopic, "开始 go build 系统"+system)
						logTxt = logTxt + "开始 go build 系统" + system + `<p/>`
						log.Traceln("开始 go build 系统" + system)
						switch system {
						case "darwin":
							build = "CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o " + request.Name
						case "linxu":
							build = "CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o " + request.Name
						default:
							build = "CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o " + request.Name
						}
						//开始编译代码
						shell := `
						cd ` + name + `
						go mod tidy
						cd ` + request.Path + `
						` + build
						s._RunInLinux(shell, func(success string) {
							s._PuseMsgToAdmin(ctx.GetToken(), define.BuildServiceTopic, success)
							logTxt = logTxt + success + `<p/>`
							log.Traceln(success)
						}, func(err string) {
							s._PuseMsgToAdmin(ctx.GetToken(), define.BuildServiceTopic, err)
							logTxt = logTxt + err + `<p/>`
							log.Traceln(err)
						})
						path := fmt.Sprintf("./%s/%s", name, request.Path)
						tarName := uuid.GetToken() + ".tar"

						s._PuseMsgToAdmin(ctx.GetToken(), define.BuildServiceTopic, "打包 tar")
						logTxt = logTxt + "打包 tar" + `<p/>`
						log.Traceln("打包 tar")
						if e := s._CreateTar(path, tarName, false); e != nil {
							s._PuseMsgToAdmin(ctx.GetToken(), define.BuildServiceTopic, e.Error())
							logTxt = logTxt + e.Error() + `<p/>`
							log.Traceln(e.Error())
						} else {
							s._PuseMsgToAdmin(ctx.GetToken(), define.BuildServiceTopic, "开始编译镜像")
							logTxt = logTxt + "开始编译镜像" + `<p/>`
							log.Traceln("开始编译镜像")
							if err := s._BuildImage("./"+tarName, "", image, func(res string) {
								s._PuseMsgToAdmin(ctx.GetToken(), define.BuildServiceTopic, res)
								logTxt = logTxt + res + `<p/>`
								log.Traceln(res)
							}); err != nil {
								s._PuseMsgToAdmin(ctx.GetToken(), define.BuildServiceTopic, err.Error())
								logTxt = logTxt + err.Error() + `<p/>`
								log.Traceln(err.Error())
							}
							s._PuseMsgToAdmin(ctx.GetToken(), define.BuildServiceTopic, "开始推送镜像")
							logTxt = logTxt + "开始推送镜像" + `<p/>`
							log.Traceln("开始推送镜像")
							if e := s._PushImage(request.Accouunt, request.Pwd, image, func(res string) {
								s._PuseMsgToAdmin(ctx.GetToken(), define.BuildServiceTopic, res)
								logTxt = logTxt + res + `<p/>`
								log.Traceln(res)
							}); e != nil {
								s._PuseMsgToAdmin(ctx.GetToken(), define.BuildServiceTopic, e.Error())
								logTxt = logTxt + e.Error() + `<p/>`
								log.Traceln(e.Error())
							}
						}
						s._PuseMsgToAdmin(ctx.GetToken(), define.BuildServiceTopic, "删除执行文件夹")
						logTxt = logTxt + "删除执行文件夹" + `<p/>`
						log.Traceln("删除执行文件夹")
						os.RemoveAll(fmt.Sprintf("%s/%s/%s", name, request.Path, request.Name))
						os.RemoveAll(tarName)
					} else {
						log.Errorln(e.Error())
						s._PuseMsgToAdmin(ctx.GetToken(), define.BuildServiceTopic, e.Error())
						logTxt = logTxt + e.Error() + `<p/>`
					}
				}
			}
			//完成
			s._PuseMsgToAdmin(ctx.GetToken(), define.BuildServiceTopic, "执行完成")
			logTxt = logTxt + "执行完成" + `<p/>`
			log.Traceln("执行完成")
			err := s.mysql.GetWriteEngine().Model(&table.BuildService{}).Where("id = ?", event.buildID).Update(
				map[string]interface{}{
					"Log":    logTxt,
					"Status": true,
				}).Error
			if err != nil {
				log.Errorln(err.Error())
			}
			s.wait.Done()
		case <-s.closeEvent:
			return
		}
	}
}

//_PuseMsgToAdmin 给admin推送消息
func (s *Service) _PuseMsgToAdmin(token, topic, msg string) error {
	admin := new(table.Admin)
	if e := s.cache.GetCache().Get(token, admin); e != nil {
		//用户下线了,不推送任何消息
		return errors.New("token fail")
	}
	ctx := context.Background()
	ctx.SetIsTest(false)
	ctx.SetTraceID(uuid.GetToken())
	ctx.SetToken(token)
	if e := s.service.GetClient().CallByAddress(context.WithTimeout(ctx, int64(time.Second*time.Duration(6))), admin.GateAddress, define.SvcGateWay, "Push", &gateParam.PushReq{
		Token: token,
		Topic: topic,
		Msg:   msg,
	}, &gateParam.PushRes{}); e != nil {
		log.Warnln(e.Error())
		return e
	}
	return nil
}
