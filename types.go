package nbmysql

import (
	"database/sql"
	"encoding/json"
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

type Field interface {
	IsValid() bool
	IsNull() bool
	sql.Scanner
	SQLVal() string
	Invalidate()
	json.Marshaler
	InsertValuePair() [2]string
	UpdateValue() string
	Where() *Where
}

type StringField struct {
	val   string
	valid bool
	null  bool
	class *StringFieldClass
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
		return []byte("\"invalid\""), nil
	}
	if f.null {
		return []byte("\"NULL\""), nil
	}
	return []byte(fmt.Sprintf("%q", f.val)), nil
}

func (s *StringField) InsertValuePair() [2]string {
	return [2]string{s.class.columnName, s.SQLVal()}
}

func (s *StringField) UpdateValue() string {
	return s.class.columnName + " = " + s.SQLVal()
}

func (s *StringField) Where() *Where {
	return &Where{s.class.columnName + " = " + s.SQLVal()}
}

type IntField struct {
	val   int64
	valid bool
	null  bool
	class *IntFieldClass
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
		return []byte("\"invalid\""), nil
	}
	if f.null {
		return []byte("\"NULL\""), nil
	}
	return []byte(strconv.FormatInt(f.val, 10)), nil
}

func (s *IntField) InsertValuePair() [2]string {
	return [2]string{s.class.columnName, s.SQLVal()}
}

func (s *IntField) UpdateValue() string {
	return s.class.columnName + " = " + s.SQLVal()
}

func (s *IntField) Where() *Where {
	return &Where{s.class.columnName + " = " + s.SQLVal()}
}

type FloatField struct {
	val   float64
	valid bool
	null  bool
	class *FloatFieldClass
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
		return []byte("\"invalid\""), nil
	}
	if f.null {
		return []byte("\"NULL\""), nil
	}
	return []byte(strconv.FormatFloat(f.val, 'f', -1, 64)), nil
}

func (s *FloatField) InsertValuePair() [2]string {
	return [2]string{s.class.columnName, s.SQLVal()}
}

func (s *FloatField) UpdateValue() string {
	return s.class.columnName + " = " + s.SQLVal()
}

func (s *FloatField) Where() *Where {
	return &Where{s.class.columnName + " = " + s.SQLVal()}
}

type BoolField struct {
	val   bool
	valid bool
	null  bool
	class *BoolFieldClass
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
		return []byte("\"invalid\""), nil
	}
	if f.null {
		return []byte("\"NULL\""), nil
	}
	return []byte(strconv.FormatBool(f.val)), nil
}

func (s *BoolField) InsertValuePair() [2]string {
	return [2]string{s.class.columnName, s.SQLVal()}
}

func (s *BoolField) UpdateValue() string {
	return s.class.columnName + " = " + s.SQLVal()
}

func (s *BoolField) Where() *Where {
	return &Where{s.class.columnName + " = " + s.SQLVal()}
}

type DateField struct {
	val   time.Time
	valid bool
	null  bool
	class *DateFieldClass
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
		return []byte("\"invalid\""), nil
	}
	if f.null {
		return []byte("\"NULL\""), nil
	}
	return []byte(fmt.Sprintf("%q", f.val.Format("2006-01-02"))), nil
}

func (s *DateField) InsertValuePair() [2]string {
	return [2]string{s.class.columnName, s.SQLVal()}
}

func (s *DateField) UpdateValue() string {
	return s.class.columnName + " = " + s.SQLVal()
}

func (s *DateField) Where() *Where {
	return &Where{s.class.columnName + " = " + s.SQLVal()}
}

type DatetimeField struct {
	val   time.Time
	valid bool
	null  bool
	class *DatetimeFieldClass
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
		return []byte("\"invalid\""), nil
	}
	if f.null {
		return []byte("\"NULL\""), nil
	}
	return []byte(fmt.Sprintf("%q", f.val.Format("2006-01-02 15:04:05"))), nil
}

func (s *DatetimeField) InsertValuePair() [2]string {
	return [2]string{s.class.columnName, s.SQLVal()}
}

func (s *DatetimeField) UpdateValue() string {
	return s.class.columnName + " = " + s.SQLVal()
}

func (s *DatetimeField) Where() *Where {
	return &Where{s.class.columnName + " = " + s.SQLVal()}
}

type SetStruct struct {
	ModelName string
	FieldName string
	Func      func(Field)
}

