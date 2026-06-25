package client

import (
	"context"
	"fmt"
	"net/http"
)

type RegistryItemID string

type RegistryItem struct {
	ID            int64                 `json:"id"`
	AlfrescoID    string                `json:"alfrescoId"`
	RegistryID    int64                 `json:"registryId"`
	RegistryTitle string                `json:"registryTitle"`
	Deleted       bool                  `json:"deleted"`
	Sections      []RegistryItemSection `json:"sections"`
}

type RegistryItemSection struct {
	SectionTitle string             `json:"sectionTitle"`
	Fields       []RegistryProperty `json:"fields"`
}

type registryItemResponse struct {
	Status  int          `json:"status"`
	Result  RegistryItem `json:"result"`
	Message string       `json:"message"`
	Trace   string       `json:"trace"`
}

type registryPlainDataResponse struct {
	Status  int            `json:"status"`
	Result  RegistryRecord `json:"result"`
	Message string         `json:"message"`
	Trace   string         `json:"trace"`
}

func (i RegistryItem) Fields() []RegistryProperty {
	var fields []RegistryProperty
	for _, section := range i.Sections {
		fields = append(fields, section.Fields...)
	}
	return fields
}

func (i RegistryItem) FieldsByName() map[string]RegistryProperty {
	return PropertiesByName(i.Fields())
}

func (c *Client) GetRegistryItem(ctx context.Context, registryID RegistryID, itemID RegistryItemID) (*RegistryItem, error) {
	if registryID == "" {
		return nil, fmt.Errorf("registry id is empty")
	}
	if itemID == "" {
		return nil, fmt.Errorf("item id is empty")
	}

	registryIDString := string(registryID)
	var data registryItemResponse
	if err := c.getJSON(ctx, "registry/"+registryIDString+"/items/"+string(itemID)+"/data", nil, registryIDString, &data); err != nil {
		return nil, err
	}
	if data.Status != http.StatusOK {
		return nil, &APIError{
			Status:  data.Status,
			Message: data.Message,
			Trace:   data.Trace,
		}
	}

	return &data.Result, nil
}

func (c *Client) GetRegistryItemPlainData(ctx context.Context, registryID RegistryID, itemID RegistryItemID) (*RegistryRecord, error) {
	if registryID == "" {
		return nil, fmt.Errorf("registry id is empty")
	}
	if itemID == "" {
		return nil, fmt.Errorf("item id is empty")
	}

	registryIDString := string(registryID)
	var data registryPlainDataResponse
	if err := c.getJSON(ctx, "registry/"+registryIDString+"/items/"+string(itemID)+"/plaindata", nil, registryIDString, &data); err != nil {
		return nil, err
	}
	if data.Status != http.StatusOK {
		return nil, &APIError{
			Status:  data.Status,
			Message: data.Message,
			Trace:   data.Trace,
		}
	}

	return &data.Result, nil
}
