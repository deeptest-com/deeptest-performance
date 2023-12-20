package serverConsts

type VuGenerator string

const (
	ConstantVus VuGenerator = "constant_vus"
)

func (e VuGenerator) String() string {
	return string(e)
}
