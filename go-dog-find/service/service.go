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
	Name   string
	Item   *Item
}

//Item 内容
type Item struct {
	Label   string
	Key     string
	Value   string
	Time    int64
	TimeOut int64
	Version uint64
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
	name       string
	exit       chan bool
	lock       sync.RWMutex
}

//NewService 初始化服务
func NewService() *Service {
	s := &Service{
		close:      0,
		registers:  make(map[string]*Register),
		discoverys: make(map[string]*Discovery),
		items:      make(map[string]*Item),
		exit:       make(chan bool),
		name:       uuid.GetToken(),
	}
	s.cfg = NewConfig()
	//创建gossip对象
	s.gossip = gossip.NewGossip(s.name, s.cfg.GetGossipPort(), s, s.cfg.GetMembers())
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
func (s *Service) Add(data *param.Data) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	update := new(Update)
	update.Action = _Add
	update.Name = s.name
	now := time.Now().Unix()
	if item, ok := s.items[data.Key]; ok {
		item.Time = now
		item.Version++
		update.Item = item
	} else {
		item := &Item{
			Label:   data.Label,
			Key:     data.Key,
			Value:   data.Value,
			Version: 1,
			Time:    now,
			TimeOut: 5,
		}
		update.Item = item
		s.items[data.Key] = item
		//log.Traceln("上线", data.Label, data.Key)
	}
	b, err := msgpack.Marshal(update)
	if err != nil {
		log.Errorln(err.Error())
		return err
	}
	s.broadcasts = append(s.broadcasts, b)
	return nil
}

//Del 删除
func (s *Service) Del(data *param.Data) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	if item, ok := s.items[data.Key]; ok {
		b, err := msgpack.Marshal(&Update{
			Action: _Del,
			Name:   s.name,
			Item:   item,
		})
		if err != nil {
			log.Errorln(err.Error())
			return err
		}
		delete(s.items, data.Key)
		s.broadcasts = append(s.broadcasts, b)
		//log.Traceln("下线", data.Label, data.Key)
	}
	return nil
}

//Get 获取
func (s *Service) Get(label string) (datas []*param.Data) {
	s.lock.RLock()
	items := s.items
	s.lock.RUnlock()
	for _, item := range items {
		if item.Label == label {
			datas = append(datas, &param.Data{
				Label: item.Label,
				Key:   item.Key,
				Value: item.Value,
			})
		}
	}
	return
}

//NodeMeta 节点数据
func (s *Service) NodeMeta(limit int) []byte {
	return []byte(s.name)
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
	if s.name == update.Name {
		return
	}
	s.lock.Lock()
	if update.Action == _Add {
		if item, ok := s.items[update.Item.Key]; ok {
			//比较版本
			if item.Version <= update.Item.Version {
				s.items[update.Item.Key] = update.Item
				s.broadcasts = append(s.broadcasts, b)
			}
		} else {
			s.items[update.Item.Key] = update.Item
			s.broadcasts = append(s.broadcasts, b)
			//log.Traceln("上线", update.Item.Label, update.Item.Key)
		}
	}
	if update.Action == _Del {
		if _, ok := s.items[update.Item.Key]; ok {
			delete(s.items, update.Item.Key)
			s.broadcasts = append(s.broadcasts, b)
			//log.Traceln("下线", update.Item.Label, update.Item.Key)
		}
	}
	s.lock.Unlock()
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
	if join == true {
		return nil
	}
	s.lock.RLock()
	b, _ := msgpack.Marshal(&s.items)
	s.lock.RUnlock()
	return b
}

//MergeRemoteState 从LocalState接口推拉过来的数据 同步
func (s *Service) MergeRemoteState(buf []byte, join bool) {
	if len(buf) == 0 {
		return
	}
	if join == true {
		return
	}
	items := make(map[string]*Item)
	if err := msgpack.Unmarshal(buf, &items); err != nil {
		return
	}
	s.lock.Lock()
	defer s.lock.Unlock()
	for key, i := range items {
		if item, ok := s.items[key]; ok {
			//比较版本
			if item.Version <= i.Version {
				s.items[key] = i
			}
		}
	}
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
	s.exit <- true
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
	go s._EventLoop()
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

//EventLoop 事件处理器
func (s *Service) _EventLoop() {
	for {
		select {
		case <-s.exit:
			close(s.exit)
			return
		case <-time.After(time.Second * 5):
			//执行操作
			s.lock.Lock()
			str := "在线服务:"
			now := time.Now().Unix()
			for key, item := range s.items {
				if item.Time+item.TimeOut < now {
					delete(s.items, key)
					continue
				}
				str = fmt.Sprintf("%s\n%s", str, item.Key)
			}
			s.lock.Unlock()
			log.Traceln(str)
		}
	}
}
