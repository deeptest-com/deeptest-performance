package service

import (
	"context"
	serverDomain "github.com/aaronchen2k/deeptest/cmd/server/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/pkg/mq"
	"github.com/aaronchen2k/deeptest/proto"
	"github.com/kataras/iris/v12"
	"google.golang.org/grpc"
	"io"
	"log"
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

	go mq.SubServerMsg(s.DealwithResult)

	err = stream.Send(&proto.PerformanceExecReq{
		Uuid:  req.Uuid,
		Title: req.Title,
		Vus:   int32(req.Vus),

		NsqServerAddress: req.NsqServerAddress,
		NsqLookupAddress: req.NsqLookupAddress,
	})
	if err != nil {
		return
	}

	for true {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		mqData := mq.MqMsg{
			Event:  "result",
			Result: *res,
		}
		mq.PubServerMsg(mqData)
	}

	stream.CloseSend()

	return
}

func (s *PerformanceTestServices) DealwithResult(mqMsg mq.MqMsg) (err error) {
	result := mqMsg.Result

	if result.Msg != "" {
		log.Printf("Msg: %s", result.Msg)
	} else {
		log.Printf("Result %s: %s", result.Uuid, result.Status)
	}

	return
}
