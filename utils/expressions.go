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

// CalendarDayHandler creates handlers for calendar day interactions
type CalendarDayHandler struct {
	calendarID  string
	signals     *SignalManager
	monthOffset int
	mode        string
	dayNumber   int
}

// NewCalendarDayHandler creates a calendar day handler
func NewCalendarDayHandler(calendarID string, signals *SignalManager, monthOffset int, mode string, dayNumber int) *CalendarDayHandler {
	return &CalendarDayHandler{
		calendarID:  calendarID,
		signals:     signals,
		monthOffset: monthOffset,
		mode:        mode,
		dayNumber:   dayNumber,
	}
}

// BuildClickHandler creates the day click handler expression
func (c *CalendarDayHandler) BuildClickHandler() string {
	expr := NewExpression()

	if c.mode == "range" {
		// Range mode: complex logic for start/end selection
		dateCalcExpr := fmt.Sprintf(
			"new Date(%s + 'T12:00:00Z').getFullYear() + '-' + (new Date(%s + 'T12:00:00Z').getMonth() + 1 + %d).toString().padStart(2, '0') + '-' + parseInt(evt.target.textContent).toString().padStart(2, '0')",
			c.signals.Signal("currentDate"),
			c.signals.Signal("currentDate"),
			c.monthOffset,
		)

		// Use functional approach with intelligent date swapping
		// If no rangeStart: set rangeStart
		// If rangeStart but no rangeEnd: 
		//   - If clicked < rangeStart: swap (clicked becomes start, old start becomes end)
		//   - If clicked > rangeStart: set as end
		// If both exist: reset and start over
		expr.Statement(fmt.Sprintf(
			"((clickedDate) => !%s ? (%s, %s) : !%s ? (clickedDate < %s ? (%s, %s) : (%s, %s)) : (%s, %s))(%s)",
			c.signals.Signal("rangeStart"),
			c.signals.Set("rangeStart", "clickedDate"),
			c.signals.Set("rangeEnd", "''"),
			c.signals.Signal("rangeEnd"),
			c.signals.Signal("rangeStart"),
			c.signals.Set("rangeEnd", c.signals.Signal("rangeStart")),
			c.signals.Set("rangeStart", "clickedDate"),
			c.signals.Set("rangeStart", c.signals.Signal("rangeStart")),
			c.signals.Set("rangeEnd", "clickedDate"),
			c.signals.Set("rangeStart", "clickedDate"),
			c.signals.Set("rangeEnd", "''"),
			dateCalcExpr,
		))
	} else {
		// Single mode: calculate date and set selectedDate
		expr.Statement("const currentDate = new Date(" + c.signals.Signal("currentDate") + " + 'T12:00:00Z')")
		expr.Statement(fmt.Sprintf("const targetDate = new Date(currentDate.getFullYear(), currentDate.getMonth() + %d, 1)", c.monthOffset))
		expr.Statement("const year = targetDate.getFullYear()")
		expr.Statement("const month = targetDate.getMonth()")
		expr.Statement("const buttonDay = parseInt(evt.target.textContent)")
		expr.Statement("const clickedDate = year + '-' + (month + 1).toString().padStart(2, '0') + '-' + buttonDay.toString().padStart(2, '0')")
		expr.Statement(c.signals.Set("selectedDate", "clickedDate"))
	}

	return expr.Build()
}

