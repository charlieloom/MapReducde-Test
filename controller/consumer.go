package controller

import (
	"context"
	"fmt"
	"log"

	"github.com/Shopify/sarama"
)

type consumerGroupHandler struct {
}

func (consumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}
func (consumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Printf("Message claimed: topic = %s, partition = %d, value = %s, offset= %d, timestamp = %v, ", msg.Topic, msg.Partition, msg.Value, msg.Offset, msg.Timestamp)
		session.MarkMessage(msg, "")
	}
	return nil
}
func ConsumeGroup(brokers []string, topics []string) {
	// 配置
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = false
	config.Version = sarama.V0_10_2_0                     // specify appropriate version
	config.Consumer.Offsets.Initial = sarama.OffsetOldest // 未找到组消费位移的时候从哪边开始消费

	//创建消费组
	consumergroup, err := sarama.NewConsumerGroup(brokers, "consumer-group", config)
	if err != nil {
		fmt.Printf("new consumergroup error: %s\n", err.Error())
		return
	}
	defer consumergroup.Close()

	consumergroup.Consume(context.Background(), topics, consumerGroupHandler{})
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
