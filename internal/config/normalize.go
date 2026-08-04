package config

func Normalize(cfg *Config) {
	if cfg.Version == 0 {
		cfg.Version = 2
	}

	if cfg.Dashboard.Title == "" ||
		cfg.Dashboard.Title == "Hublet" ||
		cfg.Dashboard.Title == "Hublet v2" {
		cfg.Dashboard.Title = "Hublet"
	}

	if cfg.Dashboard.Description == "" {
		cfg.Dashboard.Description = "My self-hosted dashboard"
		cfg.Dashboard.DescriptionVisible = true
	}

	if cfg.Dashboard.Icon.Type == "" {
		cfg.Dashboard.Icon.Type = "initial"
	}

	if cfg.Dashboard.Icon.Type == "initial" &&
		cfg.Dashboard.Icon.Value == "" {
		cfg.Dashboard.Icon.Value = "H"
	}

	if cfg.Dashboard.IconSize == "" {
		cfg.Dashboard.IconSize = "medium"
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

	if cfg.Appearance.Cards.Radius == "" {
		cfg.Appearance.Cards.Radius = "large"
	}

	if cfg.Appearance.Cards.Shadow == "" {
		cfg.Appearance.Cards.Shadow = "soft"
	}

	if cfg.Appearance.Background.Type == "" {
		cfg.Appearance.Background.Type = "solid"
	}

	if cfg.Appearance.Background.Type == "gradient" {
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

		if section.Surface == "" {
			section.Surface = "solid"
			section.SurfaceOpacity = 82
			section.SurfaceBlur = 16
			section.ShowBorder = true
		}

		switch section.LegacyLayout {
		case "list", "featured", "large":
			section.GridColumns = 1
		}

		section.LegacyLayout = ""

		switch section.Width {
		case "":
			section.Width = "medium"

		case "small":
			section.Width = "narrow"

		case "large":
			section.Width = "extra-wide"
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

	normalizeSectionGrid(cfg.Sections)

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

func normalizeSectionGrid(
	sections []Section,
) {
	occupied := map[int][24]bool{}

	for sectionIndex := range sections {
		section := &sections[sectionIndex]
		legacyGrid := section.GridColumnSpan == 0

		if legacyGrid {
			section.GridColumnSpan =
				legacySectionGridSpan(section.Width)

			if section.GridColumn > 0 {
				section.GridColumn =
					(section.GridColumn-1)*2 + 1
			} else if section.StartColumn > 0 {
				section.GridColumn =
					(section.StartColumn-1)*2 + 1
			}
		}

		if section.GridColumnSpan < 4 ||
			section.GridColumnSpan > 24 {
			section.GridColumnSpan =
				legacySectionGridSpan(section.Width)
		}

		if section.GridRowSpan < 1 {
			section.GridRowSpan = 1
		}

		if !gridPlacementFits(
			occupied,
			section.GridRow,
			section.GridColumn,
			section.GridRowSpan,
			section.GridColumnSpan,
		) {
			section.GridRow,
				section.GridColumn =
				firstGridPlacement(
					occupied,
					section.GridRowSpan,
					section.GridColumnSpan,
				)
		}

		occupyGridPlacement(
			occupied,
			section.GridRow,
			section.GridColumn,
			section.GridRowSpan,
			section.GridColumnSpan,
		)

		section.StartColumn = 0
	}
}

func legacySectionGridSpan(
	width string,
) int {
	switch width {
	case "narrow":
		return 6
	case "medium":
		return 8
	case "wide":
		return 12
	case "extra-wide":
		return 16
	case "full":
		return 24
	default:
		return 8
	}
}

func gridPlacementFits(
	occupied map[int][24]bool,
	row int,
	column int,
	rowSpan int,
	columnSpan int,
) bool {
	if row < 1 ||
		column < 1 ||
		rowSpan < 1 ||
		columnSpan < 4 ||
		column+columnSpan-1 > 24 {
		return false
	}

	for gridRow := row;
		gridRow < row+rowSpan;
		gridRow++ {
		cells := occupied[gridRow]

		for cell := column - 1;
			cell < column-1+columnSpan;
			cell++ {
			if cells[cell] {
				return false
			}
		}
	}

	return true
}

func firstGridPlacement(
	occupied map[int][24]bool,
	rowSpan int,
	columnSpan int,
) (int, int) {
	maximumRow := 1

	for row := range occupied {
		if row > maximumRow {
			maximumRow = row
		}
	}

	for row := 1; row <= maximumRow+1; row++ {
		for column := 1;
			column+columnSpan-1 <= 24;
			column++ {
			if gridPlacementFits(
				occupied,
				row,
				column,
				rowSpan,
				columnSpan,
			) {
				return row, column
			}
		}
	}

	return maximumRow + 1, 1
}

func occupyGridPlacement(
	occupied map[int][24]bool,
	row int,
	column int,
	rowSpan int,
	columnSpan int,
) {
	for gridRow := row;
		gridRow < row+rowSpan;
		gridRow++ {
		cells := occupied[gridRow]

		for cell := column - 1;
			cell < column-1+columnSpan;
			cell++ {
			cells[cell] = true
		}

		occupied[gridRow] = cells
	}
}
