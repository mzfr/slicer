package output

// Info exported
type Info struct {
	Debuggable  string
	AllowBackup string
	Activity    Components
	receiver    Components
	service     Components
	strings     []interface{}
}

// Components exported
type Components struct {
	Name          string
	Permission    string
	IntentFilters Filters
}

// Filters exported
type Filters struct {
	Data       []string
	Action     []interface{}
	Components []interface{}
}
