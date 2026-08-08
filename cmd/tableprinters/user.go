package tableprinters

import (
	"strings"

	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
)

type Named interface {
	GetName() string
}

func (t *TablePrinter) UserTable(data []*apiv1.User, wide bool) ([]string, [][]string, error) {
	var (
		rows   [][]string
		header = []string{"Login", "Name", "Email"}
	)

	if wide {
		header = []string{"Login", "Name", "Email", "Tenants", "Projects"}
	}

	for _, user := range data {
		login := user.Login
		name := user.Name
		email := user.Email

		if wide {
			rows = append(rows, []string{login, name, email, namesString(user.Tenants), namesString(user.Projects)})

		} else {
			rows = append(rows, []string{login, name, email})
		}
	}

	return header, rows, nil
}

func namesString[T Named](arr []T) string {
	names := make([]string, 0, len(arr))

	for _, t := range arr {
		name := t.GetName()
		if name != "" {
			names = append(names, name)
		}
	}

	return strings.Join(names, ", ")
}
