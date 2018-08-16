package nbmysql

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

//GenPackage Generate package statement
func GenPackage(pkg string) string {
	return fmt.Sprintf(PackageTemp, pkg)
}

//GenDb Generate database object declare statement
func GenDb(db Database) string {
	return fmt.Sprintf(DbTemp, db.ObjName)
}

//GenInitFunc Generate init function
func GenInitFunc(db Database, stmtInit string) string {
	return fmt.Sprintf(InitFuncTemp, db.Username, db.Password, db.Address, db.DatabaseName, db.ObjName, stmtInit)
}

//GenMapElemBlock Generate map element
func GenMapElemBlock(col Column) string {
	return fmt.Sprintf(MapElemTemp, "@"+col.FieldName, BackQuote(col.ColumnName))
}

//GenQueryFieldMapBlock Generate a map for mapping filed name to column name in where clause in query
func GenQueryFieldMapBlock(tab Table) string {
	elemList := make([]string, len(tab.Columns))
	for i, col := range tab.Columns {
		elemList[i] = GenMapElemBlock(col)
	}
	return fmt.Sprintf(QueryFieldMapTemp, tab.ModelName, strings.Join(elemList, "\n"))
}

//GenField Generate field definition in model struct
func GenField(col Column) string {
	return fmt.Sprintf(FieldTemp, col.FieldName, col.FieldType)
}

//GenModel Generate model defination
func GenModel(tab Table) string {
	list := make([]string, len(tab.Columns))
	for i, col := range tab.Columns {
		list[i] = GenField(col)
	}
	return fmt.Sprintf(ModelTemp, tab.ModelName, strings.Join(list, "\n"))
}

