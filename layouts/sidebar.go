package layouts

import "github.com/coreycole/datastarui/components/sidebar"

func GetSidebarSections() []sidebar.SidebarSection {
	return []sidebar.SidebarSection{
		{
			Title: "Getting Started",
			Items: []sidebar.SidebarItem{
				{Title: "Create Theme", Href: "/create"},
			},
		},
		{
			Title: "Components",
			Items: []sidebar.SidebarItem{
				{Title: "Breadcrumb", Href: "/components/breadcrumb"},
				{Title: "Button", Href: "/components/button"},
				{Title: "Calendar", Href: "/components/calendar"},
				{Title: "Card", Href: "/components/card"},
				{Title: "Checkbox", Href: "/components/checkbox"},
				{Title: "Datepicker", Href: "/components/datepicker"},
				{Title: "Dialog", Href: "/components/dialog"},
				{Title: "Dropdown", Href: "/components/dropdown"},
				{Title: "Form", Href: "/components/form"},
				{Title: "Popover", Href: "/components/popover"},
				{Title: "Select", Href: "/components/select"},
				{Title: "Sheet", Href: "/components/sheet"},
				{Title: "Sidebar", Href: "/components/sidebar"},
				{Title: "Tabs", Href: "/components/tabs"},
				{Title: "Tooltip", Href: "/components/tooltip"},
			},
		},
	}
}