type FieldClass interface {
	LessFunc() func(im, jm ModelInstance) (int, error)
	DistFunc() func(m ModelInstance) (string, error)
}

type StringFieldClass struct {
	tableName  string
	columnName string
	modelName  string
	fieldName  string
}

func NewStringFieldClass(tableName, columnName, modelName, fieldName string) *StringFieldClass {
	return &StringFieldClass{tableName, columnName, modelName, fieldName}
}

func (c *StringFieldClass) NewInstance() StringField {
	return StringField{class: c}
}

func (c *StringFieldClass) New(val string, isNull bool) [2]string {
	l := [2]string{}
	l[0] = fmt.Sprintf("%s.%s", c.tableName, c.columnName)
	if isNull {
		l[1] = "NULL"
	} else {
		l[1] = val
	}
	return l
}

func (c *StringFieldClass) Set(val string, isNull bool) SetStruct {
	return SetStruct{
		c.modelName,
		c.fieldName,
		func(f Field) {
			sf := f.(*StringField)
			sf.Set(val, isNull)
		},
	}
}

func (c *StringFieldClass) Eq(val string) *Where {
	return &Where{fmt.Sprintf("%s.%s = %q", c.tableName, c.columnName, val)}
}

func (c *StringFieldClass) Neq(val string) *Where {
	return &Where{fmt.Sprintf("%s.%s <> %q", c.tableName, c.columnName, val)}
}

func (c *StringFieldClass) Contains(val string) *Where {
	return &Where{fmt.Sprintf("%s.%s LIKE \"%%%s%%\"", c.tableName, c.columnName, val)}
}

func (c *StringFieldClass) IsNull() *Where {
	return &Where{fmt.Sprintf("%s.%s IS NULL", c.tableName, c.columnName)}
}

func (c *StringFieldClass) LessFunc() func(ModelInstance, ModelInstance) (int, error) {
	return func(mi, mj ModelInstance) (int, error) {
		fi, err := mi.GetField(c.fieldName)
		if err != nil {
			return 0, err
		}
		fj, err := mj.GetField(c.fieldName)
		if err != nil {
			return 0, err
		}
		iv, iValid, iNull := fi.(*StringField).Get()
		jv, jValid, jNull := fj.(*StringField).Get()
		if iValid && jValid {
			if !(iNull || jNull) {
				switch {
				case iv < jv:
					return -1, nil
				case iv > jv:
					return 1, nil
				default:
					return 0, nil
				}
			} else {
				switch {
				case iNull && jNull:
					return 0, nil
				case iNull:
					return -1, nil
				default:
					return 1, nil
				}
			}
		} else {
			switch {
			case !(iValid || jValid):
				return 0, nil
			case iValid:
				return 1, nil
			default:
				return -1, nil
			}
		}
	}
}

func (c *StringFieldClass) DistFunc() func(ModelInstance) (string, error) {
	return func(instance ModelInstance) (string, error) {
		field, err := instance.GetField(c.fieldName)
		if err != nil {
			return "", err
		}
		return field.SQLVal(), nil
	}
}

type IntFieldClass struct {
	tableName  string
	columnName string
	modelName  string
	fieldName  string
}

func NewIntFieldClass(tableName, columnName, modelName, fieldName string) *IntFieldClass {
	return &IntFieldClass{tableName, columnName, modelName, fieldName}
}

func (c *IntFieldClass) NewInstance() IntField {
	return IntField{class: c}
}

func (c *IntFieldClass) New(val int64, isNull bool) [2]string {
	l := [2]string{}
	l[0] = fmt.Sprintf("%s.%s", c.tableName, c.columnName)
	if isNull {
		l[1] = "NULL"
	} else {
		l[1] = fmt.Sprintf("%d", val)
	}
	return l
}

func (c *IntFieldClass) Set(val int64, isNull bool) SetStruct {
	return SetStruct{
		c.modelName,
		c.fieldName,
		func(f Field) {
			nf := f.(*IntField)
			nf.Set(val, isNull)
		},
	}
}

func (c *IntFieldClass) Eq(val int64) *Where {
	return &Where{fmt.Sprintf("%s.%s = %d", c.tableName, c.columnName, val)}
}

func (c *IntFieldClass) Neq(val int64) *Where {
	return &Where{fmt.Sprintf("%s.%s <> %d", c.tableName, c.columnName, val)}
}

func (c *IntFieldClass) Lt(val int64) *Where {
	return &Where{fmt.Sprintf("%s.%s < %d", c.tableName, c.columnName, val)}
}

