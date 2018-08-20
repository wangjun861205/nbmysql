package nbmysql

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

//ForeignKey define a foreign key relation
type ForeignKey struct {
	SrcCol *Column
	DstCol *Column
	DstTab *Table
}

//ReverseForeignKey define a reverse foreign key relation. Note: Database parser will automatic define ReverseForeignKey correspond to source table
//Foreignkey in destination table, you should not define it in database definition file manually
type ReverseForeignKey struct {
	SrcCol *Column
	DstCol *Column
	DstTab *Table
}

//ManyToMany define a many to many relation
type ManyToMany struct {
	SrcCol      *Column
	MidLeftCol  *Column
	MidRightCol *Column
	DstCol      *Column
	MidTab      *Table
	DstTab      *Table
}

//ForeignKeyInfo intermediate struct for parse foreign key relation definition
type ForeignKeyInfo struct {
	SrcColName string
	DstColName string
	DstTabName string
}

//ManyToManyInfo intermediate struct for parse many to many relation definition
type ManyToManyInfo struct {
	SrcColName string
	DstColName string
	DstTabName string
}

//GoMidMap intermediate table for many to many relation
var GoMidMap = map[string]string{
	"int":    "Int",
	"float":  "Float",
	"string": "String",
	"bool":   "Bool",
	"time":   "Time",
}

// var GoArgMap = map[string]string{
// 	"int64":     "complex128",
// 	"float64":   "complex128",
// 	"string":    "[]byte",
// 	"bool":      "int",
// 	"time.Time": "time.Time",
// }

//MySqlGoMap map MySQL type to go type
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

//MySqlMidMap map MySQL type to middle type
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

//PackageRe regexp pattern for finding package name
var packageRe = regexp.MustCompile(`package\s*(\w+?);`)

//UsernameRe regexp pattern for finding database server username
var usernameRe = regexp.MustCompile(`@Username\s*=\s*"(\w+?)";`)

//PasswordRe regexp pattern for finding database server password
var passwordRe = regexp.MustCompile(`@Password\s*=\s*"(\w+?)";`)

//AddressRe regexp pattern for finding database server address
var addressRe = regexp.MustCompile(`@Address\s*=\s*"(.*?)";`)

//NameRe regexp pattern for finding database name
var nameRe = regexp.MustCompile(`@Name\s*=\s*"(\w+?)";`)

//TableRe regexp pattern for finding table definition
var tableRe = regexp.MustCompile(`(?ms)Table\s+(\w+)\s+{(.+?)};`)

//ColumnRe regexp pattern for finding column definition in table
var columnRe = regexp.MustCompile(`Column\s+(\w+)\s+([\w\(\)]+)(.*?),`)

//NullableRe regexp pattern for finding nullable attribute in column
var nullableRe = regexp.MustCompile(`NOT NULL`)

//DefaultRe regexp pattern for finding default attribute in column
var defaultRe = regexp.MustCompile(`DEFAULT\s+(['"].*['"]|[^\s]+)`)

//AutoIncrementRe regexp pattern for finding auto incremtnt attribute in column
var autoIncrementRe = regexp.MustCompile(`AUTO_INCREMENT`)

//UniqueRe regexp pattern for finding unique attribute in column
var uniqueRe = regexp.MustCompile(`UNIQUE`)

//ForeignKeyRe regexp pattern for finding foreign key relation definition in table
var foreignKeyRe = regexp.MustCompile(`ForeignKey\s+(\w+)\s+(\w+)\s+(\w+),`)

//ManyToManyRe regexp pattern for finding many to many relation definition in table
var manyToManyRe = regexp.MustCompile(`ManyToMany\s+(\w+)\s+(\w+)\s+(\w+),`)

//PrimaryKeyRe regexp pattern for finding primary key column in table
var primaryKeyRe = regexp.MustCompile(`PRIMARY KEY\s+(.*?),`)

//UniqueKeyRe regexp pattern for finding unique key columns in table
var uniqueKeyRe = regexp.MustCompile(`UNIQUE KEY\s+\((.*?)\),`)

//OnRe regexp pattern for finding on condition in table
var onRe = regexp.MustCompile(`ON\s+(.*)`)

func findPackage(s string) (string, error) {
	pkg := packageRe.FindStringSubmatch(s)
	if len(pkg) < 2 {
		return "", errors.New("no package")
	}
	return pkg[1], nil
}

func findUsername(s string) (string, error) {
	username := usernameRe.FindStringSubmatch(s)
	if len(username) < 2 {
		return "", errors.New("no database username")
	}
	return username[1], nil
}

func findPassword(s string) (string, error) {
	password := passwordRe.FindStringSubmatch(s)
	if len(password) < 2 {
		return "", errors.New("no database password")
	}
	return password[1], nil
}

func findAddr(s string) (string, error) {
	addr := addressRe.FindStringSubmatch(s)
	if len(addr) < 2 {
		return "", errors.New("no database address")
	}
	return addr[1], nil
}

func findName(s string) (string, error) {
	name := nameRe.FindStringSubmatch(s)
	if len(name) < 2 {
		return "", errors.New("no database name")
	}
	return name[1], nil
}

func findTables(s string) ([][]string, error) {
	tables := tableRe.FindAllStringSubmatch(s, -1)
	if len(tables) == 0 {
		return nil, errors.New("no table exists")
	}
	return tables, nil
}

func findColumns(s string) ([][]string, error) {
	columns := columnRe.FindAllStringSubmatch(s, -1)
	if len(columns) == 0 {
		return nil, errors.New("no column exists")
	}
	return columns, nil
}

