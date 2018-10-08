package nbmysql

//PackageTemp package template
var packageTemp = `package {{ DB.Package }}
`

//ImportTemp import block template
const importTemp = `
import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sort"
	"github.com/wangjun861205/nbmysql"
)
`

//dbTemp *sql.DB declaration template
const dbTemp = `
var {{ DB.ObjName }} *sql.DB
`

//initFuncTemp init func template
const initFuncTemp = `
func init() {
	db, err := sql.Open("mysql", "{{ DB.Username }}:{{ DB.Password }}@tcp({{ DB.Address }})/{{ DB.DatabaseName }}")
	if err != nil {
		panic(err)
	}
	{{ DB.ObjName }} = db
}
`

const fromRowFuncTemp = `
{{ for _, tab in DB.Tables }}
	func {{ tab.ModelName }}FromRow(row *sql.Row) (*{{ tab.ModelName}}Instance, error) {
		m := {{ tab.ModelName }}Instance{}
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
{{ endfor }}
`

const fromRowsFuncTemp = `
{{ for _, tab in DB.Tables }}
	func {{ tab.ModelName }}FromRows(rows *sql.Rows) ({{ tab.ModelName }}List, error) {
		l := make({{ tab.ModelName }}List, 0, 128)
		for rows.Next() {
			m := {{ tab.ModelName }}Instance{}
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
{{ endfor }}
`

const foreignKeyMethodTemp = `
{{ for _, tab in DB.Tables }}
	{{ for _, fk in tab.ForeignKeys }}
		func (m *{{ tab.ModelName }}Instance) {{ fk.DstTab.ModelName }}() (*{{ fk.DstTab.ModelName }}Instance, error) {
			isValid, isNull := m.{{ fk.SrcCol.FieldName }}.IsValid(), m.{{ fk.SrcCol.FieldName }}.IsNull()
			if !isValid || isNull {
				return nil, errors.New("{{ tab.ModelName }}.{{ fk.SrcCol.FieldName }} cannot be null")
			}
			srcColVal, _, _ := m.{{ fk.SrcCol.FieldName }}.Get()
			return {{ fk.DstTab.ModelName }}.QueryOne({{ fk.DstTab.ModelName }}.{{ fk.DstCol.FieldName }}.Eq(srcColVal))
		}
	{{ endfor }}
{{ endfor }}
`

const reverseForeignKeyMethodTemp = `
{{ for _, tab in DB.Tables }}
	{{ for _, rfk in tab.ReverseForeignKeys }}
		func (m *{{ tab.ModelName }}Instance) {{ rfk.DstTab.ModelName }}List() struct {
			All func() ({{ rfk.DstTab.ModelName }}List, error)
			Query func(*nbmysql.Where) ({{ rfk.DstTab.ModelName }}List, error)
		} {
			return struct {
				All func() ({{ rfk.DstTab.ModelName }}List, error)
				Query func(*nbmysql.Where) ({{ rfk.DstTab.ModelName }}List, error)
				} {
					All: func() ({{ rfk.DstTab.ModelName }}List, error) {
						isValid, isNull := m.{{ rfk.SrcCol.FieldName }}.IsValid(), m.{{ rfk.SrcCol.FieldName }}.IsNull()
						if !isValid || isNull {
							return nil, errors.New("{{ tab.ModelName }}.{{ rfk.SrcCol.FieldName }} cannot be null")
						}
						srcColVal, _, _ := m.{{ rfk.SrcCol.FieldName }}.Get()
						return {{ rfk.DstTab.ModelName }}.Query({{ rfk.DstTab.ModelName }}.{{ rfk.DstCol.FieldName }}.Eq(srcColVal))
						},
					Query: func(where *nbmysql.Where) ({{ rfk.DstTab.ModelName }}List, error) {
						isValid, isNull := m.{{ rfk.SrcCol.FieldName }}.IsValid(), m.{{ rfk.SrcCol.FieldName }}.IsNull()
						if !isValid || isNull {
							return nil, errors.New("{{ tab.ModelName }}.{{ rfk.SrcCol.FieldName }} cannot be null")
						}
						srcColVal, _, _ := m.{{ rfk.SrcCol.FieldName }}.Get()
						return {{ rfk.DstTab.ModelName }}.Query({{ rfk.DstTab.ModelName }}.{{ rfk.DstCol.FieldName }}.Eq(srcColVal).And(where))
						},
					}
			}
	{{ endfor }}
{{ endfor }}
`

const manyToManyMethodTemp = `
{{ for _, tab in DB.Tables }} 
	{{ for _, mtm in tab.ManyToManys }}
		func (m *{{ tab.ModelName }}Instance) {{ mtm.DstTab.ModelName }}List() struct {
			All func() ({{ mtm.DstTab.ModelName }}List, error)
			Query func(*nbmysql.Where) ({{ mtm.DstTab.ModelName }}List, error)
			Add func(*{{ mtm.DstTab.ModelName }}Instance) error
			Remove func(*{{ mtm.DstTab.ModelName }}Instance) error
			} {
				return struct {
					All func() ({{ mtm.DstTab.ModelName }}List, error)
					Query func(*nbmysql.Where) ({{ mtm.DstTab.ModelName }}List, error)
					Add func(*{{ mtm.DstTab.ModelName }}Instance) error
					Remove func(*{{ mtm.DstTab.ModelName }}Instance) error
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
						Query: func(where *nbmysql.Where) ({{ mtm.DstTab.ModelName }}List, error) {
							isValid, isNull := m.{{ mtm.SrcCol.FieldName }}.IsValid(), m.{{ mtm.SrcCol.FieldName }}.IsNull()
							if !isValid || isNull {
								return nil, errors.New("{{ tab.ModelName }}.{{ mtm.SrcCol.FieldName }} cannot be null")
							}
							var stmtStr string
							if where == nil {
								stmtStr = fmt.Sprintf("SELECT {{ mtm.DstTab.TableName }}.* FROM {{ tab.TableName }} JOIN {{ mtm.MidTab.TableName }} ON {{ tab.TableName }}.{{ mtm.SrcCol.ColumnName }} = {{ mtm.MidTab.TableName }}.{{ mtm.MidLeftCol.ColumnName }} JOIN {{ mtm.DstTab.TableName }} ON {{ mtm.MidTab.TableName }}.{{ mtm.MidRightCol.ColumnName }} = {{ mtm.DstTab.TableName}}.{{ mtm.DstCol.ColumnName }} WHERE {{ tab.TableName }}.{{ mtm.SrcCol.ColumnName }} = %s", m.{{ mtm.SrcCol.FieldName }}.SQLVal())
							} else {
								stmtStr = fmt.Sprintf("SELECT {{ mtm.DstTab.TableName }}.* FROM {{ tab.TableName }} JOIN {{ mtm.MidTab.TableName }} ON {{ tab.TableName }}.{{ mtm.SrcCol.ColumnName }} = {{ mtm.MidTab.TableName }}.{{ mtm.MidLeftCol.ColumnName }} JOIN {{ mtm.DstTab.TableName }} ON {{ mtm.MidTab.TableName }}.{{ mtm.MidRightCol.ColumnName }} = {{ mtm.DstTab.TableName}}.{{ mtm.DstCol.ColumnName }} WHERE {{ tab.TableName }}.{{ mtm.SrcCol.ColumnName }} = %s AND %s", m.{{ mtm.SrcCol.FieldName }}.SQLVal(), where.String())
							}
							rows, err := {{ DB.ObjName }}.Query(stmtStr)
							if err != nil {
								return nil, err
							}
							return {{ mtm.DstTab.ModelName }}FromRows(rows)
						},
						Add: func(mm *{{ mtm.DstTab.ModelName }}Instance) error {
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
						Remove: func(mm *{{ mtm.DstTab.ModelName }}Instance) error {
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
{{ endfor }}
`

