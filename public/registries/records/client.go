package records

import (
	"context"
	"fmt"

	arshin "github.com/ReanSn0w/go-fgis-private-arshin/public/client"
	"github.com/ReanSn0w/go-fgis-private-arshin/public/registries"
)

type Client struct {
	arshin *arshin.Client
	spec   registries.Spec
}

func NewForRegistry(client *arshin.Client, registryID arshin.RegistryID) (*Client, error) {
	spec, ok := registries.SpecForRegistry(registryID)
	if !ok {
		return nil, ErrUnknownRegistry{RegistryID: registryID}
	}
	return NewForSpec(client, spec), nil
}

func NewForSpec(client *arshin.Client, spec registries.Spec) *Client {
	return &Client{arshin: client, spec: spec}
}

func (c *Client) RegistrySpec() registries.Spec {
	return c.spec
}

func (c *Client) List(ctx context.Context, query arshin.RegistryQuery) (*Page, error) {
	if c.arshin == nil {
		return nil, fmt.Errorf("records: arshin client is nil")
	}

	raw, err := c.arshin.ListRegistryData(ctx, c.spec.RegistryID, query)
	if err != nil {
		return nil, err
	}

	items := make([]Record, 0, len(raw.Items))
	for _, item := range raw.Items {
		items = append(items, MapRecord(c.spec.RegistryID, item))
	}

	return &Page{
		TotalCount:  raw.TotalCount,
		CurrentPage: raw.CurrentPage,
		PageSize:    raw.PageSize,
		Items:       items,
		RawPage:     *raw,
	}, nil
}

func (c *Client) Get(ctx context.Context, itemID arshin.RegistryItemID) (Record, error) {
	if c.arshin == nil {
		return nil, fmt.Errorf("records: arshin client is nil")
	}

	item, err := c.arshin.GetRegistryItem(ctx, c.spec.RegistryID, itemID)
	if err != nil {
		return nil, err
	}
	return MapItem(*item), nil
}

type ErrUnknownRegistry struct {
	RegistryID arshin.RegistryID
}

func (e ErrUnknownRegistry) Error() string {
	return "records: unknown registry " + string(e.RegistryID)
}
