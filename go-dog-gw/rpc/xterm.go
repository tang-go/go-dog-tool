package rpc

import (
	"github.com/tang-go/go-dog-tool/define"
	"github.com/tang-go/go-dog-tool/go-dog-gw/param"
	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/plugins"
)

//XtermPust 推送xterm消息
func XtermPush(client plugins.Client, ctx plugins.Context, address string, uid string, msg string) (bool, error) {
	res := new(param.XtermPushRes)
	if err := client.CallByAddress(ctx, address, define.SvcGateWay, "XtermPush", &param.XtermPushReq{Uid: uid, Msg: msg}, res); err != nil {
		log.Errorln(err.Error())
		return res.Success, err
	}
	return res.Success, nil
}
