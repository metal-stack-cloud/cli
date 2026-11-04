package config

import "slices"

const (
	DefaultApiURL         = "https://api.metalstack.cloud"
	DefaultAfterLoginPage = "https://metalstack.cloud"
	DefaultConsoleURL     = "https://console.metalstack.cloud"

	// Access level admin used for shoot kubeconfig generation
	AccessLevelAdmin = "admin"
	// Access level viewer used for shoot kubeconfig generation
	AccessLevelViewer = "viewer"
)

var AccessLevels = []string{AccessLevelAdmin, AccessLevelViewer}

func IsValidAccessLevel(level string) bool {
	return slices.Contains(AccessLevels, level)
}
