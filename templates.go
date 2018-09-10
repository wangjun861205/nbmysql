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

const insertStmtMethodTemp = `{{ for _, tab in DB.Tables }}
func (m *{{ tab.ModelName }})InsertStmt() *Stmt {
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
	return NewStmt(m, stmtStr)
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
{{ endfor }}`

const updateStmtMethodTemp = `{{ for _, tab in DB.Tables }}
	func (m *{{ tab.ModelName }}) UpdateStmt() (*Stmt, error) {
		colList := make([]string, 0, 16)
		valList := make([]string, 0, 16)
		if !m.{{ tab.PrimaryKey.FieldName }}.IsValid() || m.{{ tab.PrimaryKey.FieldName }}.IsNull() {
			return nil, errors.New("{{ tab.ModelName }}.Update(): primary key cannot be invalid or null")
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
		return NewStmt(m, stmtStr), nil
	}
{{ endfor }}`

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

// const modelListTypeTemp = `type modelList interface {
// Len() int
// Swap(i, j int)
// }
// type sortObj struct {
// 	list modelList
// 	lessFuncs []func(i, j int) int
// }

// func (so sortObj) Len() int {
// 	return so.list.Len()
// }

// func (so sortObj) Swap(i, j int) {
// 	so.list.Swap(i, j)
// }

// func (so sortObj) Less(i, j int) bool {
// 	for _, f := range so.lessFuncs {
// 		v := f(i, j)
// 		switch v {
// 		case -1:
// 			return true
// 		case 0:
// 			continue
// 		case 1:
// 			return false
// 		}
// 	}
// 	return false
// }

// type sortFunc func() func(i, j int) int

// {{ for _, tab in DB.Tables }}
// 	type {{ tab.ModelName }}List []*{{ tab.ModelName }}

// 	func (l {{ tab.ModelName }}List) Len() int {
// 		return len(l)
// 	}

// 	func (l {{ tab.ModelName }}List) Swap(i, j int) {
// 		l[i], l[j] = l[j], l[i]
// 	}

// 	{{ for _, col in tab.Columns }}
// 		func (l {{ tab.ModelName }}List) By{{ col.FieldName }}() func(i, j int) int {
// 			{{ switch col.FieldType }}
// 				{{ case "string" }}
// 					f := func(i, j int) int {
// 						ival, _, _ := l[i].Get{{ col.FieldName }}()
// 						jval, _, _ := l[j].Get{{ col.FieldName }}()
// 						switch {
// 						case ival<jval:
// 							return -1
// 						case ival > jval:
// 							return 1
// 						default:
// 							return 0
// 						}
// 					}
// 					return f
// 				{{ case "int64" }}
// 					f := func(i, j int) int {
// 						ival, _, _ := l[i].Get{{ col.FieldName }}()
// 						jval, _, _ := l[j].Get{{ col.FieldName }}()
// 						switch {
// 						case ival<jval:
// 							return -1
// 						case ival > jval:
// 							return 1
// 						default:
// 							return 0
// 						}
// 					}
// 					return f
// 				{{ case "float64" }}
// 					f := func(i, j int) int {
// 						ival, _, _ := l[i].Get{{ col.FieldName }}()
// 						jval, _, _ := l[j].Get{{ col.FieldName }}()
// 						switch {
// 						case ival<jval:
// 							return -1
// 						case ival > jval:
// 							return 1
// 						default:
// 							return 0
// 						}
// 					}
// 					return f
// 				{{ case "bool" }}
// 					f := func(i, j int) int {
// 						ival, _, _ := l[i].Get{{ col.FieldName }}()
// 						jval, _, _ := l[j].Get{{ col.FieldName }}()
// 						var ii, ji int
// 						if ival {
// 							ii = 1
// 						}
// 						if jval {
// 							ji = 1
// 						}
// 						switch {
// 						case ii < ji:
// 							return -1
// 						case ii > ji:
// 							return 1
// 						default:
// 							return 0
// 						}
// 					}
// 					return f
// 				{{ case "time.Time" }}
// 					f := func(i, j int) int {
// 						ival, _, _ := l[i].Get{{ col.FieldName }}()
// 						jval, _, _ := l[j].Get{{ col.FieldName }}()
// 						switch {
// 						case ival.Before(jval):
// 							return -1
// 						case ival.After(jval):
// 							return 1
// 						default:
// 							return 0
// 						}
// 					}
// 					return f
// 			{{ endswitch }}
// 		}
// 	{{ endfor }}
// 	func (l {{ tab.ModelName }}List) Sort(reverse bool, funcs ...sortFunc) {
// 		if len(funcs) == 0 {
// 			return
// 		}
// 		so := sortObj{list: l, lessFuncs: make([]func(i, j int) int, len(funcs))}
// 		for i := range funcs {
// 			so.lessFuncs[i] = funcs[i]()
// 		}
// 		if reverse {
// 			sort.Sort(sort.Reverse(so))
// 		} else {
// 			sort.Sort(so)
// 		}
// 	}
// 	func (l {{ tab.ModelName }}List) Filter(f func(*{{ tab.ModelName }}) bool) {{ tab.ModelName }}List {
// 		fl := make({{ tab.ModelName }}List, 0, l.Len())
// 		for _, m := range l {
// 			if f(m) {
// 				fl = append(fl, m)
// 			}
// 		}
// 		return fl
// 	}
// {{ endfor }}`

