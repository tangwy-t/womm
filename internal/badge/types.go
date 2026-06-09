package badge

type BadgeType string

const (
	Declarative BadgeType = "declarative"
	Certified   BadgeType = "certified"
)

type Rarity string

const (
	Common    Rarity = "common"
	Rare      Rarity = "rare"
	Legendary Rarity = "legendary"
)

type Badge struct {
	ID       string            `json:"id"`
	Name     map[string]string `json:"name"`
	Subtitle map[string]string `json:"subtitle"`
	Icon     string            `json:"icon"`
	Type     BadgeType         `json:"type"`
	Rarity   Rarity            `json:"rarity"`
}

func (b *Badge) LocalizedName(lang string) string {
	if name, ok := b.Name[lang]; ok {
		return name
	}
	return b.Name["zh"]
}

func (b *Badge) LocalizedSubtitle(lang string) string {
	if sub, ok := b.Subtitle[lang]; ok {
		return sub
	}
	return b.Subtitle["zh"]
}
