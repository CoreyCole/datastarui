package layouts

// RootArgs defines the arguments for the Root layout
type RootArgs struct {
	CurrentPage          string // The current page section (e.g., "components", "docs")
	CurrentPath          string // The actual URL path (e.g., "/components/button")
	InspectorEnabled     bool   // Whether to show the Datastar inspector (from env var)
	DatastarProAvailable bool   // Whether the Datastar Pro file exists locally
	StyleClass           string // Style class applied to body (e.g., "style-vega") - if empty, uses saved theme
}
