package controller

import (
	"dockermysql/infra/dao"
	model2 "dockermysql/model"
	"fmt"
	"log"

	"github.com/Shopify/sarama"
	jsoniter "github.com/json-iterator/go"
)

type ConsumerGroupHandler struct {
}

func (ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {

	return nil
}
func (ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	var address = []string{"127.0.0.1:9092"}
	topics := []string{"topic2"}
	producer := InitProducer(address)
	//关闭生产者
	defer producer.Close()

	for msg := range claim.Messages() {
		//接受到查询任务
		log.Printf("need select : topic = %s, partition = %d, value = %s, offset= %d, timestamp = %v, ", msg.Topic, msg.Partition, msg.Value, msg.Offset, msg.Timestamp)
		var queryMsg model2.QueryMsg
		err := jsoniter.Unmarshal(msg.Value, &queryMsg)
		if err != nil {
			fmt.Println("unmarshal error:", err)
			return err
		}

		productlist, err := dao.GetAllproducts(&queryMsg.Condition, queryMsg.Offset, queryMsg.Limit)
		if err != nil {
			fmt.Println(err)
			return err
		}

		value := model2.ExportMsg{
			Productlist: productlist,
			File:        queryMsg.File,
			Row:         queryMsg.Row,
		}
		jsonData, _ := jsoniter.Marshal(value)

		newmsg := &sarama.ProducerMessage{
			Topic: topics[0],
			Value: sarama.ByteEncoder(jsonData),
		}
		//发送查询后的数据
		_, _, err = producer.SendMessage(newmsg)
		if err != nil {
			log.Printf("FAILED to send message: %s\n", err)
		} else {

			log.Printf("export msg send: topic = %s, partition = %d, offset= %d, timestamp = %v, ", newmsg.Topic, newmsg.Partition, newmsg.Offset, newmsg.Timestamp)
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
