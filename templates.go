package nbmysql

//PackageTemp package template
var packageTemp = `package {{ DB.Package }}
`

//ImportTemp import block template
const importTemp = `import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
	"sort"
	"github.com/wangjun861205/nbmysql"
)`

//dbTemp *sql.DB declaration template
const dbTemp = `var {{ DB.ObjName }} *sql.DB
`

//initFuncTemp init func template
const initFuncTemp = `func init() {
	db, err := sql.Open("mysql", "{{ DB.Username }}:{{ DB.Password }}@tcp({{ DB.Address }})/{{ DB.DatabaseName }}")
	if err != nil {
		panic(err)
	}
	{{ DB.ObjName }} = db
}
`

const whereMapTemp = `{{ for _, tab in DB.Tables }}
	var {{ tab.ModelName }}Map = map[string]string {
		{{ for _, col in tab.Columns }}
			"@{{ col.FieldName }}": "{{ col.ColumnName }}",
		{{ endfor }}
	}
{{ endfor }}`

const fromRowFuncTemp = `{{ for _, tab in DB.Tables }}
	func {{ tab.ModelName }}FromRow(row *sql.Row) (*{{ tab.ModelName}}, error) {
		m := {{ tab.ModelName }}{}
		addrList := make([]interface{}, 0, 16)
		{{ for _, col in tab.Columns }}
			addrList = append(addrList, &(m.{{ col.FieldName }}))
		{{ endfor }}
		err := row.Scan(addrList...)
		if err != nil {
			return nil, err
		}
		return &m, nil
	}
{{ endfor }}`

const fromRowsFuncTemp = `{{ for _, tab in DB.Tables }}
	func {{ tab.ModelName }}FromRows(rows *sql.Rows) ({{ tab.ModelName }}List, error) {
		l := make({{ tab.ModelName }}List, 0, 128)
		for rows.Next() {
			m := {{ tab.ModelName }}{}
			addrList := make([]interface{}, 0, 16)
			{{ for _, col in tab.Columns }}
				addrList = append(addrList, &(m.{{ col.FieldName }}))
			{{ endfor }}
			err := rows.Scan(addrList...)
			if err != nil {
				return l, err
			}
			l = append(l, &m)
		}
		return l, nil
	}
{{ endfor }}`

const allModelFuncTemp = `{{ for _, tab in DB.Tables }}
	func All{{ tab.ModelName }}() ({{ tab.ModelName }}List, error) {
		rows, err := {{ DB.ObjName }}.Query("SELECT * FROM {{ tab.TableName }}")
		if err != nil {
			return nil, err
		}
		return {{ tab.ModelName }}FromRows(rows)
	}
{{ endfor }}`

const queryFuncTemp = `{{ for _, tab in DB.Tables }}
	func Query{{ tab.ModelName }}(where string) ({{ tab.ModelName }}List, error) {
		for k, v := range {{ tab.ModelName }}Map {
			where = strings.Replace(where, k, v, -1)
		}
		rows, err := {{ DB.ObjName }}.Query(fmt.Sprintf("SELECT * FROM {{ tab.TableName }} WHERE %s", where))
		if err != nil {
			return nil, err
		}
		return {{ tab.ModelName }}FromRows(rows)
	}
{{ endfor }}`

const queryOneFuncTemp = `{{ for _, tab in DB.Tables }}
	func QueryOne{{ tab.ModelName }}(where string) (*{{ tab.ModelName }}, error) {
		for k, v := range {{ tab.ModelName }}Map {
			where = strings.Replace(where, k, v, -1)
		}
		row := {{ DB.ObjName }}.QueryRow(fmt.Sprintf("SELECT * FROM {{ tab.TableName }} WHERE %s", where))
		return {{ tab.ModelName }}FromRow(row)
	}
{{ endfor }}`

const foreignKeyMethodTemp = `{{ for _, tab in DB.Tables }}
	{{ for _, fk in tab.ForeignKeys }}
		func (m *{{ tab.ModelName }}) {{ fk.DstTab.ModelName }}() (*{{ fk.DstTab.ModelName }}, error) {
			isValid, isNull := m.{{ fk.SrcCol.FieldName }}.IsValid(), m.{{ fk.SrcCol.FieldName }}.IsNull()
			if !isValid || isNull {
				return nil, errors.New("{{ tab.ModelName }}.{{ fk.SrcCol.FieldName }} cannot be null")
			}
			stmtStr := fmt.Sprintf("SELECT * FROM {{ fk.DstTab.TableName }} WHERE {{ fk.DstCol.ColumnName }} = %s", m.{{ fk.SrcCol.FieldName }}.SQLVal())
			row := {{ DB.ObjName }}.QueryRow(stmtStr)
			return {{ fk.DstTab.ModelName }}FromRow(row)
		}
	{{ endfor }}
{{ endfor }}`

