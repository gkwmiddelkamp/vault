package vault

type TokenType int

const (
	Public           = 0
	ReadOnly         = 10
	ReadWrite        = 20
	EnvironmentAdmin = 50
	MasterAdmin      = 100
)

func (t TokenType) String() string {
	return [...]string{"Public", "ReadOnly", "ReadWrite", "EnvironmentAdmin", "MasterAdmin"}[t]
}
