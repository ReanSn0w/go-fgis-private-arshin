package records

import (
	arshin "github.com/ReanSn0w/go-fgis-private-arshin/public/client"
	"github.com/ReanSn0w/go-fgis-private-arshin/public/registries"
)

type Record interface {
	BaseRecord() Base
}

type Base struct {
	ID            string                             `json:"id"`
	Type          string                             `json:"type,omitempty"`
	RegistryID    arshin.RegistryID                  `json:"registryId,omitempty"`
	RegistryTitle string                             `json:"registryTitle,omitempty"`
	RawProperties map[string]arshin.RegistryProperty `json:"-"`
	RawRecord     arshin.RegistryRecord              `json:"-"`
	RawItem       *arshin.RegistryItem               `json:"-"`
}

func (b Base) BaseRecord() Base {
	return b
}

type Page struct {
	TotalCount  int                     `json:"totalCount"`
	CurrentPage int                     `json:"currentPage"`
	PageSize    int                     `json:"pageSize"`
	Items       []Record                `json:"items"`
	RawPage     arshin.RegistryDataPage `json:"-"`
}

type NormativeDocument struct {
	Base
	Kind         string                 `json:"kind,omitempty"`
	Number       string                 `json:"number,omitempty"`
	AcceptedAt   string                 `json:"acceptedAt,omitempty"`
	Organization string                 `json:"organization,omitempty"`
	Name         string                 `json:"name,omitempty"`
	Edition      string                 `json:"edition,omitempty"`
	Status       string                 `json:"status,omitempty"`
	Documents    []arshin.AttachmentRef `json:"documents,omitempty"`
}

type CalibrationMarkCipher struct {
	Base
	Code                  string                 `json:"code,omitempty"`
	OrganizationName      string                 `json:"organizationName,omitempty"`
	RegistrationAuthority string                 `json:"registrationAuthority,omitempty"`
	Documents             []arshin.AttachmentRef `json:"documents,omitempty"`
}

type GSIStandard struct {
	Base
	DocumentType string                 `json:"documentType,omitempty"`
	Number       string                 `json:"number,omitempty"`
	Name         string                 `json:"name,omitempty"`
	IntroducedAt string                 `json:"introducedAt,omitempty"`
	Status       string                 `json:"status,omitempty"`
	Documents    []arshin.AttachmentRef `json:"documents,omitempty"`
}

type InternationalTreaty struct {
	Base
	TypeAndDesignation string                 `json:"typeAndDesignation,omitempty"`
	Date               string                 `json:"date,omitempty"`
	Venue              string                 `json:"venue,omitempty"`
	Organization       string                 `json:"organization,omitempty"`
	Name               string                 `json:"name,omitempty"`
	Documents          []arshin.AttachmentRef `json:"documents,omitempty"`
}

type InternationalComparison struct {
	Base
	StandardRefs   []arshin.RegistryItemRef `json:"standardRefs,omitempty"`
	Code           string                   `json:"code,omitempty"`
	ComparisonType string                   `json:"comparisonType,omitempty"`
	Description    string                   `json:"description,omitempty"`
	MeasuredValue  string                   `json:"measuredValue,omitempty"`
	Years          string                   `json:"years,omitempty"`
	Documents      []arshin.AttachmentRef   `json:"documents,omitempty"`
}

type ReferenceData struct {
	Base
	Number           string                 `json:"number,omitempty"`
	Name             string                 `json:"name,omitempty"`
	Storage          string                 `json:"storage,omitempty"`
	DevelopmentState string                 `json:"developmentState,omitempty"`
	Status           string                 `json:"status,omitempty"`
	Documents        []arshin.AttachmentRef `json:"documents,omitempty"`
}

type StandardUnit struct {
	Base
	RegistryNumber string                   `json:"registryNumber,omitempty"`
	Name           string                   `json:"name,omitempty"`
	OrderNumber    string                   `json:"orderNumber,omitempty"`
	OrderDate      string                   `json:"orderDate,omitempty"`
	Documents      []arshin.AttachmentRef   `json:"documents,omitempty"`
	Relations      []arshin.RegistryItemRef `json:"relations,omitempty"`
}

type PrimaryStandard struct {
	Base
	RegistryNumber    string                   `json:"registryNumber,omitempty"`
	StandardName      string                   `json:"standardName,omitempty"`
	InstituteRefs     []arshin.RegistryItemRef `json:"instituteRefs,omitempty"`
	ApprovalYear      string                   `json:"approvalYear,omitempty"`
	CertificationYear string                   `json:"certificationYear,omitempty"`
	Status            string                   `json:"status,omitempty"`
	Documents         []arshin.AttachmentRef   `json:"documents,omitempty"`
	Relations         []arshin.RegistryItemRef `json:"relations,omitempty"`
}

type InternationalDocument struct {
	Base
	Number              string                 `json:"number,omitempty"`
	RussianName         string                 `json:"russianName,omitempty"`
	RussianOrganization string                 `json:"russianOrganization,omitempty"`
	RussianDocuments    []arshin.AttachmentRef `json:"russianDocuments,omitempty"`
	EnglishDocuments    []arshin.AttachmentRef `json:"englishDocuments,omitempty"`
	Documents           []arshin.AttachmentRef `json:"documents,omitempty"`
}

type TechnicalInstrumentClassification struct {
	Base
	DocumentType string                 `json:"documentType,omitempty"`
	Number       string                 `json:"number,omitempty"`
	Date         string                 `json:"date,omitempty"`
	Name         string                 `json:"name,omitempty"`
	Status       string                 `json:"status,omitempty"`
	Documents    []arshin.AttachmentRef `json:"documents,omitempty"`
}

type RegulatedMeasurementList struct {
	Base
	Documents []arshin.AttachmentRef `json:"documents,omitempty"`
}

type StandardSampleType struct {
	Base
	RegistryNumber          string                   `json:"registryNumber,omitempty"`
	Name                    string                   `json:"name,omitempty"`
	CertifiedCharacteristic string                   `json:"certifiedCharacteristic,omitempty"`
	Expiration              string                   `json:"expiration,omitempty"`
	ProductionRefs          []arshin.RegistryItemRef `json:"productionRefs,omitempty"`
	Documents               []arshin.AttachmentRef   `json:"documents,omitempty"`
	Relations               []arshin.RegistryItemRef `json:"relations,omitempty"`
}

type ProductionNotice struct {
	Base
	RegistrationDate         string                 `json:"registrationDate,omitempty"`
	RegistrationNumber       string                 `json:"registrationNumber,omitempty"`
	OrganizationFullName     string                 `json:"organizationFullName,omitempty"`
	OrganizationINN          string                 `json:"organizationINN,omitempty"`
	OrganizationLegalAddress string                 `json:"organizationLegalAddress,omitempty"`
	Documents                []arshin.AttachmentRef `json:"documents,omitempty"`
}

type GenericRecord struct {
	Base
	Spec      registries.Spec          `json:"spec"`
	Documents []arshin.AttachmentRef   `json:"documents,omitempty"`
	Relations []arshin.RegistryItemRef `json:"relations,omitempty"`
}