// const modelInsertMethodTemp = `
// {{ for _, tab in DB.Tables }}
// func (m *{{ tab.ModelName }}Instance)Insert() error {
// 	colList := make([]string, 0, 16)
// 	valList := make([]string, 0, 16)
// 	{{ for _, col in tab.Columns }}
// 		if m.{{ col.FieldName }}.IsValid() {
// 			if m.{{ col.FieldName }}.IsNull() {
// 				colList = append(colList, {{ col.ColumnName }})
// 				valList = append(valList, "NULL")
// 			} else {
// 				colList = append(colList, {{ col.ColumnName }})
// 				valList = append(valList, m.{{ col.FieldName }}.SQLVal())
// 			}
// 		}
// 	{{ endfor }}
// 	stmtStr := fmt.Sprintf("INSERT INTO {{ tab.TableName }} (%s) VALUES (%s)", strings.Join(colList, ", "), strings.Join(valList, ", "))
// 	_, err := {{ DB.ObjName }}.Exec(stmtStr)
// 	return err
// }
// {{ endfor }}
// `
const modelInsertMethodTemp = `
{{ for _, tab in DB.Tables }}
func (m *{{ tab.ModelName }}Instance)Insert() error {
	colList := make([]string, 0, 16)
	valList := make([]string, 0, 16)
	for _, field := range m.fields() {
		if field.IsValid() {
			insertPair := field.InsertValuePair()
			colList = append(colList, insertPair[0])
			valList = append(valList, insertPair[1])
		}
	}
	stmtStr := fmt.Sprintf("INSERT INTO {{ tab.TableName }} (%s) VALUES (%s)", strings.Join(colList, ", "), strings.Join(valList, ", "))
	res, err := {{ DB.ObjName }}.Exec(stmtStr)
	if err != nil {
		return err
	}
	lastInsertId, err := res.LastInsertId()
	if err != nil {
		return err
	}
	m.{{ tab.AutoIncrement.FieldName }}.Set(lastInsertId)
	return nil
}
{{ endfor }}
`

const insertStmtMethodTemp = `
{{ for _, tab in DB.Tables }}
func (m *{{ tab.ModelName }}Instance)InsertStmt() *nbmysql.Stmt {
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
	return nbmysql.NewStmt(m, stmtStr)
}
{{ endfor }}
`

