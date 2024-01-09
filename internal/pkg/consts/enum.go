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

type MsgCategory string

const (
	MsgCategoryInstruction MsgCategory = "instruction"
	MsgCategoryResult      MsgCategory = "result"
)

func (e MsgCategory) String() string {
	return string(e)
}

type MsgInstruction string

const (
	MsgInstructionStart    MsgInstruction = "start"
	MsgInstructionEnd      MsgInstruction = "end"
	MsgInstructionTerminal MsgInstruction = "terminal"
)

func (e MsgInstruction) String() string {
	return string(e)
}

type MsgResult string

const (
	MsgResultRecord  MsgResult = "record"
	MsgResultSummary MsgResult = "summary"
)

func (e MsgResult) String() string {
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

type ExecType string

const (
	ExecStop ExecType = "stop"

	Init     ExecType = "init"
	ExecPlan ExecType = "execPlan"
)

func (e ExecType) String() string {
	return string(e)
}
