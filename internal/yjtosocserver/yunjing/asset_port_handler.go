package yunjing

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/fatih/structs"
	"github.com/golang/protobuf/proto"
	uuid "github.com/satori/go.uuid"
	"logconvert/internal/pkg/soc"
	"logconvert/internal/yjtosocserver/protos/ydyes"
	"logconvert/log"
	"time"
)

func init() {
	Register("asset_port", &AssetPortConsumer{})

}

type AssetPortConsumer struct {
}

func (c *AssetPortConsumer) Handle(messageValue []byte, producer sarama.AsyncProducer, topic string) {
	head, body, err := DecryptProto(messageValue)
	if err != nil {
		log.Errorf("Error Decrypt: %v", err)
	}

	headMap := structs.Map(head)
	recurse(headMap)

	processInfoMsg := &ydyes.OpenPortMsg{}
	if err := proto.Unmarshal(body, processInfoMsg); err != nil {
		log.Errorf("Error Unmarshal: %v", err)
	}

	for _, msg := range processInfoMsg.PortInfoLst {
		bodyMap := structs.Map(msg)
		recurse(bodyMap)
		bodyMap["BytesGuid"] = string(processInfoMsg.BytesGuid)
		bodyMap["BytesQuuid"] = string(processInfoMsg.BytesQuuid)
		bodyMap["PushResult"] = string(processInfoMsg.PushResult)
		bodyMap["BytesUuid"] = string(processInfoMsg.BytesUuid)
		bodyMap["Platform"] = processInfoMsg.ReportType

		allMap := map[string]interface{}{"Head": headMap, "OpenPortMsg": bodyMap}

		if evStr, err := json.Marshal(allMap); err != nil {
			log.Errorf("marshal error:%v", err)
		} else {
			rawEvent := &soc.RawEvent{
				LogsourceIp:        "127.0.0.1",
				LogsourceName:      "inner-probe-xdr_yunjing",
				LogsourceTimestamp: time.Now().UnixNano() / 1e6,
				LogsourceCategory:  "event",
				RawLogCharset:      "utf-8",
				RawLog:             string(evStr),
				EventUuid:          uuid.Must(uuid.NewV4()).String(),
			}

			rawByte, _ := json.Marshal(rawEvent)
			producer.Input() <- &sarama.ProducerMessage{
				Topic: topic,
				Value: sarama.ByteEncoder(rawByte),
			}
		}
	}
}
