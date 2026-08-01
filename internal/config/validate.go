package config

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

var (
	idPattern = regexp.MustCompile(
		`^[a-z0-9][a-z0-9_-]{0,63}$`,
	)

	colorPattern = regexp.MustCompile(
		`^#[0-9a-fA-F]{6}$`,
	)

	shortcutPattern = regexp.MustCompile(
		`^[a-zA-Z0-9_-]{1,24}$`,
	)
)

func Validate(cfg Config) error {
	if cfg.Version != 2 {
		return errors.New(
			"unsupported config version",
		)
	}

	title := strings.TrimSpace(
		cfg.Dashboard.Title,
	)

	if title == "" || len(title) > 80 {
		return errors.New(
			"dashboard title must contain between 1 and 80 characters",
		)
	}

	switch cfg.Appearance.Mode {
	case "minimal", "standard", "advanced":
	default:
		return fmt.Errorf(
			"unsupported appearance mode %q",
			cfg.Appearance.Mode,
		)
	}

	switch cfg.Appearance.Font.Family {
	case "system",
		"inter",
		"geist",
		"manrope",
		"ibm-plex-sans":
	default:
		return fmt.Errorf(
			"unsupported font family %q",
			cfg.Appearance.Font.Family,
		)
	}

	switch cfg.Appearance.Font.Scale {
	case "small", "medium", "large":
	default:
		return fmt.Errorf(
			"unsupported font scale %q",
			cfg.Appearance.Font.Scale,
		)
	}

	switch cfg.Appearance.Cards.Size {
	case "small", "medium", "large":
	default:
		return fmt.Errorf(
			"unsupported default card size %q",
			cfg.Appearance.Cards.Size,
		)
	}

	switch cfg.Appearance.Cards.Density {
	case "compact", "comfortable", "relaxed":
	default:
		return fmt.Errorf(
			"unsupported card density %q",
			cfg.Appearance.Cards.Density,
		)
	}

	switch cfg.Appearance.Cards.Radius {
	case "small", "medium", "large":
	default:
		return fmt.Errorf(
			"unsupported card radius %q",
			cfg.Appearance.Cards.Radius,
		)
	}

	switch cfg.Appearance.Cards.Shadow {
	case "none", "soft", "medium", "floating":
	default:
		return fmt.Errorf(
			"unsupported card shadow %q",
			cfg.Appearance.Cards.Shadow,
		)
	}

	switch cfg.Appearance.Background.Type {
	case "solid", "gradient", "wallpaper":
	default:
		return fmt.Errorf(
			"unsupported background type %q",
			cfg.Appearance.Background.Type,
		)
	}

	if !colorPattern.MatchString(
		cfg.Appearance.Background.Color,
	) {
		return errors.New(
			"invalid background color",
		)
	}

	if cfg.Appearance.Background.Blur < 0 ||
		cfg.Appearance.Background.Blur > 40 {
		return errors.New(
			"background blur must be between 0 and 40",
		)
	}

	if cfg.Appearance.Background.Brightness < 20 ||
		cfg.Appearance.Background.Brightness > 150 {
		return errors.New(
			"background brightness must be between 20 and 150",
		)
	}

	if cfg.Appearance.Background.Overlay < 0 ||
		cfg.Appearance.Background.Overlay > 100 {
		return errors.New(
			"background overlay must be between 0 and 100",
		)
	}

	switch cfg.Search.WebSearchEngine {
	case "google", "duckduckgo", "bing":
	default:
		return fmt.Errorf(
			"unsupported web search engine %q",
			cfg.Search.WebSearchEngine,
		)
	}

	switch cfg.Modules.Weather.Mode {
	case "current", "today", "five-day":
	default:
		return fmt.Errorf(
			"unsupported weather mode %q",
			cfg.Modules.Weather.Mode,
		)
	}

	switch cfg.Modules.Clock.Style {
	case "minimal", "digital", "large":
	default:
		return fmt.Errorf(
			"unsupported clock style %q",
			cfg.Modules.Clock.Style,
		)
	}

	switch cfg.Modules.Clock.TimeFormat {
	case "12h", "24h":
	default:
		return fmt.Errorf(
			"unsupported clock time format %q",
			cfg.Modules.Clock.TimeFormat,
		)
	}

	sectionIDs := map[string]struct{}{}
	itemIDs := map[string]struct{}{}

	for _, section := range cfg.Sections {
		if !idPattern.MatchString(section.ID) {
			return fmt.Errorf(
				"invalid section id %q",
				section.ID,
			)
		}

		if _, exists := sectionIDs[section.ID]; exists {
			return fmt.Errorf(
				"duplicate section id %q",
				section.ID,
			)
		}

		sectionIDs[section.ID] = struct{}{}

		if strings.TrimSpace(section.Title) == "" {
			return fmt.Errorf(
				"section %q requires a title",
				section.ID,
			)
		}

		if len(section.Title) > 80 {
			return fmt.Errorf(
				"section %q title is too long",
				section.ID,
			)
		}

		if !colorPattern.MatchString(section.Accent) {
			return fmt.Errorf(
				"section %q has invalid accent color",
				section.ID,
			)
		}

		switch section.Surface {
		case "solid",
			"transparent",
			"glass",
			"none":
		default:
			return fmt.Errorf(
				"section %q has unsupported surface %q",
				section.ID,
				section.Surface,
			)
		}

		if section.SurfaceOpacity < 0 ||
			section.SurfaceOpacity > 100 {
			return fmt.Errorf(
				"section %q surface opacity must be between 0 and 100",
				section.ID,
			)
		}

		if section.SurfaceBlur < 0 ||
			section.SurfaceBlur > 40 {
			return fmt.Errorf(
				"section %q surface blur must be between 0 and 40",
				section.ID,
			)
		}

		switch section.Layout {
		case "list",
			"grid",
			"compact",
			"featured":
		default:
			return fmt.Errorf(
				"section %q has unsupported card arrangement %q",
				section.ID,
				section.Layout,
			)
		}

		switch section.Width {
		case "narrow",
			"medium",
			"wide",
			"extra-wide",
			"full":
		default:
			return fmt.Errorf(
				"section %q has unsupported width %q",
				section.ID,
				section.Width,
			)
		}

		switch section.CardSize {
		case "inherit",
			"small",
			"medium",
			"large":
		default:
			return fmt.Errorf(
				"section %q has unsupported card size %q",
				section.ID,
				section.CardSize,
			)
		}

		if section.GridColumns < 1 ||
			section.GridColumns > 6 {
			return fmt.Errorf(
				"section %q grid columns must be between 1 and 6",
				section.ID,
			)
		}

		for _, item := range section.Items {
			if err := validateItem(
				item,
				itemIDs,
			); err != nil {
				return err
			}
		}
	}

	shortcutKeys := map[string]struct{}{}

	for _, shortcut := range cfg.Shortcuts {
		key := strings.ToLower(
			strings.TrimSpace(shortcut.Key),
		)

		if !shortcutPattern.MatchString(key) {
			return fmt.Errorf(
				"invalid shortcut key %q",
				shortcut.Key,
			)
		}

		if _, exists := shortcutKeys[key]; exists {
			return fmt.Errorf(
				"duplicate shortcut key %q",
				key,
			)
		}

		shortcutKeys[key] = struct{}{}

		if strings.TrimSpace(shortcut.Label) == "" {
			return fmt.Errorf(
				"shortcut %q requires a label",
				key,
			)
		}

		if len(shortcut.Label) > 100 {
			return fmt.Errorf(
				"shortcut %q label is too long",
				key,
			)
		}

		if err := validateHTTPURL(
			shortcut.URL,
		); err != nil {
			return fmt.Errorf(
				"shortcut %q: %w",
				key,
				err,
			)
		}

		switch shortcut.Icon.Type {
		case "auto", "local", "none":
		default:
			return fmt.Errorf(
				"shortcut %q has unsupported icon type",
				key,
			)
		}
	}

	if cfg.Modules.Weather.Enabled {
		if strings.TrimSpace(
			cfg.Modules.Weather.Location,
		) == "" &&
			(cfg.Modules.Weather.Latitude == nil ||
				cfg.Modules.Weather.Longitude == nil) {
			return errors.New(
				"enabled weather module requires a location or coordinates",
			)
		}

		if cfg.Modules.Weather.Latitude != nil {
			latitude :=
				*cfg.Modules.Weather.Latitude

			if latitude < -90 || latitude > 90 {
				return errors.New(
					"weather latitude must be between -90 and 90",
				)
			}
		}

		if cfg.Modules.Weather.Longitude != nil {
			longitude :=
				*cfg.Modules.Weather.Longitude

			if longitude < -180 ||
				longitude > 180 {
				return errors.New(
					"weather longitude must be between -180 and 180",
				)
			}
		}
	}

	return nil
}

