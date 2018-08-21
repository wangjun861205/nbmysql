package nbmysql

//PackageTemp package template
var packageTemp = `package {{Package}}`

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

//dbTemp *sql.DB declaration template
const dbTemp = `var {{Database>ObjName}} *sql.DB`

//initFuncTemp init func template
const initFuncTemp = `func init() {
	db, err := sql.Open("mysql", "{{Database>Username}}:{{Database>Password}}@tcp({{Database>Address}})/{{Database>DatabaseName}}")
	if err != nil {
		panic(err)
	}
	{{Database>ObjName}} = db
	{{Block}}
}`

//mapElemTemp template of elements of map
const mapElemTemp = `"@{{Column>FieldName}}": "{{Column>ColumnName}}",`

//queryFieldMapTemp query field map declaration template
const queryFieldMapTemp = `var {{Table>ModelName}}Map = map[string]string {
	{{Block}}
	}`

//fieldTemp field in table model template
const fieldTemp = `{{Column>FieldName}} *{{Column>FieldType}}`

//ModelTemp model template
const modelTemp = `type {{Table>ModelName}} struct {
		{{Block}}
		_IsStored bool
}`

const modelGetFieldMethodTemp = `func (m *{{Table>ModelName}}) Get{{Column>FieldName}}() ({{Column>FieldType}}, bool) {
	if m.{{Column>FieldName}} == nil {
		return {{ZeroValue}}, true
		}
	return *m.{{Column>FieldName}}, false
	}`

const modelSetFieldMethodTemp = `func (m *{{Table>ModelName}}) Set{{Column>FieldName}}(val {{Column>FieldType}}, null bool) {
	if null {
		m.{{Column>FieldName}} = nil
		m._IsStored = false
		return
		}
	m.{{Column>FieldName}} = &val
	m._IsStored = false
	}`

//funcArgTemp arguments in function signature template
const funcArgTemp = `{{Column>ArgName}} *{{Column>FieldType}}`

//newModelAsignTemp asign statement in NewXXX() template
const newModelAsignTemp = `{{Column>FieldName}}: {{Column>ArgName}}`

//newModelFuncTemp NewXXX() template
const newModelFuncTemp = `func New{{Table>ModelName}}({{Args}}) *{{Table>ModelName}} {
		return &{{Table>ModelName}}{
			{{Asigns}}, 
			_IsStored: false}
	}`

//allModelFuncTemp AllXXX() template
const allModelFuncTemp = `func All{{Table>ModelName}}() ([]*{{Table>ModelName}}, error) {
		rows, err := {{Database>ObjName}}.Query("SELECT * FROM {{Table>TableName}}")
		if err != nil {
			return nil, err
		}
		list := make([]*{{Table>ModelName}}, 0, 256)
		for rows.Next() {
			model, err := {{Table>ModelName}}FromRows(rows)
			if err != nil {
				return nil, err
			}
			model._IsStored = true
			list = append(list, model)
		}
		return list, nil
	}`

//queryModelFuncTemp QueryXXX() template
const queryModelFuncTemp = `func Query{{Table>ModelName}}(query string) ([]*{{Table>ModelName}}, error) {
		for k, v := range {{Table>ModelName}}Map {
			query = strings.Replace(query, k, v, -1)
		}
		rows, err := {{Database>ObjName}}.Query(fmt.Sprintf("SELECT * FROM {{Table>TableName}} WHERE %s", query))
		if err != nil {
			return nil, err
		}
		list := make([]*{{Table>ModelName}}, 0, 256)
		for rows.Next() {
			model, err := {{Table>ModelName}}FromRows(rows)
			if err != nil {
				return nil, err
			}
			model._IsStored = true
			list = append(list, model)
		}
		return list, nil
	}`

const queryOneFuncTemp = `func QueryOne{{Table>ModelName}}(query string) (*{{Table>ModelName}}, error) {
	for k, v := range {{Table>ModelName}}Map {
		query = strings.Replace(query, k, v, -1)
		}
	row := {{Database>ObjName}}.QueryRow(fmt.Sprintf("SELECT * FROM {{Table>TableName}} WHERE %s", query))
	return {{Table>ModelName}}FromRow(row)
	}`

