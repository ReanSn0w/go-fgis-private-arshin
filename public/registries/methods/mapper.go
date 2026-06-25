package methods

import (
	"context"
	"fmt"
	"strconv"

	arshin "github.com/ReanSn0w/go-fgis-private-arshin/public/client"
)

const (
	fieldRegistryNumber       = "foei:NumRegCMM"
	fieldRecordNumber         = "foei:number"
	fieldName                 = "foei:NameCMM"
	fieldStatus               = "foei:StatusCMM"
	fieldSystemStatus         = "foei:status"
	fieldPublishedAt          = "foei:date"
	fieldCertificateDate      = "foei:DataCertCMM"
	fieldCertificateNumber    = "foei:NumCertValCMM"
	fieldDeveloperName        = "foei:NameDevMetdCMM"
	fieldCertificationOrgName = "foei:NameCerOrgCMM"
	fieldMeasuredVariable     = "foei:MeasuredVarCMM"
	fieldMeasurementRange     = "foei:RangeMeasuCMM"
	fieldErrorCharacteristic  = "foei:ErrorCharactCMM"
	fieldMethodKind           = "foei:TypeMetCMM"
	fieldMeasurementKind      = "foei:TypeMeasuCMM"
	fieldRegistrySection      = "foei:SFCMM"
	fieldOrderNumber          = "foei:NumberOrderMetCMM"
	fieldOrderDate            = "foei:DataOrderMetCMM"
	fieldOrderDocument        = "foei:DocOrderMetCMM"
	fieldRefOrderNumber       = "foei:NumberRefOrderMetCMM"
	fieldRefOrderDate         = "foei:DataRefOrderMetCMM"
	fieldRefOrderDocument     = "foei:DocRefOrderMetCMM"
	fieldTextDocument         = "foei:DocCMM"
	fieldRangeDocument        = "foei:RMDocCMM"
	fieldErrorDocument        = "foei:ECDocCMM"
	fieldCertificateDocument  = "foei:CertValCMM"
	fieldMethodDataDocument   = "foei:MetodDocCMM"
	fieldCMM1RelationAssoc    = "foei:CMM1Relation_assoc"
	fieldRelationType         = "foei:RelationTypeCMM"
	fieldRelatedMethod        = "foei:RelatedCMM"
)

func (c *Client) GetByRegistryNumber(ctx context.Context, number string) (*Method, error) {
	page, err := c.arshin.ListRegistryData(ctx, c.spec.RegistryID, arshin.RegistryQuery{
		PageNumber: 1,
		PageSize:   2,
		Filters: []arshin.Filter{
			{Field: fieldRegistryNumber, Value: number},
		},
	})
	if err != nil {
		return nil, err
	}
	if len(page.Items) == 0 {
		return nil, fmt.Errorf("methods: method %q not found", number)
	}
	if len(page.Items) > 1 {
		return nil, fmt.Errorf("methods: method %q is not unique", number)
	}
	return MapRecord(page.Items[0]), nil
}

func (c *Client) GetByRegistryNumberDetailed(ctx context.Context, number string) (*Method, error) {
	method, err := c.GetByRegistryNumber(ctx, number)
	if err != nil {
		return nil, err
	}

	item, err := c.arshin.GetRegistryItem(ctx, c.spec.RegistryID, arshin.RegistryItemID(method.ID))
	if err != nil {
		return nil, err
	}

	detailed := MapItem(*item)
	if detailed.Type == "" {
		detailed.Type = method.Type
	}
	if detailed.RegistryID == "" {
		detailed.RegistryID = c.spec.RegistryID
	}
	detailed.RawRecord = method.RawRecord
	return detailed, nil
}

func MapRecord(record arshin.RegistryRecord) *Method {
	properties := record.PropertiesByName()
	spec, _ := SpecForItemType(record.Type)

	method := mapProperties(properties, spec)
	method.ID = record.ID
	method.Type = record.Type
	method.RegistryID = spec.RegistryID
	method.RawRecord = record
	return method
}