func validateItem(
	item Item,
	existing map[string]struct{},
) error {
	if !idPattern.MatchString(item.ID) {
		return fmt.Errorf(
			"invalid item id %q",
			item.ID,
		)
	}

	if _, exists := existing[item.ID]; exists {
		return fmt.Errorf(
			"duplicate item id %q",
			item.ID,
		)
	}

	existing[item.ID] = struct{}{}

	switch item.Type {
	case "service":
	default:
		return fmt.Errorf(
			"item %q has unsupported type %q",
			item.ID,
			item.Type,
		)
	}

	if strings.TrimSpace(item.Name) == "" {
		return fmt.Errorf(
			"item %q requires a name",
			item.ID,
		)
	}

	if len(item.Name) > 100 {
		return fmt.Errorf(
			"item %q name is too long",
			item.ID,
		)
	}

	if len(item.Description) > 240 {
		return fmt.Errorf(
			"item %q description is too long",
			item.ID,
		)
	}

	if err := validateHTTPURL(item.URL); err != nil {
		return fmt.Errorf(
			"item %q: %w",
			item.ID,
			err,
		)
	}

	switch item.Icon.Type {
	case "auto", "local", "none":
	default:
		return fmt.Errorf(
			"item %q has unsupported icon type",
			item.ID,
		)
	}

	if item.Icon.Type == "local" &&
		strings.TrimSpace(item.Icon.Value) == "" {
		return fmt.Errorf(
			"item %q local icon requires a value",
			item.ID,
		)
	}

	return nil
}

func validateHTTPURL(
	value string,
) error {
	parsed, err := url.ParseRequestURI(value)

	if err != nil {
		return errors.New(
			"invalid URL",
		)
	}

	if parsed.Scheme != "http" &&
		parsed.Scheme != "https" {
		return errors.New(
			"URL must use http or https",
		)
	}

	if parsed.Host == "" {
		return errors.New(
			"URL requires a hostname",
		)
	}

	return nil
}
