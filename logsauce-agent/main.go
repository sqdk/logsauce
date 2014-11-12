package main

import (
	"github.com/sqdk/logsauce"
)

func main() {
	config, _ := logsauce.ReadConfig("./client.conf")

	logsauce.WatchFiles(config.ClientConfiguration.FilesToWatch, true, config.ClientConfiguration.ServerAddress, config.ClientConfiguration.ClientToken)
}
