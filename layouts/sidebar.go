package layouts

import "github.com/coreycole/datastarui/components/sidebar"

func getSidebarSections() []sidebar.SidebarSection {
	return []sidebar.SidebarSection{
		{
			Title: "Components",
			Items: []sidebar.SidebarItem{
				{Title: "Accordion", Href: "/components/accordion"},
				{Title: "Alert", Href: "/components/alert"},
				{Title: "Alert Dialog", Href: "/components/alert-dialog"},
				{Title: "Aspect Ratio", Href: "/components/aspect-ratio"},
				{Title: "Avatar", Href: "/components/avatar"},
				{Title: "Badge", Href: "/components/badge"},
				{Title: "Breadcrumb", Href: "/components/breadcrumb"},
				{Title: "Button", Href: "/components/button"},
				{Title: "Calendar", Href: "/components/calendar"},
				{Title: "Card", Href: "/components/card"},
				{Title: "Carousel", Href: "/components/carousel"},
				{Title: "Chart", Href: "/components/chart"},
				{Title: "Checkbox", Href: "/components/checkbox"},
				{Title: "Collapsible", Href: "/components/collapsible"},
				{Title: "Combobox", Href: "/components/combobox"},
				{Title: "Command", Href: "/components/command"},
				{Title: "Context Menu", Href: "/components/context-menu"},
				{Title: "Data Table", Href: "/components/data-table"},
				{Title: "Date Picker", Href: "/components/date-picker"},
				{Title: "Dialog", Href: "/components/dialog"},
				{Title: "Drawer", Href: "/components/drawer"},
				{Title: "Dropdown Menu", Href: "/components/dropdown-menu"},
				{Title: "Form", Href: "/components/form"},
				{Title: "Hover Card", Href: "/components/hover-card"},
				{Title: "Input", Href: "/components/input"},
				{Title: "Input OTP", Href: "/components/input-otp"},
				{Title: "Label", Href: "/components/label"},
				{Title: "Menubar", Href: "/components/menubar"},
				{Title: "Navigation Menu", Href: "/components/navigation-menu"},
				{Title: "Pagination", Href: "/components/pagination"},
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
