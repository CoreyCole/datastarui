package tabs

import (
	"fmt"
	"github.com/coreycole/datastarui/utils"
)

// TabsHandler creates handlers for Tabs component functionality
type TabsHandler struct {
	tabsID  string
	signals *utils.SignalManager
}

// NewTabsHandler creates a tabs handler
func NewTabsHandler(tabsID string, signals *utils.SignalManager) *TabsHandler {
	return &TabsHandler{
		tabsID:  tabsID,
		signals: signals,
	}
}

// BuildTriggerClickHandler creates the tab trigger click handler
func (t *TabsHandler) BuildTriggerClickHandler(value string) string {
	return t.signals.SetString("active", value)
}

// BuildTriggerDataClass creates conditional classes for tab triggers
func (t *TabsHandler) BuildTriggerDataClass(value string) string {
	condition := fmt.Sprintf("%s === '%s'", t.signals.Signal("active"), value)

	return utils.NewDataClass().
		Add("bg-background", condition).
		Add("text-foreground", condition).
		Add("shadow-sm", condition).
		Build()
}

// BuildTriggerStateAttr creates the data-state attribute expression
func (t *TabsHandler) BuildTriggerStateAttr(value string) string {
	condition := fmt.Sprintf("%s === '%s'", t.signals.Signal("active"), value)
	return fmt.Sprintf("%s ? 'active' : 'inactive'", condition)
}

// BuildTriggerAriaSelected creates the aria-selected attribute expression
func (t *TabsHandler) BuildTriggerAriaSelected(value string) string {
	condition := fmt.Sprintf("%s === '%s'", t.signals.Signal("active"), value)
	return fmt.Sprintf("%s ? 'true' : 'false'", condition)
}

// BuildTriggerTabIndex creates the tabindex attribute expression
func (t *TabsHandler) BuildTriggerTabIndex(value string) string {
	condition := fmt.Sprintf("%s === '%s'", t.signals.Signal("active"), value)
	return fmt.Sprintf("%s ? '0' : '-1'", condition)
}

// BuildContentShowExpression creates the show expression for tab content
func (t *TabsHandler) BuildContentShowExpression(value string) string {
	return fmt.Sprintf("%s === '%s'", t.signals.Signal("active"), value)
}

// BuildContentAriaHidden creates the aria-hidden attribute expression
func (t *TabsHandler) BuildContentAriaHidden(value string) string {
	condition := fmt.Sprintf("%s === '%s'", t.signals.Signal("active"), value)
	return fmt.Sprintf("%s ? 'false' : 'true'", condition)
}