//foreignKeyMethodTemp foreign key relation method template
const foreignKeyMethodTemp = `func (m *{{FK>DstTab>ModelName}}) {{FK>DstTab>ModelName}}By{{FK>SrcCol>FieldName}}() (*{{FK>DstTab>ModelName}}, error) {
	if m.{{FK>SrcCol>FieldName}} == nil {
		return nil, errors.New("{{Table>ModelName}}.{{FK>SrcCol>FieldName}} must not be nil")
	}
	row := {{Database>ObjName}}.QueryRow("{{QueryStmt}}", m.{{FK>SrcCol>FieldName}})
	if row == nil {
		return nil, nbmysql.ErrRecordNotExists
	}
	model, err := {{FK>DstTab>ModelName}}FromRow(row)
	if err != nil {
		return nil, err
	}
	model._IsStored = true
	return model, nil
}`

//foreignKeyQuerySQLTemp query sql statement template in foreign key relation
const foreignKeyQuerySQLTemp = `SELECT * FROM {{FK>DstTab>TableName}} WHERE {{FK>DstCol>ColumnName}} = ?`

//reverseForeignKeyStructTypeTemp reverse foreign key struct definition template
const reverseForeignKeyStructTypeTemp = `type {{Table>ModelName}}To{{RFK>DstTab>ModelName}} struct {
	All func() ([]*{{RFK>DstTab>ModelName}}, error)
	Query func(query string) ([]*{{RFK>DstTab>ModelName}}, error)}`

//reverseForeignKeyAllSQLTemp sql statement template in reverse foreign key All() method
const reverseForeignKeyAllSQLTemp = `SELECT * FROM {{RFK>DstTab>TableName}} WHERE {{RFK>DstCol>ColumnName}} = ?`

//reverseForeignKeyAllMethodTemp All() method template in reverse foreign key relation struct
const reverseForeignKeyAllMethodTemp = `func() ([]*{{RFK>DstTab>ModelName}}, error) {
	if m.{{RFK>SrcCol>FieldName}} == nil {
		return nil, errors.New("{{Table>ModelName}}.{{RFK>SrcCol>FieldName}} must not be nil")
	}
	rows, err := {{Database>ObjName}}.Query("{{QueryStmt}}", *m.{{RFK>SrcCol>FieldName}})
	if err != nil {
		return nil, err
	}
	list := make([]*{{RFK>DstTab>ModelName}}, 0, 256)
	for rows.Next() {
		model, err := {{RFK>DstTab>ModelName}}FromRows(rows)
		if err != nil {
			return nil, err
		}
		model._IsStored = true
		list = append(list, model)
	}
	return list, nil
}`

//reverseForeignKeyQuerySQLTemp sql statement template in Query() method of reverse foreign key relation struct
const reverseForeignKeyQuerySQLTemp = `SELECT * FROM {{RFK>DstTab>TableName}} WHERE {{RFK>DstCol>ColumnName}} = ? AND %s`

//reverseForeignKeyQueryMethodTemp Query() method template in reverse foreign key relation struct
const reverseForeignKeyQueryMethodTemp = `func(query string) ([]*{{RFK>DstTab>ModelName}}, error) {
	if m.{{RFK>SrcCol>FieldName}} == nil {
		return nil, errors.New("{{Table>ModelName}}.{{RFK>SrcCol>FieldName}} must not be nil")
	}
	for k, v := range {{RFK>DstTab>ModelName}}Map {
		query = strings.Replace(query, k, v, -1)
	}
	rows, err := {{Database>ObjName}}.Query(fmt.Sprintf("{{QueryStmt}}", query), *m.{{RFK>SrcCol>FieldName}})
	if err != nil {
		return nil, err
	}
	list := make([]*{{RFK>DstTab>ModelName}}, 0, 256)
	for rows.Next() {
		model, err := {{RFK>DstTab>ModelName}}FromRows(rows)
		if err != nil {
			return nil, err
		}
		model._IsStored = true
		list = append(list, model)
	}
	return list, nil
}`

//reverseForeignKeyMethodTemp foreign key realtion method template
const reverseForeignKeyMethodTemp = `func (m *{{RFK>DstTab>ModelName}}) {{RFK>DstTab>ModelName}}By{{RFK>SrcCol>FieldName}}() {{Table>ModelName}}To{{RFK>DstTab>ModelName}} {
	return {{Table>ModelName}}To{{RFK>DstTab>ModelName}} {
		All: {{AllMethod}},
		Query: {{QueryMethod}},
	}
}`

