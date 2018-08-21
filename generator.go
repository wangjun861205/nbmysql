package nbmysql

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/wangjun861205/nbfmt"
)

//genPackage generate package statement
func genPackage(pkg string) string {
	s, err := nbfmt.Fmt(packageTemp, map[string]interface{}{"Package": pkg})
	if err != nil {
		fmt.Println("error: genPackage() failed")
		log.Fatal(err)
	}
	return s
}

//genDb generate database object declare statement
func genDb(db Database) string {
	s, err := nbfmt.Fmt(dbTemp, map[string]interface{}{"Database": db})
	if err != nil {
		fmt.Println("error: genDb() failed")
		log.Fatal(err)
	}
	return s
}

//genInitFunc generate init function
func genInitFunc(db Database, stmtInit string) string {
	s, err := nbfmt.Fmt(initFuncTemp, map[string]interface{}{"Database": db, "Block": stmtInit})
	if err != nil {
		fmt.Println("error: genInitFunc() failed")
		log.Fatal(err)
	}
	return s
}

//genMapElemBlock generate map element
func genMapElemBlock(col Column) string {
	s, err := nbfmt.Fmt(mapElemTemp, map[string]interface{}{"Column": col})
	if err != nil {
		fmt.Println("error: genMapElemBlock() failed")
		log.Fatal(err)
	}
	return s

}

//genQueryFieldMapBlock generate a map for mapping filed name to column name in where clause in query
func genQueryFieldMapBlock(tab Table) string {
	elemList := make([]string, len(tab.Columns))
	for i, col := range tab.Columns {
		elemList[i] = genMapElemBlock(col)
	}
	s, err := nbfmt.Fmt(queryFieldMapTemp, map[string]interface{}{"Table": tab, "Block": strings.Join(elemList, "\n")})
	if err != nil {
		fmt.Println("error: genQueryFialedMapBlock() failed")
		log.Fatal(err)
	}
	return s
}

//genField generate field definition in model struct
func genField(col Column) string {
	s, err := nbfmt.Fmt(fieldTemp, map[string]interface{}{"Column": col})
	if err != nil {
		fmt.Println("error: genField() failed")
		log.Fatal(err)
	}
	return s
}

//genModel generate model defination
func genModel(tab Table) string {
	list := make([]string, len(tab.Columns))
	for i, col := range tab.Columns {
		list[i] = genField(col)
	}
	s, err := nbfmt.Fmt(modelTemp, map[string]interface{}{"Table": tab, "Block": strings.Join(list, "\n")})
	if err != nil {
		fmt.Println("error: genModel() failed")
		log.Fatal(err)
	}
	return s
}

//genModelGetFieldMethod generate model get xxx field value method
func genModelGetFieldMethod(col Column, tab Table) string {
	var zeroFunc = func() string {
		var s string
		switch col.FieldType {
		case "int64":
			s = "0"
		case "float64":
			s = "0.0"
		case "string":
			s = `""`
		case "bool":
			s = "false"
		case "time.Time":
			s = "time.Time{}"
		}
		return s
	}
	s, err := nbfmt.Fmt(modelGetFieldMethodTemp, map[string]interface{}{
		"Table":     tab,
		"Column":    col,
		"ZeroValue": zeroFunc()})
	if err != nil {
		fmt.Println("error: genModelGetFieldMethod() failded")
		log.Fatal(err)
	}
	return s
}

//genModelSetFieldMethod generate model set xxx field value method
func genModelSetFieldMethod(col Column, tab Table) string {
	s, err := nbfmt.Fmt(modelSetFieldMethodTemp, map[string]interface{}{"Table": tab, "Column": col})
	if err != nil {
		fmt.Println("error: genModelSetFieldMethod() failed")
		log.Fatal(err)
	}
	return s
}

//genNewFuncArgList generate argument list in NewXXX() function
func genNewFuncArgList(tab Table) string {
	list := make([]string, 0, len(tab.Columns))
	for _, col := range tab.Columns {
		if !col.AutoIncrement && col.Default != "CURRENT_TIMESTAMP" && col.On != "UPDATE CURRENT_TIMESTAMP" {
			s, err := nbfmt.Fmt(funcArgTemp, map[string]interface{}{"Column": col})
			if err != nil {
				fmt.Println("error: funcArgTemp format failed")
				log.Fatal(err)
			}
			list = append(list, s)
		}
	}
	return strings.Join(list, ", ")
}

