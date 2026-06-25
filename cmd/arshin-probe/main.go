package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"

	arshin "github.com/ReanSn0w/go-fgis-private-arshin/public/client"
	"github.com/ReanSn0w/go-fgis-private-arshin/public/registries"
)

type dataBehaviorReport struct {
	GeneratedAt time.Time           `json:"generatedAt"`
	Checks      []dataBehaviorCheck `json:"checks"`
}

type dataBehaviorCheck struct {
	Name        string   `json:"name"`
	RegistryID  string   `json:"registryId"`
	OK          bool     `json:"ok"`
	Error       string   `json:"error,omitempty"`
	TotalCount  int      `json:"totalCount,omitempty"`
	CurrentPage int      `json:"currentPage,omitempty"`
	PageSize    int      `json:"pageSize,omitempty"`
	ItemCount   int      `json:"itemCount,omitempty"`
	RecordTypes []string `json:"recordTypes,omitempty"`
	FirstID     string   `json:"firstId,omitempty"`
	FirstNumber string   `json:"firstNumber,omitempty"`
	ValuesNil   *bool    `json:"valuesNil,omitempty"`
	Note        string   `json:"note,omitempty"`
}

type sessionSecurityReport struct {
	GeneratedAt time.Time              `json:"generatedAt"`
	Checks      []sessionSecurityCheck `json:"checks"`
}

type sessionSecurityCheck struct {
	Name             string   `json:"name"`
	Method           string   `json:"method"`
	URL              string   `json:"url"`
	StatusCode       int      `json:"statusCode,omitempty"`
	ContentType      string   `json:"contentType,omitempty"`
	ContentLength    int64    `json:"contentLength,omitempty"`
	BodyPrefix       string   `json:"bodyPrefix,omitempty"`
	SetCookieNames   []string `json:"setCookieNames,omitempty"`
	CookieNamesAfter []string `json:"cookieNamesAfter,omitempty"`
	AjaxTokenName    string   `json:"ajaxTokenName,omitempty"`
	LooksLikeJSON    bool     `json:"looksLikeJson,omitempty"`
	LooksLikeHTML    bool     `json:"looksLikeHtml,omitempty"`
	Error            string   `json:"error,omitempty"`
}

type relationSummary struct {
	FieldName    string                   `json:"fieldName"`
	FieldTitle   string                   `json:"fieldTitle"`
	FieldType    string                   `json:"fieldType"`
	Value        any                      `json:"value,omitempty"`
	Refs         []arshin.RegistryItemRef `json:"refs,omitempty"`
	ChildObjects []childObjectSummary     `json:"childObjects,omitempty"`
	Error        string                   `json:"error,omitempty"`
}

type childObjectSummary struct {
	ID      string                   `json:"id"`
	Type    string                   `json:"type"`
	NodeRef string                   `json:"nodeRef,omitempty"`
	Refs    []arshin.RegistryItemRef `json:"refs,omitempty"`
	Fields  []propertySummary        `json:"fields,omitempty"`
}

type registrySurvey struct {
	RegistryID  string            `json:"registryId"`
	Title       string            `json:"title"`
	OK          bool              `json:"ok"`
	Error       string            `json:"error,omitempty"`
	TotalCount  int               `json:"totalCount,omitempty"`
	CurrentPage int               `json:"currentPage,omitempty"`
	PageSize    int               `json:"pageSize,omitempty"`
	RecordTypes []string          `json:"recordTypes,omitempty"`
	Properties  []propertySummary `json:"properties,omitempty"`
}

type registryMetadata struct {
	RegistryID    string                    `json:"registryId"`
	Title         string                    `json:"title,omitempty"`
	OK            bool                      `json:"ok"`
	Error         string                    `json:"error,omitempty"`
	Details       *arshin.RegistryDetails   `json:"details,omitempty"`
	Fields        []arshin.RegistryProperty `json:"fields,omitempty"`
	DisplayFields []arshin.RegistryProperty `json:"displayFields,omitempty"`
}

type metadataDump struct {
	FieldTypes []arshin.FieldType `json:"fieldTypes,omitempty"`
	Registries []registryMetadata `json:"registries"`
}

