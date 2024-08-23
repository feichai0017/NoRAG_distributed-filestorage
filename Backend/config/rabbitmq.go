package config

const (
	//AsyncTransferEnable : 是否开启文件异步转移(默认同步)
	AsyncTransferEnable = true
	//RabbitURL : rabbitmq服务的入口url
	RabbitURL = "amqp://guest:guest@localhost:5672/"
	//TransExchangeName : 用于文件transfer的交换机
	TransExchangeName = "uploadserver.trans"
	//TransS3QueueName : s3转移队列名
	TransS3QueueName = "uploadserver.trans.s3"
	//TransS3ErrQueueName : s3转移失败后写入另一个队列的队列名
	TransS3ErrQueueName = "uploadserver.trans.s3.err"
	//TransS3RoutingKey: s3转移队列路由名
	TransS3RoutingKey = "s3"
)
