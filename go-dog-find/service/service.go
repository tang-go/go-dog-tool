package service

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/tang-go/go-dog-tool/go-dog-find/gossip"
	"github.com/tang-go/go-dog-tool/go-dog-find/param"
	"github.com/tang-go/go-dog/lib/io"
	"github.com/tang-go/go-dog/lib/uuid"
	"github.com/tang-go/go-dog/log"
	"github.com/vmihailenco/msgpack"
)

const (
	_Add int8 = iota
	_Get
	_Del
)

//Update 数据同步
type Update struct {
	Action int8
	Label  string
	Data   *param.Data
}

//Item 内容
type Item struct {
	Datas map[string]*param.Data
}

//Service 控制服务
type Service struct {
	cfg        *Config
	gossip     *gossip.Gossip
	registers  map[string]*Register
	discoverys map[string]*Discovery
	items      map[string]*Item
	broadcasts [][]byte
	close      int32
	lock       sync.RWMutex
}

//NewService 初始化服务
func NewService() *Service {
	s := &Service{
		close:      0,
		registers:  make(map[string]*Register),
		discoverys: make(map[string]*Discovery),
		items:      make(map[string]*Item),
	}
	s.cfg = NewConfig()
	//创建gossip对象
	s.gossip = gossip.NewGossip(s.cfg.GetGossipPort(), s, s.cfg.GetMembers())
	//初始化日志
	switch s.cfg.GetRunmode() {
	case "panic":
		log.SetLevel(log.PanicLevel)
		break
	case "fatal":
		log.SetLevel(log.FatalLevel)
		break
	case "error":
		log.SetLevel(log.ErrorLevel)
		break
	case "warn":
		log.SetLevel(log.WarnLevel)
		break
	case "info":
		log.SetLevel(log.InfoLevel)
		break
	case "debug":
		log.SetLevel(log.DebugLevel)
		break
	case "trace":
		log.SetLevel(log.TraceLevel)
		break
	default:
		log.SetLevel(log.TraceLevel)
		break
	}
	return s
}

//Add 增加
func (s *Service) Add(label string, data *param.Data) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	item, ok := s.items[label]
	if ok {
		if _, o := item.Datas[data.Key]; o {
			return nil
		}
		item.Datas[data.Key] = data
	} else {
		//不存在
		s.items[label] = &Item{
			Datas: map[string]*param.Data{
				data.Key: data,
			},
		}
	}
	//储存当前指令
	b, err := msgpack.Marshal(&Update{
		Action: _Add,
		Label:  label,
		Data:   data,
	})
	if err != nil {
		log.Errorln(err.Error())
		return err
	}
	s.broadcasts = append(s.broadcasts, b)
	log.Traceln("上线", label, data.Key)
	return nil
}

//Del 删除
func (s *Service) Del(label string, data *param.Data) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	item, ok := s.items[label]
	if ok {
		delete(item.Datas, data.Key)
		if len(item.Datas) <= 0 {
			delete(s.items, label)
		}
	}
	b, err := msgpack.Marshal(&Update{
		Action: _Del,
		Label:  label,
		Data:   data,
	})
	if err != nil {
		log.Errorln(err.Error())
		return err
	}
	s.broadcasts = append(s.broadcasts, b)
	log.Traceln("下线", label, data.Key)
	return nil
}

//Get 获取
func (s *Service) Get(label string) (datas []*param.Data, o bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	val, ok := s.items[label]
	if !ok {
		return nil, false
	}
	for _, data := range val.Datas {
		datas = append(datas, data)
	}
	return datas, true
}

//NodeMeta 节点数据
func (s *Service) NodeMeta(limit int) []byte {
	return []byte("test")
}

