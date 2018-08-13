package nbmysql

//PackageTemp package template
var PackageTemp = `package %s`

//ImportTemp import block template
const ImportTemp = `import (
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
const DbTemp = `var %s *sql.DB`

//InitFuncTemp init func template
const InitFuncTemp = `func init() {
	db, err := sql.Open("mysql", "%s:%s@tcp(%s)/%s")
	if err != nil {
		panic(err)
	}
	%s = db
	%s
}`

//FieldTemp field in table model template
const FieldTemp = `%s *%s`

//ModelTemp model template
const ModelTemp = `type %s struct {
		%s
		_IsStored bool
}`

//FuncArgTemp arguments in function signature template
const FuncArgTemp = `%s %s`

//FuncArgNameTemp arguments name in function body template
const FuncArgNameTemp = `%s`

// const CheckIntArgTemp = `if imag(%s) > 0 {
// 	 *%s = int64(real(%s))
// 	 } else {
// 	 %s = nil
// 	 }`

// const CheckFloatArgTemp = `if imag(%s) > 0 {
// 	*%s = float64(real(%s))
// 	} else {
// 		%s = nil
// 	}`

// const CheckStringArgTemp = `if %s != nil {
// 	*%s = string(%s)
// } else {
// 	%s = nil
// }`

// const CheckBoolArgTemp = `if %s == 0 {
// 	*%s = false
// 	} else if %s > 0 {
// 	*%s = true
// 	} else {
// 		%s = nil
// 	}`

// const CheckTimeArgTemp = `%s = %s`

//NewModelAsignTemp asign statement in NewXXX() template
const NewModelAsignTemp = `%s: %s`

//NewModelFuncTemp NewXXX() template
const NewModelFuncTemp = `func New%s(%s) *%s {
		return &%s{%s, 
			_IsStored: false}
	}`

//AllModelFuncTemp AllXXX() template
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

//QueryModelFuncTemp QueryXXX() template
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

//ForeignKeyAllSQLTemp query all sql statement in foreign key relation template
const ForeignKeyAllSQLTemp = `SELECT %s.* FROM %s JOIN %s ON %s.%s=%s.%s where %s.%s = ?`

//ForeignKeyFilterSQLTemp query sql statement in foreign key relation template
const ForeignKeyFilterSQLTemp = `SELECT %s.* FROM %s JOIN %s ON %s.%s=%s.%s where %s.%s = ? AND ?`

//InsertSQLTemp insert sql statement
const InsertSQLTemp = `INSERT INTO %s (%%s) VALUES (%%s)`

//InsertMiddleTableSQLTemp insert into middle table sql statement template
const InsertMiddleTableSQLTemp = `INSERT INTO %s (%s, %s) VALUES (?, ?)`

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

// const InsertOrUpdateArgTemp = `m.%s`

//ModelInsertOrUpdateMethodTemp XXX.InsertOrUpdate() method template
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

// const UpdateArgTemp = `m.%s`

//ModelUpdateMethodTemp XXX.Update() method template
const ModelUpdateMethodTemp = `func (m *%s) Update() error {
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
const DeleteArgTemp = `m.%s`

//ModelDeleteMethodTemp XXX.Delete() method template
const ModelDeleteMethodTemp = `func (m *%s) Delete() error {
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
const NewMiddleTypeTemp = `_%s := new(nbmysql.%s)`

//ModelFromRowsCheckNullBlockTemp check middle type IsNull in XXXFromRows() function template
const ModelFromRowsCheckNullBlockTemp = `if !_%s.IsNull {
		%s = &_%s.Value
	}`

//ModelFromRowsFuncTemp XXXFromRows() function template
const ModelFromRowsFuncTemp = `func %sFromRows(rows *sql.Rows) (*%s, error) {
		%s
		err := rows.Scan(%s)
		if err != nil {
			return nil, err
		}
		return &%s{%s, true}, nil
	}`

//ModelFromRowFuncTemp XXXFromRow() function template
const ModelFromRowFuncTemp = `func %sFromRow(row *sql.Row) (*%s, error) {
	%s
	err := row.Scan(%s)
	if err != nil {
		return nil, err
	}
	return &%s{%s, true}, nil
}`

//MapElemTemp template of elements of map
const MapElemTemp = `"%s": "%s",`

//QueryFieldMapTemp query field map declaration template
const QueryFieldMapTemp = `var %sMap = map[string]string {
	%s
	}`

//QueryByPrimaryKeySQLTemp query by primary sql statement template
const QueryByPrimaryKeySQLTemp = `SELECT * FROM %s WHERE %s = ?`

//ModelExistsMethodTemp XXX.Exists() method template
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

//ForeignKeyQuerySQLTemp query sql statement template in foreign key relation
const ForeignKeyQuerySQLTemp = `SELECT * FROM %s WHERE %s = ?`

//ForeignKeyMethodTemp foreign key relation method template
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

//ReverseForeignKeyStructTypeTemp reverse foreign key struct definition template
const ReverseForeignKeyStructTypeTemp = `type %sTo%s struct {
	All func() ([]*%s, error)
	Query func(query string) ([]*%s, error)}`

//ReverseForeignKeyAllSQLTemp sql statement template in reverse foreign key All() method
const ReverseForeignKeyAllSQLTemp = `SELECT * FROM %s WHERE %s = ?`

//ReverseForeignKeyAllMethodTemp All() method template in reverse foreign key relation struct
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

//ReverseForeignKeyQuerySQLTemp sql statement template in Query() method of reverse foreign key relation struct
const ReverseForeignKeyQuerySQLTemp = `SELECT * FROM %s WHERE %s = ? AND %%s`

//ReverseForeignKeyQueryMethodTemp Query() method template in reverse foreign key relation struct
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

//ReverseForeignKeyMethodTemp foreign key realtion method template
const ReverseForeignKeyMethodTemp = `func (m *%s) %sBy%s() %sTo%s {
	return %sTo%s {
		All: %s,
		Query: %s,
	}
}`

//ManyToManyStructTypeTemp many to many relation struct definition template
const ManyToManyStructTypeTemp = `type %sTo%s struct {
		All    func() ([]*%s, error)
		Query func(query string) ([]*%s, error)
		Add func(%s *%s) error
		Remove func(%s *%s) error
	}`

//ManyToManyMethodTemp many to many relation struct declaration template
const ManyToManyMethodTemp = `func (m *%s) %sBy%s() %sTo%s {
	return %sTo%s{
		All: %s,
		Query: %s,
		Add: %s,
		Remove: %s,
	}
}`

//ManyToManyAllSQLTemp sql statement template in All() method of many to many relation struct
const ManyToManyAllSQLTemp = `SELECT %s.* FROM %s JOIN %s ON %s.%s=%s.%s JOIN %s on %s.%s = %s.%s WHERE %s.%s = ?`

//ManyToManyAllMethodTemp All() method template in many to many relation struct
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

//ManyToManyQuerySQLTemp sql statement template in Query() method of many to many relation struct
const ManyToManyQuerySQLTemp = `SELECT %s.* FROM %s JOIN %s ON %s.%s=%s.%s JOIN %s on %s.%s = %s.%s WHERE %s.%s = ? AND %%s`

//ManyToManyQueryMethodTemp Query() method template in many to many relation struct
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

//ManyToManyAddMethodTemp Add() method template in many to many relation struct
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

//ManyToManyRemoveMethodTemp Remove() method template in many to many relation struct
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

//FieldCheckNullTemp nil check block template in Insert() or Update() method
const FieldCheckNullTemp = `if m.%s == nil {
	return errors.New("%s.%s can not be null")
	}`

//ModelCheckMethodTemp XXX.check() method template
const ModelCheckMethodTemp = `func (m *%s) check() error {
	%s
	return nil
	}`

//InsertStmtTemp insert sql.Stmt declaration template
const InsertStmtTemp = `var %sInsertStmt *sql.Stmt`

//UpdateStmtTemp update sql.Stmt declaration template
const UpdateStmtTemp = `var %sUpdateStmt *sql.Stmt`

//DeleteStmtTemp delete sql.Stmt declaration template
const DeleteStmtTemp = `var %sDeleteStmt *sql.Stmt`

//InsertOrUpdateStmtTemp insert or update sql.Stmt declaration template
const InsertOrUpdateStmtTemp = `var %sInsertOrUpdateStmt *sql.Stmt`

//ManyToManyDeleteStmtTemp delete sql.Stmt template for many to many relation
const ManyToManyDeleteStmtTemp = `var %sTo%sDeleteStmt *sql.Stmt`

//InsertStmtInitTemp insert sql.Stmt init template
const InsertStmtInitTemp = `%sInsertStmt, err = %s.Prepare("INSERT INTO %s (%s) VALUES (%s)")
if err != nil {
	log.Fatal(err)
	}`

//UpdateStmtInitTemp update sql.Stmt init template
const UpdateStmtInitTemp = `%sUpdateStmt, err = %s.Prepare("UPDATE %s SET %s WHERE %s = ?")
if err != nil {
	log.Fatal(err)
	}`

//DeleteStmtInitTemp delete sql.Stmt init template
const DeleteStmtInitTemp = `%sDeleteStmt, err = %s.Prepare("DELETE FROM %s WHERE %s = ?")
if err != nil {
	log.Fatal(err)
	}`

//InsertMidStmtInitTemp insert into middle table sql.Stmt init template
const InsertMidStmtInitTemp = `%sInsertStmt, err = %s.Prepare("INSERT INTO %s (%s, %s) VALUES (?, ?)")
if err != nil {
	log.Fatal(err)
	}`

//ManyToManyDeleteStmtInitTemp delete middle table sql.Stmt init template
const ManyToManyDeleteStmtInitTemp = `%sTo%sDeleteStmt, err = %s.Prepare("DELETE FROM %s WHERE %s = ? AND %s = ?")`

//UpdateColumnTemp update set clause template
const UpdateColumnTemp = `%s = ?`

//UpdateLastInsertIDTemp get last insert id template
const UpdateLastInsertIDTemp = `%s = LAST_INSERT_ID(%s)`

//InsertOrUpdateStmtInitTemp insert or update sql.Stmt init template
const InsertOrUpdateStmtInitTemp = `%sInsertOrUpdateStmt, err = %s.Prepare("INSERT INTO %s (%s) VALUES (%s) ON DUPLICATE KEY UPDATE %s, %s")
if err != nil {
	log.Fatal(err)
	}`

//StmtArgTemp append to argument list template
const StmtArgTemp = `argList = append(argList, m.%s)`

//StmtArgNilToDefaultTemp map nil argument to correspond default value template
const StmtArgNilToDefaultTemp = `if m.%s == nil {
	argList = append(argList, %v)
	} else {
	argList = append(argList, m.%s)
	}`

//StmtArgNilToDefaultBlockTemp nil arguments to correspond default value template
const StmtArgNilToDefaultBlockTemp = `argList := make([]interface{}, 0, %d)
%s`

//QueryOneFuncTemp QueryOneXXX() function template
const QueryOneFuncTemp = `func QueryOne%s(query string) (*%s, error) {
	for k, v := range %sMap {
		query = strings.Replace(query, k, v, -1)
		}
	row := %s.QueryRow(fmt.Sprintf("SELECT * FROM %s WHERE %%s", query))
	return %sFromRow(row)
	}`

//ModelSetFieldMethodTemp XXX.SetYYY() method template
const ModelSetFieldMethodTemp = `func (m *%s) Set%s(val %s, null bool) {
	if null {
		m.%s = nil
		m._IsStored = false
		return
		}
	m.%s = &val
	m._IsStored = false
	}`

//ModelGetFieldMethodTemp XXX.GetYYY() method template
const ModelGetFieldMethodTemp = `func (m *%s) Get%s() (%s, bool) {
	if m.%s == nil {
		return %s, true
		}
	return *m.%s, false
	}`

//ModelListTypeTemp XXXList struct definition template
const ModelListTypeTemp = `type %sList struct {
	Models []*%s
	Funcs []func(i, j int) int
	}`

//ModelCompareByIntMethodTemp compare by int field method template
const ModelCompareByIntMethodTemp = `By%s: func(i, j int) int {
	if *ml.Models[i].%s == *ml.Models[j].%s {
		return 0
	}
	if *ml.Models[i].%s < *ml.Models[j].%s {
		return -1
	}
	return 1
	},`

//ModelCompareByFloatMethodTemp compare by float field method template
const ModelCompareByFloatMethodTemp = `By%s: func(i, j int) int {
	if *ml.Models[i].%s == *ml.Models[j].%s {
		return 0
	}
	if *ml.Models[i].%s < *ml.Models[j].%s {
		return -1
	}
	return 1
	},`

//ModelCompareByStringMethodTemp compare by string field method template
const ModelCompareByStringMethodTemp = `By%s: func(i, j int) int {
	if *ml.Models[i].%s == *ml.Models[j].%s {
		return 0
	}
	if *ml.Models[i].%s < *ml.Models[j].%s {
		return -1
	}
	return 1
	},`

//ModelCompareByBoolMethodTemp compare by bool field method template
const ModelCompareByBoolMethodTemp = `By%s: func(i, j int) int {
	if *ml.Models[i].%s == *ml.Models[j].%s {
		return 0
	}
	if *ml.Models[i].%s == false && *ml.Models[j].%s == true {
		return -1
	}
	return 1
	},`

//ModelCompareByTimeMethodTemp compare by time.Time field method template
const ModelCompareByTimeMethodTemp = `By%s: func(i, j int) int {
	if ml.Models[i].%s.Equal(*ml.Models[j].%s) {
		return 0
	}
	if ml.Models[i].%s.Before(*ml.Models[j].%s) {
		return -1
	}
	return 1
	},`

//ModelSortMethodsStructFieldTypeTemp sort struct field definition template
const ModelSortMethodsStructFieldTypeTemp = `By%s func(i, j int) int`

//ModelSortMethodsStructTypeTemp sort struct definition template
const ModelSortMethodsStructTypeTemp = `type %sSortMethods struct {
	%s
	}`

//ModelSortMethodsStructFuncTemp sort struct method template
const ModelSortMethodsStructFuncTemp = `func (ml %sList) SortMethods() %sSortMethods {
	return %sSortMethods{
		%s
		}
	}`

//ModelListLenMethodTemp XXXList.Len() method template
const ModelListLenMethodTemp = `func (ml %sList) Len() int {
	return len(ml.Models)
	}`

//ModelListSwapMethodTemp XXXList().Swap() method template
const ModelListSwapMethodTemp = `func (ml %sList) Swap(i, j int) {
	ml.Models[i], ml.Models[j] = ml.Models[j], ml.Models[i]
	}`

//ModelListLessMethodTemp XXXList().Less() method template
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

//ModelSortFuncSwitchBlockTemp sort function switch block template
const ModelSortFuncSwitchBlockTemp = `case "%s": 
	%sList.Funcs = append(%sList.Funcs, sortMethods.By%s)`

//ModelSortFuncTemp XXXSortBy() function template
const ModelSortFuncTemp = `func %sSortBy(ml []*%s, desc bool, fields ...string) {
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
