package config

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

var (
	idPattern       = regexp.MustCompile(`^[a-z0-9][a-z0-9_-]{0,63}$`)
	colorPattern    = regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)
	shortcutPattern = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,24}$`)
)

func Validate(cfg Config) error {
	if cfg.Version != 2 {
		return errors.New("unsupported config version")
	}

	title := strings.TrimSpace(cfg.Dashboard.Title)
	if title == "" || len(title) > 80 {
		return errors.New("dashboard title must contain between 1 and 80 characters")
	}

	switch cfg.Search.WebSearchEngine {
	case "google", "duckduckgo", "bing":
	default:
		return fmt.Errorf(
			"unsupported web search engine %q",
			cfg.Search.WebSearchEngine,
		)
	}

	sectionIDs := map[string]struct{}{}
	itemIDs := map[string]struct{}{}

	for _, section := range cfg.Sections {
		if !idPattern.MatchString(section.ID) {
			return fmt.Errorf("invalid section id %q", section.ID)
		}

		if _, exists := sectionIDs[section.ID]; exists {
			return fmt.Errorf("duplicate section id %q", section.ID)
		}
		sectionIDs[section.ID] = struct{}{}

		if strings.TrimSpace(section.Title) == "" {
			return fmt.Errorf("section %q requires a title", section.ID)
		}

		if !colorPattern.MatchString(section.Accent) {
			return fmt.Errorf(
				"section %q has invalid accent color",
				section.ID,
			)
		}

		switch section.Layout {
		case "grid", "list", "compact", "large":
		default:
			return fmt.Errorf(
				"section %q has unsupported layout %q",
				section.ID,
				section.Layout,
			)
		}

		for _, item := range section.Items {
			if err := validateItem(item, itemIDs); err != nil {
				return err
			}
		}
	}

	shortcutKeys := map[string]struct{}{}

	for _, shortcut := range cfg.Shortcuts {
		key := strings.ToLower(strings.TrimSpace(shortcut.Key))

		if !shortcutPattern.MatchString(key) {
			return fmt.Errorf("invalid shortcut key %q", shortcut.Key)
		}

		if _, exists := shortcutKeys[key]; exists {
			return fmt.Errorf("duplicate shortcut key %q", key)
		}
		shortcutKeys[key] = struct{}{}

		if strings.TrimSpace(shortcut.Label) == "" {
			return fmt.Errorf("shortcut %q requires a label", key)
		}

		if err := validateHTTPURL(shortcut.URL); err != nil {
			return fmt.Errorf("shortcut %q: %w", key, err)
		}
	}

	return nil
}

func validateItem(item Item, existing map[string]struct{}) error {
	if !idPattern.MatchString(item.ID) {
		return fmt.Errorf("invalid item id %q", item.ID)
	}

	if _, exists := existing[item.ID]; exists {
		return fmt.Errorf("duplicate item id %q", item.ID)
	}
	existing[item.ID] = struct{}{}

	if strings.TrimSpace(item.Name) == "" {
		return fmt.Errorf("item %q requires a name", item.ID)
	}

	if err := validateHTTPURL(item.URL); err != nil {
		return fmt.Errorf("item %q: %w", item.ID, err)
	}

	switch item.Icon.Type {
	case "auto", "local", "none":
	default:
		return fmt.Errorf(
			"item %q has unsupported icon type",
			item.ID,
		)
	}

	return nil
}

func validateHTTPURL(value string) error {
	parsed, err := url.ParseRequestURI(value)
	if err != nil {
		return errors.New("invalid URL")
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return errors.New("URL must use http or https")
	}

	if parsed.Host == "" {
		return errors.New("URL requires a hostname")
	}

	return nil
}
