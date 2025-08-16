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

