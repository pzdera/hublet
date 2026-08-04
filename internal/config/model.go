package config

type Config struct {
	Version    int        `json:"version"`
	Dashboard  Dashboard  `json:"dashboard"`
	Appearance Appearance `json:"appearance"`
	Search     Search     `json:"search"`
	Sections   []Section  `json:"sections"`
	Shortcuts  []Shortcut `json:"shortcuts"`
	Modules    Modules    `json:"modules"`
}

type Dashboard struct {
	Title              string        `json:"title"`
	Description        string        `json:"description"`
	DescriptionVisible bool          `json:"descriptionVisible"`
	Icon               DashboardIcon `json:"icon"`
	IconSize           string        `json:"iconSize"`
	Theme              string        `json:"theme"`
	Wallpaper          *string       `json:"wallpaper"`
}

type DashboardIcon struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Appearance struct {
	Mode       string     `json:"mode"`
	Font       Font       `json:"font"`
	Cards      Cards      `json:"cards"`
	Background Background `json:"background"`
	Animations bool       `json:"animations"`
}

type Font struct {
	Family string `json:"family"`
	Scale  string `json:"scale"`
}

type Cards struct {
	Size    string `json:"size"`
	Density string `json:"density"`
	Radius  string `json:"radius"`
	Shadow  string `json:"shadow"`
	Border  bool   `json:"border"`
}

type Background struct {
	Type       string  `json:"type"`
	Color      string  `json:"color"`
	Blur       int     `json:"blur"`
	Brightness int     `json:"brightness"`
	Overlay    int     `json:"overlay"`
	Wallpaper  *string `json:"wallpaper"`
}

type Search struct {
	AutoFocus            bool   `json:"autoFocus"`
	OpenShortcutDirectly bool   `json:"openShortcutDirectly"`
	WebSearchEnabled     bool   `json:"webSearchEnabled"`
	WebSearchEngine      string `json:"webSearchEngine"`
}

type Section struct {
	ID             string `json:"id"`
	Title          string `json:"title"`
	Accent         string `json:"accent"`
	Surface        string `json:"surface"`
	SurfaceOpacity int    `json:"surfaceOpacity"`
	SurfaceBlur    int    `json:"surfaceBlur"`
	ShowBorder     bool   `json:"showBorder"`
	Layout         string `json:"layout"`
	Width          string `json:"width"`
	GridRow        int    `json:"gridRow"`
	GridColumn     int    `json:"gridColumn"`
	GridRowSpan    int    `json:"gridRowSpan"`
	GridColumnSpan int    `json:"gridColumnSpan"`
	StartColumn    int    `json:"startColumn,omitempty"`
	CardSize       string `json:"cardSize"`
	GridColumns    int    `json:"gridColumns"`
	FillLastRow    bool   `json:"fillLastRow"`
	Collapsed      bool   `json:"collapsed"`
	Items          []Item `json:"items"`
}

type Item struct {
	ID           string           `json:"id"`
	Type         string           `json:"type"`
	Name         string           `json:"name"`
	URL          string           `json:"url"`
	Description  string           `json:"description"`
	Icon         Icon             `json:"icon"`
	OpenInNewTab bool             `json:"openInNewTab"`
	Resources    ServiceResources `json:"resources"`
}

type ServiceResources struct {
	Enabled    bool `json:"enabled"`
	ShowStatus bool `json:"showStatus"`
	ShowCPU    bool `json:"showCpu"`
	ShowMemory bool `json:"showMemory"`
}

type Icon struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Shortcut struct {
	Key   string `json:"key"`
	Label string `json:"label"`
	URL   string `json:"url"`
	Icon  Icon   `json:"icon"`
}

type Modules struct {
	Weather    WeatherModule    `json:"weather"`
	Clock      ClockModule      `json:"clock"`
	Statistics StatisticsModule `json:"statistics"`
}

type WeatherModule struct {
	Enabled   bool     `json:"enabled"`
	Mode      string   `json:"mode"`
	Location  string   `json:"location"`
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
}

type ClockModule struct {
	Enabled    bool   `json:"enabled"`
	Style      string `json:"style"`
	TimeFormat string `json:"timeFormat"`
	ShowDate   bool   `json:"showDate"`
}

type StatisticsModule struct {
	Enabled bool `json:"enabled"`
}

func Default() Config {
	return Config{
		Version: 2,

		Dashboard: Dashboard{
			Title:              "Hublet v2",
			Description:        "My self-hosted dashboard",
			DescriptionVisible: true,
			Icon: DashboardIcon{
				Type:  "initial",
				Value: "H",
			},
			IconSize: "medium",
			Theme:    "midnight",
		},

		Appearance: Appearance{
			Mode: "minimal",

			Font: Font{
				Family: "system",
				Scale:  "medium",
			},

			Cards: Cards{
				Size:    "medium",
				Density: "comfortable",
				Radius:  "large",
				Shadow:  "soft",
				Border:  true,
			},

			Background: Background{
				Type:       "solid",
				Color:      "#090c12",
				Blur:       0,
				Brightness: 100,
				Overlay:    0,
			},

			Animations: true,
		},

		Search: Search{
			AutoFocus:            true,
			OpenShortcutDirectly: true,
			WebSearchEnabled:     true,
			WebSearchEngine:      "google",
		},

		Sections:  []Section{},
		Shortcuts: []Shortcut{},

		Modules: Modules{
			Weather: WeatherModule{
				Enabled: false,
				Mode:    "current",
			},

			Clock: ClockModule{
				Enabled:    false,
				Style:      "minimal",
				TimeFormat: "24h",
				ShowDate:   true,
			},

			Statistics: StatisticsModule{
				Enabled: false,
			},
		},
	}
}