const reverseForeignKeyMethodTemp = `{{ for _, tab in DB.Tables }}
	{{ for _, rfk in tab.ReverseForeignKeys }}
		func (m *{{ tab.ModelName }}) {{ rfk.DstTab.ModelName }}List() struct {
			All func() ({{ rfk.DstTab.ModelName }}List, error)
			Query func(string) ({{ rfk.DstTab.ModelName }}List, error)
		} {
			return struct {
				All func() ({{ rfk.DstTab.ModelName }}List, error)
				Query func(string) ({{ rfk.DstTab.ModelName }}List, error)
				} {
					All: func() ({{ rfk.DstTab.ModelName }}List, error) {
						isValid, isNull := m.{{ rfk.SrcCol.FieldName }}.IsValid(), m.{{ rfk.SrcCol.FieldName }}.IsNull()
						if !isValid || isNull {
							return nil, errors.New("{{ tab.ModelName }}.{{ rfk.SrcCol.FieldName }} cannot be null")
						}
						stmtStr := fmt.Sprintf("SELECT * FROM {{ rfk.DstTab.TableName }} WHERE {{ rfk.DstCol.ColumnName }} = %s", m.{{ rfk.SrcCol.FieldName }}.SQLVal())
						rows, err := {{ DB.ObjName }}.Query(stmtStr)
						if err != nil {
							return nil, err
						}
						return {{ rfk.DstTab.ModelName }}FromRows(rows)
						},
					Query: func(where string) ({{ rfk.DstTab.ModelName }}List, error) {
						isValid, isNull := m.{{ rfk.SrcCol.FieldName }}.IsValid(), m.{{ rfk.SrcCol.FieldName }}.IsNull()
						if !isValid || isNull {
							return nil, errors.New("{{ tab.ModelName }}.{{ rfk.SrcCol.FieldName }} cannot be null")
						}
						var stmtStr string
						for k, v := range {{ tab.ModelName }}Map {
							where = strings.Replace(where, k, v, -1)
						}
						if where != "" {
							stmtStr = fmt.Sprintf("SELECT * FROM {{ rfk.DstTab.TableName }} WHERE {{ rfk.DstCol.ColumnName }} = %s AND %s", m.{{ rfk.SrcCol.FieldName }}.SQLVal(), where)
						} else {
							stmtStr = fmt.Sprintf("SELECT * FROM {{ rfk.DstTab.TableName }} WHERE {{ rfk.DstCol.ColumnName }} = %s", m.{{ rfk.SrcCol.FieldName }}.SQLVal())
						}
						rows, err := {{ DB.ObjName }}.Query(stmtStr)
						if err != nil {
							return nil, err
						}
						return {{ rfk.DstTab.ModelName }}FromRows(rows)
						},
					}
			}
	{{ endfor }}
{{ endfor }}`

