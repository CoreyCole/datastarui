package label

import "github.com/a-h/templ"

type LabelArgs struct {
	Class      string           // Additional CSS classes
	For        string           // HTML for attribute (links to input ID)
	Attributes templ.Attributes // Additional HTML attributes
}
