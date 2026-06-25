package methods

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	arshin "github.com/ReanSn0w/go-fgis-private-arshin/public/client"
)

func TestMapRecord(t *testing.T) {
	fixturePath := filepath.Join("..", "..", "..", "testdata", "fr-1-31-2022-44733.json")
	data, err := os.ReadFile(fixturePath)
	if err != nil {
		t.Fatal(err)
	}

	var page arshin.RegistryDataPage
	if err := json.Unmarshal(data, &page); err != nil {
		t.Fatal(err)
	}
	if len(page.Items) != 1 {
		t.Fatalf("items = %d", len(page.Items))
	}

	method := MapRecord(page.Items[0])

	if method.ID != "1404476" {
		t.Fatalf("ID = %q", method.ID)
	}
	if method.Type != "foei:CMM1_type" {
		t.Fatalf("Type = %q", method.Type)
	}
	if method.RegistryID != CertifiedRegistryID {
		t.Fatalf("RegistryID = %q", method.RegistryID)
	}
	if method.RegistryNumber != "ФР.1.31.2022.44733" {
		t.Fatalf("RegistryNumber = %q", method.RegistryNumber)
	}
	if method.RecordNumber != "44733" {
		t.Fatalf("RecordNumber = %q", method.RecordNumber)
	}
	if method.Status != "Действует" {
		t.Fatalf("Status = %q", method.Status)
	}
	if method.SystemStatus != "Опубликована" {
		t.Fatalf("SystemStatus = %q", method.SystemStatus)
	}
	if method.PublishedAt != "15.12.2022" {
		t.Fatalf("PublishedAt = %q", method.PublishedAt)
	}
	if method.CertificateDate != "31.10.2022" {
		t.Fatalf("CertificateDate = %q", method.CertificateDate)
	}
	if method.CertificateNumber != "88-16207-048-RA.RU.310657-2022" {
		t.Fatalf("CertificateNumber = %q", method.CertificateNumber)
	}
	if method.CertificationOrgName != "АХУ УРО РАН" {
		t.Fatalf("CertificationOrgName = %q", method.CertificationOrgName)
	}
	if !strings.Contains(method.Name, "ПНД Ф 14.1:2:3:4.48-2022") {
		t.Fatalf("Name = %q", method.Name)
	}
	if !strings.Contains(method.DeveloperName, "ФГБУ") {
		t.Fatalf("DeveloperName = %q", method.DeveloperName)
	}
	if len(method.RawProperties) == 0 {
		t.Fatal("RawProperties is empty")
	}
}

func TestMapRecordAttachmentsAndRelations(t *testing.T) {
	fixturePath := filepath.Join("..", "..", "..", "testdata", "fr-1-31-2022-44733.json")
	data, err := os.ReadFile(fixturePath)
	if err != nil {
		t.Fatal(err)
	}

	var page arshin.RegistryDataPage
	if err := json.Unmarshal(data, &page); err != nil {
		t.Fatal(err)
	}

	method := MapRecord(page.Items[0])

	if len(method.RangeDocuments) != 1 {
		t.Fatalf("RangeDocuments = %d", len(method.RangeDocuments))
	}
	if method.RangeDocuments[0].FileID != "c11f089b-22f4-4d66-87cd-7aff99909001" {
		t.Fatalf("RangeDocuments[0].FileID = %q", method.RangeDocuments[0].FileID)
	}
	if len(method.Relations) != 1 {
		t.Fatalf("Relations = %d", len(method.Relations))
	}
	if method.Relations[0].Type != "Заменяет" {
		t.Fatalf("Relation.Type = %q", method.Relations[0].Type)
	}
	if method.Relations[0].RelatedRegistryNumber != "ФР.1.31.2013.16016" {
		t.Fatalf("Relation.RelatedRegistryNumber = %q", method.Relations[0].RelatedRegistryNumber)
	}
	if method.Relations[0].RelatedRef == nil || method.Relations[0].RelatedRef.ItemID != "286196" {
		t.Fatalf("Relation.RelatedRef = %+v", method.Relations[0].RelatedRef)
	}
}

func TestMapItem(t *testing.T) {
	fixturePath := filepath.Join("..", "..", "..", "testdata", "item-1404476-data.json")
	data, err := os.ReadFile(fixturePath)
	if err != nil {
		t.Fatal(err)
	}

	var item arshin.RegistryItem
	if err := json.Unmarshal(data, &item); err != nil {
		t.Fatal(err)
	}

	method := MapItem(item)

	if method.ID != "1404476" {
		t.Fatalf("ID = %q", method.ID)
	}
	if method.RegistryID != CertifiedRegistryID {
		t.Fatalf("RegistryID = %q", method.RegistryID)
	}
	if method.Type != "foei:CMM1_type" {
		t.Fatalf("Type = %q", method.Type)
	}
	if method.RegistryNumber != "ФР.1.31.2022.44733" {
		t.Fatalf("RegistryNumber = %q", method.RegistryNumber)
	}
	if len(method.RangeDocuments) != 1 {
		t.Fatalf("RangeDocuments = %d", len(method.RangeDocuments))
	}
	if len(method.Relations) != 1 {
		t.Fatalf("Relations = %d", len(method.Relations))
	}
	if method.RawItem == nil {
		t.Fatal("RawItem is nil")
	}
}

func TestSupportedRegistries(t *testing.T) {
	if _, err := NewForRegistry(nil, PrimaryReferenceRegistryID); err != nil {
		t.Fatalf("NewForRegistry registry 6: %v", err)
	}
	if _, err := NewForRegistry(nil, ReferenceRegistryID); err != nil {
		t.Fatalf("NewForRegistry registry 8: %v", err)
	}
	if _, err := NewForRegistry(nil, "19"); err == nil {
		t.Fatal("NewForRegistry registry 19 succeeded")
	}
}

func TestRegistrySpecs(t *testing.T) {
	spec, ok := SpecForRegistry(PrimaryReferenceRegistryID)
	if !ok {
		t.Fatal("registry 6 spec not found")
	}
	if spec.ItemType != "foei:CMM2_type" {
		t.Fatalf("ItemType = %q", spec.ItemType)
	}
	if spec.OrderFields.Number != "foei:NumberOrderMetCMM" {
		t.Fatalf("OrderFields.Number = %q", spec.OrderFields.Number)
	}

	spec, ok = SpecForItemType("foei:CMM3_type")
	if !ok {
		t.Fatal("CMM3 spec not found")
	}
	if spec.RegistryID != ReferenceRegistryID {
		t.Fatalf("RegistryID = %q", spec.RegistryID)
	}
}
