package main

import (
	"fmt"
	"github.com/sqdk/logsauce"
)

func main() {
	config, _ := logsauce.ReadConfig("./server.conf")

	logsauce.InitializeDB(config)
	logsauce.RegisterRoutes(config.ListenPort, config.Relaymode, config.ServerMode)
	fmt.Println("Starting server")

	looper := make(chan int)

	<-looper
}
