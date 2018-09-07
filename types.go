package nbmysql

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
)

//UniqueKey *Column slice for unique keys
type UniqueKey []*Column

//Column define column
type Column struct {
	FieldName     string
	FieldType     string
	ArgName       string
	ColumnName    string
	MidType       string
	MySqlType     string
	Nullable      bool
	Default       string
	Unique        bool
	AutoIncrement bool
	On            string
}

//Table define table
type Table struct {
	Columns            []Column
	ModelName          string
	TableName          string
	ArgName            string
	PrimaryKey         *Column
	AutoIncrement      *Column
	ForeignKeyInfos    []ForeignKeyInfo
	ManyToManyInfos    []ManyToManyInfo
	ForeignKeys        []ForeignKey
	ReverseForeignKeys []ReverseForeignKey
	ManyToManys        []ManyToMany
	UniqueKeys         []UniqueKey
}

//Database define database
type Database struct {
	Package      string
	ObjName      string
	Username     string
	Password     string
	Address      string
	DatabaseName string
	Tables       []Table
	MidTables    []Table
}

//CreateTableIfNotExists create table by table struct which is not exist in database
func (db *Database) CreateTableIfNotExists(tab Table) error {
	conn, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", db.Username, db.Password, db.Address, db.DatabaseName))
	if err != nil {
		return err
	}
	defer conn.Close()
	colList := make([]string, len(tab.Columns))
	for i, col := range tab.Columns {
		l := make([]string, 0, 6)
		l = append(l, col.ColumnName)
		l = append(l, col.MySqlType)
		if col.AutoIncrement {
			l = append(l, "AUTO_INCREMENT")
		}
		if !col.Nullable {
			l = append(l, "NOT NULL")
		}
		if col.Unique {
			l = append(l, "UNIQUE")
		}
		if col.Default != "" {
			l = append(l, "DEFAULT "+col.Default)
		}
		if col.On != "" {
			l = append(l, "ON "+col.On)
		}
		colList[i] = strings.Join(l, " ")
	}
	uniqueList := make([]string, 0, len(tab.UniqueKeys))
	for _, uni := range tab.UniqueKeys {
		bqList := make([]string, len(uni))
		for i := range uni {
			bqList[i] = string(uni[i].ColumnName)
		}
		ukColList := make([]string, len(uni))
		for i, col := range uni {
			ukColList[i] = strings.Trim(col.ColumnName, "`")
		}
		uniqueList = append(uniqueList, fmt.Sprintf("UNIQUE KEY `%s_unique` (%s)", strings.Join(ukColList, "_"), strings.Join(bqList, ", ")))
	}
	var uniqueClause string
	if len(uniqueList) > 0 {
		uniqueClause = ", " + strings.Join(uniqueList, ", ")
	}
	stmt := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s, PRIMARY KEY (%s)%s)", tab.TableName, strings.Join(colList, ", "),
		tab.PrimaryKey.ColumnName, uniqueClause)
	_, err = conn.Exec(stmt)
	return err
}

//AddForeignKeyConstraint add foreign key constraint to table
func (db *Database) AddForeignKeyConstraint() error {
	conn, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", db.Username, db.Password, db.Address, db.DatabaseName))
	if err != nil {
		return err
	}
	defer conn.Close()
	for _, tab := range db.Tables {
		for _, fk := range tab.ForeignKeys {
			_, err := conn.Exec(
				fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s (%s) ON DELETE CASCADE ON UPDATE CASCADE",
					tab.TableName,
					BackQuote("fk_"+strings.Trim(fk.DstTab.TableName, "`")+"__"+strings.Trim(fk.DstCol.ColumnName, "`")),
					fk.SrcCol.ColumnName,
					fk.DstTab.TableName,
					fk.DstCol.ColumnName))
			if err != nil {
				if sqlErr, ok := err.(*mysql.MySQLError); ok && sqlErr.Number == 1826 {
					log.Printf("warnning: %s", sqlErr.Error())
					continue
				}
				return err
			}
		}
		for _, mtm := range tab.ManyToManys {
			_, err := conn.Exec(
				fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s (%s) ON DELETE CASCADE ON UPDATE CASCADE",
					mtm.MidTab.TableName,
					BackQuote("midfk_"+strings.Trim(tab.TableName, "`")+"__"+strings.Trim(mtm.SrcCol.ColumnName, "`")),
					mtm.MidLeftCol.ColumnName,
					tab.TableName,
					mtm.SrcCol.ColumnName))
			if err != nil {
				if sqlErr, ok := err.(*mysql.MySQLError); ok && sqlErr.Number == 1826 {
					log.Printf("warnning: %s", sqlErr.Error())
					continue
				}
				return err
			}
		}
	}
	return nil
}

type StringField struct {
	val   string
	valid bool
	null  bool
}

func (f *StringField) IsValid() bool {
	return f.valid
}

func (f *StringField) IsNull() bool {
	return f.null
}

