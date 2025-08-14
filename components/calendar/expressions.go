package calendar

import (
	"fmt"

	"github.com/coreycole/datastarui/utils"
)

type calendarDayHandler struct {
	calendarID  string
	signals     *utils.SignalManager
	monthOffset int
	mode        string
	dayNumber   int
}

func newCalendarDayHandler(calendarID string, signals *utils.SignalManager, monthOffset int, mode string, dayNumber int) *calendarDayHandler {
	return &calendarDayHandler{
		calendarID:  calendarID,
		signals:     signals,
		monthOffset: monthOffset,
		mode:        mode,
		dayNumber:   dayNumber,
	}
}

// BuildClickHandler creates the day click handler expression
func (c *calendarDayHandler) buildClickHandler() string {
	expr := utils.NewExpression()

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
func (c *calendarDayHandler) buildDateInputSync(datePickerInputsID string) string {
	expr := utils.NewExpression()

	if c.mode == "range" {
		// Sync range signals for DateInput - use the same namespace since datePickerInputsID === calendarID
		expr.Statement(c.signals.Set("startDateValue", c.signals.Signal("rangeStart")))
		expr.Statement(c.signals.Set("endDateValue", c.signals.Signal("rangeEnd")))
		expr.Statement(c.signals.Set("startInputValue", c.signals.Signal("rangeStart")+" ? new Date("+c.signals.Signal("rangeStart")+" + 'T12:00:00Z').toLocaleDateString('en-US', {month: '2-digit', day: '2-digit', year: 'numeric', timeZone: 'UTC'}) : ''"))
		expr.Statement(c.signals.Set("endInputValue", c.signals.Signal("rangeEnd")+" ? new Date("+c.signals.Signal("rangeEnd")+" + 'T12:00:00Z').toLocaleDateString('en-US', {month: '2-digit', day: '2-digit', year: 'numeric', timeZone: 'UTC'}) : ''"))
	} else {
		// Sync single date signals for DateInput - use the same namespace since datePickerInputsID === calendarID
		expr.Statement(c.signals.Set("dateValue", c.signals.Signal("selectedDate")))
		expr.Statement(c.signals.Set("inputValue", c.signals.Signal("selectedDate")+" ? new Date("+c.signals.Signal("selectedDate")+" + 'T12:00:00Z').toLocaleDateString('en-US', {month: '2-digit', day: '2-digit', year: 'numeric', timeZone: 'UTC'}) : ''"))
	}

	return expr.Build()
}

// CalendarSelectionClasses creates conditional classes for calendar day selection
type calendarSelectionClasses struct {
	calendarID  string
	signals     *utils.SignalManager
	monthOffset int
	dayNumber   int
	mode        string
}

// NewCalendarSelectionClasses creates a calendar selection class helper
func newCalendarSelectionClasses(calendarID string, signals *utils.SignalManager, monthOffset int, dayNumber int, mode string) *calendarSelectionClasses {
	return &calendarSelectionClasses{
		calendarID:  calendarID,
		signals:     signals,
		monthOffset: monthOffset,
		dayNumber:   dayNumber,
		mode:        mode,
	}
}

// Build creates the data-class expression for calendar day selection
func (c *calendarSelectionClasses) build() string {
	dataClass := utils.NewDataClass()

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
		dataClass.Add("bg-primary text-primary-foreground",
			fmt.Sprintf("%s === (%s) || %s === (%s)",
				c.signals.Signal("rangeStart"), expectedDateExpr,
				c.signals.Signal("rangeEnd"), expectedDateExpr,
			))

		dataClass.Add("bg-accent text-accent-foreground",
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
		dataClass.Add("bg-primary text-primary-foreground",
			fmt.Sprintf("%s === (%s)", c.signals.Signal("selectedDate"), expectedDateExpr))
	}

	return dataClass.Build()
}
