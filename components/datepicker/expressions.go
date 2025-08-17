package datepicker

import (
	"fmt"

	"github.com/coreycole/datastarui/utils"
)

// datePickerPopoverHandler creates handlers for DatePicker popover functionality
type datePickerPopoverHandler struct {
	datePickerID string
	dateInputID  string
	mode         string
	signals      *utils.SignalManager
}

// NewdatePickerPopoverHandler creates a datepicker popover handler
func newDatePickerPopoverHandler(datePickerID, dateInputID, mode string, signals *utils.SignalManager) *datePickerPopoverHandler {
	return &datePickerPopoverHandler{
		datePickerID: datePickerID,
		dateInputID:  dateInputID,
		mode:         mode,
		signals:      signals,
	}
}

// BuildEscapeHandler creates clean escape key handler for popover
func (d *datePickerPopoverHandler) buildEscapeHandler() string {
	openCheck := fmt.Sprintf("document.querySelector('[data-datepicker-id=\"%s\"]') && %s", d.datePickerID, d.signals.Signal("open"))

	inputID := d.dateInputID
	if d.mode == "range" {
		inputID += "_start" // Focus on start input in range mode
	}

	return utils.NewExpression().
		Conditional(
			fmt.Sprintf("evt.key === 'Escape' && %s", openCheck),
			fmt.Sprintf("(evt.preventDefault(), evt.stopPropagation(), %s, document.getElementById('%s').focus())",
				d.signals.Set("open", "false"), inputID),
			"null",
		).
		Build()
}

// BuildTabHandler creates clean tab key handler for popover
func (d *datePickerPopoverHandler) buildTabHandler() string {
	openCheck := fmt.Sprintf("document.querySelector('[data-datepicker-id=\"%s\"]') && %s", d.datePickerID, d.signals.Signal("open"))

	return utils.NewExpression().
		Conditional(
			fmt.Sprintf("evt.key === 'Tab' && %s", openCheck),
			d.signals.Set("open", "false"),
			"null",
		).
		Build()
}

// BuildKeyboardHandler creates the combined keyboard handler
func (d *datePickerPopoverHandler) buildKeyboardHandler() string {
	escapeHandler := d.buildEscapeHandler()
	tabHandler := d.buildTabHandler()
	return escapeHandler + "; " + tabHandler
}

// BuildDateSelectHandler creates the complex date selection handler with DateInput sync
func (d *datePickerPopoverHandler) buildDateSelectHandler(closeOnSelect bool, dateInputSignals *utils.SignalManager) string {
	expr := utils.NewExpression()

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
		singleDateCondition := "evt.detail.dateValue"
		singleDateActions := fmt.Sprintf(
			"((displayDate) => (%s, %s, %s))(evt.detail.dateValue.replace(/-/g, '/'))",
			dateInputSignals.Set("inputValue", "displayDate"),
			dateInputSignals.Set("dateValue", "evt.detail.dateValue"),
			d.signals.Set("dateValue", "evt.detail.dateValue"),
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
func (d *datePickerPopoverHandler) buildMonthChangeHandler() string {
	return utils.NewExpression().
		Conditional(
			"evt.detail.displayMonth",
			d.signals.Set("displayMonth", "evt.detail.displayMonth"),
			"null",
		).
		Build()
}

// BuildOpenTriggerHandler creates the calendar icon click handler
func (d *datePickerPopoverHandler) buildOpenTriggerHandler() string {
	// Toggle open state - removed stopPropagation to allow click-outside handlers to work
	return d.signals.Toggle("open")
}

// BuildClickOutsideHandler creates the click outside handler to close popover
func (d *datePickerPopoverHandler) buildClickOutsideHandler() string {
	// Only close when open AND the click target is not THIS datepicker's calendar button
	// Allow other calendar buttons to close this popover for coordination
	datepickerSelector := fmt.Sprintf("[data-datepicker-id=\"%s\"]", d.datePickerID)
	return d.signals.Signal("open") + " && !evt.target.closest('" + datepickerSelector + " button[data-on-click*=\"open\"]') ? " + d.signals.Set("open", "false") + " : null"
}
