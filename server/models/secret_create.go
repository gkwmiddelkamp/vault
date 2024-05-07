package models

type SecretCreate struct {
	Description string `json:"description"`
	Secret      string `json:"secret"`
}
