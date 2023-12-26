package service

import (
	"context"
	"encoding/json"
	"fmt"
	serverDomain "github.com/aaronchen2k/deeptest/cmd/server/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/pkg/queue"
	"github.com/aaronchen2k/deeptest/proto"
	"github.com/kataras/iris/v12"
	"github.com/nsqio/go-nsq"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

type PerformanceTestServices struct {
	PerformanceServiceClient proto.PerformanceServiceClient
}

func (s *PerformanceTestServices) Connect(ctx iris.Context) {
	connect, err := grpc.Dial("127.0.0.1:9528", grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}
	s.PerformanceServiceClient = proto.NewPerformanceServiceClient(connect)
}

func (s *PerformanceTestServices) Exec(req serverDomain.PlanExecReq) (err error) {
	stream, err := s.PerformanceServiceClient.Exec(context.Background())
	if err != nil {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	go queue.SubServerMsg(s.DealwithResult, cancel)

	err = stream.Send(&proto.PerformanceExecReq{
		Uuid:  req.Uuid,
		Title: req.Title,
		Vus:   int32(req.Vus),

		NsqServerAddress: req.NsqServerAddress,
		NsqLookupAddress: req.NsqLookupAddress,

		Scenarios: req.Scenarios,
	})
	if err != nil {
		return
	}

	// use nsq to get results
	if req.NsqServerAddress != "" {
		go func() {
			channel := fmt.Sprintf("channel_%s", req.Uuid)
			consumer, err := nsq.NewConsumer(req.Uuid, channel, nsq.NewConfig())
			if err != nil {
				return
			}
			defer consumer.Stop()

			msgCallback := func(bytes []byte) error {
				log.Println(fmt.Sprintf("receive msg: %s", bytes))

				result := proto.PerformanceExecResult{}
				json.Unmarshal(bytes, &result)

				s.DealwithResult(result)

				return nil
			}
			consumer.AddHandler(newMsgProcessor(msgCallback))

			nsqAddr := req.NsqServerAddress
			if req.NsqLookupAddress != "" {
				nsqAddr = req.NsqLookupAddress
			}
			err = consumer.ConnectToNSQD(nsqAddr)
			if err != nil {
				return
			}

			// wait util getting stop instruction from mq
			for true {
				select {
				case <-ctx.Done():
					return

				default:
					time.Sleep(3 * time.Second)
				}
			}
		}()
	}

	for true {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		queue.PubServerMsg(*res)
	}

	stream.CloseSend()

	return
}

func (s *PerformanceTestServices) DealwithResult(result proto.PerformanceExecResult) (err error) {
	if result.Msg != "" {
		log.Printf("Msg: %s", result.Msg)
	} else {
		log.Printf("Result %s: %s", result.Uuid, result.Status)
	}

	return
}

type msgProcessor struct {
	callback func(msg []byte) error
	cancel   context.CancelFunc
}

func newMsgProcessor(callback func(msg []byte) error) *msgProcessor {
	return &msgProcessor{
		callback: callback,
	}
}

func (m *msgProcessor) HandleMessage(msg *nsq.Message) (err error) {
	err = m.callback(msg.Body)
	if err != nil {
		return
	}

	msg.Finish()

	return nil
}
