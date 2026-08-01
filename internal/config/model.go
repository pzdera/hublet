package config

type Config struct {
	Version    int         `json:"version"`
	Dashboard  Dashboard   `json:"dashboard"`
	Appearance Appearance  `json:"appearance"`
	Search     Search      `json:"search"`
	Sections   []Section   `json:"sections"`
	Shortcuts  []Shortcut  `json:"shortcuts"`
	Weather    Weather     `json:"weather"`
}

type Dashboard struct {
	Title     string  `json:"title"`
	Theme     string  `json:"theme"`
	Wallpaper *string `json:"wallpaper"`
}

type Appearance struct {
	Density   string `json:"density"`
	CardStyle string `json:"cardStyle"`
	Radius    string `json:"radius"`
}

type Search struct {
	AutoFocus           bool   `json:"autoFocus"`
	OpenShortcutDirectly bool   `json:"openShortcutDirectly"`
	WebSearchEnabled    bool   `json:"webSearchEnabled"`
	WebSearchEngine     string `json:"webSearchEngine"`
}

type Section struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Accent    string `json:"accent"`
	Layout    string `json:"layout"`
	Collapsed bool   `json:"collapsed"`
	Items     []Item `json:"items"`
}

type Item struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	URL          string `json:"url"`
	Description  string `json:"description"`
	Icon         Icon   `json:"icon"`
	OpenInNewTab bool   `json:"openInNewTab"`
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

type Weather struct {
	Enabled   bool     `json:"enabled"`
	Location  string   `json:"location"`
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
}

func Default() Config {
	return Config{
		Version: 2,
		Dashboard: Dashboard{
			Title: "Hublet",
			Theme: "midnight",
		},
		Appearance: Appearance{
			Density:   "comfortable",
			CardStyle: "glass",
			Radius:    "large",
		},
		Search: Search{
			AutoFocus:            true,
			OpenShortcutDirectly: true,
			WebSearchEnabled:     true,
			WebSearchEngine:      "google",
		},
		Sections: []Section{
			{
				ID:        "welcome",
				Title:     "Getting started",
				Accent:    "#4f8cff",
				Layout:    "grid",
				Collapsed: false,
				Items: []Item{
					{
						ID:           "github",
						Name:         "GitHub",
						URL:          "https://github.com",
						Description:  "Code hosting",
						Icon:         Icon{Type: "auto"},
						OpenInNewTab: true,
					},
					{
						ID:           "youtube",
						Name:         "YouTube",
						URL:          "https://youtube.com",
						Description:  "Video platform",
						Icon:         Icon{Type: "auto"},
						OpenInNewTab: true,
					},
				},
			},
		},
		Shortcuts: []Shortcut{
			{
				Key:   "kp",
				Label: "KupujemProdajem",
				URL:   "https://www.kupujemprodajem.com",
				Icon:  Icon{Type: "auto"},
			},
			{
				Key:   "yt",
				Label: "YouTube",
				URL:   "https://www.youtube.com",
				Icon:  Icon{Type: "auto"},
			},
		},
		Weather: Weather{
			Enabled: false,
		},
	}
}
