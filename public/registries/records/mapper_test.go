package records

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	arshin "github.com/ReanSn0w/go-fgis-private-arshin/public/client"
	"github.com/ReanSn0w/go-fgis-private-arshin/public/registries"
)

func TestMapNormativeDocument(t *testing.T) {
	record := arshin.RegistryRecord{
		ID:   "1",
		Type: registries.ND.ItemType,
		Properties: []arshin.RegistryProperty{
			{Name: "foei:TypeND", Value: "Приказ"},
			{Name: "foei:NumberND", Value: "123"},
			{Name: "foei:DataND", Value: "01.02.2026"},
			{Name: "foei:OrgND", Value: "Минпромторг"},
			{Name: "foei:NameND", Value: "О проверке"},
			{Name: "foei:EditionND", Value: "ред. 1"},
			{Name: "foei:StatusND", Value: "Действует"},
		},
	}

	got, ok := MapRecord(registries.NDRegistryID, record).(NormativeDocument)
	if !ok {
		t.Fatalf("type = %T", MapRecord(registries.NDRegistryID, record))
	}
	if got.Number != "123" || got.Name != "О проверке" || got.Status != "Действует" {
		t.Fatalf("mapped normative document = %+v", got)
	}
	if got.RawProperties["foei:NumberND"].Value != "123" {
		t.Fatalf("raw properties = %+v", got.RawProperties)
	}
}

func TestTypedClientListReturnsConcreteRecords(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/registry/11/data" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(arshin.RegistryDataResponse{
			Status: http.StatusOK,
			Result: arshin.RegistryDataPage{
				TotalCount:  1,
				CurrentPage: 1,
				PageSize:    20,
				Items: []arshin.RegistryRecord{
					{
						ID:   "su-1",
						Type: registries.SU.ItemType,
						Properties: []arshin.RegistryProperty{
							{Name: "foei:NumRegSU", Value: "ГЭТ 1"},
							{Name: "foei:NameSU", Value: "Эталон единицы величины"},
						},
					},
				},
			},
		})
	}))
	defer server.Close()

	client, err := arshin.NewClient(
		arshin.WithBaseURL(server.URL+"/api/"),
		arshin.WithRateLimit(0),
	)
	if err != nil {
		t.Fatal(err)
	}

	typedClient := NewTypedClient(client, StandardUnits)
	page, err := typedClient.List(context.Background(), arshin.RegistryQuery{})
	if err != nil {
		t.Fatal(err)
	}

	if len(page.Items) != 1 {
		t.Fatalf("items = %d", len(page.Items))
	}
	unit := page.Items[0]
	if unit.RegistryNumber != "ГЭТ 1" || unit.Name != "Эталон единицы величины" {
		t.Fatalf("unit = %+v", unit)
	}
}

func TestMapInternationalDocumentAttachments(t *testing.T) {
	record := arshin.RegistryRecord{
		ID:   "2",
		Type: registries.MD.ItemType,
		Properties: []arshin.RegistryProperty{
			{Name: "foei:NumberMD", Value: "OIML R 1"},
			{Name: "foei:NameRusMD", Value: "Документ"},
			{Name: "foei:DocRusMD", Type: "ATTACH", Value: "ru.pdf", Link: "/fundmetrology/api/downloadfile/rus", MIME: "application/pdf"},
			{Name: "foei:DocEngMD", Type: "ATTACH", Value: "en.pdf", Link: "/fundmetrology/api/downloadfile/eng", MIME: "application/pdf"},
		},
	}

	got, ok := MapRecord(registries.MDRegistryID, record).(InternationalDocument)
	if !ok {
		t.Fatalf("type = %T", MapRecord(registries.MDRegistryID, record))
	}
	if len(got.RussianDocuments) != 1 || got.RussianDocuments[0].FileID != "rus" {
		t.Fatalf("RussianDocuments = %+v", got.RussianDocuments)
	}
	if len(got.EnglishDocuments) != 1 || got.EnglishDocuments[0].FileID != "eng" {
		t.Fatalf("EnglishDocuments = %+v", got.EnglishDocuments)
	}
	if len(got.Documents) != 2 {
		t.Fatalf("Documents = %+v", got.Documents)
	}
}

func TestMapStandardSampleTypeProductionRef(t *testing.T) {
	record := arshin.RegistryRecord{
		ID:   "3",
		Type: registries.UTSO.ItemType,
		Properties: []arshin.RegistryProperty{
			{Name: "foei:NumberUTSO", Value: "ГСО 1"},
			{Name: "foei:NameUTSO", Value: "Стандартный образец"},
			{Name: "foei:NameCertCharUTSO", Value: "Массовая доля"},
			{Name: "foei:ProductionUTSO", Type: "LINK_INTERNAL", Link: "/fundmetrology/registry/47/items/99"},
		},
	}

	got, ok := MapRecord(registries.UTSORegistryID, record).(StandardSampleType)
	if !ok {
		t.Fatalf("type = %T", MapRecord(registries.UTSORegistryID, record))
	}
	if got.RegistryNumber != "ГСО 1" || got.CertifiedCharacteristic != "Массовая доля" {
		t.Fatalf("mapped UTSO = %+v", got)
	}
	if len(got.ProductionRefs) != 1 || got.ProductionRefs[0].RegistryID != "47" || got.ProductionRefs[0].ItemID != "99" {
		t.Fatalf("ProductionRefs = %+v", got.ProductionRefs)
	}
}

func TestMapProductionNotice(t *testing.T) {
	record := arshin.RegistryRecord{
		ID:   "4",
		Type: registries.P1WF.ItemType,
		Properties: []arshin.RegistryProperty{
			{Name: "gost:registrationDate", Value: "01.01.2026"},
			{Name: "gost:registrationNumber", Value: "123-У"},
			{Name: "gost:p1wfOrganizationFullName", Value: "ООО Метрология"},
			{Name: "gost:p1wfOrganizationINN", Value: "7700000000"},
		},
	}

	got, ok := MapRecord(registries.P1WFRegistryID, record).(ProductionNotice)
	if !ok {
		t.Fatalf("type = %T", MapRecord(registries.P1WFRegistryID, record))
	}
	if got.RegistrationNumber != "123-У" || got.OrganizationINN != "7700000000" {
		t.Fatalf("mapped notice = %+v", got)
	}
}
