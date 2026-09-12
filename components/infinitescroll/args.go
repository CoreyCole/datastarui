package infinitescroll

import "github.com/a-h/templ"

// InfiniteScrollArgs defines the properties for the InfiniteScroll container
type InfiniteScrollArgs struct {
	ID             string // Required: snake_case identifier
	PatchAboveExpr string // Optional: expression for loading above/scrollback (omit = no above edge)
	PatchBelowExpr string // Optional: expression for loading below/forward (omit = no below edge)
	// MorphMap overrides - empty = derived from ID
	HostID            string
	ItemsID           string
	SentinelAboveID   string
	SentinelBelowID   string
	LoadingAboveID    string
	LoadingBelowID    string
	Class             string
	Attributes        templ.Attributes
}

// Direction represents the scroll direction
type Direction string

const (
	DirectionAbove Direction = "above"
	DirectionBelow Direction = "below"
)

// HostArgs defines properties for the Host container
type HostArgs struct {
	ID         string
	Class      string
	Attributes templ.Attributes
}

// ItemsArgs defines properties for the Items container
type ItemsArgs struct {
	ID         string
	Class      string
	Attributes templ.Attributes
}

// SentinelArgs defines properties for sentinel elements
type SentinelArgs struct {
	ID          string
	Direction   Direction
	PatchExpr   string // The @get() or @post() expression to load more
	Class       string
	Attributes  templ.Attributes
}

// LoadingArgs defines properties for loading indicators
type LoadingArgs struct {
	ID         string
	Direction  Direction
	Class      string
	Attributes templ.Attributes
}

// LoadingSentinelArgs combines loading and sentinel for convenience
type LoadingSentinelArgs struct {
	ID         string
	Direction  Direction
	PatchExpr  string
	LoadingID  string
	Class      string
	Attributes templ.Attributes
}
