package printer

import (
	"sort"
	"strconv"

	adminv1 "github.com/metal-stack-cloud/api/go/admin/v1"
)

func (t *TablePrinter) UserTable(data []*adminv1.User, wide bool) ([]string, [][]string, error) {
	var (
		rows [][]string
	)

	header := []string{"ID", "Name", "Email", "Admitted"}
	if wide {
		header = []string{"ID", "Name", "Email", "Admitted"}
	}

	sort.SliceStable(data, func(i, j int) bool { return data[i].User.Login < data[j].User.Login })
	for _, user := range data {
		id := user.User.Login
		name := user.User.Name
		email := user.User.Email
		admitted := strconv.FormatBool(user.User.Admitted)

		if wide {
			rows = append(rows, []string{id, name, email, admitted})
		} else {
			rows = append(rows, []string{id, name, email, admitted})
		}
	}

	return header, rows, nil
}
