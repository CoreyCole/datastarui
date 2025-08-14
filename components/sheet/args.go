package sheet

import "github.com/a-h/templ"

// SheetArgs defines the args for the Sheet container
type SheetArgs struct {
	ID          string
	DefaultOpen bool   // Whether the sheet should be open by default
	Modal       bool   // Whether the sheet is modal (dismisses on backdrop click) or non-modal
	Side        string // Position of the sheet: "top", "right", "bottom", "left" (default: "right")
	Class       string
	Attributes  templ.Attributes
}

// SheetTriggerArgs defines the args for the SheetTrigger component
type SheetTriggerArgs struct {
	SheetID    string
	AsChild    bool
	Class      string
	Attributes templ.Attributes
}

// SheetContentArgs defines the args for the SheetContent component
type SheetContentArgs struct {
	SheetID    string // ID of the parent sheet for close button
	Class      string
	Attributes templ.Attributes
}

// SheetHeaderArgs defines the args for the SheetHeader component
type SheetHeaderArgs struct {
	Class      string
	Attributes templ.Attributes
}

// SheetFooterArgs defines the args for the SheetFooter component
type SheetFooterArgs struct {
	Class      string
	Attributes templ.Attributes
}

// SheetTitleArgs defines the args for the SheetTitle component
type SheetTitleArgs struct {
	Class      string
	Attributes templ.Attributes
}

// SheetDescriptionArgs defines the args for the SheetDescription component
type SheetDescriptionArgs struct {
	Class      string
	Attributes templ.Attributes
}

// SheetCloseArgs defines the args for the SheetClose component
type SheetCloseArgs struct {
	SheetID     string
	ReturnValue string // Optional return value when closing the sheet
	AsChild     bool   // Whether to render as a child element instead of button
	Class       string
	Attributes  templ.Attributes
}