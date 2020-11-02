package main

import "github.com/tang-go/go-dog-tool/go-dog-ctl/api"

func main() {
	api := api.NewService()
	api.Run()
}
