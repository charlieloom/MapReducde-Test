package mq

import (
	"context"
	"fmt"

	"github.com/Shopify/sarama"
)

func InitConsumer(brokers []string, topics []string, handler sarama.ConsumerGroupHandler) {
	// 配置
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = false
	config.Version = sarama.V0_10_2_0                     // specify appropriate version
	config.Consumer.Offsets.Initial = sarama.OffsetNewest // 未找到组消费位移的时候从哪边开始消费

	//创建消费组
	consumergroup, err := sarama.NewConsumerGroup(brokers, "consumer-group", config)
	if err != nil {
		panic(fmt.Sprintf("new consumergroup error: %s\n", err.Error()))
	}
	defer consumergroup.Close()

	for {
		err := consumergroup.Consume(context.Background(), topics, handler)
		if err != nil {
			fmt.Printf("consumer error=%s topics=%v\n", err.Error(), topics)
			return
		}
	}
}
