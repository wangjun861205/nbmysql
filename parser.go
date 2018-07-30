package nbmysql

import (
	"errors"
	"io/ioutil"
	"os"
	"regexp"

	_ "github.com/go-sql-driver/mysql"
)

type ForeignKey struct {
	SrcCol *Column
	DstCol *Column
	DstTab *Table
}

type ManyToMany struct {
	SrcCol      *Column
	MidLeftCol  *Column
	MidRightCol *Column
	DstCol      *Column
	MidTab      *Table
	DstTab      *Table
}

type ForeignKeyInfo struct {
	SrcColName string
	DstColName string
	DstTabName string
}

type ManyToManyInfo struct {
	SrcColName string
	DstColName string
	DstTabName string
}

var GoMidMap = map[string]string{
	"int":    "Int",
	"float":  "Float",
	"string": "String",
	"bool":   "Bool",
	"time":   "Time",
}

var MySqlGoMap = map[string]string{
	"INT":       "int64",
	"FLOAT":     "float64",
	"BOOL":      "bool",
	"DATETIME":  "time.Time",
	"DATE":      "time.Time",
	"TIMESTAMP": "time.Time",
	"VARCHAR":   "string",
	"TEXT":      "string",
}

var MySqlMidMap = map[string]string{
	"INT":       "Int",
	"FLOAT":     "Float",
	"BOOL":      "Bool",
	"DATETIME":  "Time",
	"DATE":      "Time",
	"TIMESTAMP": "Time",
	"VARCHAR":   "String",
	"TEXT":      "String",
}

var PackageRe = regexp.MustCompile(`package\s*(\w+?);`)
var UsernameRe = regexp.MustCompile(`@Username\s*=\s*"(\w+?)";`)
var PasswordRe = regexp.MustCompile(`@Password\s*=\s*"(\w+?)";`)
var AddressRe = regexp.MustCompile(`@Address\s*=\s*"(.*?)";`)
var NameRe = regexp.MustCompile(`@Name\s*=\s*"(\w+?)";`)
var TableRe = regexp.MustCompile(`(?ms)Table\s+(\w+)\s+{(.+?)};`)
var ColumnRe = regexp.MustCompile(`Column\s+(\w+)\s+([\w\(\)]+)(.*?),`)
var NullableRe = regexp.MustCompile(`NOT NULL`)
var DefaultRe = regexp.MustCompile(`DEFAULT (".*?")`)
var AutoIncrementRe = regexp.MustCompile(`AUTOINCREMENT`)
var UniqueRe = regexp.MustCompile(`UNIQUE`)
var ForeignKeyRe = regexp.MustCompile(`ForeignKey\s+(\w+)\s+(\w+)\s+(\w+),`)
var ManyToManyRe = regexp.MustCompile(`ManyToMany\s+(\w+)\s+(\w+)\s+(\w+),`)
var PrimaryKeyRe = regexp.MustCompile(`PRIMARY KEY\s+(\w+),`)

func findPackage(s string) (string, error) {
	pkg := PackageRe.FindStringSubmatch(s)
	if len(pkg) < 2 {
		return "", errors.New("no package")
	}
	return pkg[1], nil
}

func findUsername(s string) (string, error) {
	username := UsernameRe.FindStringSubmatch(s)
	if len(username) < 2 {
		return "", errors.New("no database username")
	}
	return username[1], nil
}

func findPassword(s string) (string, error) {
	password := PasswordRe.FindStringSubmatch(s)
	if len(password) < 2 {
		return "", errors.New("no database password")
	}
	return password[1], nil
}

func findAddr(s string) (string, error) {
	addr := AddressRe.FindStringSubmatch(s)
	if len(addr) < 2 {
		return "", errors.New("no database address")
	}
	return addr[1], nil
}

func findName(s string) (string, error) {
	name := NameRe.FindStringSubmatch(s)
	if len(name) < 2 {
		return "", errors.New("no database name")
	}
	return name[1], nil
}

func findTables(s string) ([][]string, error) {
	tables := TableRe.FindAllStringSubmatch(s, -1)
	if len(tables) == 0 {
		return nil, errors.New("no table exists")
	}
	return tables, nil
}

