package exec

import (
	"context"
	"fmt"
	"github.com/aaronchen2k/deeptest/internal/agent/logs"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	_httpUtils "github.com/aaronchen2k/deeptest/pkg/lib/http"
	_intUtils "github.com/aaronchen2k/deeptest/pkg/lib/int"
	"github.com/aaronchen2k/deeptest/proto"
	"log"
	"time"
)

func ExecScenario(valCtx context.Context, stream *proto.PerformanceService_ExecServer, sender logs.MessageSender) {
	startTime := time.Now()

	task := valCtx.Value("task").(domain.Task)

	for index, processor := range task.Scenario[0].Processors {
		log.Println("exec processor", processor)

		{
			bytes, err := _httpUtils.Get("http://111.231.16.35:9000/get")
			log.Println(bytes, err)
		}

		duration := _intUtils.GenRandNum(100, 1000)
		time.Sleep(time.Duration(duration) * time.Millisecond)

		status := "pass"
		if index%3 == 0 {
			status = "fail"
		}

		result := proto.PerformanceExecResult{
			Uuid:     fmt.Sprintf("%s@%s", processor.Name, task.Uuid),
			Duration: int64(duration),
			Status:   status,
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
