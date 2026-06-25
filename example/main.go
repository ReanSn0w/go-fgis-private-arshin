package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	arshin "github.com/ReanSn0w/go-fgis-private-arshin/public/client"
	"github.com/ReanSn0w/go-fgis-private-arshin/public/registries"
	"github.com/ReanSn0w/go-fgis-private-arshin/public/registries/methods"
)

const defaultDocumentField = "foei:NameCMM"

type report struct {
	GeneratedAt time.Time        `json:"generatedAt"`
	Query       string           `json:"query"`
	Field       string           `json:"field"`
	PageSize    int              `json:"pageSize"`
	MaxPages    int              `json:"maxPages"`
	Registries  []registryReport `json:"registries"`
}

type registryReport struct {
	RegistryID string   `json:"registryId"`
	Title      string   `json:"title,omitempty"`
	TotalCount int      `json:"totalCount"`
	Records    []record `json:"records"`
	Error      string   `json:"error,omitempty"`
}

type record struct {
	Record    recordSummary     `json:"record"`
	Method    *methodSummary    `json:"method,omitempty"`
	Relations []relationSummary `json:"relations,omitempty"`
	Error     string            `json:"error,omitempty"`
}

type relationSummary struct {
	FieldName             string         `json:"fieldName,omitempty"`
	FieldTitle            string         `json:"fieldTitle,omitempty"`
	Type                  string         `json:"type,omitempty"`
	RelatedRegistryNumber string         `json:"relatedRegistryNumber,omitempty"`
	RelatedRef            *itemRef       `json:"relatedRef,omitempty"`
	Record                *recordSummary `json:"record,omitempty"`
	Method                *methodSummary `json:"method,omitempty"`
	Error                 string         `json:"error,omitempty"`
}

type recordSummary struct {
	ID          string          `json:"id"`
	RegistryID  string          `json:"registryId,omitempty"`
	Type        string          `json:"type,omitempty"`
	DisplayName string          `json:"displayName,omitempty"`
	Properties  []propertyValue `json:"properties,omitempty"`
	Documents   []document      `json:"documents,omitempty"`
}

type propertyValue struct {
	Name  string `json:"name"`
	Title string `json:"title,omitempty"`
	Type  string `json:"type,omitempty"`
	Value any    `json:"value,omitempty"`
}

type methodSummary struct {
	ID                   string                `json:"id"`
	RegistryID           string                `json:"registryId,omitempty"`
	RegistryNumber       string                `json:"registryNumber,omitempty"`
	RecordNumber         string                `json:"recordNumber,omitempty"`
	Name                 string                `json:"name,omitempty"`
	Status               string                `json:"status,omitempty"`
	SystemStatus         string                `json:"systemStatus,omitempty"`
	PublishedAt          string                `json:"publishedAt,omitempty"`
	CertificateDate      string                `json:"certificateDate,omitempty"`
	CertificateNumber    string                `json:"certificateNumber,omitempty"`
	DeveloperName        string                `json:"developerName,omitempty"`
	CertificationOrgName string                `json:"certificationOrgName,omitempty"`
	OrderNumber          string                `json:"orderNumber,omitempty"`
	OrderDate            string                `json:"orderDate,omitempty"`
	OrderDocument        string                `json:"orderDocument,omitempty"`
	Documents            map[string][]document `json:"documents,omitempty"`
}

type document struct {
	FileID   string `json:"fileId,omitempty"`
	Name     string `json:"name,omitempty"`
	MIME     string `json:"mime,omitempty"`
	Download string `json:"download,omitempty"`
}

type itemRef struct {
	RegistryID string `json:"registryId"`
	ItemID     string `json:"itemId"`
	Link       string `json:"link,omitempty"`
}

