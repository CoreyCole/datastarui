package registry

// MenuAccent represents a sidebar accent styling option
type MenuAccent struct {
	Name        string
	Label       string
	Description string
}

// MENU_ACCENTS contains available sidebar accent styles
var MENU_ACCENTS = []MenuAccent{
	{Name: "subtle", Label: "Subtle", Description: "Muted accent colors"},
	{Name: "bold", Label: "Bold", Description: "Primary color as accent"},
}

var menuAccentMap = map[string]*MenuAccent{}

func init() {
	for i := range MENU_ACCENTS {
		menuAccentMap[MENU_ACCENTS[i].Name] = &MENU_ACCENTS[i]
	}
}

// GetMenuAccent returns a menu accent option by name
func GetMenuAccent(name string) *MenuAccent { return menuAccentMap[name] }

// AllMenuAccents returns all available menu accent options
func AllMenuAccents() []MenuAccent { return MENU_ACCENTS }
