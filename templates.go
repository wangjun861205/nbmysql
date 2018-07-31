package nbmysql

var PackageTemp = `package %s`

const ImportTemp = `import (
	"database/sql"
	"fmt"
	"strings"
	"time"
	"github.com/wangjun861205/nbmysql"
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
		Insert func(%s *%s) error
	}`

const FuncArgTemp = `%s%s *%s`
const FuncArgNameTemp = `%s%s`

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
		rows, err := %s.Query(fmt.Sprintf("SELECT * FROM %s WHERE %%s", query))
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

const InsertSQLTemp = `INSERT INTO %s (%%s) VALUES (%%s)`
const InsertMiddleTableSQLTemp = `INSERT INTO %s (%s, %s) VALUES (?, ?)`

const ManyToManyInsertTemp = `Insert: func(%s *%s) error {
				tx, err := %s.Begin()
				if err != nil {
					return err
				}
				colList := make([]string, 0, 32)
				valList := make([]string, 0, 32)
				%s
				res, err := tx.Exec(fmt.Sprintf("%s", strings.Join(colList, ", "), strings.Join(valList, ", ")))
				if err != nil {
					tx.Rollback()
					return err
				}
				lastInsertId, err := res.LastInsertId()
				if err != nil {
					tx.Rollback()
					return err
				}
				%s.%s = &lastInsertId
				_, err = tx.Exec("%s", *m.%s, *%s.%s)
				if err != nil {
					tx.Rollback()
					return err
				}
				return tx.Commit()
			},`

const ForeignKeyInsertTemp = `Insert: func(%s *%s) error {
				tx, err := %s.Begin()
				if err != nil {
					return err
				}
				colList := make([]string, 0, 32)
				valList := make([]string, 0, 32)
				%s
				res, err := tx.Exec("%s", strings.Join(colList, ", "), strings.Join(valList, ", "))
				if err != nil {
					tx.Rollback()
					return err
				}
				lastInsertId, err := res.LastInsertId()
				if err != nil {
					tx.Rollback()
					return err
				}
				%s.%s = &lastInsertId
				return tx.Commit()
			},`

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
			%s
		}
	}`

const ModelCheckStringBlockTemp = `if %s.%s != nil {
		colList = append(colList, "%s")
		valList = append(valList, fmt.Sprintf("%%q", *%s.%s))
	}`

const ModelCheckIntBlockTemp = `if %s.%s != nil {
		colList = append(colList, "%s")
		valList = append(valList, fmt.Sprintf("%%d", *%s.%s))
	}`
const ModelCheckFloatBlockTemp = `if %s.%s != nil {
		colList = append(colList, "%s")
		valList = append(valList, fmt.Sprintf("%%f", *%s.%s))
	}`

const ModelCheckTimeBlockTemp = `if %s.%s != nil {
		colList = append(colList, "%s")
		valList = append(valList, fmt.Sprintf("%%q", %s.%s.Format("2006-01-02 15:04:05")))
	}`

const ModelCheckBoolBlockTemp = `if %s.%s != nil {
		colList = append(colList, "%s")
		valList = append(valList, fmt.Sprintf("%%t", *%s.%s)
	}`

const ModelInsertMethodTemp = `func (m *%s) Insert() error {
		colList := make([]string, 0, 32)
		valList := make([]string, 0, 32)
		%s
		res, err := %s.Exec(fmt.Sprintf("%s", strings.Join(colList, ", "), strings.Join(valList, ", ")))
		if err != nil {
			return err
		}
		lastInsertId, err := res.LastInsertId()
		if err != nil {
			return err
		}
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
		_, err := %s.Exec(fmt.Sprintf("UPDATE %s SET %%s WHERE %s = ?", strings.Join(updateList, ", ")), *m.%s)
		return err
	}`

const DeleteSQLTemp = `DELETE FROM %s WHERE %s = ?`

const ManyToManyDeleteBlockTemp = `_, err = tx.Exec("%s", *m.%s)
	if err != nil {
		tx.Rollback()
		return err
		}`

const ModelDeleteMethodTemp = `func (m *%s) Delete() error {
		tx, err := %s.Begin()
		if err != nil {
			return err
		}
		%s
		_, err = tx.Exec("%s", *m.%s)
		if err != nil {
			tx.Rollback()
			return err
		}
		return tx.Commit()
	}`

const NewMiddleTypeTemp = `_%s := new(nbmysql.%s)`

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