const manyToManyMethodTemp = `{{ for _, tab in DB.Tables }} 
	{{ for _, mtm in tab.ManyToManys }}
		func (m *{{ tab.ModelName }}) {{ mtm.DstTab.ModelName }}List() struct {
			All func() ({{ mtm.DstTab.ModelName }}List, error)
			Query func(string) ({{ mtm.DstTab.ModelName }}List, error)
			Add func(*{{ mtm.DstTab.ModelName }}) error
			Remove func(*{{ mtm.DstTab.ModelName }}) error
			} {
				return struct {
					All func() ({{ mtm.DstTab.ModelName }}List, error)
					Query func(string) ({{ mtm.DstTab.ModelName }}List, error)
					Add func(*{{ mtm.DstTab.ModelName }}) error
					Remove func(*{{ mtm.DstTab.ModelName }}) error
					} {
						All: func() ({{ mtm.DstTab.ModelName }}List, error) {
							isValid, isNull := m.{{ mtm.SrcCol.FieldName }}.IsValid(), m.{{ mtm.SrcCol.FieldName }}.IsNull()
							if !isValid || isNull {
								return nil, errors.New("{{ tab.ModelName }}.{{ mtm.SrcCol.FieldName }} cannot be null")
							}
							stmtStr := fmt.Sprintf("SELECT {{ mtm.DstTab.TableName }}.* FROM {{ tab.TableName }} JOIN {{ mtm.MidTab.TableName }} ON {{ tab.TableName }}.{{ mtm.SrcCol.ColumnName }} = {{ mtm.MidTab.TableName }}.{{ mtm.MidLeftCol.ColumnName }} JOIN {{ mtm.DstTab.TableName }} ON {{ mtm.MidTab.TableName }}.{{ mtm.MidRightCol.ColumnName }} = {{ mtm.DstTab.TableName}}.{{ mtm.DstCol.ColumnName }} WHERE {{ tab.TableName }}.{{ mtm.SrcCol.ColumnName }} = %s", m.{{ mtm.SrcCol.FieldName }}.SQLVal())
							rows, err := {{ DB.ObjName }}.Query(stmtStr)
							if err != nil {
								return nil, err
							}
							return {{ mtm.DstTab.ModelName }}FromRows(rows)
						},
						Query: func(where string) ({{ mtm.DstTab.ModelName }}List, error) {
							isValid, isNull := m.{{ mtm.SrcCol.FieldName }}.IsValid(), m.{{ mtm.SrcCol.FieldName }}.IsNull()
							if !isValid || isNull {
								return nil, errors.New("{{ tab.ModelName }}.{{ mtm.SrcCol.FieldName }} cannot be null")
							}
							for k, v := range {{ mtm.DstTab.ModelName }}Map {
								where = strings.Replace(where, k, v, -1)
							}
							var stmtStr string
							if where == "" {
								stmtStr = fmt.Sprintf("SELECT {{ mtm.DstTab.TableName }}.* FROM {{ tab.TableName }} JOIN {{ mtm.MidTab.TableName }} ON {{ tab.TableName }}.{{ mtm.SrcCol.ColumnName }} = {{ mtm.MidTab.TableName }}.{{ mtm.MidLeftCol.ColumnName }} JOIN {{ mtm.DstTab.TableName }} ON {{ mtm.MidTab.TableName }}.{{ mtm.MidRightCol.ColumnName }} = {{ mtm.DstTab.TableName}}.{{ mtm.DstCol.ColumnName }} WHERE {{ tab.TableName }}.{{ mtm.SrcCol.ColumnName }} = %s", m.{{ mtm.SrcCol.FieldName }}.SQLVal())
							} else {
								stmtStr = fmt.Sprintf("SELECT {{ mtm.DstTab.TableName }}.* FROM {{ tab.TableName }} JOIN {{ mtm.MidTab.TableName }} ON {{ tab.TableName }}.{{ mtm.SrcCol.ColumnName }} = {{ mtm.MidTab.TableName }}.{{ mtm.MidLeftCol.ColumnName }} JOIN {{ mtm.DstTab.TableName }} ON {{ mtm.MidTab.TableName }}.{{ mtm.MidRightCol.ColumnName }} = {{ mtm.DstTab.TableName}}.{{ mtm.DstCol.ColumnName }} WHERE {{ tab.TableName }}.{{ mtm.SrcCol.ColumnName }} = %s AND %s", m.{{ mtm.SrcCol.FieldName }}.SQLVal(), where)
							}
							rows, err := {{ DB.ObjName }}.Query(stmtStr)
							if err != nil {
								return nil, err
							}
							return {{ mtm.DstTab.ModelName }}FromRows(rows)
						},
						Add: func(mm *{{ mtm.DstTab.ModelName }}) error {
							isValid, isNull := m.{{ mtm.SrcCol.FieldName }}.IsValid(), m.{{ mtm.SrcCol.FieldName }}.IsNull()
							if !isValid || isNull {
								return errors.New("{{ tab.ModelName }}.{{ mtm.SrcCol.FieldName }} cannot be null")
							}
							isValid, isNull = mm.{{ mtm.DstCol.FieldName }}.IsValid(), mm.{{ mtm.DstCol.FieldName }}.IsNull()
							if !isValid || isNull {
								return errors.New("{{ mtm.DstTab.ModelName }}.{{ mtm.DstCol.FieldName }} cannot be null")
							}
							stmtStr := fmt.Sprintf("INSERT INTO {{ mtm.MidTab.TableName }} ({{ mtm.MidLeftCol.ColumnName }}, {{ mtm.MidRightCol.ColumnName }}) VALUES (%s, %s)", m.{{ mtm.SrcCol.FieldName }}.SQLVal(), mm.{{ mtm.DstCol.FieldName }}.SQLVal())
							_, err := {{ DB.ObjName }}.Exec(stmtStr)
							return err
						},
						Remove: func(mm *{{ mtm.DstTab.ModelName }}) error {
							isValid, isNull := m.{{ mtm.SrcCol.FieldName }}.IsValid(), m.{{ mtm.SrcCol.FieldName }}.IsNull()
							if !isValid || isNull {
								return errors.New("{{ tab.ModelName }}.{{ mtm.SrcCol.FieldName }} cannot be null")
							}
							isValid, isNull = mm.{{ mtm.DstCol.FieldName }}.IsValid(), mm.{{ mtm.DstCol.FieldName }}.IsNull()
							if !isValid || isNull {
								return errors.New("{{ mtm.DstTab.ModelName }}.{{ mtm.DstCol.FieldName }} cannot be null")
							}
							stmtStr := fmt.Sprintf("DELETE FROM {{ mtm.MidTab.TableName }} WHERE {{ mtm.MidLeftCol.ColumnName }} = %s AND {{ mtm.MidRightCol.ColumnName }} = %s", m.{{ mtm.SrcCol.FieldName }}.SQLVal(), mm.{{ mtm.DstCol.FieldName }}.SQLVal())
							_, err := {{ DB.ObjName }}.Exec(stmtStr)
							return err
						},
					}
				}
	{{ endfor }}
{{ endfor }}`

const stmtArgNilToDefaultTemp = `if m.{{Column.FieldName}} == nil {
	argList = append(argList, {{DefaultValue}})
	} else {
	argList = append(argList, m.{{Column.FieldName}})
	}`

const stmtArgTemp = `argList = append(argList, m.{{Column.FieldName}})`

const stmtArgNilToDefaultBlockTemp = `argList := make([]interface{}, 0, {{Length}})
{{Block}}`

const modelInsertMethodTemp = `{{ for _, tab in DB.Tables }}
func (m *{{ tab.ModelName }})Insert() error {
	colList := make([]string, 0, 16)
	valList := make([]string, 0, 16)
	{{ for _, col in tab.Columns }}
		if m.{{ col.FieldName }}.IsValid() {
			if m.{{ col.FieldName }}.IsNull() {
				colList = append(colList, {{ col.ColumnName }})
				valList = append(valList, "NULL")
			} else {
				colList = append(colList, {{ col.ColumnName }})
				valList = append(valList, m.{{ col.FieldName }}.SQLVal())
			}
		}
	{{ endfor }}
	stmtStr := fmt.Sprintf("INSERT INTO {{ tab.TableName }} (%s) VALUES (%s)", strings.Join(colList, ", "), strings.Join(valList, ", "))
	_, err := {{ DB.ObjName }}.Exec(stmtStr)
	return err
}
{{ endfor }}`

