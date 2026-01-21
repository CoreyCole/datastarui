package registry

// IconLibrary represents an icon library option
type IconLibrary struct {
	Name        string
	Title       string
	Description string
}

// ICON_LIBRARIES contains available icon libraries
var ICON_LIBRARIES = []IconLibrary{
	{Name: "lucide", Title: "Lucide", Description: "Clean, consistent line icons"},
	{Name: "phosphor", Title: "Phosphor", Description: "Flexible icon family with multiple weights"},
	{Name: "hugeicons", Title: "HugeIcons", Description: "Modern, comprehensive icon set"},
	{Name: "tabler", Title: "Tabler", Description: "Consistent stroke-based icons"},
}

var iconLibraryMap = map[string]*IconLibrary{}

func init() {
	for i := range ICON_LIBRARIES {
		iconLibraryMap[ICON_LIBRARIES[i].Name] = &ICON_LIBRARIES[i]
	}
}

// GetIconLibrary returns an icon library by name
func GetIconLibrary(name string) *IconLibrary { return iconLibraryMap[name] }

// AllIconLibraries returns all available icon libraries
func AllIconLibraries() []IconLibrary { return ICON_LIBRARIES }