//genNewFuncAsignBlock generate middle variable asign block in NewXXX() function
func genNewFuncAsignBlock(tab Table) string {
	list := make([]string, 0, len(tab.Columns))
	for _, col := range tab.Columns {
		if !col.AutoIncrement && col.Default != "CURRENT_TIMESTAMP" && col.On != "UPDATE CURRENT_TIMESTAMP" {
			s, err := nbfmt.Fmt(newModelAsignTemp, map[string]interface{}{"Column": col})
			if err != nil {
				fmt.Println("error: newModelAsignTemp format failed")
				log.Fatal(err)
			}
			list = append(list, s)
		}
	}
	return strings.Join(list, ",\n")
}

//genNewFunc generate NewXXX() function
func genNewFunc(tab Table) string {
	s, err := nbfmt.Fmt(newModelFuncTemp, map[string]interface{}{"Table": tab, "Args": genNewFuncArgList(tab), "Asigns": genNewFuncAsignBlock(tab)})
	if err != nil {
		fmt.Println("error: genNewFunc() failed")
		log.Fatal(err)
	}
	return s
}

//genAllFunc generate AllXXX() function
func genAllFunc(tab Table, db Database) string {
	s, err := nbfmt.Fmt(allModelFuncTemp, map[string]interface{}{"Table": tab, "Database": db})
	if err != nil {
		fmt.Println("error: genAllFunc() failed")
		log.Fatal(err)
	}
	return s
}

//genQueryFunc generate QueryXXX() function
func genQueryFunc(tab Table, db Database) string {
	s, err := nbfmt.Fmt(queryModelFuncTemp, map[string]interface{}{"Table": tab, "Database": db})
	if err != nil {
		fmt.Println("error: genQueryFunc() failed")
		log.Fatal(err)
	}
	return s
}

//genQueryOneFunc generate QueryOneXXX() function
func genQueryOneFunc(tab Table, db Database) string {
	s, err := nbfmt.Fmt(queryOneFuncTemp, map[string]interface{}{"Table": tab, "Database": db})
	if err != nil {
		fmt.Println("error: genQueryOneFunc() failed")
		log.Fatal(err)
	}
	return s
}

//genForeignKeyMethod generate foreign key relation method
func genForeignKeyMethod(fk ForeignKey, srcTab Table, db Database) string {
	queryStmt, err := nbfmt.Fmt(foreignKeyQuerySQLTemp, map[string]interface{}{"FK": fk})
	if err != nil {
		log.Fatal(err)
	}
	s, err := nbfmt.Fmt(foreignKeyMethodTemp, map[string]interface{}{
		"Table":     srcTab,
		"FK":        fk,
		"QueryStmt": queryStmt,
		"Database":  db})
	if err != nil {
		fmt.Println("error: genForeignKeyMethod() failed")
		log.Fatal(err)
	}
	return s
}

//genReverseForeignKeyStruct generate reverse foreign key relation method
func genReverseForeignKeyStruct(rfk ReverseForeignKey, srcTab Table, db Database) string {
	s, err := nbfmt.Fmt(reverseForeignKeyStructTypeTemp, map[string]interface{}{"Table": srcTab, "RFK": rfk})
	if err != nil {
		fmt.Println("error: genReverseForeignKeyStruct() failed")
		log.Fatal(err)
	}
	return s
}

//genReverseForeignKeyAllMethod generate All() method in reverse foreign key relation struct
func genReverseForeignKeyAllMethod(rfk ReverseForeignKey, srcTab Table, db Database) string {
	queryStmt, err := nbfmt.Fmt(reverseForeignKeyAllSQLTemp, map[string]interface{}{"RFK": rfk})
	if err != nil {
		fmt.Println("error: reverseForgienKeyAllSQLTemp format failed")
		log.Fatal(err)
	}
	s, err := nbfmt.Fmt(reverseForeignKeyAllMethodTemp, map[string]interface{}{"RFK": rfk, "Table": srcTab, "Database": db, "QueryStmt": queryStmt})
	if err != nil {
		fmt.Println("error: genReverseForeignKeyAllMethod() failed")
		log.Fatal(err)
	}
	return s
}

//genReverseForeignKeyQueryMethod generate Query() method in reverse foreign key relation struct
func genReverseForeignKeyQueryMethod(rfk ReverseForeignKey, srcTab Table, db Database) string {
	queryStmt, err := nbfmt.Fmt(reverseForeignKeyQuerySQLTemp, map[string]interface{}{"RFK": rfk})
	if err != nil {
		fmt.Println("error: reverseForeignKeyQuerySQLTemp format failed")
		log.Fatal(err)
	}
	s, err := nbfmt.Fmt(reverseForeignKeyQueryMethodTemp, map[string]interface{}{"RFK": rfk, "Table": srcTab, "Database": db, "QueryStmt": queryStmt})
	if err != nil {
		fmt.Println("error: genReverseForeignKeyQueryMethod() failed")
		log.Fatal(err)
	}
	return s
}

