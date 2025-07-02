package sidebar

import "github.com/a-h/templ"

type SidebarItem struct {
	Title string
	Href  string
	Items []SidebarItem
	Label string // For "New" badges etc
}

type SidebarSection struct {
	Title string
	Items []SidebarItem
	Label string // For "New" badges etc
}

type SidebarProps struct {
	Sections    []SidebarSection
	CurrentPath string
	Class       string
	Attributes  templ.Attributes
}
