package utils

import (
	"fmt"
	"strings"
)

// DatastarExpression represents a structured Datastar expression builder
type DatastarExpression struct {
	statements []string
}

// NewExpression creates a new Datastar expression builder
func NewExpression() *DatastarExpression {
	return &DatastarExpression{
		statements: make([]string, 0),
	}
}

// Statement adds a statement to the expression
func (e *DatastarExpression) Statement(stmt string) *DatastarExpression {
	e.statements = append(e.statements, stmt)
	return e
}

// SetSignal adds a signal assignment statement
func (e *DatastarExpression) SetSignal(signal, value string) *DatastarExpression {
	e.statements = append(e.statements, fmt.Sprintf("$%s = %s", signal, value))
	return e
}

// Conditional adds a conditional statement
func (e *DatastarExpression) Conditional(condition, trueExpr, falseExpr string) *DatastarExpression {
	e.statements = append(e.statements, fmt.Sprintf("%s ? %s : %s", condition, trueExpr, falseExpr))
	return e
}

// PreventDefault adds preventDefault and stopPropagation calls
func (e *DatastarExpression) PreventDefault() *DatastarExpression {
	e.statements = append(e.statements, "evt.preventDefault()")
	e.statements = append(e.statements, "evt.stopPropagation()")
	return e
}

// Build returns the final Datastar expression string
func (e *DatastarExpression) Build() string {
	if len(e.statements) == 0 {
		return ""
	}
	return strings.Join(e.statements, "; ")
}

// KeyboardHandler creates a keyboard event handler expression
type KeyboardHandler struct {
	keys   []string
	action *DatastarExpression
}

// NewKeyboardHandler creates a new keyboard handler
func NewKeyboardHandler(keys ...string) *KeyboardHandler {
	return &KeyboardHandler{
		keys:   keys,
		action: NewExpression(),
	}
}

// OnKeys adds the action to perform when specified keys are pressed
func (k *KeyboardHandler) OnKeys(actionFn func(*DatastarExpression) *DatastarExpression) *KeyboardHandler {
	k.action = actionFn(k.action)
	return k
}

// Build creates the complete keyboard handler expression
func (k *KeyboardHandler) Build() string {
	if len(k.keys) == 0 || len(k.action.statements) == 0 {
		return "null"
	}

	// Build key condition
	keyConditions := make([]string, len(k.keys))
	for i, key := range k.keys {
		keyConditions[i] = fmt.Sprintf("evt.key === '%s'", key)
	}
	condition := "(" + strings.Join(keyConditions, " || ") + ")"

	// Build action as a compound expression using comma operator
	// The comma operator allows multiple statements in a single expression
	actionStatements := make([]string, len(k.action.statements))
	copy(actionStatements, k.action.statements)
	
	// Always wrap multiple statements in parentheses for the comma operator
	var action string
	if len(actionStatements) == 1 {
		action = actionStatements[0]
	} else {
		action = "(" + strings.Join(actionStatements, ", ") + ")"
	}

	return fmt.Sprintf("%s ? %s : null", condition, action)
}

// SelectItemHandler creates a standardized select item handler
type SelectItemHandler struct {
	selectID string
	value    string
}

// NewSelectItemHandler creates a select item click/keyboard handler
func NewSelectItemHandler(selectID, value string) *SelectItemHandler {
	return &SelectItemHandler{
		selectID: selectID,
		value:    value,
	}
}

// BuildClickHandler creates the click handler expression
func (s *SelectItemHandler) BuildClickHandler() string {
	return NewExpression().
		SetSignal(s.selectID+".value", fmt.Sprintf("'%s'", s.value)).
		SetSignal(s.selectID+".label", "evt.currentTarget.querySelector('.select-item-text').textContent").
		SetSignal(s.selectID+".open", "false").
		SetSignal(s.selectID+".highlighted", "-1").
		Build()
}

// BuildKeyboardHandler creates the keyboard handler expression
func (s *SelectItemHandler) BuildKeyboardHandler() string {
	return NewKeyboardHandler(" ", "Enter").
		OnKeys(func(expr *DatastarExpression) *DatastarExpression {
			return expr.
				PreventDefault().
				SetSignal(s.selectID+".value", fmt.Sprintf("'%s'", s.value)).
				SetSignal(s.selectID+".label", "evt.currentTarget.querySelector('.select-item-text').textContent").
				SetSignal(s.selectID+".open", "false").
				SetSignal(s.selectID+".highlighted", "-1")
		}).
		Build()
}

