package consts

type ResultStatus string

const (
	Pass    ResultStatus = "pass"
	Fail    ResultStatus = "fail"
	Error   ResultStatus = "error"
	Block   ResultStatus = "block"
	Unknown ResultStatus = "unknown"
)

func (e ResultStatus) String() string {
	return string(e)
}

type Instruction string

const (
	Exit   Instruction = "exit"
	Record Instruction = "record"
	Result Instruction = "result"
)

func (e Instruction) String() string {
	return string(e)
}

type GeneratorType string

const (
	GeneratorConstant GeneratorType = "constant"
	GeneratorRamp     GeneratorType = "ramp"
)

func (e GeneratorType) String() string {
	return string(e)
}

type TargetType string

const (
	TargetQps       TargetType = "qps"
	TargetDuration  TargetType = "duration"
	TargetErrorRate TargetType = "errorRate"
)

func (e TargetType) String() string {
	return string(e)
}

type WsMsgCategory string

const (
	ProgressInProgress WsMsgCategory = "in_progress"
	ProgressEnd        WsMsgCategory = "end"
	ProgressResult     WsMsgCategory = "result"
	Processor          WsMsgCategory = "processor"
	Initialize         WsMsgCategory = "initialize"
	Statistic          WsMsgCategory = "stat"
	Exception          WsMsgCategory = "exception"
)

func (e WsMsgCategory) String() string {
	return string(e)
}
