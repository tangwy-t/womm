package render

import (
	"fmt"
	"strings"
	"testing"

	"github.com/womm/womm/internal/badge"
)

func testBadge() *badge.Badge {
	return &badge.Badge{
		ID:       "test",
		Name:     map[string]string{"zh": "测试徽章", "en": "Test Badge"},
		Subtitle: map[string]string{"zh": "测试副标题", "en": "Subtitle"},
		Icon:     "checkmark",
		Type:     badge.Declarative,
		Rarity:   badge.Bronze,
	}
}

func TestRenderBadgeSVG(t *testing.T) {
	r := NewRenderer()
	svg, err := r.Render(testBadge(), "pixel", "badge", "zh")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(svg, "<svg") {
		t.Error("missing svg tag")
	}
	if !strings.Contains(svg, "测试徽章") {
		t.Error("missing badge name")
	}
}

func TestRenderEnglish(t *testing.T) {
	r := NewRenderer()
	svg, err := r.Render(testBadge(), "pixel", "badge", "en")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(svg, "Test Badge") {
		t.Error("missing English name")
	}
}

func TestRenderAllThemes(t *testing.T) {
	r := NewRenderer()
	b := testBadge()
	for _, theme := range []string{"pixel", "cyberpunk", "glitch", "clean"} {
		t.Run(theme, func(t *testing.T) {
			svg, err := r.Render(b, theme, "badge", "zh")
			if err != nil {
				t.Errorf("theme %s failed: %v", theme, err)
			}
			if !strings.Contains(svg, "<svg") {
				t.Error("no svg output")
			}
		})
	}
}

func TestRenderAllTemplates(t *testing.T) {
	r := NewRenderer()
	b := testBadge()
	for _, tmpl := range []string{"badge", "wide", "terminal", "stamp", "github"} {
		t.Run(tmpl, func(t *testing.T) {
			svg, err := r.Render(b, "pixel", tmpl, "zh")
			if err != nil {
				t.Errorf("template %s failed: %v", tmpl, err)
			}
			if !strings.Contains(svg, "<svg") {
				t.Error("no svg output")
			}
		})
	}
}

func TestRenderGithubTemplateAllTiers(t *testing.T) {
	r := NewRenderer()
	for _, rarity := range []badge.Rarity{badge.Bronze, badge.Silver, badge.Gold} {
		t.Run(string(rarity), func(t *testing.T) {
			b := &badge.Badge{
				ID:       "test-" + string(rarity),
				Name:     map[string]string{"zh": "测试", "en": "Test"},
				Subtitle: map[string]string{"zh": "副标题", "en": "Subtitle"},
				Icon:     "moon",
				Type:     badge.Certified,
				Rarity:   rarity,
			}
			svg, err := r.Render(b, "pixel", "github", "zh")
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(svg, "<svg") {
				t.Error("missing svg")
			}
			if !strings.Contains(svg, "radialGradient") {
				t.Error("missing glow effect")
			}
			if !strings.Contains(svg, fmt.Sprintf("bg-glow-test-%s", rarity)) {
				t.Errorf("missing tier-specific gradient for %s", rarity)
			}
			// github template should be pure icon medallion - no text
			if strings.Contains(svg, "<text") {
				t.Errorf("github template should not contain text elements for %s tier", rarity)
			}
			// Should be 160x160 square
			if !strings.Contains(svg, `width="160" height="160"`) {
				t.Errorf("github template should be 160x160 square for %s tier", rarity)
			}
			// Should have metallic ring gradient
			if !strings.Contains(svg, fmt.Sprintf("ring-test-%s", rarity)) {
				t.Errorf("missing metallic ring gradient for %s tier", rarity)
			}
			// Gold tier should have specular highlight overlay
			if rarity == badge.Gold {
				if !strings.Contains(svg, "ellipse") {
					t.Error("gold tier should have specular highlight overlay")
				}
			} else {
				if strings.Contains(svg, "ellipse") {
					t.Errorf("%s tier should not have specular highlight", rarity)
				}
			}
		})
	}
}
