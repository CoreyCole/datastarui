package dialog

import "github.com/a-h/templ"

// DialogProps defines the props for the Dialog container (native <dialog> element)
type DialogProps struct {
	ID         string
	Class      string
	ClosedBy   string // "none", "closerequest", "any" - controls how dialog can be dismissed
	Attributes templ.Attributes
}

// DialogTriggerProps defines the props for the DialogTrigger component
type DialogTriggerProps struct {
	DialogID   string
	Modal      bool // true for showModal(), false for show()
	AsChild    bool
	Class      string
	Attributes templ.Attributes
}

// DialogContentProps defines the props for the DialogContent component (for backwards compatibility)
type DialogContentProps struct {
	ID          string
	Class       string
	Size        string // "sm", "md", "lg", "xl", "full"
	ShowOverlay bool   // Whether to show the backdrop overlay
	Attributes  templ.Attributes
}

// DialogOverlayProps defines the props for the DialogOverlay component
type DialogOverlayProps struct {
	ID         string
	Class      string
	Attributes templ.Attributes
}

// DialogHeaderProps defines the props for the DialogHeader component
type DialogHeaderProps struct {
	Class      string
	Attributes templ.Attributes
}

// DialogFooterProps defines the props for the DialogFooter component
type DialogFooterProps struct {
	Class      string
	Attributes templ.Attributes
}

// DialogTitleProps defines the props for the DialogTitle component
type DialogTitleProps struct {
	Class      string
	Attributes templ.Attributes
}

// DialogDescriptionProps defines the props for the DialogDescription component
type DialogDescriptionProps struct {
	Class      string
	Attributes templ.Attributes
}

// DialogCloseProps defines the props for the DialogClose component
type DialogCloseProps struct {
	DialogID    string
	ReturnValue string // Optional return value when closing the dialog
	AsChild     bool
	Class       string
	Attributes  templ.Attributes
}