//genReverseForeignKeyMethod generate model reverse foreign key relation method
func genReverseForeignKeyMethod(rfk ReverseForeignKey, srcTab Table, db Database) string {
	allMethod := genReverseForeignKeyAllMethod(rfk, srcTab, db)
	queryMethod := genReverseForeignKeyQueryMethod(rfk, srcTab, db)
	s, err := nbfmt.Fmt(reverseForeignKeyMethodTemp, map[string]interface{}{"Table": srcTab, "RFK": rfk, "AllMethod": allMethod, "QueryMethod": queryMethod, "Database": db})
	if err != nil {
		fmt.Println("error: genReverseForeignKeyMethod() failed")
		log.Fatal(err)
	}
	return s
}

//genManyToManyStruct generate many to many relation struct definition
func genManyToManyStruct(mtm ManyToMany, srcTab Table, db Database) string {
	s, err := nbfmt.Fmt(manyToManyStructTypeTemp, map[string]interface{}{"Table": srcTab, "MTM": mtm})
	if err != nil {
		fmt.Println("error: genManyToManyStruct() failed")
		log.Fatal(err)
	}
	return s
}

//genManyToManyAllMethod generate All() method in many to many relation struct
func genManyToManyAllMethod(mtm ManyToMany, srcTab Table, db Database) string {
	queryStmt, err := nbfmt.Fmt(manyToManyAllSQLTemp, map[string]interface{}{"MTM": mtm, "Table": srcTab})
	if err != nil {
		fmt.Println("error: manyToManyAllSQLTemp format failed")
		log.Fatal(err)
	}
	s, err := nbfmt.Fmt(manyToManyAllMethodTemp, map[string]interface{}{"MTM": mtm, "Database": db, "QueryStmt": queryStmt})
	if err != nil {
		fmt.Println("error: genManyToManyAllMethod() failed")
		log.Fatal(err)
	}
	return s
}

//genManyToManyQueryMethod generate Query() method in many to many relation struct
func genManyToManyQueryMethod(mtm ManyToMany, srcTab Table, db Database) string {
	queryStmt, err := nbfmt.Fmt(manyToManyQuerySQLTemp, map[string]interface{}{"Table": srcTab, "MTM": mtm})
	if err != nil {
		fmt.Println("error: manyToManyQuerySQLTemp format failed")
		log.Fatal(err)
	}
	s, err := nbfmt.Fmt(manyToManyQueryMethodTemp, map[string]interface{}{"MTM": mtm, "Database": db, "QueryStmt": queryStmt})
	if err != nil {
		fmt.Println("error: genManyToManyQueryMethod() failed")
		log.Fatal(err)
	}
	return s
}

//genManyToManyAddMethod generate Add() method in many to many relation struct
func genManyToManyAddMethod(mtm ManyToMany, tab Table) string {
	s, err := nbfmt.Fmt(manyToManyAddMethodTemp, map[string]interface{}{"Table": tab, "MTM": mtm})
	if err != nil {
		fmt.Println("error: genManyToManyAddMethod() failed")
		log.Fatal(err)
	}
	return s
}

//genManyToManyRemoveMethod generate Remove() method in many to many relation struct
func genManyToManyRemoveMethod(mtm ManyToMany, tab Table) string {
	s, err := nbfmt.Fmt(manyToManyRemoveMethodTemp, map[string]interface{}{"MTM": mtm, "Table": tab})
	if err != nil {
		fmt.Println("error: genManyToManyRemoveMethod() failed")
		log.Fatal(err)
	}
	return s
}

//genManyToManyMethod generate model many to many relation method
func genManyToManyMethod(mtm ManyToMany, srcTab Table, db Database) string {
	allMethod := genManyToManyAllMethod(mtm, srcTab, db)
	queryMethod := genManyToManyQueryMethod(mtm, srcTab, db)
	addMethod := genManyToManyAddMethod(mtm, srcTab)
	removeMethod := genManyToManyRemoveMethod(mtm, srcTab)
	s, err := nbfmt.Fmt(manyToManyMethodTemp, map[string]interface{}{
		"Database":     db,
		"Table":        srcTab,
		"MTM":          mtm,
		"AllMethod":    allMethod,
		"QueryMethod":  queryMethod,
		"AddMethod":    addMethod,
		"RemoveMethod": removeMethod})
	if err != nil {
		fmt.Println("error: genManyToManyMethod() failed")
		log.Fatal(err)
	}
	return s
}