//manyToManyStructTypeTemp many to many relation struct definition template
const manyToManyStructTypeTemp = `type {{Table>ModelName}}To{{MTM>DstTab>ModelName}} struct {
		All    func() ([]*{{MTM>DstTab>ModelName}}, error)
		Query func(query string) ([]*{{MTM>DstTab>ModelName}}, error)
		Add func({{MTM>DstTab>ArgName}} *{{MTM>DstTab>ModelName}}) error
		Remove func({{MTM>DstTab>ArgName}} *{{MTM>DstTab>ModelName}}) error
	}`

//manyToManyAllSQLTemp sql statement template in All() method of many to many relation struct
const manyToManyAllSQLTemp = `SELECT {{MTM>DstTab>TableName}}.* FROM {{Table>TableName}} JOIN {{MTM>MidTab>TableName}} ON {{Table>TableName}}.{{MTM>SrcCol>ColumnName}}={{MTM>MidTab>TableName}}.{{MTM>MidLeftCol>ColumnName}} JOIN {{MTM>DstTab>TableName}} on {{MTM>MidTab>TableName}}.{{MTM>MidRightCol>ColumnName}} = {{MTM>DstTab>TableName}}.{{MTM>DstCol>ColumnName}} WHERE {{Table>TableName}}.{{MTM>DstCol>ColumnName}} = ?`

//manyToManyAllMethodTemp All() method template in many to many relation struct
const manyToManyAllMethodTemp = `func() ([]*{{MTM>DstTab>ModelName}}, error) {
	rows, err := {{Database>ObjName}}.Query("{{QueryStmt}}", *m.{{MTM>SrcCol>FieldName}})
	if err != nil {
		return nil, err
	}
	list := make([]*{{MTM>DstTab>ModelName}}, 0, 256)
	for rows.Next() {
		model, err := {{MTM>DstTab>ModelName}}FromRows(rows)
		if err != nil {
			return nil, err
		}
		model._IsStored = true
		list = append(list, model)
	}
	return list, nil
}`

//manyToManyQuerySQLTemp sql statement template in Query() method of many to many relation struct
const manyToManyQuerySQLTemp = `SELECT {{MTM>DstTab>TableName}}.* FROM {{Table>TableName}} JOIN {{MTM>MidTab>TableName}} ON {{Table>TableName}}.{{MTM>SrcCol>ColumnName}}={{MTM>MidTab>TableName}}.{{MTM>MidLeftCol>ColumnName}} JOIN {{MTM>DstTab>TableName}} on {{MTM>MidTab>TableName}}.{{MTM>MidRightCol>ColumnName}} = {{MTM>DstTab>TableName}}.{{MTM>DstCol>ColumnName}} WHERE {{Table>TableName}}.{{MTM>SrcCol>ColumnName}} = ? AND %s`

//manyToManyQueryMethodTemp Query() method template in many to many relation struct
const manyToManyQueryMethodTemp = `func(query string) ([]*{{MTM>DstTab>ModelName}}, error) {
	for k, v := range {{MTM>DstTab>ModelName}}Map {
		query = strings.Replace(query, k, v, -1)
	}
	rows, err := {{Database>ObjName}}.Query(fmt.Sprintf("{{QueryStmt}}", query), *m.{{MTM>SrcCol>FieldName}})
	if err != nil {
		return nil, err
	}
	list := make([]*{{MTM>DstTab>ModelName}}, 0, 256)
	for rows.Next() {
		model, err := {{MTM>DstTab>ModelName}}FromRows(rows)
		if err != nil {
			return nil, err
		}
		model._IsStored = true
		list = append(list, model)
	}
	return list, nil
}`

//manyToManyAddMethodTemp Add() method template in many to many relation struct
const manyToManyAddMethodTemp = `func({{MTM>DstTab>ArgName}} *{{MTM>DstTab>ModelName}}) error {
	if !m._IsStored {
		return errors.New("{{Table>ModelName}} model is not stored in database")
	}
	if !{{MTM>DstTab>ArgName}}._IsStored {
		return errors.New("{{MTM>DstTab>ModelName}} model is not stored in database")
	}
	_, err := {{Table>ModelName}}To{{MTM>DstTab>ModelName}}InsertStmt.Exec(m.{{MTM>SrcCol>FieldName}}, {{MTM>DstTab>ArgName}}.{{MTM>DstCol>FieldName}})
	return err
}`

