package api

import (
	"bytes"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/tang-go/go-dog-tool/define"
	"github.com/tang-go/go-dog-tool/go-dog-ctl/param"
	customerror "github.com/tang-go/go-dog/error"
	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/plugins"
)

//BuildService 管理员登录
func (pointer *API) BuildService(ctx plugins.Context, request param.BuildServiceReq) (response param.BuildServiceRes, err error) {
	paths := strings.Split(request.Git, "/")
	l := len(paths)
	if l <= 0 {
		err = customerror.EnCodeError(define.GetAdminInfoErr, "路径不正确")
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
		response.Result = "目前只支持linux和mac"
		return
	}
	name := strings.Replace(paths[l-1], ".git", "", -1)
	if ok, err := pointer._PathExists(name); ok && err == nil {
		response.Result = pointer._RunInLinux(`
		cd ` + name + `
		git pull 
		cd ` + request.Path + `
		go mod tidy
		` + build + `
		docker build -t ` + request.Harbor + `/` + request.Name + `:` + request.Version + ` .
		docker push ` + request.Harbor + `/` + request.Name + `:` + request.Version + ` 	
		rm -rf ` + request.Name + `
		echo 编译完成`)
	} else {
		response.Result = pointer._RunInLinux(`
		git clone ` + request.Git + `
		cd ` + name + `
		cd ` + request.Path + `
		go mod tidy
		` + build + `
		docker build -t ` + request.Harbor + `/` + request.Name + `:` + request.Version + ` .
		docker push ` + request.Harbor + `/` + request.Name + `:` + request.Version + ` 	
		rm -rf ` + request.Name + `
		echo 编译完成`)
	}
	log.Traceln(response.Result)
	return
}

func (pointer *API) _RunInLinux(cmd string) string {
	var out bytes.Buffer
	//var stderr bytes.Buffer
	c := exec.Command("sh", "-c", cmd)
	//c.Stdin = os.Stdin
	c.Stdout = &out
	//c.Stderr = &stderr
	c.Start()
	c.Wait()
	return out.String()
}

func (pointer *API) _PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
