package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type RegistryID string

type RegistryInfo struct {
	ID    RegistryID
	Title string
}

type RegistryDetails struct {
	ID              string `json:"id"`
	Title           string `json:"title"`
	Type            string `json:"type"`
	NodeRef         string `json:"nodeRef"`
	Size            *int   `json:"size"`
	AlfrescoID      string `json:"alfrescoId"`
	CreatedAt       int64  `json:"createdAt"`
	UpdatedAt       int64  `json:"updatedAt"`
	CreatedAtFormat string `json:"createdAtFormat"`
	UpdatedAtFormat string `json:"updatedAtFormat"`
	Deleted         bool   `json:"deleted"`
	Writable        *bool  `json:"writable"`
}

type RegistryQuery struct {
	PageNumber int
	PageSize   int
	OrgID      string
	Filters    []Filter
	Sorts      []Sort
}

type Filter struct {
	Field string
	Value string
}

type Sort struct {
	Field     string
	Direction string
}

type RegistryDataResponse struct {
	Status  int              `json:"status"`
	Result  RegistryDataPage `json:"result"`
	Message string           `json:"message"`
	Trace   string           `json:"trace"`
}

type RegistryDataPage struct {
	TotalCount  int              `json:"totalCount"`
	CurrentPage int              `json:"currentPage"`
	PageSize    int              `json:"pageSize"`
	Items       []RegistryRecord `json:"items"`
}

type registryDetailsResponse struct {
	Status  int             `json:"status"`
	Result  RegistryDetails `json:"result"`
	Message string          `json:"message"`
	Trace   string          `json:"trace"`
}

type registryFieldsResponse struct {
	Status  int                `json:"status"`
	Result  []RegistryProperty `json:"result"`
	Message string             `json:"message"`
	Trace   string             `json:"trace"`
}

func (c *Client) GetRegistry(ctx context.Context, registryID RegistryID) (*RegistryDetails, error) {
	if registryID == "" {
		return nil, fmt.Errorf("registry id is empty")
	}

	registryIDString := string(registryID)
	var data registryDetailsResponse
	if err := c.getJSON(ctx, "registry/"+registryIDString, nil, registryIDString, &data); err != nil {
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

func (c *Client) ListRegistryFields(ctx context.Context, registryID RegistryID) ([]RegistryProperty, error) {
	return c.listRegistryProperties(ctx, registryID, "fields")
}

func (c *Client) ListRegistryDisplayFields(ctx context.Context, registryID RegistryID) ([]RegistryProperty, error) {
	return c.listRegistryProperties(ctx, registryID, "displayfields")
}

func (c *Client) listRegistryProperties(ctx context.Context, registryID RegistryID, endpoint string) ([]RegistryProperty, error) {
	if registryID == "" {
		return nil, fmt.Errorf("registry id is empty")
	}

	registryIDString := string(registryID)
	var data registryFieldsResponse
	if err := c.getJSON(ctx, "registry/"+registryIDString+"/"+endpoint, nil, registryIDString, &data); err != nil {
		return nil, err
	}
	if data.Status != http.StatusOK {
		return nil, &APIError{
			Status:  data.Status,
			Message: data.Message,
			Trace:   data.Trace,
		}
	}

	return data.Result, nil
}

func (c *Client) ListRegistryData(ctx context.Context, registryID RegistryID, query RegistryQuery) (*RegistryDataPage, error) {
	if registryID == "" {
		return nil, fmt.Errorf("registry id is empty")
	}

	if query.PageNumber == 0 {
		query.PageNumber = 1
	}
	if query.PageSize == 0 {
		query.PageSize = 20
	}
	if query.OrgID == "" {
		query.OrgID = DefaultOrgID
	}

	values := url.Values{}
	values.Set("pageNumber", fmt.Sprintf("%d", query.PageNumber))
	values.Set("pageSize", fmt.Sprintf("%d", query.PageSize))
	values.Set("orgID", query.OrgID)

	for _, filter := range query.Filters {
		values.Add("filterBy", filter.Field)
		values.Add("filterValues", filter.Value)
	}

	for _, sort := range query.Sorts {
		values.Add("sortBy", sort.Field)
		values.Add("sortValues", sort.Direction)
	}

	registryIDString := string(registryID)
	var data RegistryDataResponse
	if err := c.getJSON(ctx, "registry/"+registryIDString+"/data", values, registryIDString, &data); err != nil {
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
