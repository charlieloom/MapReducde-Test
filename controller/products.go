package controller

import (
	model "dockermysql/dal/model"
	"dockermysql/infra/dao"
	model2 "dockermysql/model"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Shopify/sarama"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/xuri/excelize/v2"
)

func PageQuery(c *gin.Context) {
	var condition model2.Condition
	condition.Id, _ = strconv.Atoi(c.Query("id"))
	condition.Name = c.Query("name")
	condition.Category = c.Query("category")
	condition.Page, _ = strconv.Atoi(c.Query("page"))
	condition.PageSize, _ = strconv.Atoi(c.Query("pageSize"))
	condition.Sort = c.Query("sort")

	start := time.Now()
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	selectDone := 0
	var err error
	var productlists []*model.Product
	offest := (condition.Page - 1) * condition.PageSize
	for {
		productlist, err := dao.GetAllproducts(&condition, offest+selectDone, 500000)
		if err != nil {
			fmt.Println(err)
			return
		}
		if len(productlists) <= 0 {
			break
		}
		productlists = append(productlists, productlist...)
		selectDone += len(productlists)
		if selectDone >= condition.PageSize {
			break
		}
	}
	cost := time.Since(start)
	fmt.Printf("花费时间：[%s]", cost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": productlists,
	})
}

func ExportProducts(c *gin.Context) {
	var condition model2.Condition
	condition.Id, _ = strconv.Atoi(c.Query("id"))
	condition.Name = c.Query("name")
	condition.Category = c.Query("category")
	condition.Page, _ = strconv.Atoi(c.Query("page"))
	condition.PageSize, _ = strconv.Atoi(c.Query("pageSize")) //总共要查询的数量
	condition.Sort = c.Query("sort")

	var address = []string{"127.0.0.1:9092"}
	topics := []string{"topic1", "topic2", "topic3"}

	producer := InitProducer(address)
	//关闭生产者
	defer producer.Close()

	log.Println("开始时间")
	start := time.Now()
	f := excelize.NewFile()
	// 保存文件
	if err := f.SaveAs("products.xlsx"); err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	var err error
	offest := (condition.Page - 1) * condition.PageSize //初始偏移
	for i := 0; i < condition.PageSize; i += 200 {
		value := model2.QueryMsg{
			Condition: condition,
			Offset:    i + offest, //当前偏移
			Limit:     200,
			Row:       i,
			File:      "products.xlsx",
		}

		jsonData, _ := jsoniter.Marshal(value)
		msg := &sarama.ProducerMessage{
			Topic: topics[0],
			Value: sarama.ByteEncoder(jsonData),
		}
		//发送消息
		producer.Input() <- msg
	}
	cost := time.Since(start)
	fmt.Printf("花费时间：[%s]", cost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": "ok",
	})
}
