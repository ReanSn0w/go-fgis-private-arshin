package client

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDownloadFile(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/downloadfile/file-id" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", `form-data; name="attachment"; filename="%D0%A4%D0%B0%D0%B9%D0%BB.pdf"`)
		_, _ = w.Write([]byte("%PDF"))
	}))
	defer server.Close()

	client, err := NewClient(
		WithBaseURL(server.URL+"/api/"),
		WithRateLimit(0),
	)
	if err != nil {
		t.Fatal(err)
	}

	file, err := client.DownloadFile(context.Background(), "file-id")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Body.Close()

	if file.Filename != "Файл.pdf" {
		t.Fatalf("filename = %q", file.Filename)
	}
	if file.ContentType != "application/octet-stream" {
		t.Fatalf("contentType = %q", file.ContentType)
	}

	body, err := io.ReadAll(file.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != "%PDF" {
		t.Fatalf("body = %q", body)
	}
}

func TestDownloadFileJSONError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(APIError{
			Status:  http.StatusInternalServerError,
			Message: "file not found",
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

	_, err = client.DownloadFile(context.Background(), "missing")
	if err == nil {
		t.Fatal("expected error")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("error = %T %v", err, err)
	}
	if apiErr.Message != "file not found" {
		t.Fatalf("apiErr.Message = %q", apiErr.Message)
	}
}

func TestDownloadFileLink(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/downloadfile/file-id" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		_, _ = w.Write([]byte("ok"))
	}))
	defer server.Close()

	client, err := NewClient(
		WithBaseURL(server.URL+"/api/"),
		WithRateLimit(0),
	)
	if err != nil {
		t.Fatal(err)
	}

	file, err := client.DownloadFileLink(context.Background(), "/api/downloadfile/file-id")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Body.Close()

	body, err := io.ReadAll(file.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != "ok" {
		t.Fatalf("body = %q", body)
	}
}