func (f *StringField) Set(val string, nullAndValid ...bool) {
	var valid bool
	var null bool
	switch len(nullAndValid) {
	case 0:
		valid = true
		null = false
	case 1:
		valid = true
		null = nullAndValid[0]
	case 2:
		null, valid = nullAndValid[0], nullAndValid[1]
	default:
		panic(fmt.Sprintf("invalid length of nullAndValid arguments in StringField.Set(): required 2 supported %d", len(nullAndValid)))
	}
	f.val = val
	f.valid = valid
	f.null = null
}

func (f *StringField) Get() (val string, valid bool, null bool) {
	return f.val, f.valid, f.null
}

func (f *StringField) Scan(v interface{}) error {
	f.valid = true
	if v == nil {
		f.null = true
		return nil
	}
	f.null = false
	switch val := v.(type) {
	case []byte:
		f.val = string(val)
	case string:
		f.val = val
	default:
		return fmt.Errorf("not supported value type for StringField: %T", val)
	}
	return nil
}

func (f *StringField) SQLVal() string {
	return fmt.Sprintf("%q", f.val)
}

func (f *StringField) Invalidate() {
	f.valid = false
}

func (f *StringField) MarshalJSON() ([]byte, error) {
	if !f.valid {
		return []byte("invalid"), nil
	}
	if f.null {
		return []byte("NULL"), nil
	}
	return []byte(fmt.Sprintf("%q", f.val)), nil
}

type IntField struct {
	val   int64
	valid bool
	null  bool
}

func (f *IntField) IsValid() bool {
	return f.valid
}

func (f *IntField) IsNull() bool {
	return f.null
}

func (f *IntField) Set(val int64, nullAndValid ...bool) {
	var valid bool
	var null bool
	switch len(nullAndValid) {
	case 0:
		valid = true
		null = false
	case 1:
		valid = true
		null = nullAndValid[0]
	case 2:
		null, valid = nullAndValid[0], nullAndValid[1]
	default:
		panic(fmt.Sprintf("invalid length of nullAndValid arguments in IntField.Set(): required 2 supported %d", len(nullAndValid)))
	}
	f.val = val
	f.valid = valid
	f.null = null
}

func (f *IntField) Get() (val int64, valid bool, null bool) {
	return f.val, f.valid, f.null
}

func (f *IntField) Scan(v interface{}) error {
	f.valid = true
	if v == nil {
		f.null = true
		return nil
	}
	f.null = false
	switch val := v.(type) {
	case []byte:
		i64, err := strconv.ParseInt(string(val), 10, 64)
		if err != nil {
			return err
		}
		f.val = i64
	case int64:
		f.val = val
	default:
		return fmt.Errorf("not supported value type for IntField: %T", val)
	}
	return nil
}

func (f *IntField) SQLVal() string {
	return fmt.Sprintf("%d", f.val)
}

func (f *IntField) Invalidate() {
	f.valid = false
}

func (f *IntField) MarshalJSON() ([]byte, error) {
	if !f.valid {
		return []byte("invalid"), nil
	}
	if f.null {
		return []byte("NULL"), nil
	}
	return []byte(strconv.FormatInt(f.val, 10)), nil
}

type FloatField struct {
	val   float64
	valid bool
	null  bool
}

func (f *FloatField) IsValid() bool {
	return f.valid
}

func (f *FloatField) IsNull() bool {
	return f.null
}

func (f *FloatField) Set(val float64, nullAndValid ...bool) {
	var valid bool
	var null bool
	switch len(nullAndValid) {
	case 0:
		valid = true
		null = false
	case 1:
		valid = true
		null = nullAndValid[0]
	case 2:
		null, valid = nullAndValid[0], nullAndValid[1]
	default:
		panic(fmt.Sprintf("invalid length of nullAndValid arguments in FloatField.Set(): required 2 supported %d", len(nullAndValid)))
	}
	f.val = val
	f.valid = valid
	f.null = null
}

func (f *FloatField) Get() (val float64, valid bool, null bool) {
	return f.val, f.valid, f.null
}

func (f *FloatField) Scan(v interface{}) error {
	f.valid = true
	if v == nil {
		f.null = true
		return nil
	}
	f.null = false
	switch val := v.(type) {
	case []byte:
		f64, err := strconv.ParseFloat(string(val), 64)
		if err != nil {
			return err
		}
		f.val = f64
	case float64:
		f.val = val
	default:
		return fmt.Errorf("not supported value type for IntField: %T", val)
	}
	return nil
}

func (f *FloatField) SQLVal() string {
	return fmt.Sprintf("%f", f.val)
}

func (f *FloatField) Invalidate() {
	f.valid = false
}

func (f *FloatField) MarshalJSON() ([]byte, error) {
	if !f.valid {
		return []byte("invalid"), nil
	}
	if f.null {
		return []byte("NULL"), nil
	}
	return []byte(strconv.FormatFloat(f.val, 'f', -1, 64)), nil
}

type BoolField struct {
	val   bool
	valid bool
	null  bool
}

func (f *BoolField) IsValid() bool {
	return f.valid
}

func (f *BoolField) IsNull() bool {
	return f.null
}