const modelUpdateMethodTemp = `{{ for _, tab in DB.Tables }}
	func (m *{{ tab.ModelName }})Update() error {
		colList := make([]string, 0, 16)
		valList := make([]string, 0, 16)
		if !m.{{ tab.PrimaryKey.FieldName }}.IsValid() || m.{{ tab.PrimaryKey.FieldName }}.IsNull() {
			return errors.New("{{ tab.ModelName }}.Update(): primary key cannot be invalid or null")
		}
		{{ for _, col in tab.Columns }}
			{{ if col.ColumnName != tab.PrimaryKey.ColumnName }}
				if m.{{ col.FieldName }}.IsValid() {
					colList = append(colList, {{ col.ColumnName }})
					if m.{{ col.FieldName }}.IsNull() {
						valList = append(valList, "NULL")
					} else {
						valList = append(valList, m.{{ col.FieldName }}.SQLVal())
					}
				}
			{{ endif }}
		{{ endfor }}
		if len(colList) == 0 {
			return errors.New("no valid field to update")
		}
		setList := make([]string, len(colList))
		for i := range colList {
			setList[i] = colList[i] + "=" + valList[i]
		}
		stmtStr := fmt.Sprintf("UPDATE {{ tab.TableName }} SET %s WHERE {{ tab.PrimaryKey.ColumnName }} = %s", strings.Join(setList, ", "), m.{{ tab.PrimaryKey.FieldName }}.SQLVal())
		_, err := {{ DB.ObjName }}.Exec(stmtStr)
		return err
	}
{{ endfor }}
		`
const modelInsertOrUpdateMethodTemp = `{{ for _, tab in DB.Tables }}
	func (m *{{ tab.ModelName }}) InsertOrUpdate() error {
		insertColList := make([]string, 0, 16)
		insertValList := make([]string, 0, 16)
		updateColList := make([]string, 0, 16)
		updateValList := make([]string, 0, 16)
		{{ for _, col in tab.Columns }}
			{{ if col.ColumnName != tab.AutoIncrement.ColumnName }}
				if m.{{ col.FieldName }}.IsValid() {
					insertColList = append(insertColList, {{ col.ColumnName }})
					{{ if col.ColumnName != tab.PrimaryKey.ColumnName }}
						updateColList = append(updateColList, {{ col.ColumnName }})
					{{ endif }}
					if m.{{ col.FieldName }}.IsNull() {
						insertValList = append(insertValList, "NULL")
						{{ if col.ColumnName != tab.PrimaryKey.ColumnName }}
							updateValList = append(updateValList, "NULL")
						{{ endif }}
					} else {
						insertValList = append(insertValList, m.{{ col.FieldName }}.SQLVal())
						{{ if col.ColumnName != tab.PrimaryKey.ColumnName }}
							updateValList = append(updateValList, m.{{ col.FieldName }}.SQLVal())
						{{ endif }}
					}
				}
			{{ endif }}
		{{ endfor }}
		setList := make([]string, len(updateColList))
		for i := range updateColList {
			setList[i] = updateColList[i] + "=" + updateValList[i]
		}
		stmtStr := fmt.Sprintf("INSERT INTO {{ tab.TableName }} (%s) VALUES (%s) ON DUPLICATE KEY UPDATE {{ tab.AutoIncrement.ColumnName }} = LAST_INSERT_ID({{ tab.AutoIncrement.ColumnName }}), %s", strings.Join(insertColList, ", "), strings.Join(insertValList, ", "), strings.Join(setList, ", "))
		res, err := {{ DB.ObjName }}.Exec(stmtStr)
		if err != nil {
			return err
		}
		lastInsertID, err := res.LastInsertId()
		if err != nil {
			return err
		}
		m.{{ tab.AutoIncrement.FieldName }}.Set(lastInsertID)
		return nil
	}
{{ endfor }}`

//newMiddleTypeTemp craete middle type template
// const newMiddleTypeTemp = `_{{Column.ArgName}} := new(nbmysql.{{Column.MidType}})`

const modelCheckMethodTemp = `{{ for _, tab in DB.Tables }}
	func (m *{{ tab.ModelName }}) check() error {
		{{ for _, col in tab.Columns }}
			{{ if col.ColumnName != tab.AutoIncrement.ColumnName && col.Nullable == false && col.Default == "" }}
				if !m.{{ col.FieldName }}.IsValid() || m.{{ col.FieldName }}.IsNull() {
					return errors.New("the {{ col.FieldName }} of {{ tab.ModelName }} cannot be null")
				}
			{{ endif }}
		{{ endfor }}
		return nil
	}
{{ endfor }}`

const insertStmtTemp = `var {{Table.ModelName}}InsertStmt *sql.Stmt`

const updateStmtTemp = `var {{Table.ModelName}}UpdateStmt *sql.Stmt`

const deleteStmtTemp = `var {{Table.ModelName}}DeleteStmt *sql.Stmt`

const insertOrUpdateStmtTemp = `var {{Table.ModelName}}InsertOrUpdateStmt *sql.Stmt`

const insertMidStmtTemp = `var {{Table.ModelName}}To{{MTM.DstTab.ModelName}}InsertStmt *sql.Stmt`

