package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/tang-go/go-dog-tool/define"
	"github.com/tang-go/go-dog-tool/go-dog-ctl/param"
	"github.com/tang-go/go-dog-tool/go-dog-ctl/table"
	customerror "github.com/tang-go/go-dog/error"
	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/plugins"
)

//DelDocker 删除镜像
func (s *Service) DelDocker(ctx plugins.Context, request param.DelDockerReq) (response param.DelDockerRes, err error) {
	admin, e := s.GetAdmin(ctx)
	if e != nil {
		log.Errorln(e.Error())
		err = customerror.EnCodeError(define.GetAdminInfoErr, e.Error())
		return
	}
	docker := new(table.Docker)
	if s.mysql.GetReadEngine().Where("id = ?", request.DockerID).First(docker).RecordNotFound() == true {
		err = customerror.EnCodeError(define.DelDockerErr, "DockerID不正确")
		return
	}
	if admin.OwnerID != docker.OwnerID {
		err = customerror.EnCodeError(define.DelDockerErr, "DockerID不正确")
		return
	}
	//删除
	name := fmt.Sprintf("%d-%s", admin.OwnerID, docker.Name)
	s._CloseDocker(name)
	if e := s.mysql.GetWriteEngine().Delete(docker).Error; e != nil {
		log.Errorln(e.Error())
		err = customerror.EnCodeError(define.DelDockerErr, "删除失败")
		return
	}
	response.Success = true
	return
}

//RestartDocker 重启镜像
func (s *Service) RestartDocker(ctx plugins.Context, request param.RestartDockerReq) (response param.RestartDockerRes, err error) {
	admin, e := s.GetAdmin(ctx)
	if e != nil {
		log.Errorln(e.Error())
		err = customerror.EnCodeError(define.GetAdminInfoErr, e.Error())
		return
	}
	docker := new(table.Docker)
	if s.mysql.GetReadEngine().Where("id = ?", request.DockerID).First(docker).RecordNotFound() == true {
		err = customerror.EnCodeError(define.RestartDockerErr, "DockerID不正确")
		return
	}
	if admin.OwnerID != docker.OwnerID {
		err = customerror.EnCodeError(define.RestartDockerErr, "DockerID不正确")
		return
	}
	var ports []*param.Ports
	if e := json.Unmarshal([]byte(docker.Ports), &ports); e != nil {
		log.Errorln(e)
		err = customerror.EnCodeError(define.RestartDockerErr, "ports不正确")
		return
	}
	//重启
	name := fmt.Sprintf("%d-%s", admin.OwnerID, docker.Name)
	s._CloseDocker(name)
	go func() {
		e := s._StartDocker(ctx.GetToken(), docker.Image, name, docker.Account, docker.Pwd, ports)
		if e != nil {
			log.Errorln(e)
			return
		}
	}()
	//生成登录记录
	tbLog := &table.Log{
		//日志ID
		LogID: s.snowflake.GetID(),
		//类型
		Type: table.RestartDockerType,
		//操作人
		AdminID: admin.AdminID,
		//名称
		AdminName: admin.Name,
		//操作方法
		Method: "RestartDocker",
		//描述
		Description: "重启docker服务",
		//业主ID
		OwnerID: admin.OwnerID,
		//操作IP
		IP: ctx.GetAddress(),
		//操作URL
		URL: ctx.GetURL(),
		//操作时间
		Time: time.Now().Unix(),
	}
	if e := s.mysql.GetWriteEngine().Create(&tbLog).Error; e != nil {
		log.Errorln(e.Error())
	}
	response.Success = true
	return
}

