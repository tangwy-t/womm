package badge

import "testing"

func TestLookup(t *testing.T) {
	reg := NewRegistry()
	RegisterAll(reg)

	b, ok := reg.Lookup("midnight-coder")
	if !ok {
		t.Fatal("expected midnight-coder to exist")
	}
	if b.Type != Certified {
		t.Errorf("expected Certified, got %v", b.Type)
	}
	if b.Rarity != Rare {
		t.Errorf("expected Rare, got %v", b.Rarity)
	}
}

func TestLookupNotFound(t *testing.T) {
	reg := NewRegistry()
	_, ok := reg.Lookup("nonexistent")
	if ok {
		t.Fatal("expected not found")
	}
}

func TestListAll(t *testing.T) {
	reg := NewRegistry()
	RegisterAll(reg)
	all := reg.ListAll()
	if len(all) != 25 {
		t.Errorf("expected 25 badges, got %d", len(all))
	}
}

func TestListByType(t *testing.T) {
	reg := NewRegistry()
	RegisterAll(reg)
	d := reg.ListByType(Declarative)
	if len(d) != 10 {
		t.Errorf("expected 10 declarative, got %d", len(d))
	}
	c := reg.ListByType(Certified)
	if len(c) != 15 {
		t.Errorf("expected 15 certified, got %d", len(c))
	}
}

func TestI18n(t *testing.T) {
	reg := NewRegistry()
	RegisterAll(reg)
	b, _ := reg.Lookup("midnight-coder")
	if b.LocalizedName("zh") != "午夜编码者" {
		t.Errorf("wrong zh name: %s", b.LocalizedName("zh"))
	}
	if b.LocalizedName("en") != "Midnight Coder" {
		t.Errorf("wrong en name: %s", b.LocalizedName("en"))
	}
}

func TestI18nFallback(t *testing.T) {
	reg := NewRegistry()
	RegisterAll(reg)
	b, _ := reg.Lookup("midnight-coder")
	if b.LocalizedName("fr") != "午夜编码者" {
		t.Errorf("expected zh fallback for unknown lang, got: %s", b.LocalizedName("fr"))
	}
}
