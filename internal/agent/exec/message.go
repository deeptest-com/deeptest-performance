package exec

import (
	"encoding/json"
	"github.com/aaronchen2k/deeptest/proto"
	"github.com/nsqio/go-nsq"
)

type MessageSender interface {
	Send(result proto.PerformanceExecResult) error
}

type NsqSender struct {
	NsqServerAddress string
	NsqLookupAddress string

	Topic    string
	producer *nsq.Producer
}

type GrpcSender struct {
	Stream *proto.PerformanceService_ExecServer
}

func NewNsqSender(topic string, nsqServerAddress, nsqLookupAddress string) MessageSender {
	ret := NsqSender{
		NsqServerAddress: nsqServerAddress,
		NsqLookupAddress: nsqLookupAddress,

		Topic: topic,
	}

	var err error

	if ret.producer == nil {
		ret.producer, err = nsq.NewProducer(nsqServerAddress, nsq.NewConfig())
		if err != nil {
			return nil
		}
	}

	return ret
}

func NewGrpcSender(stream *proto.PerformanceService_ExecServer) MessageSender {
	ret := GrpcSender{
		Stream: stream,
	}

	return ret
}

func (s NsqSender) Send(result proto.PerformanceExecResult) (err error) {
	bytes, err := json.Marshal(result)
	if err != nil {
		return err
	}

	err = s.producer.Publish(s.Topic, bytes)
	if err != nil {
		return err
	}

	return
}

func (s GrpcSender) Send(result proto.PerformanceExecResult) (err error) {
	//mqData := mq.MqMsg{
	//	Event:  "result",
	//	Result: result,
	//}
	//mq.PubAgentMsg(mqData)

	(*s.Stream).Send(&result)

	return
}