func findForeignKeys(s string) [][]string {
	return foreignKeyRe.FindAllStringSubmatch(s, -1)
}

func findManyToMany(s string) [][]string {
	return manyToManyRe.FindAllStringSubmatch(s, -1)
}

func findPrimaryKey(s string) ([]string, error) {
	primaryKey := primaryKeyRe.FindStringSubmatch(s)
	if len(primaryKey) == 0 {
		return nil, errors.New("no primary key")
	}
	return primaryKey, nil
}

func findAutoIncrement(s string) []string {
	return autoIncrementRe.FindStringSubmatch(s)
}

//ParseDatabase parse database definition file and generate database info struct
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
	for i := range db.Tables {
		for _, info := range db.Tables[i].ForeignKeyInfos {
			err := parseForeignKey(info, &db.Tables[i], &db)
			if err != nil {
				return db, err
			}
			err = parseReverseForeignKey(info, &db.Tables[i], &db)
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
		col.Nullable = !nullableRe.MatchString(column[3])
		col.AutoIncrement = autoIncrementRe.MatchString(column[3])
		col.Unique = uniqueRe.MatchString(column[3])
		def := defaultRe.FindStringSubmatch(column[3])
		if len(def) > 0 {
			col.Default = def[1]
		}
		on := onRe.FindStringSubmatch(column[3])
		if len(on) > 0 {
			col.On = on[1]
		}
		table.Columns = append(table.Columns, col)
	}
	primaryKey, err := findPrimaryKey(t[2])
	if err != nil {
		return table, err
	}
	for i := range table.Columns {
		if table.Columns[i].ColumnName == primaryKey[1] {
			table.PrimaryKey = &table.Columns[i]
		}
	}
	if table.PrimaryKey == nil {
		return table, errors.New("primary key column not exist")
	}
	for i := range table.Columns {
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
	uniqueKeys, err := parseUniqueKeys(t[2], &table)
	if err != nil {
		return table, err
	}
	table.UniqueKeys = uniqueKeys
	return table, nil
}

func parseForeignKey(info ForeignKeyInfo, tab *Table, db *Database) error {
	var dstTab *Table
	var srcCol *Column
	var dstCol *Column
	for i := range db.Tables {
		if db.Tables[i].TableName == info.DstTabName {
			dstTab = &db.Tables[i]
			for j := range db.Tables[i].Columns {
				if db.Tables[i].Columns[j].ColumnName == info.DstColName {
					dstCol = &(db.Tables[i].Columns[j])
					break
				}
			}
		}
	}
	for i := range tab.Columns {
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

func parseReverseForeignKey(info ForeignKeyInfo, srcTab *Table, db *Database) error {
	var dstTab *Table
	for i, tab := range db.Tables {
		if tab.TableName == info.DstTabName {
			dstTab = &(db.Tables[i])
			break
		}
	}
	if dstTab == nil {
		return fmt.Errorf("parse reverse forkeign key error: %s table not exsits", info.DstTabName)
	}
	var srcCol *Column
	for i, col := range srcTab.Columns {
		if col.ColumnName == info.SrcColName {
			srcCol = &(srcTab.Columns[i])
			break
		}
	}
	if srcCol == nil {
		return fmt.Errorf("parse reverse foreign key error: %s source column not exists", info.SrcColName)
	}
	var dstCol *Column
	for i, col := range dstTab.Columns {
		if col.ColumnName == info.SrcColName {
			dstCol = &(dstTab.Columns[i])
			break
		}
	}
	if dstCol == nil {
		return fmt.Errorf("parse reverse foreign key error: %s destination column not exists", info.DstColName)
	}
	if dstTab.ReverseForeignKeys == nil {
		dstTab.ReverseForeignKeys = make([]ReverseForeignKey, 0, 8)
	}
	rfk := ReverseForeignKey{DstTab: srcTab, SrcCol: dstCol, DstCol: srcCol}
	dstTab.ReverseForeignKeys = append(dstTab.ReverseForeignKeys, rfk)
	return nil
}

func parseManyToMany(info ManyToManyInfo, tab *Table, db *Database) error {
	var srcCol *Column
	var dstCol *Column
	var dstTab *Table
	var midLeftCol *Column
	var midRightCol *Column

	for i := range tab.Columns {
		if tab.Columns[i].ColumnName == info.SrcColName {
			srcCol = &tab.Columns[i]
			break
		}
	}
	if srcCol == nil {
		return errors.New("source column not exists")
	}

	for i := range db.Tables {
		if db.Tables[i].TableName == info.DstTabName {
			dstTab = &db.Tables[i]
			break
		}
	}
	if dstTab == nil {
		return errors.New("destination table not exists")
	}

	for i := range dstTab.Columns {
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
		UniqueKey{&midTab.Columns[1], &midTab.Columns[2]},
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

func parseUniqueKeys(s string, tab *Table) ([]UniqueKey, error) {
	l := uniqueKeyRe.FindAllStringSubmatch(s, -1)
	ukList := make([]UniqueKey, len(l))
	for i, uk := range l {
		colList := strings.Split(uk[1], ",")
		uniqueKey := make(UniqueKey, len(colList))
		for j, colName := range colList {
			col := getColumn(strings.Trim(colName, " "), tab)
			if col == nil {
				return nil, fmt.Errorf("unique key %s not exists", strings.Trim(colName, " "))
			}
			uniqueKey[j] = col
		}
		ukList[i] = uniqueKey
	}
	return ukList, nil
}
