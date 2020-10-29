package api

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/tang-go/go-dog-tool/define"
	"github.com/tang-go/go-dog-tool/go-dog-ctl/param"
	"github.com/tang-go/go-dog-tool/go-dog-ctl/table"
	customerror "github.com/tang-go/go-dog/error"
	"github.com/tang-go/go-dog/plugins"
)

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
		p := ""
		for _, port := range container.Ports {
			p = fmt.Sprintf("%s -p %d:%d/%s", p, port.PublicPort, port.PrivatePort, port.Type)
		}
		response.Data = append(response.Data, param.Docker{
			ID:      container.ID,
			Name:    container.Names[0],
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
