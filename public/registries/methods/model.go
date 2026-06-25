package methods

import arshin "github.com/ReanSn0w/go-fgis-private-arshin/public/client"

type Method struct {
	ID                   string
	Type                 string
	RegistryID           arshin.RegistryID
	RegistryNumber       string
	RecordNumber         string
	Name                 string
	Status               string
	SystemStatus         string
	PublishedAt          string
	CertificateDate      string
	CertificateNumber    string
	DeveloperName        string
	CertificationOrgName string
	MeasuredVariable     string
	MeasurementRange     string
	ErrorCharacteristic  string
	MethodKind           string
	MeasurementKind      string
	RegistrySection      string
	OrderNumber          string
	OrderDate            string
	OrderDocument        string
	OrderDocuments       []arshin.AttachmentRef
	TextDocuments        []arshin.AttachmentRef
	RangeDocuments       []arshin.AttachmentRef
	ErrorDocuments       []arshin.AttachmentRef
	CertificateDocuments []arshin.AttachmentRef
	MethodDataDocuments  []arshin.AttachmentRef
	Relations            []Relation
	RawProperties        map[string]arshin.RegistryProperty
	RawRecord            arshin.RegistryRecord
	RawItem              *arshin.RegistryItem
}

type Relation struct {
	Type                  string
	RelatedRegistryNumber string
	RelatedRef            *arshin.RegistryItemRef
}