type propertySummary struct {
	Name     string `json:"name"`
	Title    string `json:"title"`
	Type     string `json:"type"`
	Multiple bool   `json:"multiple"`
	HasLink  bool   `json:"hasLink,omitempty"`
	HasMIME  bool   `json:"hasMime,omitempty"`
}

func main() {
	registryID := flag.String("registry", "16", "registry id")
	pageNumber := flag.Int("page", 1, "page number")
	pageSize := flag.Int("size", 5, "page size")
	itemID := flag.String("item", "", "registry item id")
	plain := flag.Bool("plain", false, "with -item, fetch /plaindata instead of sectioned /data")
	relations := flag.Bool("relations", false, "with -item, print link and child-object relations")
	download := flag.String("download", "", "download file id or /api/downloadfile/... link")
	filterBy := flag.String("filter-by", "", "filter field name")
	filterValue := flag.String("filter-value", "", "filter value")
	printProperties := flag.Bool("properties", false, "print unique properties instead of raw JSON")
	all := flag.Bool("all", false, "survey all known registries")
	metadata := flag.Bool("metadata", false, "fetch registry metadata instead of data")
	behavior := flag.Bool("behavior", false, "run generic /data behavior checks")
	session := flag.Bool("session", false, "run session/header/cookie behavior checks")
	out := flag.String("out", "", "write JSON output to file")
	flag.Parse()

	client, err := arshin.NewClient()
	if err != nil {
		exit(err)
	}

	query := arshin.RegistryQuery{
		PageNumber: *pageNumber,
		PageSize:   *pageSize,
	}
	if *filterBy != "" {
		query.Filters = append(query.Filters, arshin.Filter{
			Field: *filterBy,
			Value: *filterValue,
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var output any
	if *download != "" {
		var file *arshin.DownloadedFile
		if strings.Contains(*download, "/") {
			file, err = client.DownloadFileLink(ctx, *download)
		} else {
			file, err = client.DownloadFile(ctx, *download)
		}
		if err != nil {
			exit(err)
		}
		defer file.Body.Close()

		if *out == "" {
			if file.Filename != "" {
				*out = file.Filename
			} else {
				*out = "download.bin"
			}
		}

		if err := writeStream(*out, file.Body); err != nil {
			exit(err)
		}
		fmt.Fprintf(os.Stderr, "saved %s (%s, %d bytes)\n", *out, file.ContentType, file.ContentLength)
		return
	}

	if *itemID != "" {
		if *relations {
			item, err := client.GetRegistryItem(ctx, arshin.RegistryID(*registryID), arshin.RegistryItemID(*itemID))
			if err != nil {
				exit(err)
			}
			output = summarizeRelations(*item)
			writeJSON(*out, output)
			return
		}

		if *plain {
			output, err = client.GetRegistryItemPlainData(ctx, arshin.RegistryID(*registryID), arshin.RegistryItemID(*itemID))
		} else {
			output, err = client.GetRegistryItem(ctx, arshin.RegistryID(*registryID), arshin.RegistryItemID(*itemID))
		}
		if err != nil {
			exit(err)
		}
		writeJSON(*out, output)
		return
	}

	if *behavior {
		output = runDataBehaviorChecks(ctx, client)
		writeJSON(*out, output)
		return
	}

	if *session {
		output = runSessionSecurityChecks(ctx)
		writeJSON(*out, output)
		return
	}

	if *metadata {
		if *all {
			fieldTypes, err := client.ListFieldTypes(ctx)
			if err != nil {
				exit(err)
			}
			output = metadataDump{
				FieldTypes: fieldTypes,
				Registries: fetchKnownRegistriesMetadata(ctx, client),
			}
		} else {
			output = fetchRegistryMetadata(ctx, client, arshin.RegistryID(*registryID), "")
		}
		writeJSON(*out, output)
		return
	}

	if *all {
		output = surveyKnownRegistries(ctx, client, *pageSize)
		writeJSON(*out, output)
		return
	}

	page, err := client.ListRegistryData(ctx, arshin.RegistryID(*registryID), query)
	if err != nil {
		exit(err)
	}

	if *printProperties {
		printPropertySummary(os.Stdout, page.Items)
		return
	}

	output = page
	writeJSON(*out, output)
}

func surveyKnownRegistries(ctx context.Context, client *arshin.Client, pageSize int) []registrySurvey {
	knownRegistries := registries.Known()
	result := make([]registrySurvey, 0, len(knownRegistries))

	for _, registry := range knownRegistries {
		item := registrySurvey{
			RegistryID: string(registry.RegistryID),
			Title:      registry.Title,
		}

		page, err := client.ListRegistryData(ctx, registry.RegistryID, arshin.RegistryQuery{
			PageNumber: 1,
			PageSize:   pageSize,
		})
		if err != nil {
			item.Error = err.Error()
			result = append(result, item)
			continue
		}

		item.OK = true
		item.TotalCount = page.TotalCount
		item.CurrentPage = page.CurrentPage
		item.PageSize = page.PageSize
		item.RecordTypes = uniqueRecordTypes(page.Items)
		item.Properties = summarizeProperties(page.Items)
		result = append(result, item)
	}

	return result
}

func summarizeRelations(item arshin.RegistryItem) []relationSummary {
	var result []relationSummary
	for _, field := range item.Fields() {
		if field.Type != "LINK_INTERNAL" && field.Type != "CHILD_OBJECT" {
			continue
		}

		summary := relationSummary{
			FieldName:  field.Name,
			FieldTitle: field.Title,
			FieldType:  field.Type,
			Value:      field.Value,
			Refs:       field.ItemRefs(),
		}

		if field.Type == "CHILD_OBJECT" {
			children, err := field.ChildObjects()
			if err != nil {
				summary.Error = err.Error()
			} else {
				summary.ChildObjects = summarizeChildObjects(children)
			}
		}

		result = append(result, summary)
	}
	return result
}

func summarizeChildObjects(children []arshin.ChildObject) []childObjectSummary {
	result := make([]childObjectSummary, 0, len(children))
	for _, child := range children {
		result = append(result, childObjectSummary{
			ID:      child.ID,
			Type:    child.Type,
			NodeRef: child.NodeRef,
			Refs:    child.ItemRefs(),
			Fields:  summarizeProperties([]arshin.RegistryRecord{{Properties: child.Fields}}),
		})
	}
	return result
}

func runDataBehaviorChecks(ctx context.Context, client *arshin.Client) dataBehaviorReport {
	checks := []dataBehaviorCheck{
		runDataCheck(ctx, client, "first_page_registry_16", "16", arshin.RegistryQuery{PageNumber: 1, PageSize: 3}),
		runDataCheck(ctx, client, "pagination_registry_16_page_2", "16", arshin.RegistryQuery{PageNumber: 2, PageSize: 3}),
		runDataCheck(ctx, client, "empty_filter_registry_16", "16", arshin.RegistryQuery{
			PageNumber: 1,
			PageSize:   3,
			Filters: []arshin.Filter{
				{Field: "foei:NumRegCMM", Value: "__missing_value__"},
			},
		}),
		runDataCheck(ctx, client, "single_filter_registry_16", "16", arshin.RegistryQuery{
			PageNumber: 1,
			PageSize:   3,
			Filters: []arshin.Filter{
				{Field: "foei:NumRegCMM", Value: "ФР.1.31.2022.44733"},
			},
		}),
		runDataCheck(ctx, client, "multiple_filters_registry_16", "16", arshin.RegistryQuery{
			PageNumber: 1,
			PageSize:   3,
			Filters: []arshin.Filter{
				{Field: "foei:NumRegCMM", Value: "ФР.1.31.2022.44733"},
				{Field: "foei:StatusCMM", Value: "Действует"},
			},
		}),
		runDataCheck(ctx, client, "sort_asc_registry_16", "16", arshin.RegistryQuery{
			PageNumber: 1,
			PageSize:   3,
			Sorts: []arshin.Sort{
				{Field: "foei:number", Direction: "asc"},
			},
		}),
		runDataCheck(ctx, client, "sort_desc_registry_16", "16", arshin.RegistryQuery{
			PageNumber: 1,
			PageSize:   3,
			Sorts: []arshin.Sort{
				{Field: "foei:number", Direction: "desc"},
			},
		}),
		runDataCheck(ctx, client, "larger_page_size_registry_16", "16", arshin.RegistryQuery{PageNumber: 1, PageSize: 100}),
		runDataCheck(ctx, client, "large_page_size_500_registry_16", "16", arshin.RegistryQuery{PageNumber: 1, PageSize: 500}),
		runDataCheck(ctx, client, "small_registry_18", "18", arshin.RegistryQuery{PageNumber: 1, PageSize: 10}),
		runDataCheck(ctx, client, "largest_registry_11", "11", arshin.RegistryQuery{PageNumber: 1, PageSize: 3}),
		runDataCheck(ctx, client, "registry_47_namespace", "47", arshin.RegistryQuery{PageNumber: 1, PageSize: 3}),
		runDataCheck(ctx, client, "invalid_registry", "999999", arshin.RegistryQuery{PageNumber: 1, PageSize: 3}),
	}

	for i := range checks {
		switch checks[i].Name {
		case "invalid_registry":
			if !checks[i].OK {
				checks[i].OK = true
				checks[i].Note = "Expected failure for invalid registry id."
			}
		case "empty_filter_registry_16":
			checks[i].OK = checks[i].OK && checks[i].TotalCount == 0 && checks[i].ItemCount == 0
		case "single_filter_registry_16", "multiple_filters_registry_16":
			checks[i].OK = checks[i].OK && checks[i].TotalCount == 1 && checks[i].FirstID == "1404476"
		}
	}

	return dataBehaviorReport{
		GeneratedAt: time.Now().UTC(),
		Checks:      checks,
	}
}

func runDataCheck(ctx context.Context, client *arshin.Client, name string, registryID arshin.RegistryID, query arshin.RegistryQuery) dataBehaviorCheck {
	check := dataBehaviorCheck{
		Name:       name,
		RegistryID: string(registryID),
	}

	page, err := client.ListRegistryData(ctx, registryID, query)
	if err != nil {
		check.Error = err.Error()
		return check
	}

	check.OK = true
	check.TotalCount = page.TotalCount
	check.CurrentPage = page.CurrentPage
	check.PageSize = page.PageSize
	check.ItemCount = len(page.Items)
	check.RecordTypes = uniqueRecordTypes(page.Items)

	if len(page.Items) > 0 {
		first := page.Items[0]
		check.FirstID = first.ID
		check.FirstNumber = propertyString(first, "foei:number")
		valuesNil := first.Values == nil
		check.ValuesNil = &valuesNil
	}

	return check
}

func runSessionSecurityChecks(ctx context.Context) sessionSecurityReport {
	jar, _ := cookiejar.New(nil)
	jarClient := &http.Client{
		Timeout: 30 * time.Second,
		Jar:     jar,
	}
	plainClient := &http.Client{Timeout: 30 * time.Second}

	registryPage := "https://fgis.gost.ru/fundmetrology/registry/16"
	dataURL := "https://fgis.gost.ru/fundmetrology/api/registry/16/data?pageNumber=1&pageSize=1&orgID=CURRENT_ORG&filterBy=foei%3ANumRegCMM&filterValues=%D0%A4%D0%A0.1.31.2022.44733"
	itemURL := "https://fgis.gost.ru/fundmetrology/api/registry/16/items/1404476/data"
	downloadURL := "https://fgis.gost.ru/fundmetrology/api/downloadfile/c11f089b-22f4-4d66-87cd-7aff99909001"
	invalidURL := "https://fgis.gost.ru/fundmetrology/api/registry/999999/data?pageNumber=1&pageSize=1&orgID=CURRENT_ORG"
	htmlAsJSONURL := "https://fgis.gost.ru/fundmetrology/registry/16"

	var checks []sessionSecurityCheck

	checks = append(checks, runHTTPCheck(ctx, plainClient, "data_plain_no_browser_headers_no_cookie", dataURL, nil))
	checks = append(checks, runHTTPCheck(ctx, plainClient, "data_browser_headers_no_cookie", dataURL, browserHeaders(registryPage, "")))

	pageCheck := runHTTPCheck(ctx, jarClient, "registry_page_preflight_sets_cookie_and_token", registryPage, map[string]string{
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
		"Accept-Language": "ru-RU,ru;q=0.9",
		"User-Agent":      arshin.DefaultUserAgent,
	})
	tokenName := tokenNameFromHTML(pageCheck.BodyPrefix)
	pageCheck.AjaxTokenName = tokenName
	checks = append(checks, pageCheck)

	checks = append(checks, runHTTPCheck(ctx, jarClient, "data_after_preflight_cookie", dataURL, browserHeaders(registryPage, "")))
	if tokenName != "" {
		checks = append(checks, runHTTPCheck(ctx, jarClient, "data_after_preflight_cookie_with_ajax_token_header", dataURL, browserHeaders(registryPage, tokenName)))
	}
	checks = append(checks, runHTTPCheck(ctx, jarClient, "item_after_preflight_cookie", itemURL, browserHeaders(registryPage, "")))
	checks = append(checks, runHTTPCheck(ctx, jarClient, "download_after_preflight_cookie", downloadURL, browserHeaders(registryPage, "")))
	checks = append(checks, runHTTPCheck(ctx, plainClient, "invalid_registry_json_error", invalidURL, browserHeaders(registryPage, "")))
	checks = append(checks, runHTTPCheck(ctx, plainClient, "html_page_when_json_expected", htmlAsJSONURL, map[string]string{
		"Accept":          "application/json, text/plain, */*",
		"Accept-Language": "ru-RU,ru;q=0.9",
		"User-Agent":      arshin.DefaultUserAgent,
	}))

	return sessionSecurityReport{
		GeneratedAt: time.Now().UTC(),
		Checks:      checks,
	}
}

func runHTTPCheck(ctx context.Context, client *http.Client, name string, rawURL string, headers map[string]string) sessionSecurityCheck {
	check := sessionSecurityCheck{
		Name:   name,
		Method: http.MethodGet,
		URL:    rawURL,
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		check.Error = err.Error()
		return check
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		check.Error = err.Error()
		return check
	}
	defer resp.Body.Close()

	check.StatusCode = resp.StatusCode
	check.ContentType = resp.Header.Get("Content-Type")
	check.ContentLength = resp.ContentLength
	check.SetCookieNames = cookieNames(resp.Cookies())

	body, err := io.ReadAll(io.LimitReader(resp.Body, 512))
	if err != nil {
		check.Error = err.Error()
		return check
	}
	if shouldKeepBodyPrefix(check.ContentType) {
		check.BodyPrefix = string(body)
	}
	trimmed := strings.TrimSpace(check.BodyPrefix)
	check.LooksLikeJSON = strings.HasPrefix(trimmed, "{") || strings.HasPrefix(trimmed, "[")
	check.LooksLikeHTML = strings.HasPrefix(strings.ToLower(trimmed), "<!doctype") || strings.HasPrefix(strings.ToLower(trimmed), "<html")

	if client.Jar != nil {
		if parsed, err := url.Parse(rawURL); err == nil {
			check.CookieNamesAfter = cookieNames(client.Jar.Cookies(parsed))
		}
	}

	if strings.Contains(strings.ToLower(check.ContentType), "html") {
		check.AjaxTokenName = tokenNameFromHTML(check.BodyPrefix)
	}

	return check
}

func shouldKeepBodyPrefix(contentType string) bool {
	contentType = strings.ToLower(contentType)
	return strings.Contains(contentType, "json") ||
		strings.Contains(contentType, "html") ||
		strings.Contains(contentType, "text")
}

func browserHeaders(referer string, ajaxTokenName string) map[string]string {
	headers := map[string]string{
		"Accept":          "application/json, text/plain, */*",
		"Accept-Language": "ru-RU,ru;q=0.9",
		"Referer":         referer,
		"User-Agent":      arshin.DefaultUserAgent,
	}
	if ajaxTokenName != "" {
		headers[ajaxTokenName] = "probe"
	}
	return headers
}

func cookieNames(cookies []*http.Cookie) []string {
	names := make([]string, 0, len(cookies))
	for _, cookie := range cookies {
		names = append(names, cookie.Name)
	}
	sort.Strings(names)
	return names
}

func tokenNameFromHTML(html string) string {
	const marker = "name='hmac-token-name' content='"
	start := strings.Index(html, marker)
	if start == -1 {
		return ""
	}
	start += len(marker)
	end := strings.Index(html[start:], "'")
	if end == -1 {
		return ""
	}
	return html[start : start+end]
}

func fetchKnownRegistriesMetadata(ctx context.Context, client *arshin.Client) []registryMetadata {
	knownRegistries := registries.Known()
	result := make([]registryMetadata, 0, len(knownRegistries))

	for _, registry := range knownRegistries {
		result = append(result, fetchRegistryMetadata(ctx, client, registry.RegistryID, registry.Title))
	}

	return result
}

func fetchRegistryMetadata(ctx context.Context, client *arshin.Client, registryID arshin.RegistryID, title string) registryMetadata {
	item := registryMetadata{
		RegistryID: string(registryID),
		Title:      title,
	}

	details, err := client.GetRegistry(ctx, registryID)
	if err != nil {
		item.Error = err.Error()
		return item
	}
	item.Details = details
	item.OK = true
	if item.Title == "" {
		item.Title = details.Title
	}

	fields, err := client.ListRegistryFields(ctx, registryID)
	if err != nil {
		item.OK = false
		item.Error = err.Error()
		return item
	}
	item.Fields = fields

	displayFields, err := client.ListRegistryDisplayFields(ctx, registryID)
	if err != nil {
		item.OK = false
		item.Error = err.Error()
		return item
	}
	item.DisplayFields = displayFields

	return item
}

func uniqueRecordTypes(records []arshin.RegistryRecord) []string {
	seen := map[string]struct{}{}
	for _, record := range records {
		if record.Type == "" {
			continue
		}
		seen[record.Type] = struct{}{}
	}

	types := make([]string, 0, len(seen))
	for recordType := range seen {
		types = append(types, recordType)
	}
	sort.Strings(types)
	return types
}

func summarizeProperties(records []arshin.RegistryRecord) []propertySummary {
	seen := map[string]propertySummary{}
	for _, record := range records {
		for _, property := range record.Properties {
			item := seen[property.Name]
			if item.Name == "" {
				item = propertySummary{
					Name:     property.Name,
					Title:    property.Title,
					Type:     property.Type,
					Multiple: property.Multiple,
				}
			}
			if property.Link != nil {
				item.HasLink = true
			}
			if property.MIME != nil {
				item.HasMIME = true
			}
			seen[property.Name] = item
		}
	}

	names := make([]string, 0, len(seen))
	for name := range seen {
		names = append(names, name)
	}
	sort.Strings(names)

	properties := make([]propertySummary, 0, len(names))
	for _, name := range names {
		properties = append(properties, seen[name])
	}
	return properties
}

func propertyString(record arshin.RegistryRecord, name string) string {
	property, ok := record.PropertiesByName()[name]
	if !ok || property.Value == nil {
		return ""
	}
	value, ok := property.Value.(string)
	if ok {
		return value
	}
	return fmt.Sprint(property.Value)
}

func printPropertySummary(w io.Writer, records []arshin.RegistryRecord) {
	for _, item := range summarizeProperties(records) {
		fmt.Fprintf(w, "%s\t%s\t%s\t%t\n", item.Name, item.Title, item.Type, item.Multiple)
	}
}

func writeJSON(path string, value any) {
	var w io.Writer = os.Stdout
	var file *os.File
	var err error

	if path != "" {
		file, err = os.Create(path)
		if err != nil {
			exit(err)
		}
		defer file.Close()
		w = file
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(value); err != nil {
		exit(err)
	}
}

func writeStream(path string, r io.Reader) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, r)
	return err
}

func exit(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
