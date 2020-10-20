package main

import (
	"github.com/tang-go/go-dog-tool/go-dog-ctl/api"
)

func main() {
	s := api.NewService()
	s.Run()
}
