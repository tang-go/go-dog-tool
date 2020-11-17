package main

import (
	"github.com/tang-go/go-dog-tool/go-dog-ctl/router"
	"github.com/tang-go/go-dog-tool/go-dog-ctl/service"
)

func main() {
	s := service.NewService(router.GETRouter, router.POSTRouter, router.RPCRouter)
	s.Run()
}
