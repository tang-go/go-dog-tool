package api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/tang-go/go-dog-tool/define"
	"github.com/tang-go/go-dog-tool/go-dog-ctl/param"
	"github.com/tang-go/go-dog-tool/go-dog-ctl/table"
	gateParam "github.com/tang-go/go-dog-tool/go-dog-gw/param"
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

//BuildImage 编译镜像
func (pointer *API) BuildImage(tarFile, project, imageName string) error {
	dockerBuildContext, err := os.Open(tarFile)
	if err != nil {
		return err
	}
	defer dockerBuildContext.Close()

	buildOptions := types.ImageBuildOptions{
		Dockerfile: "Dockerfile", // optional, is the default
		Tags:       []string{imageName},
		Labels: map[string]string{
			project: "project",
		},
	}
	output, err := pointer.docker.ImageBuild(context.Background(), dockerBuildContext, buildOptions)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(output.Body)
	if err != nil {
		return err
	}
	if strings.Contains(string(body), "error") {
		return fmt.Errorf("build image to docker error")
	}
	return nil
}

func PushImage(cli *client.Client, registryUser, registryPassword, image string) error {
	authConfig := types.AuthConfig{
		Username: registryUser,
		Password: registryPassword,
	}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		return err
	}
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)
	out, err := cli.ImagePush(context.TODO(), image, types.ImagePushOptions{RegistryAuth: authStr})
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(out)
	if err != nil {
		return err
	}
	if strings.Contains(string(body), "error") {
		return fmt.Errorf("push image to docker error")
	}
	return nil
}
