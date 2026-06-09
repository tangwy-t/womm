package render

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"strings"
	"unicode/utf8"

	"github.com/womm/womm/internal/badge"
)

//go:embed templates/*.svg.tmpl
var templateFiles embed.FS

type Renderer struct {
	templates map[string]*template.Template
}

type tmplData struct {
	Width       int
	RectWidth   int
	Theme       Theme
	Icon        template.HTML
	Name        string
	Subtitle    string
	RarityLabel string
}

func NewRenderer() *Renderer {
	r := &Renderer{templates: make(map[string]*template.Template)}
	for _, name := range []string{"badge", "wide", "terminal", "stamp"} {
		data, err := templateFiles.ReadFile("templates/" + name + ".svg.tmpl")
		if err != nil {
			continue
		}
		tmpl, err := template.New(name).Parse(string(data))
		if err != nil {
			continue
		}
		r.templates[name] = tmpl
	}
	return r
}

func (r *Renderer) Render(b *badge.Badge, themeName, templateName, lang string) (string, error) {
	tmpl, ok := r.templates[templateName]
	if !ok {
		return "", fmt.Errorf("unknown template: %s", templateName)
	}
	theme := GetTheme(themeName)
	name := b.LocalizedName(lang)
	subtitle := b.LocalizedSubtitle(lang)
	charLen := utf8.RuneCountInString(name)
	width := 30 + charLen*theme.TitleSize + 20
	if width < 120 {
		width = 120
	}
	data := tmplData{
		Width:       width,
		RectWidth:   width - 2,
		Theme:       theme,
		Icon:        template.HTML(GetIcon(b.Icon)),
		Name:        name,
		Subtitle:    subtitle,
		RarityLabel: strings.ToUpper(string(b.Rarity)),
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("template execute: %w", err)
	}
	return buf.String(), nil
}