const modelUpdateMethodTemp = `
{{ for _, tab in DB.Tables }}
	func (m *{{ tab.ModelName }}Instance)Update() error {
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

const updateStmtMethodTemp = `
{{ for _, tab in DB.Tables }}
	func (m *{{ tab.ModelName }}Instance) UpdateStmt() (*nbmysql.Stmt, error) {
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
			return nil, errors.New("no valid field to update")
		}
		setList := make([]string, len(colList))
		for i := range colList {
			setList[i] = colList[i] + "=" + valList[i]
		}
		stmtStr := fmt.Sprintf("UPDATE {{ tab.TableName }} SET %s WHERE {{ tab.PrimaryKey.ColumnName }} = %s", strings.Join(setList, ", "), m.{{ tab.PrimaryKey.FieldName }}.SQLVal())
		return nbmysql.NewStmt(m, stmtStr), nil
	}
{{ endfor }}
`

const modelInsertOrUpdateMethodTemp = `
{{ for _, tab in DB.Tables }}
	func (m *{{ tab.ModelName }}Instance) InsertOrUpdate() error {
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
{{ endfor }}
`

const modelCheckMethodTemp = `
{{ for _, tab in DB.Tables }}
	func (m *{{ tab.ModelName }}Instance) check() error {
		{{ for _, col in tab.Columns }}
			{{ if col.ColumnName != tab.AutoIncrement.ColumnName && col.Nullable == false && col.Default == "" }}
				if !m.{{ col.FieldName }}.IsValid() || m.{{ col.FieldName }}.IsNull() {
					return errors.New("the {{ col.FieldName }} of {{ tab.ModelName }} cannot be null")
				}
			{{ endif }}
		{{ endfor }}
		return nil
	}
{{ endfor }}
`

const modelListTypeTemp = `
{{ for _, tab in DB.Tables }}
	type {{ tab.ModelName }}List []*{{ tab.ModelName }}Instance

	func (l {{ tab.ModelName }}List) Len() int {
		return len(l)
	}

	func (l {{ tab.ModelName }}List) Swap(i, j int) {
		l[i], l[j] = l[j], l[i]
	}

	type {{ tab.ArgName }}SortObj struct {
		list {{ tab.ModelName }}List
		lessFuncs []func(im, jm nbmysql.ModelInstance) (int, error)
	}

	func (so {{ tab.ArgName }}SortObj) Len() int {
		return so.list.Len()
	}

	func (so {{ tab.ArgName }}SortObj) Swap(i, j int) {
		so.list.Swap(i, j)
	}

	func (so {{ tab.ArgName }}SortObj) Less(i, j int) bool {
		iModel, jModel := so.list[i], so.list[j]
		for _, f := range so.lessFuncs {
			v, err := f(iModel, jModel)
			if err != nil {
				panic(err)
			}
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
{{ endfor }}
`

const countFuncTemp = `
{{ for _, tab in DB.Tables }}
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
{{ endfor }}
`

const modelInstanceTypeTemp = `
{{ for _, tab in DB.Tables }}
	type {{ tab.ModelName }}Instance struct {
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
		class *{{ tab.ModelName }}Class
	}
{{ endfor }}
`

// const getInstanceClassMethodTemp = `
// {{ for _, tab in DB.Tables }}
// 	func (m *{{ tab.ModelName }}Instance) GetClass() *{{ tab.ModelName }} {
// 		return &{{ tab.ModelName }}
// 		}
// {{ endfor }}
// `

const getFieldMethodTemp = `
{{ for _, tab in DB.Tables }}
	func (m *{{ tab.ModelName }}Instance) GetField(fieldName string) (nbmysql.Field, error) {
		switch fieldName {
		{{ for _, col in tab.Columns }}
		case "{{ col.FieldName }}":
			return &m.{{ col.FieldName }}, nil
		{{ endfor }}
		default:
			return nil, fmt.Errorf("{{ DB.Package }}.{{ tab.ModelName }}Instance.GetField() error: invalid field name (%s)", fieldName)
		}
	}
{{ endfor }}
`

const modelInvalidateMethodTemp = `
{{ for _, tab in DB.Tables }}
	func (m *{{ tab.ModelName }}Instance) Invalidate() {
		{{ for _, col in tab.Columns }}
			m.{{ col.FieldName }}.Invalidate()
		{{ endfor }}
	}
{{ endfor }}
`

const modelDeleteMethodTemp = `
{{ for _, tab in DB.Tables }}
	func (m *{{ tab.ModelName }}Instance) Delete() error {
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
{{ endfor }}
`

const deleteStmtMethodTemp = `
{{ for _, tab in DB.Tables }}
	func (m *{{ tab.ModelName }}Instance) DeleteStmt() (*nbmysql.Stmt, error) {
		if !m.{{ tab.PrimaryKey.FieldName }}.IsValid() || m.{{ tab.PrimaryKey.FieldName }}.IsNull() {
			return nil, errors.New("{{ tab.ModelName }}.Delete() error: primary key is not valid")
		}
		stmtStr := fmt.Sprintf("DELETE FROM {{ tab.TableName }} WHERE {{ tab.PrimaryKey.ColumnName }} = %s", m.{{ tab.PrimaryKey.FieldName }}.SQLVal())
		return nbmysql.NewStmt(m, stmtStr), nil
	}
{{ endfor }}
`

const setLastInsertIDMethodTemp = `
{{ for _, tab in DB.Tables }}
	func (m *{{ tab.ModelName }}Instance) SetLastInsertID(id int64) {
		m.{{ tab.AutoIncrement.FieldName }}.Set(id)
	}
{{ endfor }}
`

const modelListDistinctMethodTemp = `
{{ for _, tab in DB.Tables }}
	func (l {{ tab.ModelName }}List) Distinct(fcs ...nbmysql.FieldClass) {{ tab.ModelName }}List {
		if len(fcs) == 0 {
			fcs = append(fcs, {{ tab.ModelName }}.{{ tab.PrimaryKey.FieldName }})
		}
		m := make(map[string]bool)
		filteredList := make({{ tab.ModelName }}List, 0, l.Len())
		builder := strings.Builder{}
		for _, model := range l {
			for _, fc := range fcs {
				s, err := fc.DistFunc()(model)
				if err != nil {
					panic(err)
				}
				builder.WriteString(s)
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
{{ endfor }}
`

const modelListSortMethodTemp = `
{{ for _, tab in DB.Tables }}
	func (l {{ tab.ModelName }}List) Sort(reverse bool, fcs ...nbmysql.FieldClass) {
		lessFuncs := make([]func(im, jm nbmysql.ModelInstance) (int, error), len(fcs))
		for i, fc := range fcs {
			lessFuncs[i] = fc.LessFunc()
		}
		so := {{ tab.ArgName }}SortObj{list: l, lessFuncs: lessFuncs}
		if reverse {
			sort.Sort(sort.Reverse(so))
		} else {
			sort.Sort(so)
		}
	}
{{ endfor }}
`

const modelExistsMethodTemp = `
{{ for _, tab in DB.Tables }}
	func (m *{{ tab.ModelName }}Instance) Exists() (bool, error) {
		colList, valList := make([]string, 0, 16), make([]string, 0, 16)	
		{{ for _, col in tab.Columns }}
			if m.{{ col.FieldName }}.IsValid() {
				colList = append(colList, {{ col.ColumnName }})
				if m.{{ col.FieldName }}.IsNull() {
					valList = append(valList, "NULL")
				} else {
					valList = append(valList, m.{{ col.FieldName }}.SQLVal())
				}
			}
		{{ endfor }}
		whereList := make([]string, len(colList))
		for i := range colList {
			whereList[i] = colList[i] + "=" + valList[i]
		}
		stmtStr := fmt.Sprintf("SELECT EXISTS(SELECT * FROM {{ tab.TableName }} WHERE %s)", strings.Join(whereList, " and "))
		row := {{ DB.ObjName }}.QueryRow(stmtStr)
		var result int
		err := row.Scan(&result)
		if err != nil {
			return false, err
		}
		if result == 1 {
			return true, nil
		}
		return false, nil
	}
{{ endfor }}
`

const fieldsMethodTemp = `
{{ for _, tab in DB.Tables }}
	func (m *{{ tab.ModelName }}Instance) fields() []nbmysql.Field {
		fieldnames := m.class.fieldnames()
		fieldList := make([]nbmysql.Field, len(fieldnames))
		for i, fname := range fieldnames {
			field, err := m.GetField(fname)
			if err != nil {
				panic(err)
			}
			fieldList[i] = field
		}
		return fieldList
	}
{{ endfor }}
`

const classesFileImportTemp = `
import (
	"errors"
	"fmt"
	"strings"

	"github.com/wangjun861205/nbmysql"
)
`

// const modelClassTypeTemp = `
// {{ for _, tab in DB.Tables }}
// 	var {{ tab.ModelName }} struct {
// 		{{ for _, col in tab.Columns }}
// 			{{ switch col.FieldType }}
// 				{{ case "string" }}
// 					{{ col.FieldName }} *nbmysql.StringFieldClass
// 				{{ case "int64" }}
// 					{{ col.FieldName }} *nbmysql.IntFieldClass
// 				{{ case "float64" }}
// 					{{ col.FieldName }} *nbmysql.FloatFieldClass
// 				{{ case "bool" }}
// 					{{ col.FieldName }} *nbmysql.BoolFieldClass
// 				{{ case "time.Time" }}
// 					{{ switch col.MySqlType }}
// 						{{ case "DATETIME", "TIMESTAMP" }}
// 							{{ col.FieldName }} *nbmysql.DatetimeFieldClass
// 						{{ case "DATE" }}
// 							{{ col.FieldName }} *nbmysql.DateFieldClass
// 					{{ endswitch }}
// 			{{ endswitch }}
// 		{{ endfor }}
// 		New func(...nbmysql.SetStruct) (*{{ tab.ModelName }}Instance, error)
// 		All func() ({{ tab.ModelName }}List, error)
// 		Query func(*nbmysql.Where) ({{ tab.ModelName }}List, error)
// 		QueryOne func(*nbmysql.Where) (*{{ tab.ModelName }}Instance, error)
// 		Insert func(...[2]string) (int64, error)
// 		Update func(*nbmysql.Where, ...[2]string) (int64, error)
// 		Delete func(*nbmysql.Where) (int64, error)
// 		InsertOrUpdate func(...[2]string) (int64, error)
// 		Exists func(*nbmysql.Where) (bool, error)
// 		Count func(*nbmysql.Where) (int64, error)
// 		InsertStmt func(...[2]string) (*nbmysql.Stmt, error)
// 		UpdateStmt func(*nbmysql.Where, ...[2]string) (*nbmysql.Stmt, error)
// 		DeleteStmt func(*nbmysql.Where) (*nbmysql.Stmt, error)
// 	}
// {{ endfor }}
// `
const modelClassTypeTemp = `
{{ for _, tab in DB.Tables }}
	type {{ tab.ModelName }}Class struct {
		{{ for _, col in tab.Columns }}
			{{ switch col.FieldType }}
				{{ case "string" }}
					{{ col.FieldName }} *nbmysql.StringFieldClass
				{{ case "int64" }}
					{{ col.FieldName }} *nbmysql.IntFieldClass
				{{ case "float64" }}
					{{ col.FieldName }} *nbmysql.FloatFieldClass
				{{ case "bool" }}
					{{ col.FieldName }} *nbmysql.BoolFieldClass
				{{ case "time.Time" }}
					{{ switch col.MySqlType }}
						{{ case "DATETIME", "TIMESTAMP" }}
							{{ col.FieldName }} *nbmysql.DatetimeFieldClass
						{{ case "DATE" }}
							{{ col.FieldName }} *nbmysql.DateFieldClass
					{{ endswitch }}
			{{ endswitch }}
		{{ endfor }}
	}

	func (c *{{ tab.ModelName }}Class) New(ss ...nbmysql.SetStruct) (*{{ tab.ModelName }}Instance, error) {
		m := &{{ tab.ModelName }}Instance{class: c}
		{{ for _, col in tab.Columns }}
			m.{{ col.FieldName }} = {{ tab.ModelName }}.{{ col.FieldName }}.NewInstance()
		{{ endfor }}
		for _, s := range ss {
			if s.ModelName == "{{ tab.ModelName }}" {
				switch s.FieldName {
				{{ for _, col in tab.Columns }}
				case "{{ col.FieldName }}":
					s.Func(&m.{{ col.FieldName }})
				{{ endfor }}
				default:
					return nil, fmt.Errorf("{{ DB.Package }}.{{ tab.ModelName }}.New() error: invalid field name (%s)", s.FieldName)
				}
			} else {
				return nil, fmt.Errorf("{{ DB.Package }}.{{ tab.ModelName }}.New() error: invalid model name (%s)", s.ModelName)
			}
		}
		return m, nil
	}

	func (c *{{ tab.ModelName }}Class) All() ({{ tab.ModelName }}List, error) {
		stmtStr := "SELECT * FROM {{ tab.TableName }}"
		rows, err := {{ DB.ObjName }}.Query(stmtStr)
		if err != nil {
			return nil, err
		}
		return {{ tab.ModelName }}FromRows(rows)
	}

	func (c *{{ tab.ModelName }}Class) Query(where *nbmysql.Where) ({{ tab.ModelName }}List, error) {
		if where == nil {
			return c.All()
		}
		stmtStr := fmt.Sprintf("SELECT * FROM {{ tab.TableName }} WHERE %s", where.String())
		rows, err := {{ DB.ObjName }}.Query(stmtStr)
		if err != nil {
			return nil, err
		}
		return {{ tab.ModelName }}FromRows(rows)
	}

	func (c *{{ tab.ModelName }}Class) QueryOne(where *nbmysql.Where) (*{{ tab.ModelName }}Instance, error) {
		var stmtStr string
		if where == nil {
			stmtStr = "SELECT * FROM {{ tab.TableName }}"
		} else {
			stmtStr = fmt.Sprintf("SELECT * FROM {{ tab.TableName }} WHERE %s", where.String())
		}
		row := {{ DB.ObjName }}.QueryRow(stmtStr)
		return {{ tab.ModelName }}FromRow(row)
	}

	func (c *{{ tab.ModelName }}Class) Insert(args ...[2]string) (int64, error) {
		if len(args) == 0 {
			return -1, errors.New("{{ DB.Package }}.{{ tab.ModelName }}.Insert() error: empty arg list")
		}
		colList := make([]string, len(args))
		valList := make([]string, len(args))
		for i, arg := range args {
			colList[i] = arg[0]
			valList[i] = arg[1]
		}
		stmtStr := fmt.Sprintf("INSERT INTO {{ tab.TableName }} (%s) VALUES (%s)", strings.Join(colList, ", "), strings.Join(valList, ", "))
		res, err := {{ DB.ObjName }}.Exec(stmtStr)
		if err != nil {
			return -1, err
		}
		return res.LastInsertId()
	}

	func (c *{{ tab.ModelName }}Class) Update(where *nbmysql.Where, args ...[2]string) (int64, error) {
		if len(args) == 0 {
			return -1, errors.New("{{ DB.Package }}.{{ tab.ModelName }}.Update() error: empty arg list")
		}
		l := make([]string, len(args))
		for i, arg := range args {
			l[i] = arg[0] + " = " + arg[1]
		}
		var stmtStr string
		if where == nil {
			stmtStr = fmt.Sprintf("UPDATE {{ tab.TableName }} SET %s", strings.Join(l, ", "))
		} else {
			stmtStr = fmt.Sprintf("UPDATE {{ tab.TableName }} SET %s WHERE %s", strings.Join(l, ", "), where.String())
		}
		res, err := {{ DB.ObjName }}.Exec(stmtStr)
		if err != nil {
			return -1, err
		}
		return res.RowsAffected()
	}

	func (c *{{ tab.ModelName }}Class) Delete(where *nbmysql.Where) (int64, error) {
		var stmtStr string
		if where == nil {
			stmtStr = "DELETE FROM {{ tab.TableName }}"
		} else {
			stmtStr = fmt.Sprintf("DELETE FROM {{ tab.TableName }} WHERE %s", where.String())
		}
		res, err := {{ DB.ObjName }}.Exec(stmtStr)
		if err != nil {
			return -1, err
		}
		return res.RowsAffected()
	}

	func (c *{{ tab.ModelName }}Class) InsertOrUpdate(args ...[2]string) (int64, error) {
		if len(args) == 0 {
			return -1, errors.New("{{ DB.Package }}.{{ tab.ModelName }}.InsertOrUpdate() error: empty args")
		}
		colList := make([]string, len(args))
		valList := make([]string, len(args))
		setList := make([]string, len(args))
		for i, arg := range args {
			colList[i] = arg[0]
			valList[i] = arg[1]
			setList[i] = arg[0] + " = " + arg[1]
		}
		stmtStr := fmt.Sprintf("INSERT INTO {{ tab.TableName }} (%s) VALUES (%s) ON DUPLICATE KEY UPDATE {{ tab.AutoIncrement.ColumnName }} = LAST_INSERT_ID({{ tab.AutoIncrement.ColumnName }}), %s",
		strings.Join(colList, ", "), strings.Join(valList, ", "), strings.Join(setList, ", "))
		res, err := {{ DB.ObjName }}.Exec(stmtStr)
		if err != nil {
			return -1, err
		}
		return res.LastInsertId()
	}

	func (c *{{ tab.ModelName }}Class) Exists(where *nbmysql.Where) (bool, error) {
		var stmtStr string
		if where == nil {
			stmtStr = "SELECT EXISTS( SELECT * FROM {{ tab.TableName }}"
		} else {
			stmtStr = fmt.Sprintf("SELECT EXISTS( SELECT * FROM {{ tab.TableName }} WHERE %s )", where.String())
		}
		row := {{ DB.ObjName }}.QueryRow(stmtStr)
		var num int
		err := row.Scan(&num)
		if err != nil {
			return false, err
		}
		if num == 1 {
			return true, nil
		}
		return false, nil
	}

	func (c *{{ tab.ModelName }}Class) Count(where *nbmysql.Where) (int64, error) {
		var stmtStr string
		if where == nil {
			stmtStr = "SELECT COUNT(*) FROM {{ tab.TableName }}"
		} else {
			stmtStr = fmt.Sprintf("SELECT COUNT(*) FROM {{ tab.TableName }} WHERE %s", where.String())
		}
		row := {{ DB.ObjName }}.QueryRow(stmtStr)
		var num int64
		err := row.Scan(&num)
		if err != nil {
			return -1, err
		}
		return num, nil
	}

	func (c *{{ tab.ModelName }}Class) InsertStmt(args ...[2]string) (*nbmysql.Stmt, error) {
		if len(args) == 0 {
			return nil, errors.New("{{ DB.Package }}.{{ tab.ModelName }}.InsertStmt() error: empty args")
		}
		colList := make([]string, len(args))
		valList := make([]string, len(args))
		for i, arg := range args {
			colList[i] = arg[0]
			valList[i] = arg[1]
		}
		stmtStr := fmt.Sprintf("INSERT INTO {{ tab.TableName }} (%s) VALUES (%s)", strings.Join(colList, ", "), strings.Join(valList, ", "))
		return nbmysql.NewStmt(nil, stmtStr), nil
	}

	func (c *{{ tab.ModelName }}Class) UpdateStmt(where *nbmysql.Where, args ...[2]string) (*nbmysql.Stmt, error) {
		if len(args) == 0 {
			return nil, errors.New("{{ DB.Package }}.{{ tab.ModelName }}.InsertStmt() error: empty args")
		}
		setList := make([]string, len(args))
		for i, arg := range args {
			setList[i] = arg[0] + " = " + arg[1]
		}
		var stmtStr string
		if where == nil {
			stmtStr = fmt.Sprintf("UPDATE {{ tab.TableName }} SET %s", strings.Join(setList, ", "))
		} else {
			stmtStr = fmt.Sprintf("UPDATE {{ tab.TableName }} SET %s WHERE %s", strings.Join(setList, ", "), where.String())
		}
		return nbmysql.NewStmt(nil, stmtStr), nil
	}

	func (c *{{ tab.ModelName }}Class) DeleteStmt(where *nbmysql.Where) (*nbmysql.Stmt, error) {
		var stmtStr string
		if where == nil {
			stmtStr = "DELETE FROM {{ tab.TableName }}"
		} else {
			stmtStr = fmt.Sprintf("DELETE FROM {{ tab.TableName }} WHERE %s", where.String())
		}
		return nbmysql.NewStmt(nil, stmtStr), nil
	}

	func (c *{{ tab.ModelName }}Class) fieldnames() []string {
		return []string{
			{{ for _, col in tab.Columns }}
				"{{ col.FieldName }}",
			{{ endfor }}
		}
	}

	func (c *{{ tab.ModelName }}Class) primaryKeyFieldName() string {
		return "{{ tab.PrimaryKey.FieldName }}"
	}

	func (c *{{ tab.ModelName }}Class) autoIncrementFieldName() string {
		return "{{ tab.AutoIncrement.FieldName }}"
	}
{{ endfor }}
`

// const modelClassInitTemp = `
// func init() {
// {{ for _, tab in DB.Tables }}
// 	{{ tab.ModelName }} = struct {
// 		{{ for _, col in tab.Columns }}
// 			{{ switch col.FieldType }}
// 				{{ case "string" }}
// 					{{ col.FieldName }} *nbmysql.StringFieldClass
// 				{{ case "int64" }}
// 					{{ col.FieldName }} *nbmysql.IntFieldClass
// 				{{ case "float64" }}
// 					{{ col.FieldName }} *nbmysql.FloatFieldClass
// 				{{ case "bool" }}
// 					{{ col.FieldName }} *nbmysql.BoolFieldClass
// 				{{ case "time.Time" }}
// 					{{ switch col.MySqlType }}
// 						{{ case "DATETIME", "TIMESTAMP" }}
// 							{{ col.FieldName }} *nbmysql.DatetimeFieldClass
// 						{{ case "DATE" }}
// 							{{ col.FieldName }} *nbmysql.DateFieldClass
// 					{{ endswitch }}
// 			{{ endswitch }}
// 		{{ endfor }}
// 		New func(...nbmysql.SetStruct) (*{{ tab.ModelName }}Instance, error)
// 		All func() ({{ tab.ModelName }}List, error)
// 		Query func(*nbmysql.Where) ({{ tab.ModelName }}List, error)
// 		QueryOne func(*nbmysql.Where) (*{{ tab.ModelName }}Instance, error)
// 		Insert func(...[2]string) (int64, error)
// 		Update func(*nbmysql.Where, ...[2]string) (int64, error)
// 		Delete func(*nbmysql.Where) (int64, error)
// 		InsertOrUpdate func(...[2]string) (int64, error)
// 		Exists func(*nbmysql.Where) (bool, error)
// 		Count func(*nbmysql.Where) (int64, error)
// 		InsertStmt func(...[2]string) (*nbmysql.Stmt, error)
// 		UpdateStmt func(*nbmysql.Where, ...[2]string) (*nbmysql.Stmt, error)
// 		DeleteStmt func(*nbmysql.Where) (*nbmysql.Stmt, error)
// 	}{
// 		{{ for _, col in tab.Columns }}
// 			{{ switch col.FieldType }}
// 				{{ case "string" }}
// 					{{ col.FieldName }}: nbmysql.NewStringFieldClass("{{ tab.TableName }}", "{{ col.ColumnName }}", "{{ tab.ModelName }}", "{{ col.FieldName }}"),
// 				{{ case "int64" }}
// 					{{ col.FieldName }}: nbmysql.NewIntFieldClass("{{ tab.TableName }}", "{{ col.ColumnName }}", "{{ tab.ModelName }}", "{{ col.FieldName }}"),
// 				{{ case "float64" }}
// 					{{ col.FieldName }}: nbmysql.NewFloatFieldClass("{{ tab.TableName }}", "{{ col.ColumnName }}", "{{ tab.ModelName }}", "{{ col.FieldName }}"),
// 				{{ case "bool" }}
// 					{{ col.FieldName }}: nbmysql.NewBoolFieldClass("{{ tab.TableName }}", "{{ col.ColumnName }}", "{{ tab.ModelName }}", "{{ col.FieldName }}"),
// 				{{ case "time.Time" }}
// 					{{ switch col.MySqlType }}
// 						{{ case "DATETIME", "TIMESTAMP" }}
// 							{{ col.FieldName }}: nbmysql.NewDatetimeFieldClass("{{ tab.TableName }}", "{{ col.ColumnName }}", "{{ tab.ModelName }}", "{{ col.FieldName }}"),
// 						{{ case "DATE" }}
// 							{{ col.FieldName }}: nbmysql.NewDateFieldClass("{{ tab.TableName }}", "{{ col.ColumnName }}", "{{ tab.ModelName }}", "{{ col.FieldName }}"),
// 				{{ endswitch }}
// 			{{ endswitch }}
// 		{{ endfor }}
// 		New: func(ss ...nbmysql.SetStruct) (*{{ tab.ModelName }}Instance, error) {
// 			m := &{{ tab.ModelName }}Instance{}
// 			{{ for _, col in tab.Columns }}
// 				m.{{ col.FieldName }} = {{ tab.ModelName }}.{{ col.FieldName }}.NewInstance()
// 			{{ endfor }}
// 			for _, s := range ss {
// 				if s.ModelName == "{{ tab.ModelName }}" {
// 					switch s.FieldName {
// 					{{ for _, col in tab.Columns }}
// 					case "{{ col.FieldName }}":
// 						s.Func(&m.{{ col.FieldName }})
// 					{{ endfor }}
// 					default:
// 						return nil, fmt.Errorf("{{ DB.Package }}.{{ tab.ModelName }}.New() error: invalid field name (%s)", s.FieldName)
// 					}
// 				} else {
// 					return nil, fmt.Errorf("{{ DB.Package }}.{{ tab.ModelName }}.New() error: invalid model name (%s)", s.ModelName)
// 				}
// 			}
// 			return m, nil
// 		},
// 		All: func() ({{ tab.ModelName }}List, error) {
// 			stmtStr := "SELECT * FROM {{ tab.TableName }}"
// 			rows, err := {{ DB.ObjName }}.Query(stmtStr)
// 			if err != nil {
// 				return nil, err
// 			}
// 			return {{ tab.ModelName }}FromRows(rows)
// 		},
// 		Query: func(where *nbmysql.Where) ({{ tab.ModelName }}List, error) {
// 			if where == nil {
// 				return {{ tab.ModelName }}.All()
// 			}
// 			stmtStr := fmt.Sprintf("SELECT * FROM {{ tab.TableName }} WHERE %s", where.String())
// 			rows, err := {{ DB.ObjName }}.Query(stmtStr)
// 			if err != nil {
// 				return nil, err
// 			}
// 			return {{ tab.ModelName }}FromRows(rows)
// 		},
// 		QueryOne: func(where *nbmysql.Where) (*{{ tab.ModelName }}Instance, error) {
// 			var stmtStr string
// 			if where == nil {
// 				stmtStr = "SELECT * FROM {{ tab.TableName }}"
// 			} else {
// 				stmtStr = fmt.Sprintf("SELECT * FROM {{ tab.TableName }} WHERE %s", where.String())
// 			}
// 			row := {{ DB.ObjName }}.QueryRow(stmtStr)
// 			return {{ tab.ModelName }}FromRow(row)
// 		},
// 		Insert: func(args ...[2]string) (int64, error) {
// 			if len(args) == 0 {
// 				return -1, errors.New("{{ DB.Package }}.{{ tab.ModelName }}.Insert() error: empty arg list")
// 			}
// 			colList := make([]string, len(args))
// 			valList := make([]string, len(args))
// 			for i, arg := range args {
// 				colList[i] = arg[0]
// 				valList[i] = arg[1]
// 			}
// 			stmtStr := fmt.Sprintf("INSERT INTO {{ tab.TableName }} (%s) VALUES (%s)", strings.Join(colList, ", "), strings.Join(valList, ", "))
// 			res, err := {{ DB.ObjName }}.Exec(stmtStr)
// 			if err != nil {
// 				return -1, err
// 			}
// 			return res.LastInsertId()
// 		},
// 		Update: func(where *nbmysql.Where, args ...[2]string) (int64, error) {
// 			if len(args) == 0 {
// 				return -1, errors.New("{{ DB.Package }}.{{ tab.ModelName }}.Update() error: empty arg list")
// 			}
// 			l := make([]string, len(args))
// 			for i, arg := range args {
// 				l[i] = arg[0] + " = " + arg[1]
// 			}
// 			var stmtStr string
// 			if where == nil {
// 				stmtStr = fmt.Sprintf("UPDATE {{ tab.TableName }} SET %s", strings.Join(l, ", "))
// 			} else {
// 				stmtStr = fmt.Sprintf("UPDATE {{ tab.TableName }} SET %s WHERE %s", strings.Join(l, ", "), where.String())
// 			}
// 			res, err := {{ DB.ObjName }}.Exec(stmtStr)
// 			if err != nil {
// 				return -1, err
// 			}
// 			return res.RowsAffected()
// 		},
// 		Delete: func(where *nbmysql.Where) (int64, error) {
// 			var stmtStr string
// 			if where == nil {
// 				stmtStr = "DELETE FROM {{ tab.TableName }}"
// 			} else {
// 				stmtStr = fmt.Sprintf("DELETE FROM {{ tab.TableName }} WHERE %s", where.String())
// 			}
// 			res, err := {{ DB.ObjName }}.Exec(stmtStr)
// 			if err != nil {
// 				return -1, err
// 			}
// 			return res.RowsAffected()
// 		},
// 		InsertOrUpdate: func(args ...[2]string) (int64, error) {
// 			if len(args) == 0 {
// 				return -1, errors.New("{{ DB.Package }}.{{ tab.ModelName }}.InsertOrUpdate() error: empty args")
// 			}
// 			colList := make([]string, len(args))
// 			valList := make([]string, len(args))
// 			setList := make([]string, len(args))
// 			for i, arg := range args {
// 				colList[i] = arg[0]
// 				valList[i] = arg[1]
// 				setList[i] = arg[0] + " = " + arg[1]
// 			}
// 			stmtStr := fmt.Sprintf("INSERT INTO {{ tab.TableName }} (%s) VALUES (%s) ON DUPLICATE KEY UPDATE {{ tab.AutoIncrement.ColumnName }} = LAST_INSERT_ID({{ tab.AutoIncrement.ColumnName }}), %s",
// 			strings.Join(colList, ", "), strings.Join(valList, ", "), strings.Join(setList, ", "))
// 			res, err := {{ DB.ObjName }}.Exec(stmtStr)
// 			if err != nil {
// 				return -1, err
// 			}
// 			return res.LastInsertId()
// 		},
// 		Exists: func(where *nbmysql.Where) (bool, error) {
// 			var stmtStr string
// 			if where == nil {
// 				stmtStr = "SELECT EXISTS( SELECT * FROM {{ tab.TableName }}"
// 			} else {
// 				stmtStr = fmt.Sprintf("SELECT EXISTS( SELECT * FROM {{ tab.TableName }} WHERE %s )", where.String())
// 			}
// 			row := {{ DB.ObjName }}.QueryRow(stmtStr)
// 			var num int
// 			err := row.Scan(&num)
// 			if err != nil {
// 				return false, err
// 			}
// 			if num == 1 {
// 				return true, nil
// 			}
// 			return false, nil
// 		},
// 		Count: func(where *nbmysql.Where) (int64, error) {
// 			var stmtStr string
// 			if where == nil {
// 				stmtStr = "SELECT COUNT(*) FROM {{ tab.TableName }}"
// 			} else {
// 				stmtStr = fmt.Sprintf("SELECT COUNT(*) FROM {{ tab.TableName }} WHERE %s", where.String())
// 			}
// 			row := {{ DB.ObjName }}.QueryRow(stmtStr)
// 			var num int64
// 			err := row.Scan(&num)
// 			if err != nil {
// 				return -1, err
// 			}
// 			return num, nil
// 		},
// 		InsertStmt: func(args ...[2]string) (*nbmysql.Stmt, error) {
// 			if len(args) == 0 {
// 				return nil, errors.New("{{ DB.Package }}.{{ tab.ModelName }}.InsertStmt() error: empty args")
// 			}
// 			colList := make([]string, len(args))
// 			valList := make([]string, len(args))
// 			for i, arg := range args {
// 				colList[i] = arg[0]
// 				valList[i] = arg[1]
// 			}
// 			stmtStr := fmt.Sprintf("INSERT INTO {{ tab.TableName }} (%s) VALUES (%s)", strings.Join(colList, ", "), strings.Join(valList, ", "))
// 			return nbmysql.NewStmt(nil, stmtStr), nil
// 		},
// 		UpdateStmt: func(where *nbmysql.Where, args ...[2]string) (*nbmysql.Stmt, error) {
// 			if len(args) == 0 {
// 				return nil, errors.New("{{ DB.Package }}.{{ tab.ModelName }}.InsertStmt() error: empty args")
// 			}
// 			setList := make([]string, len(args))
// 			for i, arg := range args {
// 				setList[i] = arg[0] + " = " + arg[1]
// 			}
// 			var stmtStr string
// 			if where == nil {
// 				stmtStr = fmt.Sprintf("UPDATE {{ tab.TableName }} SET %s", strings.Join(setList, ", "))
// 			} else {
// 				stmtStr = fmt.Sprintf("UPDATE {{ tab.TableName }} SET %s WHERE %s", strings.Join(setList, ", "), where.String())
// 			}
// 			return nbmysql.NewStmt(nil, stmtStr), nil
// 		},
// 		DeleteStmt: func(where *nbmysql.Where) (*nbmysql.Stmt, error) {
// 			var stmtStr string
// 			if where == nil {
// 				stmtStr = "DELETE FROM {{ tab.TableName }}"
// 			} else {
// 				stmtStr = fmt.Sprintf("DELETE FROM {{ tab.TableName }} WHERE %s", where.String())
// 			}
// 			return nbmysql.NewStmt(nil, stmtStr), nil
// 		},
// 	}
// {{ endfor }}
// }
// `
const modelClassInitTemp = `
{{ for _, tab in DB.Tables }}
	var {{ tab.ModelName }} *{{ tab.ModelName }}Class
{{ endfor }}

func init() {
{{ for _, tab in DB.Tables }}
	{{ tab.ModelName }} = &{{ tab.ModelName }}Class {
		{{ for _, col in tab.Columns }}
			{{ switch col.FieldType }}
				{{ case "string" }}
					{{ col.FieldName }}: nbmysql.NewStringFieldClass("{{ tab.TableName }}", "{{ col.ColumnName }}", "{{ tab.ModelName }}", "{{ col.FieldName }}"),
				{{ case "int64" }}
					{{ col.FieldName }}: nbmysql.NewIntFieldClass("{{ tab.TableName }}", "{{ col.ColumnName }}", "{{ tab.ModelName }}", "{{ col.FieldName }}"),
				{{ case "float64" }}
					{{ col.FieldName }}: nbmysql.NewFloatFieldClass("{{ tab.TableName }}", "{{ col.ColumnName }}", "{{ tab.ModelName }}", "{{ col.FieldName }}"),
				{{ case "bool" }}
					{{ col.FieldName }}: nbmysql.NewBoolFieldClass("{{ tab.TableName }}", "{{ col.ColumnName }}", "{{ tab.ModelName }}", "{{ col.FieldName }}"),
				{{ case "time.Time" }}
					{{ switch col.MySqlType }}
						{{ case "DATETIME", "TIMESTAMP" }}
							{{ col.FieldName }}: nbmysql.NewDatetimeFieldClass("{{ tab.TableName }}", "{{ col.ColumnName }}", "{{ tab.ModelName }}", "{{ col.FieldName }}"),
						{{ case "DATE" }}
							{{ col.FieldName }}: nbmysql.NewDateFieldClass("{{ tab.TableName }}", "{{ col.ColumnName }}", "{{ tab.ModelName }}", "{{ col.FieldName }}"),
				{{ endswitch }}
			{{ endswitch }}
		{{ endfor }}
	}
{{ endfor }}
}
`

// const newModelFuncTemp = `
// {{ for _, tab in DB.Tables }}
// func New{{ tab.ModelName }}(fields ...FieldArg) (*{{ tab.ModelName }}, error) {
// 	m := &{{ tab.ModelName }}{}
// 	for _, field := range fields {
// 		switch f := field.(type) {
// 		{{ for _, col in tab.Columns }}
// 		case *{{ tab.ModelName }}{{ col.FieldName }}:
// 			if f == nil {
// 				{{ switch col.FieldType }}
// 					{{ case "string" }}
// 						m.{{ col.FieldName }}.Set("", true)
// 					{{ case "int64" }}
// 						m.{{ col.FieldName }}.Set(0, true)
// 					{{ case "float64" }}
// 						m.{{ col.FieldName }}.Set(0.0, true)
// 					{{ case "bool" }}
// 						m.{{ col.FieldName }}.Set(false, false)
// 					{{ case "time.Time" }}
// 						m.{{ col.FieldName }}.Set(time.Time{}, false)
// 				{{ endswitch }}
// 			}
// 			m.{{ col.FieldName }}.Set({{ col.FieldType }}(*f))
// 		{{ endfor }}
// 		default:
// 			return nil, errors.New("invalid field in New{{ tab.ModelName }}()")
// 		}
// 	}
// 	return m, nil
// }
// {{ endfor }}
// `

// const updateFuncTemp = `
// {{ for _, tab in DB.Tables }}
// 	func Update{{ tab.ModelName }}(where string, update ...FieldArg) error {
// 		if len(update) == 0 {
// 			return errors.New("Update{{ tab.ModelName }}() error: update list cannot be empty")
// 		}
// 		for k, v := range {{ tab.ModelName }}Map {
// 			where = strings.Replace(where, k, v, -1)
// 		}
// 		setList := make([]string, len(update))
// 		for i, f := range update {
// 			if f.tableName() != "{{ tab.TableName }}" {
// 				return fmt.Errorf("Update{{ tab.ModelName }}() error: %T not belong to {{ tab.ModelName }}", f)
// 			}
// 			if f == nil {
// 				setList[i] = f.columnName() + "=" + "NULL"
// 			} else {
// 				setList[i] = f.columnName() + "=" + f.sqlValue()
// 			}
// 		}
// 		stmtStr := fmt.Sprintf("UPDATE {{ tab.TableName }} SET %s WHERE %s", strings.Join(setList, ", "), where)
// 		_, err := {{ DB.ObjName }}.Exec(stmtStr)
// 		return err
// 	}
// {{ endfor }}
// `

// const updateStmtFuncTemp = `
// {{ for _, tab in DB.Tables }}
// 	func Update{{ tab.ModelName }}Stmt(where string, update ...FieldArg) (*nbmysql.Stmt, error) {
// 		if len(update) == 0 {
// 			return nil, errors.New("Update{{ tab.ModelName }}() error: update list cannot be empty")
// 		}
// 		for k, v := range {{ tab.ModelName }}Map {
// 			where = strings.Replace(where, k, v, -1)
// 		}
// 		setList := make([]string, len(update))
// 		for i, f := range update {
// 			if f.tableName() != "{{ tab.TableName }}" {
// 				return nil, fmt.Errorf("Update{{ tab.ModelName }}() error: %T not belong to {{ tab.ModelName }}", f)
// 			}
// 			if f == nil {
// 				setList[i] = f.columnName() + "=" + "NULL"
// 			} else {
// 				setList[i] = f.columnName() + "=" + f.sqlValue()
// 			}
// 		}
// 		stmtStr := fmt.Sprintf("UPDATE {{ tab.TableName }} SET %s WHERE %s", strings.Join(setList, ", "), where)
// 		return nbmysql.NewStmt(nil, stmtStr), nil
// 	}
// {{ endfor }}
// `

// const modelInsertFuncTemp = `
// {{ for _, tab in DB.Tables }}
// 	func Insert{{ tab.ModelName }}(fields ...FieldArg) (*{{ tab.ModelName }}, error) {
// 		m, err := New{{ tab.ModelName }}(fields...)
// 		if err != nil {
// 			return nil, err
// 		}
// 		colList := make([]string, len(fields))
// 		valList := make([]string, len(fields))
// 		for i, f := range fields {
// 			colList[i] = f.columnName()
// 			valList[i] = f.sqlValue()
// 		}
// 		stmtStr := fmt.Sprintf("INSERT INTO {{ tab.TableName }} (%s) VALUES (%s)", strings.Join(colList, ", "), strings.Join(valList, ", "))
// 		res, err := {{ DB.ObjName }}.Exec(stmtStr)
// 		if err != nil {
// 			return nil, err
// 			}
// 		lastInsertId, err := res.LastInsertId()
// 		if err != nil {
// 			return nil, err
// 		}
// 		m.{{ tab.AutoIncrement.FieldName }}.Set(lastInsertId)
// 		return m, nil
// 		}
// {{ endfor }}
// `

// const insertStmtFuncTemp = `
// {{ for _, tab in DB.Tables }}
// 	func Insert{{ tab.ModelName }}Stmt(fields ...FieldArg) *nbmysql.Stmt {
// 		colList := make([]string, len(fields))
// 		valList := make([]string, len(fields))
// 		for i, f := range fields {
// 			colList[i] = f.columnName()
// 			valList[i] = f.sqlValue()
// 		}
// 		stmtStr := fmt.Sprintf("INSERT INTO {{ tab.TableName }} (%s) VALUES (%s)", strings.Join(colList, ", "), strings.Join(valList, ", "))
// 		return nbmysql.NewStmt(nil, stmtStr)
// 		}
// {{ endfor }}
// `

// const allModelFuncTemp = `
// {{ for _, tab in DB.Tables }}
// 	func All{{ tab.ModelName }}() ({{ tab.ModelName }}List, error) {
// 		rows, err := {{ DB.ObjName }}.Query("SELECT * FROM {{ tab.TableName }}")
// 		if err != nil {
// 			return nil, err
// 		}
// 		return {{ tab.ModelName }}FromRows(rows)
// 	}
// {{ endfor }}
// `

// const queryFuncTemp = `
// {{ for _, tab in DB.Tables }}
// 	func Query{{ tab.ModelName }}(where string) ({{ tab.ModelName }}List, error) {
// 		for k, v := range {{ tab.ModelName }}Map {
// 			where = strings.Replace(where, k, v, -1)
// 		}
// 		rows, err := {{ DB.ObjName }}.Query(fmt.Sprintf("SELECT * FROM {{ tab.TableName }} WHERE %s", where))
// 		if err != nil {
// 			return nil, err
// 		}
// 		return {{ tab.ModelName }}FromRows(rows)
// 	}
// {{ endfor }}
// `

// const queryOneFuncTemp = `
// {{ for _, tab in DB.Tables }}
// 	func QueryOne{{ tab.ModelName }}(where string) (*{{ tab.ModelName }}, error) {
// 		for k, v := range {{ tab.ModelName }}Map {
// 			where = strings.Replace(where, k, v, -1)
// 		}
// 		row := {{ DB.ObjName }}.QueryRow(fmt.Sprintf("SELECT * FROM {{ tab.TableName }} WHERE %s", where))
// 		return {{ tab.ModelName }}FromRow(row)
// 	}
// {{ endfor }}
// `

// const deleteFuncTemp = `
// {{ for _, tab in DB.Tables }}
// 	func Delete{{ tab.ModelName }}(where string) (int64, error) {
// 		for k, v := range {{ tab.ModelName }}Map {
// 			where = strings.Replace(where, k, v, -1)
// 		}
// 		var stmtStr string
// 		if where != "" {
// 			stmtStr = fmt.Sprintf("DELETE FROM {{ tab.TableName }} WHERE %s", where)
// 		} else {
// 			stmtStr = fmt.Sprintf("DELETE FROM {{ tab.TableName }}")
// 		}
// 		res, err := {{ DB.ObjName }}.Exec(stmtStr)
// 		if err != nil {
// 			return -1, err
// 		}
// 		return res.RowsAffected()
// 	}
// {{ endfor }}
// `

// const deleteStmtFuncTemp = `
// {{ for _, tab in DB.Tables }}
// 	func Delete{{ tab.ModelName }}Stmt(where string) *nbmysql.Stmt {
// 		for k, v := range {{ tab.ModelName }}Map {
// 			where = strings.Replace(where, k, v, -1)
// 		}
// 		var stmtStr string
// 		if where != "" {
// 			stmtStr = fmt.Sprintf("DELETE FROM {{ tab.TableName }} WHERE %s", where)
// 		} else {
// 			stmtStr = fmt.Sprintf("DELETE FROM {{ tab.TableName }}")
// 		}
// 		return nbmysql.NewStmt(nil, stmtStr)
// 	}
// {{ endfor }}
// `

// const whereMapTemp = `
// {{ for _, tab in DB.Tables }}
// 	var {{ tab.ModelName }}Map = map[string]string {
// 		{{ for _, col in tab.Columns }}
// 			"@{{ col.FieldName }}": "{{ col.ColumnName }}",
// 		{{ endfor }}
// 	}
// {{ endfor }}
// `

// const fieldTypeTemp = `
// type FieldArg interface {
// 	columnName() string
// 	tableName() string
// 	sqlValue() string
// }
// {{ for _, tab in DB.Tables }}
// 	{{ for _, col in tab.Columns }}
// 		{{ switch col.FieldType }}
// 			{{ case "int64" }}
// 				type {{ tab.ModelName }}{{ col.FieldName }} int64
// 				func New{{ tab.ModelName }}{{ col.FieldName }}(val int64) *{{ tab.ModelName }}{{ col.FieldName }} {
// 					f := {{ tab.ModelName }}{{ col.FieldName }}(val)
// 					return &f
// 				}
// 				func (f *{{ tab.ModelName }}{{ col.FieldName }}) columnName() string {
// 					return "{{ col.ColumnName }}"
// 				}
// 				func (f *{{ tab.ModelName }}{{ col.FieldName }}) sqlValue() string {
// 					return fmt.Sprintf("%d", *f)
// 				}
// 			{{ case "string" }}
// 				type {{ tab.ModelName }}{{ col.FieldName }} string
// 				func New{{ tab.ModelName }}{{ col.FieldName}}(val string) *{{ tab.ModelName }}{{ col.FieldName}} {
// 					f := {{ tab.ModelName }}{{ col.FieldName }}(val)
// 					return &f
// 				}
// 				func (f *{{ tab.ModelName }}{{ col.FieldName }}) columnName() string {
// 					return "{{ col.ColumnName }}"
// 				}
// 				func (f *{{ tab.ModelName }}{{ col.FieldName }}) sqlValue() string {
// 					return fmt.Sprintf("%q", *f)
// 				}
// 			{{ case "float64" }}
// 				type {{ table.ModelName }}{{ col.FieldName }} float64
// 				func New{{ tab.ModelName }}{{ col.FieldName }}(val float64) *{{ tab.ModelName }}{{ col.FieldName }} {
// 					f := {{ tab.ModelName }}{{ col.FieldName }}(val)
// 					return &f
// 				}
// 				func (f *{{ tab.ModelName }}{{ col.FieldName }}) columnName() string {
// 					return "{{ col.ColumnName }}"
// 				}
// 				func (f *{{ tab.ModelName }}{{ col.FieldName }}) sqlValue() string {
// 					return fmt.Sprintf("%f", *f)
// 				}
// 			{{ case "bool" }}
// 				type {{ tab.ModelName }}{{ col.FieldName }} bool
// 				func New{{ tab.ModelName }}{{ col.FieldName }}(val bool) *{{ tab.ModelName }}{{ col.FieldName }} {
// 					f := {{ tab.ModelName }}{{ col.FieldName }}(val)
// 					return &f
// 				}
// 				func (f *{{ tab.ModelName }}{{ col.FieldName }}) columnName() string {
// 					return "{{ col.ColumnName }}"
// 				}
// 				func (f *{{ tab.ModelName }}{{ col.FieldName }}) sqlValue() string {
// 					return fmt.Sprintf("%t", *f)
// 				}
// 			{{ case "time.Time" }}
// 				type {{ tab.ModelName }}{{ col.FieldName }} time.Time
// 				func New{{ tab.ModelName }}{{ col.FieldName }}(val time.Time) *{{ tab.ModelName }}{{ col.FieldName }}{
// 					f := {{ tab.ModelName }}{{ col.FieldName }}(val)
// 					return &f
// 				}
// 				func (f *{{ tab.ModelName }}{{ col.FieldName }}) columnName() string {
// 					return "{{ col.ColumnName }}"
// 				}
// 				func (f *{{ tab.ModelName }}{{ col.FieldName }}) sqlValue() string {
// 					{{ switch col.MySqlType }}
// 						{{ case "DATE" }}
// 							return time.Time(*f).Format("2006-01-02")
// 						{{ case "DATETIME" }}
// 							return time.Time(*f).Format("2006-01-02 15:04:05")
// 						{{ case "TIMESTAMP" }}
// 							return time.Time(*f).Format("2006-01-02 15:04:05")
// 					{{ endswitch }}
// 				}
// 		{{ endswitch }}
// 		func (f *{{ tab.ModelName}}{{ col.FieldName}}) tableName() string {
// 			return "{{ tab.TableName }}"
// 		}
// 	{{ endfor }}
// {{ endfor }}
// `

// const stmtTypeTemp = `
// type stmtType int
// const (
// 	insert stmtType = iota
// 	update
// 	delete
// 	other
// )

// type Stmt struct {
// 	model nbmysql.ModelInstance
// 	typ stmtType
// 	stmt string
// 	lastInsertID int64
// }

// func nbmysql.NewStmt(model nbmysql.ModelInstance, stmtStr string) *nbmysql.Stmt {
// 	stmt := &Stmt{model: model, stmt: stmtStr}
// 	l := strings.Split(strings.ToUpper(stmtStr), " ")
// 	switch l[0] {
// 	case "INSERT":
// 		stmt.typ = insert
// 	case "UPDATE":
// 		stmt.typ = update
// 	case "DELETE":
// 		stmt.typ = delete
// 	default:
// 		stmt.typ = other
// 	}
// 	return stmt
// }

// func (s *nbmysql.Stmt) Join(stmts ...*nbmysql.Stmt) StmtList {
// 	l := make(StmtList, len(stmts) + 1)
// 	l[0] = s
// 	for i := 0; i < len(stmts); i++ {
// 		l[i+1] = stmts[i]
// 	}
// 	return l
// }

// func (s *nbmysql.Stmt) Exec() error {
// 	res, err := {{ DB.ObjName }}.Exec(s.stmt)
// 	if err != nil {
// 		return err
// 	}
// 	if s.model != nil {
// 		switch s.typ {
// 		case insert:
// 			lastInsertID, err := res.LastInsertId()
// 			if err != nil {
// 				return err
// 			}
// 			s.model.SetLastInsertID(lastInsertID)
// 		case delete:
// 			s.model.Invalidate()
// 		}
// 	}
// 	return nil
// }

// type StmtList []*nbmysql.Stmt

// func (sl StmtList) Exec() error {
// 	tx, err := {{ DB.ObjName }}.Begin()
// 	if err != nil {
// 		return err
// 	}
// 	for _, s := range sl {
// 		res, err := tx.Exec(s.stmt)
// 		if err != nil {
// 			tx.Rollback()
// 			return err
// 		}
// 		if s.typ == insert && s.model != nil {
// 			lastInsertID, err := res.LastInsertId()
// 			if err != nil {
// 				tx.Rollback()
// 				return err
// 			}
// 			s.lastInsertID = lastInsertID
// 		}
// 	}
// 	err = tx.Commit()
// 	if err != nil {
// 		return err
// 	}
// 	for _, s := range sl {
// 		if s.model != nil {
// 			switch s.typ {
// 			case insert:
// 				s.model.SetLastInsertID(s.lastInsertID)
// 			case delete:
// 				s.model.Invalidate()
// 			}
// 		}
// 	}
// 	return nil
// }

// func (sl StmtList) Join(sls ...StmtList) StmtList {
// 	for _, l := range sls {
// 		sl = append(sl, l...)
// 	}
// 	return sl
// 	}
// `
