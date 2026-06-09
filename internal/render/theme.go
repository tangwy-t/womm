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
