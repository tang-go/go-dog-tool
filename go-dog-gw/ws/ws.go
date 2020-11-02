package ws

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/tang-go/go-dog-tool/go-dog-gw/param"
	"github.com/tang-go/go-dog/cache"
	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/plugins"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	HandshakeTimeout: 5 * time.Second,
	// 取消ws跨域校验
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

//Ws websocket客户端
type Ws struct {
	cache   *cache.Cache
	clitens sync.Map
}

//NewWs 新建ws
func NewWs(cfg plugins.Cfg) *Ws {
	return &Ws{
		//初始化缓存
		cache: cache.NewCache(cfg),
	}
}

// Connect websocket链接
func (pointer *Ws) Connect(w http.ResponseWriter, r *http.Request, c *gin.Context) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Warnln("ws链接失败 ", err.Error())
		return
	}
	defer conn.Close()
	token := c.Query("token")
	if token == "" {
		return
	}
	if _, ok := pointer.clitens.Load(token); ok {
		log.Tracef("重复登录 | %s |", token)
		return
	}
	info := new(interface{})
	if e := pointer.cache.GetCache().Get(token, info); e != nil {
		log.Errorln(e.Error())
		return
	}
	log.Tracef("玩家上线 | %s |", token)
	client := newclient(token, conn)
	pointer.clitens.Store(token, client)
	client.run()
	pointer.clitens.Delete(token)
}

//Push 消息推送
func (pointer *Ws) Push(ctx plugins.Context, request param.PushReq) (response param.PushRes, err error) {
	log.Tracef("推送消息 | %s | %s | %s |", request.Token, request.Topic, request.Msg)
	value, ok := pointer.clitens.Load(request.Token)
	if !ok {
		return
	}
	if msg, err := json.Marshal(request); err == nil {
		if cli, o := value.(*client); o {
			cli.push(msg)
		}
	}
	return
}

type client struct {
	conn  *websocket.Conn
	token string
}

func newclient(token string, conn *websocket.Conn) *client {
	return &client{
		conn:  conn,
		token: token,
	}
}

func (c *client) run() {
	for {
		_, err := c.read(c.conn, time.Now().Add(time.Second*5))
		if err != nil {
			c.conn.Close()
			log.Errorln(err.Error())
			return
		}
	}
}

func (c *client) read(conn *websocket.Conn, t time.Time) ([]byte, error) {
	err := conn.SetReadDeadline(t)
	if err != nil {
		return nil, err
	}
	defer conn.SetReadDeadline(time.Time{})
	_, message, err := conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (c *client) push(message []byte) {
	c.write(c.conn, message, time.Now().Add(time.Second*5))
}

func (c *client) write(conn *websocket.Conn, message []byte, t time.Time) error {
	err := conn.SetWriteDeadline(t)
	if err != nil {
		return err
	}
	defer conn.SetWriteDeadline(time.Time{})
	err = conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		return err
	}
	return nil
}