func (c *IntFieldClass) Gt(val int64) *Where {
	return &Where{fmt.Sprintf("%s.%s > %d", c.tableName, c.columnName, val)}
}

func (c *IntFieldClass) Lte(val int64) *Where {
	return &Where{fmt.Sprintf("%s.%s <= %d", c.tableName, c.columnName, val)}
}

func (c *IntFieldClass) Gte(val int64) *Where {
	return &Where{fmt.Sprintf("%s.%s >= %d", c.tableName, c.columnName, val)}
}

func (c *IntFieldClass) IsNull() *Where {
	return &Where{fmt.Sprintf("%s.%s IS NULL", c.tableName, c.columnName)}
}

func (c *IntFieldClass) LessFunc() func(ModelInstance, ModelInstance) (int, error) {
	return func(mi, mj ModelInstance) (int, error) {
		fi, err := mi.GetField(c.fieldName)
		if err != nil {
			return 0, err
		}
		fj, err := mj.GetField(c.fieldName)
		if err != nil {
			return 0, err
		}
		iv, iValid, iNull := fi.(*IntField).Get()
		jv, jValid, jNull := fj.(*IntField).Get()
		if iValid && jValid {
			if !(iNull || jNull) {
				switch {
				case iv < jv:
					return -1, nil
				case iv > jv:
					return 1, nil
				default:
					return 0, nil
				}
			} else {
				switch {
				case iNull && jNull:
					return 0, nil
				case iNull:
					return -1, nil
				default:
					return 1, nil
				}
			}
		} else {
			switch {
			case !(iValid || jValid):
				return 0, nil
			case iValid:
				return 1, nil
			default:
				return -1, nil
			}
		}
	}
}

func (c *IntFieldClass) DistFunc() func(ModelInstance) (string, error) {
	return func(instance ModelInstance) (string, error) {
		field, err := instance.GetField(c.fieldName)
		if err != nil {
			return "", err
		}
		return field.SQLVal(), nil
	}
}

type FloatFieldClass struct {
	tableName  string
	columnName string
	modelName  string
	fieldName  string
}

func NewFloatFieldClass(tableName, columnName, modelName, fieldName string) *FloatFieldClass {
	return &FloatFieldClass{tableName, columnName, modelName, fieldName}
}

func (c *FloatFieldClass) NewInstance() FloatField {
	return FloatField{class: c}
}

func (c *FloatFieldClass) New(val float64, isNull bool) [2]string {
	l := [2]string{}
	l[0] = fmt.Sprintf("%s.%s", c.tableName, c.columnName)
	if isNull {
		l[1] = "NULL"
	} else {
		l[1] = fmt.Sprintf("%f", val)
	}
	return l
}

func (c *FloatFieldClass) Set(val float64, isNull bool) SetStruct {
	return SetStruct{
		c.modelName,
		c.fieldName,
		func(f Field) {
			ff := f.(*FloatField)
			ff.Set(val, isNull)
		},
	}
}

func (c *FloatFieldClass) Eq(val float64) *Where {
	return &Where{fmt.Sprintf("%s.%s = %f", c.tableName, c.columnName, val)}
}

func (c *FloatFieldClass) Neq(val float64) *Where {
	return &Where{fmt.Sprintf("%s.%s <> %f", c.tableName, c.columnName, val)}
}

func (c *FloatFieldClass) Lt(val float64) *Where {
	return &Where{fmt.Sprintf("%s.%s < %f", c.tableName, c.columnName, val)}
}

func (c *FloatFieldClass) Gt(val float64) *Where {
	return &Where{fmt.Sprintf("%s.%s > %f", c.tableName, c.columnName, val)}
}

func (c *FloatFieldClass) Lte(val float64) *Where {
	return &Where{fmt.Sprintf("%s.%s <= %f", c.tableName, c.columnName, val)}
}

func (c *FloatFieldClass) Gte(val float64) *Where {
	return &Where{fmt.Sprintf("%s.%s >= %f", c.tableName, c.columnName, val)}
}

func (c *FloatFieldClass) IsNull() *Where {
	return &Where{fmt.Sprintf("%s.%s IS NULL", c.tableName, c.columnName)}
}

