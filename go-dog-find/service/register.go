package service

import (
	"net"
	"time"

	"github.com/tang-go/go-dog-tool/go-dog-find/param"
	"github.com/tang-go/go-dog/lib/io"
	"github.com/tang-go/go-dog/log"
)

//Register 注册
type Register struct {
	conn        net.Conn
	offlinefunc func()
	datas       []*param.Data
	service     *Service
}

//NewRegister 新建一个服务注册
func NewRegister(service *Service, conn net.Conn, offlinefunc func()) *Register {
	return &Register{
		conn:        conn,
		offlinefunc: offlinefunc,
		service:     service,
	}
}

//Run 启动
func (r *Register) Run() {
	defer r.offlinefunc()
	for {
		_, buff, err := io.ReadByTime(r.conn, time.Now().Add(time.Second*5))
		if err != nil {
			for _, data := range r.datas {
				r.service.Del(data)
			}
			r.conn.Close()
			log.Errorln(err.Error())
			return
		}
		event := new(param.Event)
		if err := event.DeCode(buff, event); err != nil {
			log.Errorln(err.Error())
			continue
		}
		if event.Cmd == param.Reg {
			//注册事件
			r.datas = append(r.datas, event.Data)
		}
		for _, data := range r.datas {
			//上线服务
			data.Time = time.Now().Unix()
			r.service.Add(data)
		}
	}
}
