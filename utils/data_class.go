package utils

import (
	"fmt"
	"strings"
)

// DataClass helps build Datastar data-class attribute expressions
// It generates object syntax for conditional CSS classes
type DataClass struct {
	classes map[string]string // className -> condition
}

// NewDataClass creates a new DataClass builder
func NewDataClass() *DataClass {
	return &DataClass{
		classes: make(map[string]string),
	}
}

// Add adds a conditional class
func (d *DataClass) Add(className, condition string) *DataClass {
	d.classes[className] = condition
	return d
}

// AddMultiple adds multiple classes with the same condition
func (d *DataClass) AddMultiple(classNames []string, condition string) *DataClass {
	for _, className := range classNames {
		d.classes[className] = condition
	}
	return d
}

// Build creates the data-class object expression
func (d *DataClass) Build() string {
	if len(d.classes) == 0 {
		return "{}"
	}

	var parts []string
	for className, condition := range d.classes {
		parts = append(parts, fmt.Sprintf("'%s': %s", className, condition))
	}

	return "{" + strings.Join(parts, ", ") + "}"
}

// HighlightedItem creates a data-class expression for highlighted select/list items
func HighlightedItem(signalPath string, index int) string {
	return NewDataClass().
		Add("bg-accent", fmt.Sprintf("$%s === %d", signalPath, index)).
		Add("text-accent-foreground", fmt.Sprintf("$%s === %d", signalPath, index)).
		Build()
}

// SelectedItem creates a data-class expression for selected items
func SelectedItem(signalPath, value string) string {
	return NewDataClass().
		Add("bg-primary", fmt.Sprintf("$%s === '%s'", signalPath, value)).
		Add("text-primary-foreground", fmt.Sprintf("$%s === '%s'", signalPath, value)).
		Build()
}

// ActiveTab creates a data-class expression for active tab styling
func ActiveTab(signalPath, value string) string {
	return NewDataClass().
		Add("border-b-2", fmt.Sprintf("$%s === '%s'", signalPath, value)).
		Add("border-primary", fmt.Sprintf("$%s === '%s'", signalPath, value)).
		Build()
}

// OpenState creates a data-class expression based on open/closed state
func OpenState(signalPath string, openClasses, closedClasses []string) string {
	dc := NewDataClass()
	
	// Add open state classes
	for _, className := range openClasses {
		dc.Add(className, fmt.Sprintf("$%s", signalPath))
	}
	
	// Add closed state classes
	for _, className := range closedClasses {
		dc.Add(className, fmt.Sprintf("!$%s", signalPath))
	}
	
	return dc.Build()
}