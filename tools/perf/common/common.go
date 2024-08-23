package common

import (
	"github.com/Layr-Labs/eigenda/api/clients"
	"github.com/Layr-Labs/eigenda/core/auth"
	"github.com/streadway/amqp"
	"log"
)

var ConfigInstance = DevConfig

type SendMsg struct {
	SendRequestTime    int64  `json:"send_request_time"`
	ReceiveRequestTime int64  `json:"receive_request_time"`
	RequestId          string `json:"request_id"`
}

type PrintMsg struct {
	SendRequestTime    int64  `json:"send_request_time"`
	ReceiveRequestTime int64  `json:"receive_request_time"`
	RequestId          string `json:"request_id"`
	ConfirmTime        int64  `json:"confirm_time"`
}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func ConnectMQ() (*amqp.Channel, amqp.Queue) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	FailOnError(err, "Failed to connect to RabbitMQ")

	// 创建一个 channel
	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")

	// 声明一个队列
	q, err := ch.QueueDeclare(
		ConfigInstance.QueueName, // 队列名称
		false,                    // 消息是否持久化
		false,                    // 是否自动删除
		false,                    // 是否排他
		false,                    // 是否阻塞
		nil,                      // 其他参数
	)
	FailOnError(err, "Failed to declare a queue")
	return ch, q
}

func NewClient() clients.DisperserClient {

	privateKeyHex := "0x554337ed11c6f083f87c0459b798e684b36e2aeb2b9d386a01925e6118334600"
	signer := auth.NewLocalBlobRequestSigner(privateKeyHex)

	disp := clients.NewDisperserClient(&clients.Config{

		Hostname:          ConfigInstance.Hostname,
		Port:              ConfigInstance.Port,
		UseSecureGrpcFlag: ConfigInstance.UseSecureGrpcFlag,
		Timeout:           ConfigInstance.Timeout,
	}, signer)

	return disp
}
