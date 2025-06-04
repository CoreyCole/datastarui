package utils

import (
	"encoding/json"
	"fmt"
	"strings"
)

// SignalManager provides a structured way to manage Datastar signals
// It namespaces signals by ID and provides serialization capabilities
type SignalManager struct {
	ID          string      `json:"-"`
	Signals     interface{} `json:"-"`
	DataSignals string      `json:"-"`
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
func Signals(id string, signalsStruct interface{}) *SignalManager {
	// Sanitize ID by replacing hyphens with underscores for JavaScript compatibility
	sanitizedID := strings.ReplaceAll(id, "-", "_")

	// Create the nested structure: {[sanitizedID]: signalsStruct}
	nested := map[string]interface{}{
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
