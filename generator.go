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

type genInfo struct {
	column   Column
	table    Table
	database Database
}

//genPackage generate package statement
func genPackage(pkg string) string {
	return fmt.Sprintf(packageTemp, pkg)
}

//genDb generate database object declare statement
func genDb(db Database) string {
	return fmt.Sprintf(dbTemp, db.ObjName)
}

//genInitFunc generate init function
func genInitFunc(db Database, stmtInit string) string {
	return fmt.Sprintf(initFuncTemp, db.Username, db.Password, db.Address, db.DatabaseName, db.ObjName, stmtInit)
}

//genMapElemBlock generate map element
func genMapElemBlock(col Column) string {
	return fmt.Sprintf(mapElemTemp, "@"+col.FieldName, BackQuote(col.ColumnName))
}

//genQueryFieldMapBlock generate a map for mapping filed name to column name in where clause in query
func genQueryFieldMapBlock(tab Table) string {
	elemList := make([]string, len(tab.Columns))
	for i, col := range tab.Columns {
		elemList[i] = genMapElemBlock(col)
	}
	return fmt.Sprintf(queryFieldMapTemp, tab.ModelName, strings.Join(elemList, "\n"))
}

//genField generate field definition in model struct
func genField(col Column) string {
	return fmt.Sprintf(fieldTemp, col.FieldName, col.FieldType)
}

//genModel generate model defination
func genModel(tab Table) string {
	list := make([]string, len(tab.Columns))
	for i, col := range tab.Columns {
		list[i] = genField(col)
	}
	return fmt.Sprintf(modelTemp, tab.ModelName, strings.Join(list, "\n"))
}

//genModelGetFieldMethod generate model get xxx field value method
func genModelGetFieldMethod(col Column, tab Table) string {
	return fmt.Sprintf(modelGetFieldMethodTemp,
		tab.ModelName,
		col.FieldName,
		col.FieldType,
		col.FieldName,
		func() string {
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
		}(),
		col.FieldName)
}

//genModelSetFieldMethod generate model set xxx field value method
func genModelSetFieldMethod(col Column, tab Table) string {
	return fmt.Sprintf(modelSetFieldMethodTemp,
		tab.ModelName,
		col.FieldName,
		col.FieldType,
		col.FieldName,
		col.FieldName)
}

//genNewFuncArgList generate argument list in NewXXX() function
func genNewFuncArgList(tab Table) string {
	list := make([]string, 0, len(tab.Columns))
	for _, col := range tab.Columns {
		if !col.AutoIncrement && col.Default != "CURRENT_TIMESTAMP" && col.On != "UPDATE CURRENT_TIMESTAMP" {
			list = append(list, fmt.Sprintf(funcArgTemp, col.ArgName, "*"+col.FieldType))
		}
	}
	return strings.Join(list, ", ")
}

//genNewFuncAsignBlock generate middle variable asign block in NewXXX() function
func genNewFuncAsignBlock(tab Table) string {
	list := make([]string, 0, len(tab.Columns))
	for _, col := range tab.Columns {
		if !col.AutoIncrement && col.Default != "CURRENT_TIMESTAMP" && col.On != "UPDATE CURRENT_TIMESTAMP" {
			list = append(list, fmt.Sprintf(newModelAsignTemp, col.FieldName, col.ArgName))
		}
	}
	return strings.Join(list, ",\n")
}

//genNewFunc generate NewXXX() function
func genNewFunc(tab Table) string {
	return fmt.Sprintf(newModelFuncTemp,
		tab.ModelName,
		genNewFuncArgList(tab),
		tab.ModelName,
		tab.ModelName,
		genNewFuncAsignBlock(tab))
}

//genAllFunc generate AllXXX() function
func genAllFunc(tab Table, db Database) string {
	return fmt.Sprintf(allModelFuncTemp, tab.ModelName, tab.ModelName, db.ObjName, BackQuote(tab.TableName), tab.ModelName, tab.ModelName)
}

//genQueryFunc generate QueryXXX() function
func genQueryFunc(tab Table, db Database) string {
	return fmt.Sprintf(queryModelFuncTemp, tab.ModelName, tab.ModelName, tab.ModelName, db.ObjName, BackQuote(tab.TableName), tab.ModelName, tab.ModelName)
}

