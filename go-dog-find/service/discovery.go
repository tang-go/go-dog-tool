package service

import (
	"net"
	"sync"
	"time"

	"github.com/tang-go/go-dog-tool/go-dog-find/param"
	"github.com/tang-go/go-dog/lib/io"
	"github.com/tang-go/go-dog/log"
)

//Discovery 服务发现
type Discovery struct {
	conn        net.Conn
	topics      map[string]string
	offlinefunc func()
	service     *Service
	close       chan bool
	lock        sync.RWMutex
}

//NewDiscovery 新建一个服务发现
func NewDiscovery(service *Service, conn net.Conn, offlinefunc func()) *Discovery {
	return &Discovery{
		conn:        conn,
		offlinefunc: offlinefunc,
		service:     service,
		close:       make(chan bool, 1),
		topics:      make(map[string]string),
	}
}

//Run 启动
func (d *Discovery) Run() {
	defer d.offlinefunc()
	go d.EventLoop()
	for {
		_, buff, err := io.ReadByTime(d.conn, time.Now().Add(time.Second*5))
		if err != nil {
			d.conn.Close()
			d.close <- true
			log.Errorln(err.Error())
			return
		}
		//处理各种事件
		event := new(param.Event)
		if err := event.DeCode(buff, event); err != nil {
			log.Errorln(err.Error())
			continue
		}
		switch event.Cmd {
		//监听事件
		case param.Listen:
			d.Listen(event.Label)
		}
	}
}

//Listen 监听
func (d *Discovery) Listen(label string) {
	d.lock.Lock()
	d.topics[label] = label
	log.Traceln("监听", label)
	d.lock.Unlock()
}

//PushEvent 推送事件
func (d *Discovery) PushEvent() {
	d.lock.RLock()
	for _, label := range d.topics {
		all := new(param.All)
		all.Label = label
		datas, ok := d.service.Get(label)
		if ok {
			for _, data := range datas {
				all.Datas = append(all.Datas, data)
			}
		}
		buff, err := all.EnCode(all)
		if err != nil {
			log.Errorln(err.Error())
			continue
		}
		//推送消息
		_, err = io.WriteByTime(d.conn, buff, time.Now().Add(time.Second*5))
		if err != nil {
			d.conn.Close()
			log.Errorln(err.Error())
			return
		}
	}
	d.lock.RUnlock()
}

//EventLoop 事件处理器
func (d *Discovery) EventLoop() {
	for {
		select {
		case <-d.close:
			close(d.close)
			return
		case <-time.After(time.Second * 2):
			//执行操作
			d.PushEvent()
		}
	}
}