const manyToManyRemoveMethodTemp = `func({{MTM>DstTab>ArgName}} *{{MTM>DstTab>ModelName}}) error {
	if !m._IsStored {
		return errors.New("{{Table>ModelName}} model is not stored in database")
	}
	if !{{MTM>DstTab>ArgName}}._IsStored {
		return errors.New("{{MTM>DstTab>ModelName}} model is not stored in database")
	}
	_, err := {{Table>ModelName}}To{{MTM>DstTab>ModelName}}DeleteStmt.Exec(m.{{MTM>SrcCol>FieldName}}, {{MTM>DstTab>ArgName}}.{{MTM>DstCol>FieldName}})
	return err
}`

//manyToManyMethodTemp many to many relation struct declaration template
const manyToManyMethodTemp = `func (m *{{Table>ModelName}}) {{MTM>DstTab>ModelName}}By{{MTM>SrcCol>FieldName}}() {{Table>ModelName}}To{{MTM>DstTab>ModelName}} {
	return {{Table>ModelName}}To{{MTM>DstTab>ModelName}}{
		All: {{AllMethod}},
		Query: {{QueryMethod}},
		Add: {{AddMethod}},
		Remove: {{RemoveMethod}},
	}
}`

const stmtArgNilToDefaultTemp = `if m.{{Column>FieldName}} == nil {
	argList = append(argList, {{DefaultValue}})
	} else {
	argList = append(argList, m.{{Column>FieldName}})
	}`

const stmtArgTemp = `argList = append(argList, m.{{Column>FieldName}})`

const stmtArgNilToDefaultBlockTemp = `argList := make([]interface{}, 0, {{Length}})
{{Block}}`

//modelInsertMethodTemp XXX.Insert() method template
const modelInsertMethodTemp = `func (m *{{Table>ModelName}}) Insert() error {
		err := m.check()
		if err != nil {
			return err
		}
		{{Block}}
		res, err := {{Table>ModelName}}InsertStmt.Exec(argList...)
		if err != nil {
			return err
		}
		lastInsertId, err := res.LastInsertId()
		if err != nil {
			return err
		}
		m.{{Table>AutoIncrement>FieldName}} = &lastInsertId
		m._IsStored = true
		return nil
}`

//modelUpdateMethodTemp XXX.Update() method template
const modelUpdateMethodTemp = `func (m *{{Table>ModelName}}) Update() error {
	if !m._IsStored {
		return nbmysql.ErrModelNotStoredInDB
	}
	err := m.check()
	if err != nil {
		return err
	}
	{{Block}}
	argList = append(argList, m.{{Table>AutoIncrement>FieldName}})
	_, err = {{Table>ModelName}}UpdateStmt.Exec(argList...)
	if err != nil {
		return err
	}
	return nil
}`

//modelInsertOrUpdateMethodTemp XXX.InsertOrUpdate() method template
const modelInsertOrUpdateMethodTemp = `func (m *{{Table>ModelName}}) InsertOrUpdate() error {
	err := m.check()
	if err != nil {
		return err
	}
	{{Block}}
	argList = append(argList, argList...)
	res, err := {{Table>ModelName}}InsertOrUpdateStmt.Exec(argList...)
	if err != nil {
		return err
	}
	lastInsertId, err := res.LastInsertId()
	if err != nil {
		return err
	}
	m.{{Table>AutoIncrement>FieldName}} = &lastInsertId
	m._IsStored = true
	return nil
}`

//modelDeleteMethodTemp XXX.Delete() method template
const modelDeleteMethodTemp = `func (m *{{Table>ModelName}}) Delete() error {
	if !m._IsStored {
		return nbmysql.ErrModelNotStoredInDB
	}
	_, err := {{Table>ModelName}}DeleteStmt.Exec(m.{{Table>AutoIncrement>FieldName}})
	if err != nil {
		return err
	}
	m._IsStored = false
	return nil
	}`

//newMiddleTypeTemp craete middle type template
const newMiddleTypeTemp = `_{{Column>ArgName}} := new(nbmysql.{{Column>MidType}})`

//modelFromRowsFuncTemp XXXFromRows() function template
const modelFromRowsFuncTemp = `func {{Table>ModelName}}FromRows(rows *sql.Rows) (*{{Table>ModelName}}, error) {
		{{Block}}
		err := rows.Scan({{MidArgs}})
		if err != nil {
			return nil, err
		}
		return &{{Table>ModelName}}{
			{{MidToGos}}, true}, nil
	}`

