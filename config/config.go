package config

// Configurations exported
type Configurations struct {
	paths PathsConfiguration
	URLs  map[string][]URLConfiguration
}

// PathsConfiguration exported
type PathsConfiguration struct {
	manifest string
	strings  string
	raw      string
	xml      string
}

// URLConfiguration exported
type URLConfiguration struct {
	maps       string
	streetview string
	directions string
	places     string
	geocoding  string
}
