package sidebar

import (
	"github.com/coreycole/datastarui/utils"
)

// SidebarHandler creates handlers for Sidebar component functionality
type SidebarHandler struct {
	signals *utils.SignalManager
}

// NewSidebarHandler creates a sidebar handler
func NewSidebarHandler(signals *utils.SignalManager) *SidebarHandler {
	return &SidebarHandler{
		signals: signals,
	}
}

// BuildMobileToggleHandler creates handler to open mobile sidebar
func (s *SidebarHandler) BuildMobileToggleHandler() string {
	return s.signals.Set("mobileOpen", "true")
}

// BuildMobileCloseHandler creates handler to close mobile sidebar
func (s *SidebarHandler) BuildMobileCloseHandler() string {
	return s.signals.Set("mobileOpen", "false")
}

// BuildKeyboardToggleHandler creates Cmd+B keyboard shortcut handler
func (s *SidebarHandler) BuildKeyboardToggleHandler() string {
	condition := "(evt.metaKey || evt.ctrlKey) && evt.key === 'b'"
	toggleValue := "!" + s.signals.Signal("mobileOpen")
	return s.signals.ConditionalAction(condition, "mobileOpen", toggleValue)
}

// BuildMobileToggleWithPreventDefaultHandler creates handler that prevents default and toggles
func (s *SidebarHandler) BuildMobileToggleWithPreventDefaultHandler() string {
	expr := utils.NewExpression().
		Statement("evt.preventDefault()").
		Statement(s.signals.Set("mobileOpen", "true"))
	return expr.Build()
}

// BuildConditionalCloseHandler creates handler that only closes if mobile is open
func (s *SidebarHandler) BuildConditionalCloseHandler() string {
	condition := s.signals.Signal("mobileOpen")
	return s.signals.ConditionalAction(condition, "mobileOpen", "false")
}