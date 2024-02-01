package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/Shopify/sarama"
	"github.com/gin-gonic/gin"
)

func Produce(c *gin.Context) {

	var address = []string{"127.0.0.1:9092"}
	topics := []string{"topic1", "topic2", "topic3"}

	// 配置
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Timeout = 5 * time.Second
	producer, err := sarama.NewSyncProducer(address, config)
	if err != nil {
		log.Printf("new sync producer error : %s \n", err.Error())
		return
	}
	//关闭生产者
	defer producer.Close()

	//循环发送消息
	for i := 0; i < 10; i++ {
		//创建消息
		pos := rand.Intn(3)
		type Mess struct {
			Id   int    `json:"id"`
			Name string `json:"name"`
		}
		value := Mess{i, "name"}
		jsonData, _ := json.Marshal(value)
		msg := &sarama.ProducerMessage{
			Topic: topics[pos],
			Value: sarama.ByteEncoder(jsonData),
		}

		//发送消息
		part, offset, err := producer.SendMessage(msg)
		if err != nil {
			log.Printf("send message error :%s \n", err.Error())
		} else {
			fmt.Printf("SUCCESS:topic=%s  patition=%d, value=%v, offset=%d \n", topics[pos], part, value, offset)
		}
		time.Sleep(2 * time.Second)
	}
}
