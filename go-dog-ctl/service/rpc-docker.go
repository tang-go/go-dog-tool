package service

import (
	"bufio"
	"net"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/tang-go/go-dog-tool/define"
	"github.com/tang-go/go-dog-tool/go-dog-ctl/param"
	"github.com/tang-go/go-dog-tool/go-dog-ctl/table"
	"github.com/tang-go/go-dog-tool/go-dog-gw/rpc"
	customerror "github.com/tang-go/go-dog/error"
	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/pkg/context"
	"github.com/tang-go/go-dog/plugins"
)

//StartListDockerLog 开始监听docker日志
func (s *Service) StartListDockerLog(ctx plugins.Context, request param.StartListDockerLogReq) (response param.StartListDockerLogRes, err error) {
	admin := new(table.Admin)
	if e := s.cache.GetCache().Get(ctx.GetToken(), admin); e != nil {
		log.Errorln(e.Error())
		err = customerror.EnCodeError(define.AdminTokenErr, "token失效或者不正确")
		return
	}
	if _, ok := s.dockerListn.Load(request.UID); ok {
		err = customerror.EnCodeError(define.StartListDockerLogErr, "Uid已经存在了")
		return
	}
	reader, e := s.docker.ContainerLogs(ctx, request.ID, types.ContainerLogsOptions{ShowStderr: true, ShowStdout: true, Follow: true, Tail: "100"})
	if e != nil {
		log.Errorln("获取容器ID ", request.ID, " 错误 ", e.Error())
		err = customerror.EnCodeError(define.StartListDockerLogErr, e.Error())
		return
	}
	logs, e := s.docker.ContainerAttach(ctx, request.ID, types.ContainerAttachOptions{Stderr: true, Stdout: true, Stream: true, Logs: false})
	if e != nil {
		log.Errorln("获取容器ID ", request.ID, " 错误 ", e.Error())
		err = customerror.EnCodeError(define.StartListDockerLogErr, e.Error())
		return
	}
	go func() {
		read := bufio.NewScanner(reader)
		for read.Scan() {
			if _, e := rpc.XtermPush(s.service.GetClient(), context.WithTimeout(ctx, int64(time.Second*time.Duration(6))), request.Address, request.UID, read.Text()); e != nil {
				log.Warnln("推送错误，退出", e.Error())
				reader.Close()
				return
			}
		}
		log.Traceln("获取实时日志")
		reader.Close()
		s.dockerListn.Store(request.UID, logs.Conn)
		scanner := bufio.NewScanner(logs.Conn)
		for scanner.Scan() {
			if _, e := rpc.XtermPush(s.service.GetClient(), context.WithTimeout(ctx, int64(time.Second*time.Duration(6))), request.Address, request.UID, scanner.Text()); e != nil {
				logs.Conn.Close()
				log.Warnln("推送错误，退出", e.Error())
				break
			}
		}
		s.dockerListn.Delete(request.UID)
		log.Warnln("链接关闭，正常退出", request.UID)
	}()
	response.Success = true
	return
}

//EndListDockerLog 结束监听docker日志
func (s *Service) EndListDockerLog(ctx plugins.Context, request param.EndListDockerLogReq) (response param.EndListDockerLogRes, err error) {
	value, ok := s.dockerListn.Load(request.UID)
	if !ok {
		//此处表示已经先退出的链接
		response.Success = true
		return
	}
	if conn, ok := value.(net.Conn); ok {
		conn.Close()
		log.Warnln("链接关闭，正常退出", request.UID)
	}
	response.Success = true
	return
}
