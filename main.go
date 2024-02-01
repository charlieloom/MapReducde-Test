package main

import (
	"dockermysql/controller"
	"dockermysql/dal"
	"dockermysql/dal/query"
	"dockermysql/routers"
	"os"
)

const MySQLDSN = "root:123456@tcp(127.0.0.1:3307)/Shop?charset=utf8mb4&parseTime=True"

func init() {
	dal.DB = dal.ConnectDB(MySQLDSN).Debug()
}

func main() {
	port := os.Args[1]
	query.SetDefault(dal.DB)
	var address = []string{"127.0.0.1:9092"}
	topics := []string{"topic1", "topic2", "topic3"}
	// go controller.Produce(address, topics)

	go controller.ConsumeGroup(address, topics)

	// go controller.Consume(address, topics[0])
	r := routers.SetupRouter()
	r.Run(":" + port)
}
