package dateinput

import (
	"fmt"
	"strings"

	"github.com/coreycole/datastarui/utils"
)

// DateInputHandler creates handlers for DateInput component functionality
type dateInputHandler struct {
	inputID    string
	signals    *utils.SignalManager
	calendarID string // Optional calendar coordination
}

// NewDateInputHandler creates a DateInput handler
func newDateInputHandler(inputID string, signals *utils.SignalManager, calendarID string) *dateInputHandler {
	return &dateInputHandler{
		inputID:    inputID,
		signals:    signals,
		calendarID: calendarID,
	}
}

// BuildInputHandler creates a simpler input handler that works with data-bind
func (d *dateInputHandler) buildInputHandler(inputSignal, dateSignal string) string {
	expr := utils.NewExpression()

	// Get the raw input value
	expr.Statement("const rawValue = evt.target.value")

	// Apply formatting to convert to MM/DD/YYYY format
	expr.Statement(`const formatValue = (value) => {
		// If the value already contains slashes, handle carefully
		if (value.includes('/')) {
			// Split by slashes to get month, day, year parts
			const parts = value.split('/');
			let month = parts[0] || '';
			let day = parts[1] || '';
			let year = parts[2] || '';
			
			// Remove non-digits from each part
			month = month.replace(/[^\d]/g, '');
			day = day.replace(/[^\d]/g, '');
			year = year.replace(/[^\d]/g, '');
			
			// Special case: if typing in day field and it gets too long, move overflow to year
			// This handles the case where user types "11/112025" (meaning they typed year after day)
			if (day.length > 2 && !year) {
				year = day.substring(2);
				day = day.substring(0, 2);
			}
			
			// Reconstruct with slashes
			if (!month) return '';
			if (!day && !year) return month;
			if (!year) return month + '/' + day;
			return month + '/' + day + '/' + year;
		}
		
		// For values without slashes (like initial typing), apply standard formatting
		const digits = value.replace(/[^\d]/g, '');
		if (digits.length === 0) return '';
		if (digits.length <= 2) return digits;
		if (digits.length <= 4) return digits.substring(0,2) + '/' + digits.substring(2);
		return digits.substring(0,2) + '/' + digits.substring(2,4) + '/' + digits.substring(4,8);
	}`)

	// Format the value
	expr.Statement("const formattedValue = formatValue(rawValue)")

	// Update the input signal with formatted value
	expr.Statement(d.signals.Set(inputSignal, "formattedValue"))

	// Update date signal based on valid date patterns
	fourDigitYear := d.signals.Set(dateSignal, "formattedValue.split('/')[2] + '-' + formattedValue.split('/')[0].padStart(2, '0') + '-' + formattedValue.split('/')[1].padStart(2, '0')")
	twoDigitYear := d.signals.Set(dateSignal, "'20' + formattedValue.split('/')[2] + '-' + formattedValue.split('/')[0].padStart(2, '0') + '-' + formattedValue.split('/')[1].padStart(2, '0')")

	// Set date signal based on completeness
	expr.Statement(fmt.Sprintf(
		"formattedValue.split('/').length === 3 && formattedValue.split('/')[2].length >= 1 ? (formattedValue.split('/')[2].length === 4 ? %s : formattedValue.split('/')[2].length === 2 ? %s : null) : %s",
		fourDigitYear, twoDigitYear, d.signals.Set(dateSignal, "''")),
	)

	// Add calendar coordination if provided - reference existing calendar signals without creating new ones
	if d.calendarID != "" {
		// Create a signal manager that references the existing calendar signals (without creating new ones)
		calendarSignals := &utils.SignalManager{ID: d.calendarID}
		
		expr.Statement(fmt.Sprintf(
			"formattedValue.split('/').length === 3 && (formattedValue.split('/')[2].length === 4 || formattedValue.split('/')[2].length === 2) ? (fullYear => (%s, %s, %s))(formattedValue.split('/')[2].length === 4 ? formattedValue.split('/')[2] : '20' + formattedValue.split('/')[2]) : %s",
			calendarSignals.Set("dateValue", "fullYear + '-' + formattedValue.split('/')[0].padStart(2, '0') + '-' + formattedValue.split('/')[1].padStart(2, '0')"),
			calendarSignals.Set("inputValue", "formattedValue"),
			calendarSignals.Set("currentDate", "fullYear + '-' + formattedValue.split('/')[0].padStart(2, '0') + '-01'"),
			calendarSignals.Set("dateValue", "''"),
		))
	}

	return expr.Build()
}

