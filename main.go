package main

import (
	"admin-server/database"
	"admin-server/pkg/config"
	"admin-server/pkg/util"
	"admin-server/routers"
	"fmt"
	"github.com/fvbock/endless"
	"log"
	"net/http"
	"runtime"
	"syscall"
)

func main() {
	config.Setup()
	database.Setup()
	defer database.Eloquent.Close()
	util.Setup()

	routersInit := routers.InitRouter()
	readTimeout := config.Conf.Server.ReadTimeout
	writeTimeout := config.Conf.Server.WriteTimeout
	endPoint := fmt.Sprintf(":%d", config.Conf.Server.HttpPort)
	maxHeaderBytes := 1 << 20

	if runtime.GOOS == "windows" {
		server := &http.Server{
			Addr:           endPoint,
			Handler:        routersInit,
			ReadTimeout:    readTimeout,
			WriteTimeout:   writeTimeout,
			MaxHeaderBytes: maxHeaderBytes,
		}

		server.ListenAndServe()
		return
	}

	endless.DefaultReadTimeOut = readTimeout
	endless.DefaultWriteTimeOut = writeTimeout
	endless.DefaultMaxHeaderBytes = maxHeaderBytes
	server := endless.NewServer(endPoint, routersInit)
	server.BeforeBegin = func(add string) {
		log.Printf("Actual pid is %d", syscall.Getpid())
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Printf("Server err: %v", err)
	}
}