//genQueryOneFunc generate QueryOneXXX() function
func genQueryOneFunc(tab Table, db Database) string {
	return fmt.Sprintf(queryOneFuncTemp,
		tab.ModelName,
		tab.ModelName,
		tab.ModelName,
		db.ObjName,
		BackQuote(tab.TableName),
		tab.ModelName)
}

//genForeignKeyMethod generate foreign key relation method
func genForeignKeyMethod(fk ForeignKey, srcTab Table, db Database) string {
	return fmt.Sprintf(foreignKeyMethodTemp,
		srcTab.ModelName,
		fk.DstTab.ModelName,
		fk.SrcCol.FieldName,
		fk.DstTab.ModelName,
		fk.SrcCol.FieldName,
		srcTab.ModelName,
		fk.SrcCol.FieldName,
		db.ObjName,
		fmt.Sprintf(foreignKeyQuerySQLTemp, BackQuote(fk.DstTab.TableName), BackQuote(fk.DstCol.ColumnName)),
		fk.SrcCol.FieldName,
		fk.DstTab.ModelName)
}

//genReverseForeignKeyStruct generate reverse foreign key relation method
func genReverseForeignKeyStruct(rfk ReverseForeignKey, srcTab Table, db Database) string {
	return fmt.Sprintf(reverseForeignKeyStructTypeTemp,
		srcTab.ModelName,
		rfk.DstTab.ModelName,
		rfk.DstTab.ModelName,
		rfk.DstTab.ModelName)
}

//genReverseForeignKeyAllMethod generate All() method in reverse foreign key relation struct
func genReverseForeignKeyAllMethod(rfk ReverseForeignKey, srcTab Table, db Database) string {
	return fmt.Sprintf(reverseForeignKeyAllMethodTemp,
		rfk.DstTab.ModelName,
		rfk.SrcCol.FieldName,
		srcTab.ModelName,
		rfk.SrcCol.FieldName,
		db.ObjName,
		fmt.Sprintf(reverseForeignKeyAllSQLTemp, BackQuote(rfk.DstTab.TableName), BackQuote(rfk.DstCol.ColumnName)),
		rfk.SrcCol.FieldName,
		rfk.DstTab.ModelName,
		rfk.DstTab.ModelName)
}

//genReverseForeignKeyQueryMethod generate Query() method in reverse foreign key relation struct
func genReverseForeignKeyQueryMethod(rfk ReverseForeignKey, srcTab Table, db Database) string {
	return fmt.Sprintf(reverseForeignKeyQueryMethodTemp,
		rfk.DstTab.ModelName,
		rfk.SrcCol.FieldName,
		srcTab.ModelName,
		rfk.SrcCol.FieldName,
		rfk.DstTab.ModelName,
		db.ObjName,
		fmt.Sprintf(reverseForeignKeyQuerySQLTemp, BackQuote(rfk.DstTab.TableName), BackQuote(rfk.DstCol.ColumnName)),
		rfk.SrcCol.FieldName,
		rfk.DstTab.ModelName,
		rfk.DstTab.ModelName)
}

//genReverseForeignKeyMethod generate model reverse foreign key relation method
func genReverseForeignKeyMethod(rfk ReverseForeignKey, srcTab Table, db Database) string {
	allMethod := genReverseForeignKeyAllMethod(rfk, srcTab, db)
	queryMethod := genReverseForeignKeyQueryMethod(rfk, srcTab, db)
	return fmt.Sprintf(reverseForeignKeyMethodTemp,
		srcTab.ModelName,
		rfk.DstTab.ModelName,
		rfk.SrcCol.FieldName,
		srcTab.ModelName,
		rfk.DstTab.ModelName,
		srcTab.ModelName,
		rfk.DstTab.ModelName,
		allMethod,
		queryMethod)
}

//genManyToManyStruct generate many to many relation struct definition
func genManyToManyStruct(mtm ManyToMany, srcTab Table, db Database) string {
	return fmt.Sprintf(manyToManyStructTypeTemp,
		srcTab.ModelName,
		mtm.DstTab.ModelName,
		mtm.DstTab.ModelName,
		mtm.DstTab.ModelName,
		mtm.DstTab.ArgName,
		mtm.DstTab.ModelName,
		mtm.DstTab.ArgName,
		mtm.DstTab.ModelName)
}

