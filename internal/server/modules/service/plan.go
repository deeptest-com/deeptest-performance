package service

import (
	"context"
	"github.com/aaronchen2k/deeptest/proto"
	"github.com/kataras/iris/v12"
	"google.golang.org/grpc"
	"io"
	"log"
	"strconv"
)

type PlanService struct {
	StreamClient proto.StreamServiceClient
}

func (s *PlanService) Connect(ctx iris.Context) {
	connect, err := grpc.Dial("127.0.0.1:9528", grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}
	s.StreamClient = proto.NewStreamServiceClient(connect)
}

func (s *PlanService) OrderList(ctx iris.Context) {
	stream, err := s.StreamClient.OrderList(context.Background(), &proto.OrderSearchParams{})
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

func (s *PlanService) UploadImage(ctx iris.Context) {
	stream, err := s.StreamClient.UploadFile(context.Background())
	if err != nil {
		ctx.JSON(map[string]string{
			"err": err.Error(),
		})
		return
	}
	for i := 1; i <= 10; i++ {
		img := &proto.Image{FileName: "image" + strconv.Itoa(i), File: "file data"}
		images := &proto.StreamImageList{Image: img}
		err := stream.Send(images)
		if err != nil {
			ctx.JSON(map[string]string{
				"err": err.Error(),
			})
			return
		}
	}
	//发送完毕 关闭并获取服务端返回的消息
	resp, err := stream.CloseAndRecv()
	if err != nil {
		ctx.JSON(map[string]string{
			"err": err.Error(),
		})
		return
	}
	ctx.JSON(map[string]interface{}{"result": resp, "message": "success"})
	log.Println(resp)
}

func (s *PlanService) SumData(ctx iris.Context) {
	stream, err := s.StreamClient.SumData(context.Background())
	if err != nil {
		ctx.JSON(map[string]string{
			"err": err.Error(),
		})
		return
	}

	for i := 1; i <= 10; i++ {
		err = stream.Send(&proto.StreamSumData{Number: int32(i)})
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

		log.Printf("res number:%d", res.Number)
	}

	stream.CloseSend()

	return
}
