package service

import (
	"context"
	"encoding/json"
	"fmt"
	serverDomain "github.com/aaronchen2k/deeptest/cmd/server/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/aaronchen2k/deeptest/internal/pkg/queue"
	websocketHelper "github.com/aaronchen2k/deeptest/internal/pkg/websocket"
	"github.com/aaronchen2k/deeptest/proto"
	"github.com/kataras/iris/v12/websocket"
	"github.com/nsqio/go-nsq"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

type PerformanceTestServices struct {
	execCtx    context.Context
	execCancel context.CancelFunc

	PerformanceServiceClient proto.PerformanceServiceClient
}

func (s *PerformanceTestServices) Connect() {
	connect, err := grpc.Dial("127.0.0.1:9528", grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}
	s.PerformanceServiceClient = proto.NewPerformanceServiceClient(connect)
}

func (s *PerformanceTestServices) ExecStart(req serverDomain.PlanExecReq, wsMsg *websocket.Message) (err error) {
	if s.PerformanceServiceClient == nil {
		s.Connect()
	}

	s.execCtx, s.execCancel = context.WithCancel(context.Background())

	// stop execution in 2 ways:
	// 1. call cancel in this method by websocket client
	// 2. sub cancel instruction from agent via grpc

	if req.NsqServerAddress != "" { // agent send logs via nsq MQ
		// check ctx.Done
		go s.HandleAgentNsqMsg(req, s.execCtx, req.Uuid, wsMsg)

	} else { // agent send logs via grpc, server store msgs in queue
		// check ctx.Done
		// may cancel ctx by instruction from agent
		go queue.SubAgentGrpcMsg(s.DealwithResult, s.execCtx, s.execCancel, req.Uuid, wsMsg)
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

func (s *PerformanceTestServices) ExecStop(req serverDomain.PlanExecReq, wsMsg *websocket.Message) (err error) {
	if s.execCancel != nil {
		s.execCancel()
	}

	websocketHelper.SendExecInstructionToClient("", "", consts.MsgInstructionEnd, req.Uuid, wsMsg)

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

		queue.PubGrpcMsg(*res)
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

func (s *PerformanceTestServices) HandleAgentNsqMsg(req serverDomain.PlanExecReq, ctx context.Context,
	uuid string, wsMsg *websocket.Message) (err error) {

	channel := fmt.Sprintf("channel_%s", req.Uuid)
	consumer, err := nsq.NewConsumer(req.Uuid, channel, nsq.NewConfig())
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer consumer.Stop()

	consumer.AddHandler(newNsqMsgProcessor(s.nsqMsgCallback, req.Uuid, wsMsg))

	nsqAddr := req.NsqServerAddress
	if req.NsqLookupAddress != "" {
		nsqAddr = req.NsqLookupAddress
	}
	err = consumer.ConnectToNSQD(nsqAddr)
	if err != nil {
		return
	}

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

func (s *PerformanceTestServices) nsqMsgCallback(bytes []byte, execUUid string, wsMsg *websocket.Message) error {
	log.Println(fmt.Sprintf("receive msg: %s", bytes))

	result := proto.PerformanceExecResult{}
	json.Unmarshal(bytes, &result)

	result.ExecUUid = execUUid

	s.DealwithResult(result, execUUid, wsMsg)

	return nil
}

func (s *PerformanceTestServices) DealwithResult(result proto.PerformanceExecResult, execUUid string, wsMsg *websocket.Message) (err error) {
	if result.Instruction != "" { // only dealwith results msg, agent will send Instruction via grpc
		return
	}

	if result.Record != nil {
		if result.Record.Msg != "" {
			log.Printf("Msg: %s", result.Record.Msg)
		} else {
			log.Printf("Result %d: %s", result.Record.RecordId, result.Record.Status)
		}
	}

	if wsMsg != nil {
		websocketHelper.SendExecResultToClient(result, consts.MsgResultRecord, execUUid, wsMsg)
	}

	return
}

type NsqMsgProcessor struct {
	callback func(msg []byte, execUUid string, wsMsg *websocket.Message) error
	cancel   context.CancelFunc
	execUuid string
	wsMsg    *websocket.Message
}

func newNsqMsgProcessor(callback func(msg []byte, execUUid string, wsMsg *websocket.Message) error, execUUid string, wsMsg *websocket.Message) *NsqMsgProcessor {
	return &NsqMsgProcessor{
		execUuid: execUUid,
		callback: callback,
		wsMsg:    wsMsg,
	}
}

func (m *NsqMsgProcessor) HandleMessage(msg *nsq.Message) (err error) {
	err = m.callback(msg.Body, m.execUuid, m.wsMsg)
	if err != nil {
		return
	}

	msg.Finish()

	return nil
}
