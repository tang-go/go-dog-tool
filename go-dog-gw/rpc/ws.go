package rpc

import (
	"github.com/tang-go/go-dog-tool/define"
	"github.com/tang-go/go-dog-tool/go-dog-gw/param"
	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/plugins"
)

//Push 推送消息
func Push(client plugins.Client, ctx plugins.Context, address string, token, topic, msg string) (bool, error) {
	res := new(param.PushRes)
	if err := client.CallByAddress(ctx, address, define.SvcGateWay, "Push", &param.PushReq{Token: token, Topic: topic, Msg: msg}, res); err != nil {
		log.Errorln(err.Error())
		return res.Success, err
	}
	return res.Success, nil
}
