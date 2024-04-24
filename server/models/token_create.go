package models

type TokenCreate struct {
	Description   string `json:"description"`
	ExpiresAt     string `json:"expiresAt,omitempty"`
	TokenType     string `json:"tokenType"`
	EnvironmentId string `json:"environmentId,omitempty"`
}
