package methods

import (
	arshin "github.com/ReanSn0w/go-fgis-private-arshin/public/client"
	"github.com/ReanSn0w/go-fgis-private-arshin/public/registries"
)

const (
	RegistryID                 arshin.RegistryID = registries.CMM1RegistryID
	CertifiedRegistryID        arshin.RegistryID = registries.CMM1RegistryID
	PrimaryReferenceRegistryID arshin.RegistryID = registries.CMM2RegistryID
	ReferenceRegistryID        arshin.RegistryID = registries.CMM3RegistryID
)

type Client struct {
	arshin *arshin.Client
	spec   RegistrySpec
}

func New(client *arshin.Client) *Client {
	spec, _ := SpecForRegistry(CertifiedRegistryID)
	return &Client{arshin: client, spec: spec}
}

func NewForRegistry(client *arshin.Client, registryID arshin.RegistryID) (*Client, error) {
	spec, ok := SpecForRegistry(registryID)
	if !ok {
		return nil, ErrUnsupportedRegistry{RegistryID: registryID}
	}
	return &Client{arshin: client, spec: spec}, nil
}

func NewForSpec(client *arshin.Client, spec registries.Spec) (*Client, error) {
	methodsSpec, ok := SpecForRegistry(spec.RegistryID)
	if !ok {
		return nil, ErrUnsupportedRegistry{RegistryID: spec.RegistryID}
	}
	return &Client{arshin: client, spec: methodsSpec}, nil
}

func IsSupportedRegistry(registryID arshin.RegistryID) bool {
	_, ok := SpecForRegistry(registryID)
	return ok
}

func (c *Client) RegistrySpec() RegistrySpec {
	return c.spec
}

type ErrUnsupportedRegistry struct {
	RegistryID arshin.RegistryID
}

func (e ErrUnsupportedRegistry) Error() string {
	return "methods: unsupported registry " + string(e.RegistryID)
}
