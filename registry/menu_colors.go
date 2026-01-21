package registry

// MenuColor represents a sidebar color scheme option
type MenuColor struct {
	Name        string
	Label       string
	Description string
}

// MENU_COLORS contains available sidebar color schemes
var MENU_COLORS = []MenuColor{
	{Name: "default", Label: "Default", Description: "Standard sidebar colors"},
	{Name: "inverted", Label: "Inverted", Description: "Dark sidebar with light text"},
}

var menuColorMap = map[string]*MenuColor{}

func init() {
	for i := range MENU_COLORS {
		menuColorMap[MENU_COLORS[i].Name] = &MENU_COLORS[i]
	}
}

// GetMenuColor returns a menu color option by name
func GetMenuColor(name string) *MenuColor { return menuColorMap[name] }

// AllMenuColors returns all available menu color options
func AllMenuColors() []MenuColor { return MENU_COLORS }
