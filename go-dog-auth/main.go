package main

import (
	"github.com/tang-go/go-dog-tool/go-dog-auth/router"
	"github.com/tang-go/go-dog-tool/go-dog-auth/service"

	"github.com/tang-go/go-dog/log"
)

func main() {
	s := service.NewService(router.RPCRouter)
	if e := s.Run(); e != nil {
		log.Errorln(e.Error())
	}
}
