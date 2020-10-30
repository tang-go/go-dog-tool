package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/tang-go/go-dog-tool/define"
	"github.com/tang-go/go-dog-tool/go-dog-ctl/param"
	"github.com/tang-go/go-dog-tool/go-dog-ctl/table"
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
	//加密登录信息方式
	auth, e := pointer.docker.RegistryLogin(context.Background(), types.AuthConfig{
		Username:      request.Account,
		Password:      request.Pwd,
		ServerAddress: "qa.game.com",
	})
	if e != nil {
		log.Errorln(e.Error())
		err = customerror.EnCodeError(define.StartDockerErr, "登录失败")
		return
	}
	events, e := pointer.docker.ImagePull(context.Background(), request.Images, types.ImagePullOptions{
		RegistryAuth: auth.IdentityToken,
	})
	if e != nil {
		log.Errorln(e.Error())
		err = customerror.EnCodeError(define.StartDockerErr, fmt.Sprintf("获取镜像失败:%s", request.Images))
		return
	}

	d := json.NewDecoder(events)
	event := make(map[string]interface{})
	for {
		if e := d.Decode(&event); e != nil {
			if err == io.EOF {
				break
			}
			log.Errorln(e.Error())
			break
		}
		fmt.Println(event)
	}
	name := fmt.Sprintf("%d-%s", admin.OwnerID, request.Name)
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
		err = customerror.EnCodeError(define.StartDockerErr, "创建docker服务失败，请查看端口和服务名称是否已经存在")
		return
	}
	if e := pointer.docker.ContainerStart(ctx, containerResp.ID, types.ContainerStartOptions{}); e != nil {
		log.Errorln(e.Error())
		err = customerror.EnCodeError(define.StartDockerErr, "创建docker服务失败，请查看端口和服务名称是否已经存在")
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

	containers, err := pointer.docker.ContainerList(context.Background(), types.ContainerListOptions{})
	total := len(containers)
	response.PageNo = 1
	response.PageSize = total
	if total%response.PageSize > 0 {
		response.TotalPage = total/response.PageSize + 1
	}
	if total%response.PageSize < 0 {
		response.TotalPage = total / response.PageSize
	}
	response.TotalCount = total
	for _, container := range containers {
		names := container.Names[0]
		index := strings.Index(names, "-")
		if index <= 0 {
			return
		}
		ownerID := names[1:index]
		if ownerID != strconv.FormatInt(admin.OwnerID, 10) {
			return
		}
		name := names[index+1:]
		p := ""
		for _, port := range container.Ports {
			p = fmt.Sprintf("%s -p %d:%d/%s", p, port.PublicPort, port.PrivatePort, port.Type)
		}

		response.Data = append(response.Data, param.Docker{
			ID:      container.ID,
			Name:    name,
			AdminID: admin.AdminID,
			Image:   container.Image,
			Status:  container.Status,
			OwnerID: admin.OwnerID,
			Command: container.Command,
			Ports:   p,
			Time:    time.Unix(container.Created, 0).Format("2006-01-02 15:04:05"),
		})
	}
	return
}
