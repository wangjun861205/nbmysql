package nbmysql

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

var removeParenthesisRe = regexp.MustCompile(`\(.*?\)`)

// func GetLastInsertId(db *sql.DB) (int64, error) {
// 	id := new(Int)
// 	row := db.QueryRow("SELECT LAST_INSERT_ID()")
// 	err := row.Scan(id)
// 	if err != nil {
// 		return -1, err
// 	}
// 	return id.Value, nil
// }

//BackQuote wrap string within "`"
func BackQuote(s string) string {
	return fmt.Sprintf("`%s`", s)
}

//ToCap convert xxx_yyy to XxxYyy
func ToCap(name string) string {
	l := strings.Split(name, "_")
	for i, n := range l {
		l[i] = strings.Title(n)
	}
	return strings.Join(l, "")
}

//ToLocal convert xxx_yyy to xxxYyy
func ToLocal(name string) string {
	l := strings.Split(name, "_")
	for i, n := range l[1:] {
		l[i+1] = strings.Title(n)
	}
	return strings.Join(l, "")
}

//RmParenthesis remove parentheist from string
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

//StringArg shortcut for string arguments in NewXXX() function
func StringArg(s string) *string {
	return &s
}

//IntArg shortcut for int arguments in NewXXX() function
func IntArg(i int64) *int64 {
	return &i
}

//FloatArg shortcut for float arguments in NewXXX() function
func FloatArg(f float64) *float64 {
	return &f
}

//BoolArg shortcut for bool arguments in NewXXX() function
func BoolArg(b bool) *bool {
	return &b
}

//TimeArg shortcut for time.Time arguments in NewXXX() function
func TimeArg(t time.Time) *time.Time {
	return &t
}
