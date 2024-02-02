package main

import (
	"dockermysql/controller"
	"dockermysql/infra/dao"
	"dockermysql/infra/mq"
	"dockermysql/routers"
	"fmt"
	"os"

	"github.com/xuri/excelize/v2"
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
	_, err := excelize.OpenFile("products.xlsx")
	// f := excelize.NewFile()
	if err != nil {
		fmt.Println("mian Error:", err)
		fmt.Println("file")
		return
	}
	Init(os.Args[2])
	port := os.Args[1]
	r := routers.SetupRouter()
	r.Run(":" + port)
}