//genStmtArgNilToDefault generate block for replacing nil argument to default value in insert or update method when not null column value is nil
func genStmtArgNilToDefault(col Column) string {
	var value interface{}
	var err error
	switch col.FieldType {
	case "int64":
		value, err = strconv.ParseInt(col.Default, 10, 64)
	case "float64":
		value, err = strconv.ParseFloat(col.Default, 64)
	case "bool":
		value, err = strconv.ParseInt(col.Default, 10, 64)
	case "string":
		value = col.Default
	case "time.Time":
		if col.MySqlType == "DATE" {
			value = "time.Now()"
		} else {
			value = "time.Now()"
		}
	}
	if err != nil {
		log.Fatal(err)
	}
	s, err := nbfmt.Fmt(stmtArgNilToDefaultTemp, map[string]interface{}{"Column": col, "DefaultValue": value})
	if err != nil {
		fmt.Println("error: genStmtArgNilToDefault() failed")
		log.Fatal(err)
	}
	return s
}

//genStmtArgNilToDefaultBlock generate block to replacing nil arguments to default value
func genStmtArgNilToDefaultBlock(tab Table) string {
	list := make([]string, 0, len(tab.Columns))
	for _, col := range tab.Columns {
		if !col.AutoIncrement && col.Default != "CURRENT_TIMESTAMP" && col.On != "UPDATE CURRENT_TIMESTAMP" {
			if !col.Nullable && col.Default != "" {
				list = append(list, genStmtArgNilToDefault(col))
			} else {
				s, err := nbfmt.Fmt(stmtArgTemp, map[string]interface{}{"Column": col})
				if err != nil {
					fmt.Println("error: stmtArgTemp format failed")
					log.Fatal(err)
				}
				list = append(list, s)
			}
		}
	}
	s, err := nbfmt.Fmt(stmtArgNilToDefaultBlockTemp, map[string]interface{}{"Length": len(list), "Block": strings.Join(list, "\n")})
	if err != nil {
		fmt.Println("error: genStmtArgNilToDefaultBlock() failed")
		log.Fatal(err)
	}
	return s
}

//genModelInsertMethod generate XXX.Insert() method
func genModelInsertMethod(tab Table) string {
	s, err := nbfmt.Fmt(modelInsertMethodTemp, map[string]interface{}{"Table": tab, "Block": genStmtArgNilToDefaultBlock(tab)})
	if err != nil {
		fmt.Println("error: genModelInsertMethod() failed")
		log.Fatal(err)
	}
	return s
}

//genModelUpdateMethod generate XXX.Update() method
func genModelUpdateMethod(tab Table) string {
	s, err := nbfmt.Fmt(modelUpdateMethodTemp, map[string]interface{}{"Table": tab, "Block": genStmtArgNilToDefaultBlock(tab)})
	if err != nil {
		fmt.Println("error: genModelUpdateMethod() failed")
		log.Fatal(err)
	}
	return s
}

//genModelInsertOrUpdateMethod generate XXX.InsertOrUpdate() method
func genModelInsertOrUpdateMethod(tab Table) string {
	s, err := nbfmt.Fmt(modelInsertOrUpdateMethodTemp, map[string]interface{}{"Table": tab, "Block": genStmtArgNilToDefaultBlock(tab)})
	if err != nil {
		fmt.Println("error: genModelInsertOrUpdateMethod() failed")
		log.Fatal(err)
	}
	return s
}

//genModelDeleteMethod generate XXX.Delete() method
func genModelDeleteMethod(tab Table) string {
	s, err := nbfmt.Fmt(modelDeleteMethodTemp, map[string]interface{}{"Table": tab})
	if err != nil {
		fmt.Println("error: genModelDeleteMethod() failed")
		log.Fatal(err)
	}
	return s
}

//genNewMiddleTypeBlock generate new middle type consponsed to field type block
func genNewMiddleTypeBlock(col Column) string {
	s, err := nbfmt.Fmt(newMiddleTypeTemp, map[string]interface{}{"Column": col})
	if err != nil {
		fmt.Println("error: genNewMiddleTypeBlock() failed")
		log.Fatal(err)
	}
	return s
}

//genFromRowsFunc generate XXXFromRows(sql.Rows) ([]*XXX, error) function
func genFromRowsFunc(tab Table) string {
	midList := make([]string, len(tab.Columns))
	midNameList := make([]string, len(tab.Columns))
	finalArgList := make([]string, len(tab.Columns))
	for i, col := range tab.Columns {
		midList[i] = genNewMiddleTypeBlock(col)
		midNameList[i] = "_" + col.ArgName
		finalArgList[i] = "_" + col.ArgName + ".ToGo()"
	}
	s, err := nbfmt.Fmt(modelFromRowsFuncTemp, map[string]interface{}{"Table": tab,
		"Block":    strings.Join(midList, "\n"),
		"MidArgs":  strings.Join(midNameList, ", "),
		"MidToGos": strings.Join(finalArgList, ", ")})
	if err != nil {
		fmt.Println("error: genFromRowsFunc() failed")
		log.Fatal(err)
	}
	return s
}

