package gateway

import (
	"encoding/json"
	"net"
	"sync"
	"time"

	"github.com/tang-go/go-dog/lib/io"
	"github.com/tang-go/go-dog/lib/rand"
	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/pkg/discovery/param"
	"github.com/tang-go/go-dog/serviceinfo"
)

//ServcieAPI api列表
type ServcieAPI struct {
	Method *serviceinfo.API
	Name   string
	Count  int32
}

//GoDogDiscovery 服务发现
type GoDogDiscovery struct {
	address    []string
	conn       net.Conn
	ttl        time.Duration
	pos        int
	count      int
	close      bool
	closeheart chan bool
	apidata    map[string]*serviceinfo.APIServiceInfo
	rpcdata    map[string]*serviceinfo.RPCServiceInfo
	apis       map[string]*ServcieAPI
	lock       sync.RWMutex
}

//NewGoDogDiscovery  新建发现服务
func NewGoDogDiscovery(address []string) *GoDogDiscovery {
	dis := &GoDogDiscovery{
		address:    address,
		ttl:        2 * time.Second,
		count:      len(address),
		pos:        0,
		close:      false,
		closeheart: make(chan bool),
		apidata:    make(map[string]*serviceinfo.APIServiceInfo),
		rpcdata:    make(map[string]*serviceinfo.RPCServiceInfo),
		apis:       make(map[string]*ServcieAPI),
	}
	if err := dis._ConnectClient(); err != nil {
		panic(err)
	}
	//等待一个心跳时间
	//time.Sleep(dis.ttl)
	return dis
}

//GetAllAPIService 获取所有API服务
func (d *GoDogDiscovery) GetAllAPIService() (services []*serviceinfo.APIServiceInfo) {
	d.lock.RLock()
	for _, service := range d.apidata {
		services = append(services, service)
	}
	d.lock.RUnlock()
	return
}

//GetAllRPCService 获取所有RPC服务
func (d *GoDogDiscovery) GetAllRPCService() (services []*serviceinfo.RPCServiceInfo) {
	d.lock.RLock()
	for _, service := range d.rpcdata {
		services = append(services, service)
	}
	d.lock.RUnlock()
	return
}

//GetRPCServiceByName 通过名称获取RPC服务
func (d *GoDogDiscovery) GetRPCServiceByName(name string) (services []*serviceinfo.RPCServiceInfo) {
	d.lock.RLock()
	for _, service := range d.rpcdata {
		if service.Name == name {
			services = append(services, service)
		}
	}
	d.lock.RUnlock()
	return
}

//GetAPIByURL 通过RUL获取API服务
func (d *GoDogDiscovery) GetAPIByURL(url string) (*ServcieAPI, bool) {
	d.lock.RLock()
	defer d.lock.RUnlock()
	s, ok := d.apis[url]
	return s, ok
}

//RangeAPI 遍历api
func (d *GoDogDiscovery) RangeAPI(f func(url string, api *ServcieAPI)) {
	d.lock.RLock()
	for url, api := range d.apis {
		f(url, api)
	}
	d.lock.RUnlock()
}

//GetAPIServiceByName 通过名称获取API服务
func (d *GoDogDiscovery) GetAPIServiceByName(name string) (services []*serviceinfo.APIServiceInfo) {
	d.lock.RLock()
	for _, service := range d.apidata {
		if service.Name == name {
			services = append(services, service)
		}
	}
	d.lock.RUnlock()
	return
}

//_ConnectClient 建立链接
func (d *GoDogDiscovery) _ConnectClient() error {
	d.lock.Lock()
	defer d.lock.Unlock()
	if d.close {
		return nil
	}
	//随机链接
	index := rand.IntRand(0, d.count)
	address := d.address[index]
	tcpAddr, err := net.ResolveTCPAddr("tcp4", address)
	if err != nil {
		return err
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return err
	}
	//发送登陆请求
	login := new(param.LoginReq)
	login.Type = param.DisType
	buff, err := login.EnCode(login)
	if err != nil {
		conn.Close()
		log.Errorln(err.Error())
		return err
	}
	if _, err := io.WriteByTime(conn, buff, time.Now().Add(d.ttl)); err != nil {
		//断线开启重新链接
		conn.Close()
		log.Errorln(err.Error())
		return err
	}
	d.conn = conn
	//开启心跳
	go d._Heart()
	//开启监听
	go d._Watch()
	//默认监听rpc服务消息
	d._WatchRPCService()
	d._WatchAPIService()
	log.Traceln("链接成功注册中心", address)
	return nil
}

