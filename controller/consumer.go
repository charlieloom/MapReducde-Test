package controller

import (
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
		fmt.Printf("start to get message from partition %v \n", partition)
		for msg := range partitionconsumer.Messages() {
			fmt.Printf("topic=%s partition=%d, offset=%d, key=%v, value=%s \n", msg.Topic, msg.Partition, msg.Offset, msg.Key, string(msg.Value))
		}
	}

}
