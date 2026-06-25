package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestRegistryItemEndpoints(t *testing.T) {
	var gotPaths []string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPaths = append(gotPaths, r.URL.Path)
		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Path {
		case "/api/registry/16/items/1404476/data":
			_ = json.NewEncoder(w).Encode(registryItemResponse{
				Status: http.StatusOK,
				Result: RegistryItem{
					ID:         1404476,
					AlfrescoID: "alfresco-id",
					RegistryID: 16,
					Sections: []RegistryItemSection{
						{
							SectionTitle: "Сведения",
							Fields: []RegistryProperty{
								{Name: "foei:NumRegCMM", Type: "text", Title: "Номер в реестре", Value: "ФР.1"},
							},
						},
					},
				},
			})
		case "/api/registry/16/items/1404476/plaindata":
			_ = json.NewEncoder(w).Encode(registryPlainDataResponse{
				Status: http.StatusOK,
				Result: RegistryRecord{
					ID:   "1404476",
					Type: "foei:CMM1_type",
					Properties: []RegistryProperty{
						{Name: "foei:NumRegCMM", Type: "text", Title: "Номер в реестре", Value: "ФР.1"},
					},
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

	item, err := client.GetRegistryItem(context.Background(), "16", "1404476")
	if err != nil {
		t.Fatal(err)
	}
	if item.ID != 1404476 {
		t.Fatalf("item.ID = %d", item.ID)
	}
	if len(item.Sections) != 1 || len(item.Fields()) != 1 {
		t.Fatalf("item = %+v", item)
	}
	if item.FieldsByName()["foei:NumRegCMM"].Value != "ФР.1" {
		t.Fatalf("fields = %+v", item.FieldsByName())
	}

	record, err := client.GetRegistryItemPlainData(context.Background(), "16", "1404476")
	if err != nil {
		t.Fatal(err)
	}
	if record.ID != "1404476" || record.Type != "foei:CMM1_type" {
		t.Fatalf("record = %+v", record)
	}

	wantPaths := []string{
		"/api/registry/16/items/1404476/data",
		"/api/registry/16/items/1404476/plaindata",
	}
	if !reflect.DeepEqual(gotPaths, wantPaths) {
		t.Fatalf("paths = %#v, want %#v", gotPaths, wantPaths)
	}
}