//_Watch 开始监听
func (d *GoDogDiscovery) _Watch() {
	for {
		_, buff, err := io.Read(d.conn)
		if err != nil {
			d.closeheart <- true
			d.conn.Close()
			log.Errorln(err.Error())
			break
		}
		all := new(param.All)
		if err := all.DeCode(buff, all); err != nil {
			log.Errorln(err.Error())
			continue
		}
		d.lock.Lock()
		if all.Label == "/rpc" {
			mp := make(map[string]string)
			for _, data := range all.Datas {
				if _, ok := d.rpcdata[data.Key]; !ok {
					info := new(serviceinfo.RPCServiceInfo)
					if err := json.Unmarshal([]byte(data.Value), info); err != nil {
						log.Errorln(err.Error(), data.Key, data.Value)
						continue
					}
					d.rpcdata[data.Key] = info
					log.Tracef("rpc 上线 | %s | %s | %s ", info.Name, data.Key, info.Address)
				}
				mp[data.Key] = data.Value
			}
			for key, info := range d.rpcdata {
				if _, ok := mp[key]; !ok {
					delete(d.rpcdata, key)
					log.Tracef("rpc 下线 | %s | %s | %s ", info.Name, key, info.Address)
				}
			}
		}
		if all.Label == "/api" {
			mp := make(map[string]string)
			for _, data := range all.Datas {
				if _, ok := d.apidata[data.Key]; !ok {
					info := new(serviceinfo.APIServiceInfo)
					if err := json.Unmarshal([]byte(data.Value), info); err != nil {
						log.Errorln(err.Error(), data.Key, data.Value)
						continue
					}
					for _, method := range info.API {
						url := "/api/" + info.Name + "/" + method.Version + "/" + method.Path
						if api, ok := d.apis[url]; ok {
							api.Count++
						} else {
							d.apis[url] = &ServcieAPI{
								Method: method,
								Name:   info.Name,
								Count:  1,
							}
							log.Tracef(" 上线 | %s | %s | %s ", info.Name, data.Key, url)
						}
					}
					d.apidata[data.Key] = info
				}
				mp[data.Key] = data.Value
			}
			for key, data := range d.apidata {
				if _, ok := mp[key]; !ok {
					for _, method := range data.API {
						url := "/api/" + data.Name + "/" + method.Version + "/" + method.Path
						if api, ok := d.apis[url]; ok {
							api.Count--
							if api.Count <= 0 {
								delete(d.apis, url)
								log.Tracef(" 下线 | %s | %s | %s ", data.Name, data.Key, url)
							}
						}
					}
					delete(d.apidata, key)
				}
			}
		}
		d.lock.Unlock()
	}

	for {
		time.Sleep(d.ttl)
		log.Traceln("断线重链注册中心....")
		if d._ConnectClient() == nil {
			return
		}
	}
}

//_Heart 心跳
func (d *GoDogDiscovery) _Heart() {
	heart := &param.Event{
		Cmd: param.Heart,
	}
	buff, _ := heart.EnCode(heart)
	for {
		select {
		case <-d.closeheart:
			return
		case <-time.After(d.ttl):
			if _, err := io.WriteByTime(d.conn, buff, time.Now().Add(d.ttl)); err != nil {
				//断线开启重新链接
				d.conn.Close()
				log.Errorln(err.Error())
				break
			}
		}
	}
}

//WatchRPCService 开始RPC服务发现
func (d *GoDogDiscovery) _WatchRPCService() {
	//开启监听
	listen := &param.Event{
		Cmd:   param.Listen,
		Label: "/rpc",
	}
	buff, err := listen.EnCode(listen)
	if err != nil {
		panic(err.Error())
	}
	if _, err := io.WriteByTime(d.conn, buff, time.Now().Add(d.ttl)); err != nil {
		panic(err.Error())
	}
	log.Traceln("watch /rpc")
}

//WatchAPIService 开始API服务发现
func (d *GoDogDiscovery) _WatchAPIService() {
	//开启监听
	listen := &param.Event{
		Cmd:   param.Listen,
		Label: "/api",
	}
	buff, err := listen.EnCode(listen)
	if err != nil {
		panic(err.Error())
	}
	if _, err := io.WriteByTime(d.conn, buff, time.Now().Add(d.ttl)); err != nil {
		panic(err.Error())
	}
	log.Traceln("watch /api")
}

//Close 关闭服务
func (d *GoDogDiscovery) Close() error {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.close = true
	if d.conn != nil {
		return d.conn.Close()
	}
	return nil
}
