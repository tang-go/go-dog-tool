package api

import (
	"os"
	"os/exec"

	"github.com/tang-go/go-dog-tool/go-dog-ctl/param"
	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/plugins"
)

// docker build -t ` + harbor + ` .
// docker push ` + harbor + `

//BuildService 管理员登录
func (pointer *API) BuildService(ctx plugins.Context, request param.BuildServiceReq) (response param.BuildServiceRes, err error) {
	log.Traceln("收到编译推送镜像需求")
	pointer._RunInLinux(`
		git clone ` + request.Git + `
		cd ` + request.Path + `
		go build -o ` + request.Name + `
		echo "开始打包docker"
		docker build -t ` + request.Harbor + `/` + request.Name + `:` + request.Version + ` .
		rm -rf ` + request.Name + `
		echo 编译完成！！！！！！`)
	return
}

func (pointer *API) _RunInLinux(cmd string) {
	c := exec.Command("sh", "-c", cmd)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Start()
	c.Wait()
}
