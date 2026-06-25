package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListFieldTypes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/fieldtypes" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(fieldTypesResponse{
			Status: http.StatusOK,
			Result: []FieldType{
				{Title: "ATTACH", Attach: true},
				{Title: "LINK_INTERNAL", Link: true},
			},
		})
	}))
	defer server.Close()

	client, err := NewClient(
		WithBaseURL(server.URL+"/api/"),
		WithRateLimit(0),
	)
	if err != nil {
		t.Fatal(err)
	}

	fieldTypes, err := client.ListFieldTypes(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(fieldTypes) != 2 {
		t.Fatalf("fieldTypes = %d", len(fieldTypes))
	}
	if !fieldTypes[0].Attach {
		t.Fatalf("fieldTypes[0] = %+v", fieldTypes[0])
	}
	if !fieldTypes[1].Link {
		t.Fatalf("fieldTypes[1] = %+v", fieldTypes[1])
	}
}