// BuildDateInputSync creates the DateInput synchronization expression
func (c *CalendarDayHandler) BuildDateInputSync(datePickerInputsID string) string {
	expr := NewExpression()

	if c.mode == "range" {
		// Sync range signals for DateInput - use the same namespace since datePickerInputsID === calendarID
		expr.Statement(c.signals.Set("startDateValue", c.signals.Signal("rangeStart")))
		expr.Statement(c.signals.Set("endDateValue", c.signals.Signal("rangeEnd")))
		expr.Statement(c.signals.Set("startInputValue", c.signals.Signal("rangeStart") + " ? new Date(" + c.signals.Signal("rangeStart") + " + 'T12:00:00Z').toLocaleDateString('en-US', {month: '2-digit', day: '2-digit', year: 'numeric', timeZone: 'UTC'}) : ''"))
		expr.Statement(c.signals.Set("endInputValue", c.signals.Signal("rangeEnd") + " ? new Date(" + c.signals.Signal("rangeEnd") + " + 'T12:00:00Z').toLocaleDateString('en-US', {month: '2-digit', day: '2-digit', year: 'numeric', timeZone: 'UTC'}) : ''"))
	} else {
		// Sync single date signals for DateInput - use the same namespace since datePickerInputsID === calendarID  
		expr.Statement(c.signals.Set("dateValue", c.signals.Signal("selectedDate")))
		expr.Statement(c.signals.Set("inputValue", c.signals.Signal("selectedDate") + " ? new Date(" + c.signals.Signal("selectedDate") + " + 'T12:00:00Z').toLocaleDateString('en-US', {month: '2-digit', day: '2-digit', year: 'numeric', timeZone: 'UTC'}) : ''"))
	}

	return expr.Build()
}

// CalendarSelectionClasses creates conditional classes for calendar day selection
type CalendarSelectionClasses struct {
	calendarID  string
	signals     *SignalManager
	monthOffset int
	dayNumber   int
	mode        string
}

// NewCalendarSelectionClasses creates a calendar selection class helper
func NewCalendarSelectionClasses(calendarID string, signals *SignalManager, monthOffset int, dayNumber int, mode string) *CalendarSelectionClasses {
	return &CalendarSelectionClasses{
		calendarID:  calendarID,
		signals:     signals,
		monthOffset: monthOffset,
		dayNumber:   dayNumber,
		mode:        mode,
	}
}

// Build creates the data-class expression for calendar day selection
func (c *CalendarSelectionClasses) Build() string {
	classExpr := NewClassExpression()

	// Calculate expected date for this day
	expectedDateExpr := fmt.Sprintf(
		"new Date(%s + 'T12:00:00Z').getFullYear() + '-' + (new Date(%s + 'T12:00:00Z').getMonth() + 1 + %d).toString().padStart(2, '0') + '-' + (%d).toString().padStart(2, '0')",
		c.signals.Signal("currentDate"),
		c.signals.Signal("currentDate"),
		c.monthOffset,
		c.dayNumber,
	)

	if c.mode == "range" {
		// Range mode: highlight start/end dates and in-between days
		classExpr.AddClass("bg-primary text-primary-foreground",
			fmt.Sprintf("%s === (%s) || %s === (%s)",
				c.signals.Signal("rangeStart"), expectedDateExpr,
				c.signals.Signal("rangeEnd"), expectedDateExpr,
			))

		classExpr.AddClass("bg-accent text-accent-foreground",
			fmt.Sprintf("%s && %s && (%s) > %s && (%s) < %s",
				c.signals.Signal("rangeStart"),
				c.signals.Signal("rangeEnd"),
				expectedDateExpr,
				c.signals.Signal("rangeStart"),
				expectedDateExpr,
				c.signals.Signal("rangeEnd"),
			))
	} else {
		// Single mode: highlight selected date
		classExpr.AddClass("bg-primary text-primary-foreground",
			fmt.Sprintf("%s === (%s)", c.signals.Signal("selectedDate"), expectedDateExpr))
	}

	return classExpr.Build()
}

// DatePickerPopoverHandler creates handlers for DatePicker popover functionality  
type DatePickerPopoverHandler struct {
	datePickerID string
	dateInputID  string
	mode         string
	signals      *SignalManager
}

// NewDatePickerPopoverHandler creates a datepicker popover handler
func NewDatePickerPopoverHandler(datePickerID, dateInputID, mode string, signals *SignalManager) *DatePickerPopoverHandler {
	return &DatePickerPopoverHandler{
		datePickerID: datePickerID,
		dateInputID:  dateInputID,
		mode:         mode,
		signals:      signals,
	}
}

