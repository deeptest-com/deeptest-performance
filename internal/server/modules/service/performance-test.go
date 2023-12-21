package service

import (
	"context"
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

func (s *PerformanceTestServices) Exec(ctx iris.Context) {
	stream, err := s.PerformanceServiceClient.Exec(context.Background())
	if err != nil {
		ctx.JSON(map[string]string{
			"err": err.Error(),
		})
		return
	}

	go mq.SubMsg(s.DealwithResult)

	for i := 1; i <= 10; i++ {
		err = stream.Send(&proto.PerformanceExecReq{
			Uuid:  "UUID-123",
			Title: "Performance Testing Task UUID-123",
			Vus:   10,
		})
		if err == io.EOF {
			break
		}
		if err != nil {
			return
		}

		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return
		}

		mqData := mq.MqMsg{
			Event:  "result",
			Result: *res,
		}
		mq.PubMsg(mqData)
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
