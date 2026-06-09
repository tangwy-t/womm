package render

type Theme struct {
	Name        string
	BgColor     string
	FgColor     string
	AccentColor string
	BorderColor string
	DimColor    string
	FontFamily  string
	TitleSize   int
	SubSize     int
}

var themes = map[string]Theme{
	"pixel": {
		Name: "pixel", BgColor: "#1a1a2e", FgColor: "#00ff41",
		AccentColor: "#00ff41", BorderColor: "#00ff41", DimColor: "#336633",
		FontFamily: "monospace", TitleSize: 11, SubSize: 8,
	},
	"cyberpunk": {
		Name: "cyberpunk", BgColor: "#0d0221", FgColor: "#05d9e8",
		AccentColor: "#ff2a6d", BorderColor: "#d300c5", DimColor: "#3d1a5e",
		FontFamily: "monospace", TitleSize: 11, SubSize: 8,
	},
	"glitch": {
		Name: "glitch", BgColor: "#111111", FgColor: "#ffffff",
		AccentColor: "#ff3333", BorderColor: "#666666", DimColor: "#888888",
		FontFamily: "monospace", TitleSize: 13, SubSize: 8,
	},
	"clean": {
		Name: "clean", BgColor: "#ffffff", FgColor: "#333333",
		AccentColor: "#e056a0", BorderColor: "#dddddd", DimColor: "#999999",
		FontFamily: "sans-serif", TitleSize: 11, SubSize: 8,
	},
}

func GetTheme(name string) Theme {
	if t, ok := themes[name]; ok {
		return t
	}
	return themes["pixel"]
}

type TierConfig struct {
	Name      string
	Accent    string
	Glow      string
	RingLight string
	RingDark  string
	Edge      string
}

var tiers = map[string]TierConfig{
	"bronze": {
		Name:      "bronze",
		Accent:    "#cd7f32",
		Glow:      "#e0a070",
		RingLight: "#e8b888",
		RingDark:  "#8b4513",
		Edge:      "#3a2818",
	},
	"silver": {
		Name:      "silver",
		Accent:    "#c0c0c0",
		Glow:      "#e0e0e0",
		RingLight: "#e0e0e0",
		RingDark:  "#808080",
		Edge:      "#1a1a1a",
	},
	"gold": {
		Name:      "gold",
		Accent:    "#ffd700",
		Glow:      "#ffed4e",
		RingLight: "#ffed4e",
		RingDark:  "#b8860b",
		Edge:      "#3a3000",
	},
}

func GetTier(rarity string) TierConfig {
	if t, ok := tiers[rarity]; ok {
		return t
	}
	return tiers["bronze"]
}