// BuildEscapeHandler creates clean escape key handler for popover
func (d *DatePickerPopoverHandler) BuildEscapeHandler() string {
	openCheck := fmt.Sprintf("document.querySelector('[data-datepicker-id=\"%s\"]') && %s", d.datePickerID, d.signals.Signal("open"))
	
	inputID := d.dateInputID
	if d.mode == "range" {
		inputID += "_start" // Focus on start input in range mode
	}
	
	return NewExpression().
		Conditional(
			fmt.Sprintf("evt.key === 'Escape' && %s", openCheck),
			fmt.Sprintf("(evt.preventDefault(), evt.stopPropagation(), %s, document.getElementById('%s').focus())", 
				d.signals.Set("open", "false"), inputID),
			"null",
		).
		Build()
}

// BuildTabHandler creates clean tab key handler for popover
func (d *DatePickerPopoverHandler) BuildTabHandler() string {
	openCheck := fmt.Sprintf("document.querySelector('[data-datepicker-id=\"%s\"]') && %s", d.datePickerID, d.signals.Signal("open"))
	
	return NewExpression().
		Conditional(
			fmt.Sprintf("evt.key === 'Tab' && %s", openCheck),
			d.signals.Set("open", "false"),
			"null",
		).
		Build()
}

// BuildKeyboardHandler combines escape and tab handlers
func (d *DatePickerPopoverHandler) BuildKeyboardHandler() string {
	escapeHandler := d.BuildEscapeHandler()
	tabHandler := d.BuildTabHandler()
	return escapeHandler + "; " + tabHandler
}

// BuildDateSelectHandler creates the complex date selection handler with DateInput sync
func (d *DatePickerPopoverHandler) BuildDateSelectHandler(closeOnSelect bool, dateInputSignals *SignalManager) string {
	expr := NewExpression()
	
	if d.mode == "range" {
		// Range mode: handle both start and end date selection
		rangeCompleteCondition := "evt.detail.rangeStart && evt.detail.rangeEnd"
		rangeCompleteActions := fmt.Sprintf(
			"((startDisplay, endDisplay) => (%s, %s, %s, %s, %s, %s))(evt.detail.rangeStart.replace(/-/g, '/'), evt.detail.rangeEnd.replace(/-/g, '/'))",
			dateInputSignals.Set("startInputValue", "startDisplay"),
			dateInputSignals.Set("startDateValue", "evt.detail.rangeStart"),
			dateInputSignals.Set("endInputValue", "endDisplay"),
			dateInputSignals.Set("endDateValue", "evt.detail.rangeEnd"),
			d.signals.Set("rangeStart", "evt.detail.rangeStart"),
			d.signals.Set("rangeEnd", "evt.detail.rangeEnd"),
		)
		
		rangeStartCondition := "evt.detail.rangeStart"
		rangeStartActions := fmt.Sprintf(
			"((startDisplay) => (%s, %s, %s, %s, %s, %s))(evt.detail.rangeStart.replace(/-/g, '/'))",
			dateInputSignals.Set("startInputValue", "startDisplay"),
			dateInputSignals.Set("startDateValue", "evt.detail.rangeStart"),
			dateInputSignals.Set("endInputValue", "''"),
			dateInputSignals.Set("endDateValue", "''"),
			d.signals.Set("rangeStart", "evt.detail.rangeStart"),
			d.signals.Set("rangeEnd", "''"),
		)
		
		selectAction := fmt.Sprintf("%s ? %s : %s ? %s : null", 
			rangeCompleteCondition, rangeCompleteActions,
			rangeStartCondition, rangeStartActions)
		
		if closeOnSelect {
			expr.Statement(fmt.Sprintf("(%s, %s)", selectAction, d.signals.Set("open", "false")))
		} else {
			expr.Statement(selectAction)
		}
	} else {
		// Single mode: simpler date selection
		singleDateCondition := "evt.detail.selectedDate"
		singleDateActions := fmt.Sprintf(
			"((displayDate) => (%s, %s, %s))(evt.detail.selectedDate.replace(/-/g, '/'))",
			dateInputSignals.Set("inputValue", "displayDate"),
			dateInputSignals.Set("dateValue", "evt.detail.selectedDate"),
			d.signals.Set("selectedDate", "evt.detail.selectedDate"),
		)
		
		selectAction := fmt.Sprintf("%s ? %s : null", singleDateCondition, singleDateActions)
		
		if closeOnSelect {
			expr.Statement(fmt.Sprintf("(%s, %s)", selectAction, d.signals.Set("open", "false")))
		} else {
			expr.Statement(selectAction)
		}
	}
	
	return expr.Build()
}