//genFromRowFunc generate XXXFromRow(sql.Rows) (*XXX, error) function
func genFromRowFunc(tab Table) string {
	midList := make([]string, len(tab.Columns))
	midNameList := make([]string, len(tab.Columns))
	finalArgList := make([]string, len(tab.Columns))
	for i, col := range tab.Columns {
		midList[i] = genNewMiddleTypeBlock(col)
		midNameList[i] = "_" + col.ArgName
		finalArgList[i] = "_" + col.ArgName + ".ToGo()"
	}
	s, err := nbfmt.Fmt(modelFromRowFuncTemp, map[string]interface{}{
		"Table":    tab,
		"Block":    strings.Join(midList, "\n"),
		"MidArgs":  strings.Join(midNameList, ", "),
		"MidToGos": strings.Join(finalArgList, ", "),
	})
	if err != nil {
		fmt.Println("error: genFromRowFunc() failed")
		log.Fatal(err)
	}
	return s
}

//genModelCheckMethod generate XXX.check() method.This method is for excluding AUTO_INCREMTN column and AUTO_TIMESTAMP column from insert and update method
func genModelCheckMethod(tab Table) string {
	list := make([]string, 0, len(tab.Columns))
	for _, col := range tab.Columns {
		if !col.Nullable && !col.AutoIncrement && col.Default == "" {
			s, err := nbfmt.Fmt(fieldCheckNullTemp, map[string]interface{}{"Column": col, "Table": tab})
			if err != nil {
				fmt.Println("error: fieldCheckNullTemp format failed")
				log.Fatal(err)
			}
			list = append(list, s)
		}
	}
	s, err := nbfmt.Fmt(modelCheckMethodTemp, map[string]interface{}{"Table": tab, "Block": strings.Join(list, "\n")})
	if err != nil {
		fmt.Println("error: genModelCheckMethod() failed")
		log.Fatal(err)
	}
	return s
}

//genInsertStmt generate insert sql.Stmt declaration
func genInsertStmt(tab Table) string {
	s, err := nbfmt.Fmt(insertStmtTemp, map[string]interface{}{"Table": tab})
	if err != nil {
		fmt.Println("error: genInsertStmt() failed")
		log.Fatal(err)
	}
	return s
}

//genUpdateStmt generate update sql.Stmt declaration
func genUpdateStmt(tab Table) string {
	s, err := nbfmt.Fmt(updateStmtTemp, map[string]interface{}{"Table": tab})
	if err != nil {
		fmt.Println("error: genUpdateStmt() failed")
		log.Fatal(err)
	}
	return s
}

//genDeleteStmt generate delete sql.Stmt declaration
func genDeleteStmt(tab Table) string {
	s, err := nbfmt.Fmt(deleteStmtTemp, map[string]interface{}{"Table": tab})
	if err != nil {
		fmt.Println("error: genDeleteStmt() failed")
		log.Fatal(err)
	}
	return s
}

//genInsertOrUpdateStmt generate insert or update sql.Stmt declaration
func genInsertOrUpdateStmt(tab Table) string {
	s, err := nbfmt.Fmt(insertOrUpdateStmtTemp, map[string]interface{}{"Table": tab})
	if err != nil {
		fmt.Println("error: genInsertOrUpdateStmt() failed")
		log.Fatal(err)
	}
	return s
}

//genInsertMidStmt generate insert into middle table sql.Stmt declaration for many to many relation
func genInsertMidStmt(tab Table, mtm ManyToMany) string {
	s, err := nbfmt.Fmt(insertMidStmtTemp, map[string]interface{}{"Table": tab, "MTM": mtm})
	if err != nil {
		fmt.Println("error: genInsertMidStmt() failed")
		log.Fatal(err)
	}
	return s
}

//genManyToManyDeleteStmt generate delete sql.Stmt declaration for many to many relation
func genManyToManyDeleteStmt(mtm ManyToMany, tab Table) string {
	s, err := nbfmt.Fmt(manyToManyDeleteStmtTemp, map[string]interface{}{"Table": tab, "MTM": mtm})
	if err != nil {
		fmt.Println("error: genManyToManyDeleteStmt() failed")
		log.Fatal(err)
	}
	return s
}

//genStmtVar generate prepared sql.Stmt declaration
func genStmtVar(tab Table) string {
	list := make([]string, 0, 8)
	list = append(list, genInsertStmt(tab))
	list = append(list, genUpdateStmt(tab))
	list = append(list, genDeleteStmt(tab))
	list = append(list, genInsertOrUpdateStmt(tab))
	list = append(list, genModelCountStmtDeclare(tab))
	for _, mtm := range tab.ManyToManys {
		list = append(list, genInsertMidStmt(tab, mtm))
		list = append(list, genManyToManyDeleteStmt(mtm, tab))
	}
	return strings.Join(list, "\n")
}