//GenModelGetFieldMethod Generate model get xxx field value method
func GenModelGetFieldMethod(col Column, tab Table) string {
	return fmt.Sprintf(ModelGetFieldMethodTemp,
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

//GenModelSetFieldMethod Generate model set xxx field value method
func GenModelSetFieldMethod(col Column, tab Table) string {
	return fmt.Sprintf(ModelSetFieldMethodTemp,
		tab.ModelName,
		col.FieldName,
		col.FieldType,
		col.FieldName,
		col.FieldName)
}

//GenNewFuncArgList Generate argument list in NewXXX() function
func GenNewFuncArgList(tab Table) string {
	list := make([]string, 0, len(tab.Columns))
	for _, col := range tab.Columns {
		if !col.AutoIncrement && col.Default != "CURRENT_TIMESTAMP" && col.On != "UPDATE CURRENT_TIMESTAMP" {
			list = append(list, fmt.Sprintf(FuncArgTemp, col.ArgName, "*"+col.FieldType))
		}
	}
	return strings.Join(list, ", ")
}

//GenNewFuncAsignBlock Generate middle variable asign block in NewXXX() function
func GenNewFuncAsignBlock(tab Table) string {
	list := make([]string, 0, len(tab.Columns))
	for _, col := range tab.Columns {
		if !col.AutoIncrement && col.Default != "CURRENT_TIMESTAMP" && col.On != "UPDATE CURRENT_TIMESTAMP" {
			list = append(list, fmt.Sprintf(NewModelAsignTemp, col.FieldName, col.ArgName))
		}
	}
	return strings.Join(list, ",\n")
}

//GenNewFunc Generate NewXXX() function
func GenNewFunc(tab Table) string {
	return fmt.Sprintf(NewModelFuncTemp,
		tab.ModelName,
		GenNewFuncArgList(tab),
		tab.ModelName,
		tab.ModelName,
		GenNewFuncAsignBlock(tab))
}

//GenAllFunc Generate AllXXX() function
func GenAllFunc(tab Table, db Database) string {
	return fmt.Sprintf(AllModelFuncTemp, tab.ModelName, tab.ModelName, db.ObjName, BackQuote(tab.TableName), tab.ModelName, tab.ModelName)
}

//GenQueryFunc Generate QueryXXX() function
func GenQueryFunc(tab Table, db Database) string {
	return fmt.Sprintf(QueryModelFuncTemp, tab.ModelName, tab.ModelName, tab.ModelName, db.ObjName, BackQuote(tab.TableName), tab.ModelName, tab.ModelName)
}

//GenQueryOneFunc Generate QueryOneXXX() function
func GenQueryOneFunc(tab Table, db Database) string {
	return fmt.Sprintf(QueryOneFuncTemp,
		tab.ModelName,
		tab.ModelName,
		tab.ModelName,
		db.ObjName,
		BackQuote(tab.TableName),
		tab.ModelName)
}

//GenForeignKeyMethod Generate foreign key relation method
func GenForeignKeyMethod(fk ForeignKey, srcTab Table, db Database) string {
	return fmt.Sprintf(ForeignKeyMethodTemp,
		srcTab.ModelName,
		fk.DstTab.ModelName,
		fk.SrcCol.FieldName,
		fk.DstTab.ModelName,
		fk.SrcCol.FieldName,
		srcTab.ModelName,
		fk.SrcCol.FieldName,
		db.ObjName,
		fmt.Sprintf(ForeignKeyQuerySQLTemp, BackQuote(fk.DstTab.TableName), BackQuote(fk.DstCol.ColumnName)),
		fk.SrcCol.FieldName,
		fk.DstTab.ModelName)
}

//GenReverseForeignKeyStruct Generate reverse foreign key relation method
func GenReverseForeignKeyStruct(rfk ReverseForeignKey, srcTab Table, db Database) string {
	return fmt.Sprintf(ReverseForeignKeyStructTypeTemp,
		srcTab.ModelName,
		rfk.DstTab.ModelName,
		rfk.DstTab.ModelName,
		rfk.DstTab.ModelName)
}

//GenReverseForeignKeyAllMethod Generate All() method in reverse foreign key relation struct
func GenReverseForeignKeyAllMethod(rfk ReverseForeignKey, srcTab Table, db Database) string {
	return fmt.Sprintf(ReverseForeignKeyAllMethodTemp,
		rfk.DstTab.ModelName,
		rfk.SrcCol.FieldName,
		srcTab.ModelName,
		rfk.SrcCol.FieldName,
		db.ObjName,
		fmt.Sprintf(ReverseForeignKeyAllSQLTemp, BackQuote(rfk.DstTab.TableName), BackQuote(rfk.DstCol.ColumnName)),
		rfk.SrcCol.FieldName,
		rfk.DstTab.ModelName,
		rfk.DstTab.ModelName)
}

//GenReverseForeignKeyQueryMethod Generate Query() method in reverse foreign key relation struct
func GenReverseForeignKeyQueryMethod(rfk ReverseForeignKey, srcTab Table, db Database) string {
	return fmt.Sprintf(ReverseForeignKeyQueryMethodTemp,
		rfk.DstTab.ModelName,
		rfk.SrcCol.FieldName,
		srcTab.ModelName,
		rfk.SrcCol.FieldName,
		rfk.DstTab.ModelName,
		db.ObjName,
		fmt.Sprintf(ReverseForeignKeyQuerySQLTemp, BackQuote(rfk.DstTab.TableName), BackQuote(rfk.DstCol.ColumnName)),
		rfk.SrcCol.FieldName,
		rfk.DstTab.ModelName,
		rfk.DstTab.ModelName)
}

//GenReverseForeignKeyMethod Generate model reverse foreign key relation method
func GenReverseForeignKeyMethod(rfk ReverseForeignKey, srcTab Table, db Database) string {
	allMethod := GenReverseForeignKeyAllMethod(rfk, srcTab, db)
	queryMethod := GenReverseForeignKeyQueryMethod(rfk, srcTab, db)
	return fmt.Sprintf(ReverseForeignKeyMethodTemp,
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

//GenManyToManyStruct Generate many to many relation struct definition
func GenManyToManyStruct(mtm ManyToMany, srcTab Table, db Database) string {
	return fmt.Sprintf(ManyToManyStructTypeTemp,
		srcTab.ModelName,
		mtm.DstTab.ModelName,
		mtm.DstTab.ModelName,
		mtm.DstTab.ModelName,
		mtm.DstTab.ArgName,
		mtm.DstTab.ModelName,
		mtm.DstTab.ArgName,
		mtm.DstTab.ModelName)
}

//GenManyToManyAllMethod Generate All() method in many to many relation struct
func GenManyToManyAllMethod(mtm ManyToMany, srcTab Table, db Database) string {
	return fmt.Sprintf(ManyToManyAllMethodTemp,
		mtm.DstTab.ModelName,
		db.ObjName,
		fmt.Sprintf(ManyToManyAllSQLTemp, BackQuote(mtm.DstTab.TableName), BackQuote(srcTab.TableName), BackQuote(mtm.MidTab.TableName),
			BackQuote(srcTab.TableName), BackQuote(mtm.SrcCol.ColumnName), BackQuote(mtm.MidTab.TableName), BackQuote(mtm.MidLeftCol.ColumnName),
			BackQuote(mtm.DstTab.TableName), BackQuote(mtm.MidTab.TableName), BackQuote(mtm.MidRightCol.ColumnName), BackQuote(mtm.DstTab.TableName),
			BackQuote(mtm.DstCol.ColumnName), BackQuote(srcTab.TableName), BackQuote(mtm.SrcCol.ColumnName)),
		mtm.SrcCol.FieldName,
		mtm.DstTab.ModelName,
		mtm.DstTab.ModelName)
}

//GenManyToManyQueryMethod Generate Query() method in many to many relation struct
func GenManyToManyQueryMethod(mtm ManyToMany, srcTab Table, db Database) string {
	return fmt.Sprintf(ManyToManyQueryMethodTemp,
		mtm.DstTab.ModelName,
		mtm.DstTab.ModelName,
		db.ObjName,
		fmt.Sprintf(ManyToManyQuerySQLTemp, BackQuote(mtm.DstTab.TableName), BackQuote(srcTab.TableName), BackQuote(mtm.MidTab.TableName),
			BackQuote(srcTab.TableName), BackQuote(mtm.SrcCol.ColumnName), BackQuote(mtm.MidTab.TableName), BackQuote(mtm.MidLeftCol.ColumnName),
			BackQuote(mtm.DstTab.TableName), BackQuote(mtm.MidTab.TableName), BackQuote(mtm.MidRightCol.ColumnName), BackQuote(mtm.DstTab.TableName),
			BackQuote(mtm.DstCol.ColumnName), BackQuote(srcTab.TableName), BackQuote(mtm.SrcCol.ColumnName)),
		mtm.SrcCol.FieldName,
		mtm.DstTab.ModelName,
		mtm.DstTab.ModelName)
}

//GenManyToManyAddMethod Generate Add() method in many to many relation struct
func GenManyToManyAddMethod(mtm ManyToMany, tab Table) string {
	return fmt.Sprintf(ManyToManyAddMethodTemp,
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

//GenManyToManyRemoveMethod Generate Remove() method in many to many relation struct
func GenManyToManyRemoveMethod(mtm ManyToMany, tab Table) string {
	return fmt.Sprintf(ManyToManyRemoveMethodTemp,
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

//GenManyToManyMethod Generate model many to many relation method
func GenManyToManyMethod(mtm ManyToMany, srcTab Table, db Database) string {
	allMethod := GenManyToManyAllMethod(mtm, srcTab, db)
	queryMethod := GenManyToManyQueryMethod(mtm, srcTab, db)
	addMethod := GenManyToManyAddMethod(mtm, srcTab)
	removeMethod := GenManyToManyRemoveMethod(mtm, srcTab)
	return fmt.Sprintf(ManyToManyMethodTemp,
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

// func GenCheckStringBlock(col Column, localModelName string) string {
// 	return fmt.Sprintf(ModelCheckStringBlockTemp, localModelName, col.FieldName, BackQuote(col.ColumnName), localModelName, col.FieldName)
// }

// func GenCheckIntBlock(col Column, localModelName string) string {
// 	return fmt.Sprintf(ModelCheckIntBlockTemp, localModelName, col.FieldName, BackQuote(col.ColumnName), localModelName, col.FieldName)
// }

// func GenCheckFloatBlock(col Column, localModelName string) string {
// 	return fmt.Sprintf(ModelCheckFloatBlockTemp, localModelName, col.FieldName, BackQuote(col.ColumnName), localModelName, col.FieldName)
// }

// func GenCheckTimeBlock(col Column, localModelName string) string {
// 	return fmt.Sprintf(ModelCheckTimeBlockTemp, localModelName, col.FieldName, BackQuote(col.ColumnName), localModelName, col.FieldName)
// }

// func GenCheckBoolBlock(col Column, localModelName string) string {
// 	return fmt.Sprintf(ModelCheckBoolBlockTemp, localModelName, col.FieldName, BackQuote(col.ColumnName), localModelName, col.FieldName)
// }

// func GenInsertSQL(tab Table) string {
// 	return fmt.Sprintf(InsertSQLTemp, BackQuote(tab.TableName))
// }

// func GenInsertArgsBlock(tab Table) string {
// 	list := make([]string, 0, len(tab.Columns))
// 	for _, col := range tab.Columns {
// 		if !col.AutoIncrement && col.Default != "CURRENT_TIMESTAMP" && col.On != "UPDATE CURRENT_TIMESTAMP" {
// 			list = append(list, fmt.Sprintf(ModelInsertArgTemp, col.FieldName))
// 		}
// 	}
// 	return strings.Join(list, ", ")
// }

//GenModelInsertMethod Generate XXX.Insert() method
func GenModelInsertMethod(tab Table) string {
	return fmt.Sprintf(ModelInsertMethodTemp,
		tab.ModelName,
		GenStmtArgNilToDefaultBlock(tab),
		tab.ModelName,
		tab.AutoIncrement.FieldName)
}

// func GenUpdateArgs(tab Table) string {
// 	list := make([]string, 0, len(tab.Columns))
// 	for _, col := range tab.Columns {
// 		if !col.AutoIncrement && col.Default != "CURRENT_TIMESTAMP" && col.On != "UPDATE CURRENT_TIMESTAMP" {
// 			list = append(list, fmt.Sprintf(UpdateArgTemp, col.FieldName))
// 		}
// 	}
// 	return strings.Join(list, ", ")
// }

//GenModelUpdateMethod Generate XXX.Update() method
func GenModelUpdateMethod(tab Table) string {
	return fmt.Sprintf(ModelUpdateMethodTemp,
		tab.ModelName,
		GenStmtArgNilToDefaultBlock(tab),
		tab.AutoIncrement.FieldName,
		tab.ModelName)
}

//GenModelExistsMethod Generate XXX.Exists() method
func GenModelExistsMethod(tab Table, db Database) string {
	return fmt.Sprintf(ModelExistsMethodTemp,
		tab.ModelName,
		tab.PrimaryKey.FieldName,
		tab.ModelName,
		tab.PrimaryKey.FieldName,
		db.ObjName,
		fmt.Sprintf(QueryByPrimaryKeySQLTemp, BackQuote(tab.TableName), BackQuote(tab.PrimaryKey.ColumnName)),
		tab.PrimaryKey.FieldName)
}

// func GenInsertOrUpdateArgs(tab Table) string {
// 	list := make([]string, 0, len(tab.Columns))
// 	for _, col := range tab.Columns {
// 		if !col.AutoIncrement && col.Default != "CURRENT_TIMESTAMP" && col.On != "UPDATE CURRENT_TIMESTAMP" {
// 			list = append(list, fmt.Sprintf(InsertOrUpdateArgTemp, col.FieldName))
// 		}
// 	}
// 	finalList := make([]string, len(list)*2)
// 	copy(finalList[:len(list)], list)
// 	copy(finalList[len(list):], list)
// 	return strings.Join(finalList, ", ")
// }

//GenModelInsertOrUpdateMethod Generate XXX.InsertOrUpdate() method
func GenModelInsertOrUpdateMethod(tab Table) string {
	return fmt.Sprintf(ModelInsertOrUpdateMethodTemp,
		tab.ModelName,
		GenStmtArgNilToDefaultBlock(tab),
		tab.ModelName,
		tab.AutoIncrement.FieldName,
	)
}

// func GenDeleteSQL(tab Table) string {
// 	return fmt.Sprintf(DeleteSQLTemp, BackQuote(tab.TableName), BackQuote(tab.PrimaryKey.ColumnName))
// }

// func GenManyToManyDeleteSQL(info ManyToMany) string {
// 	return fmt.Sprintf(DeleteSQLTemp, BackQuote(info.MidTab.TableName), BackQuote(info.MidLeftCol.ColumnName))
// }

// func GenManyToManyDeleteBlock(info ManyToMany) string {
// 	return fmt.Sprintf(ManyToManyDeleteBlockTemp, GenManyToManyDeleteSQL(info), info.SrcCol.FieldName)
// }

// func GenCasecadeDeleteLoop(rfk ReverseForeignKey, db Database) string {
// 	return fmt.Sprintf(CascadeDeleteLoopTemp,
// 		fmt.Sprintf(CascadeDeleteSQLTemp, BackQuote(rfk.DstTab.TableName), BackQuote(rfk.DstCol.ColumnName)),
// 		rfk.SrcCol.FieldName)
// }

// func GenDeleteFunc(tab Table, db Database) string {
// 	list := make([]string, len(tab.ManyToManys))
// 	for i, mtm := range tab.ManyToManys {
// 		list[i] = GenManyToManyDeleteBlock(mtm)
// 	}
// 	rfkList := make([]string, len(tab.ReverseForeignKeys))
// 	for i, rfk := range tab.ReverseForeignKeys {
// 		rfkList[i] = GenCasecadeDeleteLoop(rfk, db)
// 	}
// 	var cascadeDelete string
// 	if len(rfkList) > 0 {
// 		cascadeDelete = fmt.Sprintf(CascadeDeleteBlockTemp, strings.Join(rfkList, "\n"))
// 	}
// 	return fmt.Sprintf(ModelDeleteMethodTemp, tab.ModelName, db.ObjName, strings.Join(list, "\n"), cascadeDelete, GenDeleteSQL(tab), tab.PrimaryKey.FieldName)
// }

//GenModelDeleteMethod Generate XXX.Delete() method
func GenModelDeleteMethod(tab Table) string {
	return fmt.Sprintf(ModelDeleteMethodTemp,
		tab.ModelName,
		tab.ModelName,
		fmt.Sprintf(DeleteArgTemp, tab.AutoIncrement.FieldName))
}

//GenNewMiddleTypeBlock Generate new middle type consponsed to field type block
func GenNewMiddleTypeBlock(col Column) string {
	return fmt.Sprintf(NewMiddleTypeTemp, col.ArgName, col.MidType)
}

//GenFromRowsCheckBlock Generate check block in XXXFromRows() function
func GenFromRowsCheckBlock(col Column) string {
	return fmt.Sprintf(ModelFromRowsCheckNullBlockTemp, col.ArgName, col.ArgName, col.ArgName)
}

//GenFromRowsFunc Generate XXXFromRows(sql.Rows) ([]*XXX, error) function
func GenFromRowsFunc(tab Table) string {
	midList := make([]string, len(tab.Columns))
	midNameList := make([]string, len(tab.Columns))
	finalArgList := make([]string, len(tab.Columns))
	for i, col := range tab.Columns {
		midList[i] = GenNewMiddleTypeBlock(col)
		midNameList[i] = "_" + col.ArgName
		finalArgList[i] = "_" + col.ArgName + ".ToGo()"
	}
	return fmt.Sprintf(ModelFromRowsFuncTemp, tab.ModelName, tab.ModelName, strings.Join(midList, "\n"), strings.Join(midNameList, ", "),
		tab.ModelName, strings.Join(finalArgList, ", "))
}

//GenFromRowFunc Generate XXXFromRow(sql.Rows) (*XXX, error) function
func GenFromRowFunc(tab Table) string {
	midList := make([]string, len(tab.Columns))
	midNameList := make([]string, len(tab.Columns))
	finalArgList := make([]string, len(tab.Columns))
	for i, col := range tab.Columns {
		midList[i] = GenNewMiddleTypeBlock(col)
		midNameList[i] = "_" + col.ArgName
		finalArgList[i] = "_" + col.ArgName + ".ToGo()"
	}
	return fmt.Sprintf(ModelFromRowFuncTemp, tab.ModelName, tab.ModelName, strings.Join(midList, "\n"), strings.Join(midNameList, ", "),
		tab.ModelName, strings.Join(finalArgList, ", "))
}

//GenModelCheckMethod Generate XXX.check() method.This method is for excluding AUTO_INCREMTN column and AUTO_TIMESTAMP column from insert and update method
func GenModelCheckMethod(tab Table) string {
	list := make([]string, 0, len(tab.Columns))
	for _, col := range tab.Columns {
		if !col.Nullable && !col.AutoIncrement && col.Default == "" {
			list = append(list, fmt.Sprintf(FieldCheckNullTemp, col.FieldName, tab.ModelName, col.FieldName))
		}
	}
	return fmt.Sprintf(ModelCheckMethodTemp, tab.ModelName, strings.Join(list, "\n"))
}

//GenInsertStmt Generate insert sql.Stmt declaration
func GenInsertStmt(tab Table) string {
	return fmt.Sprintf(InsertStmtTemp, tab.ModelName)
}

//GenUpdateStmt Generate update sql.Stmt declaration
func GenUpdateStmt(tab Table) string {
	return fmt.Sprintf(UpdateStmtTemp, tab.ModelName)
}

//GenDeleteStmt Generate delete sql.Stmt declaration
func GenDeleteStmt(tab Table) string {
	return fmt.Sprintf(DeleteStmtTemp, tab.ModelName)
}

//GenInsertOrUpdateStmt Generate insert or update sql.Stmt declaration
func GenInsertOrUpdateStmt(tab Table) string {
	return fmt.Sprintf(InsertOrUpdateStmtTemp, tab.ModelName)
}

//GenInsertMidStmt Generate insert into middle table sql.Stmt declaration for many to many relation
func GenInsertMidStmt(tab Table, mtm ManyToMany) string {
	return fmt.Sprintf(InsertStmtTemp, tab.ModelName+"To"+mtm.DstTab.ModelName)
}

//GenManyToManyDeleteStmt Generate delete sql.Stmt declaration for many to many relation
func GenManyToManyDeleteStmt(mtm ManyToMany, tab Table) string {
	return fmt.Sprintf(ManyToManyDeleteStmtTemp, tab.ModelName, mtm.DstTab.ModelName)
}

//GenStmtVar Generate prepared sql.Stmt declaration
func GenStmtVar(tab Table) string {
	list := make([]string, 0, 8)
	list = append(list, GenInsertStmt(tab))
	list = append(list, GenUpdateStmt(tab))
	list = append(list, GenDeleteStmt(tab))
	list = append(list, GenInsertOrUpdateStmt(tab))
	list = append(list, genModelCountStmtDeclare(tab))
	for _, mtm := range tab.ManyToManys {
		list = append(list, GenInsertMidStmt(tab, mtm))
		list = append(list, GenManyToManyDeleteStmt(mtm, tab))
	}
	return strings.Join(list, "\n")
}

//GenInsertStmtInitBlock Generate insert sql.Stmt in init function
func GenInsertStmtInitBlock(tab Table, db Database) string {
	list := make([]string, 0, len(tab.Columns))
	for _, col := range tab.Columns {
		if !col.AutoIncrement && col.Default != "CURRENT_TIMESTAMP" && col.On != "UPDATE CURRENT_TIMESTAMP" {
			list = append(list, BackQuote(col.ColumnName))
		}
	}
	return fmt.Sprintf(InsertStmtInitTemp,
		tab.ModelName,
		db.ObjName,
		BackQuote(tab.TableName),
		strings.Join(list, ", "),
		strings.Trim(strings.Repeat("?, ", len(list)), ", "))
}

//GenUpdateStmtInitBlock Generate update sql.Stmt in init function
func GenUpdateStmtInitBlock(tab Table, db Database) string {
	list := make([]string, 0, len(tab.Columns))
	for _, col := range tab.Columns {
		if !col.AutoIncrement && col.Default != "CURRENT_TIMESTAMP" && col.On != "UPDATE CURRENT_TIMESTAMP" {
			list = append(list, BackQuote(col.ColumnName)+" = ?")
		}
	}
	return fmt.Sprintf(UpdateStmtInitTemp,
		tab.ModelName,
		db.ObjName,
		BackQuote(tab.TableName),
		strings.Join(list, ", "),
		BackQuote(tab.AutoIncrement.ColumnName))
}

//GenDeleteStmtInitBlock Generate delete sql.Stmt in init function
func GenDeleteStmtInitBlock(tab Table, db Database) string {
	return fmt.Sprintf(DeleteStmtInitTemp,
		tab.ModelName,
		db.ObjName,
		BackQuote(tab.TableName),
		BackQuote(tab.AutoIncrement.ColumnName))
}

//GenInsertOrUpdateStmtInitBlock Generate insert or update sql.Stmt in init function
func GenInsertOrUpdateStmtInitBlock(tab Table, db Database) string {
	insertList := make([]string, 0, len(tab.Columns))
	argList := make([]string, 0, len(tab.Columns))
	updateList := make([]string, 0, len(tab.Columns))
	for _, col := range tab.Columns {
		if !col.AutoIncrement && col.Default != "CURRENT_TIMESTAMP" && col.On != "UPDATE CURRENT_TIMESTAMP" {
			insertList = append(insertList, BackQuote(col.ColumnName))
			argList = append(argList, "?")
			updateList = append(updateList, fmt.Sprintf(UpdateColumnTemp, BackQuote(col.ColumnName)))
		}
	}
	return fmt.Sprintf(InsertOrUpdateStmtInitTemp,
		tab.ModelName,
		db.ObjName,
		BackQuote(tab.TableName),
		strings.Join(insertList, ", "),
		strings.Join(argList, ", "),
		fmt.Sprintf(UpdateLastInsertIDTemp, BackQuote(tab.AutoIncrement.ColumnName), BackQuote(tab.AutoIncrement.ColumnName)),
		strings.Join(updateList, ", "))
}

//GenInsertMidStmtInitBlock Generate insert into middle table sql.Stmt for many to many relation in init function
func GenInsertMidStmtInitBlock(mtm ManyToMany, tab Table, db Database) string {
	return fmt.Sprintf(InsertMidStmtInitTemp,
		tab.ModelName+"To"+mtm.DstTab.ModelName,
		db.ObjName,
		BackQuote(mtm.MidTab.TableName),
		BackQuote(mtm.MidLeftCol.ColumnName),
		BackQuote(mtm.MidRightCol.ColumnName))
}

//GenManyToManyDeleteStmtInitBlock Generate delete middle table sql.Stmt for many to many relation in init function
func GenManyToManyDeleteStmtInitBlock(mtm ManyToMany, tab Table, db Database) string {
	return fmt.Sprintf(ManyToManyDeleteStmtInitTemp,
		tab.ModelName,
		mtm.DstTab.ModelName,
		db.ObjName,
		BackQuote(mtm.MidTab.TableName),
		BackQuote(mtm.MidLeftCol.ColumnName),
		BackQuote(mtm.MidRightCol.ColumnName))
}

//GenStmtInitBlock Generate sql.Stmt init block in init function
func GenStmtInitBlock(tab Table, db Database) string {
	list := make([]string, 0, 8)
	list = append(list, GenInsertStmtInitBlock(tab, db))
	list = append(list, GenUpdateStmtInitBlock(tab, db))
	list = append(list, GenDeleteStmtInitBlock(tab, db))
	list = append(list, GenInsertOrUpdateStmtInitBlock(tab, db))
	list = append(list, genModelCountStmtInit(tab, db))
	for _, mtm := range tab.ManyToManys {
		list = append(list, GenInsertMidStmtInitBlock(mtm, tab, db))
		list = append(list, GenManyToManyDeleteStmtInitBlock(mtm, tab, db))
	}
	return strings.Join(list, "\n")
}

//GenStmtArgNilToDefault Generate block for replacing nil argument to default value in insert or update method when not null column value is nil
func GenStmtArgNilToDefault(col Column) string {
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
	return fmt.Sprintf(StmtArgNilToDefaultTemp, col.FieldName, value, col.FieldName)
}

//GenStmtArgNilToDefaultBlock Generate block to replacing nil arguments to default value
func GenStmtArgNilToDefaultBlock(tab Table) string {
	list := make([]string, 0, len(tab.Columns))
	for _, col := range tab.Columns {
		if !col.AutoIncrement && col.Default != "CURRENT_TIMESTAMP" && col.On != "UPDATE CURRENT_TIMESTAMP" {
			if !col.Nullable && col.Default != "" {
				list = append(list, GenStmtArgNilToDefault(col))
			} else {
				list = append(list, fmt.Sprintf(StmtArgTemp, col.FieldName))
			}
		}
	}
	return fmt.Sprintf(StmtArgNilToDefaultBlockTemp, len(list), strings.Join(list, "\n"))
}

//GenModelListType Generate XXXList struct definition
func GenModelListType(tab Table) string {
	return fmt.Sprintf(ModelListTypeTemp, tab.ModelName, tab.ModelName)
}

//GenModelCompareMethod Generate field compare method for sort function
func GenModelCompareMethod(col Column) string {
	var temp string
	switch col.FieldType {
	case "int64":
		temp = ModelCompareByIntMethodTemp
	case "float64":
		temp = ModelCompareByFloatMethodTemp
	case "string":
		temp = ModelCompareByStringMethodTemp
	case "bool":
		temp = ModelCompareByBoolMethodTemp
	case "time.Time":
		temp = ModelCompareByTimeMethodTemp
	}
	return fmt.Sprintf(temp, col.FieldName, col.FieldName, col.FieldName, col.FieldName, col.FieldName)
}

//GenModelSortMethodsStructType Generate model sort struct
func GenModelSortMethodsStructType(tab Table) string {
	list := make([]string, len(tab.Columns))
	for i, col := range tab.Columns {
		list[i] = fmt.Sprintf(ModelSortMethodsStructFieldTypeTemp, col.FieldName)
	}
	return fmt.Sprintf(ModelSortMethodsStructTypeTemp, tab.ModelName, strings.Join(list, "\n"))
}

//GenModelSortMethodsFunc Generate model sort method
func GenModelSortMethodsFunc(tab Table) string {
	list := make([]string, len(tab.Columns))
	for i, col := range tab.Columns {
		list[i] = GenModelCompareMethod(col)
	}
	return fmt.Sprintf(ModelSortMethodsStructFuncTemp, tab.ModelName, tab.ModelName, tab.ModelName, strings.Join(list, "\n"))
}

//GenModelListLenMethod Generate XXXList.Len() method
func GenModelListLenMethod(tab Table) string {
	return fmt.Sprintf(ModelListLenMethodTemp, tab.ModelName)
}

//GenModelListSwapMethod Generate XXXList.Swap() method
func GenModelListSwapMethod(tab Table) string {
	return fmt.Sprintf(ModelListSwapMethodTemp, tab.ModelName)
}

//GenModelListLessMethod Generate XXXList.Less() method
func GenModelListLessMethod(tab Table) string {
	return fmt.Sprintf(ModelListLessMethodTemp, tab.ModelName)
}

//GenModelSortFuncSwitchBlock Generate XXXSortBy() switch block
func GenModelSortFuncSwitchBlock(col Column, tab Table) string {
	return fmt.Sprintf(ModelSortFuncSwitchBlockTemp,
		col.FieldName,
		tab.ArgName,
		tab.ArgName,
		col.FieldName)
}

//GenModelSortFunc Generate XXXSortBy() function
func GenModelSortFunc(tab Table) string {
	list := make([]string, len(tab.Columns))
	for i, col := range tab.Columns {
		list[i] = GenModelSortFuncSwitchBlock(col, tab)
	}
	return fmt.Sprintf(ModelSortFuncTemp,
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

func genModelCountStmtInit(tab Table, db Database) string {
	return fmt.Sprintf(modelCountStmtInitTemp, tab.ModelName, db.ObjName, BackQuote(tab.TableName))
}

func genModelCountFunc(tab Table) string {
	return fmt.Sprintf(modelCountFuncTemp, tab.ModelName, tab.ModelName)
}

//Gen Generate database definition
func Gen(db Database, outName string) error {
	buf := bytes.NewBuffer([]byte{})
	buf.WriteString(GenPackage(db.Package) + "\n\n")
	buf.WriteString(ImportTemp + "\n\n")
	buf.WriteString(GenDb(db) + "\n\n")
	initStmtBuf := bytes.NewBuffer([]byte{})
	for _, tab := range db.Tables {
		buf.WriteString(GenStmtVar(tab) + "\n\n")
		initStmtBuf.WriteString(GenStmtInitBlock(tab, db) + "\n\n")
	}
	buf.WriteString(GenInitFunc(db, initStmtBuf.String()) + "\n\n")
	for _, tab := range db.Tables {
		buf.WriteString(GenQueryFieldMapBlock(tab) + "\n\n")
		buf.WriteString(GenModel(tab) + "\n\n")
		for _, col := range tab.Columns {
			buf.WriteString(GenModelGetFieldMethod(col, tab) + "\n\n")
			buf.WriteString(GenModelSetFieldMethod(col, tab) + "\n\n")
		}
		for _, fk := range tab.ForeignKeys {
			buf.WriteString(GenForeignKeyMethod(fk, tab, db) + "\n\n")
		}
		for _, rfk := range tab.ReverseForeignKeys {
			buf.WriteString(GenReverseForeignKeyStruct(rfk, tab, db) + "\n\n")
			buf.WriteString(GenReverseForeignKeyMethod(rfk, tab, db) + "\n\n")
		}
		for _, mtm := range tab.ManyToManys {
			buf.WriteString(GenManyToManyStruct(mtm, tab, db) + "\n\n")
			buf.WriteString(GenManyToManyMethod(mtm, tab, db) + "\n\n")
		}
		buf.WriteString(GenNewFunc(tab) + "\n\n")
		buf.WriteString(GenAllFunc(tab, db) + "\n\n")
		buf.WriteString(GenQueryFunc(tab, db) + "\n\n")
		buf.WriteString(GenQueryOneFunc(tab, db) + "\n\n")
		buf.WriteString(GenModelInsertMethod(tab) + "\n\n")
		buf.WriteString(GenModelInsertOrUpdateMethod(tab) + "\n\n")
		buf.WriteString(GenModelUpdateMethod(tab) + "\n\n")
		buf.WriteString(GenModelDeleteMethod(tab) + "\n\n")
		buf.WriteString(genModelCountFunc(tab) + "\n\n")
		buf.WriteString(GenFromRowsFunc(tab) + "\n\n")
		buf.WriteString(GenFromRowFunc(tab) + "\n\n")
		buf.WriteString(GenModelExistsMethod(tab, db) + "\n\n")
		buf.WriteString(GenModelCheckMethod(tab) + "\n\n")
		buf.WriteString(GenModelListType(tab) + "\n\n")
		buf.WriteString(GenModelListLenMethod(tab) + "\n\n")
		buf.WriteString(GenModelListSwapMethod(tab) + "\n\n")
		buf.WriteString(GenModelListLessMethod(tab) + "\n\n")
		buf.WriteString(GenModelSortMethodsStructType(tab) + "\n\n")
		buf.WriteString(GenModelSortMethodsFunc(tab) + "\n\n")
		buf.WriteString(GenModelSortFunc(tab) + "\n\n")
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