const modelListTypeTemp = `{{ for _, tab in DB.Tables }}
	type {{ tab.ModelName }}List []*{{ tab.ModelName }}

	func (l {{ tab.ModelName }}List) Len() int {
		return len(l)
	}

	func (l {{ tab.ModelName }}List) Swap(i, j int) {
		l[i], l[j] = l[j], l[i]
	}

	type {{ tab.ArgName }}SortObj struct {
		list {{ tab.ModelName }}List
		lessFuncs []func(l {{ tab.ModelName }}List, i, j int) int
	}

	func (so {{ tab.ArgName }}SortObj) Len() int {
		return so.list.Len()
	}

	func (so {{ tab.ArgName }}SortObj) Swap(i, j int) {
		so.list.Swap(i, j)
	}

	func (so {{ tab.ArgName }}SortObj) Less(i, j int) bool {
		for _, f := range so.lessFuncs {
			v := f(so.list, i, j)
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
{{ endfor }}`

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
	func Insert{{ tab.ModelName }}(fields ...FieldArg) (*{{ tab.ModelName }}, error) {
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

const insertStmtFuncTemp = `{{ for _, tab in DB.Tables }}
	func Insert{{ tab.ModelName }}Stmt(fields ...FieldArg) *Stmt {
		colList := make([]string, len(fields))
		valList := make([]string, len(fields))
		for i, f := range fields {
			colList[i] = f.columnName()
			valList[i] = f.sqlValue()
		}
		stmtStr := fmt.Sprintf("INSERT INTO {{ tab.TableName }} (%s) VALUES (%s)", strings.Join(colList, ", "), strings.Join(valList, ", "))
		return NewStmt(nil, stmtStr)
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

const updateStmtFuncTemp = `{{ for _, tab in DB.Tables }}
	func Update{{ tab.ModelName }}Stmt(where string, update ...FieldArg) (*Stmt, error) {
		if len(update) == 0 {
			return nil, errors.New("Update{{ tab.ModelName }}() error: update list cannot be empty")
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
		return NewStmt(nil, stmtStr), nil
	}
{{ endfor }}`

const modelInvalidateMethodTemp = `{{ for _, tab in DB.Tables }}
	func (m *{{ tab.ModelName }}) Invalidate() {
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
		m.Invalidate()
		return nil
	}
{{ endfor }}`

const deleteStmtMethodTemp = `{{ for _, tab in DB.Tables }}
	func (m *{{ tab.ModelName }}) DeleteStmt() (*Stmt, error) {
		if !m.{{ tab.PrimaryKey.FieldName }}.IsValid() || m.{{ tab.PrimaryKey.FieldName }}.IsNull() {
			return nil, errors.New("{{ tab.ModelName }}.Delete() error: primary key is not valid")
		}
		stmtStr := fmt.Sprintf("DELETE FROM {{ tab.TableName }} WHERE {{ tab.PrimaryKey.ColumnName }} = %s", m.{{ tab.PrimaryKey.FieldName }}.SQLVal())
		return NewStmt(m, stmtStr), nil
	}
{{ endfor }}`

const deleteFuncTemp = `{{ for _, tab in DB.Tables }}
	func Delete{{ tab.ModelName }}(where string) (int64, error) {
		for k, v := range {{ tab.ModelName }}Map {
			where = strings.Replace(where, k, v, -1)
		}
		var stmtStr string
		if where != "" {
			stmtStr = fmt.Sprintf("DELETE FROM {{ tab.TableName }} WHERE %s", where)
		} else {
			stmtStr = fmt.Sprintf("DELETE FROM {{ tab.TableName }}") 
		}
		res, err := {{ DB.ObjName }}.Exec(stmtStr)
		if err != nil {
			return -1, err
		}
		return res.RowsAffected()
	}
{{ endfor }}`

const deleteStmtFuncTemp = `{{ for _, tab in DB.Tables }}
	func Delete{{ tab.ModelName }}Stmt(where string) *Stmt {
		for k, v := range {{ tab.ModelName }}Map {
			where = strings.Replace(where, k, v, -1)
		}
		var stmtStr string
		if where != "" {
			stmtStr = fmt.Sprintf("DELETE FROM {{ tab.TableName }} WHERE %s", where)
		} else {
			stmtStr = fmt.Sprintf("DELETE FROM {{ tab.TableName }}") 
		}
		return NewStmt(nil, stmtStr)
	}
{{ endfor }}`

const modelInterfaceTypeTemp = `type model interface {
	Invalidate()
	SetLastInsertID(int64)
}`

const setLastInsertIDMethodTemp = `{{ for _, tab in DB.Tables }}
	func (m *{{ tab.ModelName }}) SetLastInsertID(id int64) {
		m.{{ tab.AutoIncrement.FieldName }}.Set(id)
	}
{{ endfor }}`

const stmtTypeTemp = `type stmtType int
const (
	insert stmtType = iota
	update
	delete
	other
)

type Stmt struct {
	model model
	typ stmtType
	stmt string
	lastInsertID int64
}

func NewStmt(model model, stmtStr string) *Stmt {
	stmt := &Stmt{model: model, stmt: stmtStr}
	l := strings.Split(strings.ToUpper(stmtStr), " ")
	switch l[0] {
	case "INSERT":
		stmt.typ = insert
	case "UPDATE":
		stmt.typ = update
	case "DELETE":
		stmt.typ = delete
	default:
		stmt.typ = other
	}
	return stmt
}

func (s *Stmt) Join(stmts ...*Stmt) StmtList {
	l := make(StmtList, len(stmts) + 1)
	l[0] = s
	for i := 0; i < len(stmts); i++ {
		l[i+1] = stmts[i]
	}
	return l
}

func (s *Stmt) Exec() error {
	res, err := {{ DB.ObjName }}.Exec(s.stmt)
	if err != nil {
		return err
	}
	if s.model != nil {
		switch s.typ {
		case insert:
			lastInsertID, err := res.LastInsertId()
			if err != nil {
				return err
			}
			s.model.SetLastInsertID(lastInsertID)
		case delete:
			s.model.Invalidate()
		}
	}
	return nil
}

type StmtList []*Stmt

func (sl StmtList) Exec() error {
	tx, err := {{ DB.ObjName }}.Begin()
	if err != nil {
		return err
	}
	for _, s := range sl {
		res, err := tx.Exec(s.stmt)
		if err != nil {
			tx.Rollback()
			return err
		}
		if s.typ == insert && s.model != nil {
			lastInsertID, err := res.LastInsertId()
			if err != nil {
				tx.Rollback()
				return err
			}
			s.lastInsertID = lastInsertID
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	for _, s := range sl {
		if s.model != nil {
			switch s.typ {
			case insert:
				s.model.SetLastInsertID(s.lastInsertID)
			case delete:
				s.model.Invalidate()
			}
		}
	}
	return nil
}

func (sl StmtList) Join(sls ...StmtList) StmtList {
	for _, l := range sls {
		sl = append(sl, l...)
	}
	return sl
	}`

const modelDistinctMethodTemp = `{{ for _, tab in DB.Tables }}
	{{ for _, col in tab.Columns }}
		func (m *{{ tab.ModelName }}) DistBy{{ col.FieldName }}() string {
			return m.{{ col.FieldName }}.SQLVal()
		}
	{{ endfor }}
{{ endfor }}`

const modelListDistinctMethodTemp = `{{ for _, tab in DB.Tables }}
	func (l {{ tab.ModelName }}List) Distinct(fs ...func(*{{ tab.ModelName }}) string) {{ tab.ModelName }}List {
		if len(fs) == 0 {
			fs = append(fs, func(mm *{{ tab.ModelName }}) string {
				return mm.{{ tab.PrimaryKey.FieldName }}.SQLVal()
			})
		}
		m := make(map[string]bool)
		filteredList := make({{ tab.ModelName }}List, 0, l.Len())
		builder := strings.Builder{}
		for _, model := range l {
			for _, f := range fs {
				builder.WriteString(f(model))
			}
			if exists := m[builder.String()]; exists {
				builder.Reset()
				continue
			} else {
				m[builder.String()] = true
				filteredList = append(filteredList, model)
				builder.Reset()
			}
		}
		return filteredList
	}
{{ endfor }}`

const modelListSortByMethodTemp = `{{ for _, tab in DB.Tables }}
	{{ for _, col in tab.Columns }}
		func (l {{ tab.ModelName }}List) SortBy{{ col.FieldName }}(i, j int) int {
			var iInt, jInt int
			iVal, iValid, iNull := l[i].{{ col.FieldName }}.Get()
			jVal, jValid, jNull := l[j].{{ col.FieldName }}.Get()
			if !iValid || !jValid {
				if iValid {
					iInt = 1
				}
				if jValid {
					jInt = 1
				}
				return iInt - jInt
			}
			if iNull || jNull {
				if !iNull {
					iInt = 1
				}
				if !jNull {
					jInt = 1
				}
				return iInt - jInt
			}
			{{ switch col.FieldType }}
				{{ case "string", "int64", "float64" }}
					switch {
					case iVal < jVal:
						return -1
					case iVal > jVal:
						return 1
					default:
						return 0
					}
				{{ case "bool" }}
					if iVal {
						iInt = 1
					}
					if jVal {
						jInt = 1
					}
					return iInt - jInt
				{{ case "time.Time" }}
					switch {
					case iVal.Before(jVal):
						return -1
					case iVal.After(jVal):
						return 1
					default:
						return 0
					}
			{{ endswitch }}
		}
	{{ endfor }}
{{ endfor }}`

const modelListSortMethodTemp = `{{ for _, tab in DB.Tables }}
	func (l {{ tab.ModelName }}List) Sort(reverse bool, fs ...func({{ tab.ModelName }}List, int, int) int) {
		so := {{ tab.ArgName }}SortObj{list: l, lessFuncs: fs}
		if reverse {
			sort.Sort(sort.Reverse(so))
		} else {
			sort.Sort(so)
		}
	}
{{ endfor }}`
