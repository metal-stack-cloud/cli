package helpers

import (
	"fmt"
	"strings"
)

// The genericcli convert interface returns the uuid of an entity from parsed documents and internally uses this id
// for calling the delete function when using with the "file" flag.
// However, our api requires the project id for deletion as well and when parsed from file the project id gets lost this way.
// To workaround this, we encode the project into the id and in the delete function decode this in case the "file" flag is set.

func EncodeProject(uuid, project string) string {
	return fmt.Sprintf("%s (%s)", uuid, project)
}

func DecodeProject(encoded string) (string, string, error) {
	if !strings.HasSuffix(encoded, ")") {
		return "", "", fmt.Errorf("project id is not encoded into id string")
	}

	uuid, project, ok := strings.Cut(encoded, " (")
	if !ok {
		return "", "", fmt.Errorf("project id is not encoded into id string")
	}

	// cut away the ")"
	project = project[:len(project)-1]

	return uuid, project, nil
}