//StartDocker 以docker模式启动docker
func (s *Service) StartDocker(ctx plugins.Context, request param.StartDockerReq) (response param.StartDockerRes, err error) {
	admin, e := s.GetAdmin(ctx)
	if e != nil {
		log.Errorln(e.Error())
		err = customerror.EnCodeError(define.GetAdminInfoErr, e.Error())
		return
	}
	//获取镜像仓库
	image := new(table.Image)
	if e := s.mysql.GetReadEngine().Where("owner_id = ? AND id = ?", admin.OwnerID, request.Image).First(image).Error; e != nil {
		err = customerror.EnCodeError(define.StartDockerErr, "镜像仓库ID不正确")
		return
	}
	//重启
	name := fmt.Sprintf("%d-%s", admin.OwnerID, request.Name)
	imagesAddress := image.Address + `/` + request.Name + `:` + request.Version
	s._CloseDocker(name)
	go func() {
		e := s._StartDocker(ctx.GetToken(), imagesAddress, name, image.Account, image.Pwd, request.Ports)
		if e != nil {
			log.Errorln(e)
			return
		}
		if s.mysql.GetReadEngine().Where("name = ? AND owner_id = ?", request.Name, admin.OwnerID).First(&table.Docker{}).RecordNotFound() == true {
			ps, _ := json.Marshal(request.Ports)
			//添加记录
			docker := &table.Docker{
				//唯一主键
				ID: s.snowflake.GetID(),
				//Name
				Name: request.Name,
				//编译发布的管理员
				AdminID: admin.AdminID,
				//发布镜像
				Image: imagesAddress,
				//账号
				Account: image.Account,
				//密码
				Pwd: image.Pwd,
				//端口
				Ports: string(ps),
				//业主ID
				OwnerID: admin.OwnerID,
				//注册事件
				Time: time.Now().Unix(),
			}
			if e := s.mysql.GetWriteEngine().Create(&docker).Error; e != nil {
				log.Errorln(e.Error())
			}
		}
	}()
	//生成登录记录
	tbLog := &table.Log{
		//日志ID
		LogID: s.snowflake.GetID(),
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
		URL: ctx.GetURL(),
		//操作时间
		Time: time.Now().Unix(),
	}
	if e := s.mysql.GetWriteEngine().Create(&tbLog).Error; e != nil {
		log.Errorln(e.Error())
	}
	response.Success = true
	return
}

//CloseDocker 关闭docker容器
func (s *Service) CloseDocker(ctx plugins.Context, request param.CloseDockerReq) (response param.CloseDockerRes, err error) {
	admin, e := s.GetAdmin(ctx)
	if e != nil {
		log.Errorln(e.Error())
		err = customerror.EnCodeError(define.GetAdminInfoErr, e.Error())
		return
	}
	if e := s._CloseDocker(request.ID); e != nil {
		log.Errorln(e.Error())
		err = customerror.EnCodeError(define.CloseDockerErr, e.Error())
		return
	}
	//生成登录记录
	tbLog := &table.Log{
		//日志ID
		LogID: s.snowflake.GetID(),
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
		URL: ctx.GetURL(),
		//操作时间
		Time: time.Now().Unix(),
	}
	if e := s.mysql.GetWriteEngine().Create(&tbLog).Error; e != nil {
		log.Errorln(e.Error())
		err = customerror.EnCodeError(define.CloseDockerErr, "关闭服务失败")
		return
	}
	response.Success = true
	return
}

//GetDockerList  获取docker列表
func (s *Service) GetDockerList(ctx plugins.Context, request param.GetDockerListReq) (response param.GetDockerListRes, err error) {
	admin, ok := ctx.GetShareByKey("Admin").(*table.Admin)
	if ok == false {
		err = customerror.EnCodeError(define.GetBuildServiceListErr, "管理员信息失败")
		return
	}
	var dockers []table.Docker
	if e := s.mysql.GetReadEngine().Where("owner_id = ?", admin.OwnerID).Order("time DESC").Find(&dockers).Error; e != nil {
		err = customerror.EnCodeError(define.GetBuildServiceListErr, e.Error())
		return
	}
	mp := make(map[string]types.Container)
	containers, err := s.docker.ContainerList(context.Background(), types.ContainerListOptions{})
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
	}

	total := len(dockers)
	response.PageNo = 1
	response.PageSize = total
	response.TotalCount = total
	if total <= 0 {
		response.TotalPage = 0
		return
	}
	if total%response.PageSize > 0 {
		response.TotalPage = total/response.PageSize + 1
	}
	if total%response.PageSize < 0 {
		response.TotalPage = total / response.PageSize
	}
	for _, docker := range dockers {
		d := param.Docker{
			DockerID: strconv.FormatInt(docker.ID, 10),
			Name:     docker.Name,
			Image:    docker.Image,
			Time:     time.Unix(docker.Time, 0).Format("2006-01-02 15:04:05"),
		}
		container, ok := mp[docker.Name]
		if ok {
			d.RunStatus = true
			d.ID = container.ID
			d.Status = container.Status
			d.Command = container.Command
		}
		p := ""
		for _, port := range container.Ports {
			p = fmt.Sprintf("%s -p %d:%d/%s", p, port.PublicPort, port.PrivatePort, port.Type)
		}
		d.Ports = p
		response.Data = append(response.Data, d)
	}
	return
}
