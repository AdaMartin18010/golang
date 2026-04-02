# Kafka 客户端

> **分类**: 开源技术堆栈

---

## sarama

```go
import "github.com/IBM/sarama"
```

---

## 生产者

```go
config := sarama.NewConfig()
config.Producer.RequiredAcks = sarama.WaitForAll
config.Producer.Retry.Max = 5
config.Producer.Return.Successes = true

producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, config)
if err != nil {
    log.Fatal(err)
}
defer producer.Close()

msg := &sarama.ProducerMessage{
    Topic: "my-topic",
    Key:   sarama.StringEncoder("key"),
    Value: sarama.StringEncoder("hello world"),
}

partition, offset, err := producer.SendMessage(msg)
```

---

## 消费者

```go
config := sarama.NewConfig()
config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin

consumer, err := sarama.NewConsumerGroup([]string{"localhost:9092"}, "my-group", config)
if err != nil {
    log.Fatal(err)
}

type handler struct{}

func (h handler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (h handler) Cleanup(sarama.ConsumerGroupSession) error { return nil }
func (h handler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
    for msg := range claim.Messages() {
        fmt.Printf("Message: %s\n", string(msg.Value))
        session.MarkMessage(msg, "")
    }
    return nil
}

ctx := context.Background()
consumer.Consume(ctx, []string{"my-topic"}, handler{})
```
