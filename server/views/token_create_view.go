package views

import (
	"vault/vault"
)

type TokenCreateView struct {
	TokenView
	Secret string `json:"secret"`
}

func NewTokenCreateView(token vault.Token) TokenCreateView {
	tokenViewResult := NewTokenView(token)
	result := TokenCreateView{TokenView: tokenViewResult, Secret: token.GetSecret()}

	return result
}
