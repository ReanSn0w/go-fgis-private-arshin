package records

import (
	"fmt"
	"strconv"

	arshin "github.com/ReanSn0w/go-fgis-private-arshin/public/client"
	"github.com/ReanSn0w/go-fgis-private-arshin/public/registries"
)

func MapRecord(registryID arshin.RegistryID, record arshin.RegistryRecord) Record {
	spec, _ := registries.SpecForRegistry(registryID)
	base := baseFromRecord(spec, record)
	return mapProperties(base, spec, record.Properties)
}

func MapItem(item arshin.RegistryItem) Record {
	registryID := arshin.RegistryID(strconv.FormatInt(item.RegistryID, 10))
	spec, _ := registries.SpecForRegistry(registryID)
	base := baseFromItem(spec, item)
	return mapProperties(base, spec, item.Fields())
}

func baseFromRecord(spec registries.Spec, record arshin.RegistryRecord) Base {
	return Base{
		ID:            record.ID,
		Type:          record.Type,
		RegistryID:    spec.RegistryID,
		RegistryTitle: spec.Title,
		RawProperties: record.PropertiesByName(),
		RawRecord:     record,
	}
}

func baseFromItem(spec registries.Spec, item arshin.RegistryItem) Base {
	return Base{
		ID:            strconv.FormatInt(item.ID, 10),
		Type:          spec.ItemType,
		RegistryID:    spec.RegistryID,
		RegistryTitle: item.RegistryTitle,
		RawProperties: item.FieldsByName(),
		RawItem:       &item,
	}
}

