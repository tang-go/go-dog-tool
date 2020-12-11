package ws

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/tang-go/go-dog-tool/go-dog-ctl/rpc"
	"github.com/tang-go/go-dog-tool/go-dog-gw/param"
	"github.com/tang-go/go-dog/lib/uuid"
	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/pkg/context"
	"github.com/tang-go/go-dog/plugins"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	HandshakeTimeout: 10 * time.Second,
	// 取消ws跨域校验
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

//Ws websocket客户端
type Ws struct {
	service plugins.Service
	clitens sync.Map
}

//NewWs 新建ws
func NewWs(service plugins.Service) *Ws {
	ws := &Ws{
		service: service,
	}
	service.RPC("Push", 3, false, "推送消息", ws.Push)
	return ws
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
	//调用上线
	if e := pointer.online(c.ClientIP(), token); e != nil {
		log.Errorln(e.Error())
		return
	}
	log.Tracef("玩家上线 | %s | %s |", token, c.ClientIP())
	client := newclient(pointer, c.ClientIP(), token, conn)
	pointer.clitens.Store(token, client)
	client.run()
	pointer.clitens.Delete(token)
	if e := pointer.offline(c.ClientIP(), token); e != nil {
		log.Errorln(e.Error())
	}
	log.Tracef("玩家下线 | %s | %s |", token, c.ClientIP())
}

//Push 消息推送
func (pointer *Ws) Push(ctx plugins.Context, request param.PushReq) (response param.PushRes, err error) {
	log.Tracef("推送消息 | %s | %s | %s |", request.Token, request.Topic, request.Msg)
	value, ok := pointer.clitens.Load(request.Token)
	if !ok {
		err = errors.New("token 不存在")
		return
	}
	msg, e := json.Marshal(request)
	if e != nil {
		err = e
		return
	}
	if cli, o := value.(*client); o {
		cli.push(msg)
	}
	return
}

//_Online 上线
func (pointer *Ws) online(address, token string) error {
	ctx := context.Background()
	ctx.SetIsTest(false)
	ctx.SetTraceID(uuid.GetToken())
	ctx.SetToken(token)
	ctx.SetAddress(address)
	_, err := rpc.AdminOnline(pointer.service.GetClient(), context.WithTimeout(ctx, int64(time.Second*time.Duration(6))), fmt.Sprintf("%s:%d", pointer.service.GetCfg().GetHost(), pointer.service.GetCfg().GetPort()))
	return err
}

//_Offline 下线
func (pointer *Ws) offline(address, token string) error {
	ctx := context.Background()
	ctx.SetIsTest(false)
	ctx.SetTraceID(uuid.GetToken())
	ctx.SetToken(token)
	ctx.SetAddress(address)
	_, err := rpc.AdminOffline(pointer.service.GetClient(), context.WithTimeout(ctx, int64(time.Second*time.Duration(6))), fmt.Sprintf("%s:%d", pointer.service.GetCfg().GetHost(), pointer.service.GetCfg().GetPort()))
	return err
}

type client struct {
	ws      *Ws
	conn    *websocket.Conn
	address string
	token   string
}

func newclient(ws *Ws, address string, token string, conn *websocket.Conn) *client {
	return &client{
		ws:      ws,
		address: address,
		conn:    conn,
		token:   token,
	}
}

func (c *client) run() {
	count := 0
	for {
		_, err := c.read(c.conn, time.Now().Add(time.Second*10))
		if err != nil {
			c.conn.Close()
			log.Errorln(err.Error())
			return
		}
		if count%4 == 0 {
			count = 0
			if err := c.ws.online(c.address, c.token); err != nil {
				c.conn.Close()
				log.Errorln(err.Error())
				return
			}
		}
		count++
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
	c.write(c.conn, message, time.Now().Add(time.Second*10))
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