//genManyToManyAllMethod generate All() method in many to many relation struct
func genManyToManyAllMethod(mtm ManyToMany, srcTab Table, db Database) string {
	return fmt.Sprintf(manyToManyAllMethodTemp,
		mtm.DstTab.ModelName,
		db.ObjName,
		fmt.Sprintf(manyToManyAllSQLTemp, BackQuote(mtm.DstTab.TableName), BackQuote(srcTab.TableName), BackQuote(mtm.MidTab.TableName),
			BackQuote(srcTab.TableName), BackQuote(mtm.SrcCol.ColumnName), BackQuote(mtm.MidTab.TableName), BackQuote(mtm.MidLeftCol.ColumnName),
			BackQuote(mtm.DstTab.TableName), BackQuote(mtm.MidTab.TableName), BackQuote(mtm.MidRightCol.ColumnName), BackQuote(mtm.DstTab.TableName),
			BackQuote(mtm.DstCol.ColumnName), BackQuote(srcTab.TableName), BackQuote(mtm.SrcCol.ColumnName)),
		mtm.SrcCol.FieldName,
		mtm.DstTab.ModelName,
		mtm.DstTab.ModelName)
}

//genManyToManyQueryMethod generate Query() method in many to many relation struct
func genManyToManyQueryMethod(mtm ManyToMany, srcTab Table, db Database) string {
	return fmt.Sprintf(manyToManyQueryMethodTemp,
		mtm.DstTab.ModelName,
		mtm.DstTab.ModelName,
		db.ObjName,
		fmt.Sprintf(manyToManyQuerySQLTemp, BackQuote(mtm.DstTab.TableName), BackQuote(srcTab.TableName), BackQuote(mtm.MidTab.TableName),
			BackQuote(srcTab.TableName), BackQuote(mtm.SrcCol.ColumnName), BackQuote(mtm.MidTab.TableName), BackQuote(mtm.MidLeftCol.ColumnName),
			BackQuote(mtm.DstTab.TableName), BackQuote(mtm.MidTab.TableName), BackQuote(mtm.MidRightCol.ColumnName), BackQuote(mtm.DstTab.TableName),
			BackQuote(mtm.DstCol.ColumnName), BackQuote(srcTab.TableName), BackQuote(mtm.SrcCol.ColumnName)),
		mtm.SrcCol.FieldName,
		mtm.DstTab.ModelName,
		mtm.DstTab.ModelName)
}

//genManyToManyAddMethod generate Add() method in many to many relation struct
func genManyToManyAddMethod(mtm ManyToMany, tab Table) string {
	return fmt.Sprintf(manyToManyAddMethodTemp,
		mtm.DstTab.ArgName,
		mtm.DstTab.ModelName,
		tab.ModelName,
		mtm.DstTab.ArgName,
		mtm.DstTab.ModelName,
		tab.ModelName,
		mtm.DstTab.ModelName,
		mtm.SrcCol.FieldName,
		mtm.DstTab.ArgName,
		mtm.DstCol.FieldName)
}

//genManyToManyRemoveMethod generate Remove() method in many to many relation struct
func genManyToManyRemoveMethod(mtm ManyToMany, tab Table) string {
	return fmt.Sprintf(manyToManyRemoveMethodTemp,
		mtm.DstTab.ArgName,
		mtm.DstTab.ModelName,
		tab.ModelName,
		mtm.DstTab.ArgName,
		mtm.DstTab.ModelName,
		tab.ModelName,
		mtm.DstTab.ModelName,
		mtm.SrcCol.FieldName,
		mtm.DstTab.ArgName,
		mtm.DstCol.FieldName)
}

//genManyToManyMethod generate model many to many relation method
func genManyToManyMethod(mtm ManyToMany, srcTab Table, db Database) string {
	allMethod := genManyToManyAllMethod(mtm, srcTab, db)
	queryMethod := genManyToManyQueryMethod(mtm, srcTab, db)
	addMethod := genManyToManyAddMethod(mtm, srcTab)
	removeMethod := genManyToManyRemoveMethod(mtm, srcTab)
	return fmt.Sprintf(manyToManyMethodTemp,
		srcTab.ModelName,
		mtm.DstTab.ModelName,
		mtm.SrcCol.FieldName,
		srcTab.ModelName,
		mtm.DstTab.ModelName,
		srcTab.ModelName,
		mtm.DstTab.ModelName,
		allMethod,
		queryMethod,
		addMethod,
		removeMethod)
}

