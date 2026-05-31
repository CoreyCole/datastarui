package selectors

import (
	"fmt"
	"os"

	"github.com/coreycole/datastarui/e2e/appconfig"
	"gopkg.in/yaml.v3"
)

type Key string

type Entry struct {
	Key         Key    `yaml:"key"`
	CSS         string `yaml:"css"`
	Description string `yaml:"description"`
	StableID    string `yaml:"stable_id"`
}

type Catalog struct{ entries map[Key]Entry }

func NewCatalog(entries []Entry) Catalog {
	out := Catalog{entries: map[Key]Entry{}}
	for _, entry := range entries {
		if entry.Key == "" || entry.CSS == "" {
			continue
		}
		out.entries[entry.Key] = entry
	}
	return out
}

func FromMap(values map[string]string) Catalog {
	entries := make([]Entry, 0, len(values))
	for key, css := range values {
		entries = append(entries, Entry{Key: Key(key), CSS: css})
	}
	return NewCatalog(entries)
}

func LoadCatalog(path string) (Catalog, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Catalog{}, err
	}

	var entries []Entry
	if err := yaml.Unmarshal(data, &entries); err == nil && len(entries) > 0 {
		return NewCatalog(entries), nil
	}

	var keyed map[string]string
	if err := yaml.Unmarshal(data, &keyed); err != nil {
		return Catalog{}, err
	}
	return FromMap(keyed), nil
}

func LoadCatalogFromConfig(cfg appconfig.Config) (Catalog, error) {
	catalog := FromMap(cfg.Selectors)
	if cfg.SelectorFile == "" {
		return catalog, nil
	}
	fromFile, err := LoadCatalog(cfg.ResolvePath(cfg.SelectorFile))
	if err != nil {
		return Catalog{}, err
	}
	return catalog.Merge(fromFile), nil
}

func (c Catalog) Merge(other Catalog) Catalog {
	out := Catalog{entries: map[Key]Entry{}}
	for key, entry := range c.entries {
		out.entries[key] = entry
	}
	for key, entry := range other.entries {
		out.entries[key] = entry
	}
	return out
}

func (c Catalog) Resolve(key string) (Entry, error) {
	entry, ok := c.entries[Key(key)]
	if !ok {
		return Entry{}, fmt.Errorf("unknown selector key %q", key)
	}
	return entry, nil
}

func (c Catalog) ValidateKeys(keys []string) error {
	for _, key := range keys {
		if _, err := c.Resolve(key); err != nil {
			return err
		}
	}
	return nil
}

func (c Catalog) Len() int {
	return len(c.entries)
}
