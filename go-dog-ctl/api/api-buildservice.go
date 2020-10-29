package api

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/tang-go/go-dog-tool/define"
	"github.com/tang-go/go-dog-tool/go-dog-ctl/param"
	"github.com/tang-go/go-dog-tool/go-dog-ctl/table"
	gateParam "github.com/tang-go/go-dog-tool/go-dog-gw/param"
	customerror "github.com/tang-go/go-dog/error"
	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/plugins"
)

//BuildService 管理员登录
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
	system := runtime.GOOS
	build := ""
	switch system {
	case "darwin":
		build = "CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o " + request.Name
	case "linxu":
		build = "go build -o " + request.Name
	default:
		err = customerror.EnCodeError(define.BuildServiceErr, "目前只支持linux和mac")
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

	name := strings.Replace(paths[l-1], ".git", "", -1)
	if ok, err := pointer._PathExists(name); ok && err == nil {
		log.Traceln("目录存在更新", name)
		go func() {
			logTxt := ""
			pointer._RunInLinux(
				`
				cd `+name+`
				git pull 
				cd `+request.Path+`
				go mod tidy
				`+build+`
				docker build -t `+request.Harbor+`/`+request.Name+`:`+request.Version+` .
				docker push `+request.Harbor+`/`+request.Name+`:`+request.Version+` 	
				rm -rf `+request.Name+`
				echo 编译完成`, func(success string) {
					res := new(gateParam.PushRes)
					ctx.GetClient().Broadcast(ctx, define.SvcGateWay, "Push",
						&gateParam.PushReq{
							Token: ctx.GetToken(),
							Topic: define.BuildServiceTopic,
							Msg:   success,
						},
						res)
					logTxt = logTxt + success + `<p/>`
					fmt.Println(success)
				}, func(err string) {
					res := new(gateParam.PushRes)
					ctx.GetClient().Broadcast(
						ctx,
						define.SvcGateWay,
						"Push",
						&gateParam.PushReq{
							Token: ctx.GetToken(),
							Topic: define.BuildServiceTopic,
							Msg:   err,
						},
						res)
					logTxt = logTxt + err + `<p/>`
					fmt.Println(err)
				})
			//完成
			err := pointer.mysql.GetWriteEngine().Model(&table.BuildService{}).Where("id = ?", tbBuild.ID).Update(
				map[string]interface{}{
					"Log":    logTxt,
					"Status": true,
				}).Error
			if err != nil {
				log.Errorln(err.Error())
			}
		}()

	} else {
		log.Traceln("目录不存在直接下载", name)
		go func() {
			logTxt := ""
			pointer._RunInLinux(
				`
				git clone `+request.Git+`
				cd `+name+`
				cd `+request.Path+`
				go mod tidy
				`+build+`
				docker build -t `+request.Harbor+`/`+request.Name+`:`+request.Version+` .
				docker push `+request.Harbor+`/`+request.Name+`:`+request.Version+` 	
				rm -rf `+request.Name+`
				echo 编译完成`, func(success string) {
					res := new(gateParam.PushRes)
					ctx.GetClient().Broadcast(ctx, define.SvcGateWay, "Push",
						&gateParam.PushReq{
							Token: ctx.GetToken(),
							Topic: define.BuildServiceTopic,
							Msg:   success,
						},
						res)
					fmt.Println(success)
					logTxt = logTxt + success + `<p/>`
				}, func(err string) {
					res := new(gateParam.PushRes)
					ctx.GetClient().Broadcast(
						ctx,
						define.SvcGateWay,
						"Push",
						&gateParam.PushReq{
							Token: ctx.GetToken(),
							Topic: define.BuildServiceTopic,
							Msg:   err,
						},
						res)
					fmt.Println(err)
					logTxt = logTxt + err + `<p/>`
				})
			//完成
			err := pointer.mysql.GetWriteEngine().Model(&table.BuildService{}).Where("id = ?", tbBuild.ID).Update(
				map[string]interface{}{
					"Log":    logTxt,
					"Status": true,
				}).Error
			if err != nil {
				log.Errorln(err.Error())
			}
		}()
	}
	response.Success = true
	return
}