const manyToManyDeleteStmtTemp = `var {{Table.ModelName}}To{{MTM.DstTab.ModelName}}DeleteStmt *sql.Stmt`

const insertStmtInitTemp = `{{Table.ModelName}}InsertStmt, err = {{Database.ObjName}}.Prepare("INSERT INTO {{Table.TableName}} ({{Columns}}) VALUES ({{Values}})")
if err != nil {
	log.Fatal(err)
	}`

const updateStmtInitTemp = `{{Table.ModelName}}UpdateStmt, err = {{Database.ObjName}}.Prepare("UPDATE {{Table.TableName}} SET {{Updates}} WHERE {{Table.AutoIncrement.ColumnName}} = ?")
if err != nil {
	log.Fatal(err)
	}`

const deleteStmtInitTemp = `{{Table.ModelName}}DeleteStmt, err = {{Database.ObjName}}.Prepare("DELETE FROM {{Table.TableName}} WHERE {{Table.AutoIncrement.ColumnName}} = ?")
if err != nil {
	log.Fatal(err)
	}`

const updateLastInsertIDTemp = `{{Table.AutoIncrement.ColumnName}} = LAST_INSERT_ID({{Table.AutoIncrement.ColumnName}})`

const insertOrUpdateStmtInitTemp = `{{Table.ModelName}}InsertOrUpdateStmt, err = {{Database.ObjName}}.Prepare("INSERT INTO {{Table.TableName}} ({{Columns}}) VALUES ({{Values}}) ON DUPLICATE KEY UPDATE {{UpdateLastInsertID}}, {{Updates}}")
if err != nil {
	log.Fatal(err)
	}`

const insertMidStmtInitTemp = `{{Table.ModelName}}To{{MTM.DstTab.ModelName}}InsertStmt, err = {{Database.ObjName}}.Prepare("INSERT INTO {{MTM.MidTab.TableName}} ({{MTM.MidLeftCol.ColumnName}}, {{MTM.MidRightCol.ColumnName}}) VALUES (?, ?)")
if err != nil {
	log.Fatal(err)
	}`

const manyToManyDeleteStmtInitTemp = `{{Table.ModelName}}To{{MTM.DstTab.ModelName}}DeleteStmt, err = {{Database.ObjName}}.Prepare("DELETE FROM {{MTM.MidTab.TableName}} WHERE {{MTM.MidLeftCol.ColumnName}} = ? AND {{MTM.MidRightCol.ColumnName}} = ?")`

//FuncArgNameTemp arguments name in function body template
const funcArgNameTemp = `{{Column.ArgName}}`

const updateColumnTemp = `{{Column.ColumnName}} = ?`

const modelListTypeTemp = `type modelList interface {
Len() int
Swap(i, j int)
}
type sortObj struct {
	list modelList
	lessFuncs []func(i, j int) int
}

func (so sortObj) Len() int {
	return so.list.Len()
}

func (so sortObj) Swap(i, j int) {
	so.list.Swap(i, j)
}

func (so sortObj) Less(i, j int) bool {
	for _, f := range so.lessFuncs {
		v := f(i, j)
		switch v {
		case -1:
			return true
		case 0:
			continue
		case 1:
			return false
		}
	}
	return false
}

type sortFunc func() func(i, j int) int 

{{ for _, tab in DB.Tables }}
	type {{ tab.ModelName }}List []*{{ tab.ModelName }}

	func (l {{ tab.ModelName }}List) Len() int {
		return len(l)
	}

	func (l {{ tab.ModelName }}List) Swap(i, j int) {
		l[i], l[j] = l[j], l[i]
	}

	{{ for _, col in tab.Columns }}
		func (l {{ tab.ModelName }}List) By{{ col.FieldName }}() func(i, j int) int {
			{{ switch col.FieldType }}
				{{ case "string" }}
					f := func(i, j int) int {
						ival, _, _ := l[i].Get{{ col.FieldName }}()
						jval, _, _ := l[j].Get{{ col.FieldName }}()
						switch {
						case ival<jval:
							return -1
						case ival > jval:
							return 1
						default:
							return 0
						}
					}
					return f
				{{ case "int64" }}
					f := func(i, j int) int {
						ival, _, _ := l[i].Get{{ col.FieldName }}()
						jval, _, _ := l[j].Get{{ col.FieldName }}()
						switch {
						case ival<jval:
							return -1
						case ival > jval:
							return 1
						default:
							return 0
						}
					}
					return f
				{{ case "float64" }}
					f := func(i, j int) int {
						ival, _, _ := l[i].Get{{ col.FieldName }}()
						jval, _, _ := l[j].Get{{ col.FieldName }}()
						switch {
						case ival<jval:
							return -1
						case ival > jval:
							return 1
						default:
							return 0
						}
					}
					return f
				{{ case "bool" }}
					f := func(i, j int) int {
						ival, _, _ := l[i].Get{{ col.FieldName }}()
						jval, _, _ := l[j].Get{{ col.FieldName }}()
						var ii, ji int
						if ival {
							ii = 1
						}
						if jval {
							ji = 1
						}
						switch {
						case ii < ji:
							return -1
						case ii > ji:
							return 1
						default:
							return 0
						}
					}
					return f
				{{ case "time.Time" }}
					f := func(i, j int) int {
						ival, _, _ := l[i].Get{{ col.FieldName }}()
						jval, _, _ := l[j].Get{{ col.FieldName }}()
						switch {
						case ival.Before(jval):
							return -1
						case ival.After(jval):
							return 1
						default:
							return 0
						}
					}
					return f
			{{ endswitch }}
		}
	{{ endfor }}
	func (l {{ tab.ModelName }}List) Sort(reverse bool, funcs ...sortFunc) {
		if len(funcs) == 0 {
			return
		}
		so := sortObj{list: l, lessFuncs: make([]func(i, j int) int, len(funcs))}
		for i := range funcs {
			so.lessFuncs[i] = funcs[i]()
		}
		if reverse {
			sort.Sort(sort.Reverse(so))
		} else {
			sort.Sort(so)
		}
	}
{{ endfor }}`

