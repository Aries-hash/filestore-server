package main

import (
	"filestore-server/config"
	"filestore-server/mq"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/micro/cli"
	micro "github.com/micro/go-micro"
	_ "github.com/micro/go-plugins/registry/kubernetes"

	"filestore-server/common"
	dbproxy "filestore-server/service/dbproxy/client"
	cfg "filestore-server/service/upload/config"
	upProto "filestore-server/service/upload/proto"
	"filestore-server/service/upload/route"
	upRpc "filestore-server/service/upload/rpc"
)

func startRPCService() {
	service := micro.NewService(
		micro.Name("go.micro.service.upload"), // 服务名称
		micro.RegisterTTL(time.Second*10),     // TTL指定从上一次心跳间隔起，超过这个时间服务会被服务发现移除
		micro.RegisterInterval(time.Second*5), // 让服务在指定时间内重新注册，保持TTL获取的注册时间有效
		micro.Flags(common.CustomFlags...),
	)
	service.Init(
		micro.Action(func(c *cli.Context) {
			// 检查是否指定mqhost
			mqhost := c.String("mqhost")
			if len(mqhost) > 0 {
				log.Println("custom mq address: " + mqhost)
				mq.UpdateRabbitHost(mqhost)
			}
		}),
	)

	// 初始化dbproxy client
	dbproxy.Init(service)
	// 初始化mq client
	mq.Init()

	upProto.RegisterUploadServiceHandler(service.Server(), new(upRpc.Upload))
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}

func startAPIService() {
	router := route.Router()
	router.Run(cfg.UploadServiceHost)
	// service := web.NewService(
	// 	web.Name("go.micro.web.upload"),
	// 	web.Handler(router),
	// 	web.RegisterTTL(10*time.Second),
	// 	web.RegisterInterval(5*time.Second),
	// )
	// if err := service.Init(); err != nil {
	// 	log.Fatal(err)
	// }

	// if err := service.Run(); err != nil {
	// 	log.Fatal(err)
	// }
}

func main() {
	os.MkdirAll(config.TempLocalRootDir, 0777)
	os.MkdirAll(config.TempPartRootDir, 0777)

	// api 服务
	go startAPIService()

	// rpc 服务
	startRPCService()
}
