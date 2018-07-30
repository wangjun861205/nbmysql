package nbmysql

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"
)

var removeParenthesisRe = regexp.MustCompile(`\(.*?\)`)

func GetLastId(db *sql.DB) int64 {
	id := new(Int)
	row := db.QueryRow("SELECT LAST_INSERT_ID()")
	err := row.Scan(id)
	if err != nil {
		panic(err)
	}
	return id.Value
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
