package rabbitmq

import (
	"github.com/streadway/amqp"
	"log"
	"commerce-hsz/datamodels"
	"encoding/json"
	"commerce-hsz/services"
	"sync"
)

const MQURL = "amqp://hszz:hszz@127.0.0.1:5672/hszz"

type RabbitMQ struct {
	conn *amqp.Connection
	channel *amqp.Channel
	// 队列名
	QueueName string
	// 交换机名
	Exchange string
	// key
	Key string
	MqURL string
	sync.Mutex
}

// 创建结构体实例
func NewRabbitMQ(queueName, exchange, key string) *RabbitMQ {
	return &RabbitMQ{
		QueueName: queueName,
		Exchange: exchange,
		Key: key,
		MqURL: MQURL,
	}
}

// 断开channel 和 connection
func (r *RabbitMQ) Destory() {
	r.channel.Close()
	r.conn.Close()
}

// 创建简单模式下的RabbitMQ实例
func NewRabbitMQSimple(queue string) *RabbitMQ {
	rmq := NewRabbitMQ(queue, "", "")
	// 获取connection
	var err error
	rmq.conn, err = amqp.Dial(rmq.MqURL)
	if err != nil {
		log.Printf("failed to connect rmq:%s", err)
		return nil
	}
	// 获取channel
	rmq.channel, err = rmq.conn.Channel()
	if err != nil {
		log.Printf("failed to open channel:%s", err)
		return nil
	}
	return rmq
}
// 默认模式(default)队列生产
func (r *RabbitMQ) PublishSimple(message string) error {
	// 申请队列, 如果队列不存在会自动创建, 在则跳过创建
	_, err := r.channel.QueueDeclare(
		r.QueueName,
		// 是否持久化
		false,
		// 是否自动删除
		false,
		// 是否具有排他性
		false,
		// 是否阻塞处理
		false,
		// 额外的属性
		nil,
	)
	if err != nil {
		log.Println("PublicSimple")
		log.Printf("申请队列失败:%s",err)
		return err
	}
	// 调用channel 发送消息到队列中
	r.channel.Publish(
		r.Exchange,
		r.QueueName,
		// 如果为true, 根据自身exchange类型和routekey规则无法找到符合条件的队列会把消息返还给发送者
		false,
		//如果为true，当exchange发送消息到队列后发现队列上没有消费者，则会把消息返还给发送者
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body: []byte(message),
		})
	return nil
}
// simple模式下消费
func (r *RabbitMQ) ConsumeSimple(orderService services.IOrderService, productService services.IProductService)  {
	// 申请队列, 如果队列不存在会自动创建, 在则跳过创建
	q, err := r.channel.QueueDeclare(
		r.QueueName,
		// 是否持久化
		false,
		// 是否自动删除
		false,
		// 是否具有排他性
		false,
		// 是否阻塞处理
		false,
		// 额外的属性
		nil,
	)
	if err != nil {
		log.Printf("申请队列失败:%s",err)
		return
	}
	// 消费者流量控制
	r.channel.Qos(
		1, // 当前消费者一次能接受的最大消息数量
		0, // 服务器传递的最大容量(以8字节为单位)
		false, // 如果设置为true 对channel可用
	)

	// 接收消息
	msgs, err := r.channel.Consume(
		// queue
		q.Name,
		// 用来区分多个消费者
		"",
		// 是否自动答应
		// 改为手动应答
		false,
		// 是否排他
		false,
		// 如果为true， 表示生产者和消费者不能是一个connect
		false,
		// 是否阻塞
		false,
		nil,
	)
	if err != nil {
		log.Printf("接收消息失败:%s",err)
	}

	forever := make(chan bool)
	// 启用协程处理消息
	go func() {
		for d := range msgs {
			// 消息逻辑处理(暂时简单打印)
			log.Printf("接收到消息: %s", d.Body)

			message := &datamodels.Message{}
			err := json.Unmarshal([]byte(d.Body), message)
			if err != nil {
				log.Println(err)
			}

			// 插入订单
			_, err = orderService.InsertOrderByMessage(message)
			if err != nil {
				log.Println(err)
			}

			// 扣除商品数量
			err = productService.SubNumberOne(message.ProductID)
			if err != nil {
				log.Println(err)
			}

			// 如果为true表示确认所有未确认的消息
			// 为false表示确认当前消息
			d.Ack(false)
		}
	}()

	log.Printf("按 CTRL+C 退出")
	<- forever
}