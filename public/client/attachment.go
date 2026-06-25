package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"path"
	"strings"
)

type DownloadedFile struct {
	Body          io.ReadCloser
	Filename      string
	ContentType   string
	ContentLength int64
}

func (c *Client) DownloadFile(ctx context.Context, fileID string) (*DownloadedFile, error) {
	fileID = strings.TrimSpace(fileID)
	if fileID == "" {
		return nil, fmt.Errorf("file id is empty")
	}
	return c.DownloadFileLink(ctx, "/api/downloadfile/"+fileID)
}

func (c *Client) DownloadFileLink(ctx context.Context, link string) (*DownloadedFile, error) {
	link = strings.TrimSpace(link)
	if link == "" {
		return nil, fmt.Errorf("download link is empty")
	}

	if err := c.limiter.Wait(ctx); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.downloadURL(link), nil)
	if err != nil {
		return nil, err
	}
	c.setDefaultHeaders(req, "")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		defer resp.Body.Close()
		if strings.Contains(strings.ToLower(resp.Header.Get("Content-Type")), "json") {
			var apiErr APIError
			if err := json.NewDecoder(resp.Body).Decode(&apiErr); err == nil && apiErr.Status != 0 {
				return nil, &apiErr
			}
		}
		return nil, NewHTTPError(resp)
	}

	return &DownloadedFile{
		Body:          resp.Body,
		Filename:      filenameFromContentDisposition(resp.Header.Get("Content-Disposition")),
		ContentType:   resp.Header.Get("Content-Type"),
		ContentLength: resp.ContentLength,
	}, nil
}

func (c *Client) downloadURL(link string) string {
	if parsed, err := url.Parse(link); err == nil && parsed.IsAbs() {
		return parsed.String()
	}

	trimmed := strings.TrimPrefix(link, "/")
	if strings.HasPrefix(trimmed, "api/") {
		trimmed = strings.TrimPrefix(trimmed, "api/")
	}
	return c.newAPIURL(trimmed, nil)
}

func filenameFromContentDisposition(contentDisposition string) string {
	_, params, err := mime.ParseMediaType(contentDisposition)
	if err != nil {
		return ""
	}

	if filename := params["filename"]; filename != "" {
		return cleanFilename(filename)
	}
	if filename := params["filename*"]; filename != "" {
		return cleanFilename(filename)
	}
	return ""
}

func cleanFilename(filename string) string {
	if decoded, err := url.PathUnescape(filename); err == nil {
		filename = decoded
	}
	return path.Base(filename)
}