const modelCompareByIntMethodTemp = `By{{Column.FieldName}}: func(i, j int) int {
	if *ml.Models[i].{{Column.FieldName}} == *ml.Models[j].{{Column.FieldName}} {
		return 0
	}
	if *ml.Models[i].{{Column.FieldName}} < *ml.Models[j].{{Column.FieldName}} {
		return -1
	}
	return 1
	},`

const modelCompareByFloatMethodTemp = `By{{Column.FieldName}}: func(i, j int) int {
	if *ml.Models[i].{{Column.FieldName}} == *ml.Models[j].{{Column.FieldName}} {
		return 0
	}
	if *ml.Models[i].{{Column.FieldName}} < *ml.Models[j].{{Column.FieldName}} {
		return -1
	}
	return 1
	},`

const modelCompareByStringMethodTemp = `By{{Column.FieldName}}: func(i, j int) int {
	if *ml.Models[i].{{Column.FieldName}} == *ml.Models[j].{{Column.FieldName}} {
		return 0
	}
	if *ml.Models[i].{{Column.FieldName}} < *ml.Models[j].{{Column.FieldName}} {
		return -1
	}
	return 1
	},`

const modelCompareByBoolMethodTemp = `By{{Column.FieldName}}: func(i, j int) int {
	if *ml.Models[i].{{Column.FieldName}} == *ml.Models[j].{{Column.FieldName}} {
		return 0
	}
	if *ml.Models[i].{{Column.FieldName}} == false && *ml.Models[j].{{Column.FieldName}} == true {
		return -1
	}
	return 1
	},`

const modelCompareByTimeMethodTemp = `By{{Column.FieldName}}: func(i, j int) int {
	if ml.Models[i].{{Column.FieldName}}.Equal(*ml.Models[j].{{Column.FieldName}}) {
		return 0
	}
	if ml.Models[i].{{Column.FieldName}}.Before(*ml.Models[j].{{Column.FieldName}}) {
		return -1
	}
	return 1
	},`

const modelSortMethodsStructFieldTypeTemp = `By{{Column.FieldName}} func(i, j int) int`

const modelSortMethodsStructTypeTemp = `type {{Table.ModelName}}SortMethods struct {
	{{FieldTypeBlock}}
	}`

const modelSortMethodsStructFuncTemp = `func (ml {{Table.ModelName}}List) SortMethods() {{Table.ModelName}}SortMethods {
	return {{Table.ModelName}}SortMethods{
		{{CompareMethodBlock}}
		}
	}`

const modelListLenMethodTemp = `func (ml {{Table.ModelName}}List) Len() int {
	return len(ml.Models)
	}`

const modelListSwapMethodTemp = `func (ml {{Table.ModelName}}List) Swap(i, j int) {
	ml.Models[i], ml.Models[j] = ml.Models[j], ml.Models[i]
	}`

