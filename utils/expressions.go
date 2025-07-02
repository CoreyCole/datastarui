package utils

import (
	"fmt"
	"strings"
)

// DatastarExpression represents a structured Datastar expression builder
type DatastarExpression struct {
	statements []string
	separator  string
}

// NewExpression creates a new Datastar expression builder with semicolon separator
func NewExpression() *DatastarExpression {
	return &DatastarExpression{
		statements: make([]string, 0),
		separator:  "; ",
	}
}

// Statement adds a statement to the expression
func (e *DatastarExpression) Statement(stmt string) *DatastarExpression {
	if stmt != "" {
		e.statements = append(e.statements, stmt)
	}
	return e
}

// SetSignal adds a signal assignment statement
func (e *DatastarExpression) SetSignal(signal, value string) *DatastarExpression {
	return e.Statement(fmt.Sprintf("$%s = %s", signal, value))
}

// Conditional adds a conditional statement
func (e *DatastarExpression) Conditional(condition, trueExpr, falseExpr string) *DatastarExpression {
	if falseExpr == "" {
		falseExpr = "null"
	}
	return e.Statement(fmt.Sprintf("%s ? %s : %s", condition, trueExpr, falseExpr))
}

// Build returns the final expression string
func (e *DatastarExpression) Build() string {
	if len(e.statements) == 0 {
		return ""
	}
	if len(e.statements) == 1 {
		return e.statements[0]
	}
	// Wrap multiple statements in parentheses if using comma separator
	if e.separator == ", " {
		return "(" + strings.Join(e.statements, e.separator) + ")"
	}
	return strings.Join(e.statements, e.separator)
}


// BuildConditional creates a conditional expression
func BuildConditional(condition, trueExpr, falseExpr string) string {
	if falseExpr == "" {
		falseExpr = "null"
	}
	return fmt.Sprintf("%s ? %s : %s", condition, trueExpr, falseExpr)
}

// FocusCapture creates capture handlers
type FocusCapture struct {
	targetSelector string
	action         string
}

// NewFocusCapture creates a focus capture handler
func NewFocusCapture() *FocusCapture {
	return &FocusCapture{}
}

// OnlyInputs restricts to input elements only
func (f *FocusCapture) OnlyInputs() *FocusCapture {
	f.targetSelector = "evt.target.tagName.toLowerCase() === 'input'"
	return f
}

// OnSelector restricts to specific selector
func (f *FocusCapture) OnSelector(selector string) *FocusCapture {
	f.targetSelector = fmt.Sprintf("evt.target.matches('%s')", selector)
	return f
}

// SetSignal sets a signal when focus condition is met
func (f *FocusCapture) SetSignal(signals *SignalManager, signal, value string) *FocusCapture {
	f.action = signals.Set(signal, value)
	return f
}

// Build creates the focus capture handler
func (f *FocusCapture) Build() string {
	if f.targetSelector == "" {
		return f.action
	}
	return BuildConditional(f.targetSelector, f.action, "null")
}
