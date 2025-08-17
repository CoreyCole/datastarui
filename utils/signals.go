package utils

import (
	"encoding/json"
	"fmt"
	"strings"
)

// SignalManager provides a structured way to manage Datastar signals
// It namespaces signals by ID and provides serialization capabilities
type SignalManager struct {
	ID          string      `json:"id"`
	Signals     any         `json:"signals"`
	DataSignals string      `json:"dataSignals"`
}

// Signals creates a new SignalManager instance with the given ID and signals struct
// The signalsStruct should be any struct that has json tags for each property
// The ID will be sanitized to replace hyphens with underscores for JavaScript compatibility
// Example:
//
//	type MySignals struct {
//	    Open     bool   `json:"open"`
//	    Value    string `json:"value"`
//	    Count    int    `json:"count"`
//	}
//	signals := Signals("myComponent", MySignals{Open: false, Value: "", Count: 0})
//	// Use in templ: data-signals={ signals.DataSignals }
func Signals(id string, signalsStruct any) *SignalManager {
	// Sanitize ID by replacing hyphens with underscores for JavaScript compatibility
	sanitizedID := strings.ReplaceAll(id, "-", "_")

	// Create the nested structure: {[sanitizedID]: signalsStruct}
	nested := map[string]any{
		sanitizedID: signalsStruct,
	}

	// Marshal to JSON
	jsonBytes, err := json.Marshal(nested)
	if err != nil {
		// Fallback to empty object if marshaling fails
		jsonBytes = []byte("{}")
	}

	return &SignalManager{
		ID:          sanitizedID,
		Signals:     signalsStruct,
		DataSignals: string(jsonBytes),
	}
}

// Signal returns a reference to a specific signal property
// Example: signals.Signal("open") returns "$myComponent.open"
func (sm *SignalManager) Signal(property string) string {
	return fmt.Sprintf("$%s.%s", sm.ID, property)
}

// Toggle returns a toggle expression for a boolean signal property
// Example: signals.Toggle("open") returns "$myComponent.open = !$myComponent.open"
func (sm *SignalManager) Toggle(property string) string {
	ref := sm.Signal(property)
	return fmt.Sprintf("%s = !%s", ref, ref)
}

// Set returns a set expression for a signal property
// Example: signals.Set("value", "'hello'") returns "$myComponent.value = 'hello'"
func (sm *SignalManager) Set(property, value string) string {
	return fmt.Sprintf("%s = %s", sm.Signal(property), value)
}

// SetString returns a set expression for a string signal property with proper quoting
// Example: signals.SetString("value", "hello") returns "$myComponent.value = 'hello'"
func (sm *SignalManager) SetString(property, value string) string {
	return fmt.Sprintf("%s = '%s'", sm.Signal(property), value)
}

// Conditional returns a conditional expression for a signal property
// Example: signals.Conditional("loading", "Saving...", "Save") returns "$myComponent.loading ? 'Saving...' : 'Save'"
func (sm *SignalManager) Conditional(property, trueValue, falseValue string) string {
	return fmt.Sprintf("%s ? %s : %s", sm.Signal(property), trueValue, falseValue)
}

// ConditionalAction creates a safe conditional action expression using ternary operator
// Example: signals.ConditionalAction("evt.target === evt.currentTarget", "open", "false")
// Returns: "evt.target === evt.currentTarget ? ($component.open = false) : void 0"
func (sm *SignalManager) ConditionalAction(condition, property, value string) string {
	return fmt.Sprintf("%s ? (%s) : void 0", condition, sm.Set(property, value))
}

// ConditionalMultiAction creates a safe conditional expression with multiple actions using ternary operator
// Example: signals.ConditionalMultiAction("condition", []string{"action1", "action2"})
// Returns: "condition ? (action1, action2) : void 0"
func (sm *SignalManager) ConditionalMultiAction(condition string, actions ...string) string {
	if len(actions) == 0 {
		return ""
	}
	actionsStr := ""
	for i, action := range actions {
		if i > 0 {
			actionsStr += ", "
		}
		actionsStr += action
	}
	return fmt.Sprintf("%s ? (%s) : void 0", condition, actionsStr)
}


// MultiStateConditional creates a chain of conditional expressions for handling multiple states
//
//	Example: signals.MultiStateConditional([]StateAction{
//	  {Condition: "!$comp.start", Actions: []string{"setStart", "clearEnd"}},
//	  {Condition: "!$comp.end", Actions: []string{"setEnd", "dispatchEvent"}},
//	  {Condition: "true", Actions: []string{"resetRange"}},
//	})
type StateAction struct {
	Condition string
	Actions   []string
}

func (sm *SignalManager) MultiStateConditional(states []StateAction) string {
	if len(states) == 0 {
		return ""
	}

	result := ""
	for i, state := range states {
		if i > 0 {
			result += " : "
		}

		// Format actions
		actionsStr := ""
		for j, action := range state.Actions {
			if j > 0 {
				actionsStr += ", "
			}
			actionsStr += action
		}

		// Handle the final condition specially - if it's "true", just execute the actions
		isLastCondition := i == len(states)-1
		if isLastCondition && state.Condition == "true" {
			// For the final "true" condition, just execute the actions without ternary
			if len(state.Actions) == 1 {
				result += actionsStr
			} else {
				result += fmt.Sprintf("(%s)", actionsStr)
			}
		} else {
			// Build conditional normally
			if len(state.Actions) == 1 {
				result += fmt.Sprintf("%s ? %s", state.Condition, actionsStr)
			} else {
				result += fmt.Sprintf("%s ? (%s)", state.Condition, actionsStr)
			}
		}
	}

	return result
}

