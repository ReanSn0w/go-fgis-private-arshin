package client

type RegistryProperty struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	Multiple   bool   `json:"multiple"`
	Title      string `json:"title"`
	Constraint any    `json:"constraint"`
	Value      any    `json:"value"`
	LongValue  any    `json:"longValue"`
	Link       any    `json:"link"`
	MIME       any    `json:"mime"`
}

func PropertiesByName(properties []RegistryProperty) map[string]RegistryProperty {
	index := make(map[string]RegistryProperty, len(properties))
	for _, property := range properties {
		index[property.Name] = property
	}
	return index
}