const modelListLessMethodTemp = `func (ml {{Table.ModelName}}List) Less(i, j int) bool {
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

const modelSortFuncSwitchBlockTemp = `case "{{Column.FieldName}}": 
	{{Table.ArgName}}List.Funcs = append({{Table.ArgName}}List.Funcs, sortMethods.By{{Column.FieldName}})`

const modelSortFuncTemp = `func {{Table.ModelName}}SortBy(ml []*{{Table.ModelName}}, desc bool, fields ...string) {
	{{Table.ArgName}}List := {{Table.ModelName}}List {
		Models: ml,
		Funcs: make([]func(i, j int) int, 0, {{ColumnNumber}}),
		}
	sortMethods := {{Table.ArgName}}List.SortMethods()
	for _, field := range fields {
		switch field {
			{{SwitchBlock}}
		}
	}
	if desc {
		sort.Sort(sort.Reverse({{Table.ArgName}}List))
	} else {
		sort.Sort({{Table.ArgName}}List)
	}
}`

const countFuncTemp = `{{ for _, tab in DB.Tables }}
	func Count{{ tab.ModelName }}(where ...string) (int64, error) {
		var count int64
		var row *sql.Row
		if len(where) == 0 {
			row = {{ DB.ObjName }}.QueryRow("SELECT COUNT(*) FROM {{ tab.TableName }}")
		} else {
			row = {{ DB.ObjName }}.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM {{ tab.TableName }} WHERE %s", where[0]))
		}
		err := row.Scan(&count)
		if err != nil {
			return -1, err
		}
		return count, nil
	}
{{ endfor }}`

const fieldTypeTemp = `type FieldArg interface {
	columnName() string
	tableName() string
	sqlValue() string
}
{{ for _, tab in DB.Tables }}
	{{ for _, col in tab.Columns }}
		{{ switch col.FieldType }}
			{{ case "int64" }}
				type {{ tab.ModelName }}{{ col.FieldName }} int64
				func New{{ tab.ModelName }}{{ col.FieldName }}(val int64) *{{ tab.ModelName }}{{ col.FieldName }} {
					f := {{ tab.ModelName }}{{ col.FieldName }}(val)
					return &f
				}
				func (f *{{ tab.ModelName }}{{ col.FieldName }}) columnName() string {
					return "{{ col.ColumnName }}"
				}
				func (f *{{ tab.ModelName }}{{ col.FieldName }}) sqlValue() string {
					return fmt.Sprintf("%d", *f)
				}
			{{ case "string" }}
				type {{ tab.ModelName }}{{ col.FieldName }} string
				func New{{ tab.ModelName }}{{ col.FieldName}}(val string) *{{ tab.ModelName }}{{ col.FieldName}} {
					f := {{ tab.ModelName }}{{ col.FieldName }}(val)
					return &f
				}
				func (f *{{ tab.ModelName }}{{ col.FieldName }}) columnName() string {
					return "{{ col.ColumnName }}"
				}
				func (f *{{ tab.ModelName }}{{ col.FieldName }}) sqlValue() string {
					return fmt.Sprintf("%q", *f)
				}
			{{ case "float64" }}
				type {{ table.ModelName }}{{ col.FieldName }} float64
				func New{{ tab.ModelName }}{{ col.FieldName }}(val float64) *{{ tab.ModelName }}{{ col.FieldName }} {
					f := {{ tab.ModelName }}{{ col.FieldName }}(val)
					return &f
				}
				func (f *{{ tab.ModelName }}{{ col.FieldName }}) columnName() string {
					return "{{ col.ColumnName }}"
				}
				func (f *{{ tab.ModelName }}{{ col.FieldName }}) sqlValue() string {
					return fmt.Sprintf("%f", *f)
				}
			{{ case "bool" }}
				type {{ tab.ModelName }}{{ col.FieldName }} bool
				func New{{ tab.ModelName }}{{ col.FieldName }}(val bool) *{{ tab.ModelName }}{{ col.FieldName }} {
					f := {{ tab.ModelName }}{{ col.FieldName }}(val)
					return &f
				}
				func (f *{{ tab.ModelName }}{{ col.FieldName }}) columnName() string {
					return "{{ col.ColumnName }}"
				}
				func (f *{{ tab.ModelName }}{{ col.FieldName }}) sqlValue() string {
					return fmt.Sprintf("%t", *f)
				}
			{{ case "time.Time" }}
				type {{ tab.ModelName }}{{ col.FieldName }} time.Time
				func New{{ tab.ModelName }}{{ col.FieldName }}(val time.Time) *{{ tab.ModelName }}{{ col.FieldName }}{
					f := {{ tab.ModelName }}{{ col.FieldName }}(val)
					return &f
				}
				func (f *{{ tab.ModelName }}{{ col.FieldName }}) columnName() string {
					return "{{ col.ColumnName }}"
				}
				func (f *{{ tab.ModelName }}{{ col.FieldName }}) sqlValue() string {
					{{ switch col.MySqlType }}
						{{ case "DATE" }}
							return time.Time(*f).Format("2006-01-02")
						{{ case "DATETIME" }}
							return time.Time(*f).Format("2006-01-02 15:04:05")
						{{ case "TIMESTAMP" }}
							return time.Time(*f).Format("2006-01-02 15:04:05")
					{{ endswitch }}
				}
		{{ endswitch }}
		func (f *{{ tab.ModelName}}{{ col.FieldName}}) tableName() string {
			return "{{ tab.TableName }}"
		}
	{{ endfor }}
{{ endfor }}
`

const modelInsertFuncTemp = `{{ for _, tab in DB.Tables }}
	func {{ tab.ModelName }}Insert(fields ...FieldArg) (*{{ tab.ModelName }}, error) {
		m, err := New{{ tab.ModelName }}(fields...)
		if err != nil {
			return nil, err
		}
		colList := make([]string, len(fields))
		valList := make([]string, len(fields))
		for i, f := range fields {
			colList[i] = f.columnName()
			valList[i] = f.sqlValue()
		}
		stmtStr := fmt.Sprintf("INSERT INTO {{ tab.TableName }} (%s) VALUES (%s)", strings.Join(colList, ", "), strings.Join(valList, ", "))
		res, err := {{ DB.ObjName }}.Exec(stmtStr)
		if err != nil {
			return nil, err
			}
		lastInsertId, err := res.LastInsertId()
		if err != nil {
			return nil, err
		}
		m.{{ tab.AutoIncrement.FieldName }}.Set(lastInsertId)
		return m, nil
		}
{{ endfor }}`

const newModelFuncTemp = `{{ for _, tab in DB.Tables }}
func New{{ tab.ModelName }}(fields ...FieldArg) (*{{ tab.ModelName }}, error) {
	m := &{{ tab.ModelName }}{}
	for _, field := range fields {
		switch f := field.(type) {
		{{ for _, col in tab.Columns }}
		case *{{ tab.ModelName }}{{ col.FieldName }}:
			if f == nil {
				{{ switch col.FieldType }}
					{{ case "string" }}
						m.{{ col.FieldName }}.Set("", true)
					{{ case "int64" }}
						m.{{ col.FieldName }}.Set(0, true)
					{{ case "float64" }}
						m.{{ col.FieldName }}.Set(0.0, true)
					{{ case "bool" }}
						m.{{ col.FieldName }}.Set(false, false)
					{{ case "time.Time" }}
						m.{{ col.FieldName }}.Set(time.Time{}, false)
				{{ endswitch }}
			}
			m.{{ col.FieldName }}.Set({{ col.FieldType }}(*f))
		{{ endfor }}
		default:
			return nil, errors.New("invalid field in New{{ tab.ModelName }}()")
		}
	}
	return m, nil
}
{{ endfor }}`

const modelTypeTemp = `{{ for _, tab in DB.Tables }}
	type {{ tab.ModelName }} struct {
		{{ for _, col in tab.Columns }}
			{{ switch col.FieldType }}
				{{ case "string" }}
					{{ col.FieldName }} nbmysql.StringField
				{{ case "int64" }}
					{{ col.FieldName }} nbmysql.IntField
				{{ case "float64" }}
					{{ col.FieldName }} nbmysql.FloatField
				{{ case "bool" }}
					{{ col.FieldName }} nbmysql.BoolField
				{{ case "time.Time" }}
					{{ switch col.MySqlType }}
						{{ case "DATE" }}
							{{ col.FieldName }} nbmysql.DateField
						{{ case "DATETIME" }}
							{{ col.FieldName }} nbmysql.DatetimeField
						{{ case "TIMESTAMP" }}
							{{ col.FieldName }} nbmysql.DatetimeField
					{{ endswitch }}
			{{ endswitch }}
		{{ endfor }}
	}
	{{ for _, col in tab.Columns }}
		func (m *{{ tab.ModelName }})Get{{ col.FieldName }}() (val {{ col.FieldType }}, valid bool, null bool) {
			return m.{{ col.FieldName }}.Get()
		}
		func (m *{{ tab.ModelName }})Set{{ col.FieldName }}(val {{ col.FieldType }}, nullAndValid ...bool) {
			m.{{ col.FieldName }}.Set(val, nullAndValid...)
		}
	{{ endfor }}
{{ endfor }}
`
const updateFuncTemp = `{{ for _, tab in DB.Tables }}
	func Update{{ tab.ModelName }}(where string, update ...FieldArg) error {
		if len(update) == 0 {
			return errors.New("Update{{ tab.ModelName }}() error: update list cannot be empty")
		}
		for k, v := range {{ tab.ModelName }}Map {
			where = strings.Replace(where, k, v, -1)
		}
		setList := make([]string, len(update))
		for i, f := range update {
			if f.tableName() != "{{ tab.TableName }}" {
				return fmt.Errorf("Update{{ tab.ModelName }}() error: %T not belong to {{ tab.ModelName }}", f)
			}
			if f == nil {
				setList[i] = f.columnName() + "=" + "NULL"
			} else {
				setList[i] = f.columnName() + "=" + f.sqlValue()
			}
		}
		stmtStr := fmt.Sprintf("UPDATE {{ tab.TableName }} SET %s WHERE %s", strings.Join(setList, ", "), where)
		_, err := {{ DB.ObjName }}.Exec(stmtStr)
		return err
	}
{{ endfor }}`

const modelInvalidateMethodTemp = `{{ for _, tab in DB.Tables }}
	func (m *{{ tab.ModelName }}) invalidate() {
		{{ for _, col in tab.Columns }}
			m.{{ col.FieldName }}.Invalidate()
		{{ endfor }}
	}
{{ endfor }}`

const modelDeleteMethodTemp = `{{ for _, tab in DB.Tables }}
	func (m *{{ tab.ModelName }}) Delete() error {
		if !m.{{ tab.PrimaryKey.FieldName }}.IsValid() || m.{{ tab.PrimaryKey.FieldName }}.IsNull() {
			return errors.New("{{ tab.ModelName }}.Delete() error: primary key is not valid")
		}
		stmtStr := fmt.Sprintf("DELETE FROM {{ tab.TableName }} WHERE {{ tab.PrimaryKey.ColumnName }} = %s", m.{{ tab.PrimaryKey.FieldName }}.SQLVal())
		res, err := {{ DB.ObjName }}.Exec(stmtStr)
		if err != nil {
			return err
		}
		affectedRows, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if affectedRows == 0 {
			return fmt.Errorf("{{ tab.ModelName }}.Delete() error: row not exists in {{ tab.TableName }}, primary Key(%s)", m.{{ tab.PrimaryKey.FieldName }}.SQLVal())
		}
		m.invalidate()
		return nil
	}
{{ endfor }}`

const deleteFuncTemp = `{{ for _, tab in DB.Tables }}
	func Delete{{ tab.ModelName }}(where string) (int64, error) {
		for k, v := range {{ tab.ModelName }}Map {
			where = strings.Replace(where, k, v, -1)
		}
		stmtStr := fmt.Sprintf("DELETE FROM {{ tab.TableName }} WHERE %s", where)
		res, err := {{ DB.ObjName }}.Exec(stmtStr)
		if err != nil {
			return -1, err
		}
		return res.RowsAffected()
	}
{{ endfor }}`