// BuildMonthChangeHandler creates the month navigation handler
func (d *DatePickerPopoverHandler) BuildMonthChangeHandler() string {
	return NewExpression().
		Conditional(
			"evt.detail.displayMonth",
			d.signals.Set("displayMonth", "evt.detail.displayMonth"),
			"null",
		).
		Build()
}

// BuildOpenTriggerHandler creates the calendar icon click handler
func (d *DatePickerPopoverHandler) BuildOpenTriggerHandler() string {
	return NewExpression().
		Statement("evt.preventDefault()").
		Statement("evt.stopPropagation()").
		Statement(d.signals.Toggle("open")).
		Build()
}

// BuildClickOutsideHandler creates the click outside handler to close popover
func (d *DatePickerPopoverHandler) BuildClickOutsideHandler() string {
	return NewExpression().
		Conditional(
			d.signals.Signal("open"),
			d.signals.Set("open", "false"),
			"null",
		).
		Build()
}

// DateInputHandler creates handlers for DateInput component functionality
type DateInputHandler struct {
	inputID    string
	signals    *SignalManager
	calendarID string // Optional calendar coordination
}

// NewDateInputHandler creates a DateInput handler
func NewDateInputHandler(inputID string, signals *SignalManager, calendarID string) *DateInputHandler {
	return &DateInputHandler{
		inputID:    inputID,
		signals:    signals,
		calendarID: calendarID,
	}
}

// BuildInputHandler creates the complex input formatting handler
func (d *DateInputHandler) BuildInputHandler(inputSignal, dateSignal string) string {
	expr := NewExpression()
	
	// Core date formatting logic
	expr.Statement("const value = evt.target.value")
	expr.Statement("const lastChar = value.slice(-1)")
	expr.Statement("const beforeSlash = value.slice(0, -1)")
	
	// Check cursor position and text selection state
	expr.Statement("const cursorAtEnd = evt.target.selectionStart === evt.target.value.length")
	expr.Statement("const hasSelection = evt.target.selectionStart !== evt.target.selectionEnd")
	expr.Statement("const shouldFormat = cursorAtEnd && evt.target.value.replace(/[^\\d]/g, '').length <= evt.target.value.length && !hasSelection")
	
	// Handle slash typing and auto-format single digits - only when appropriate
	expr.Statement(`if (shouldFormat) {
		if (lastChar === '/') {
			const parts = beforeSlash.split('/');
			if (parts.length === 1 && parts[0].length === 1) {
				evt.target.value = '0' + parts[0] + '/';
			} else if (parts.length === 2 && parts[1].length === 1) {
				evt.target.value = parts[0] + '/0' + parts[1] + '/';
			} else {
				evt.target.value = value;
			}
		} else {
			evt.target.value = (digits => digits.length >= 1 ? digits.substring(0,2) + (digits.length >= 3 ? '/' + digits.substring(2,4) + (digits.length >= 5 ? '/' + digits.substring(4,8) : '') : '') : '')(evt.target.value.replace(/[^\\d]/g, ''));
		}
	}`)
	
	// Always update input signal
	expr.Statement(d.signals.Set(inputSignal, "evt.target.value"))
	
	// Update date signal based on valid date patterns - be more lenient during input
	fourDigitYear := d.signals.Set(dateSignal, "evt.target.value.split('/')[2] + '-' + evt.target.value.split('/')[0].padStart(2, '0') + '-' + evt.target.value.split('/')[1].padStart(2, '0')")
	twoDigitYear := d.signals.Set(dateSignal, "'20' + evt.target.value.split('/')[2] + '-' + evt.target.value.split('/')[0].padStart(2, '0') + '-' + evt.target.value.split('/')[1].padStart(2, '0')")
	
	// Only clear date signal if we have less than 3 parts or empty year, allow partial editing
	expr.Statement(fmt.Sprintf(
		"evt.target.value.split('/').length === 3 && evt.target.value.split('/')[2].length >= 1 ? (evt.target.value.split('/')[2].length === 4 ? %s : evt.target.value.split('/')[2].length === 2 ? %s : null) : %s",
		fourDigitYear, twoDigitYear, d.signals.Set(dateSignal, "''")),
	)
	
	// Add calendar coordination if provided
	if d.calendarID != "" {
		d.addCalendarCoordination(expr, dateSignal)
	}
	
	return expr.Build()
}

