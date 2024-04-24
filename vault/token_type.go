package vault

import "errors"

type TokenType int

var tokenTypeNames = [...]string{"Public", "ReadOnly", "ReadWrite", "EnvironmentAdmin", "MasterAdmin"}

const (
	Public           = 0
	ReadOnly         = 1
	ReadWrite        = 2
	EnvironmentAdmin = 3
	MasterAdmin      = 4
)

func (t TokenType) String() string {
	return tokenTypeNames[t]
}

func TokenTypeFromString(input string) (TokenType, error) {
	for k, v := range tokenTypeNames {
		if input == v {
			return TokenType(k), nil
		}
	}
	return 0, errors.New("invalid token type")
}
