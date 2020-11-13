package api

import (
	"strconv"
	"strings"
	"time"

	"github.com/tang-go/go-dog-tool/define"
	"github.com/tang-go/go-dog-tool/go-dog-ctl/param"
	"github.com/tang-go/go-dog-tool/go-dog-ctl/table"
	customerror "github.com/tang-go/go-dog/error"
	"github.com/tang-go/go-dog/log"
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
			ID: strconv.FormatInt(build.ID, 10),
			//发布镜像
			Image: build.Image,
			//状态
			Status: build.Status,
			//执行日志
			Log: build.Log,
			//角色创建时间
			Time: time.Unix(build.Time, 0).Format("2006-01-02 15:04:05"),
		})
	}
	return
}

//BuildService 编译发布docker镜像
func (pointer *API) BuildService(ctx plugins.Context, request param.BuildServiceReq) (response param.BuildServiceRes, err error) {
	admin, ok := ctx.GetShareByKey("Admin").(*table.Admin)
	if ok == false {
		err = customerror.EnCodeError(define.GetAdminInfoErr, "管理员信息失败")
		return
	}
	paths := strings.Split(request.Git, "/")
	l := len(paths)
	if l <= 0 {
		err = customerror.EnCodeError(define.BuildServiceErr, "路径不正确")
		return
	}
	//添加编译记录
	tbBuild := table.BuildService{
		ID:      pointer.snowflake.GetID(),
		AdminID: admin.AdminID,
		Status:  false,
		Image:   request.Harbor + "/" + request.Name + ":" + request.Version,
		OwnerID: admin.OwnerID,
		Time:    time.Now().Unix(),
	}
	//生成登录记录
	tbLog := &table.Log{
		//日志ID
		LogID: pointer.snowflake.GetID(),
		//类型
		Type: table.BuildServiceType,
		//操作人
		AdminID: admin.AdminID,
		//名称
		AdminName: admin.Name,
		//操作方法
		Method: "BuildService",
		//描述
		Description: "编译发布服务",
		//业主ID
		OwnerID: admin.OwnerID,
		//操作IP
		IP: ctx.GetAddress(),
		//操作URL
		URL: ctx.GetDataByKey("URL").(string),
		//操作时间
		Time: time.Now().Unix(),
	}
	//开启数据库操作
	tx := pointer.mysql.GetWriteEngine().Begin()
	if e := tx.Create(&tbBuild).Error; e != nil {
		tx.Rollback()
		log.Errorln(e.Error())
		err = customerror.EnCodeError(define.BuildServiceErr, "编译发布失败")
		return
	}
	if e := tx.Create(&tbLog).Error; e != nil {
		tx.Rollback()
		log.Errorln(e.Error())
		err = customerror.EnCodeError(define.BuildServiceErr, "编译发布失败")
		return
	}
	tx.Commit()
	go pointer._SendEvent(tbBuild.ID, ctx, &request)
	response.Success = true
	return
}
