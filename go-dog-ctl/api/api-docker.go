package api

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/tang-go/go-dog-tool/define"
	"github.com/tang-go/go-dog-tool/go-dog-ctl/param"
	"github.com/tang-go/go-dog-tool/go-dog-ctl/table"
	gateParam "github.com/tang-go/go-dog-tool/go-dog-gw/param"
	customerror "github.com/tang-go/go-dog/error"
	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/plugins"
)

//StartDocker 以docker模式启动docker
func (pointer *API) StartDocker(ctx plugins.Context, request param.StartDockerReq) (response param.StartDockerRes, err error) {
	admin, ok := ctx.GetShareByKey("Admin").(*table.Admin)
	if ok == false {
		err = customerror.EnCodeError(define.GetAdminInfoErr, "管理员信息失败")
		return
	}
	if pointer.mysql.GetReadEngine().Where("name = ? AND owner_id = ?", request.Name, admin.OwnerID).First(&table.Docker{}).RecordNotFound() == false {
		err = customerror.EnCodeError(define.StartDockerErr, "已经存在次名称服务")
		return
	}
	name := fmt.Sprintf("%d-%s", admin.OwnerID, request.Name)
	pointer._CloseDocker(name)
	go func() {
		if e := pointer._PullImage(request.Account, request.Pwd, request.Images, func(res string) {
			ctx.GetClient().Broadcast(ctx, define.SvcGateWay, "Push", &gateParam.PushReq{
				Token: ctx.GetToken(),
				Topic: define.RunDockerTopic,
				Msg:   res,
			}, &gateParam.PushRes{})
		}); e != nil {
			ctx.GetClient().Broadcast(ctx, define.SvcGateWay, "Push", &gateParam.PushReq{
				Token: ctx.GetToken(),
				Topic: define.RunDockerTopic,
				Msg:   e.Error(),
			}, &gateParam.PushRes{})
			return
		}
		config := &container.Config{
			Image:      request.Images,
			Domainname: name,
		}
		portSet := make(map[nat.Port]struct{})
		portBindings := make(map[nat.Port][]nat.PortBinding)
		for _, port := range request.Ports {
			portSet[nat.Port(port.InsidePort)] = struct{}{}
			portBindings[nat.Port(port.InsidePort)] = []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: port.ExternalPort,
				},
			}
		}
		config.ExposedPorts = portSet
		hostConfig := &container.HostConfig{
			PortBindings: portBindings,
		}
		containerResp, e := pointer.docker.ContainerCreate(ctx, config, hostConfig, nil, name)
		if e != nil {
			log.Errorln(e.Error())
			ctx.GetClient().Broadcast(ctx, define.SvcGateWay, "Push", &gateParam.PushReq{
				Token: ctx.GetToken(),
				Topic: define.RunDockerTopic,
				Msg:   e.Error(),
			}, &gateParam.PushRes{})
			return
		}
		if e := pointer.docker.ContainerStart(ctx, containerResp.ID, types.ContainerStartOptions{}); e != nil {
			log.Errorln(e.Error())
			ctx.GetClient().Broadcast(ctx, define.SvcGateWay, "Push", &gateParam.PushReq{
				Token: ctx.GetToken(),
				Topic: define.RunDockerTopic,
				Msg:   e.Error(),
			}, &gateParam.PushRes{})
			return
		}
		ctx.GetClient().Broadcast(ctx, define.SvcGateWay, "Push", &gateParam.PushReq{
			Token: ctx.GetToken(),
			Topic: define.RunDockerTopic,
			Msg:   "启动成功",
		}, &gateParam.PushRes{})

		//添加记录
		docker := &table.Docker{
			//唯一主键
			ID: pointer.snowflake.GetID(),
			//Name
			Name: request.Name,
			//编译发布的管理员
			AdminID: admin.AdminID,
			//发布镜像
			Image: request.Images,
			//账号
			Account: request.Account,
			//密码
			Pwd: request.Pwd,
			//业主ID
			OwnerID: admin.OwnerID,
			//注册事件
			Time: time.Now().Unix(),
		}
		if e := pointer.mysql.GetWriteEngine().Create(&docker).Error; e != nil {
			log.Errorln(e.Error())
		}
	}()
	//生成登录记录
	tbLog := &table.Log{
		//日志ID
		LogID: pointer.snowflake.GetID(),
		//类型
		Type: table.StartDockerType,
		//操作人
		AdminID: admin.AdminID,
		//名称
		AdminName: admin.Name,
		//操作方法
		Method: "StartDocker",
		//描述
		Description: "docker启动服务",
		//业主ID
		OwnerID: admin.OwnerID,
		//操作IP
		IP: ctx.GetAddress(),
		//操作URL
		URL: ctx.GetDataByKey("URL").(string),
		//操作时间
		Time: time.Now().Unix(),
	}
	if e := pointer.mysql.GetWriteEngine().Create(&tbLog).Error; e != nil {
		log.Errorln(e.Error())
		err = customerror.EnCodeError(define.StartDockerErr, "插入数据库记录失败")
		return
	}
	response.Success = true
	return
}

