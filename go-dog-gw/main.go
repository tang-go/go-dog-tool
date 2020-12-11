package main

import (
	"github.com/gin-gonic/gin"
	"github.com/tang-go/go-dog-tool/define"
	"github.com/tang-go/go-dog-tool/go-dog-ctl/rpc"
	"github.com/tang-go/go-dog-tool/go-dog-gw/gateway"
	"github.com/tang-go/go-dog-tool/go-dog-gw/ws"
	"github.com/tang-go/go-dog-tool/go-dog-gw/xterm"
	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/plugins"
)

func main() {
	gate := gateway.NewGateway(define.SvcGateWay, define.SvcController)
	//网关验证权限
	gate.Auth(func(client plugins.Client, ctx plugins.Context, token, url string) error {
		admin, err := rpc.Auth(client, ctx, token, url)
		if err != nil {
			log.Errorln(err.Error())
			return err
		}
		return ctx.SetData("Admin", admin)
	})
	//初始化websocket客户端
	ws := ws.NewWs(gate.GetService())
	gate.OpenWebSocket("/ws", func(c *gin.Context) {
		ws.Connect(c.Writer, c.Request, c)
	})
	//初始化xterm客户端
	xtermWs := xterm.NewWs(gate.GetService())
	gate.OpenWebSocket("/xtermws", func(c *gin.Context) {
		xtermWs.Connect(c.Writer, c.Request, c)
	})
	gate.Run(8080)
}
