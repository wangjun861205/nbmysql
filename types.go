package nbmysql

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

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
