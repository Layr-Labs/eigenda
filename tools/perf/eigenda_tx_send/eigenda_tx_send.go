package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"external/common"
	"fmt"
	"github.com/Layr-Labs/eigenda/api/clients"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/streadway/amqp"
	"log"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

const clientTimeout = 10 * time.Second
const sendTimeout = 10 * time.Second
const retrieveTimeout = 10 * time.Second

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func SendData(disp clients.DisperserClient) (string, error) {
	data := make([]byte, 190*1024*1024/100)
	_, err := rand.Read(data)
	if err != nil {
		return "", err
	}

	paddedData := codec.ConvertByPaddingEmptyByte(data)
	ctxTimeout, cancel := context.WithTimeout(context.Background(), sendTimeout)
	defer cancel()

	blobStatus1, key1, err := disp.DisperseBlobAuthenticated(ctxTimeout, paddedData, []uint8{})
	if err != nil {
		return "", err
	}
	if blobStatus1 == nil {
		return "", errors.New("blob status == nil")
	}
	if key1 == nil {
		return "", errors.New("key == nil")
	}
	encoded := base64.StdEncoding.EncodeToString(key1)
	return encoded, nil
}

func main() {
	{
		disp := common.NewClient()
		requestId, err := SendData(disp)
		if err != nil {
			panic(err)
		}
		fmt.Println("requestId:", requestId)
		return
	}
	file, err := os.OpenFile("send.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("无法打开日志文件: %v", err)
	}
	defer file.Close()
	log.SetOutput(file)

	// 连接到 RabbitMQ 服务器
	ch, q := common.ConnectMQ()
	defer ch.Close()

	sendInterval := int64(1 * 1000)
	sendBlobNumber := 30

	sendMutex := sync.Mutex{}

	sendInterval = sendInterval * 1 * int64(time.Millisecond)
	sendTicker := time.NewTicker(time.Duration(sendInterval))

	ctx, cancel := context.WithCancel(context.Background())
	disp := common.NewClient()
	defer cancel()

	beginTime := time.Now()
	totalSendNumber := int64(0)
	totalSuccessSendNumber := int64(0)

	lastRequestId := ""
	done := make(chan bool, 1)
	{
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			sig := <-sigs
			log.Println("receive signal:", sig)
			log.Println("beginTime:", beginTime, ", endTime:", time.Now(), ", totalSendNumber:", totalSendNumber,
				", totalSuccessSendNumber:", totalSuccessSendNumber, ", lastRequestId:", lastRequestId)
			// 收到信号后，向done通道发送true，通知程序退出
			done <- true
		}()

	}

	log.Println("begin loop")
	for {
		select {
		case <-ctx.Done():
			break
		case <-sendTicker.C:
			go func() {

				if !sendMutex.TryLock() {
					log.Println("send ongoing, time:", time.Now())
					return
				}

				var wg sync.WaitGroup
				successSendNumber := int64(0)

				for i := 0; i < sendBlobNumber; i++ {
					wg.Add(1)

					go func() {
						defer wg.Done()
						atomic.AddInt64(&totalSendNumber, 1)

						body := common.SendMsg{
							SendRequestTime: time.Now().Unix(),
						}

						requestId, err := SendData(disp)
						if err != nil {
							log.Println("send failed:", err)
						} else {

							lastRequestId = requestId

							fmt.Println("debug00000000, requestId:", requestId)
							body.RequestId = requestId
							body.ReceiveRequestTime = time.Now().Unix()

							atomic.AddInt64(&successSendNumber, 1)
							atomic.AddInt64(&totalSuccessSendNumber, 1)
						}

						bodyBytes, _ := json.Marshal(body)

						err = ch.Publish(
							"",
							q.Name,
							false,
							false,
							amqp.Publishing{
								ContentType: "application/json",
								Body:        bodyBytes,
							})
						failOnError(err, "Failed to publish a message")

					}()
				}
				wg.Wait()
				log.Println("success send :", successSendNumber)

				sendMutex.Unlock()
			}()
		case <-done:
			log.Println("program exit")
			os.Exit(0)
		}
	}
}