func findColumns(s string) ([][]string, error) {
	columns := ColumnRe.FindAllStringSubmatch(s, -1)
	if len(columns) == 0 {
		return nil, errors.New("no column exists")
	}
	return columns, nil
}

func findForeignKeys(s string) [][]string {
	return ForeignKeyRe.FindAllStringSubmatch(s, -1)
}

func findManyToMany(s string) [][]string {
	return ManyToManyRe.FindAllStringSubmatch(s, -1)
}

func findPrimaryKey(s string) ([]string, error) {
	primaryKey := PrimaryKeyRe.FindStringSubmatch(s)
	if len(primaryKey) == 0 {
		return nil, errors.New("no primary key")
	}
	return primaryKey, nil
}

func findAutoIncrement(s string) []string {
	return AutoIncrementRe.FindStringSubmatch(s)
}

func ParseDatabase(file string) (Database, error) {
	db := Database{}
	f, err := os.Open(file)
	if err != nil {
		return db, err
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return db, err
	}
	s := string(b)
	pkg, err := findPackage(s)
	if err != nil {
		return db, err
	}
	username, err := findUsername(s)
	if err != nil {
		return db, err
	}
	password, err := findPassword(s)
	if err != nil {
		return db, err
	}
	addr, err := findAddr(s)
	if err != nil {
		return db, err
	}
	name, err := findName(s)
	if err != nil {
		return db, err
	}
	ts, err := findTables(s)
	if err != nil {
		return db, err
	}
	db.Package = pkg
	db.Username = username
	db.Password = password
	db.Address = addr
	db.DatabaseName = name
	db.ObjName = ToCap(name)
	for _, t := range ts {
		table, err := parseTable(t, db)
		if err != nil {
			return db, err
		}
		db.Tables = append(db.Tables, table)
	}
	for i, _ := range db.Tables {
		for _, info := range db.Tables[i].ForeignKeyInfos {
			err := parseForeignKey(info, &db.Tables[i], &db)
			if err != nil {
				return db, err
			}
		}
		for _, info := range db.Tables[i].ManyToManyInfos {
			err := parseManyToMany(info, &db.Tables[i], &db)
			if err != nil {
				return db, err
			}
		}
	}
	return db, nil
}

func parseTable(t []string, db Database) (Table, error) {
	table := Table{}
	table.TableName = t[1]
	table.ModelName = ToCap(t[1])
	table.ArgName = ToLocal(t[1])
	columns, err := findColumns(t[2])
	if err != nil {
		return table, err
	}
	for _, column := range columns {
		col := Column{}
		col.ColumnName = column[1]
		col.ArgName = ToLocal(column[1])
		col.FieldName = ToCap(column[1])
		col.MySqlType = column[2]
		col.FieldType = MySqlGoMap[RmParenthesis(column[2])]
		col.MidType = MySqlMidMap[RmParenthesis(column[2])]
		col.Nullable = !NullableRe.MatchString(column[3])
		col.AutoIncrement = AutoIncrementRe.MatchString(column[3])
		col.Unique = UniqueRe.MatchString(column[3])
		def := DefaultRe.FindStringSubmatch(column[3])
		if len(def) > 0 {
			col.Default = def[1]
		}
		table.Columns = append(table.Columns, col)
	}
	primaryKey, err := findPrimaryKey(t[2])
	if err != nil {
		return table, err
	}
	for i, _ := range table.Columns {
		if table.Columns[i].ColumnName == primaryKey[1] {
			table.PrimaryKey = &table.Columns[i]
		}
	}
	if table.PrimaryKey == nil {
		return table, errors.New("primary key column not exist")
	}
	for i, _ := range table.Columns {
		if table.Columns[i].AutoIncrement {
			table.AutoIncrement = &table.Columns[i]
			break
		}
	}
	foreignKeys := findForeignKeys(t[2])
	table.ForeignKeyInfos = make([]ForeignKeyInfo, len(foreignKeys))
	for i, fk := range foreignKeys {
		table.ForeignKeyInfos[i] = ForeignKeyInfo{
			DstTabName: fk[1],
			SrcColName: fk[2],
			DstColName: fk[3],
		}
	}
	manyToManys := findManyToMany(t[2])
	table.ManyToManyInfos = make([]ManyToManyInfo, len(manyToManys))
	for i, mtm := range manyToManys {
		table.ManyToManyInfos[i] = ManyToManyInfo{
			DstTabName: mtm[1],
			SrcColName: mtm[2],
			DstColName: mtm[3],
		}
	}
	return table, nil
}

