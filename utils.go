package nbmysql

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"
)

var removeParenthesisRe = regexp.MustCompile(`\(.*?\)`)

func GetLastInsertId(db *sql.DB) (int64, error) {
	id := new(Int)
	row := db.QueryRow("SELECT LAST_INSERT_ID()")
	err := row.Scan(id)
	if err != nil {
		return -1, err
	}
	return id.Value, nil
}

func BackQuote(s string) string {
	return fmt.Sprintf("`%s`", s)
}

func ToCap(name string) string {
	l := strings.Split(name, "_")
	for i, n := range l {
		l[i] = strings.Title(n)
	}
	return strings.Join(l, "")
}

func ToLocal(name string) string {
	l := strings.Split(name, "_")
	for i, n := range l[1:] {
		l[i+1] = strings.Title(n)
	}
	return strings.Join(l, "")
}

func RmParenthesis(s string) string {
	return removeParenthesisRe.ReplaceAllString(s, "")
}

func getColumn(colName string, tab *Table) *Column {
	for i, col := range tab.Columns {
		if col.ColumnName == colName {
			return &tab.Columns[i]
		}
	}
	return nil
}