//modelFromRowFuncTemp XXXFromRow() function template
const modelFromRowFuncTemp = `func {{Table>ModelName}}FromRow(row *sql.Row) (*{{Table>ModelName}}, error) {
	{{Block}}
	err := row.Scan({{MidArgs}})
	if err != nil {
		return nil, err
	}
	return &{{Table>ModelName}}{
		{{MidToGos}}, true}, nil
}`

const fieldCheckNullTemp = `if m.{{Column>FieldName}} == nil {
	return errors.New("{{Table>ModelName}}.{{Column>FieldName}} can not be null")
	}`

const modelCheckMethodTemp = `func (m *{{Table>ModelName}}) check() error {
	{{Block}}
	return nil
	}`

const insertStmtTemp = `var {{Table>ModelName}}InsertStmt *sql.Stmt`

const updateStmtTemp = `var {{Table>ModelName}}UpdateStmt *sql.Stmt`

const deleteStmtTemp = `var {{Table>ModelName}}DeleteStmt *sql.Stmt`

const insertOrUpdateStmtTemp = `var {{Table>ModelName}}InsertOrUpdateStmt *sql.Stmt`

const insertMidStmtTemp = `var {{Table>ModelName}}To{{MTM>DstTab>ModelName}}InsertStmt *sql.Stmt`

const manyToManyDeleteStmtTemp = `var {{Table>ModelName}}To{{MTM>DstTab>ModelName}}DeleteStmt *sql.Stmt`

const insertStmtInitTemp = `{{Table>ModelName}}InsertStmt, err = {{Database>ObjName}}.Prepare("INSERT INTO {{Table>TableName}} ({{Columns}}) VALUES ({{Values}})")
if err != nil {
	log.Fatal(err)
	}`

const updateStmtInitTemp = `{{Table>ModelName}}UpdateStmt, err = {{Database>ObjName}}.Prepare("UPDATE {{Table>TableName}} SET {{Updates}} WHERE {{Table>AutoIncrement>ColumnName}} = ?")
if err != nil {
	log.Fatal(err)
	}`

const deleteStmtInitTemp = `{{Table>ModelName}}DeleteStmt, err = {{Database>ObjName}}.Prepare("DELETE FROM {{Table>TableName}} WHERE {{Table>AutoIncrement>ColumnName}} = ?")
if err != nil {
	log.Fatal(err)
	}`

const updateLastInsertIDTemp = `{{Table>AutoIncrement>ColumnName}} = LAST_INSERT_ID({{Table>AutoIncrement>ColumnName}})`

const insertOrUpdateStmtInitTemp = `{{Table>ModelName}}InsertOrUpdateStmt, err = {{Database>ObjName}}.Prepare("INSERT INTO {{Table>TableName}} ({{Columns}}) VALUES ({{Values}}) ON DUPLICATE KEY UPDATE {{UpdateLastInsertID}}, {{Updates}}")
if err != nil {
	log.Fatal(err)
	}`

const insertMidStmtInitTemp = `{{Table>ModelName}}To{{MTM>DstTab>ModelName}}InsertStmt, err = {{Database>ObjName}}.Prepare("INSERT INTO {{MTM>MidTab>TableName}} ({{MTM>MidLeftCol>ColumnName}}, {{MTM>MidRightCol>ColumnName}}) VALUES (?, ?)")
if err != nil {
	log.Fatal(err)
	}`

const manyToManyDeleteStmtInitTemp = `{{Table>ModelName}}To{{MTM>DstTab>ModelName}}DeleteStmt, err = {{Database>ObjName}}.Prepare("DELETE FROM {{MTM>MidTab>TableName}} WHERE {{MTM>MidLeftCol>ColumnName}} = ? AND {{MTM>MidRightCol>ColumnName}} = ?")`

//FuncArgNameTemp arguments name in function body template
const funcArgNameTemp = `{{Column>ArgName}}`

const updateColumnTemp = `{{Column>ColumnName}} = ?`

const modelListTypeTemp = `type {{Table>ModelName}}List struct {
	Models []*{{Table>ModelName}}
	Funcs []func(i, j int) int
	}`

