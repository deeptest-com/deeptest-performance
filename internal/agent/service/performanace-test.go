package agentService

import (
	"github.com/aaronchen2k/deeptest/proto"
	"io"
)

type PerformanceTestServices struct{}

func (services *PerformanceTestServices) Exec(stream proto.PerformanceService_ExecServer) error {
	i := 0

	for {
		err := stream.Send(&proto.PerformanceExecResp{
			Title:  "",
			Status: "pass",
			Result: "data",
		})
		if err != nil {
			return err
		}

		//
		res, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		err = stream.Send(&proto.PerformanceExecResp{
			Title:  res.Title,
			Status: "pass",
			Result: "data",
		})
		if err != nil {
			return err
		}

		i++
	}
}
