package logs

import (
	"github.com/aaronchen2k/deeptest/proto"
)

type MessageSender interface {
	Send(result proto.PerformanceExecResult) error
}
