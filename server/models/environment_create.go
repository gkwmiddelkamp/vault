package models

type EnvironmentCreate struct {
	Name    string `json:"name"`
	Contact string `json:"contact,omitempty"`
	Active  bool   `json:"active"`
}