// RangeSelection creates the standard range selection logic pattern
// This handles the common pattern of: no start -> set start, has start -> complete range, has both -> reset
// startProp, endProp are the signal property names (e.g., "rangeStart", "rangeEnd")
// clickedValue is the expression for the clicked value (e.g., "dayDateStr")
// eventDetail is the detail object for the custom event (optional)
func (sm *SignalManager) RangeSelection(startProp, endProp, clickedValue, eventDetail string) string {
	startRef := sm.Signal(startProp)
	endRef := sm.Signal(endProp)

	// Default event detail if not provided
	if eventDetail == "" {
		eventDetail = fmt.Sprintf("{ %s: %s, %s: %s }", startProp, startRef, endProp, endRef)
	}

	states := []StateAction{
		{
			Condition: fmt.Sprintf("!%s", startRef),
			Actions: []string{
				sm.Set(startProp, clickedValue),
				sm.Set(endProp, "''"),
			},
		},
		{
			Condition: fmt.Sprintf("!%s", endRef),
			Actions: []string{
				fmt.Sprintf("%s < %s ? (oldStart = %s, %s, %s) : (%s)",
					clickedValue, startRef,
					startRef, sm.Set(startProp, clickedValue), sm.Set(endProp, "oldStart"),
					sm.Set(endProp, clickedValue)),
				fmt.Sprintf("this.dispatchEvent(new CustomEvent('calendar-change', { bubbles: true, detail: %s }))", eventDetail),
			},
		},
		{
			Condition: "true",
			Actions: []string{
				sm.Set(startProp, clickedValue),
				sm.Set(endProp, "''"),
			},
		},
	}

	return sm.MultiStateConditional(states)
}

// SingleOrRange creates a conditional that handles both single and range selection modes
// modeProp is the signal property for the mode (e.g., "mode")
// singleActions are the actions for single mode
// rangeExpression is the complete range selection expression
func (sm *SignalManager) SingleOrRange(modeProp string, singleActions []string, rangeExpression string) string {
	modeRef := sm.Signal(modeProp)

	// Format single actions
	singleActionsStr := ""
	for i, action := range singleActions {
		if i > 0 {
			singleActionsStr += ", "
		}
		singleActionsStr += action
	}

	if len(singleActions) == 1 {
		return fmt.Sprintf("%s === 'single' ? %s : (%s)", modeRef, singleActionsStr, rangeExpression)
	} else {
		return fmt.Sprintf("%s === 'single' ? (%s) : (%s)", modeRef, singleActionsStr, rangeExpression)
	}
}

// DateComparison creates date comparison expressions with proper handling
// Useful for calendar date logic where you need to compare date strings
func (sm *SignalManager) DateComparison(date1, operator, date2 string) string {
	switch operator {
	case "<", "<=", ">", ">=":
		return fmt.Sprintf("%s %s %s", date1, operator, date2)
	case "==", "===":
		return fmt.Sprintf("%s === %s", date1, date2)
	case "!=", "!==":
		return fmt.Sprintf("%s !== %s", date1, date2)
	default:
		return fmt.Sprintf("%s %s %s", date1, operator, date2)
	}
}

// DataClass creates a clean JSON object for data-class attributes from a map of class names to conditions
// This allows for more maintainable conditional class logic
// Example:
//
//	classes := signals.DataClass(map[string]string{
//	  "bg-primary text-white": "$component.active",
//	  "opacity-50": "!$component.enabled",
//	  "hidden": "$component.mode === 'hidden'",
//	})
//	// Use in templ: data-class={ classes }
func (sm *SignalManager) DataClass(classConditions map[string]string) string {
	if len(classConditions) == 0 {
		return "{}"
	}

	var parts []string
	for className, condition := range classConditions {
		// Escape single quotes in class names
		escapedClass := strings.ReplaceAll(className, "'", "\\'")
		parts = append(parts, fmt.Sprintf("'%s': %s", escapedClass, condition))
	}

	return fmt.Sprintf("{%s}", strings.Join(parts, ", "))
}

// Equals creates a comparison expression between a signal and a value
// Example: signals.Equals("value", "option1") returns "$component.value === 'option1'"
func (sm *SignalManager) Equals(property, value string) string {
	return fmt.Sprintf("%s === '%s'", sm.Signal(property), value)
}

// NotEquals creates a not-equals comparison expression
// Example: signals.NotEquals("value", "option1") returns "$component.value !== 'option1'"
func (sm *SignalManager) NotEquals(property, value string) string {
	return fmt.Sprintf("%s !== '%s'", sm.Signal(property), value)
}

// TernaryClass creates a ternary expression for conditional CSS classes
// Example: signals.TernaryClass("checked", "bg-primary", "bg-secondary")
// Returns: "$component.checked ? 'bg-primary' : 'bg-secondary'"
func (sm *SignalManager) TernaryClass(property, trueClass, falseClass string) string {
	return fmt.Sprintf("%s ? '%s' : '%s'", sm.Signal(property), trueClass, falseClass)
}

// TernaryStyle creates a ternary expression for conditional inline styles
// Example: signals.TernaryStyle("visible", "opacity: 1", "opacity: 0")
// Returns: "$component.visible ? 'opacity: 1' : 'opacity: 0'"
func (sm *SignalManager) TernaryStyle(property, trueStyle, falseStyle string) string {
	return fmt.Sprintf("%s ? '%s' : '%s'", sm.Signal(property), trueStyle, falseStyle)
}