// SelectTriggerHandler creates handlers for select trigger buttons
type SelectTriggerHandler struct {
	selectID string
	signals  *SignalManager
}

// NewSelectTriggerHandler creates a select trigger handler
func NewSelectTriggerHandler(selectID string, signals *SignalManager) *SelectTriggerHandler {
	return &SelectTriggerHandler{
		selectID: selectID,
		signals:  signals,
	}
}

// BuildClickHandler creates the trigger click handler
func (s *SelectTriggerHandler) BuildClickHandler() string {
	// Find current selection index expression
	findIndexExpr := fmt.Sprintf(
		"Math.max(0, Array.from(document.querySelectorAll('[data-select-id=\"%s\"] [data-select-item]:not([data-disabled])')).findIndex(el => el.dataset.value === %s))",
		s.selectID,
		s.signals.Signal("value"),
	)

	return NewExpression().
		Statement(s.signals.Toggle("open")).
		SetSignal(s.selectID+".highlighted", findIndexExpr).
		Build()
}

// BuildKeyboardHandler creates the trigger keyboard handler
func (s *SelectTriggerHandler) BuildKeyboardHandler() string {
	findIndexExpr := fmt.Sprintf(
		"Math.max(0, Array.from(document.querySelectorAll('[data-select-id=\"%s\"] [data-select-item]:not([data-disabled])')).findIndex(el => el.dataset.value === %s))",
		s.selectID,
		s.signals.Signal("value"),
	)

	// Build the compound action with proper parentheses for multiple statements
	action := fmt.Sprintf("(evt.preventDefault(), evt.stopPropagation(), %s, %s)",
		s.signals.Set("open", "true"),
		s.signals.Set("highlighted", findIndexExpr),
	)

	return NewExpression().
		Conditional(
			fmt.Sprintf("(evt.key === 'ArrowDown' || evt.key === 'ArrowUp' || evt.key === ' ' || evt.key === 'Enter') && !%s", s.signals.Signal("open")),
			action,
			"null",
		).
		Build()
}

// SelectContentHandler creates handlers for select dropdown content
type SelectContentHandler struct {
	selectID string
	signals  *SignalManager
}

// NewSelectContentHandler creates a select content handler
func NewSelectContentHandler(selectID string, signals *SignalManager) *SelectContentHandler {
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

	return NewExpression().
		Statement(arrowDown).
		Statement(arrowUp).
		Statement(enterSpace).
		Statement(escape).
		Statement(tab).
		Build()
}

// ClassExpression helps build data-class expressions
type ClassExpression struct {
	classes map[string]string // className -> condition
}

// NewClassExpression creates a new class expression builder
func NewClassExpression() *ClassExpression {
	return &ClassExpression{
		classes: make(map[string]string),
	}
}

// AddClass adds a conditional class
func (c *ClassExpression) AddClass(className, condition string) *ClassExpression {
	c.classes[className] = condition
	return c
}

// AddClasses adds multiple conditional classes
func (c *ClassExpression) AddClasses(classes map[string]string) *ClassExpression {
	for className, condition := range classes {
		c.classes[className] = condition
	}
	return c
}

// Build creates the data-class object expression
func (c *ClassExpression) Build() string {
	if len(c.classes) == 0 {
		return "{}"
	}

	var parts []string
	for className, condition := range c.classes {
		parts = append(parts, fmt.Sprintf("'%s': %s", className, condition))
	}

	return "{" + strings.Join(parts, ", ") + "}"
}

// ConditionalClasses is a helper for common conditional class patterns
type ConditionalClasses struct {
	signalPath string
}

// NewConditionalClasses creates a conditional class helper for a signal
func NewConditionalClasses(signalPath string) *ConditionalClasses {
	return &ConditionalClasses{signalPath: signalPath}
}

// Equals adds a class when signal equals a value
func (cc *ConditionalClasses) Equals(value, className string) *ClassExpression {
	condition := fmt.Sprintf("$%s === '%s'", cc.signalPath, value)
	return NewClassExpression().AddClass(className, condition)
}

// Boolean adds a class when signal is truthy
func (cc *ConditionalClasses) Boolean(className string) *ClassExpression {
	condition := fmt.Sprintf("$%s", cc.signalPath)
	return NewClassExpression().AddClass(className, condition)
}

// NotBoolean adds a class when signal is falsy
func (cc *ConditionalClasses) NotBoolean(className string) *ClassExpression {
	condition := fmt.Sprintf("!$%s", cc.signalPath)
	return NewClassExpression().AddClass(className, condition)
}

