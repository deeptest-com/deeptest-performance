package exec

import (
	"context"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	"github.com/aaronchen2k/deeptest/proto"
	"log"
)

func ExecTask(valCtx context.Context, stream *proto.PerformanceService_ExecServer) {
	task := valCtx.Value("task").(domain.Task)
	task.Scenario.Processors = []string{"p1", "p2", "p3", "p4", "p5", "p6"}
	log.Println(task)

	ExecScenario(task.Scenario, valCtx, stream)

	return
}
