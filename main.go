package main

import (
	"fmt"
	"github.com/fvbock/endless"
	"go-admin-starter/middleware/jwt"
	"go-admin-starter/routers"
	"go-admin-starter/utils"
	"go-admin-starter/utils/config"
	"log"
	"net/http"
	"runtime"
	"syscall"
)

func main() {
	conf := config.New()
	utils.LogSetup()
	jwt.SetSignKey(conf.App.JwtSecret)

	routersInit := routers.InitRouter()
	readTimeout := conf.Server.ReadTimeout
	writeTimeout := conf.Server.WriteTimeout
	endPoint := fmt.Sprintf(":%d", conf.Server.HttpPort)
	maxHeaderBytes := 1 << 20

	log.Printf("Server start at http port: %d", conf.Server.HttpPort)

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
