package client

import (
	"encoding/json"
	"fmt"
	"regexp"
)

var registryItemLinkPattern = regexp.MustCompile(`(?:^|/)registry/([^/]+)/items/([^/?#]+)`)

type RegistryItemRef struct {
	RegistryID RegistryID     `json:"registryId"`
	ItemID     RegistryItemID `json:"itemId"`
	Link       string         `json:"link"`
}

type ChildObject struct {
	ID      string              `json:"id"`
	Type    string              `json:"type"`
	NodeRef string              `json:"nodeRef"`
	Columns []ChildObjectColumn `json:"columns"`
	Fields  []RegistryProperty  `json:"fields"`
}

type ChildObjectColumn struct {
	Name  string `json:"name"`
	Title string `json:"title"`
}

func ParseRegistryItemLink(link string) (RegistryItemRef, bool) {
	matches := registryItemLinkPattern.FindStringSubmatch(link)
	if len(matches) != 3 {
		return RegistryItemRef{}, false
	}
	return RegistryItemRef{
		RegistryID: RegistryID(matches[1]),
		ItemID:     RegistryItemID(matches[2]),
		Link:       link,
	}, true
}

func (p RegistryProperty) LinkStrings() []string {
	return stringsFromRaw(p.Link)
}

func (p RegistryProperty) ItemRefs() []RegistryItemRef {
	links := p.LinkStrings()
	refs := make([]RegistryItemRef, 0, len(links))
	for _, link := range links {
		ref, ok := ParseRegistryItemLink(link)
		if ok {
			refs = append(refs, ref)
		}
	}
	return refs
}

func (p RegistryProperty) ChildObjects() ([]ChildObject, error) {
	if p.Value == nil {
		return nil, nil
	}

	data, err := json.Marshal(p.Value)
	if err != nil {
		return nil, fmt.Errorf("marshal child objects: %w", err)
	}

	var children []ChildObject
	if err := json.Unmarshal(data, &children); err != nil {
		return nil, fmt.Errorf("unmarshal child objects: %w", err)
	}
	return children, nil
}

func (c ChildObject) FieldsByName() map[string]RegistryProperty {
	return PropertiesByName(c.Fields)
}

func (c ChildObject) ItemRefs() []RegistryItemRef {
	var refs []RegistryItemRef
	for _, field := range c.Fields {
		refs = append(refs, field.ItemRefs()...)
	}
	return refs
}

func stringsFromRaw(value any) []string {
	switch v := value.(type) {
	case nil:
		return nil
	case string:
		if v == "" {
			return nil
		}
		return []string{v}
	case []string:
		return v
	case []any:
		values := make([]string, 0, len(v))
		for _, item := range v {
			switch typed := item.(type) {
			case string:
				if typed != "" {
					values = append(values, typed)
				}
			}
		}
		return values
	default:
		return nil
	}
}
