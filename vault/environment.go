package vault

const EnvironmentCollection = "environment"

type Environment struct {
	Name    string `bson:"name,omitempty"`
	Contact string `bson:"contact,omitempty"`
	Active  bool   `bson:"active,omitempty"`
}

func NewEnvironment(name string, contact string, active bool) Environment {
	return Environment{
		Name:    name,
		Contact: contact,
		Active:  active,
	}
}
