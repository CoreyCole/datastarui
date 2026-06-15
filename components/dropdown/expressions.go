package dropdown

import (
	"fmt"

	"github.com/coreycole/datastarui/utils"
)

// DropdownHandler creates handlers for dropdown components
type DropdownHandler struct {
	signals *utils.SignalManager
}

// NewDropdownHandler creates a dropdown handler
func newDropdownHandler(signals *utils.SignalManager) *DropdownHandler {
	return &DropdownHandler{
		signals: signals,
	}
}

// BuildClickOutsideHandler creates a click outside handler for closing dropdown
func (d *DropdownHandler) buildClickOutsideHandler() string {
	return d.signals.ConditionalAction(d.signals.Signal("open"), "open", "false")
}

// BuildEscapeHandler creates an escape key handler for closing dropdown
func (d *DropdownHandler) buildEscapeHandler() string {
	condition := fmt.Sprintf("evt.key === 'Escape' && %s", d.signals.Signal("open"))
	return d.signals.ConditionalAction(condition, "open", "false")
}

func dropdownMenuContentID(id string) string {
	return id + "-content"
}

func dropdownMenuTriggerClickExpr(id string) string {
	signals := utils.Signals(id, DropdownSignals{})
	contentID := dropdownMenuContentID(id)
	positionExpr := fmt.Sprintf(`(function(){const c=document.getElementById(%q);if(c){const r=el.getBoundingClientRect();c.style.setProperty('--dui-dropdown-trigger-top',r.top+'px');c.style.setProperty('--dui-dropdown-trigger-right',r.right+'px');c.style.setProperty('--dui-dropdown-trigger-bottom',r.bottom+'px');c.style.setProperty('--dui-dropdown-trigger-left',r.left+'px');c.style.setProperty('--dui-dropdown-trigger-center-x',(r.left+r.width/2)+'px');c.style.setProperty('--dui-dropdown-trigger-center-y',(r.top+r.height/2)+'px');}})()`, contentID)
	return positionExpr + "; " + signals.Toggle("open")
}

func dropdownMenuItemClickExpr(id, onClick string) string {
	signals := utils.Signals(id, DropdownSignals{})
	hideExpr := signals.Set("open", "false")
	if onClick == "" {
		return hideExpr
	}
	return hideExpr + "; " + onClick
}

func dropdownMenuContentPositionStyle(side, align string, offset int) string {
	if offset == 0 {
		offset = 4 // Default offset like shadcn/ui spacing units.
	}
	offsetCSS := fmt.Sprintf("%.2frem", float64(offset)*0.25)

	var top, left string
	transforms := make([]string, 0, 2)

	switch side {
	case "top":
		top = fmt.Sprintf("calc(var(--dui-dropdown-trigger-top, 0px) - %s)", offsetCSS)
		transforms = append(transforms, "translateY(-100%)")
	case "left":
		left = fmt.Sprintf("calc(var(--dui-dropdown-trigger-left, 0px) - %s)", offsetCSS)
		transforms = append(transforms, "translateX(-100%)")
	case "right":
		left = fmt.Sprintf("calc(var(--dui-dropdown-trigger-right, 0px) + %s)", offsetCSS)
	default:
		top = fmt.Sprintf("calc(var(--dui-dropdown-trigger-bottom, 0px) + %s)", offsetCSS)
	}

	if side == "left" || side == "right" {
		switch align {
		case "end":
			top = "var(--dui-dropdown-trigger-bottom, 0px)"
			transforms = append(transforms, "translateY(-100%)")
		case "center":
			top = "var(--dui-dropdown-trigger-center-y, 0px)"
			transforms = append(transforms, "translateY(-50%)")
		default:
			top = "var(--dui-dropdown-trigger-top, 0px)"
		}
	} else {
		switch align {
		case "end":
			left = "var(--dui-dropdown-trigger-right, 0px)"
			transforms = append(transforms, "translateX(-100%)")
		case "center":
			left = "var(--dui-dropdown-trigger-center-x, 0px)"
			transforms = append(transforms, "translateX(-50%)")
		default:
			left = "var(--dui-dropdown-trigger-left, 0px)"
		}
	}

	style := fmt.Sprintf("display: none; top: %s; left: %s;", top, left)
	if len(transforms) > 0 {
		style += " transform: "
		for i, transform := range transforms {
			if i > 0 {
				style += " "
			}
			style += transform
		}
		style += ";"
	}
	return style
}
