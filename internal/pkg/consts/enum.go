package consts

type Instruction string

const (
	Exit   Instruction = "exit"
	Result Instruction = "result"
)

func (e Instruction) String() string {
	return string(e)
}
