package layouts

import "github.com/coreycole/datastarui/components/sidebar"

func getSidebarSections() []sidebar.SidebarSection {
	return []sidebar.SidebarSection{
		{
			Title: "Components",
			Items: []sidebar.SidebarItem{
				{Title: "Breadcrumb", Href: "/components/breadcrumb"},
				{Title: "Button", Href: "/components/button"},
				{Title: "Card", Href: "/components/card"},
				{Title: "Checkbox", Href: "/components/checkbox"},
				{Title: "Dropdown", Href: "/components/dropdown"},
				{Title: "Form", Href: "/components/form"},
				{Title: "Tabs", Href: "/components/tabs"},
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
