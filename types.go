package nbmysql

import (
	"database/sql"
	"fmt"
	"strings"
)

type UniqueKey []string

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
}

type Table struct {
	Columns         []Column
	ModelName       string
	TableName       string
	ArgName         string
	PrimaryKey      *Column
	AutoIncrement   *Column
	ForeignKeyInfos []ForeignKeyInfo
	ManyToManyInfos []ManyToManyInfo
	ForeignKeys     []ForeignKey
	ManyToManys     []ManyToMany
	UniqueKeys      []UniqueKey
}

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

func (db *Database) CreateTableIfNotExists(tab Table) error {
	conn, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", db.Username, db.Password, db.Address, db.DatabaseName))
	if err != nil {
		return err
	}
	defer conn.Close()
	colList := make([]string, len(tab.Columns))
	for i, col := range tab.Columns {
		l := make([]string, 0, 6)
		l = append(l, BackQuote(col.ColumnName))
		l = append(l, col.MySqlType)
		if col.AutoIncrement {
			l = append(l, "AUTO_INCREMENT")
		}
		if !col.Nullable {
			l = append(l, "NOT NULL")
		}
		if col.Default != "" {
			l = append(l, "DEFAULT "+col.Default)
		}
		if col.Unique {
			l = append(l, "UNIQUE")
		}
		colList[i] = strings.Join(l, " ")
	}
	uniqueList := make([]string, 0, len(tab.UniqueKeys))
	for _, uni := range tab.UniqueKeys {
		bqList := make([]string, len(uni))
		for i, _ := range uni {
			bqList[i] = BackQuote(string(uni[i]))
		}
		uniqueList = append(uniqueList, fmt.Sprintf("UNIQUE KEY `%s_unique` (%s)", strings.Join(uni, "_"), strings.Join(bqList, ", ")))
	}
	stmt := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s, PRIMARY KEY (%s), %s)", BackQuote(tab.TableName), strings.Join(colList, ", "),
		BackQuote(tab.PrimaryKey.ColumnName), strings.Join(uniqueList, ", "))
	_, err = conn.Exec(stmt)
	return err
}