// addCalendarCoordination adds calendar synchronization logic
func (d *DateInputHandler) addCalendarCoordination(expr *DatastarExpression, dateSignal string) {
	if strings.Contains(dateSignal, "startDateValue") {
		// Start date coordination
		expr.Statement(fmt.Sprintf(
			"evt.target.value.split('/').length === 3 && (evt.target.value.split('/')[2].length === 4 || evt.target.value.split('/')[2].length === 2) ? (fullYear => ($%s.rangeStart = fullYear + '-' + evt.target.value.split('/')[0].padStart(2, '0') + '-' + evt.target.value.split('/')[1].padStart(2, '0'), $%s.currentDate = fullYear + '-' + evt.target.value.split('/')[0].padStart(2, '0') + '-01'))(evt.target.value.split('/')[2].length === 4 ? evt.target.value.split('/')[2] : '20' + evt.target.value.split('/')[2]) : ($%s.rangeStart = '')",
			d.calendarID, d.calendarID, d.calendarID,
		))
	} else if strings.Contains(dateSignal, "endDateValue") {
		// End date coordination
		expr.Statement(fmt.Sprintf(
			"evt.target.value.split('/').length === 3 && (evt.target.value.split('/')[2].length === 4 || evt.target.value.split('/')[2].length === 2) ? (fullYear => ($%s.rangeEnd = fullYear + '-' + evt.target.value.split('/')[0].padStart(2, '0') + '-' + evt.target.value.split('/')[1].padStart(2, '0'), $%s.currentDate = fullYear + '-' + evt.target.value.split('/')[0].padStart(2, '0') + '-01'))(evt.target.value.split('/')[2].length === 4 ? evt.target.value.split('/')[2] : '20' + evt.target.value.split('/')[2]) : ($%s.rangeEnd = '')",
			d.calendarID, d.calendarID, d.calendarID,
		))
	} else if strings.Contains(dateSignal, "dateValue") {
		// Single date coordination
		expr.Statement(fmt.Sprintf(
			"evt.target.value.split('/').length === 3 && (evt.target.value.split('/')[2].length === 4 || evt.target.value.split('/')[2].length === 2) ? (fullYear => (%s, %s))(evt.target.value.split('/')[2].length === 4 ? evt.target.value.split('/')[2] : '20' + evt.target.value.split('/')[2]) : %s",
			d.signals.Set("selectedDate", "fullYear + '-' + evt.target.value.split('/')[0].padStart(2, '0') + '-' + evt.target.value.split('/')[1].padStart(2, '0')"),
			d.signals.Set("currentDate", "fullYear + '-' + evt.target.value.split('/')[0].padStart(2, '0') + '-01'"),
			d.signals.Set("selectedDate", "''"),
		))
	}
}