func MapItem(item arshin.RegistryItem) *Method {
	registryID := arshin.RegistryID(strconv.FormatInt(item.RegistryID, 10))
	properties := item.FieldsByName()
	spec, _ := SpecForRegistry(registryID)

	method := mapProperties(properties, spec)
	method.ID = strconv.FormatInt(item.ID, 10)
	method.RegistryID = registryID
	method.Type = spec.ItemType
	method.RawItem = &item
	return method
}

func mapProperties(properties map[string]arshin.RegistryProperty, spec RegistrySpec) *Method {
	orderFields := resolveOrderFields(spec, properties)
	orderNumber := stringProperty(properties, orderFields.Number)
	orderDate := stringProperty(properties, orderFields.Date)
	orderDocument := stringProperty(properties, orderFields.Document)

	return &Method{
		RegistryNumber:       stringProperty(properties, fieldRegistryNumber),
		RecordNumber:         stringProperty(properties, fieldRecordNumber),
		Name:                 stringProperty(properties, fieldName),
		Status:               stringProperty(properties, fieldStatus),
		SystemStatus:         stringProperty(properties, fieldSystemStatus),
		PublishedAt:          stringProperty(properties, fieldPublishedAt),
		CertificateDate:      stringProperty(properties, fieldCertificateDate),
		CertificateNumber:    stringProperty(properties, fieldCertificateNumber),
		DeveloperName:        stringProperty(properties, fieldDeveloperName),
		CertificationOrgName: stringProperty(properties, fieldCertificationOrgName),
		MeasuredVariable:     stringProperty(properties, fieldMeasuredVariable),
		MeasurementRange:     stringProperty(properties, fieldMeasurementRange),
		ErrorCharacteristic:  stringProperty(properties, fieldErrorCharacteristic),
		MethodKind:           stringProperty(properties, fieldMethodKind),
		MeasurementKind:      stringProperty(properties, fieldMeasurementKind),
		RegistrySection:      stringProperty(properties, fieldRegistrySection),
		OrderNumber:          orderNumber,
		OrderDate:            orderDate,
		OrderDocument:        orderDocument,
		OrderDocuments:       attachments(properties, orderFields.Document),
		TextDocuments:        attachments(properties, fieldTextDocument),
		RangeDocuments:       attachments(properties, fieldRangeDocument),
		ErrorDocuments:       attachments(properties, fieldErrorDocument),
		CertificateDocuments: attachments(properties, fieldCertificateDocument),
		MethodDataDocuments:  attachments(properties, fieldMethodDataDocument),
		Relations:            relations(properties, fieldCMM1RelationAssoc),
		RawProperties:        properties,
	}
}

func resolveOrderFields(spec RegistrySpec, properties map[string]arshin.RegistryProperty) OrderFields {
	if spec.OrderFields != (OrderFields{}) {
		return spec.OrderFields
	}
	if _, ok := properties[fieldOrderNumber]; ok {
		return OrderFields{
			Number:   fieldOrderNumber,
			Date:     fieldOrderDate,
			Document: fieldOrderDocument,
		}
	}
	if _, ok := properties[fieldRefOrderNumber]; ok {
		return OrderFields{
			Number:   fieldRefOrderNumber,
			Date:     fieldRefOrderDate,
			Document: fieldRefOrderDocument,
		}
	}
	return OrderFields{}
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

func attachments(properties map[string]arshin.RegistryProperty, name string) []arshin.AttachmentRef {
	property, ok := properties[name]
	if !ok {
		return nil
	}
	return property.Attachments()
}

func relations(properties map[string]arshin.RegistryProperty, name string) []Relation {
	property, ok := properties[name]
	if !ok {
		return nil
	}

	children, err := property.ChildObjects()
	if err != nil {
		return nil
	}

	result := make([]Relation, 0, len(children))
	for _, child := range children {
		fields := child.FieldsByName()
		relation := Relation{
			Type:                  stringProperty(fields, fieldRelationType),
			RelatedRegistryNumber: stringProperty(fields, fieldRelatedMethod),
		}
		refs := fields[fieldRelatedMethod].ItemRefs()
		if len(refs) > 0 {
			relation.RelatedRef = &refs[0]
		}
		if relation.Type != "" || relation.RelatedRegistryNumber != "" || relation.RelatedRef != nil {
			result = append(result, relation)
		}
	}
	return result
}