//genModelInsertMethod generate XXX.Insert() method
func genModelInsertMethod(tab Table) string {
	return fmt.Sprintf(modelInsertMethodTemp,
		tab.ModelName,
		genStmtArgNilToDefaultBlock(tab),
		tab.ModelName,
		tab.AutoIncrement.FieldName)
}

//genModelUpdateMethod generate XXX.Update() method
func genModelUpdateMethod(tab Table) string {
	return fmt.Sprintf(modelUpdateMethodTemp,
		tab.ModelName,
		genStmtArgNilToDefaultBlock(tab),
		tab.AutoIncrement.FieldName,
		tab.ModelName)
}

//genModelExistsMethod generate XXX.Exists() method
func genModelExistsMethod(tab Table, db Database) string {
	return fmt.Sprintf(modelExistsMethodTemp,
		tab.ModelName,
		tab.PrimaryKey.FieldName,
		tab.ModelName,
		tab.PrimaryKey.FieldName,
		db.ObjName,
		fmt.Sprintf(queryByPrimaryKeySQLTemp, BackQuote(tab.TableName), BackQuote(tab.PrimaryKey.ColumnName)),
		tab.PrimaryKey.FieldName)
}

//genModelInsertOrUpdateMethod generate XXX.InsertOrUpdate() method
func genModelInsertOrUpdateMethod(tab Table) string {
	return fmt.Sprintf(modelInsertOrUpdateMethodTemp,
		tab.ModelName,
		genStmtArgNilToDefaultBlock(tab),
		tab.ModelName,
		tab.AutoIncrement.FieldName,
	)
}

//genModelDeleteMethod generate XXX.Delete() method
func genModelDeleteMethod(tab Table) string {
	return fmt.Sprintf(modelDeleteMethodTemp,
		tab.ModelName,
		tab.ModelName,
		fmt.Sprintf(deleteArgTemp, tab.AutoIncrement.FieldName))
}

//genNewMiddleTypeBlock generate new middle type consponsed to field type block
func genNewMiddleTypeBlock(col Column) string {
	return fmt.Sprintf(newMiddleTypeTemp, col.ArgName, col.MidType)
}

//genFromRowsCheckBlock generate check block in XXXFromRows() function
func genFromRowsCheckBlock(col Column) string {
	return fmt.Sprintf(modelFromRowsCheckNullBlockTemp, col.ArgName, col.ArgName, col.ArgName)
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
	return fmt.Sprintf(modelFromRowsFuncTemp, tab.ModelName, tab.ModelName, strings.Join(midList, "\n"), strings.Join(midNameList, ", "),
		tab.ModelName, strings.Join(finalArgList, ", "))
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
	return fmt.Sprintf(modelFromRowFuncTemp, tab.ModelName, tab.ModelName, strings.Join(midList, "\n"), strings.Join(midNameList, ", "),
		tab.ModelName, strings.Join(finalArgList, ", "))
}

//genModelCheckMethod generate XXX.check() method.This method is for excluding AUTO_INCREMTN column and AUTO_TIMESTAMP column from insert and update method
func genModelCheckMethod(tab Table) string {
	list := make([]string, 0, len(tab.Columns))
	for _, col := range tab.Columns {
		if !col.Nullable && !col.AutoIncrement && col.Default == "" {
			list = append(list, fmt.Sprintf(fieldCheckNullTemp, col.FieldName, tab.ModelName, col.FieldName))
		}
	}
	return fmt.Sprintf(modelCheckMethodTemp, tab.ModelName, strings.Join(list, "\n"))
}

//genInsertStmt generate insert sql.Stmt declaration
func genInsertStmt(tab Table) string {
	return fmt.Sprintf(insertStmtTemp, tab.ModelName)
}

//genUpdateStmt generate update sql.Stmt declaration
func genUpdateStmt(tab Table) string {
	return fmt.Sprintf(updateStmtTemp, tab.ModelName)
}

//genDeleteStmt generate delete sql.Stmt declaration
func genDeleteStmt(tab Table) string {
	return fmt.Sprintf(deleteStmtTemp, tab.ModelName)
}

//genInsertOrUpdateStmt generate insert or update sql.Stmt declaration
func genInsertOrUpdateStmt(tab Table) string {
	return fmt.Sprintf(insertOrUpdateStmtTemp, tab.ModelName)
}

//genInsertMidStmt generate insert into middle table sql.Stmt declaration for many to many relation
func genInsertMidStmt(tab Table, mtm ManyToMany) string {
	return fmt.Sprintf(insertStmtTemp, tab.ModelName+"To"+mtm.DstTab.ModelName)
}

