package selectcomponent

import (
	"fmt"

	"github.com/coreycole/datastarui/utils"
)

// SelectTriggerHandler creates handlers for select trigger functionality
type SelectTriggerHandler struct {
	selectID string
	signals  *utils.SignalManager
}

// NewSelectTriggerHandler creates a select trigger handler
func NewSelectTriggerHandler(selectID string, signals *utils.SignalManager) *SelectTriggerHandler {
	return &SelectTriggerHandler{
		selectID: selectID,
		signals:  signals,
	}
}

// BuildClickHandler creates the trigger click handler
func (s *SelectTriggerHandler) BuildClickHandler() string {
	return s.signals.Toggle("open") + "; " + s.signals.Set("highlighted", "-1")
}

// BuildKeyboardHandler creates the trigger keyboard navigation handler
func (s *SelectTriggerHandler) BuildKeyboardHandler() string {
	expr := utils.NewExpression()

	// ArrowDown/Up: open dropdown
	expr.Statement(fmt.Sprintf("(evt.key === 'ArrowDown' || evt.key === 'ArrowUp') && !%s ? (%s, evt.preventDefault()) : null",
		s.signals.Signal("open"),
		s.signals.Set("open", "true")))

	// Space/Enter: toggle dropdown
	expr.Statement(fmt.Sprintf("(evt.key === ' ' || evt.key === 'Enter') ? (%s, %s, evt.preventDefault()) : null",
		s.signals.Toggle("open"),
		s.signals.Set("highlighted", "-1")))

	return expr.Build()
}

// SelectContentHandler creates handlers for select content functionality
type SelectContentHandler struct {
	selectID string
	signals  *utils.SignalManager
}

// NewSelectContentHandler creates a select content handler
func NewSelectContentHandler(selectID string, signals *utils.SignalManager) *SelectContentHandler {
	return &SelectContentHandler{
		selectID: selectID,
		signals:  signals,
	}
}

// BuildKeyboardHandler creates the content keyboard navigation handler
func (s *SelectContentHandler) BuildKeyboardHandler() string {
	maxItemsExpr := fmt.Sprintf(
		"document.querySelector('[data-select-id=\"%s\"]').querySelectorAll('[data-select-item]:not([data-disabled])').length - 1",
		s.selectID,
	)

	selectOpenCheck := fmt.Sprintf(
		"document.querySelector('[data-select-id=\"%s\"]') && %s",
		s.selectID,
		s.signals.Signal("open"),
	)

	// Build individual handlers using comma operator for multiple statements
	arrowDown := fmt.Sprintf(
		"evt.key === 'ArrowDown' && %s ? (evt.preventDefault(), evt.stopPropagation(), %s) : null",
		selectOpenCheck,
		s.signals.Set("highlighted", fmt.Sprintf("Math.min(%s, %s + 1)", maxItemsExpr, s.signals.Signal("highlighted"))),
	)

	arrowUp := fmt.Sprintf(
		"evt.key === 'ArrowUp' && %s ? (evt.preventDefault(), evt.stopPropagation(), %s) : null",
		selectOpenCheck,
		s.signals.Set("highlighted", fmt.Sprintf("Math.max(0, %s - 1)", s.signals.Signal("highlighted"))),
	)

	enterSpace := fmt.Sprintf(
		"(evt.key === 'Enter' || evt.key === ' ') && %s && %s >= 0 ? (evt.preventDefault(), evt.stopPropagation(), document.querySelector('[data-select-id=\"%s\"]').querySelector('[data-select-item][data-index=\"' + %s + '\"]')?.click()) : null",
		selectOpenCheck,
		s.signals.Signal("highlighted"),
		s.selectID,
		s.signals.Signal("highlighted"),
	)

	escape := fmt.Sprintf(
		"evt.key === 'Escape' && %s ? (evt.preventDefault(), evt.stopPropagation(), %s) : null",
		selectOpenCheck,
		s.signals.Set("open", "false"),
	)

	tab := fmt.Sprintf(
		"evt.key === 'Tab' && %s ? %s : null",
		selectOpenCheck,
		s.signals.Set("open", "false"),
	)

	return utils.NewExpression().
		Statement(arrowDown).
		Statement(arrowUp).
		Statement(enterSpace).
		Statement(escape).
		Statement(tab).
		Build()
}

// SelectItemHandler creates handlers for select item functionality
type SelectItemHandler struct {
	selectID string
	value    string
}

// NewSelectItemHandler creates a select item handler
func NewSelectItemHandler(selectID, value string) *SelectItemHandler {
	return &SelectItemHandler{
		selectID: selectID,
		value:    value,
	}
}

// BuildClickHandler creates the item click handler
func (s *SelectItemHandler) BuildClickHandler() string {
	// Extract label from clicked item
	labelExpr := "evt.currentTarget.querySelector('.select-item-text')?.textContent.trim() || ''"

	return fmt.Sprintf(`$%s.value = '%s'; $%s.label = %s; $%s.open = false`,
		s.selectID, s.value,
		s.selectID, labelExpr,
		s.selectID,
	)
}

// BuildKeyboardHandler creates the item keyboard handler
func (s *SelectItemHandler) BuildKeyboardHandler() string {
	return fmt.Sprintf(`(evt.key === 'Enter' || evt.key === ' ') ? (evt.preventDefault(), evt.target.click()) : null`)
}