//CloseDocker 关闭docker容器
func (pointer *API) CloseDocker(ctx plugins.Context, request param.CloseDockerReq) (response param.CloseDockerRes, err error) {
	admin, ok := ctx.GetShareByKey("Admin").(*table.Admin)
	if ok == false {
		err = customerror.EnCodeError(define.GetAdminInfoErr, "管理员信息失败")
		return
	}
	if e := pointer._CloseDocker(request.ID); e != nil {
		log.Errorln(e.Error())
		err = customerror.EnCodeError(define.CloseDockerErr, e.Error())
		return
	}
	//生成登录记录
	tbLog := &table.Log{
		//日志ID
		LogID: pointer.snowflake.GetID(),
		//类型
		Type: table.CloseDockerType,
		//操作人
		AdminID: admin.AdminID,
		//名称
		AdminName: admin.Name,
		//操作方法
		Method: "CloseDocker",
		//描述
		Description: "关闭docker启动的服务",
		//业主ID
		OwnerID: admin.OwnerID,
		//操作IP
		IP: ctx.GetAddress(),
		//操作URL
		URL: ctx.GetDataByKey("URL").(string),
		//操作时间
		Time: time.Now().Unix(),
	}
	if e := pointer.mysql.GetWriteEngine().Create(&tbLog).Error; e != nil {
		log.Errorln(e.Error())
		err = customerror.EnCodeError(define.CloseDockerErr, "关闭服务失败")
		return
	}
	response.Success = true
	return
}

//GetDockerList  获取docker列表
func (pointer *API) GetDockerList(ctx plugins.Context, request param.GetDockerListReq) (response param.GetDockerListRes, err error) {
	admin, ok := ctx.GetShareByKey("Admin").(*table.Admin)
	if ok == false {
		err = customerror.EnCodeError(define.GetBuildServiceListErr, "管理员信息失败")
		return
	}
	var dockers []table.Docker
	if e := pointer.mysql.GetReadEngine().Where("owner_id = ?", admin.OwnerID).Order("time DESC").Find(&dockers).Error; e != nil {
		err = customerror.EnCodeError(define.GetBuildServiceListErr, e.Error())
		return
	}
	mp := make(map[string]types.Container)
	containers, err := pointer.docker.ContainerList(context.Background(), types.ContainerListOptions{})
	for _, container := range containers {
		names := container.Names[0]
		index := strings.Index(names, "-")
		if index <= 0 {
			continue
		}
		ownerID := names[1:index]
		if ownerID != strconv.FormatInt(admin.OwnerID, 10) {
			continue
		}
		name := names[index+1:]
		mp[name] = container
		// p := ""
		// for _, port := range container.Ports {
		// 	p = fmt.Sprintf("%s -p %d:%d/%s", p, port.PublicPort, port.PrivatePort, port.Type)
		// }
	}

	total := len(dockers)
	response.PageNo = 1
	response.PageSize = total
	if total%response.PageSize > 0 {
		response.TotalPage = total/response.PageSize + 1
	}
	if total%response.PageSize < 0 {
		response.TotalPage = total / response.PageSize
	}
	response.TotalCount = total

	for _, docker := range dockers {
		d := param.Docker{
			Name:    docker.Name,
			AdminID: admin.AdminID,
			Image:   docker.Image,
			OwnerID: admin.OwnerID,
			Time: time.Unix(docker.Time, 0).Format("2006-01-02 15:04:05"),
		}
		container, ok := mp[docker.Name]
		if ok {

		}
		d.ID = container.ID
		d.Status = container.Status
		d.Command = container.Command
		p := ""
		for _, port := range container.Ports {
			p = fmt.Sprintf("%s -p %d:%d/%s", p, port.PublicPort, port.PrivatePort, port.Type)
		}
		d.Ports = p
		response.Data = append(response.Data, d)
	}
	return
}
