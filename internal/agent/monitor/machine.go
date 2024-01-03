package monitor

import (
	"context"
	"fmt"
	statUtils "github.com/aaronchen2k/deeptest/internal/pkg/utils/stat"
	"github.com/aaronchen2k/deeptest/proto"
	"log"
	"time"
)

func Monitor(stream *proto.PerformanceService_ExecServer, valCtx context.Context) {
	plan := valCtx.Value("plan").(*proto.PerformanceExecReq)
	log.Println(plan)

	log.Println(">>>>>> start monitor")

	for true {
		time.Sleep(3 * time.Second)

		data := statUtils.GetAll(stream)

		result := proto.PerformanceExecResult{
			Uuid: fmt.Sprintf("%s", plan.Uuid),
			Msg:  fmt.Sprintf("cpu usaged: %.2f%%", data.CpuUsage),
		}

		(*stream).Send(&result)

		// 每个场景处理器完成后，检测是否有终止执行的信号
		select {
		case <-valCtx.Done():
			fmt.Println("<<<<<<< stop monitor")

			// 中止执行该场景后续处理器
			return

		default:
		}
	}
}
