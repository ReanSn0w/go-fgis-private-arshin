package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestListRegistryDataEncodesQuery(t *testing.T) {
	var gotPath string
	var gotQuery map[string][]string
	var gotReferer string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotQuery = r.URL.Query()
		gotReferer = r.Header.Get("Referer")

		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		_ = json.NewEncoder(w).Encode(RegistryDataResponse{
			Status: http.StatusOK,
			Result: RegistryDataPage{
				TotalCount:  1,
				CurrentPage: 2,
				PageSize:    10,
				Items: []RegistryRecord{
					{
						ID:   "record-1",
						Type: "foei:TEST_type",
						Properties: []RegistryProperty{
							{Name: "foei:number", Type: "text", Title: "Номер записи", Value: "1"},
						},
					},
				},
			},
		})
	}))
	defer server.Close()

	client, err := NewClient(
		WithBaseURL(server.URL+"/fundmetrology/api/"),
		WithPublicURL(server.URL+"/fundmetrology/"),
		WithRateLimit(0),
	)
	if err != nil {
		t.Fatal(err)
	}

	page, err := client.ListRegistryData(context.Background(), "16", RegistryQuery{
		PageNumber: 2,
		PageSize:   10,
		OrgID:      "ORG",
		Filters: []Filter{
			{Field: "foei:NumRegCMM", Value: "ФР.1"},
			{Field: "foei:status", Value: "Опубликована"},
		},
		Sorts: []Sort{
			{Field: "foei:number", Direction: "asc"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if gotPath != "/fundmetrology/api/registry/16/data" {
		t.Fatalf("path = %q", gotPath)
	}

	assertQueryValues(t, gotQuery, "pageNumber", []string{"2"})
	assertQueryValues(t, gotQuery, "pageSize", []string{"10"})
	assertQueryValues(t, gotQuery, "orgID", []string{"ORG"})
	assertQueryValues(t, gotQuery, "filterBy", []string{"foei:NumRegCMM", "foei:status"})
	assertQueryValues(t, gotQuery, "filterValues", []string{"ФР.1", "Опубликована"})
	assertQueryValues(t, gotQuery, "sortBy", []string{"foei:number"})
	assertQueryValues(t, gotQuery, "sortValues", []string{"asc"})

	if gotReferer != server.URL+"/fundmetrology/registry/16" {
		t.Fatalf("referer = %q", gotReferer)
	}

	if page.TotalCount != 1 || page.CurrentPage != 2 || page.PageSize != 10 {
		t.Fatalf("unexpected page: %+v", page)
	}
	if len(page.Items) != 1 || page.Items[0].ID != "record-1" {
		t.Fatalf("unexpected items: %+v", page.Items)
	}
}

func TestListRegistryDataDefaults(t *testing.T) {
	var gotQuery map[string][]string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.Query()
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(RegistryDataResponse{
			Status: http.StatusOK,
			Result: RegistryDataPage{},
		})
	}))
	defer server.Close()

	client, err := NewClient(
		WithBaseURL(server.URL+"/"),
		WithRateLimit(0),
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.ListRegistryData(context.Background(), "1", RegistryQuery{})
	if err != nil {
		t.Fatal(err)
	}

	assertQueryValues(t, gotQuery, "pageNumber", []string{"1"})
	assertQueryValues(t, gotQuery, "pageSize", []string{"20"})
	assertQueryValues(t, gotQuery, "orgID", []string{DefaultOrgID})
}

func TestRegistryMetadataEndpoints(t *testing.T) {
	var gotPaths []string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPaths = append(gotPaths, r.URL.Path)
		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Path {
		case "/api/registry/16":
			_ = json.NewEncoder(w).Encode(registryDetailsResponse{
				Status: http.StatusOK,
				Result: RegistryDetails{
					ID:         "16",
					Title:      "Аттестованные методики (методы) измерений",
					Type:       "foei:CMM1_type",
					AlfrescoID: "alfresco-id",
				},
			})
		case "/api/registry/16/fields":
			_ = json.NewEncoder(w).Encode(registryFieldsResponse{
				Status: http.StatusOK,
				Result: []RegistryProperty{
					{Name: "foei:NumRegCMM", Type: "text", Title: "Номер в реестре"},
				},
			})
		case "/api/registry/16/displayfields":
			_ = json.NewEncoder(w).Encode(registryFieldsResponse{
				Status: http.StatusOK,
				Result: []RegistryProperty{
					{Name: "foei:NameCMM", Type: "text", Title: "Наименование документа на методику"},
				},
			})
		default:
			t.Fatalf("unexpected path %q", r.URL.Path)
		}
	}))
	defer server.Close()

	client, err := NewClient(
		WithBaseURL(server.URL+"/api/"),
		WithRateLimit(0),
	)
	if err != nil {
		t.Fatal(err)
	}

	details, err := client.GetRegistry(context.Background(), "16")
	if err != nil {
		t.Fatal(err)
	}
	if details.Type != "foei:CMM1_type" {
		t.Fatalf("details.Type = %q", details.Type)
	}

	fields, err := client.ListRegistryFields(context.Background(), "16")
	if err != nil {
		t.Fatal(err)
	}
	if len(fields) != 1 || fields[0].Name != "foei:NumRegCMM" {
		t.Fatalf("fields = %+v", fields)
	}

	displayFields, err := client.ListRegistryDisplayFields(context.Background(), "16")
	if err != nil {
		t.Fatal(err)
	}
	if len(displayFields) != 1 || displayFields[0].Name != "foei:NameCMM" {
		t.Fatalf("displayFields = %+v", displayFields)
	}

	wantPaths := []string{"/api/registry/16", "/api/registry/16/fields", "/api/registry/16/displayfields"}
	if !reflect.DeepEqual(gotPaths, wantPaths) {
		t.Fatalf("paths = %#v, want %#v", gotPaths, wantPaths)
	}
}

func TestListRegistryDataRejectsUnexpectedContentType(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("<html></html>"))
	}))
	defer server.Close()

	client, err := NewClient(
		WithBaseURL(server.URL+"/"),
		WithRateLimit(0),
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.ListRegistryData(context.Background(), "1", RegistryQuery{})
	if err == nil {
		t.Fatal("expected error")
	}
	if _, ok := err.(*UnexpectedContentTypeError); !ok {
		t.Fatalf("error = %T %v", err, err)
	}
}

func TestListRegistryDataParsesJSONAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(APIError{
			Status:  http.StatusInternalServerError,
			Message: "registry exploded politely",
			Trace:   "trace",
		})
	}))
	defer server.Close()

	client, err := NewClient(
		WithBaseURL(server.URL+"/"),
		WithRateLimit(0),
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.ListRegistryData(context.Background(), "1", RegistryQuery{})
	if err == nil {
		t.Fatal("expected error")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("error = %T %v", err, err)
	}
	if apiErr.Message != "registry exploded politely" {
		t.Fatalf("message = %q", apiErr.Message)
	}
}

func TestRateLimiterDisabledDoesNotWait(t *testing.T) {
	limiter := NewRateLimiter(0)

	started := time.Now()
	if err := limiter.Wait(context.Background()); err != nil {
		t.Fatal(err)
	}
	if time.Since(started) > 50*time.Millisecond {
		t.Fatal("disabled rate limiter waited")
	}
}

func assertQueryValues(t *testing.T, query map[string][]string, key string, want []string) {
	t.Helper()
	got := query[key]
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("query[%s] = %#v, want %#v", key, got, want)
	}
}
