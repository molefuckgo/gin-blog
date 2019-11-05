package main

import (
	"fmt"
	"github.com/fvbock/endless"
	"github.com/molefuckgo/gin-blog/pkg/setting"
	"github.com/molefuckgo/gin-blog/routers"
	"log"
	"syscall"
)

func main() {
	//router := routers.InitRouter()
	//
	//s := &http.Server{
	//	Addr:         fmt.Sprintf(":%d", setting.HTTPPort),
	//	Handler:      router,
	//	ReadTimeout:  setting.ReadTimeout,
	//	WriteTimeout: setting.WriteTimeout,
	//	// MaxHeaderBytes控制服务器读取的最大字节数，以解析请求标头的键和值，
	//	// 包括请求行。 它不限制请求主体的大小。
	//	// 如果为零，则使用DefaultMaxHeaderBytes。1 << 20 // 1 MB
	//	MaxHeaderBytes: 1 << 20,
	//}
	//
	//s.ListenAndServe()

	endless.DefaultReadTimeOut = setting.ReadTimeout
	endless.DefaultWriteTimeOut = setting.WriteTimeout
	endless.DefaultMaxHeaderBytes = 1 << 20
	endPoint := fmt.Sprintf(":%d", setting.HTTPPort)

	server := endless.NewServer(endPoint, routers.InitRouter())
	server.BeforeBegin = func(add string) {
		log.Printf("Actual pid is %d", syscall.Getpid())
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Printf("Server err %v", err)
	}
}
