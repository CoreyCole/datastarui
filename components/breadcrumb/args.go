package breadcrumb

import "github.com/a-h/templ"

// BreadcrumbArgs defines the properties for the Breadcrumb nav component
type BreadcrumbArgs struct {
	// Class allows additional CSS classes to be added
	Class string

	// Attributes allows additional HTML attributes to be added
	Attributes templ.Attributes
}

// BreadcrumbListArgs defines the properties for the BreadcrumbList ol component
type BreadcrumbListArgs struct {
	// Class allows additional CSS classes to be added
	Class string

	// Attributes allows additional HTML attributes to be added
	Attributes templ.Attributes
}

// BreadcrumbItemArgs defines the properties for the BreadcrumbItem li component
type BreadcrumbItemArgs struct {
	// Class allows additional CSS classes to be added
	Class string

	// Attributes allows additional HTML attributes to be added
	Attributes templ.Attributes
}

// BreadcrumbLinkArgs defines the properties for the BreadcrumbLink a component
type BreadcrumbLinkArgs struct {
	// AsChild renders the link as a child element (for composition)
	AsChild bool

	// Class allows additional CSS classes to be added
	Class string

	// Attributes allows additional HTML attributes to be added
	Attributes templ.Attributes

	// Href specifies the URL for the link
	Href string
}

// BreadcrumbPageArgs defines the properties for the BreadcrumbPage span component
type BreadcrumbPageArgs struct {
	// Class allows additional CSS classes to be added
	Class string

	// Attributes allows additional HTML attributes to be added
	Attributes templ.Attributes
}

// BreadcrumbSeparatorArgs defines the properties for the BreadcrumbSeparator li component
type BreadcrumbSeparatorArgs struct {
	// Class allows additional CSS classes to be added
	Class string

	// Attributes allows additional HTML attributes to be added
	Attributes templ.Attributes

	// CustomIcon allows overriding the default ChevronRight icon
	CustomIcon bool
}

// BreadcrumbEllipsisArgs defines the properties for the BreadcrumbEllipsis span component
type BreadcrumbEllipsisArgs struct {
	// Class allows additional CSS classes to be added
	Class string

	// Attributes allows additional HTML attributes to be added
	Attributes templ.Attributes
}
