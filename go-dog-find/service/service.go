package service

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/tang-go/go-dog-tool/go-dog-find/param"
	"github.com/tang-go/go-dog-tool/go-dog-find/redis"
	"github.com/tang-go/go-dog/lib/io"
	"github.com/tang-go/go-dog/lib/uuid"
	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/pkg/config"
	"github.com/tang-go/go-dog/plugins"
)

//Service 控制服务
type Service struct {
	cfg        plugins.Cfg
	redis      *redis.Redis
	registers  sync.Map
	discoverys sync.Map
	lock       sync.RWMutex
}

//NewService 初始化服务
func NewService() *Service {
	s := &Service{
		cfg: config.NewConfig(),
	}
	s.redis = redis.NewReids(s.cfg)
	return s
}

//Run 启动
func (s *Service) Run() error {
	c := make(chan os.Signal)
	//监听指定信号
	signal.Notify(c, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)
	go func() {
		err := s.run()
		if err != nil {
			panic(err.Error())
		}
	}()
	msg := <-c
	return fmt.Errorf("收到kill信号:%s", msg)
}

//run
func (s *Service) run() error {
	address := fmt.Sprintf("0.0.0.0:%d", s.cfg.GetPort())
	l, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	log.Traceln("服务发现中心启动:", address)
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Errorln(err.Error())
			continue
		}
		go s.serveConn(conn)
	}
}

// ServeConn 拦截一个链接
func (s *Service) serveConn(conn net.Conn) {
	defer conn.Close()
	//读取第一个事件
	_, buff, err := io.ReadByTime(conn, time.Now().Add(time.Second*5))
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	event := new(param.Event)
	//解析事件
	if err := event.DeCode(buff, event); err != nil {
		log.Errorln(err.Error())
		return
	}
	switch event.Cmd {
	//登陆事件
	case param.Login:
		login := new(param.LoginReq)
		if err := login.DeCode(event.Data, login); err != nil {
			log.Errorln(err.Error())
			return
		}
		switch login.Type {
		//服务注册类型
		case param.RegType:
			token := uuid.GetToken()
			register := NewRegister(s, conn)
			s.registers.Store(token, register)
			register.Run()
			s.registers.Delete(token)
		//服务发现类型
		case param.DisType:
			token := uuid.GetToken()
			discovery := NewDiscovery(s, conn)
			s.discoverys.Store(token, discovery)
			discovery.Run()
			s.discoverys.Delete(token)
		//不合法类型
		default:
			log.Errorln("类型不合法")
			return
		}
	//不合法事件
	default:
		log.Errorln("事件不合法")
		return
	}
}
