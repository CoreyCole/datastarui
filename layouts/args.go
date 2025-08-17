package layouts

// RootArgs defines the arguments for the Root layout
type RootArgs struct {
	CurrentPage string // The current page section (e.g., "components", "docs")
	CurrentPath string // The actual URL path (e.g., "/components/button")
}