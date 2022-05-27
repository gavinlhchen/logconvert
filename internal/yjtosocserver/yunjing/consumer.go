package yunjing

import (
	"encoding/binary"
	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
	"logconvert/errors"
	"logconvert/internal/yjtosocserver/protos/ydyes"
	"net"
)

func recurse(m map[string]interface{}) {
	for field, val := range m {
		if v, ok := val.(map[string]interface{}); ok {
			recurse(v)
		} else if v1, ok1 := val.([]byte); ok1 {
			m[field] = string(v1)
		} else if v2, ok2 := val.([]interface{}); ok2 {
			for _, sub := range v2 {
				if v3, ok3 := sub.(map[string]interface{}); ok3 {
					recurse(v3)
				}
			}
		}
	}
}

type YjToSocMsgHandler interface {
	Handle(messageValue []byte, producer sarama.AsyncProducer, topic string)
}

type ConsumerGroupHandler struct {
	Ready        chan bool
	Producer     sarama.AsyncProducer
	ProduceTopic string
	MsgHandler   YjToSocMsgHandler
}

func (consumer *ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	close(consumer.Ready)
	return nil
}

func (consumer *ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		consumer.MsgHandler.Handle(message.Value, consumer.Producer, consumer.ProduceTopic)
		session.MarkMessage(message, "")
	}
	return nil
}

const (
	ProtoAdditionLen uint32 = 1 + 4 + 4 + 1
)

func DecryptProto(protoBuf []byte) (head *ydyes.Head, body []byte, err error) {
	if protoBuf == nil || uint32(len(protoBuf)) < ProtoAdditionLen {
		return nil, nil, errors.New("invalid proto length")
	}

	offset := uint32(1)
	lenSize := uint32(4)
	headLen := binary.BigEndian.Uint32(protoBuf[offset : offset+lenSize])

	offset += lenSize
	bodyLen := binary.BigEndian.Uint32(protoBuf[offset : offset+lenSize])

	protoLen := uint32(len(protoBuf))
	if protoLen != (headLen + bodyLen + ProtoAdditionLen) {
		return nil, nil, errors.New("invalid proto length")
	}

	if protoBuf[0] != uint8(ydyes.MsgCmd_TCP_STX_C) ||
		protoBuf[protoLen-1] != uint8(ydyes.MsgCmd_TCP_ETX_C) {
		return nil, nil, errors.New("invalid proto. check mask fail")
	}

	offset += lenSize

	head = new(ydyes.Head)
	err = proto.Unmarshal(protoBuf[offset:headLen+offset], head)
	if err != nil {
		return nil, nil, err
	}

	var hostIp string
	if string(head.GetHostip()) != "" {
		hostIp = string(head.GetHostip())
	} else if head.GetUint32Localip() != 0 {
		hostIp = IpUintToString(head.GetUint32Localip())
	}
	head.Hostip = []byte(hostIp)

	body = protoBuf[headLen+9 : headLen+9+bodyLen]
	return head, body, nil
}

func IpUintToString(ip uint32) string {
	var bytes [4]byte
	bytes[0] = byte(ip & 0xFF)
	bytes[1] = byte((ip >> 8) & 0xFF)
	bytes[2] = byte((ip >> 16) & 0xFF)
	bytes[3] = byte((ip >> 24) & 0xFF)
	return net.IPv4(bytes[3], bytes[2], bytes[1], bytes[0]).String()
}
