package config

// Configurations exported
type Configurations struct {
	cmd      string
	checks   ChecksConfiguration
	patterns PatternConfiguration
}

// ChecksConfiguration exported
type ChecksConfiguration struct {
	firebase bool
	gmap     bool
}

// PatternConfiguration exported
type PatternConfiguration struct {
	// TODO: Add pattern for general API and google map API key
	URLS       string
	AwsKeys    string
	firbaseURL string
	google     string
}
