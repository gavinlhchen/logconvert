package yunjing

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/fatih/structs"
	"github.com/gavinlhchen/logconvert/internal/pkg/soc"
	"github.com/gavinlhchen/logconvert/internal/yjtosocserver/protos/ydkafka"
	"github.com/gavinlhchen/logconvert/log"
	"github.com/golang/protobuf/proto"
	uuid "github.com/satori/go.uuid"
	"time"
)

func init() {
	Register("event_msg", &EventMsgConsumer{})

}

type EventMsgConsumer struct {
}

func (c *EventMsgConsumer) Handle(messageValue []byte, producer sarama.AsyncProducer, topic string) {
	evt := &ydkafka.EventMsg{}
	if err := proto.Unmarshal(messageValue, evt); err != nil {
		log.Errorf("Unmarshal err:%v", err)
	}
	headMap := structs.Map(evt.Head)
	recurse(headMap)

	var bodyMap map[string]interface{}
	var bodyName string
	if evt.Body.Malware != nil {
		bodyMap = structs.Map(evt.Body.Malware)
		bodyName = "Malware"
	} else if evt.Body.Login != nil {
		bodyMap = structs.Map(evt.Body.Login)
		bodyName = "Login"
	} else if evt.Body.Bruteforce != nil {
		bodyMap = structs.Map(evt.Body.Bruteforce)
		bodyName = "Bruteforce"
	} else if evt.Body.Vul != nil {
		bodyMap = structs.Map(evt.Body.Vul)
		bodyName = "Vul"
	} else if evt.Body.Bash != nil {
		bodyMap = structs.Map(evt.Body.Bash)
		bodyName = "Bash"
	} else if evt.Body.ReverseShell != nil {
		bodyMap = structs.Map(evt.Body.ReverseShell)
		bodyName = "ReverseShell"
	} else if evt.Body.PrivilegeEscalation != nil {
		bodyMap = structs.Map(evt.Body.PrivilegeEscalation)
		bodyName = "PrivilegeEscalation"
	} else if evt.Body.NetworkAttack != nil {
		bodyMap = structs.Map(evt.Body.NetworkAttack)
		bodyName = "NetworkAttack"
	} else if evt.Body.RiskDns != nil {
		bodyMap = structs.Map(evt.Body.RiskDns)
		bodyName = "RiskDns"
	} else if evt.Body.BaseLine != nil {
		bodyMap = structs.Map(evt.Body.BaseLine)
		bodyName = "BaseLine"
	}
	if nil == bodyMap {
		return
	}
	recurse(bodyMap)
	allMap := map[string]interface{}{"Head": headMap, bodyName: bodyMap}
	if evBytes, err := json.Marshal(allMap); err != nil {
		log.Errorf("marshal error:%v", err)
	} else {
		eveStr := string(evBytes)
		rawEvent := &soc.RawEvent{
			LogsourceIp:        "127.0.0.1",
			LogsourceName:      "yunjing_event_msg",
			LogsourceTimestamp: time.Now().UnixNano() / 1e6,
			LogsourceCategory:  "event",
			RawLogCharset:      "utf-8",
			RawLog:             eveStr,
			EventUuid:          uuid.Must(uuid.NewV4()).String(),
		}

		rawByte, _ := json.Marshal(rawEvent)
		producer.Input() <- &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.ByteEncoder(rawByte),
		}
	}
}
