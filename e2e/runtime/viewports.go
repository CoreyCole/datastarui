package runtime

import "fmt"

type ViewportClass string

const (
	ViewportMobile      ViewportClass = "mobile"
	ViewportDesktopHalf ViewportClass = "desktop-half"
	ViewportDesktopFull ViewportClass = "desktop-full"
)

type Viewport struct {
	Class         ViewportClass
	Width, Height int
}

func DefaultViewports() map[ViewportClass]Viewport {
	return map[ViewportClass]Viewport{
		ViewportMobile:      {Class: ViewportMobile, Width: 390, Height: 844},
		ViewportDesktopHalf: {Class: ViewportDesktopHalf, Width: 900, Height: 900},
		ViewportDesktopFull: {Class: ViewportDesktopFull, Width: 1440, Height: 1000},
	}
}

func DefaultVerifyViewports() []ViewportClass {
	return []ViewportClass{ViewportMobile, ViewportDesktopHalf, ViewportDesktopFull}
}

func ResolveViewports(names []string) ([]Viewport, error) {
	defs := DefaultViewports()
	if len(names) == 0 {
		names = []string{string(ViewportDesktopFull)}
	}
	out := make([]Viewport, 0, len(names))
	for _, name := range names {
		vp, ok := defs[ViewportClass(name)]
		if !ok {
			return nil, fmt.Errorf("unknown viewport %q", name)
		}
		out = append(out, vp)
	}
	return out, nil
}
