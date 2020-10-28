package api

import (
	"time"

	"github.com/tang-go/go-dog-tool/define"
	"github.com/tang-go/go-dog-tool/go-dog-ctl/param"
	"github.com/tang-go/go-dog-tool/go-dog-ctl/table"
	customerror "github.com/tang-go/go-dog/error"
	"github.com/tang-go/go-dog/plugins"
)

//GetBuildServiceList 获取编译发布记录
func (pointer *API) GetBuildServiceList(ctx plugins.Context, request param.GetBuildServiceReq) (response param.GetBuildServiceRes, err error) {
	admin, ok := ctx.GetShareByKey("Admin").(*table.Admin)
	if ok == false {
		err = customerror.EnCodeError(define.GetBuildServiceListErr, "管理员信息失败")
		return
	}
	var builds []table.BuildService
	if e := pointer.mysql.GetReadEngine().Where("owner_id = ?", admin.OwnerID).Limit(request.PageSize).Offset((request.PageNo - 1) * request.PageSize).Order("time DESC").Find(&builds).Error; e != nil {
		err = customerror.EnCodeError(define.GetBuildServiceListErr, e.Error())
		return
	}
	total := 0
	if e := pointer.mysql.GetReadEngine().Model(&table.BuildService{}).Where("owner_id = ?", admin.OwnerID).Order("time DESC").Count(&total).Error; e != nil {
		err = customerror.EnCodeError(define.GetBuildServiceListErr, e.Error())
		return
	}
	response.PageNo = request.PageNo
	response.PageSize = request.PageSize
	if total%request.PageSize > 0 {
		response.TotalPage = total/request.PageSize + 1
	}
	if total%request.PageSize < 0 {
		response.TotalPage = total / request.PageSize
	}
	response.TotalCount = total
	for _, build := range builds {
		response.Data = append(response.Data, param.BuildService{
			//唯一主键
			ID: build.ID,
			//编译发布的管理员
			AdminID: build.AdminID,
			//发布镜像
			Image: build.Image,
			//状态
			Status: build.Status,
			//执行日志
			Log: build.Log,
			//业主ID
			OwnerID: build.OwnerID,
			//角色创建时间
			Time: time.Unix(build.Time, 0).Format("2006-01-02 15:04:05"),
		})
	}
	return
}
