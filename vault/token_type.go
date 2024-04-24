package vault

type TokenType int

const (
	Public           = 0
	ReadOnly         = 1
	ReadWrite        = 2
	EnvironmentAdmin = 3
	MasterAdmin      = 4
)

func (t TokenType) String() string {
	return [...]string{"Public", "ReadOnly", "ReadWrite", "EnvironmentAdmin", "MasterAdmin"}[t]
}
