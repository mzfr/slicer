package config

// Configurations exported
type Configurations struct {
	paths    PathsConfiguration
	patterns PatternConfiguration
}

// PathsConfiguration exported
type PathsConfiguration struct {
	manifest string
	strings  string
}

// PatternConfiguration exported
type PatternConfiguration struct {
	// TODO: Add pattern for general API and google map API key
	URLS       string
	AwsKeys    string
	firbaseURL string
	google     string
}
