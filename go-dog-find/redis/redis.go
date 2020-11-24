package redis

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/tang-go/go-dog/cache"
	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/plugins"
)

const (
	RPCKey = "go-dog-find-rpc-key"
	APIKey = "go-dog-find-api-key"
)

//Item 内容
type Item struct {
	Key   string
	Value string
}

//Redis 基于redis做数据一致性
type Redis struct {
	rpclist sync.Map
	apilist sync.Map
	ttl     int64
	close   chan bool
	cache   *cache.Cache
}

//NewReids 新建一个redis
func NewReids(cfg plugins.Cfg) *Redis {
	r := new(Redis)
	r.ttl = 10
	r.cache = cache.NewCache(cfg)
	r.close = make(chan bool)
	go r.eventloop()
	return r
}

//Close 关闭
func (r *Redis) Close() {
	r.close <- true
}

//eventloop 事件处理器
func (r *Redis) eventloop() {
	for {
		select {
		case <-r.close:
			//收到关闭
			return
		case <-time.After(5 * time.Second):
			//更新内存数据
			now := time.Now().Unix()
			old := now - r.ttl
			r.updateRPC(now, old)
			r.updateAPI(now, old)
		}
	}
}

//updateRPC 更新rpc
func (r *Redis) updateRPC(max, min int64) {
	items, err := r.cache.GetCache().ZRevRangeByScore(RPCKey, max, min)
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	//更新最新数据
	rpcs := make(map[interface{}]interface{})
	for _, vali := range items {
		item := new(Item)
		if err := json.Unmarshal([]byte(vali), item); err != nil {
			log.Errorln(err.Error())
			continue
		}
		if _, ok := r.rpclist.Load(item.Key); !ok {
			r.rpclist.Store(item.Key, item)
			log.Traceln("rpc 服务 ", item.Key, "上线")
		}
		rpcs[item.Key] = item

	}
	r.rpclist.Range(func(key, vali interface{}) bool {
		if _, ok := rpcs[key]; !ok {
			r.rpclist.Delete(key)
			log.Traceln("rpc 服务 ", key, "下线")
		}
		return true
	})
	_, err = r.cache.GetCache().ZRemRangeByScore(RPCKey, 0, min)
	if err != nil {
		log.Errorln(err.Error())
		return
	}
}

//updateAPI 更新api
func (r *Redis) updateAPI(max, min int64) {
	items, err := r.cache.GetCache().ZRevRangeByScore(APIKey, max, min)
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	//更新最新数据
	apis := make(map[interface{}]interface{})
	for _, vali := range items {
		item := new(Item)
		if err := json.Unmarshal([]byte(vali), item); err != nil {
			log.Errorln(err.Error())
			continue
		}
		if _, ok := r.apilist.Load(item.Key); !ok {
			r.apilist.Store(item.Key, item)
			log.Traceln("api 服务 ", item.Key, "上线")
		}
		apis[item.Key] = item
	}
	r.apilist.Range(func(key, vali interface{}) bool {
		if _, ok := apis[key]; !ok {
			r.apilist.Delete(key)
			log.Traceln("api 服务 ", key, "下线")
		}
		return true
	})
	_, err = r.cache.GetCache().ZRemRangeByScore(RPCKey, 0, min)
	if err != nil {
		log.Errorln(err.Error())
		return
	}
}

//RangeRPC 便利rpc服务
func (r *Redis) RangeRPC(f func(key string, value string) bool) {
	r.rpclist.Range(func(k, v interface{}) bool {
		if item, ok := v.(*Item); ok {
			return f(item.Key, item.Value)
		}
		return true
	})
}

//RangeAPI 便利api服务
func (r *Redis) RangeAPI(f func(key string, value string) bool) {
	r.apilist.Range(func(k, v interface{}) bool {
		if item, ok := v.(*Item); ok {
			return f(item.Key, item.Value)
		}
		return true
	})
}

//RegisterRPC 注册RPC服务
func (r *Redis) RegisterRPC(key string, value string) {
	//储存 并且更新时间
	item := &Item{Key: key, Value: value}
	r.rpclist.Store(key, item)
	vali, err := json.Marshal(item)
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	back, err := r.cache.GetCache().Zadd(RPCKey, time.Now().Unix(), string(vali))
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	if back > 0 {
		log.Traceln("rpc 服务 ", key, "上线")
	}
}

//DelRPC 删除RPC服务
func (r *Redis) DelRPC(key string, value string) {
	item := &Item{Key: key, Value: value}
	r.rpclist.Delete(key)
	vali, err := json.Marshal(item)
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	_, err = r.cache.GetCache().ZRem(RPCKey, string(vali))
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	log.Traceln("rpc 服务 ", key, "下线")
}

//RegisterRPC 注册RPC服务
func (r *Redis) RegisterAPI(key string, value string) {
	//储存 并且更新时间
	item := &Item{Key: key, Value: value}
	r.apilist.Store(key, item)
	vali, err := json.Marshal(item)
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	back, err := r.cache.GetCache().Zadd(APIKey, time.Now().Unix(), string(vali))
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	if back > 0 {
		log.Traceln("api 服务 ", key, "上线")
	}
}

//DelRPC 删除RPC服务
func (r *Redis) DelAPI(key string, value string) {
	item := &Item{Key: key, Value: value}
	r.apilist.Delete(key)
	vali, err := json.Marshal(item)
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	_, err = r.cache.GetCache().ZRem(APIKey, string(vali))
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	log.Traceln("api 服务 ", key, "下线")
}