func (c *FloatFieldClass) LessFunc() func(ModelInstance, ModelInstance) (int, error) {
	return func(mi, mj ModelInstance) (int, error) {
		fi, err := mi.GetField(c.fieldName)
		if err != nil {
			return 0, err
		}
		fj, err := mj.GetField(c.fieldName)
		if err != nil {
			return 0, err
		}
		iv, iValid, iNull := fi.(*IntField).Get()
		jv, jValid, jNull := fj.(*IntField).Get()
		if iValid && jValid {
			if !(iNull || jNull) {
				switch {
				case iv < jv:
					return -1, nil
				case iv > jv:
					return 1, nil
				default:
					return 0, nil
				}
			} else {
				switch {
				case iNull && jNull:
					return 0, nil
				case iNull:
					return -1, nil
				default:
					return 1, nil
				}
			}
		} else {
			switch {
			case !(iValid || jValid):
				return 0, nil
			case iValid:
				return 1, nil
			default:
				return -1, nil
			}
		}
	}
}

func (c *FloatFieldClass) DistFunc() func(ModelInstance) (string, error) {
	return func(instance ModelInstance) (string, error) {
		field, err := instance.GetField(c.fieldName)
		if err != nil {
			return "", err
		}
		return field.SQLVal(), nil
	}
}

type BoolFieldClass struct {
	tableName  string
	columnName string
	modelName  string
	fieldName  string
}

func NewBoolFieldClass(tableName, columnName, modelName, fieldName string) *BoolFieldClass {
	return &BoolFieldClass{tableName, columnName, modelName, fieldName}
}

func (c *BoolFieldClass) NewInstance() BoolField {
	return BoolField{class: c}
}

func (c *BoolFieldClass) New(val, isNull bool) [2]string {
	l := [2]string{}
	l[0] = fmt.Sprintf("%s.%s", c.tableName, c.columnName)
	if isNull {
		l[1] = "NULL"
	} else {
		if val {
			l[1] = "1"
		} else {
			l[1] = "0"
		}
	}
	return l
}

func (c *BoolFieldClass) Set(val, isNull bool) SetStruct {
	return SetStruct{
		c.modelName,
		c.fieldName,
		func(f Field) {
			bf := f.(*BoolField)
			bf.Set(val, isNull)
		},
	}
}

func (c *BoolFieldClass) Eq(val bool) *Where {
	if val {
		return &Where{fmt.Sprintf("%s.%s = 1", c.tableName, c.columnName)}
	} else {
		return &Where{fmt.Sprintf("%s.%s = 0", c.tableName, c.columnName)}
	}
}

func (c *BoolFieldClass) Neq(val bool) *Where {
	if val {
		return &Where{fmt.Sprintf("%s.%s <> 1", c.tableName, c.columnName)}
	} else {
		return &Where{fmt.Sprintf("%s.%s <> 0", c.tableName, c.columnName)}
	}
}

func (c *BoolFieldClass) IsNull(val bool) *Where {
	return &Where{fmt.Sprintf("%s.%s IS NULL", c.tableName, c.columnName)}
}

func (c *BoolFieldClass) LessFunc() func(ModelInstance, ModelInstance) (int, error) {
	return func(mi, mj ModelInstance) (int, error) {
		fi, err := mi.GetField(c.fieldName)
		if err != nil {
			return 0, err
		}
		fj, err := mj.GetField(c.fieldName)
		if err != nil {
			return 0, err
		}
		iv, iValid, iNull := fi.(*BoolField).Get()
		jv, jValid, jNull := fj.(*BoolField).Get()
		if iValid && jValid {
			if !(iNull || jNull) {
				switch {
				case iv && jv:
					return 0, nil
				case iv:
					return 1, nil
				default:
					return -1, nil
				}
			} else {
				switch {
				case iNull && jNull:
					return 0, nil
				case iNull:
					return -1, nil
				default:
					return 1, nil
				}
			}
		} else {
			switch {
			case !(iValid || jValid):
				return 0, nil
			case iValid:
				return 1, nil
			default:
				return -1, nil
			}
		}
	}
}

func (c *BoolFieldClass) DistFunc() func(ModelInstance) (string, error) {
	return func(instance ModelInstance) (string, error) {
		field, err := instance.GetField(c.fieldName)
		if err != nil {
			return "", err
		}
		return field.SQLVal(), nil
	}
}

type DateFieldClass struct {
	tableName  string
	columnName string
	modelName  string
	fieldName  string
}

func NewDateFieldClass(tableName, columnName, modelName, fieldName string) *DateFieldClass {
	return &DateFieldClass{tableName, columnName, modelName, fieldName}
}

