package rpc

import (
	"github.com/tang-go/go-dog-tool/define"
	"github.com/tang-go/go-dog-tool/go-dog-ctl/param"
	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/plugins"
)

//StartListDockerLog 开始监听docker日志
func StartListDockerLog(client plugins.Client, ctx plugins.Context, uid, address, id string) (bool, error) {
	res := new(param.StartListDockerLogRes)
	if err := client.Call(ctx, plugins.RandomMode, define.SvcController, "StartListDockerLog", &param.StartListDockerLogReq{Uid: uid, ID: id, Address: address}, res); err != nil {
		log.Errorln(err.Error())
		return res.Success, err
	}
	return res.Success, nil
}

//EndListDockerLog 结束监听docker日志（此处为广播所有服务）
func EndListDockerLog(client plugins.Client, ctx plugins.Context, uid string) (bool, error) {
	res := new(param.EndListDockerLogRes)
	if err := client.Broadcast(ctx, define.SvcController, "EndListDockerLog", &param.EndListDockerLogReq{Uid: uid}, res); err != nil {
		log.Errorln(err.Error())
		return res.Success, err
	}
	return res.Success, nil
}
