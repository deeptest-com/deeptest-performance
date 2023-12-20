package service

import (
	"context"
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

	for i := 1; i <= 10; i++ {
		err = stream.Send(&proto.PerformanceExecReq{
			ExecUuid: "UUID-123",
			Title:    "任务 UUID-123",
			Vus:      10,
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

		log.Printf("%s status: %s", res.Title, res.Status)
	}

	stream.CloseSend()

	return
}