func mapProperties(base Base, spec registries.Spec, properties []arshin.RegistryProperty) Record {
	index := arshin.PropertiesByName(properties)

	switch spec.RegistryID {
	case registries.NDRegistryID:
		return NormativeDocument{
			Base:         base,
			Kind:         stringProperty(index, "foei:TypeND"),
			Number:       stringProperty(index, "foei:NumberND"),
			AcceptedAt:   stringProperty(index, "foei:DataND"),
			Organization: stringProperty(index, "foei:OrgND"),
			Name:         stringProperty(index, "foei:NameND"),
			Edition:      stringProperty(index, "foei:EditionND"),
			Status:       stringProperty(index, "foei:StatusND"),
			Documents:    attachments(properties),
		}
	case registries.SCMRegistryID:
		return CalibrationMarkCipher{
			Base:                  base,
			Code:                  stringProperty(index, "foei:SCM"),
			OrganizationName:      stringProperty(index, "foei:NameOrgSCM"),
			RegistrationAuthority: stringProperty(index, "foei:RegSCM"),
			Documents:             attachments(properties),
		}
	case registries.GSIRegistryID:
		return GSIStandard{
			Base:         base,
			DocumentType: stringProperty(index, "foei:TypeGSI"),
			Number:       stringProperty(index, "foei:NumberGSI"),
			Name:         stringProperty(index, "foei:NameGSI"),
			IntroducedAt: stringProperty(index, "foei:DateGSI"),
			Status:       stringProperty(index, "foei:StatusGSI"),
			Documents:    attachments(properties),
		}
	case registries.MDGRegistryID:
		return InternationalTreaty{
			Base:               base,
			TypeAndDesignation: stringProperty(index, "foei:TypeMDG"),
			Date:               stringProperty(index, "foei:DateMDG"),
			Venue:              stringProperty(index, "foei:VenueMDG"),
			Organization:       stringProperty(index, "foei:OrganizationMDG"),
			Name:               stringProperty(index, "foei:NameMDG"),
			Documents:          attachments(properties),
		}
	case registries.ICRegistryID:
		return InternationalComparison{
			Base:           base,
			StandardRefs:   itemRefs(index, "foei:StandartIC"),
			Code:           stringProperty(index, "foei:CodeIC"),
			ComparisonType: stringProperty(index, "foei:TypeIC"),
			Description:    stringProperty(index, "foei:DescriptionIC"),
			MeasuredValue:  stringProperty(index, "foei:MeasuredValueIC"),
			Years:          stringProperty(index, "foei:YearIC"),
			Documents:      attachments(properties),
		}
	case registries.SSDRegistryID:
		return ReferenceData{
			Base:             base,
			Number:           stringProperty(index, "foei:NumberSSD"),
			Name:             stringProperty(index, "foei:NameSSD"),
			Storage:          stringProperty(index, "foei:StorageSSD"),
			DevelopmentState: stringProperty(index, "foei:StateDevSSD"),
			Status:           stringProperty(index, "foei:StatusSSD"),
			Documents:        attachments(properties),
		}
	case registries.SURegistryID:
		return StandardUnit{
			Base:           base,
			RegistryNumber: stringProperty(index, "foei:NumRegSU"),
			Name:           stringProperty(index, "foei:NameSU"),
			OrderNumber:    stringProperty(index, "foei:NumberOrderSU"),
			OrderDate:      stringProperty(index, "foei:DataOrderSU"),
			Documents:      attachments(properties),
			Relations:      relations(properties),
		}
	case registries.GPSRegistryID:
		return PrimaryStandard{
			Base:              base,
			RegistryNumber:    stringProperty(index, "foei:NumRegGPS"),
			StandardName:      stringProperty(index, "foei:StandNameGPS"),
			InstituteRefs:     itemRefs(index, "foei:InstGuardGPS"),
			ApprovalYear:      stringProperty(index, "foei:YearApprovGPS"),
			CertificationYear: stringProperty(index, "foei:YearCertifGPS"),
			Status:            stringProperty(index, "foei:StatusGPS"),
			Documents:         attachments(properties),
			Relations:         relations(properties),
		}
	case registries.MDRegistryID:
		russianDocuments := attachmentsByField(index, "foei:DocRusMD")
		englishDocuments := attachmentsByField(index, "foei:DocEngMD")
		return InternationalDocument{
			Base:                base,
			Number:              stringProperty(index, "foei:NumberMD"),
			RussianName:         stringProperty(index, "foei:NameRusMD"),
			RussianOrganization: stringProperty(index, "foei:OrgRusMD"),
			RussianDocuments:    russianDocuments,
			EnglishDocuments:    englishDocuments,
			Documents:           append(append([]arshin.AttachmentRef{}, russianDocuments...), englishDocuments...),
		}
	case registries.TSSIRegistryID:
		return TechnicalInstrumentClassification{
			Base:         base,
			DocumentType: stringProperty(index, "foei:TypeTSSI"),
			Number:       stringProperty(index, "foei:NumberTSSI"),
			Date:         stringProperty(index, "foei:DateTSSI"),
			Name:         stringProperty(index, "foei:NameTSSI"),
			Status:       stringProperty(index, "foei:StatusTSSI"),
			Documents:    attachments(properties),
		}
	case registries.EPIRegistryID:
		return RegulatedMeasurementList{
			Base:      base,
			Documents: attachments(properties),
		}
	case registries.UTSORegistryID:
		return StandardSampleType{
			Base:                    base,
			RegistryNumber:          stringProperty(index, "foei:NumberUTSO"),
			Name:                    stringProperty(index, "foei:NameUTSO"),
			CertifiedCharacteristic: stringProperty(index, "foei:NameCertCharUTSO"),
			Expiration:              stringProperty(index, "foei:ExpDateUTSO"),
			ProductionRefs:          itemRefs(index, "foei:ProductionUTSO"),
			Documents:               attachments(properties),
			Relations:               relations(properties),
		}
	case registries.P1WFRegistryID:
		return ProductionNotice{
			Base:                     base,
			RegistrationDate:         stringProperty(index, "gost:registrationDate"),
			RegistrationNumber:       stringProperty(index, "gost:registrationNumber"),
			OrganizationFullName:     stringProperty(index, "gost:p1wfOrganizationFullName"),
			OrganizationINN:          stringProperty(index, "gost:p1wfOrganizationINN"),
			OrganizationLegalAddress: stringProperty(index, "gost:p1wfOrganizationLegalAddress"),
			Documents:                attachments(properties),
		}
	default:
		return GenericRecord{
			Base:      base,
			Spec:      spec,
			Documents: attachments(properties),
			Relations: relations(properties),
		}
	}
}

func stringProperty(properties map[string]arshin.RegistryProperty, name string) string {
	property, ok := properties[name]
	if !ok || property.Value == nil {
		return ""
	}
	value, ok := property.Value.(string)
	if ok {
		return value
	}
	return fmt.Sprint(property.Value)
}

func attachments(properties []arshin.RegistryProperty) []arshin.AttachmentRef {
	var result []arshin.AttachmentRef
	for _, property := range properties {
		result = append(result, property.Attachments()...)
	}
	return result
}

func attachmentsByField(properties map[string]arshin.RegistryProperty, name string) []arshin.AttachmentRef {
	property, ok := properties[name]
	if !ok {
		return nil
	}
	return property.Attachments()
}

func itemRefs(properties map[string]arshin.RegistryProperty, name string) []arshin.RegistryItemRef {
	property, ok := properties[name]
	if !ok {
		return nil
	}
	return property.ItemRefs()
}

func relations(properties []arshin.RegistryProperty) []arshin.RegistryItemRef {
	var result []arshin.RegistryItemRef
	for _, property := range properties {
		result = append(result, property.ItemRefs()...)
		if property.Type != "CHILD_OBJECT" {
			continue
		}
		children, err := property.ChildObjects()
		if err != nil {
			continue
		}
		for _, child := range children {
			result = append(result, child.ItemRefs()...)
		}
	}
	return result
}
