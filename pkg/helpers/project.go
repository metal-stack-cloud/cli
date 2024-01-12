package helpers

import "strings"

func TrimProvider(subject string) string {
	if !strings.Contains(subject, "@") {
		return subject
	}

	parts := strings.Split(subject, "@")

	return strings.Join(parts[:len(parts)-1], "@")
}
