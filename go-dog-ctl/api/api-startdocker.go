package api

import (
	"fmt"

	"github.com/tang-go/go-dog-tool/go-dog-ctl/param"
	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/plugins"
)

//StartDocker 以docker模式启动docker
func (pointer *API) StartDocker(ctx plugins.Context, request param.StartDockerReq) (response param.StartDockerRes, err error) {
	docker := fmt.Sprintf("docker run -d --restart=always --name=%s", request.Name)
	for _, port := range request.Ports {
		docker = fmt.Sprintf("%s -p %d:%d", docker, port.ExternalPort, port.InsidePort)
	}
	docker = fmt.Sprintf("%s %s", docker, request.Images)
	log.Traceln(docker)
	// go pointer._RunInLinux(
	// 	`
	// 	docker kill `+request.Name+`
	// 	docker rm `+request.Name+`
	// 	`+docker+`
	// 	`, func(success string) {

	// 	}, func(err string) {

	// 	})
	response.Success = true
	return
}