//NotifyMsg 操作
func (s *Service) NotifyMsg(b []byte) {
	if len(b) == 0 {
		return
	}
	update := new(Update)
	if err := msgpack.Unmarshal(b, &update); err != nil {
		log.Errorln(err.Error())
		return
	}
	s.lock.Lock()
	defer s.lock.Unlock()
	label := update.Label
	data := update.Data
	if update.Action == _Add {
		item, ok := s.items[label]
		if ok {
			if _, o := item.Datas[data.Key]; o {
				return
			}
			item.Datas[data.Key] = data
		} else {
			s.items[label] = &Item{
				Datas: map[string]*param.Data{
					data.Key: data,
				},
			}
		}
		log.Traceln("上线", label, data.Key)
	}
	if update.Action == _Del {
		item, ok := s.items[label]
		if ok {
			delete(item.Datas, data.Key)
			if len(item.Datas) <= 0 {
				delete(s.items, label)
			}
		}
		log.Traceln("下线", label, data.Key)
	}
}

//GetBroadcasts 广播
func (s *Service) GetBroadcasts(overhead, limit int) [][]byte {
	s.lock.Lock()
	buffs := s.broadcasts
	s.broadcasts = make([][]byte, 0)
	s.lock.Unlock()
	return buffs
}

//LocalState tcp推拉数据接口
func (s *Service) LocalState(join bool) []byte {
	s.lock.RLock()
	m := s.items
	s.lock.RUnlock()
	b, _ := msgpack.Marshal(m)
	return b
}

//MergeRemoteState 从LocalState接口推拉过来的数据 同步
func (s *Service) MergeRemoteState(buf []byte, join bool) {
	if len(buf) == 0 {
		return
	}
	s.lock.Lock()
	defer s.lock.Unlock()
	items := make(map[string]*Item)
	msgpack.Unmarshal(buf, &items)
	for label, item := range items {
		if i, ok := s.items[label]; ok {
			for _, data := range item.Datas {

			}
		} else {
			s.items[label] = items[label]
		}

	}
	log.Traceln("全体同步", join)
}

//Run 启动
func (s *Service) Run() error {
	c := make(chan os.Signal)
	//监听指定信号
	signal.Notify(c, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		err := s._Run()
		if err != nil {
			panic(err.Error())
		}
	}()
	msg := <-c
	s.Close()
	return fmt.Errorf("收到kill信号:%s", msg)
}

//Close 关闭服务
func (s *Service) Close() {
	atomic.AddInt32(&s.close, 1)
}

//_Run
func (s *Service) _Run() error {
	address := fmt.Sprintf("0.0.0.0:%d", s.cfg.GetPort())
	l, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	log.Traceln("服务发现中心启动:", address)
	defer l.Close()
	for {
		if atomic.LoadInt32(&s.close) > 0 {
			return nil
		}
		conn, err := l.Accept()
		if err != nil {
			log.Errorln(err.Error())
			continue
		}
		go s._ServeConn(conn)
	}
}

// ServeConn 拦截一个链接
func (s *Service) _ServeConn(conn net.Conn) {
	_, buff, err := io.ReadByTime(conn, time.Now().Add(time.Second*5))
	if err != nil {
		conn.Close()
		log.Errorln(err.Error())
		return
	}
	login := new(param.LoginReq)
	if err := login.DeCode(buff, login); err == nil {
		//服务发现
		if login.Type == param.DisType {
			token := uuid.GetToken()
			discovery := NewDiscovery(s, conn, func() {
				s.lock.Lock()
				delete(s.discoverys, token)
				s.lock.Unlock()
			})
			go discovery.Run()
			s.lock.Lock()
			s.discoverys[token] = discovery
			s.lock.Unlock()
		}
		//服务注册
		if login.Type == param.RegType {
			token := uuid.GetToken()
			register := NewRegister(s, conn, func() {
				s.lock.Lock()
				delete(s.registers, token)
				s.lock.Unlock()
			})
			go register.Run()
			s.lock.Lock()
			s.registers[token] = register
			s.lock.Unlock()
		}
	}
}
