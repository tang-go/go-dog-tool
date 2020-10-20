package main

import "github.com/tang-go/go-dog-tool/go-dog-gw/gateway"

func main() {
	gate := gateway.NewGateway()
	gate.Run()
}
