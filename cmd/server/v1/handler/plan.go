package handler

import (
	"context"
	"github.com/aaronchen2k/deeptest/internal/server/modules/service"
	_domain "github.com/aaronchen2k/deeptest/pkg/domain"
	"github.com/aaronchen2k/deeptest/proto"
	"github.com/kataras/iris/v12"
	"google.golang.org/grpc"
	"io"
	"log"
)

type PlanCtrl struct {
	streamClient proto.StreamServiceClient

	PlanService *service.PlanService `inject:""`
}

func (c *PlanCtrl) Get(ctx iris.Context) {
	c.requestGrpc()

	data := iris.Map{}
	ctx.JSON(_domain.Response{Code: _domain.NoErr.Code, Data: data, Msg: _domain.NoErr.Msg})
}

func (c *PlanCtrl) requestGrpc() {
	connect, err := grpc.Dial("127.0.0.1:9528", grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}
	c.streamClient = proto.NewStreamServiceClient(connect)

	stream, err := c.streamClient.OrderList(context.Background(), &proto.OrderSearchParams{})
	if err != nil {
		return
	}

	for {
		orderList, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return
		}

		log.Println(orderList.Order)
	}

}