const modelCompareByIntMethodTemp = `By{{Column>FieldName}}: func(i, j int) int {
	if *ml.Models[i].{{Column>FieldName}} == *ml.Models[j].{{Column>FieldName}} {
		return 0
	}
	if *ml.Models[i].{{Column>FieldName}} < *ml.Models[j].{{Column>FieldName}} {
		return -1
	}
	return 1
	},`

const modelCompareByFloatMethodTemp = `By{{Column>FieldName}}: func(i, j int) int {
	if *ml.Models[i].{{Column>FieldName}} == *ml.Models[j].{{Column>FieldName}} {
		return 0
	}
	if *ml.Models[i].{{Column>FieldName}} < *ml.Models[j].{{Column>FieldName}} {
		return -1
	}
	return 1
	},`

const modelCompareByStringMethodTemp = `By{{Column>FieldName}}: func(i, j int) int {
	if *ml.Models[i].{{Column>FieldName}} == *ml.Models[j].{{Column>FieldName}} {
		return 0
	}
	if *ml.Models[i].{{Column>FieldName}} < *ml.Models[j].{{Column>FieldName}} {
		return -1
	}
	return 1
	},`

const modelCompareByBoolMethodTemp = `By{{Column>FieldName}}: func(i, j int) int {
	if *ml.Models[i].{{Column>FieldName}} == *ml.Models[j].{{Column>FieldName}} {
		return 0
	}
	if *ml.Models[i].{{Column>FieldName}} == false && *ml.Models[j].{{Column>FieldName}} == true {
		return -1
	}
	return 1
	},`

const modelCompareByTimeMethodTemp = `By{{Column>FieldName}}: func(i, j int) int {
	if ml.Models[i].{{Column>FieldName}}.Equal(*ml.Models[j].{{Column>FieldName}}) {
		return 0
	}
	if ml.Models[i].{{Column>FieldName}}.Before(*ml.Models[j].{{Column>FieldName}}) {
		return -1
	}
	return 1
	},`

const modelSortMethodsStructFieldTypeTemp = `By{{Column>FieldName}} func(i, j int) int`

const modelSortMethodsStructTypeTemp = `type {{Table>ModelName}}SortMethods struct {
	{{FieldTypeBlock}}
	}`

const modelSortMethodsStructFuncTemp = `func (ml {{Table>ModelName}}List) SortMethods() {{Table>ModelName}}SortMethods {
	return {{Table>ModelName}}SortMethods{
		{{CompareMethodBlock}}
		}
	}`

const modelListLenMethodTemp = `func (ml {{Table>>ModelName}}List) Len() int {
	return len(ml.Models)
	}`

const modelListSwapMethodTemp = `func (ml {{Table>ModelName}}List) Swap(i, j int) {
	ml.Models[i], ml.Models[j] = ml.Models[j], ml.Models[i]
	}`

const modelListLessMethodTemp = `func (ml {{Table>ModelName}}List) Less(i, j int) bool {
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

const modelSortFuncSwitchBlockTemp = `case "{{Column>FieldName}}": 
	{{Table>ArgName}}List.Funcs = append({{Table>ArgName}}List.Funcs, sortMethods.By{{Column>FieldName}})`

const modelSortFuncTemp = `func {{Table>ModelName}}SortBy(ml []*{{Table>ModelName}}, desc bool, fields ...string) {
	{{Table>ArgName}}List := {{Table>ModelName}}List {
		Models: ml,
		Funcs: make([]func(i, j int) int, 0, {{ColumnNumber}}),
		}
	sortMethods := {{Table>ArgName}}List.SortMethods()
	for _, field := range fields {
		switch field {
			{{SwitchBlock}}
		}
	}
	if desc {
		sort.Sort(sort.Reverse({{Table>ArgName}}List))
	} else {
		sort.Sort({{Table>ArgName}}List)
	}
}`

const modelCountStmtDeclareTemp = `var {{Table>ModelName}}CountStmt *sql.Stmt`

const modelCountStmtInitTemp = `{{Table>ModelName}}CountStmt, err = {{Database>ObjName}}.Prepare("SELECT COUNT(*) FROM {{Table>TableName}}")
if err != nil {
	log.Fatal(err)
	}`

const modelCountFuncTemp = `func {{Table>ModelName}}Count() (int64, error) {
	var count int64
	err := {{Table>ModelName}}CountStmt.QueryRow().Scan(&count)
	if err != nil {
		return -1, err
	}
	return count, nil
}`
