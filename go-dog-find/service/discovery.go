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
	conn    net.Conn
	topics  sync.Map
	service *Service
	close   chan bool
}

//NewDiscovery 新建一个服务发现
func NewDiscovery(service *Service, conn net.Conn) *Discovery {
	return &Discovery{
		conn:    conn,
		service: service,
		close:   make(chan bool, 1),
	}
}

//Run 启动
func (d *Discovery) Run() {
	go d._EventLoop()
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
			listen := new(param.ListenReq)
			if err := listen.DeCode(event.Data, listen); err != nil {
				log.Errorln(err.Error())
				d.conn.Close()
				continue
			}
			//推送第一个listen消息
			if err := d.PushEvent(listen.Label); err != nil {
				log.Errorln(err.Error())
				d.conn.Close()
				continue
			}
			d.topics.Store(listen.Label, listen.Label)
		}
	}
}

//PushEvent 推送事件
func (d *Discovery) PushEvent(label param.Label) error {
	listen := new(param.ListenRes)
	listen.Label = label
	//推送rpc消息
	if label == param.RPCLabel {
		d.service.redis.RangeRPC(func(key, value string) bool {
			listen.Data = append(listen.Data, param.Data{
				Key:   key,
				Value: value,
			})
			return true
		})
	}
	//推送api消息
	if label == param.APILabel {
		d.service.redis.RangeAPI(func(key, value string) bool {
			listen.Data = append(listen.Data, param.Data{
				Key:   key,
				Value: value,
			})
			return true
		})
	}
	buff, err := listen.EnCode(listen)
	if err != nil {
		log.Errorln(err.Error())
		return err
	}
	//发送消息
	if err := d._SendMsg(d.conn, param.Listen, buff); err != nil {
		log.Errorln(err.Error())
		return err
	}
	return nil
}

//_SendMsg 发送消息
func (d *Discovery) _SendMsg(conn net.Conn, cmd int8, buff []byte) error {
	event := new(param.Event)
	event.Cmd = cmd
	event.Data = buff
	data, err := event.EnCode(event)
	if err != nil {
		log.Errorln(err.Error())
		return err
	}
	if _, err := io.WriteByTime(conn, data, time.Now().Add(time.Second*5)); err != nil {
		log.Errorln(err.Error())
		return err
	}
	return nil
}

//EventLoop 事件处理器
func (d *Discovery) _EventLoop() {
	for {
		select {
		case <-d.close:
			close(d.close)
			return
		case <-time.After(time.Second * 5):
			//执行操作
			d.topics.Range(func(k, v interface{}) bool {
				if lable, ok := k.(param.Label); ok {
					if err := d.PushEvent(lable); err != nil {
						log.Errorln(err.Error())
						d.conn.Close()
						return false
					}
				}
				return true
			})

		}
	}
}