func (c *DateFieldClass) NewInstance() DateField {
	return DateField{class: c}
}

func (c *DateFieldClass) New(val time.Time, isNull bool) [2]string {
	l := [2]string{}
	l[0] = fmt.Sprintf("%s.%s", c.tableName, c.columnName)
	if isNull {
		l[1] = "NULL"
	} else {
		l[1] = fmt.Sprintf("%q", val.Format("2006-01-02"))
	}
	return l
}

func (c *DateFieldClass) Set(val time.Time, isNull bool) SetStruct {
	return SetStruct{
		c.modelName,
		c.fieldName,
		func(f Field) {
			df := f.(*DateField)
			df.Set(val, isNull)
		},
	}
}

func (c *DateFieldClass) Eq(val time.Time) *Where {
	return &Where{fmt.Sprintf("%s.%s = %q", c.tableName, c.columnName, val.Format("2006-01-02"))}
}

func (c *DateFieldClass) Neq(val time.Time) *Where {
	return &Where{fmt.Sprintf("%s.%s <> %q", c.tableName, c.columnName, val.Format("2006-01-02"))}
}

func (c *DateFieldClass) Lt(val time.Time) *Where {
	return &Where{fmt.Sprintf("%s.%s < %q", c.tableName, c.columnName, val.Format("2006-01-02"))}
}

func (c *DateFieldClass) Gt(val time.Time) *Where {
	return &Where{fmt.Sprintf("%s.%s > %q", c.tableName, c.columnName, val.Format("2006-01-02"))}
}

func (c *DateFieldClass) Lte(val time.Time) *Where {
	return &Where{fmt.Sprintf("%s.%s <= %q", c.tableName, c.columnName, val.Format("2006-01-02"))}
}

func (c *DateFieldClass) Gte(val time.Time) *Where {
	return &Where{fmt.Sprintf("%s.%s >= %q", c.tableName, c.columnName, val.Format("2006-01-02"))}
}

func (c *DateFieldClass) IsNull() Where {
	return Where{fmt.Sprintf("%s.%s IS NULL", c.tableName, c.columnName)}
}

func (c *DateFieldClass) LessFunc() func(ModelInstance, ModelInstance) (int, error) {
	return func(mi, mj ModelInstance) (int, error) {
		fi, err := mi.GetField(c.fieldName)
		if err != nil {
			return 0, err
		}
		fj, err := mj.GetField(c.fieldName)
		if err != nil {
			return 0, err
		}
		iv, iValid, iNull := fi.(*DateField).Get()
		jv, jValid, jNull := fj.(*DateField).Get()
		if iValid && jValid {
			if !(iNull || jNull) {
				switch {
				case iv.Before(jv):
					return -1, nil
				case iv.After(jv):
					return 1, nil
				default:
					return 0, nil
				}
			} else {
				switch {
				case iNull && jNull:
					return 0, nil
				case iNull:
					return -1, nil
				default:
					return 1, nil
				}
			}
		} else {
			switch {
			case !(iValid || jValid):
				return 0, nil
			case iValid:
				return 1, nil
			default:
				return -1, nil
			}
		}
	}
}

func (c *DateFieldClass) DistFunc() func(ModelInstance) (string, error) {
	return func(instance ModelInstance) (string, error) {
		field, err := instance.GetField(c.fieldName)
		if err != nil {
			return "", err
		}
		return field.SQLVal(), nil
	}
}

type DatetimeFieldClass struct {
	tableName  string
	columnName string
	modelName  string
	fieldName  string
}

func NewDatetimeFieldClass(tableName, columnName, modelName, fieldName string) *DatetimeFieldClass {
	return &DatetimeFieldClass{tableName, columnName, modelName, fieldName}
}

func (c *DatetimeFieldClass) NewInstance() DatetimeField {
	return DatetimeField{class: c}
}

func (c *DatetimeFieldClass) New(val time.Time, isNull bool) [2]string {
	l := [2]string{}
	l[0] = fmt.Sprintf("%s.%s", c.tableName, c.columnName)
	if isNull {
		l[1] = "NULL"
	} else {
		l[1] = fmt.Sprintf("%q", val.Format("2006-01-02 15:04:05"))
	}
	return l
}

func (c *DatetimeFieldClass) Set(val time.Time, isNull bool) SetStruct {
	return SetStruct{
		c.modelName,
		c.fieldName,
		func(f Field) {
			dtf := f.(*DatetimeField)
			dtf.Set(val, isNull)
		},
	}
}

