package main

import (
	"dockermysql/controller"
	"dockermysql/infra/dao"
	"dockermysql/infra/mq"
	"dockermysql/routers"
	"os"
)

func Init() {
	dao.Init()

	go mq.InitConsumer([]string{"127.0.0.1:9092"}, []string{"topic1", "topic2", "topic3"}, &controller.ConsumerGroupHandler{})
	// other consumer
	// go mq.InitConsumer([]string{"127.0.0.1:9092"}, []string{"topic1", "topic2", "topic3"},)
}

func main() {
	Init()

	port := os.Args[1]
	r := routers.SetupRouter()
	r.Run(":" + port)
}
