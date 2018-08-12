package nbmysql

var PackageTemp = `package %s`

const ImportTemp = `import (
	"database/sql"
	"errors"
	"fmt"
	"log"
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
	%s
}`

const FieldTemp = `%s *%s`

const ModelTemp = `type %s struct {
		%s
		_IsStored bool
}`

const FuncArgTemp = `%s %s`
const FuncArgNameTemp = `%s`

const CheckIntArgTemp = `if imag(%s) > 0 {
	 *%s = int64(real(%s)) 
	 } else { 
	 %s = nil 
	 }`

const CheckFloatArgTemp = `if imag(%s) > 0 {
	*%s = float64(real(%s))
	} else {
		%s = nil
	}`

const CheckStringArgTemp = `if %s != nil {
	*%s = string(%s)
} else {
	%s = nil
}`

const CheckBoolArgTemp = `if %s == 0 {
	*%s = false 
	} else if %s > 0 {
	*%s = true
	} else {
		%s = nil
	}`

const CheckTimeArgTemp = `%s = %s`

const NewModelAsignTemp = `%s: %s`

const NewModelFuncTemp = `func New%s(%s) *%s {
		return &%s{%s, 
			_IsStored: false}
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
			model._IsStored = true
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
			model._IsStored = true
			list = append(list, model)
		}
		return list, nil
	}`

const ForeignKeyAllSQLTemp = `SELECT %s.* FROM %s JOIN %s ON %s.%s=%s.%s where %s.%s = ?`
const ForeignKeyFilterSQLTemp = `SELECT %s.* FROM %s JOIN %s ON %s.%s=%s.%s where %s.%s = ? AND ?`

