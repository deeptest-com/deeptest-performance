package logs

import (
	"github.com/aaronchen2k/deeptest/internal/agent/store"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	"github.com/aaronchen2k/deeptest/proto"
	"github.com/jinzhu/copier"
	"time"
)

func GetSummary() (summary proto.PerformanceExecSummary) {
	data := store.GetData()

	copier.CopyWithOption(&summary, data, copier.Option{
		DeepCopy: true,
	})

	return
}

func Count(result proto.PerformanceExecResult) (err error) {
	status := consts.ResultStatus(result.Record.Status)

	// deal with request
	request := domain.RequestItem{
		Status: status,
		Dur:    int(result.Record.Duration),
	}
	store.AddRequests(request)

	// deal with status
	if status == consts.Pass {
		store.AddPass(1)
	} else if status == consts.Fail {
		store.AddFail(1)
	} else if status == consts.Error {
		store.AddError(1)
	}

	currTime := time.Now().UnixMilli()
	startTime := store.GetStartTime()
	store.UpdateEndTime(currTime)

	duration := currTime - startTime
	store.UpdateDuration(duration) // 毫秒

	// count average duration
	avgDuration := computeAvgDuration(result.Record.Duration)
	store.UpdateAvgDuration(avgDuration)

	// count average qps, must put after all other actions
	avgQps := computeAvgQps()
	store.UpdateAvgQps(avgQps)

	return
}

func computeAvgDuration(requestDur int64) (ret int) {
	requestNum := len(store.GetRequests())
	oldVal := store.GetAvgDuration()

	ret = (oldVal*requestNum+int(requestDur))/requestNum + 1
	store.UpdateAvgDuration(ret)

	return
}

func computeAvgQps() (ret int) {
	successNum := store.GetPass()

	duration := store.GetDuration()

	ret = successNum * 1000 / int(duration)

	store.UpdateAvgQps(ret)

	return
}
