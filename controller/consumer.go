package controller

import (
	"dockermysql/infra/dao"
	model2 "dockermysql/model"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Shopify/sarama"
)

type ConsumerGroupHandler struct {
}

func (ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {

	return nil
}
func (ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Printf("Message claimed: topic = %s, partition = %d, value = %s, offset= %d, timestamp = %v, ", msg.Topic, msg.Partition, msg.Value, msg.Offset, msg.Timestamp)
		var queryMsg model2.QueryMsg
		err := json.Unmarshal([]byte(msg.Value), &queryMsg)
		if err != nil {
			fmt.Println("unmarshal error:", err)
			return err
		}
		productlist, err := dao.GetAllproducts(&queryMsg.Condition, queryMsg.Offset, queryMsg.Limit)
		if err != nil {
			fmt.Println(err)
			return err
		}
		var address = []string{"127.0.0.1:9092"}
		topics := []string{"topic2"}
		producer := InitProducer(address)
		//关闭生产者
		defer producer.Close()
		value := model2.ExportMsg{
			Productlist: productlist,
			File:        queryMsg.File,
			Row:         queryMsg.Row,
		}
		jsonData, _ := json.Marshal(value)
		newmsg := &sarama.ProducerMessage{
			Topic: topics[0],
			Value: sarama.ByteEncoder(jsonData),
		}
		//发送消息
		part, offset, err := producer.SendMessage(newmsg)
		if err != nil {
			log.Printf("send querymsgage error :%s \n", err.Error())
		} else {
			fmt.Printf("export msg send SUCCESS:topic=%s  patition=%d, offset=%d \n", topics[0], part, offset)
		}
		session.MarkMessage(msg, "")
	}
	return nil
}

func Consume(brokers []string, topic string) {
	consumer, err := sarama.NewConsumer(brokers, sarama.NewConfig())

	if err != nil {
		panic(err)
	}

	defer consumer.Close()

	partitionList, err := consumer.Partitions(topic)
	if err != nil {
		fmt.Printf("fail to get list of partition:err %v \n", err)
		return
	}
	fmt.Printf("topic=%s get partition: %v", topic, partitionList)

	//遍历所有分区
	for partition := range partitionList {
		//创建分区消费者
		partitionconsumer, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			fmt.Printf("failed to start consumer for partition %d, err:%v \n", partition, err)
		}

		defer partitionconsumer.AsyncClose()

		//从每个分区消费信息
		fmt.Printf("start to get querymsgage from partition %v \n", partition)
		for msg := range partitionconsumer.Messages() {
			fmt.Printf("topic=%s partition=%d, offset=%d, key=%v, value=%s \n", msg.Topic, msg.Partition, msg.Offset, msg.Key, string(msg.Value))
		}
	}

}