const InsertSQLTemp = `INSERT INTO %s (%%s) VALUES (%%s)`
const InsertMiddleTableSQLTemp = `INSERT INTO %s (%s, %s) VALUES (?, ?)`

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
		valList = append(valList, fmt.Sprintf("%%t", *%s.%s))
	}`

// const ModelInsertMethodTemp = `func (m *%s) Insert() error {
// 		colList := make([]string, 0, 32)
// 		valList := make([]string, 0, 32)
// 		%s
// 		res, err := %s.Exec(fmt.Sprintf("%s", strings.Join(colList, ", "), strings.Join(valList, ", ")))
// 		if err != nil {
// 			if sqlErr, ok := err.(*mysql.MySQLError); ok && (sqlErr.Number == 1022 || sqlErr.Number == 1062){
// 				m._IsStored = true
// 				return nbmysql.ErrDupKey
// 			}
// 			return err
// 		}
// 		lastInsertId, err := res.LastInsertId()
// 		if err != nil {
// 			return err
// 		}
// 		m.%s = &lastInsertId
// 		m._IsStored = true
// 		return nil
// }`

const ModelInsertArgTemp = `m.%s`

const ModelInsertMethodTemp = `func (m *%s) Insert() error {
		err := m.check()
		if err != nil {
			return err
		}
		%s
		res, err := %sInsertStmt.Exec(argList...)
		if err != nil {
			return err
		}
		lastInsertId, err := res.LastInsertId()
		if err != nil {
			return err
		}
		m.%s = &lastInsertId
		m._IsStored = true
		return nil
}`

const InsertOrUpdateArgTemp = `m.%s`
const ModelInsertOrUpdateMethodTemp = `func (m *%s) InsertOrUpdate() error {
	err := m.check()
	if err != nil {
		return err
	}
	%s
	argList = append(argList, argList...)
	res, err := %sInsertOrUpdateStmt.Exec(argList...)
	if err != nil {
		return err
	}
	lastInsertId, err := res.LastInsertId()
	if err != nil {
		return err
	}
	m.%s = &lastInsertId
	m._IsStored = true
	return nil
}`

const UpdateArgTemp = `m.%s`
const ModelUpdateMethodTemp = `func (m *%s) Update() error {
	err := m.check()
	if err != nil {
		return err
	}
	%s
	argList = append(argList, m.%s)
	_, err = %sUpdateStmt.Exec(argList...)
	if err != nil {
		return err
	}
	return nil
}`

const DeleteArgTemp = `m.%s`
const ModelDeleteMethodTemp = `func (m *%s) Delete() error {
	_, err := %sDeleteStmt.Exec(%s)
	if err != nil {
		return err
	}
	m._IsStored = false
	return nil
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
		return &%s{%s, true}, nil
	}`

const ModelFromRowFuncTemp = `func %sFromRow(row *sql.Row) (*%s, error) {
	%s
	err := row.Scan(%s)
	if err != nil {
		return nil, err
	}
	return &%s{%s, true}, nil
}`

const MapElemTemp = `"%s": "%s",`

const QueryFieldMapTemp = `var %sMap = map[string]string {
	%s
	}`

const QueryByPrimaryKeySQLTemp = `SELECT * FROM %s WHERE %s = ?`
const ModelExistsMethodTemp = `func (m *%s) Exists() (bool, error) {
	if m.%s == nil {
		return false, errors.New("%s.%s must not be nil")
	}
	row := %s.QueryRow("%s", m.%s)
	if row == nil {
		return false, nil
	}
	m._IsStored = true
	return true, nil
}`

const ForeignKeyQuerySQLTemp = `SELECT * FROM %s WHERE %s = ?`
const ForeignKeyMethodTemp = `func (m *%s) %sBy%s() (*%s, error) {
	if m.%s == nil {
		return nil, errors.New("%s.%s must not be nil")
	}
	row := %s.QueryRow("%s", m.%s)
	if row == nil {
		return nil, nbmysql.ErrRecordNotExists
	}
	model, err := %sFromRow(row)
	if err != nil {
		return nil, err
	}
	model._IsStored = true
	return model, nil
}`

const ReverseForeignKeyStructTypeTemp = `type %sTo%s struct {
	All func() ([]*%s, error)
	Query func(query string) ([]*%s, error)}`

const ReverseForeignKeyAllSQLTemp = `SELECT * FROM %s WHERE %s = ?`
const ReverseForeignKeyAllMethodTemp = `func() ([]*%s, error) {
	if m.%s == nil {
		return nil, errors.New("%s.%s must not be nil")
	}
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
		model._IsStored = true
		list = append(list, model)
	}
	return list, nil
}`

const ReverseForeignKeyQuerySQLTemp = `SELECT * FROM %s WHERE %s = ? AND %%s`
const ReverseForeignKeyQueryMethodTemp = `func(query string) ([]*%s, error) {
	if m.%s == nil {
		return nil, errors.New("%s.%s must not be nil")
	}
	for k, v := range %sMap {
		query = strings.Replace(query, k, v, -1)
	}
	rows, err := %s.Query(fmt.Sprintf("%s", query), *m.%s)
	if err != nil {
		return nil, err
	}
	list := make([]*%s, 0, 256)
	for rows.Next() {
		model, err := %sFromRows(rows)
		if err != nil {
			return nil, err
		}
		model._IsStored = true
		list = append(list, model)
	}
	return list, nil
}`

const ReverseForeignKeyMethodTemp = `func (m *%s) %sBy%s() %sTo%s {
	return %sTo%s {
		All: %s,
		Query: %s,
	}
}`

const ManyToManyStructTypeTemp = `type %sTo%s struct {
		All    func() ([]*%s, error)
		Query func(query string) ([]*%s, error)
		Add func(%s *%s) error
		Remove func(%s *%s) error
	}`

const ManyToManyMethodTemp = `func (m *%s) %sBy%s() %sTo%s {
	return %sTo%s{
		All: %s,
		Query: %s,
		Add: %s,
		Remove: %s,
	}
}`

const ManyToManyAllSQLTemp = `SELECT %s.* FROM %s JOIN %s ON %s.%s=%s.%s JOIN %s on %s.%s = %s.%s WHERE %s.%s = ?`
const ManyToManyAllMethodTemp = `func() ([]*%s, error) {
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
		model._IsStored = true
		list = append(list, model)
	}
	return list, nil
}`

const ManyToManyQuerySQLTemp = `SELECT %s.* FROM %s JOIN %s ON %s.%s=%s.%s JOIN %s on %s.%s = %s.%s WHERE %s.%s = ? AND %%s`
const ManyToManyQueryMethodTemp = `func(query string) ([]*%s, error) {
	for k, v := range %sMap {
		query = strings.Replace(query, k, v, -1)
	}
	rows, err := %s.Query(fmt.Sprintf("%s", query), *m.%s)
	if err != nil {
		return nil, err
	}
	list := make([]*%s, 0, 256)
	for rows.Next() {
		model, err := %sFromRows(rows)
		if err != nil {
			return nil, err
		}
		model._IsStored = true
		list = append(list, model)
	}
	return list, nil
}`

// const ManyToManyAddSQLTemp = `INSERT INTO %s (%s, %s) VALUES (?, ?)`
const ManyToManyAddMethodTemp = `func(%s *%s) error {
	if !m._IsStored {
		return errors.New("%s model is not stored in database")
	}
	if !%s._IsStored {
		return errors.New("%s model is not stored in database")
	}
	_, err := %sTo%sInsertStmt.Exec(m.%s, %s.%s)
	return err
}`

// const ManyToManyRemoveSQLTemp = `DELETE FROM %s WHERE %s = ? and %s = ?`
const ManyToManyRemoveMethodTemp = `func(%s *%s) error {
	if !m._IsStored {
		return errors.New("%s model is not stored in database")
	}
	if !%s._IsStored {
		return errors.New("%s model is not stored in database")
	}
	_, err := %sTo%sDeleteStmt.Exec(m.%s, %s.%s)
	return err
}`

const FieldCheckNullTemp = `if m.%s == nil {
	return errors.New("%s.%s can not be null")
	}`

const ModelCheckMethodTemp = `func (m *%s) check() error {
	%s
	return nil
	}`

const InsertStmtTemp = `var %sInsertStmt *sql.Stmt`
const UpdateStmtTemp = `var %sUpdateStmt *sql.Stmt`
const DeleteStmtTemp = `var %sDeleteStmt *sql.Stmt`
const InsertOrUpdateStmtTemp = `var %sInsertOrUpdateStmt *sql.Stmt`
const ManyToManyDeleteStmtTemp = `var %sTo%sDeleteStmt *sql.Stmt`

const InsertStmtInitTemp = `%sInsertStmt, err = %s.Prepare("INSERT INTO %s (%s) VALUES (%s)")
if err != nil {
	log.Fatal(err)
	}`
const UpdateStmtInitTemp = `%sUpdateStmt, err = %s.Prepare("UPDATE %s SET %s WHERE %s = ?")
if err != nil {
	log.Fatal(err)
	}`
const DeleteStmtInitTemp = `%sDeleteStmt, err = %s.Prepare("DELETE FROM %s WHERE %s = ?")
if err != nil {
	log.Fatal(err)
	}`
const InsertMidStmtInitTemp = `%sInsertStmt, err = %s.Prepare("INSERT INTO %s (%s, %s) VALUES (?, ?)")
if err != nil {
	log.Fatal(err)
	}`
const ManyToManyDeleteStmtInitTemp = `%sTo%sDeleteStmt, err = %s.Prepare("DELETE FROM %s WHERE %s = ? AND %s = ?")`

const UpdateColumnTemp = `%s = ?`
const UpdateLastInsertIDTemp = `%s = LAST_INSERT_ID(%s)`
const InsertOrUpdateStmtInitTemp = `%sInsertOrUpdateStmt, err = %s.Prepare("INSERT INTO %s (%s) VALUES (%s) ON DUPLICATE KEY UPDATE %s, %s")
if err != nil {
	log.Fatal(err)
	}`

const StmtArgTemp = `argList = append(argList, m.%s)`
const StmtArgNilToDefaultTemp = `if m.%s == nil {
	argList = append(argList, %v)
	} else {
	argList = append(argList, m.%s)
	}`
const StmtArgNilToDefaultBlockTemp = `argList := make([]interface{}, 0, %d)
%s`

const QueryOneFuncTemp = `func QueryOne%s(query string) (*%s, error) {
	for k, v := range %sMap {
		query = strings.Replace(query, k, v, -1)
		}
	row := %s.QueryRow(fmt.Sprintf("SELECT * FROM %s WHERE %%s", query))
	return %sFromRow(row)
	}`

const ModelSetFieldMethodTemp = `func (m *%s) Set%s(val %s, null bool) {
	if null {
		m.%s = nil
		m._IsStored = false
		return
		}
	m.%s = &val
	m._IsStored = false
	}`

const ModelGetFieldMethodTemp = `func (m *%s) Get%s() (%s, bool) {
	if m.%s == nil {
		return %s, true
		}
	return *m.%s, false
	}`

const ModelListTypeTemp = `type %sList struct {
	Models []*%s
	Funcs []func(i, j int) int
	}`
const ModelCompareByIntMethodTemp = `By%s: func(i, j int) int {
	if *ml.Models[i].%s == *ml.Models[j].%s {
		return 0
	}
	if *ml.Models[i].%s < *ml.Models[j].%s {
		return -1
	}
	return 1
	},`
const ModelCompareByFloatMethodTemp = `By%s: func(i, j int) int {
	if *ml.Models[i].%s == *ml.Models[j].%s {
		return 0
	}
	if *ml.Models[i].%s < *ml.Models[j].%s {
		return -1
	}
	return 1
	},`
const ModelCompareByStringMethodTemp = `By%s: func(i, j int) int {
	if *ml.Models[i].%s == *ml.Models[j].%s {
		return 0
	}
	if *ml.Models[i].%s < *ml.Models[j].%s {
		return -1
	}
	return 1
	},`
const ModelCompareByBoolMethodTemp = `By%s: func(i, j int) int {
	if *ml.Models[i].%s == *ml.Models[j].%s {
		return 0
	}
	if *ml.Models[i].%s == false && *ml.Models[j].%s == true {
		return -1
	}
	return 1
	},`
const ModelCompareByTimeMethodTemp = `By%s: func(i, j int) int {
	if ml.Models[i].%s.Equal(*ml.Models[j].%s) {
		return 0
	}
	if ml.Models[i].%s.Before(*ml.Models[j].%s) {
		return -1
	}
	return 1
	},`
const ModelSortMethodsStructFieldTypeTemp = `By%s func(i, j int) int`
const ModelSortMethodsStructTypeTemp = `type %sSortMethods struct {
	%s
	}`
const ModelSortMethodsStructFuncTemp = `func (ml %sList) SortMethods() %sSortMethods {
	return %sSortMethods{
		%s
		}
	}`
const ModelListLenMethodTemp = `func (ml %sList) Len() int {
	return len(ml.Models)
	}`
const ModelListSwapMethodTemp = `func (ml %sList) Swap(i, j int) {
	ml.Models[i], ml.Models[j] = ml.Models[j], ml.Models[i]
	}`
const ModelListLessMethodTemp = `func (ml %sList) Less(i, j int) bool {
	var less bool
	for _, f := range ml.Funcs {
		res := f(i, j)
		if res == -1 {
			less = true
			break
		} else if res == 1 {
			break
		}
		continue
	}
	return less
}`

const ModelSortFuncSwitchBlockTemp = `case "%s": 
	%sList.Funcs = append(%sList.Funcs, sortMethods.By%s)`
const ModelSortFuncTemp = `func %sSortBy(ml []*%s, fields ...string) {
	%sList := %sList {
		Models: ml,
		Funcs: make([]func(i, j int) int, 0, %d),
		}
	sortMethods := %sList.SortMethods()
	for _, field := range fields {
		switch field {
			%s
		}
	}
	sort.Sort(%sList)
}`
