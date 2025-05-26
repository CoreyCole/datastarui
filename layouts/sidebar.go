package layouts

import "github.com/coreycole/datastarui/components/sidebar"

func getSidebarSections() []sidebar.SidebarSection {
	return []sidebar.SidebarSection{
		{
			Title: "Components",
			Items: []sidebar.SidebarItem{
				{Title: "Breadcrumb", Href: "/components/breadcrumb"},
				{Title: "Button", Href: "/components/button"},
				{Title: "Dropdown", Href: "/components/dropdown"},
			},
		},
	}
}

func getCurrentPath(currentPage string) string {
	switch currentPage {
	case "docs":
		return "/docs"
	case "components":
		return "/components"
	default:
		return "/"
	}
}
