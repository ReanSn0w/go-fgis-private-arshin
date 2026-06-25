package methods

import (
	arshin "github.com/ReanSn0w/go-fgis-private-arshin/public/client"
	"github.com/ReanSn0w/go-fgis-private-arshin/public/registries"
)

type RegistrySpec struct {
	RegistryID  arshin.RegistryID
	Title       string
	ItemType    string
	OrderFields OrderFields
}

type OrderFields struct {
	Number   string
	Date     string
	Document string
}

var registrySpecs = map[arshin.RegistryID]RegistrySpec{
	CertifiedRegistryID: {
		RegistryID: CertifiedRegistryID,
		Title:      registries.CMM1.Title,
		ItemType:   registries.CMM1.ItemType,
	},
	PrimaryReferenceRegistryID: {
		RegistryID: PrimaryReferenceRegistryID,
		Title:      registries.CMM2.Title,
		ItemType:   registries.CMM2.ItemType,
		OrderFields: OrderFields{
			Number:   fieldOrderNumber,
			Date:     fieldOrderDate,
			Document: fieldOrderDocument,
		},
	},
	ReferenceRegistryID: {
		RegistryID: ReferenceRegistryID,
		Title:      registries.CMM3.Title,
		ItemType:   registries.CMM3.ItemType,
		OrderFields: OrderFields{
			Number:   fieldRefOrderNumber,
			Date:     fieldRefOrderDate,
			Document: fieldRefOrderDocument,
		},
	},
}

func SpecForRegistry(registryID arshin.RegistryID) (RegistrySpec, bool) {
	spec, ok := registrySpecs[registryID]
	return spec, ok
}

func SpecForItemType(itemType string) (RegistrySpec, bool) {
	for _, spec := range registrySpecs {
		if spec.ItemType == itemType {
			return spec, true
		}
	}
	return RegistrySpec{}, false
}

func SupportedRegistries() []RegistrySpec {
	return []RegistrySpec{
		registrySpecs[PrimaryReferenceRegistryID],
		registrySpecs[ReferenceRegistryID],
		registrySpecs[CertifiedRegistryID],
	}
}