//genInsertStmtInitBlock generate insert sql.Stmt in init function
func genInsertStmtInitBlock(tab Table, db Database) string {
	list := make([]string, 0, len(tab.Columns))
	for _, col := range tab.Columns {
		if !col.AutoIncrement && col.Default != "CURRENT_TIMESTAMP" && col.On != "UPDATE CURRENT_TIMESTAMP" {
			list = append(list, col.ColumnName)
		}
	}
	s, err := nbfmt.Fmt(insertStmtInitTemp, map[string]interface{}{
		"Table":    tab,
		"Database": db,
		"Columns":  strings.Join(list, ", "),
		"Values":   strings.Trim(strings.Repeat("?, ", len(list)), ", "),
	})
	if err != nil {
		fmt.Println("error: genInsertStmtInitBlock() failed")
		log.Fatal(err)
	}
	return s

}

//genUpdateStmtInitBlock generate update sql.Stmt in init function
func genUpdateStmtInitBlock(tab Table, db Database) string {
	list := make([]string, 0, len(tab.Columns))
	for _, col := range tab.Columns {
		if !col.AutoIncrement && col.Default != "CURRENT_TIMESTAMP" && col.On != "UPDATE CURRENT_TIMESTAMP" {
			list = append(list, col.ColumnName+" = ?")
		}
	}
	s, err := nbfmt.Fmt(updateStmtInitTemp, map[string]interface{}{
		"Table":    tab,
		"Database": db,
		"Updates":  strings.Join(list, ", "),
	})
	if err != nil {
		fmt.Println("error: genUpdateStmtInitBlock() failed")
		log.Fatal(err)
	}
	return s
}

//genDeleteStmtInitBlock generate delete sql.Stmt in init function
func genDeleteStmtInitBlock(tab Table, db Database) string {
	s, err := nbfmt.Fmt(deleteStmtInitTemp, map[string]interface{}{"Table": tab, "Database": db})
	if err != nil {
		fmt.Println("error: genDeleteStmtInitBlock() failed")
		log.Fatal(err)
	}
	return s
}

//genInsertOrUpdateStmtInitBlock generate insert or update sql.Stmt in init function
func genInsertOrUpdateStmtInitBlock(tab Table, db Database) string {
	insertList := make([]string, 0, len(tab.Columns))
	argList := make([]string, 0, len(tab.Columns))
	updateList := make([]string, 0, len(tab.Columns))
	for _, col := range tab.Columns {
		if !col.AutoIncrement && col.Default != "CURRENT_TIMESTAMP" && col.On != "UPDATE CURRENT_TIMESTAMP" {
			insertList = append(insertList, col.ColumnName)
			argList = append(argList, "?")
			s, err := nbfmt.Fmt(updateColumnTemp, map[string]interface{}{"Column": col})
			if err != nil {
				fmt.Println("error: updateColumnTemp format failed")
				log.Fatal(err)
			}
			updateList = append(updateList, s)
		}
	}
	updateLastInsertId, err := nbfmt.Fmt(updateLastInsertIDTemp, map[string]interface{}{"Table": tab})
	if err != nil {
		fmt.Println("error: updateLastInsertIDTemp format failed")
		log.Fatal(err)
	}
	s, err := nbfmt.Fmt(insertOrUpdateStmtInitTemp, map[string]interface{}{
		"Table":              tab,
		"Database":           db,
		"Columns":            strings.Join(insertList, ", "),
		"Values":             strings.Join(argList, ", "),
		"UpdateLastInsertID": updateLastInsertId,
		"Updates":            strings.Join(updateList, ", "),
	})
	if err != nil {
		fmt.Println("error: genInsertOrUpdateStmtInitBlock() failed")
		log.Fatal(err)
	}
	return s
}

//genInsertMidStmtInitBlock generate insert into middle table sql.Stmt for many to many relation in init function
func genInsertMidStmtInitBlock(mtm ManyToMany, tab Table, db Database) string {
	s, err := nbfmt.Fmt(insertMidStmtInitTemp, map[string]interface{}{"Table": tab, "Database": db, "MTM": mtm})
	if err != nil {
		fmt.Println("error: genInsertMidStmtInitBlock() failed")
		log.Fatal(err)
	}
	return s

}

//genManyToManyDeleteStmtInitBlock generate delete middle table sql.Stmt for many to many relation in init function
func genManyToManyDeleteStmtInitBlock(mtm ManyToMany, tab Table, db Database) string {
	s, err := nbfmt.Fmt(manyToManyDeleteStmtInitTemp, map[string]interface{}{"Table": tab, "MTM": mtm, "Database": db})
	if err != nil {
		fmt.Println("error: genManyToManyDeleteStmtInitBlock() failed")
		log.Fatal(err)
	}
	return s
}

