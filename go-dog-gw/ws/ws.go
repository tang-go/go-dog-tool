package ws

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/tang-go/go-dog-tool/define"
	ctlParam "github.com/tang-go/go-dog-tool/go-dog-ctl/param"
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
	return &Ws{
		service: service,
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
	//调用上线
	if e := pointer.online(token); e != nil {
		log.Errorln(e.Error())
		return
	}
	log.Tracef("玩家上线 | %s |", token)
	client := newclient(pointer, token, conn)
	pointer.clitens.Store(token, client)
	client.run()
	pointer.clitens.Delete(token)
	if e := pointer.offline(token); e != nil {
		log.Errorln(e.Error())
	}
	log.Tracef("玩家下线 | %s |", token)
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

//_Online 上线
func (pointer *Ws) online(token string) error {
	ctx := context.Background()
	ctx.SetIsTest(false)
	ctx.SetTraceID(uuid.GetToken())
	ctx.SetToken(token)
	return pointer.service.GetClient().Call(
		context.WithTimeout(ctx, int64(time.Second*time.Duration(6))),
		plugins.RandomMode,
		define.SvcController,
		"AdminOnline",
		&ctlParam.AdminOnlineReq{
			Address: fmt.Sprintf("%s:%d", pointer.service.GetCfg().GetHost(), pointer.service.GetCfg().GetPort()),
		}, &ctlParam.AdminOnlineRes{})
}

//_Offline 下线
func (pointer *Ws) offline(token string) error {
	ctx := context.Background()
	ctx.SetIsTest(false)
	ctx.SetTraceID(uuid.GetToken())
	ctx.SetToken(token)
	//调用下线
	return pointer.service.GetClient().Call(
		context.WithTimeout(ctx, int64(time.Second*time.Duration(6))),
		plugins.RandomMode,
		define.SvcController,
		"AdminOffline",
		&ctlParam.AdminOfflineReq{
			Address: fmt.Sprintf("%s:%d", pointer.service.GetCfg().GetHost(), pointer.service.GetCfg().GetPort()),
		}, &ctlParam.AdminOfflineRes{})
}

type client struct {
	ws    *Ws
	conn  *websocket.Conn
	token string
}

func newclient(ws *Ws, token string, conn *websocket.Conn) *client {
	return &client{
		ws:    ws,
		conn:  conn,
		token: token,
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
			if err := c.ws.online(c.token); err != nil {
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
