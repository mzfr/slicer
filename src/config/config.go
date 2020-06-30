package config

// Configurations exported
type Configurations struct {
	paths PathsConfiguration
	URLs  URLConfiguration
}

// PathsConfiguration exported
type PathsConfiguration struct {
	manifest string
	strings  string
}

// URLConfiguration exported
type URLConfiguration struct {
	maps       string
	streetview string
	directions string
	places     string
	geocoding  string
}
