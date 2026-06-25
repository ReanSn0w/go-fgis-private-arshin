package client

type RegistryRecord struct {
	Values      any                `json:"values"`
	Properties  []RegistryProperty `json:"properties"`
	ID          string             `json:"id"`
	AlfrescoID  string             `json:"alfrescoId"`
	NodeRef     any                `json:"nodeRef"`
	Type        string             `json:"type"`
	Permissions any                `json:"permissions"`
}

func (r RegistryRecord) PropertiesByName() map[string]RegistryProperty {
	return PropertiesByName(r.Properties)
}
