package client

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestRegistryItemFixtures(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "..", "testdata", "item-1404476-data.json"))
	if err != nil {
		t.Fatal(err)
	}

	var item RegistryItem
	if err := json.Unmarshal(data, &item); err != nil {
		t.Fatal(err)
	}

	if item.ID != 1404476 {
		t.Fatalf("item.ID = %d", item.ID)
	}
	if item.RegistryID != 16 {
		t.Fatalf("item.RegistryID = %d", item.RegistryID)
	}
	if len(item.Sections) != 5 {
		t.Fatalf("sections = %d", len(item.Sections))
	}
	fields := item.FieldsByName()
	if fields["foei:NumRegCMM"].Value != "ФР.1.31.2022.44733" {
		t.Fatalf("foei:NumRegCMM = %+v", fields["foei:NumRegCMM"])
	}
	if fields["foei:RMDocCMM"].Type != "ATTACH" {
		t.Fatalf("foei:RMDocCMM = %+v", fields["foei:RMDocCMM"])
	}
	if fields["foei:CMM1Relation_assoc"].Type != "CHILD_OBJECT" {
		t.Fatalf("foei:CMM1Relation_assoc = %+v", fields["foei:CMM1Relation_assoc"])
	}
}

func TestRegistryItemPlainDataFixture(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "..", "testdata", "item-1404476-plaindata.json"))
	if err != nil {
		t.Fatal(err)
	}

	var record RegistryRecord
	if err := json.Unmarshal(data, &record); err != nil {
		t.Fatal(err)
	}

	if record.ID != "1404476" {
		t.Fatalf("record.ID = %q", record.ID)
	}
	if record.Type != "foei:CMM1_type" {
		t.Fatalf("record.Type = %q", record.Type)
	}
	if len(record.Properties) != 30 {
		t.Fatalf("properties = %d", len(record.Properties))
	}
	if record.PropertiesByName()["foei:NumRegCMM"].Value != "ФР.1.31.2022.44733" {
		t.Fatalf("foei:NumRegCMM = %+v", record.PropertiesByName()["foei:NumRegCMM"])
	}
}
