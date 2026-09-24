package core

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type KeyGroup struct {
	Name        string
	PrivatePath string
	PublicPath  string
}

func (g *KeyGroup) HasPrivate() bool { return g.PrivatePath != "" }
func (g *KeyGroup) HasPublic() bool  { return g.PublicPath != "" }

func (g *KeyGroup) Display() string {
	tag := ""
	switch {
	case g.HasPrivate() && g.HasPublic():
		tag = " [私+公]"
	case g.HasPrivate():
		tag = " [仅私钥]"
	case g.HasPublic():
		tag = " [仅公钥]"
	}
	return g.Name + tag
}

func ListSSHGroups(dir string) ([]KeyGroup, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	groups := map[string]*KeyGroup{}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		path := filepath.Join(dir, name)
		if strings.HasSuffix(name, ".pub") {
			base := strings.TrimSuffix(name, ".pub")
			g := groups[base]
			if g == nil {
				g = &KeyGroup{Name: base}
				groups[base] = g
			}
			g.PublicPath = path
		} else {
			g := groups[name]
			if g == nil {
				g = &KeyGroup{Name: name}
				groups[name] = g
			}
			g.PrivatePath = path
		}
	}
	return sortGroups(groups), nil
}

func ListGPGGroups(dir string) ([]KeyGroup, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	groups := map[string]*KeyGroup{}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		path := filepath.Join(dir, name)
		if strings.HasSuffix(name, "-public.asc") {
			base := strings.TrimSuffix(name, "-public.asc")
			g := groups[base]
			if g == nil {
				g = &KeyGroup{Name: base}
				groups[base] = g
			}
			g.PublicPath = path
		} else if strings.HasSuffix(name, "-private.asc") {
			base := strings.TrimSuffix(name, "-private.asc")
			g := groups[base]
			if g == nil {
				g = &KeyGroup{Name: base}
				groups[base] = g
			}
			g.PrivatePath = path
		} else if strings.HasSuffix(name, ".asc") {
			base := strings.TrimSuffix(name, ".asc")
			g := groups[base]
			if g == nil {
				g = &KeyGroup{Name: base}
				groups[base] = g
			}
			if g.PublicPath == "" {
				g.PublicPath = path
			}
		}
	}
	return sortGroups(groups), nil
}

func sortGroups(m map[string]*KeyGroup) []KeyGroup {
	var out []KeyGroup
	for _, g := range m {
		out = append(out, *g)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Name < out[j].Name
	})
	return out
}
