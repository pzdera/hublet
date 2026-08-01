package config

func Normalize(cfg *Config) {
	if cfg.Version == 0 {
		cfg.Version = 2
	}

	if cfg.Dashboard.Title == "" {
		cfg.Dashboard.Title = "Hublet"
	}

	if cfg.Dashboard.Theme == "" {
		cfg.Dashboard.Theme = "midnight"
	}

	if cfg.Appearance.Mode == "" {
		cfg.Appearance.Mode = "minimal"
	}

	if cfg.Appearance.Font.Family == "" {
		cfg.Appearance.Font.Family = "system"
	}

	if cfg.Appearance.Font.Scale == "" {
		cfg.Appearance.Font.Scale = "medium"
	}

	if cfg.Appearance.Cards.Size == "" {
		cfg.Appearance.Cards.Size = "medium"
	}

	if cfg.Appearance.Cards.Density == "" {
		cfg.Appearance.Cards.Density = "comfortable"
	}

	if cfg.Appearance.Cards.Radius == "" {
		cfg.Appearance.Cards.Radius = "large"
	}

	if cfg.Appearance.Cards.Shadow == "" {
		cfg.Appearance.Cards.Shadow = "soft"
	}

	if cfg.Appearance.Background.Type == "" {
		cfg.Appearance.Background.Type = "solid"
	}

	if cfg.Appearance.Background.Color == "" {
		cfg.Appearance.Background.Color = "#090c12"
	}

	if cfg.Appearance.Background.Brightness == 0 {
		cfg.Appearance.Background.Brightness = 100
	}

	if cfg.Search.WebSearchEngine == "" {
		cfg.Search.WebSearchEngine = "google"
	}

	if cfg.Modules.Weather.Mode == "" {
		cfg.Modules.Weather.Mode = "current"
	}

	if cfg.Modules.Clock.Style == "" {
		cfg.Modules.Clock.Style = "minimal"
	}

	if cfg.Modules.Clock.TimeFormat == "" {
		cfg.Modules.Clock.TimeFormat = "24h"
	}

	for sectionIndex := range cfg.Sections {
		section := &cfg.Sections[sectionIndex]

		if section.Accent == "" {
			section.Accent = "#4f8cff"
		}

		switch section.Layout {
		case "":
			section.Layout = "list"
		case "large":
			section.Layout = "featured"
		}

		switch section.Width {
		case "":
			section.Width = "medium"
		case "small":
			section.Width = "narrow"
		case "large":
			section.Width = "extra-wide"
		}

		if section.CardSize == "" {
			section.CardSize = "inherit"
		}

		if section.GridColumns == 0 {
			section.GridColumns = 2
		}

		if section.Items == nil {
			section.Items = []Item{}
		}

		for itemIndex := range section.Items {
			item := &section.Items[itemIndex]

			if item.Type == "" {
				item.Type = "service"
			}

			if item.Icon.Type == "" {
				item.Icon.Type = "auto"
			}
		}
	}

	if cfg.Sections == nil {
		cfg.Sections = []Section{}
	}

	if cfg.Shortcuts == nil {
		cfg.Shortcuts = []Shortcut{}
	}

	for shortcutIndex := range cfg.Shortcuts {
		shortcut := &cfg.Shortcuts[shortcutIndex]

		if shortcut.Icon.Type == "" {
			shortcut.Icon.Type = "auto"
		}
	}
}
