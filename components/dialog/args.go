package dialog

import "github.com/a-h/templ"

// DialogArgs defines the args for the Dialog container (using Datastar signals)
type DialogArgs struct {
	ID          string
	DefaultOpen bool // Whether the dialog should be open by default
	Class       string
	Attributes  templ.Attributes
}

// DialogTriggerArgs defines the args for the DialogTrigger component
type DialogTriggerArgs struct {
	DialogID   string
	AsChild    bool
	Class      string
	Attributes templ.Attributes
}

// DialogContentArgs defines the args for the DialogContent component (for backwards compatibility)
type DialogContentArgs struct {
	Class      string
	Attributes templ.Attributes
}

// DialogOverlayArgs defines the args for the DialogOverlay component
type DialogOverlayArgs struct {
	ID         string
	Class      string
	Attributes templ.Attributes
}

// DialogHeaderArgs defines the args for the DialogHeader component
type DialogHeaderArgs struct {
	Class      string
	Attributes templ.Attributes
}

// DialogFooterArgs defines the args for the DialogFooter component
type DialogFooterArgs struct {
	Class      string
	Attributes templ.Attributes
}

// DialogTitleArgs defines the args for the DialogTitle component
type DialogTitleArgs struct {
	Class      string
	Attributes templ.Attributes
}

// DialogDescriptionArgs defines the args for the DialogDescription component
type DialogDescriptionArgs struct {
	Class      string
	Attributes templ.Attributes
}

// DialogCloseArgs defines the args for the DialogClose component
type DialogCloseArgs struct {
	DialogID    string
	ReturnValue string // Optional return value when closing the dialog
	Variant     string // Button variant: "default", "destructive", "outline", "secondary", "ghost", "link"
	Class       string
	Attributes  templ.Attributes
}
