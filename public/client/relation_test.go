package client

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestParseRegistryItemLink(t *testing.T) {
	ref, ok := ParseRegistryItemLink("/registry/16/items/286196")
	if !ok {
		t.Fatal("link was not parsed")
	}
	if ref.RegistryID != "16" {
		t.Fatalf("RegistryID = %q", ref.RegistryID)
	}
	if ref.ItemID != "286196" {
		t.Fatalf("ItemID = %q", ref.ItemID)
	}
}

func TestRegistryPropertyItemRefs(t *testing.T) {
	property := RegistryProperty{
		Type: "LINK_INTERNAL",
		Link: "/registry/13/items/394779",
	}

	refs := property.ItemRefs()
	if len(refs) != 1 {
		t.Fatalf("refs = %d", len(refs))
	}
	if refs[0].RegistryID != "13" || refs[0].ItemID != "394779" {
		t.Fatalf("ref = %+v", refs[0])
	}
}

func TestRegistryPropertyChildObjects(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "..", "testdata", "item-1404476-data.json"))
	if err != nil {
		t.Fatal(err)
	}

	var item RegistryItem
	if err := json.Unmarshal(data, &item); err != nil {
		t.Fatal(err)
	}

	relationProperty := item.FieldsByName()["foei:CMM1Relation_assoc"]
	children, err := relationProperty.ChildObjects()
	if err != nil {
		t.Fatal(err)
	}
	if len(children) != 1 {
		t.Fatalf("children = %d", len(children))
	}

	child := children[0]
	if child.Type != "foei:CMM1Relation" {
		t.Fatalf("child.Type = %q", child.Type)
	}
	if child.FieldsByName()["foei:RelationTypeCMM"].Value != "Заменяет" {
		t.Fatalf("relation type = %+v", child.FieldsByName()["foei:RelationTypeCMM"])
	}
	if child.FieldsByName()["foei:RelatedCMM"].Value != "ФР.1.31.2013.16016" {
		t.Fatalf("related method = %+v", child.FieldsByName()["foei:RelatedCMM"])
	}

	refs := child.ItemRefs()
	if len(refs) != 1 {
		t.Fatalf("refs = %d", len(refs))
	}
	if refs[0].RegistryID != "16" || refs[0].ItemID != "286196" {
		t.Fatalf("ref = %+v", refs[0])
	}
}
