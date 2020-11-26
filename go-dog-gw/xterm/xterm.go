package xterm

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
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
	method := c.Query("method")
	if method == "" {
		return
	}
	id := c.Query("id")
	if id == "" {
		return
	}
	uid := uuid.GetToken()
	log.Tracef("玩家%s调用方法%s获取%s信息", token, method, id)
	client := newclient(pointer.service.GetClient(), fmt.Sprintf("%s:%d", pointer.service.GetCfg().GetHost(), pointer.service.GetCfg().GetPort()), c.ClientIP(), token, uid, method, id, conn)
	pointer.clitens.Store(uid, client)
	client.run()
	pointer.clitens.Delete(uid)
	log.Tracef("玩家%s调用方法%s获取%s信息关闭", token, method, id)
}

//Push 消息推送
func (pointer *Ws) XtermPush(ctx plugins.Context, request param.XtermPushReq) (response param.XtermPushRes, err error) {
	log.Tracef("推送消息 | %s |", request.Uid)
	value, ok := pointer.clitens.Load(request.Uid)
	if !ok {
		err = errors.New("uid 不存在")
		return
	}
	if cli, o := value.(*client); o {
		buff, _ := json.Marshal(request.Msg)
		cli.push(buff)
	}
	return
}

type client struct {
	conn   *websocket.Conn
	client plugins.Client
	serip  string
	ip     string
	token  string
	method string
	id     string
	uid    string
}

func newclient(cli plugins.Client, serip, ip, token, uid, method, id string, conn *websocket.Conn) *client {
	return &client{
		serip:  serip,
		ip:     ip,
		client: cli,
		conn:   conn,
		token:  token,
		method: method,
		uid:    uid,
		id:     id,
	}
}

func (c *client) run() {
	ctx := context.Background()
	ctx.SetIsTest(false)
	ctx.SetTraceID(c.uid)
	ctx.SetToken(c.token)
	ctx.SetAddress(c.ip)
	//调用控制中心方法
	go func() {
		switch strings.ToLower(c.method) {
		case strings.ToLower("StartListDockerLog"):
			if _, err := rpc.StartListDockerLog(c.client, context.WithTimeout(ctx, int64(time.Second*time.Duration(6))), c.uid, c.serip, c.id); err != nil {
				log.Errorln(err.Error())
				c.conn.Close()
			}
		default:
			log.Warnln("方法不正确")
			c.conn.Close()
		}
	}()
	for {
		_, err := c.read(c.conn, time.Now().Add(time.Second*10))
		if err != nil {
			c.conn.Close()
			log.Errorln(err.Error())
			break
		}
	}
	switch strings.ToLower(c.method) {
	case strings.ToLower("StartListDockerLog"):
		if _, err := rpc.EndListDockerLog(c.client, context.WithTimeout(ctx, int64(time.Second*time.Duration(6))), c.uid); err != nil {
			log.Errorln(err.Error())
		}
	default:
		log.Errorln("正常退出")
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
