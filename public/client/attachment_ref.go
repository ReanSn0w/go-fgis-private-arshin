package client

import (
	"net/url"
	"path"
	"strings"
)

type AttachmentRef struct {
	Name   string `json:"name"`
	Link   string `json:"link"`
	MIME   string `json:"mime"`
	FileID string `json:"fileId"`
}

func ParseDownloadFileLink(link string) (fileID string, ok bool) {
	link = strings.TrimSpace(link)
	if link == "" {
		return "", false
	}

	parsed, err := url.Parse(link)
	if err == nil {
		link = parsed.Path
	}

	parts := strings.Split(strings.Trim(link, "/"), "/")
	for i := 0; i < len(parts)-1; i++ {
		if parts[i] == "downloadfile" && parts[i+1] != "" {
			return parts[i+1], true
		}
	}

	base := path.Base(link)
	if base == "." || base == "/" || base == "" {
		return "", false
	}
	return base, true
}

func (p RegistryProperty) Attachments() []AttachmentRef {
	names := stringsFromRaw(p.Value)
	links := stringsFromRaw(p.Link)
	mimes := stringsFromRaw(p.MIME)

	count := max(len(names), len(links), len(mimes))
	refs := make([]AttachmentRef, 0, count)
	for i := 0; i < count; i++ {
		ref := AttachmentRef{
			Name: valueAt(names, i),
			Link: valueAt(links, i),
			MIME: valueAt(mimes, i),
		}
		if ref.MIME == "" && len(mimes) == 1 {
			ref.MIME = mimes[0]
		}
		if fileID, ok := ParseDownloadFileLink(ref.Link); ok {
			ref.FileID = fileID
		}
		if ref.Name != "" || ref.Link != "" || ref.MIME != "" {
			refs = append(refs, ref)
		}
	}
	return refs
}

func valueAt(values []string, index int) string {
	if index < 0 || index >= len(values) {
		return ""
	}
	return values[index]
}