//genStmtInitBlock generate sql.Stmt init block in init function
func genStmtInitBlock(tab Table, db Database) string {
	list := make([]string, 0, 8)
	list = append(list, genInsertStmtInitBlock(tab, db))
	list = append(list, genUpdateStmtInitBlock(tab, db))
	list = append(list, genDeleteStmtInitBlock(tab, db))
	list = append(list, genInsertOrUpdateStmtInitBlock(tab, db))
	list = append(list, genModelCountStmtInit(tab, db))
	for _, mtm := range tab.ManyToManys {
		list = append(list, genInsertMidStmtInitBlock(mtm, tab, db))
		list = append(list, genManyToManyDeleteStmtInitBlock(mtm, tab, db))
	}
	return strings.Join(list, "\n")
}

//genModelListType generate XXXList struct definition
func genModelListType(tab Table) string {
	s, err := nbfmt.Fmt(modelListTypeTemp, map[string]interface{}{"Table": tab})
	if err != nil {
		fmt.Println("error: genModelListType() failed")
		log.Fatal(err)
	}
	return s
}

func genModelCompareMethod(col Column) string {
	var temp string
	switch col.FieldType {
	case "int64":
		temp = modelCompareByIntMethodTemp
	case "float64":
		temp = modelCompareByFloatMethodTemp
	case "string":
		temp = modelCompareByStringMethodTemp
	case "bool":
		temp = modelCompareByBoolMethodTemp
	case "time.Time":
		temp = modelCompareByTimeMethodTemp
	}
	s, err := nbfmt.Fmt(temp, map[string]interface{}{"Column": col})
	if err != nil {
		fmt.Println("error: genModelCompareMethod() failed")
		log.Fatal(err)
	}
	return s
}

func genModelSortMethodsStructType(tab Table) string {
	list := make([]string, len(tab.Columns))
	var err error
	for i, col := range tab.Columns {
		list[i], err = nbfmt.Fmt(modelSortMethodsStructFieldTypeTemp, map[string]interface{}{"Column": col})
		if err != nil {
			fmt.Println("error: modelSortMethodStructFieldTypeTemp format failed")
			log.Fatal(err)
		}
	}
	s, err := nbfmt.Fmt(modelSortMethodsStructTypeTemp, map[string]interface{}{"Table": tab, "FieldTypeBlock": strings.Join(list, "\n")})
	if err != nil {
		fmt.Println("error: genModelSortMethodStructType() failed")
		log.Fatal(err)
	}
	return s
}

func genModelSortMethodsFunc(tab Table) string {
	list := make([]string, len(tab.Columns))
	for i, col := range tab.Columns {
		list[i] = genModelCompareMethod(col)
	}
	s, err := nbfmt.Fmt(modelSortMethodsStructFuncTemp, map[string]interface{}{"Table": tab, "CompareMethodBlock": strings.Join(list, "\n")})
	if err != nil {
		fmt.Println("error: genModelSortMethodFunc() failed")
		log.Fatal(err)
	}
	return s
}

func genModelListLenMethod(tab Table) string {
	s, err := nbfmt.Fmt(modelListLenMethodTemp, map[string]interface{}{"Table": tab})
	if err != nil {
		fmt.Println("error: genModelListLenMethod() failed")
		log.Fatal(err)
	}
	return s
}

func genModelListSwapMethod(tab Table) string {
	s, err := nbfmt.Fmt(modelListSwapMethodTemp, map[string]interface{}{"Table": tab})
	if err != nil {
		fmt.Println("error: genModelListSwapMethod() failed")
		log.Fatal(err)
	}
	return s
}

func genModelListLessMethod(tab Table) string {
	s, err := nbfmt.Fmt(modelListLessMethodTemp, map[string]interface{}{"Table": tab})
	if err != nil {
		fmt.Println("error: genModelListLessMethod() failed")
		log.Fatal(err)
	}
	return s
}

func genModelSortFuncSwitchBlock(col Column, tab Table) string {
	s, err := nbfmt.Fmt(modelSortFuncSwitchBlockTemp, map[string]interface{}{"Column": col, "Table": tab})
	if err != nil {
		fmt.Println("error: genModelSortFuncSwitchBlock() failed")
		log.Fatal(err)
	}
	return s
}

