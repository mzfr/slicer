package output

// Info exported
type Info struct {
	permission, action, category string
}

// JSONOutput exported
type JSONOutput struct {
	packageName, Debuggable, AllowBackup string
	activites, receivers, services       []Info
}