//genManyToManyDeleteStmt generate delete sql.Stmt declaration for many to many relation
func genManyToManyDeleteStmt(mtm ManyToMany, tab Table) string {
	return fmt.Sprintf(manyToManyDeleteStmtTemp, tab.ModelName, mtm.DstTab.ModelName)
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
			list = append(list, BackQuote(col.ColumnName))
		}
	}
	return fmt.Sprintf(insertStmtInitTemp,
		tab.ModelName,
		db.ObjName,
		BackQuote(tab.TableName),
		strings.Join(list, ", "),
		strings.Trim(strings.Repeat("?, ", len(list)), ", "))
}

//genUpdateStmtInitBlock generate update sql.Stmt in init function
func genUpdateStmtInitBlock(tab Table, db Database) string {
	list := make([]string, 0, len(tab.Columns))
	for _, col := range tab.Columns {
		if !col.AutoIncrement && col.Default != "CURRENT_TIMESTAMP" && col.On != "UPDATE CURRENT_TIMESTAMP" {
			list = append(list, BackQuote(col.ColumnName)+" = ?")
		}
	}
	return fmt.Sprintf(updateStmtInitTemp,
		tab.ModelName,
		db.ObjName,
		BackQuote(tab.TableName),
		strings.Join(list, ", "),
		BackQuote(tab.AutoIncrement.ColumnName))
}

//genDeleteStmtInitBlock generate delete sql.Stmt in init function
func genDeleteStmtInitBlock(tab Table, db Database) string {
	return fmt.Sprintf(deleteStmtInitTemp,
		tab.ModelName,
		db.ObjName,
		BackQuote(tab.TableName),
		BackQuote(tab.AutoIncrement.ColumnName))
}

//genInsertOrUpdateStmtInitBlock generate insert or update sql.Stmt in init function
func genInsertOrUpdateStmtInitBlock(tab Table, db Database) string {
	insertList := make([]string, 0, len(tab.Columns))
	argList := make([]string, 0, len(tab.Columns))
	updateList := make([]string, 0, len(tab.Columns))
	for _, col := range tab.Columns {
		if !col.AutoIncrement && col.Default != "CURRENT_TIMESTAMP" && col.On != "UPDATE CURRENT_TIMESTAMP" {
			insertList = append(insertList, BackQuote(col.ColumnName))
			argList = append(argList, "?")
			updateList = append(updateList, fmt.Sprintf(updateColumnTemp, BackQuote(col.ColumnName)))
		}
	}
	return fmt.Sprintf(insertOrUpdateStmtInitTemp,
		tab.ModelName,
		db.ObjName,
		BackQuote(tab.TableName),
		strings.Join(insertList, ", "),
		strings.Join(argList, ", "),
		fmt.Sprintf(updateLastInsertIDTemp, BackQuote(tab.AutoIncrement.ColumnName), BackQuote(tab.AutoIncrement.ColumnName)),
		strings.Join(updateList, ", "))
}

//genInsertMidStmtInitBlock generate insert into middle table sql.Stmt for many to many relation in init function
func genInsertMidStmtInitBlock(mtm ManyToMany, tab Table, db Database) string {
	return fmt.Sprintf(insertMidStmtInitTemp,
		tab.ModelName+"To"+mtm.DstTab.ModelName,
		db.ObjName,
		BackQuote(mtm.MidTab.TableName),
		BackQuote(mtm.MidLeftCol.ColumnName),
		BackQuote(mtm.MidRightCol.ColumnName))
}

//genManyToManyDeleteStmtInitBlock generate delete middle table sql.Stmt for many to many relation in init function
func genManyToManyDeleteStmtInitBlock(mtm ManyToMany, tab Table, db Database) string {
	return fmt.Sprintf(manyToManyDeleteStmtInitTemp,
		tab.ModelName,
		mtm.DstTab.ModelName,
		db.ObjName,
		BackQuote(mtm.MidTab.TableName),
		BackQuote(mtm.MidLeftCol.ColumnName),
		BackQuote(mtm.MidRightCol.ColumnName))
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
	return fmt.Sprintf(stmtArgNilToDefaultTemp, col.FieldName, value, col.FieldName)
}

