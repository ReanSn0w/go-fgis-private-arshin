package client

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestParseDownloadFileLink(t *testing.T) {
	fileID, ok := ParseDownloadFileLink("/api/downloadfile/c11f089b-22f4-4d66-87cd-7aff99909001")
	if !ok {
		t.Fatal("download link was not parsed")
	}
	if fileID != "c11f089b-22f4-4d66-87cd-7aff99909001" {
		t.Fatalf("fileID = %q", fileID)
	}
}

func TestRegistryPropertyAttachments(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "..", "testdata", "item-1404476-data.json"))
	if err != nil {
		t.Fatal(err)
	}

	var item RegistryItem
	if err := json.Unmarshal(data, &item); err != nil {
		t.Fatal(err)
	}

	property := item.FieldsByName()["foei:RMDocCMM"]
	attachments := property.Attachments()
	if len(attachments) != 1 {
		t.Fatalf("attachments = %d", len(attachments))
	}
	if attachments[0].Name != "ПРИЛОЖЕНИЕ_048.pdf" {
		t.Fatalf("Name = %q", attachments[0].Name)
	}
	if attachments[0].MIME != "application/pdf" {
		t.Fatalf("MIME = %q", attachments[0].MIME)
	}
	if attachments[0].FileID != "c11f089b-22f4-4d66-87cd-7aff99909001" {
		t.Fatalf("FileID = %q", attachments[0].FileID)
	}
}
