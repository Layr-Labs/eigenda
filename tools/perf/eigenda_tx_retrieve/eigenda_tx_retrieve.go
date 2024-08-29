package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/Layr-Labs/eigenda/api/clients"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"external/common"
)

func RetrieveData(disp clients.DisperserClient, requestId string) (string, int64, error) {
	key1, err := base64.StdEncoding.DecodeString(requestId)
	if err != nil {
		return "", 0, err
	}

	ctxTimeout, cancel := context.WithTimeout(context.Background(), common.ConfigInstance.RetrieveTimeout)
	defer cancel()

	ret, err := disp.GetBlobStatus(ctxTimeout, key1)
	if err != nil {
		return "", 0, err
	}

	blockTime := uint64(0)
	if ret.Status.String() == "CONFIRMED" || ret.Status.String() == "FINALIZED" {
		confirmBlockNumber := ret.Info.BlobVerificationProof.BatchMetadata.ConfirmationBlockNumber
		client, err := ethclient.Dial(common.ConfigInstance.L1URL)
		if err != nil {
			return "", 0, err
		}

		blockNumber := big.NewInt(int64(confirmBlockNumber))
		block, err := client.BlockByNumber(context.Background(), blockNumber)
		if err != nil {
			return "", 0, err
		}
		blockTime = block.Time()
	}

	return ret.Status.String(), int64(blockTime), nil
}

func WriteFile(line string, fileName string) {
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("打开文件时出错:", err)
		return
	}
	defer file.Close()

	if _, err := file.WriteString(line); err != nil {
		log.Println("写入文件时出错:", err)
		return
	}
}

func main() {
	file, err := os.OpenFile("retrieve.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("无法打开日志文件: %v", err)
	}
	defer file.Close()
	//log.SetOutput(file)

	ch, q := common.ConnectMQ()

	// 订阅队列消息
	msgs, err := ch.Consume(
		q.Name, // 队列名称
		"",     // 消费者名称（空字符串表示由服务器生成唯一名称）
		false,  // 自动应答
		false,  // 非独占
		false,  // 非本地
		false,  // 无等待
		nil,    // 额外属性
	)
	common.FailOnError(err, "Failed to register a consumer")

	// 启动一个协程，处理消息
	disp := common.NewClient()
	totalMsg := int64(0)
	failMsg := int64(0)
	successMsg := int64(0)
	processingMsg := int64(0)

	{
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			sig := <-sigs
			log.Println("receive signal:", sig)
			log.Println("toalMsg:", totalMsg, ", successMsg:", successMsg, ", processingMsg:", processingMsg, "failMsg: ", failMsg)
			os.Exit(0)
		}()
	}

	fileName := "eigenda" + time.Now().Format("20060102150405") + ".txt"
	tmp := 0
	for d := range msgs {
		tmp += 1
		fmt.Println("tmp:", tmp)

		var body common.SendMsg
		err := json.Unmarshal(d.Body, &body)
		if err != nil {
			log.Printf("Failed to decode JSON: %s", err.Error())
		} else {
			fmt.Printf("receive msg, %s\n", string(d.Body))
			var printMsg = common.PrintMsg{
				SendRequestTime:    body.SendRequestTime,
				ReceiveRequestTime: body.ReceiveRequestTime,
				RequestId:          body.RequestId,
			}

			if body.RequestId == "" {
				d.Ack(false)
			} else {

				status, confirmTime, err := RetrieveData(disp, body.RequestId)
				if err != nil {
					fmt.Println("retrieve fail, err:", err.Error())
					continue
				}
				if status == "CONFIRMED" || status == "FINALIZED" {
					atomic.AddInt64(&successMsg, 1)
					atomic.AddInt64(&totalMsg, 1)
					d.Ack(false)
					printMsg.ConfirmTime = confirmTime
				} else if status == "FAILED" {
					atomic.AddInt64(&failMsg, 1)
					atomic.AddInt64(&totalMsg, 1)
					d.Ack(false)
					log.Println("status: FAILED")
					continue
				} else if status == "PROCESSING" {
					atomic.AddInt64(&processingMsg, 1)
					atomic.AddInt64(&totalMsg, 1)
					d.Ack(false)
					continue
				} else {
					d.Ack(false)
					continue
				}
			}
			printMsgJson, _ := json.Marshal(printMsg)
			WriteFile(string(printMsgJson)+"\n", fileName)
		}
	}
}
