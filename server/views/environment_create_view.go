package views

import (
	"vault/vault"
)

type EnvironmentCreateView struct {
	EnvironmentView
	Secret string `json:"secret"`
}

func NewEnvironmentCreateView(environment *vault.Environment, secret string) EnvironmentCreateView {
	environmentViewResult := NewEnvironmentView(environment)
	result := EnvironmentCreateView{EnvironmentView: environmentViewResult, Secret: secret}

	return result
}