func genModelSortFunc(tab Table) string {
	list := make([]string, len(tab.Columns))
	for i, col := range tab.Columns {
		list[i] = genModelSortFuncSwitchBlock(col, tab)
	}
	s, err := nbfmt.Fmt(modelSortFuncTemp, map[string]interface{}{"Table": tab, "ColumnNumber": len(tab.Columns), "SwitchBlock": strings.Join(list, "\n")})
	if err != nil {
		fmt.Println("error: genModelSortFunc() failed")
		log.Fatal(err)
	}
	return s
}

func genModelCountStmtDeclare(tab Table) string {
	s, err := nbfmt.Fmt(modelCountStmtDeclareTemp, map[string]interface{}{"Table": tab})
	if err != nil {
		fmt.Println("error: genModelCountStmtDeclare() failed")
		log.Fatal(err)
	}
	return s
}

func genModelCountStmtInit(tab Table, db Database) string {
	s, err := nbfmt.Fmt(modelCountStmtInitTemp, map[string]interface{}{"Table": tab, "Database": db})
	if err != nil {
		fmt.Println("error: genModelCountStmtInit() failed")
		log.Fatal(err)
	}
	return s
}

func genModelCountFunc(tab Table) string {
	s, err := nbfmt.Fmt(modelCountFuncTemp, map[string]interface{}{"Table": tab})
	if err != nil {
		fmt.Println("error: genModelCountFunc() failed")
		log.Fatal(err)
	}
	return s
}

//Gen generate database definition
func Gen(db Database, outName string) error {
	buf := bytes.NewBuffer([]byte{})
	buf.WriteString(genPackage(db.Package) + "\n\n")
	buf.WriteString(importTemp + "\n\n")
	buf.WriteString(genDb(db) + "\n\n")
	initStmtBuf := bytes.NewBuffer([]byte{})
	for _, tab := range db.Tables {
		buf.WriteString(genStmtVar(tab) + "\n\n")
		initStmtBuf.WriteString(genStmtInitBlock(tab, db) + "\n\n")
	}
	buf.WriteString(genInitFunc(db, initStmtBuf.String()) + "\n\n")
	for _, tab := range db.Tables {
		buf.WriteString(genQueryFieldMapBlock(tab) + "\n\n")
		buf.WriteString(genModel(tab) + "\n\n")
		for _, col := range tab.Columns {
			buf.WriteString(genModelGetFieldMethod(col, tab) + "\n\n")
			buf.WriteString(genModelSetFieldMethod(col, tab) + "\n\n")
		}
		for _, fk := range tab.ForeignKeys {
			buf.WriteString(genForeignKeyMethod(fk, tab, db) + "\n\n")
		}
		for _, rfk := range tab.ReverseForeignKeys {
			buf.WriteString(genReverseForeignKeyStruct(rfk, tab, db) + "\n\n")
			buf.WriteString(genReverseForeignKeyMethod(rfk, tab, db) + "\n\n")
		}
		for _, mtm := range tab.ManyToManys {
			buf.WriteString(genManyToManyStruct(mtm, tab, db) + "\n\n")
			buf.WriteString(genManyToManyMethod(mtm, tab, db) + "\n\n")
		}
		buf.WriteString(genNewFunc(tab) + "\n\n")
		buf.WriteString(genAllFunc(tab, db) + "\n\n")
		buf.WriteString(genQueryFunc(tab, db) + "\n\n")
		buf.WriteString(genQueryOneFunc(tab, db) + "\n\n")
		buf.WriteString(genModelInsertMethod(tab) + "\n\n")
		buf.WriteString(genModelInsertOrUpdateMethod(tab) + "\n\n")
		buf.WriteString(genModelUpdateMethod(tab) + "\n\n")
		buf.WriteString(genModelDeleteMethod(tab) + "\n\n")
		buf.WriteString(genModelCountFunc(tab) + "\n\n")
		buf.WriteString(genFromRowsFunc(tab) + "\n\n")
		buf.WriteString(genFromRowFunc(tab) + "\n\n")
		buf.WriteString(genModelCheckMethod(tab) + "\n\n")
		buf.WriteString(genModelListType(tab) + "\n\n")
		buf.WriteString(genModelListLenMethod(tab) + "\n\n")
		buf.WriteString(genModelListSwapMethod(tab) + "\n\n")
		buf.WriteString(genModelListLessMethod(tab) + "\n\n")
		buf.WriteString(genModelSortMethodsStructType(tab) + "\n\n")
		buf.WriteString(genModelSortMethodsFunc(tab) + "\n\n")
		buf.WriteString(genModelSortFunc(tab) + "\n\n")
	}
	f, err := os.OpenFile(outName, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	defer f.Close()
	n, err := f.Write(buf.Bytes())
	if err != nil {
		return err
	}
	err = f.Truncate(int64(n))
	if err != nil {
		return err
	}
	cmd := exec.Command("go", "fmt", outName)
	return cmd.Run()
}
