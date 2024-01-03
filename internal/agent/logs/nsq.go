package logs

import (
	"encoding/json"
	"github.com/aaronchen2k/deeptest/proto"
	"github.com/nsqio/go-nsq"
)

type NsqSender struct {
	NsqServerAddress string
	NsqLookupAddress string

	Topic    string
	producer *nsq.Producer
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

func (s NsqSender) Send(result proto.PerformanceExecResult) (err error) {
	Count(result)

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
