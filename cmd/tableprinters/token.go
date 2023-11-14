package tableprinters

import (
	"time"

	apiv1 "github.com/metal-stack-cloud/api/go/api/v1"
)

func (t *TablePrinter) TokenTable(data []*apiv1.Token, _ bool) ([]string, [][]string, error) {
	var (
		rows [][]string
	)
	header := []string{"ID", "User", "Description", "Expires"}

	for _, token := range data {
		row := []string{token.Uuid, token.UserId, token.Description, "in " + time.Until(token.Expires.AsTime()).String()}

		rows = append(rows, row)
	}

	return header, rows, nil
}