// addCalendarCoordination adds calendar synchronization logic using signal utilities
func (d *dateInputHandler) addCalendarCoordination(expr *utils.DatastarExpression, dateSignal string) {
	if d.calendarID == "" {
		return
	}
	
	// Create a signal manager that references the existing calendar signals
	calendarSignals := &utils.SignalManager{ID: d.calendarID}
	
	if strings.Contains(dateSignal, "startDateValue") {
		// Start date coordination
		expr.Statement(fmt.Sprintf(
			"evt.target.value.split('/').length === 3 && (evt.target.value.split('/')[2].length === 4 || evt.target.value.split('/')[2].length === 2) ? (fullYear => (%s, %s))(evt.target.value.split('/')[2].length === 4 ? evt.target.value.split('/')[2] : '20' + evt.target.value.split('/')[2]) : %s",
			calendarSignals.Set("rangeStart", "fullYear + '-' + evt.target.value.split('/')[0].padStart(2, '0') + '-' + evt.target.value.split('/')[1].padStart(2, '0')"),
			calendarSignals.Set("currentDate", "fullYear + '-' + evt.target.value.split('/')[0].padStart(2, '0') + '-01'"),
			calendarSignals.Set("rangeStart", "''"),
		))
	} else if strings.Contains(dateSignal, "endDateValue") {
		// End date coordination
		expr.Statement(fmt.Sprintf(
			"evt.target.value.split('/').length === 3 && (evt.target.value.split('/')[2].length === 4 || evt.target.value.split('/')[2].length === 2) ? (fullYear => (%s, %s))(evt.target.value.split('/')[2].length === 4 ? evt.target.value.split('/')[2] : '20' + evt.target.value.split('/')[2]) : %s",
			calendarSignals.Set("rangeEnd", "fullYear + '-' + evt.target.value.split('/')[0].padStart(2, '0') + '-' + evt.target.value.split('/')[1].padStart(2, '0')"),
			calendarSignals.Set("currentDate", "fullYear + '-' + evt.target.value.split('/')[0].padStart(2, '0') + '-01'"),
			calendarSignals.Set("rangeEnd", "''"),
		))
	} else if strings.Contains(dateSignal, "dateValue") {
		// Single date coordination
		expr.Statement(fmt.Sprintf(
			"evt.target.value.split('/').length === 3 && (evt.target.value.split('/')[2].length === 4 || evt.target.value.split('/')[2].length === 2) ? (fullYear => (%s, %s, %s))(evt.target.value.split('/')[2].length === 4 ? evt.target.value.split('/')[2] : '20' + evt.target.value.split('/')[2]) : %s",
			calendarSignals.Set("dateValue", "fullYear + '-' + evt.target.value.split('/')[0].padStart(2, '0') + '-' + evt.target.value.split('/')[1].padStart(2, '0')"),
			calendarSignals.Set("inputValue", "evt.target.value"),
			calendarSignals.Set("currentDate", "fullYear + '-' + evt.target.value.split('/')[0].padStart(2, '0') + '-01'"),
			calendarSignals.Set("dateValue", "''"),
		))
	}
}

