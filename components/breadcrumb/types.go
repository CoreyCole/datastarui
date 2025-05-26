package breadcrumb

import "github.com/a-h/templ"

// BreadcrumbProps defines the properties for the Breadcrumb nav component
type BreadcrumbProps struct {
	// Class allows additional CSS classes to be added
	Class string

	// Attributes allows additional HTML attributes to be added
	Attributes templ.Attributes
}

// BreadcrumbListProps defines the properties for the BreadcrumbList ol component
type BreadcrumbListProps struct {
	// Class allows additional CSS classes to be added
	Class string

	// Attributes allows additional HTML attributes to be added
	Attributes templ.Attributes
}

// BreadcrumbItemProps defines the properties for the BreadcrumbItem li component
type BreadcrumbItemProps struct {
	// Class allows additional CSS classes to be added
	Class string

	// Attributes allows additional HTML attributes to be added
	Attributes templ.Attributes
}

// BreadcrumbLinkProps defines the properties for the BreadcrumbLink a component
type BreadcrumbLinkProps struct {
	// AsChild renders the link as a child element (for composition)
	AsChild bool

	// Class allows additional CSS classes to be added
	Class string

	// Attributes allows additional HTML attributes to be added
	Attributes templ.Attributes

	// Href specifies the URL for the link
	Href string
}

// BreadcrumbPageProps defines the properties for the BreadcrumbPage span component
type BreadcrumbPageProps struct {
	// Class allows additional CSS classes to be added
	Class string

	// Attributes allows additional HTML attributes to be added
	Attributes templ.Attributes
}

// BreadcrumbSeparatorProps defines the properties for the BreadcrumbSeparator li component
type BreadcrumbSeparatorProps struct {
	// Class allows additional CSS classes to be added
	Class string

	// Attributes allows additional HTML attributes to be added
	Attributes templ.Attributes

	// CustomIcon allows overriding the default ChevronRight icon
	CustomIcon bool
}

// BreadcrumbEllipsisProps defines the properties for the BreadcrumbEllipsis span component
type BreadcrumbEllipsisProps struct {
	// Class allows additional CSS classes to be added
	Class string

	// Attributes allows additional HTML attributes to be added
	Attributes templ.Attributes
}
