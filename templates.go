package nbmysql

var PackageTemp = `package %s`

const ImportTemp = `import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)`

const DbTemp = `var %s *sql.DB`

const InitFuncTemp = `func init() {
	db, err := sql.Open("mysql", "%s:%s@tcp(%s)/%s")
	if err != nil {
		panic(err)
	}
	%s = db
}`

const FieldTemp = `%s *%s`

const ModelTemp = `type %s struct {
		%s
}`

const ModelRelationTemp = `type %sTo%s struct {
		All    func() ([]*%s, error)
		Filter func(query string) ([]*%s, error)
	}`

const FuncArgTemp = `%s%s *%s`
const FuncArgNameTemp = `%s, `

const NewModelFuncTemp = `func New%s(%s) *%s {
		%s := &%s{%s}
		return %s
	}`

const AllModelFuncTemp = `func All%s() ([]*%s, error) {
		rows, err := %s.Query("SELECT * FROM %s")
		if err != nil {
			return nil, err
		}
		list := make([]*%s, 0, 256)
		for rows.Next() {
			model, err := %sFromRows(rows)
			if err != nil {
				return nil, err
			}
			list = append(list, model)
		}
		return list, nil
	}`

const QueryModelFuncTemp = `func Query%s(query string) ([]*%s, error) {
		for k, v := range %sMap {
			query = strings.Replace(query, k, v, -1)
		}
		rows, err := %s.Query("SELECT * FROM %s WHERE ?", query)
		if err != nil {
			return nil, err
		}
		list := make([]*%s, 0, 256)
		for rows.Next() {
			model, err := %sFromRows(rows)
			if err != nil {
				return nil, err
			}
			list = append(list, model)
		}
		return list, nil
	}`

const ManyToManyAllSQLTemp = `SELECT %s.* FROM %s JOIN %s ON %s.%s=%s.%s JOIN %s on %s.%s = %s.%s WHERE %s.%s = ?`
const ManyToManyFilterSQLTemp = `SELECT %s.* FROM %s JOIN %s ON %s.%s=%s.%s JOIN %s on %s.%s = %s.%s WHERE %s.%s = ? AND ?`

const ForeignKeyAllSQLTemp = `SELECT %s.* FROM %s JOIN %s ON %s.%s=%s.%s where %s.%s = ?`
const ForeignKeyFilterSQLTemp = `SELECT %s.* FROM %s JOIN %s ON %s.%s=%s.%s where %s.%s = ? AND ?`

const ModelRelationFuncTemp = `func (m *%s) %sBy%s() %sTo%s {
		return %sTo%s{
			All: func() ([]*%s, error) {
				rows, err := %s.Query("%s", *m.%s)
				if err != nil {
					return nil, err
				}
				list := make([]*%s, 0, 256)
				for rows.Next() {
					model, err := %sFromRows(rows)
					if err != nil {
						return nil, err
					}
					list = append(list, model)
				}
				return list, nil
			},
			Filter: func(query string) ([]*%s, error) {
				for k, v := range %sMap {
					query = strings.Replace(query, k, v, -1)
				}
				rows, err := %s.Query("%s", *m.%s, query)
				if err != nil {
					return nil, err
				}
				list := make([]*%s, 0, 256)
				for rows.Next() {
					model, err := %sFromRows(rows)
					if err != nil {
						return nil, err
					}
					list = append(list, model)
				}
				return list, nil
			},
		}
	}`

const ModelCheckStringBlockTemp = `if m.%s != nil {
		colList = append(colList, "%s")
		valList = append(valList, fmt.Sprintf("%%q", *m.%s))
	}`

const ModelCheckIntBlockTemp = `if m.%s != nil {
		colList = append(colList, "%s")
		valList = append(valList, fmt.Sprintf("%%d", *m.%s))
	}`
const ModelCheckFloatBlockTemp = `if m.%s != nil {
		colList = append(colList, "%s")
		valList = append(valList, fmt.Sprintf("%%f", *m.%s))
	}`

const ModelCheckTimeBlockTemp = `if m.%s != nil {
		colList = append(colList, "%s")
		valList = append(valList, fmt.Sprintf("%%q", m.%s.Format("2006-01-02 15:04:05")))
	}`

const ModelCheckBoolBlockTemp = `if m.%s != nil {
		colList = append(colList, "%s")
		valList = append(valList, fmt.Sprintf("%%t", *m.%s)
	}`

const ModelInsertMethodTemp = `func (m *%s) Insert() error {
		colList := make([]string, 0, 32)
		valList := make([]string, 0, 32)
		%s
		_, err := %s.Exec("INSERT INTO %s (?) VALUES (?)", strings.Join(colList, ", "), strings.Join(valList, ", "))
		if err != nil {
			return nil
		}
		lastInsertId := GetLastId(%s)
		m.%s = &lastInsertId
		return nil
}`

const ModelUpdateMethodTemp = `func (m *%s) Update() error {
		colList := make([]string, 0, 32)
		valList := make([]string, 0, 32)
		%s
		updateList := make([]string, 0, 32)
		for i := 0; i < len(colList); i++ {
			updateList = append(updateList, fmt.Sprintf("%%s=%%s", colList[i], valList[i]))
		}
		_, err := %s.Exec("UPDATE %s SET ? WHERE id = ?", strings.Join(updateList, ", "), *m.%s)
		return err
	}`

const ModelDeleteMethodTemp = `func (m *%s) Delete() error {
		_, err := %s.Exec("DELETE FROM %s where 'id' = ?", *m.%s)
		return err
	}`

const NewMiddleTypeTemp = `_%s := new(%s)`

const ModelFromRowsCheckNullBlockTemp = `if !_%s.IsNull {
		%s = &_%s.Value
	}`

const ModelFromRowsFuncTemp = `func %sFromRows(rows *sql.Rows) (*%s, error) {
		%s
		err := rows.Scan(%s)
		if err != nil {
			return nil, err
		}
		var (
			%s
		)
		%s
		return New%s(%s), nil
	}`

const MapElemTemp = `"%s": "%s",`

const QueryFieldMapTemp = `var %sMap = map[string]string {
	%s
	}`
