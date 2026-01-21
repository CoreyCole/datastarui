package registry

// Radius represents a border-radius option
type Radius struct {
	Name  string
	Label string
	Value string
}

// RADII contains all available border-radius options
var RADII = []Radius{
	{Name: "default", Label: "Default", Value: ""},
	{Name: "none", Label: "None", Value: "0"},
	{Name: "small", Label: "Small", Value: "0.45rem"},
	{Name: "medium", Label: "Medium", Value: "0.625rem"},
	{Name: "large", Label: "Large", Value: "0.875rem"},
}

var radiiMap = map[string]*Radius{}

func init() {
	for i := range RADII {
		radiiMap[RADII[i].Name] = &RADII[i]
	}
}

// GetRadius returns a radius option by name
func GetRadius(name string) *Radius { return radiiMap[name] }

// AllRadii returns all available radius options
func AllRadii() []Radius { return RADII }
