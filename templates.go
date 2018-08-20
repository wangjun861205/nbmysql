package nbmysql

//PackageTemp package template
var packageTemp = `package %s`

//ImportTemp import block template
const importTemp = `import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
	"sort"
	"github.com/wangjun861205/nbmysql"
)`

//DbTemp *sql.DB declaration template
const dbTemp = `var %s *sql.DB`

//InitFuncTemp init func template
const initFuncTemp = `func init() {
	db, err := sql.Open("mysql", "%s:%s@tcp(%s)/%s")
	if err != nil {
		panic(err)
	}
	%s = db
	%s
}`

//FieldTemp field in table model template
const fieldTemp = `%s *%s`

//ModelTemp model template
const modelTemp = `type %s struct {
		%s
		_IsStored bool
}`

//FuncArgTemp arguments in function signature template
const funcArgTemp = `%s %s`

//FuncArgNameTemp arguments name in function body template
const funcArgNameTemp = `%s`

//NewModelAsignTemp asign statement in NewXXX() template
const newModelAsignTemp = `%s: %s`

//NewModelFuncTemp NewXXX() template
const newModelFuncTemp = `func New%s(%s) *%s {
		return &%s{%s, 
			_IsStored: false}
	}`

//AllModelFuncTemp AllXXX() template
const allModelFuncTemp = `func All%s() ([]*%s, error) {
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

//QueryModelFuncTemp QueryXXX() template
const queryModelFuncTemp = `func Query%s(query string) ([]*%s, error) {
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

//ForeignKeyAllSQLTemp query all sql statement in foreign key relation template
const foreignKeyAllSQLTemp = `SELECT %s.* FROM %s JOIN %s ON %s.%s=%s.%s where %s.%s = ?`

//ForeignKeyFilterSQLTemp query sql statement in foreign key relation template
const foreignKeyFilterSQLTemp = `SELECT %s.* FROM %s JOIN %s ON %s.%s=%s.%s where %s.%s = ? AND ?`

//InsertSQLTemp insert sql statement
const insertSQLTemp = `INSERT INTO %s (%%s) VALUES (%%s)`

//InsertMiddleTableSQLTemp insert into middle table sql statement template
const insertMiddleTableSQLTemp = `INSERT INTO %s (%s, %s) VALUES (?, ?)`

// const ModelCheckStringBlockTemp = `if %s.%s != nil {
// 		colList = append(colList, "%s")
// 		valList = append(valList, fmt.Sprintf("%%q", *%s.%s))
// 	}`

// const ModelCheckIntBlockTemp = `if %s.%s != nil {
// 		colList = append(colList, "%s")
// 		valList = append(valList, fmt.Sprintf("%%d", *%s.%s))
// 	}`
// const ModelCheckFloatBlockTemp = `if %s.%s != nil {
// 		colList = append(colList, "%s")
// 		valList = append(valList, fmt.Sprintf("%%f", *%s.%s))
// 	}`

// const ModelCheckTimeBlockTemp = `if %s.%s != nil {
// 		colList = append(colList, "%s")
// 		valList = append(valList, fmt.Sprintf("%%q", %s.%s.Format("2006-01-02 15:04:05")))
// 	}`

// const ModelCheckBoolBlockTemp = `if %s.%s != nil {
// 		colList = append(colList, "%s")
// 		valList = append(valList, fmt.Sprintf("%%t", *%s.%s))
// 	}`

// const ModelInsertArgTemp = `m.%s`

//ModelInsertMethodTemp XXX.Insert() method template
const modelInsertMethodTemp = `func (m *%s) Insert() error {
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

// const InsertOrUpdateArgTemp = `m.%s`

//ModelInsertOrUpdateMethodTemp XXX.InsertOrUpdate() method template
const modelInsertOrUpdateMethodTemp = `func (m *%s) InsertOrUpdate() error {
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

// const UpdateArgTemp = `m.%s`

//ModelUpdateMethodTemp XXX.Update() method template
const modelUpdateMethodTemp = `func (m *%s) Update() error {
	if !m._IsStored {
		return nbmysql.ErrModelNotStoredInDB
	}
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

//DeleteArgTemp argument template in Delete() method
const deleteArgTemp = `m.%s`

//ModelDeleteMethodTemp XXX.Delete() method template
const modelDeleteMethodTemp = `func (m *%s) Delete() error {
	if !m._IsStored {
		return nbmysql.ErrModelNotStoredInDB
	}
	_, err := %sDeleteStmt.Exec(%s)
	if err != nil {
		return err
	}
	m._IsStored = false
	return nil
	}`

//NewMiddleTypeTemp craete middle type template
const newMiddleTypeTemp = `_%s := new(nbmysql.%s)`

//ModelFromRowsCheckNullBlockTemp check middle type IsNull in XXXFromRows() function template
const modelFromRowsCheckNullBlockTemp = `if !_%s.IsNull {
		%s = &_%s.Value
	}`

//ModelFromRowsFuncTemp XXXFromRows() function template
const modelFromRowsFuncTemp = `func %sFromRows(rows *sql.Rows) (*%s, error) {
		%s
		err := rows.Scan(%s)
		if err != nil {
			return nil, err
		}
		return &%s{%s, true}, nil
	}`

//ModelFromRowFuncTemp XXXFromRow() function template
const modelFromRowFuncTemp = `func %sFromRow(row *sql.Row) (*%s, error) {
	%s
	err := row.Scan(%s)
	if err != nil {
		return nil, err
	}
	return &%s{%s, true}, nil
}`

//MapElemTemp template of elements of map
const mapElemTemp = `"%s": "%s",`

//QueryFieldMapTemp query field map declaration template
const queryFieldMapTemp = `var %sMap = map[string]string {
	%s
	}`

//QueryByPrimaryKeySQLTemp query by primary sql statement template
const queryByPrimaryKeySQLTemp = `SELECT * FROM %s WHERE %s = ?`

//ModelExistsMethodTemp XXX.Exists() method template
const modelExistsMethodTemp = `func (m *%s) Exists() (bool, error) {
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

//ForeignKeyQuerySQLTemp query sql statement template in foreign key relation
const foreignKeyQuerySQLTemp = `SELECT * FROM %s WHERE %s = ?`

//ForeignKeyMethodTemp foreign key relation method template
const foreignKeyMethodTemp = `func (m *%s) %sBy%s() (*%s, error) {
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

//ReverseForeignKeyStructTypeTemp reverse foreign key struct definition template
const reverseForeignKeyStructTypeTemp = `type %sTo%s struct {
	All func() ([]*%s, error)
	Query func(query string) ([]*%s, error)}`

//ReverseForeignKeyAllSQLTemp sql statement template in reverse foreign key All() method
const reverseForeignKeyAllSQLTemp = `SELECT * FROM %s WHERE %s = ?`

//ReverseForeignKeyAllMethodTemp All() method template in reverse foreign key relation struct
const reverseForeignKeyAllMethodTemp = `func() ([]*%s, error) {
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

//ReverseForeignKeyQuerySQLTemp sql statement template in Query() method of reverse foreign key relation struct
const reverseForeignKeyQuerySQLTemp = `SELECT * FROM %s WHERE %s = ? AND %%s`

//ReverseForeignKeyQueryMethodTemp Query() method template in reverse foreign key relation struct
const reverseForeignKeyQueryMethodTemp = `func(query string) ([]*%s, error) {
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

//ReverseForeignKeyMethodTemp foreign key realtion method template
const reverseForeignKeyMethodTemp = `func (m *%s) %sBy%s() %sTo%s {
	return %sTo%s {
		All: %s,
		Query: %s,
	}
}`

//ManyToManyStructTypeTemp many to many relation struct definition template
const manyToManyStructTypeTemp = `type %sTo%s struct {
		All    func() ([]*%s, error)
		Query func(query string) ([]*%s, error)
		Add func(%s *%s) error
		Remove func(%s *%s) error
	}`

//ManyToManyMethodTemp many to many relation struct declaration template
const manyToManyMethodTemp = `func (m *%s) %sBy%s() %sTo%s {
	return %sTo%s{
		All: %s,
		Query: %s,
		Add: %s,
		Remove: %s,
	}
}`

//ManyToManyAllSQLTemp sql statement template in All() method of many to many relation struct
const manyToManyAllSQLTemp = `SELECT %s.* FROM %s JOIN %s ON %s.%s=%s.%s JOIN %s on %s.%s = %s.%s WHERE %s.%s = ?`

//ManyToManyAllMethodTemp All() method template in many to many relation struct
const manyToManyAllMethodTemp = `func() ([]*%s, error) {
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

//ManyToManyQuerySQLTemp sql statement template in Query() method of many to many relation struct
const manyToManyQuerySQLTemp = `SELECT %s.* FROM %s JOIN %s ON %s.%s=%s.%s JOIN %s on %s.%s = %s.%s WHERE %s.%s = ? AND %%s`

//ManyToManyQueryMethodTemp Query() method template in many to many relation struct
const manyToManyQueryMethodTemp = `func(query string) ([]*%s, error) {
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

//ManyToManyAddMethodTemp Add() method template in many to many relation struct
const manyToManyAddMethodTemp = `func(%s *%s) error {
	if !m._IsStored {
		return errors.New("%s model is not stored in database")
	}
	if !%s._IsStored {
		return errors.New("%s model is not stored in database")
	}
	_, err := %sTo%sInsertStmt.Exec(m.%s, %s.%s)
	return err
}`

//ManyToManyRemoveMethodTemp Remove() method template in many to many relation struct
const manyToManyRemoveMethodTemp = `func(%s *%s) error {
	if !m._IsStored {
		return errors.New("%s model is not stored in database")
	}
	if !%s._IsStored {
		return errors.New("%s model is not stored in database")
	}
	_, err := %sTo%sDeleteStmt.Exec(m.%s, %s.%s)
	return err
}`

//FieldCheckNullTemp nil check block template in Insert() or Update() method
const fieldCheckNullTemp = `if m.%s == nil {
	return errors.New("%s.%s can not be null")
	}`

//ModelCheckMethodTemp XXX.check() method template
const modelCheckMethodTemp = `func (m *%s) check() error {
	%s
	return nil
	}`

//InsertStmtTemp insert sql.Stmt declaration template
const insertStmtTemp = `var %sInsertStmt *sql.Stmt`

//UpdateStmtTemp update sql.Stmt declaration template
const updateStmtTemp = `var %sUpdateStmt *sql.Stmt`

//DeleteStmtTemp delete sql.Stmt declaration template
const deleteStmtTemp = `var %sDeleteStmt *sql.Stmt`

//InsertOrUpdateStmtTemp insert or update sql.Stmt declaration template
const insertOrUpdateStmtTemp = `var %sInsertOrUpdateStmt *sql.Stmt`

//ManyToManyDeleteStmtTemp delete sql.Stmt template for many to many relation
const manyToManyDeleteStmtTemp = `var %sTo%sDeleteStmt *sql.Stmt`

//InsertStmtInitTemp insert sql.Stmt init template
const insertStmtInitTemp = `%sInsertStmt, err = %s.Prepare("INSERT INTO %s (%s) VALUES (%s)")
if err != nil {
	log.Fatal(err)
	}`

//UpdateStmtInitTemp update sql.Stmt init template
const updateStmtInitTemp = `%sUpdateStmt, err = %s.Prepare("UPDATE %s SET %s WHERE %s = ?")
if err != nil {
	log.Fatal(err)
	}`

//DeleteStmtInitTemp delete sql.Stmt init template
const deleteStmtInitTemp = `%sDeleteStmt, err = %s.Prepare("DELETE FROM %s WHERE %s = ?")
if err != nil {
	log.Fatal(err)
	}`

//InsertMidStmtInitTemp insert into middle table sql.Stmt init template
const insertMidStmtInitTemp = `%sInsertStmt, err = %s.Prepare("INSERT INTO %s (%s, %s) VALUES (?, ?)")
if err != nil {
	log.Fatal(err)
	}`

//ManyToManyDeleteStmtInitTemp delete middle table sql.Stmt init template
const manyToManyDeleteStmtInitTemp = `%sTo%sDeleteStmt, err = %s.Prepare("DELETE FROM %s WHERE %s = ? AND %s = ?")`

//UpdateColumnTemp update set clause template
const updateColumnTemp = `%s = ?`

//UpdateLastInsertIDTemp get last insert id template
const updateLastInsertIDTemp = `%s = LAST_INSERT_ID(%s)`

//InsertOrUpdateStmtInitTemp insert or update sql.Stmt init template
const insertOrUpdateStmtInitTemp = `%sInsertOrUpdateStmt, err = %s.Prepare("INSERT INTO %s (%s) VALUES (%s) ON DUPLICATE KEY UPDATE %s, %s")
if err != nil {
	log.Fatal(err)
	}`

//StmtArgTemp append to argument list template
const stmtArgTemp = `argList = append(argList, m.%s)`

//StmtArgNilToDefaultTemp map nil argument to correspond default value template
const stmtArgNilToDefaultTemp = `if m.%s == nil {
	argList = append(argList, %v)
	} else {
	argList = append(argList, m.%s)
	}`

//StmtArgNilToDefaultBlockTemp nil arguments to correspond default value template
const stmtArgNilToDefaultBlockTemp = `argList := make([]interface{}, 0, %d)
%s`

//QueryOneFuncTemp QueryOneXXX() function template
const queryOneFuncTemp = `func QueryOne%s(query string) (*%s, error) {
	for k, v := range %sMap {
		query = strings.Replace(query, k, v, -1)
		}
	row := %s.QueryRow(fmt.Sprintf("SELECT * FROM %s WHERE %%s", query))
	return %sFromRow(row)
	}`

//ModelSetFieldMethodTemp XXX.SetYYY() method template
const modelSetFieldMethodTemp = `func (m *%s) Set%s(val %s, null bool) {
	if null {
		m.%s = nil
		m._IsStored = false
		return
		}
	m.%s = &val
	m._IsStored = false
	}`

//ModelGetFieldMethodTemp XXX.GetYYY() method template
const modelGetFieldMethodTemp = `func (m *%s) Get%s() (%s, bool) {
	if m.%s == nil {
		return %s, true
		}
	return *m.%s, false
	}`

//ModelListTypeTemp XXXList struct definition template
const modelListTypeTemp = `type %sList struct {
	Models []*%s
	Funcs []func(i, j int) int
	}`

//ModelCompareByIntMethodTemp compare by int field method template
const modelCompareByIntMethodTemp = `By%s: func(i, j int) int {
	if *ml.Models[i].%s == *ml.Models[j].%s {
		return 0
	}
	if *ml.Models[i].%s < *ml.Models[j].%s {
		return -1
	}
	return 1
	},`

//ModelCompareByFloatMethodTemp compare by float field method template
const modelCompareByFloatMethodTemp = `By%s: func(i, j int) int {
	if *ml.Models[i].%s == *ml.Models[j].%s {
		return 0
	}
	if *ml.Models[i].%s < *ml.Models[j].%s {
		return -1
	}
	return 1
	},`

//ModelCompareByStringMethodTemp compare by string field method template
const modelCompareByStringMethodTemp = `By%s: func(i, j int) int {
	if *ml.Models[i].%s == *ml.Models[j].%s {
		return 0
	}
	if *ml.Models[i].%s < *ml.Models[j].%s {
		return -1
	}
	return 1
	},`

//ModelCompareByBoolMethodTemp compare by bool field method template
const modelCompareByBoolMethodTemp = `By%s: func(i, j int) int {
	if *ml.Models[i].%s == *ml.Models[j].%s {
		return 0
	}
	if *ml.Models[i].%s == false && *ml.Models[j].%s == true {
		return -1
	}
	return 1
	},`

//ModelCompareByTimeMethodTemp compare by time.Time field method template
const modelCompareByTimeMethodTemp = `By%s: func(i, j int) int {
	if ml.Models[i].%s.Equal(*ml.Models[j].%s) {
		return 0
	}
	if ml.Models[i].%s.Before(*ml.Models[j].%s) {
		return -1
	}
	return 1
	},`

//ModelSortMethodsStructFieldTypeTemp sort struct field definition template
const modelSortMethodsStructFieldTypeTemp = `By%s func(i, j int) int`

//ModelSortMethodsStructTypeTemp sort struct definition template
const modelSortMethodsStructTypeTemp = `type %sSortMethods struct {
	%s
	}`

//ModelSortMethodsStructFuncTemp sort struct method template
const modelSortMethodsStructFuncTemp = `func (ml %sList) SortMethods() %sSortMethods {
	return %sSortMethods{
		%s
		}
	}`

//ModelListLenMethodTemp XXXList.Len() method template
const modelListLenMethodTemp = `func (ml %sList) Len() int {
	return len(ml.Models)
	}`

//ModelListSwapMethodTemp XXXList().Swap() method template
const modelListSwapMethodTemp = `func (ml %sList) Swap(i, j int) {
	ml.Models[i], ml.Models[j] = ml.Models[j], ml.Models[i]
	}`

//ModelListLessMethodTemp XXXList().Less() method template
const modelListLessMethodTemp = `func (ml %sList) Less(i, j int) bool {
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

//ModelSortFuncSwitchBlockTemp sort function switch block template
const modelSortFuncSwitchBlockTemp = `case "%s": 
	%sList.Funcs = append(%sList.Funcs, sortMethods.By%s)`

//ModelSortFuncTemp XXXSortBy() function template
const modelSortFuncTemp = `func %sSortBy(ml []*%s, desc bool, fields ...string) {
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
	if desc {
		sort.Sort(sort.Reverse(%sList))
	} else {
		sort.Sort(%sList)
	}
}`

const modelCountStmtDeclareTemp = `var %sCountStmt *sql.Stmt`

// const modelCountStmtInitTemp = `%sCountStmt, err = %s.Prepare("SELECT COUNT(*) FROM %s")
// if err != nil {
// 	log.Fatal(err)
// 	}`

const modelCountStmtInitTemp = `{{Table>ModelName}}CountStmt, err = {{Database>ObjName}}.Prepare("SELECT COUNT(*) FROM {{Table>TableName}}")
if err != nil {
	log.Fatal(err)
	}`

// const modelCountFuncTemp = `func %sCount() (int64, error) {
// 	var count int64
// 	err := %sCountStmt.QueryRow().Scan(&count)
// 	if err != nil {
// 		return -1, err
// 	}
// 	return count, nil
// }`

const modelCountFuncTemp = `func {{Table>ModelName}}Count() (int64, error) {
	var count int64
	err := {{Table>ModelName}}CountStmt.QueryRow().Scan(&count)
	if err != nil {
		return -1, err
	}
	return count, nil
}`