func (f *BoolField) Set(val bool, nullAndValid ...bool) {
	var valid bool
	var null bool
	switch len(nullAndValid) {
	case 0:
		valid = true
		null = false
	case 1:
		valid = true
		null = nullAndValid[0]
	case 2:
		null, valid = nullAndValid[0], nullAndValid[1]
	default:
		panic(fmt.Sprintf("invalid length of nullAndValid arguments in BoolField.Set(): required 2 supported %d", len(nullAndValid)))
	}
	f.val = val
	f.valid = valid
	f.null = null
}

func (f *BoolField) Get() (val bool, valid bool, null bool) {
	return f.val, f.valid, f.null
}

func (f *BoolField) Scan(v interface{}) error {
	f.valid = true
	if v == nil {
		f.null = true
		return nil
	}
	f.null = false
	switch val := v.(type) {
	case []byte:
		b, err := strconv.ParseBool(string(val))
		if err != nil {
			return err
		}
		f.val = b
	case bool:
		f.val = val
	default:
		return fmt.Errorf("not supported value type for BoolField: %T", val)
	}
	return nil
}

func (f *BoolField) SQLVal() string {
	return fmt.Sprintf("%t", f.val)
}

func (f *BoolField) Invalidate() {
	f.valid = false
}

func (f *BoolField) MarshalJSON() ([]byte, error) {
	if !f.valid {
		return []byte("invalid"), nil
	}
	if f.null {
		return []byte("NULL"), nil
	}
	return []byte(strconv.FormatBool(f.val)), nil
}

type DateField struct {
	val   time.Time
	valid bool
	null  bool
}

func (f *DateField) IsValid() bool {
	return f.valid
}

func (f *DateField) IsNull() bool {
	return f.null
}

func (f *DateField) Set(val time.Time, nullAndValid ...bool) {
	var valid bool
	var null bool
	switch len(nullAndValid) {
	case 0:
		valid = true
		null = false
	case 1:
		valid = true
		null = nullAndValid[0]
	case 2:
		null, valid = nullAndValid[0], nullAndValid[1]
	default:
		panic(fmt.Sprintf("invalid length of nullAndValid arguments in DateField.Set(): required 2 supported %d", len(nullAndValid)))
	}
	f.val = val
	f.valid = valid
	f.null = null
}

func (f *DateField) Get() (val time.Time, valid bool, null bool) {
	return f.val, f.valid, f.null
}

func (f *DateField) Scan(v interface{}) error {
	f.valid = true
	if v == nil {
		f.null = true
		return nil
	}
	f.null = false
	switch val := v.(type) {
	case []byte:
		t, err := time.Parse("2006-01-02", string(val))
		if err != nil {
			return err
		}
		f.val = t
	case time.Time:
		f.val = val
	default:
		return fmt.Errorf("not supported value type for DateField: %T", val)
	}
	return nil
}

func (f *DateField) SQLVal() string {
	return f.val.Format("2006-01-02")
}

func (f *DateField) Invalidate() {
	f.valid = false
}

func (f *DateField) MarshalJSON() ([]byte, error) {
	if !f.valid {
		return []byte("invalid"), nil
	}
	if f.null {
		return []byte("NULL"), nil
	}
	return []byte(fmt.Sprintf("%q", f.val.Format("2006-01-02"))), nil
}

type DatetimeField struct {
	val   time.Time
	valid bool
	null  bool
}

func (f *DatetimeField) IsValid() bool {
	return f.valid
}

func (f *DatetimeField) IsNull() bool {
	return f.null
}

func (f *DatetimeField) Set(val time.Time, nullAndValid ...bool) {
	var valid bool
	var null bool
	switch len(nullAndValid) {
	case 0:
		valid = true
		null = false
	case 1:
		valid = true
		null = nullAndValid[0]
	case 2:
		null, valid = nullAndValid[0], nullAndValid[1]
	default:
		panic(fmt.Sprintf("invalid length of nullAndValid arguments in DatetimeField.Set(): required 2 supported %d", len(nullAndValid)))
	}
	f.val = val
	f.valid = valid
	f.null = null
}

func (f *DatetimeField) Get() (val time.Time, valid bool, null bool) {
	return f.val, f.valid, f.null
}

func (f *DatetimeField) Scan(v interface{}) error {
	f.valid = true
	if v == nil {
		f.null = true
		return nil
	}
	f.null = false
	switch val := v.(type) {
	case []byte:
		t, err := time.Parse("2006-01-02 15:04:05", string(val))
		if err != nil {
			return err
		}
		f.val = t
	case time.Time:
		f.val = val
	default:
		return fmt.Errorf("not supported value type for DatetimeField: %T", val)
	}
	return nil
}

func (f *DatetimeField) SQLVal() string {
	return f.val.Format("2006-01-02 15:04:05")
}

func (f *DatetimeField) Invalidate() {
	f.valid = false
}

func (f *DatetimeField) MarshalJSON() ([]byte, error) {
	if !f.valid {
		return []byte("invalid"), nil
	}
	if f.null {
		return []byte("NULL"), nil
	}
	return []byte(fmt.Sprintf("%q", f.val.Format("2006-01-02 15:04:05"))), nil
}
