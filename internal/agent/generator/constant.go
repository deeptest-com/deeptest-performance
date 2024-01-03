package generator

import (
	"context"
	"github.com/aaronchen2k/deeptest/internal/agent/exec"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	"github.com/aaronchen2k/deeptest/proto"
	"sync"
)

func GenerateConstant(tmplTask domain.Task, stages []*proto.Stage, stream proto.PerformanceService_ExecServer,
	planCtx context.Context, wg *sync.WaitGroup) (err error) {

	if len(stages) != 1 {
		return
	}

	for i := int32(1); i <= stages[0].Target; i++ {
		task := genTask(tmplTask, int(i))

		timeoutCtx, _ := context.WithTimeout(planCtx, consts.ExecTimeout)
		taskCtx := context.WithValue(timeoutCtx, "task", task)

		wg.Add(1)
		go func() {
			defer wg.Done()
			exec.ExecTaskWithVu(taskCtx, &stream)
		}()
	}

	return
}