// BuildBlurHandler creates the blur completion handler
func (d *dateInputHandler) buildBlurHandler(inputSignal, dateSignal string) string {
	expr := utils.NewExpression()

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

	// Add calendar coordination if needed - use signal utilities
	if d.calendarID != "" {
		// Create a signal manager that references the existing calendar signals
		calendarSignals := &utils.SignalManager{ID: d.calendarID}
		
		if strings.Contains(dateSignal, "startDateValue") {
			dateConversionActions = append(dateConversionActions,
				calendarSignals.Set("rangeStart", isoDateExpr),
				calendarSignals.Set("currentDate", "evt.target.value.split('/')[2] + '-' + evt.target.value.split('/')[0] + '-01'"),
			)
		} else if strings.Contains(dateSignal, "endDateValue") {
			dateConversionActions = append(dateConversionActions,
				calendarSignals.Set("rangeEnd", isoDateExpr),
				calendarSignals.Set("currentDate", "evt.target.value.split('/')[2] + '-' + evt.target.value.split('/')[0] + '-01'"),
			)
		} else if strings.Contains(dateSignal, "dateValue") {
			dateConversionActions = append(dateConversionActions,
				calendarSignals.Set("dateValue", isoDateExpr),
				calendarSignals.Set("inputValue", "evt.target.value"),
				calendarSignals.Set("currentDate", "evt.target.value.split('/')[2] + '-' + evt.target.value.split('/')[0] + '-01'"),
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

// addCalendarBlurCoordination adds calendar coordination for blur events using signal utilities
func (d *dateInputHandler) addCalendarBlurCoordination(expr *utils.DatastarExpression, dateSignal string) {
	if d.calendarID == "" {
		return
	}
	
	// Create a signal manager that references the existing calendar signals
	calendarSignals := &utils.SignalManager{ID: d.calendarID}
	
	if strings.Contains(dateSignal, "startDateValue") {
		calendarSync := fmt.Sprintf(`if (evt.target.value.split('/').length === 3 && evt.target.value.split('/')[2]) {
			%s;
			%s;
		}`,
			calendarSignals.Set("rangeStart", "evt.target.value.split('/')[2] + '-' + evt.target.value.split('/')[0] + '-' + evt.target.value.split('/')[1]"),
			calendarSignals.Set("currentDate", "evt.target.value.split('/')[2] + '-' + evt.target.value.split('/')[0] + '-01'"))
		expr.Statement(calendarSync)
	} else if strings.Contains(dateSignal, "endDateValue") {
		calendarSync := fmt.Sprintf(`if (evt.target.value.split('/').length === 3 && evt.target.value.split('/')[2]) {
			%s;
			%s;
		}`,
			calendarSignals.Set("rangeEnd", "evt.target.value.split('/')[2] + '-' + evt.target.value.split('/')[0] + '-' + evt.target.value.split('/')[1]"),
			calendarSignals.Set("currentDate", "evt.target.value.split('/')[2] + '-' + evt.target.value.split('/')[0] + '-01'"))
		expr.Statement(calendarSync)
	} else if strings.Contains(dateSignal, "dateValue") {
		calendarSync := fmt.Sprintf(`if (evt.target.value.split('/').length === 3 && evt.target.value.split('/')[2]) {
			%s;
			%s;
			%s;
		}`,
			calendarSignals.Set("dateValue", "evt.target.value.split('/')[2] + '-' + evt.target.value.split('/')[0] + '-' + evt.target.value.split('/')[1]"),
			calendarSignals.Set("inputValue", "evt.target.value"),
			calendarSignals.Set("currentDate", "evt.target.value.split('/')[2] + '-' + evt.target.value.split('/')[0] + '-01'"))
		expr.Statement(calendarSync)
	}
}

// BuildTabHandler creates the tab completion handler for range mode
func (d *dateInputHandler) buildTabHandler(inputSignal, dateSignal string) string {
	return utils.NewExpression().
		Conditional(
			"evt.key === 'Tab' && !evt.shiftKey",
			d.buildBlurHandler(inputSignal, dateSignal),
			"null",
		).
		Build()
}

// BuildCheckboxChangeHandler creates the checkbox change handler for range mode
func (d *dateInputHandler) buildCheckboxChangeHandler(endInputID string) string {
	// When coordinating with calendar, use direct signal references without creating new ones
	if d.calendarID != "" {
		// Create a signal manager that references the existing calendar signals
		calendarSignals := &utils.SignalManager{ID: d.calendarID}
		return utils.NewExpression().
			Conditional(
				calendarSignals.Signal("endDateEnabled"),
				fmt.Sprintf("document.getElementById('%s').focus()", endInputID),
				fmt.Sprintf("(%s, %s, document.getElementById('%s').value = '')",
					calendarSignals.Set("endInputValue", "''"),
					calendarSignals.Set("endDateValue", "''"),
					endInputID,
				),
			).
			Build()
	}
	
	// Standalone mode - create own signals
	return utils.NewExpression().
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
