package constants

type VerifyConst string

const (
	PhoneVerifyConst   VerifyConst = "phv"
	EmailVerifyConst   VerifyConst = "emv"
	ProvideVerifyConst VerifyConst = "prv"
	NoneVerifyConst    VerifyConst = "nonv"
)

func (v VerifyConst) String() string {
	return string(v)
}