func parseForeignKey(info ForeignKeyInfo, tab *Table, db *Database) error {
	var dstTab *Table
	var srcCol *Column
	var dstCol *Column
	for i, _ := range db.Tables {
		if db.Tables[i].TableName == info.DstTabName {
			dstTab = &db.Tables[i]
			for j, _ := range db.Tables[i].Columns {
				if db.Tables[i].Columns[j].ColumnName == info.DstColName {
					dstCol = &(db.Tables[i].Columns[j])
					break
				}
			}
		}
	}
	for i, _ := range tab.Columns {
		if tab.Columns[i].ColumnName == info.SrcColName {
			srcCol = &tab.Columns[i]
			break
		}
	}
	if dstTab == nil {
		return errors.New("destination table not exists")
	}
	if srcCol == nil {
		return errors.New("source column not exists")
	}
	if dstCol == nil {
		return errors.New("destination column not exists")
	}
	if tab.ForeignKeys == nil {
		tab.ForeignKeys = make([]ForeignKey, 0, 8)
	}
	tab.ForeignKeys = append(tab.ForeignKeys, ForeignKey{
		DstTab: dstTab,
		SrcCol: srcCol,
		DstCol: dstCol,
	})
	return nil
}

func parseManyToMany(info ManyToManyInfo, tab *Table, db *Database) error {
	var srcCol *Column
	var dstCol *Column
	var dstTab *Table
	var midLeftCol *Column
	var midRightCol *Column

	for i, _ := range tab.Columns {
		if tab.Columns[i].ColumnName == info.SrcColName {
			srcCol = &tab.Columns[i]
			break
		}
	}
	if srcCol == nil {
		return errors.New("source column not exists")
	}

	for i, _ := range db.Tables {
		if db.Tables[i].TableName == info.DstTabName {
			dstTab = &db.Tables[i]
			break
		}
	}
	if dstTab == nil {
		return errors.New("destination table not exists")
	}

	for i, _ := range dstTab.Columns {
		if dstTab.Columns[i].ColumnName == info.DstColName {
			dstCol = &dstTab.Columns[i]
		}
	}
	if dstCol == nil {
		return errors.New("destination column not exists")
	}

	midTab := Table{
		TableName: tab.TableName + "__" + dstTab.TableName,
		Columns: []Column{
			Column{ColumnName: "id", MySqlType: "INT", AutoIncrement: true},
			Column{ColumnName: tab.TableName + "__" + info.SrcColName, MySqlType: srcCol.MySqlType},
			Column{ColumnName: info.DstTabName + "__" + info.DstColName, MySqlType: dstCol.MySqlType},
		},
	}
	midTab.AutoIncrement = &midTab.Columns[0]
	midTab.PrimaryKey = &midTab.Columns[0]
	midTab.UniqueKeys = []UniqueKey{
		UniqueKey{midTab.Columns[1].ColumnName, midTab.Columns[2].ColumnName},
	}
	err := db.CreateTableIfNotExists(midTab)
	if err != nil {
		return err
	}
	db.MidTables = append(db.MidTables, midTab)
	midLeftCol = &midTab.Columns[1]
	midRightCol = &midTab.Columns[2]
	mtm := ManyToMany{
		SrcCol:      srcCol,
		MidLeftCol:  midLeftCol,
		MidRightCol: midRightCol,
		DstCol:      dstCol,
		MidTab:      &midTab,
		DstTab:      dstTab,
	}
	dstMtm := ManyToMany{
		SrcCol:      dstCol,
		MidLeftCol:  midRightCol,
		MidRightCol: midLeftCol,
		DstCol:      srcCol,
		MidTab:      &midTab,
		DstTab:      tab,
	}
	tab.ManyToManys = append(tab.ManyToManys, mtm)
	dstTab.ManyToManys = append(dstTab.ManyToManys, dstMtm)
	return nil
}