// BuildBlurHandler creates the blur completion handler
func (d *DateInputHandler) BuildBlurHandler(inputSignal, dateSignal string) string {
	expr := NewExpression()
	
	// Format date padding logic
	formatValue := "evt.target.value.split('/').map((p,i) => i < 2 ? p.padStart(2, '0') : (p.length === 2 ? '20' + p : p)).join('/')"
	
	// Build the actions for when we have at least 2 parts (MM/DD or MM/DD/YY)
	formatActions := []string{
		"evt.target.value = " + formatValue,
		d.signals.Set(inputSignal, "evt.target.value"),
	}
	
	// Build actions for when we have complete date (MM/DD/YYYY)
	var dateConversionActions []string
	isoDateExpr := "evt.target.value.split('/')[2] + '-' + evt.target.value.split('/')[0] + '-' + evt.target.value.split('/')[1]"
	dateConversionActions = append(dateConversionActions, d.signals.Set(dateSignal, isoDateExpr))
	
	// Add calendar coordination if needed
	if d.calendarID != "" {
		if strings.Contains(dateSignal, "startDateValue") {
			dateConversionActions = append(dateConversionActions,
				fmt.Sprintf("$%s.rangeStart = %s", d.calendarID, isoDateExpr),
				fmt.Sprintf("$%s.currentDate = evt.target.value.split('/')[2] + '-' + evt.target.value.split('/')[0] + '-01'", d.calendarID),
			)
		} else if strings.Contains(dateSignal, "endDateValue") {
			dateConversionActions = append(dateConversionActions,
				fmt.Sprintf("$%s.rangeEnd = %s", d.calendarID, isoDateExpr),
				fmt.Sprintf("$%s.currentDate = evt.target.value.split('/')[2] + '-' + evt.target.value.split('/')[0] + '-01'", d.calendarID),
			)
		} else if strings.Contains(dateSignal, "dateValue") {
			dateConversionActions = append(dateConversionActions,
				d.signals.Set("selectedDate", isoDateExpr),
				d.signals.Set("currentDate", "evt.target.value.split('/')[2] + '-' + evt.target.value.split('/')[0] + '-01'"),
			)
		}
	}
	
	// Build nested conditional: if has 2+ parts, format, then if has 3 parts, convert to ISO
	nestedCondition := fmt.Sprintf("evt.target.value.split('/').length === 3 && evt.target.value.split('/')[2] ? (%s) : null",
		strings.Join(dateConversionActions, ", "))
	
	formatActions = append(formatActions, nestedCondition)
	
	// Build the main conditional
	return expr.Conditional(
		"evt.target.value.split('/').length >= 2",
		fmt.Sprintf("(%s)", strings.Join(formatActions, ", ")),
		"null",
	).Build()
}

// addCalendarBlurCoordination adds calendar coordination for blur events
func (d *DateInputHandler) addCalendarBlurCoordination(expr *DatastarExpression, dateSignal string) {
	if strings.Contains(dateSignal, "startDateValue") {
		calendarSync := fmt.Sprintf(`if (evt.target.value.split('/').length === 3 && evt.target.value.split('/')[2]) {
			$%s.rangeStart = evt.target.value.split('/')[2] + '-' + evt.target.value.split('/')[0] + '-' + evt.target.value.split('/')[1];
			$%s.currentDate = evt.target.value.split('/')[2] + '-' + evt.target.value.split('/')[0] + '-01';
		}`, d.calendarID, d.calendarID)
		expr.Statement(calendarSync)
	} else if strings.Contains(dateSignal, "endDateValue") {
		calendarSync := fmt.Sprintf(`if (evt.target.value.split('/').length === 3 && evt.target.value.split('/')[2]) {
			$%s.rangeEnd = evt.target.value.split('/')[2] + '-' + evt.target.value.split('/')[0] + '-' + evt.target.value.split('/')[1];
			$%s.currentDate = evt.target.value.split('/')[2] + '-' + evt.target.value.split('/')[0] + '-01';
		}`, d.calendarID, d.calendarID)
		expr.Statement(calendarSync)
	} else if strings.Contains(dateSignal, "dateValue") {
		calendarSync := fmt.Sprintf(`if (evt.target.value.split('/').length === 3 && evt.target.value.split('/')[2]) {
			%s;
			%s;
		}`, 
			d.signals.Set("selectedDate", "evt.target.value.split('/')[2] + '-' + evt.target.value.split('/')[0] + '-' + evt.target.value.split('/')[1]"),
			d.signals.Set("currentDate", "evt.target.value.split('/')[2] + '-' + evt.target.value.split('/')[0] + '-01'"))
		expr.Statement(calendarSync)
	}
}

