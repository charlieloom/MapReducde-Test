package main

import (
	"dockermysql/controller"
	"dockermysql/infra/dao"
	"dockermysql/infra/mq"
	"dockermysql/routers"
	"os"
)

func Init(isExport string) {
	dao.Init()
	if isExport == "1" {
		go mq.InitConsumer([]string{"127.0.0.1:9092"}, []string{"topic2"}, &controller.ExportData{})
	} else {
		go mq.InitConsumer([]string{"127.0.0.1:9092"}, []string{"topic1"}, &controller.ConsumerGroupHandler{})
	}
}

func main() {
	Init(os.Args[2])
	port := os.Args[1]
	r := routers.SetupRouter()
	r.Run(":" + port)
}
