package service

import (
	"context"
	"encoding/json"
	"fmt"
	serverDomain "github.com/aaronchen2k/deeptest/cmd/server/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/aaronchen2k/deeptest/internal/pkg/queue"
	"github.com/aaronchen2k/deeptest/proto"
	"github.com/kataras/iris/v12/websocket"
	"github.com/nsqio/go-nsq"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

type PerformanceTestServices struct {
	PerformanceServiceClient proto.PerformanceServiceClient
}

func (s *PerformanceTestServices) Connect() {
	connect, err := grpc.Dial("127.0.0.1:9528", grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}
	s.PerformanceServiceClient = proto.NewPerformanceServiceClient(connect)
}

func (s *PerformanceTestServices) Exec(req serverDomain.PlanExecReq, wsMsg *websocket.Message) (err error) {
	if s.PerformanceServiceClient == nil {
		s.Connect()
	}

	ctx, cancel := context.WithCancel(context.Background())

	if req.NsqServerAddress == "" {
		go queue.SubMsg(s.DealwithResult, cancel)
	} else {
		go s.HandleNsqMsg(req, ctx)
	}

	// send exec request to agent
	stream, err := s.SendExecReqToAgent(req)
	if err != nil {
		return
	}

	s.HandleAndForwardGrpcMsg(stream)

	stream.CloseSend()

	return
}

func (s *PerformanceTestServices) HandleAndForwardGrpcMsg(stream proto.PerformanceService_ExecClient) (err error) {
	for true {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		queue.PubMsg(*res)
	}

	return
}

func (s *PerformanceTestServices) SendExecReqToAgent(req serverDomain.PlanExecReq) (
	stream proto.PerformanceService_ExecClient, err error) {

	stream, err = s.PerformanceServiceClient.Exec(context.Background())
	if err != nil {
		return
	}

	err = stream.Send(&proto.PerformanceExecReq{
		Uuid:  req.Uuid,
		Title: req.Title,

		GenerateType: req.GenerateType.String(),
		Stages:       req.Stages,

		NsqServerAddress: req.NsqServerAddress,
		NsqLookupAddress: req.NsqLookupAddress,

		Scenarios: req.Scenarios,
	})

	return
}

func (s *PerformanceTestServices) HandleNsqMsg(req serverDomain.PlanExecReq, ctx context.Context) (err error) {
	channel := fmt.Sprintf("channel_%s", req.Uuid)
	consumer, err := nsq.NewConsumer(req.Uuid, channel, nsq.NewConfig())
	if err != nil {
		return
	}
	defer consumer.Stop()

	consumer.AddHandler(newNsqMsgProcessor(s.NsqMsgCallback))

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

	return nil
}

func (s *PerformanceTestServices) NsqMsgCallback(bytes []byte) error {
	log.Println(fmt.Sprintf("receive msg: %s", bytes))

	result := proto.PerformanceExecResult{}
	json.Unmarshal(bytes, &result)

	s.DealwithResult(result)

	return nil
}

func (s *PerformanceTestServices) DealwithResult(result proto.PerformanceExecResult) (err error) {
	if result.Instruction != consts.Result.String() {
		return
	}

	if result.Record != nil {
		if result.Record.Msg != "" {
			log.Printf("Msg: %s", result.Record.Msg)
		} else {
			log.Printf("Result %s: %s", result.Record.Uuid, result.Record.Status)
		}
	}

	return
}

type NsqMsgProcessor struct {
	callback func(msg []byte) error
	cancel   context.CancelFunc
}

func newNsqMsgProcessor(callback func(msg []byte) error) *NsqMsgProcessor {
	return &NsqMsgProcessor{
		callback: callback,
	}
}

func (m *NsqMsgProcessor) HandleMessage(msg *nsq.Message) (err error) {
	err = m.callback(msg.Body)
	if err != nil {
		return
	}

	msg.Finish()

	return nil
}