func main() {
	query := flag.String("query", "", "document name or search string")
	registryIDs := flag.String("registry", "", "comma-separated registry ids; empty means all known registries")
	field := flag.String("field", defaultDocumentField, "registry field used for document search")
	pageSize := flag.Int("page-size", 20, "page size for search requests")
	maxPages := flag.Int("max-pages", 5, "maximum number of search pages to read")
	timeout := flag.Duration("timeout", 2*time.Minute, "request timeout")
	out := flag.String("out", "", "write JSON output to file instead of stdout")
	flag.Parse()

	if *query == "" {
		*query = strings.TrimSpace(strings.Join(flag.Args(), " "))
	}
	if strings.TrimSpace(*query) == "" {
		exit(fmt.Errorf("query is required: use -query or pass it as positional arguments"))
	}
	if *pageSize <= 0 {
		exit(fmt.Errorf("page-size must be positive"))
	}
	if *maxPages <= 0 {
		exit(fmt.Errorf("max-pages must be positive"))
	}

	registries, err := parseRegistries(*registryIDs)
	if err != nil {
		exit(err)
	}

	client, err := arshin.NewClient()
	if err != nil {
		exit(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	result := findDocumentRelations(ctx, client, registries, *field, *query, *pageSize, *maxPages)

	if err := writeJSON(*out, result); err != nil {
		exit(err)
	}
}

func parseRegistries(value string) ([]registries.Spec, error) {
	if strings.TrimSpace(value) == "" {
		return registries.Known(), nil
	}

	parts := strings.Split(value, ",")
	result := make([]registries.Spec, 0, len(parts))
	for _, part := range parts {
		registryID := arshin.RegistryID(strings.TrimSpace(part))
		if registryID == "" {
			continue
		}

		spec, ok := registries.SpecForRegistry(registryID)
		if !ok {
			return nil, fmt.Errorf("unknown registry %q", registryID)
		}
		result = append(result, spec)
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("registry list is empty")
	}
	return result, nil
}

func findDocumentRelations(ctx context.Context, client *arshin.Client, registries []registries.Spec, field string, query string, pageSize int, maxPages int) report {
	result := report{
		GeneratedAt: time.Now().UTC(),
		Query:       query,
		Field:       field,
		PageSize:    pageSize,
		MaxPages:    maxPages,
	}

	for _, registry := range registries {
		result.Registries = append(result.Registries, findRegistryDocumentRelations(ctx, client, registry, field, query, pageSize, maxPages))
	}

	return result
}

func findRegistryDocumentRelations(ctx context.Context, client *arshin.Client, registry registries.Spec, field string, query string, pageSize int, maxPages int) registryReport {
	result := registryReport{
		RegistryID: string(registry.RegistryID),
		Title:      registry.Title,
	}

	for pageNumber := 1; pageNumber <= maxPages; pageNumber++ {
		page, err := client.ListRegistryData(ctx, registry.RegistryID, arshin.RegistryQuery{
			PageNumber: pageNumber,
			PageSize:   pageSize,
			Filters: []arshin.Filter{
				{Field: field, Value: query},
			},
		})
		if err != nil {
			result.Error = err.Error()
			return result
		}
		if pageNumber == 1 {
			result.TotalCount = page.TotalCount
		}
		if len(page.Items) == 0 {
			break
		}

		for _, item := range page.Items {
			result.Records = append(result.Records, buildRecord(ctx, client, registry.RegistryID, item))
		}

		if pageNumber*pageSize >= page.TotalCount {
			break
		}
	}

	return result
}

func buildRecord(ctx context.Context, client *arshin.Client, registryID arshin.RegistryID, source arshin.RegistryRecord) record {
	result := record{
		Record: summarizeRecord(registryID, source.ID, source.Type, source.Properties),
		Method: summarizeMethodIfSupported(registryID, source),
	}

	detailedItem, err := client.GetRegistryItem(ctx, registryID, arshin.RegistryItemID(source.ID))
	if err != nil {
		result.Error = err.Error()
		return result
	}

	result.Record = summarizeRecord(registryID, source.ID, source.Type, detailedItem.Fields())
	result.Method = summarizeMethodItemIfSupported(registryID, *detailedItem)
	result.Relations = resolveRelations(ctx, client, detailedItem.Fields())
	return result
}

func resolveRelations(ctx context.Context, client *arshin.Client, properties []arshin.RegistryProperty) []relationSummary {
	var result []relationSummary
	for _, property := range properties {
		for _, ref := range property.ItemRefs() {
			result = append(result, resolveRef(ctx, client, property, "", "", ref))
		}
		if property.Type != "CHILD_OBJECT" {
			continue
		}

		children, err := property.ChildObjects()
		if err != nil {
			result = append(result, relationSummary{
				FieldName:  property.Name,
				FieldTitle: property.Title,
				Error:      err.Error(),
			})
			continue
		}
		for _, child := range children {
			fields := child.FieldsByName()
			relationType := stringProperty(fields, "foei:RelationTypeCMM")
			relatedNumber := stringProperty(fields, "foei:RelatedCMM")
			for _, ref := range child.ItemRefs() {
				result = append(result, resolveRef(ctx, client, property, relationType, relatedNumber, ref))
			}
		}
	}
	return result
}

func resolveRef(ctx context.Context, client *arshin.Client, property arshin.RegistryProperty, relationType string, relatedNumber string, ref arshin.RegistryItemRef) relationSummary {
	summary := relationSummary{
		FieldName:             property.Name,
		FieldTitle:            property.Title,
		Type:                  relationType,
		RelatedRegistryNumber: relatedNumber,
		RelatedRef: &itemRef{
			RegistryID: string(ref.RegistryID),
			ItemID:     string(ref.ItemID),
			Link:       ref.Link,
		},
	}

	relatedItem, err := client.GetRegistryItem(ctx, ref.RegistryID, ref.ItemID)
	if err != nil {
		summary.Error = err.Error()
		return summary
	}

	summary.Record = ptr(summarizeRecord(ref.RegistryID, string(ref.ItemID), "", relatedItem.Fields()))
	summary.Method = summarizeMethodItemIfSupported(ref.RegistryID, *relatedItem)
	return summary
}

func summarizeMethodIfSupported(registryID arshin.RegistryID, record arshin.RegistryRecord) *methodSummary {
	if !methods.IsSupportedRegistry(registryID) {
		return nil
	}
	return ptr(summarizeMethod(methods.MapRecord(record)))
}

func summarizeMethodItemIfSupported(registryID arshin.RegistryID, item arshin.RegistryItem) *methodSummary {
	if !methods.IsSupportedRegistry(registryID) {
		return nil
	}
	return ptr(summarizeMethod(methods.MapItem(item)))
}

func summarizeRecord(registryID arshin.RegistryID, id string, recordType string, properties []arshin.RegistryProperty) recordSummary {
	if recordType == "" {
		if spec, ok := registries.SpecForRegistry(registryID); ok {
			recordType = spec.ItemType
		}
	}

	summary := recordSummary{
		ID:         id,
		RegistryID: string(registryID),
		Type:       recordType,
		Properties: summarizeProperties(properties),
		Documents:  summarizePropertyAttachments(properties),
	}
	summary.DisplayName = firstStringProperty(properties, []string{
		"foei:NameCMM",
		"foei:NameND",
		"foei:NameGSI",
		"foei:NameMdg",
		"foei:NameMD",
		"foei:NameSI",
		"foei:NameSO",
		"foei:FullName",
		"foei:name",
		"cm:name",
	})
	return summary
}

func summarizeMethod(method *methods.Method) methodSummary {
	if method == nil {
		return methodSummary{}
	}

	return methodSummary{
		ID:                   method.ID,
		RegistryID:           string(method.RegistryID),
		RegistryNumber:       method.RegistryNumber,
		RecordNumber:         method.RecordNumber,
		Name:                 method.Name,
		Status:               method.Status,
		SystemStatus:         method.SystemStatus,
		PublishedAt:          method.PublishedAt,
		CertificateDate:      method.CertificateDate,
		CertificateNumber:    method.CertificateNumber,
		DeveloperName:        method.DeveloperName,
		CertificationOrgName: method.CertificationOrgName,
		OrderNumber:          method.OrderNumber,
		OrderDate:            method.OrderDate,
		OrderDocument:        method.OrderDocument,
		Documents:            summarizeDocuments(method),
	}
}

func summarizeDocuments(method *methods.Method) map[string][]document {
	documents := map[string][]document{
		"order":       summarizeAttachments(method.OrderDocuments),
		"text":        summarizeAttachments(method.TextDocuments),
		"range":       summarizeAttachments(method.RangeDocuments),
		"error":       summarizeAttachments(method.ErrorDocuments),
		"certificate": summarizeAttachments(method.CertificateDocuments),
		"methodData":  summarizeAttachments(method.MethodDataDocuments),
	}

	for name, attachments := range documents {
		if len(attachments) == 0 {
			delete(documents, name)
		}
	}
	if len(documents) == 0 {
		return nil
	}
	return documents
}

func summarizeAttachments(attachments []arshin.AttachmentRef) []document {
	result := make([]document, 0, len(attachments))
	for _, attachment := range attachments {
		result = append(result, document{
			FileID:   attachment.FileID,
			Name:     attachment.Name,
			MIME:     attachment.MIME,
			Download: attachment.Link,
		})
	}
	return result
}

func summarizeProperties(properties []arshin.RegistryProperty) []propertyValue {
	result := make([]propertyValue, 0, len(properties))
	for _, property := range properties {
		if property.Value == nil {
			continue
		}
		result = append(result, propertyValue{
			Name:  property.Name,
			Title: property.Title,
			Type:  property.Type,
			Value: property.Value,
		})
	}
	return result
}

func summarizePropertyAttachments(properties []arshin.RegistryProperty) []document {
	var result []document
	for _, property := range properties {
		result = append(result, summarizeAttachments(property.Attachments())...)
	}
	return result
}

func firstStringProperty(properties []arshin.RegistryProperty, names []string) string {
	index := arshin.PropertiesByName(properties)
	for _, name := range names {
		if value := stringProperty(index, name); value != "" {
			return value
		}
	}
	for _, property := range properties {
		if value, ok := property.Value.(string); ok && value != "" {
			return value
		}
	}
	return ""
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

func writeJSON(path string, value any) error {
	var out io.Writer = os.Stdout
	var file *os.File
	var err error

	if path != "" {
		file, err = os.Create(path)
		if err != nil {
			return err
		}
		defer file.Close()
		out = file
	}

	encoder := json.NewEncoder(out)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	return encoder.Encode(value)
}

func ptr[T any](value T) *T {
	return &value
}

func exit(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
