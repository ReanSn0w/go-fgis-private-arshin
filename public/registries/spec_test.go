package registries_test

import (
	"testing"

	"github.com/ReanSn0w/go-fgis-private-arshin/public/registries"
)

func TestKnownRegistries(t *testing.T) {
	known := registries.Known()
	if len(known) != 16 {
		t.Fatalf("known registries = %d", len(known))
	}

	spec, ok := registries.SpecForRegistry("19")
	if !ok {
		t.Fatal("registry 19 spec not found")
	}
	if spec.ItemType != registries.UTSO.ItemType {
		t.Fatalf("registry 19 item type = %q", spec.ItemType)
	}

	spec, ok = registries.SpecForItemType("gost:p1wfRequestType4")
	if !ok {
		t.Fatal("p1wf spec not found by item type")
	}
	if spec.RegistryID != registries.P1WF.RegistryID {
		t.Fatalf("p1wf registry id = %q", spec.RegistryID)
	}
}

func TestNamedRegistrySpecs(t *testing.T) {
	spec := registries.CMM2
	if spec.RegistryID != registries.CMM2RegistryID {
		t.Fatalf("cmm2 registry id = %q", spec.RegistryID)
	}
	if spec.Title == "" {
		t.Fatalf("cmm2 title = %q", spec.Title)
	}
	if spec.ItemType != "foei:CMM2_type" {
		t.Fatalf("cmm2 item type = %q", spec.ItemType)
	}
}
