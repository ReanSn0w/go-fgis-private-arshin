package client

import (
	"context"
	"net/http"
)

type FieldType struct {
	Title  string `json:"title"`
	Link   bool   `json:"link"`
	Attach bool   `json:"attach"`
	ID     bool   `json:"id"`
}

type fieldTypesResponse struct {
	Status  int         `json:"status"`
	Result  []FieldType `json:"result"`
	Message string      `json:"message"`
	Trace   string      `json:"trace"`
}

func (c *Client) ListFieldTypes(ctx context.Context) ([]FieldType, error) {
	var data fieldTypesResponse
	if err := c.getJSON(ctx, "fieldtypes", nil, "", &data); err != nil {
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