// BuildTabHandler creates the tab completion handler for range mode
func (d *DateInputHandler) BuildTabHandler(inputSignal, dateSignal string) string {
	return NewExpression().
		Conditional(
			"evt.key === 'Tab' && !evt.shiftKey",
			d.BuildBlurHandler(inputSignal, dateSignal),
			"null",
		).
		Build()
}

// BuildCheckboxChangeHandler creates the checkbox change handler for range mode
func (d *DateInputHandler) BuildCheckboxChangeHandler(endInputID string) string {
	return NewExpression().
		Conditional(
			d.signals.Signal("endDateEnabled"),
			fmt.Sprintf("document.getElementById('%s').focus()", endInputID),
			fmt.Sprintf("(%s, %s, document.getElementById('%s').value = '')",
				d.signals.Set("endInputValue", "''"),
				d.signals.Set("endDateValue", "''"),
				endInputID,
			),
		).
		Build()
}

// DropdownHandler creates handlers for dropdown components
type DropdownHandler struct {
	signals *SignalManager
}

// NewDropdownHandler creates a dropdown handler
func NewDropdownHandler(signals *SignalManager) *DropdownHandler {
	return &DropdownHandler{
		signals: signals,
	}
}

// BuildClickOutsideHandler creates a click outside handler for closing dropdown
func (d *DropdownHandler) BuildClickOutsideHandler() string {
	return d.signals.ConditionalAction(d.signals.Signal("open"), "open", "false")
}

// BuildEscapeHandler creates an escape key handler for closing dropdown
func (d *DropdownHandler) BuildEscapeHandler() string {
	condition := fmt.Sprintf("evt.key === 'Escape' && %s", d.signals.Signal("open"))
	return d.signals.ConditionalAction(condition, "open", "false")
}

// CreateSideClasses generates positioning classes for dropdown sides
func CreateSideClasses(side string, offset int) string {
	if offset == 0 {
		offset = 4 // Default offset like shadcn/ui
	}

	switch side {
	case "top":
		return fmt.Sprintf("bottom-full mb-%d", offset)
	case "bottom":
		return fmt.Sprintf("top-full mt-%d", offset)
	case "left":
		return fmt.Sprintf("right-full mr-%d", offset)
	case "right":
		return fmt.Sprintf("left-full ml-%d", offset)
	default:
		return fmt.Sprintf("top-full mt-%d", offset)
	}
}

// TabsHandler creates handlers for Tabs component functionality
type TabsHandler struct {
	tabsID  string
	signals *SignalManager
}

// NewTabsHandler creates a tabs handler
func NewTabsHandler(tabsID string, signals *SignalManager) *TabsHandler {
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
	
	return NewDataClass().
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

// DialogHandler creates handlers for Dialog component functionality
type DialogHandler struct {
	signals *SignalManager
}

// NewDialogHandler creates a dialog handler
func NewDialogHandler(signals *SignalManager) *DialogHandler {
	return &DialogHandler{
		signals: signals,
	}
}

// BuildBackdropClickHandler creates the backdrop click handler for closing dialog
func (d *DialogHandler) BuildBackdropClickHandler() string {
	return d.signals.ConditionalAction("evt.target === evt.currentTarget", "open", "false")
}

// BuildEscapeHandler creates an escape key handler for closing dialog
func (d *DialogHandler) BuildEscapeHandler() string {
	condition := fmt.Sprintf("evt.key === 'Escape' && %s", d.signals.Signal("open"))
	return d.signals.ConditionalAction(condition, "open", "false")
}

// BuildCloseHandler creates a close handler with optional return value
func (d *DialogHandler) BuildCloseHandler(returnValue string) string {
	expr := NewExpression().Statement(d.signals.Set("open", "false"))
	
	if returnValue != "" {
		expr.Statement(d.signals.SetString("returnValue", returnValue))
	}
	
	return expr.Build()
}

