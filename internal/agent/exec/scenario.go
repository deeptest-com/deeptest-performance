package exec

import (
	"context"
	"fmt"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	_intUtils "github.com/aaronchen2k/deeptest/pkg/lib/int"
	"github.com/aaronchen2k/deeptest/proto"
	"log"
	"time"
)

func ExecScenario(valCtx context.Context, stream *proto.PerformanceService_ExecServer, sender MessageSender) {
	startTime := time.Now()

	task := valCtx.Value("task").(domain.Task)

	for _, processor := range task.Scenario[0].Processors {
		log.Println("exec processor", processor)

		// 此处为场景处理器的耗时操作
		time.Sleep(time.Duration(_intUtils.GenRandNum(1, 10)) * time.Second)

		result := proto.PerformanceExecResult{
			Uuid:   fmt.Sprintf("%s@%s", processor.Name, task.Uuid),
			Status: "pass",
		}

		sender.Send(result)

		// 每个场景处理器完成后，检测是否有终止执行的信号
		select {
		case <-valCtx.Done():
			fmt.Println("exit left scenario processors by signal", task.VuNo)

			// 中止执行该场景后续处理器
			goto Label_END_SCENARIO

		default:
		}
	}

	fmt.Println("complete scenario normally", task.VuNo)

Label_END_SCENARIO:

	endTime := time.Now()

	result := proto.PerformanceExecResult{
		Uuid:     fmt.Sprintf("scenario_%s", task.Uuid),
		Status:   "pass",
		Duration: endTime.Unix() - startTime.Unix(),
	}
	sender.Send(result)

	return
}