//genStmtArgNilToDefaultBlock generate block to replacing nil arguments to default value
func genStmtArgNilToDefaultBlock(tab Table) string {
	list := make([]string, 0, len(tab.Columns))
	for _, col := range tab.Columns {
		if !col.AutoIncrement && col.Default != "CURRENT_TIMESTAMP" && col.On != "UPDATE CURRENT_TIMESTAMP" {
			if !col.Nullable && col.Default != "" {
				list = append(list, genStmtArgNilToDefault(col))
			} else {
				list = append(list, fmt.Sprintf(stmtArgTemp, col.FieldName))
			}
		}
	}
	return fmt.Sprintf(stmtArgNilToDefaultBlockTemp, len(list), strings.Join(list, "\n"))
}

//genModelListType generate XXXList struct definition
func genModelListType(tab Table) string {
	return fmt.Sprintf(modelListTypeTemp, tab.ModelName, tab.ModelName)
}

//genModelCompareMethod generate field compare method for sort function
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
	return fmt.Sprintf(temp, col.FieldName, col.FieldName, col.FieldName, col.FieldName, col.FieldName)
}

//genModelSortMethodsStructType generate model sort struct
func genModelSortMethodsStructType(tab Table) string {
	list := make([]string, len(tab.Columns))
	for i, col := range tab.Columns {
		list[i] = fmt.Sprintf(modelSortMethodsStructFieldTypeTemp, col.FieldName)
	}
	return fmt.Sprintf(modelSortMethodsStructTypeTemp, tab.ModelName, strings.Join(list, "\n"))
}

//genModelSortMethodsFunc generate model sort method
func genModelSortMethodsFunc(tab Table) string {
	list := make([]string, len(tab.Columns))
	for i, col := range tab.Columns {
		list[i] = genModelCompareMethod(col)
	}
	return fmt.Sprintf(modelSortMethodsStructFuncTemp, tab.ModelName, tab.ModelName, tab.ModelName, strings.Join(list, "\n"))
}

//genModelListLenMethod generate XXXList.Len() method
func genModelListLenMethod(tab Table) string {
	return fmt.Sprintf(modelListLenMethodTemp, tab.ModelName)
}

//genModelListSwapMethod generate XXXList.Swap() method
func genModelListSwapMethod(tab Table) string {
	return fmt.Sprintf(modelListSwapMethodTemp, tab.ModelName)
}

//genModelListLessMethod generate XXXList.Less() method
func genModelListLessMethod(tab Table) string {
	return fmt.Sprintf(modelListLessMethodTemp, tab.ModelName)
}

//genModelSortFuncSwitchBlock generate XXXSortBy() switch block
func genModelSortFuncSwitchBlock(col Column, tab Table) string {
	return fmt.Sprintf(modelSortFuncSwitchBlockTemp,
		col.FieldName,
		tab.ArgName,
		tab.ArgName,
		col.FieldName)
}

//genModelSortFunc generate XXXSortBy() function
func genModelSortFunc(tab Table) string {
	list := make([]string, len(tab.Columns))
	for i, col := range tab.Columns {
		list[i] = genModelSortFuncSwitchBlock(col, tab)
	}
	return fmt.Sprintf(modelSortFuncTemp,
		tab.ModelName,
		tab.ModelName,
		tab.ArgName,
		tab.ModelName,
		len(tab.Columns),
		tab.ArgName,
		strings.Join(list, "\n"),
		tab.ArgName,
		tab.ArgName)
}

func genModelCountStmtDeclare(tab Table) string {
	return fmt.Sprintf(modelCountStmtDeclareTemp, tab.ModelName)
}

// func genModelCountStmtInit(tab Table, db Database) string {
// 	return fmt.Sprintf(modelCountStmtInitTemp, tab.ModelName, db.ObjName, BackQuote(tab.TableName))
// }

func genModelCountStmtInit(tab Table, db Database) string {
	s, err := nbfmt.Fmt(modelCountStmtInitTemp, map[string]interface{}{"Table": tab, "Database": db})
	if err != nil {
		log.Fatal(err)
	}
	return s
}

// func genModelCountFunc(tab Table) string {
// 	return fmt.Sprintf(modelCountFuncTemp, tab.ModelName, tab.ModelName)
// }
func genModelCountFunc(tab Table) string {
	s, err := nbfmt.Fmt(modelCountFuncTemp, map[string]interface{}{"Table": tab})
	if err != nil {
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
		buf.WriteString(genModelExistsMethod(tab, db) + "\n\n")
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