func (c *DatetimeFieldClass) Eq(val time.Time) *Where {
	return &Where{fmt.Sprintf("%s.%s = %q", c.tableName, c.columnName, val.Format("2006-01-02 15:04:05"))}
}

func (c *DatetimeFieldClass) Neq(val time.Time) *Where {
	return &Where{fmt.Sprintf("%s.%s <> %q", c.tableName, c.columnName, val.Format("2006-01-02 15:04:05"))}
}

func (c *DatetimeFieldClass) Lt(val time.Time) *Where {
	return &Where{fmt.Sprintf("%s.%s < %q", c.tableName, c.columnName, val.Format("2006-01-02 15:04:05"))}
}

func (c *DatetimeFieldClass) Gt(val time.Time) *Where {
	return &Where{fmt.Sprintf("%s.%s > %q", c.tableName, c.columnName, val.Format("2006-01-02 15:04:05"))}
}

func (c *DatetimeFieldClass) Lte(val time.Time) *Where {
	return &Where{fmt.Sprintf("%s.%s <= %q", c.tableName, c.columnName, val.Format("2006-01-02 15:04:05"))}
}

func (c *DatetimeFieldClass) Gte(val time.Time) *Where {
	return &Where{fmt.Sprintf("%s.%s > %q", c.tableName, c.columnName, val.Format("2006-01-02 15:04:05"))}
}

func (c *DatetimeFieldClass) IsNull() *Where {
	return &Where{fmt.Sprintf("%s.%s IS NULL", c.tableName, c.columnName)}
}

func (c *DatetimeFieldClass) LessFunc() func(ModelInstance, ModelInstance) (int, error) {
	return func(mi, mj ModelInstance) (int, error) {
		fi, err := mi.GetField(c.fieldName)
		if err != nil {
			return 0, err
		}
		fj, err := mj.GetField(c.fieldName)
		if err != nil {
			return 0, err
		}
		iv, iValid, iNull := fi.(*DatetimeField).Get()
		jv, jValid, jNull := fj.(*DatetimeField).Get()
		if iValid && jValid {
			if !(iNull || jNull) {
				switch {
				case iv.Before(jv):
					return -1, nil
				case iv.After(jv):
					return 1, nil
				default:
					return 0, nil
				}
			} else {
				switch {
				case iNull && jNull:
					return 0, nil
				case iNull:
					return -1, nil
				default:
					return 1, nil
				}
			}
		} else {
			switch {
			case !(iValid || jValid):
				return 0, nil
			case iValid:
				return 1, nil
			default:
				return -1, nil
			}
		}
	}
}

func (c *DatetimeFieldClass) DistFunc() func(ModelInstance) (string, error) {
	return func(instance ModelInstance) (string, error) {
		field, err := instance.GetField(c.fieldName)
		if err != nil {
			return "", err
		}
		return field.SQLVal(), nil
	}
}

type ModelInstance interface {
	Invalidate()
	SetLastInsertID(int64)
	GetField(string) (Field, error)
}

type Where struct {
	str string
}

func (w *Where) String() string {
	return w.str
}

func (w *Where) And(other *Where) *Where {
	if w == nil && other != nil {
		return other
	} else if w != nil && other == nil {
		return w
	} else if w == nil && other == nil {
		return nil
	}
	return &Where{w.str + " AND " + other.str}
}

func (w *Where) Or(other *Where) *Where {
	if w == nil && other != nil {
		return other
	} else if w != nil && other == nil {
		return w
	} else if w == nil && other == nil {
		return nil
	}
	return &Where{w.str + " OR " + other.str}
}

func Group(w *Where) *Where {
	if w == nil {
		return nil
	}
	return &Where{"(" + w.str + ")"}
}

type stmtType int

const (
	insert stmtType = iota
	update
	delete
	other
)

type Stmt struct {
	model        ModelInstance
	typ          stmtType
	stmt         string
	lastInsertID int64
}

func NewStmt(instance ModelInstance, stmtStr string) *Stmt {
	stmt := &Stmt{model: instance, stmt: stmtStr}
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
	l := make(StmtList, len(stmts)+1)
	l[0] = s
	for i := 0; i < len(stmts); i++ {
		l[i+1] = stmts[i]
	}
	return l
}

func (s *Stmt) Exec(db sql.DB) error {
	res, err := db.Exec(s.stmt)
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

func (sl StmtList) Exec(db sql.DB) error {
	tx, err := db.Begin()
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
}
