package service

import (
	"net"
	"sync"
	"time"

	"github.com/tang-go/go-dog-tool/go-dog-find/param"
	"github.com/tang-go/go-dog/lib/io"
	"github.com/tang-go/go-dog/log"
)

//Register 注册
type Register struct {
	conn        net.Conn
	offlinefunc func()
	datas       sync.Map
	service     *Service
}

//NewRegister 新建一个服务注册
func NewRegister(service *Service, conn net.Conn) *Register {
	return &Register{
		conn:    conn,
		service: service,
	}
}

//Run 启动
func (r *Register) Run() {
	for {
		_, buff, err := io.ReadByTime(r.conn, time.Now().Add(time.Second*5))
		if err != nil {
			r.datas.Range(func(key, value interface{}) bool {
				label, ok := key.(param.Label)
				if !ok {
					return true
				}
				data, ok := value.(*param.Data)
				if !ok {
					return true
				}
				if label == param.RPCLabel {
					r.service.redis.DelRPC(data.Key, data.Value)
				}
				if label == param.APILabel {
					r.service.redis.DelAPI(data.Key, data.Value)
				}
				return true
			})
			r.conn.Close()
			log.Errorln(err.Error())
			return
		}
		event := new(param.Event)
		if err := event.DeCode(buff, event); err != nil {
			log.Errorln(err.Error())
			continue
		}
		switch event.Cmd {
		case param.Reg:
			reg := new(param.RegReq)
			if err := reg.DeCode(event.Data, reg); err != nil {
				log.Errorln(err.Error())
			} else {
				r.datas.Store(reg.Label, &reg.Data)
			}
		}
		r.datas.Range(func(key, value interface{}) bool {
			label, ok := key.(param.Label)
			if !ok {
				return true
			}
			data, ok := value.(*param.Data)
			if !ok {
				return true
			}
			if label == param.RPCLabel {
				r.service.redis.RegisterRPC(data.Key, data.Value)
			}
			if label == param.APILabel {
				r.service.redis.RegisterAPI(data.Key, data.Value)
			}
			return true
		})
	}
}
