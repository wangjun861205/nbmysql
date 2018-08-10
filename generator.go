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

func GenPackage(pkg string) string {
	return fmt.Sprintf(PackageTemp, pkg)
}

func GenDb(db Database) string {
	return fmt.Sprintf(DbTemp, db.ObjName)
}

func GenInitFunc(db Database, stmtInit string) string {
	return fmt.Sprintf(InitFuncTemp, db.Username, db.Password, db.Address, db.DatabaseName, db.ObjName, stmtInit)
}

func GenMapElemBlock(col Column) string {
	return fmt.Sprintf(MapElemTemp, "@"+col.FieldName, BackQuote(col.ColumnName))
}

func GenQueryFieldMapBlock(tab Table) string {
	elemList := make([]string, len(tab.Columns))
	for i, col := range tab.Columns {
		elemList[i] = GenMapElemBlock(col)
	}
	return fmt.Sprintf(QueryFieldMapTemp, tab.ModelName, strings.Join(elemList, "\n"))
}

func GenField(col Column) string {
	return fmt.Sprintf(FieldTemp, col.FieldName, col.FieldType)
}

func GenModel(tab Table) string {
	list := make([]string, len(tab.Columns))
	for i, col := range tab.Columns {
		list[i] = GenField(col)
	}
	return fmt.Sprintf(ModelTemp, tab.ModelName, strings.Join(list, "\n"))
}

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

func GenModelSetFieldMethod(col Column, tab Table) string {
	return fmt.Sprintf(ModelSetFieldMethodTemp,
		tab.ModelName,
		col.FieldName,
		col.FieldType,
		col.FieldName,
		col.FieldName)
}

func GenNewFuncArgList(tab Table) string {
	list := make([]string, 0, len(tab.Columns))
	for _, col := range tab.Columns {
		if !col.AutoIncrement && col.Default != "CURRENT_TIMESTAMP" && col.On != "UPDATE CURRENT_TIMESTAMP" {
			list = append(list, fmt.Sprintf(FuncArgTemp, col.ArgName, "*"+col.FieldType))
		}
	}
	return strings.Join(list, ", ")
}

func GenNewFuncAsignBlock(tab Table) string {
	list := make([]string, 0, len(tab.Columns))
	for _, col := range tab.Columns {
		if !col.AutoIncrement && col.Default != "CURRENT_TIMESTAMP" && col.On != "UPDATE CURRENT_TIMESTAMP" {
			list = append(list, fmt.Sprintf(NewModelAsignTemp, col.FieldName, col.ArgName))
		}
	}
	return strings.Join(list, ",\n")
}

func GenNewFunc(tab Table) string {
	return fmt.Sprintf(NewModelFuncTemp,
		tab.ModelName,
		GenNewFuncArgList(tab),
		tab.ModelName,
		tab.ModelName,
		GenNewFuncAsignBlock(tab))
}

func GenAllFunc(tab Table, db Database) string {
	return fmt.Sprintf(AllModelFuncTemp, tab.ModelName, tab.ModelName, db.ObjName, BackQuote(tab.TableName), tab.ModelName, tab.ModelName)
}

func GenQueryFunc(tab Table, db Database) string {
	return fmt.Sprintf(QueryModelFuncTemp, tab.ModelName, tab.ModelName, tab.ModelName, db.ObjName, BackQuote(tab.TableName), tab.ModelName, tab.ModelName)
}

func GenQueryOneFunc(tab Table, db Database) string {
	return fmt.Sprintf(QueryOneFuncTemp,
		tab.ModelName,
		tab.ModelName,
		tab.ModelName,
		db.ObjName,
		BackQuote(tab.TableName),
		tab.ModelName)
}

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

func GenReverseForeignKeyStruct(rfk ReverseForeignKey, srcTab Table, db Database) string {
	return fmt.Sprintf(ReverseForeignKeyStructTypeTemp,
		srcTab.ModelName,
		rfk.DstTab.ModelName,
		rfk.DstTab.ModelName,
		rfk.DstTab.ModelName)
}

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

func GenCheckStringBlock(col Column, localModelName string) string {
	return fmt.Sprintf(ModelCheckStringBlockTemp, localModelName, col.FieldName, BackQuote(col.ColumnName), localModelName, col.FieldName)
}

func GenCheckIntBlock(col Column, localModelName string) string {
	return fmt.Sprintf(ModelCheckIntBlockTemp, localModelName, col.FieldName, BackQuote(col.ColumnName), localModelName, col.FieldName)
}

func GenCheckFloatBlock(col Column, localModelName string) string {
	return fmt.Sprintf(ModelCheckFloatBlockTemp, localModelName, col.FieldName, BackQuote(col.ColumnName), localModelName, col.FieldName)
}

func GenCheckTimeBlock(col Column, localModelName string) string {
	return fmt.Sprintf(ModelCheckTimeBlockTemp, localModelName, col.FieldName, BackQuote(col.ColumnName), localModelName, col.FieldName)
}

func GenCheckBoolBlock(col Column, localModelName string) string {
	return fmt.Sprintf(ModelCheckBoolBlockTemp, localModelName, col.FieldName, BackQuote(col.ColumnName), localModelName, col.FieldName)
}

func GenInsertSQL(tab Table) string {
	return fmt.Sprintf(InsertSQLTemp, BackQuote(tab.TableName))
}

func GenInsertArgsBlock(tab Table) string {
	list := make([]string, 0, len(tab.Columns))
	for _, col := range tab.Columns {
		if !col.AutoIncrement && col.Default != "CURRENT_TIMESTAMP" && col.On != "UPDATE CURRENT_TIMESTAMP" {
			list = append(list, fmt.Sprintf(ModelInsertArgTemp, col.FieldName))
		}
	}
	return strings.Join(list, ", ")
}

func GenModelInsertMethod(tab Table) string {
	return fmt.Sprintf(ModelInsertMethodTemp,
		tab.ModelName,
		GenStmtArgNilToDefaultBlock(tab),
		tab.ModelName,
		tab.AutoIncrement.FieldName)
}

// func GenUpdateFunc(tab Table, db Database) string {
// 	list := make([]string, len(tab.Columns))
// 	for i, col := range tab.Columns {
// 		if !col.AutoIncrement && col.ColumnName != tab.PrimaryKey.ColumnName {
// 			switch col.FieldType {
// 			case "string":
// 				list[i] = GenCheckStringBlock(col, "m")
// 			case "int64":
// 				list[i] = GenCheckIntBlock(col, "m")
// 			case "float64":
// 				list[i] = GenCheckFloatBlock(col, "m")
// 			case "time.Time":
// 				list[i] = GenCheckTimeBlock(col, "m")
// 			case "bool":
// 				list[i] = GenCheckBoolBlock(col, "m")

// 			}
// 		}
// 	}
// 	return fmt.Sprintf(ModelUpdateMethodTemp, tab.ModelName, strings.Join(list, "\n"), db.ObjName, BackQuote(tab.TableName),
// 		BackQuote(tab.PrimaryKey.ColumnName), tab.PrimaryKey.FieldName)
// }

func GenUpdateArgs(tab Table) string {
	list := make([]string, 0, len(tab.Columns))
	for _, col := range tab.Columns {
		if !col.AutoIncrement && col.Default != "CURRENT_TIMESTAMP" && col.On != "UPDATE CURRENT_TIMESTAMP" {
			list = append(list, fmt.Sprintf(UpdateArgTemp, col.FieldName))
		}
	}
	return strings.Join(list, ", ")
}

func GenModelUpdateMethod(tab Table) string {
	return fmt.Sprintf(ModelUpdateMethodTemp,
		tab.ModelName,
		GenStmtArgNilToDefaultBlock(tab),
		tab.AutoIncrement.FieldName,
		tab.ModelName)
}

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

// func GenInsertOrUpdateMethod(tab Table) string {
// 	return fmt.Sprintf(ModelInsertOrUpdateMethodTemp, tab.ModelName)
// }

func GenInsertOrUpdateArgs(tab Table) string {
	list := make([]string, 0, len(tab.Columns))
	for _, col := range tab.Columns {
		if !col.AutoIncrement && col.Default != "CURRENT_TIMESTAMP" && col.On != "UPDATE CURRENT_TIMESTAMP" {
			list = append(list, fmt.Sprintf(InsertOrUpdateArgTemp, col.FieldName))
		}
	}
	finalList := make([]string, len(list)*2)
	copy(finalList[:len(list)], list)
	copy(finalList[len(list):], list)
	return strings.Join(finalList, ", ")
}
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

func GenModelDeleteMethod(tab Table) string {
	return fmt.Sprintf(ModelDeleteMethodTemp,
		tab.ModelName,
		tab.ModelName,
		fmt.Sprintf(DeleteArgTemp, tab.AutoIncrement.FieldName))
}

func GenNewMiddleTypeBlock(col Column) string {
	return fmt.Sprintf(NewMiddleTypeTemp, col.ArgName, col.MidType)
}

func GenFromRowsCheckBlock(col Column) string {
	return fmt.Sprintf(ModelFromRowsCheckNullBlockTemp, col.ArgName, col.ArgName, col.ArgName)
}

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

func GenModelCheckMethod(tab Table) string {
	list := make([]string, 0, len(tab.Columns))
	for _, col := range tab.Columns {
		if !col.Nullable && !col.AutoIncrement && col.Default == "" {
			list = append(list, fmt.Sprintf(FieldCheckNullTemp, col.FieldName, tab.ModelName, col.FieldName))
		}
	}
	return fmt.Sprintf(ModelCheckMethodTemp, tab.ModelName, strings.Join(list, "\n"))
}

func GenInsertStmt(tab Table) string {
	return fmt.Sprintf(InsertStmtTemp, tab.ModelName)
}

func GenUpdateStmt(tab Table) string {
	return fmt.Sprintf(UpdateStmtTemp, tab.ModelName)
}

func GenDeleteStmt(tab Table) string {
	return fmt.Sprintf(DeleteStmtTemp, tab.ModelName)
}

func GenInsertOrUpdateStmt(tab Table) string {
	return fmt.Sprintf(InsertOrUpdateStmtTemp, tab.ModelName)
}

func GenInsertMidStmt(tab Table, mtm ManyToMany) string {
	return fmt.Sprintf(InsertStmtTemp, tab.ModelName+"To"+mtm.DstTab.ModelName)
}

func GenManyToManyDeleteStmt(mtm ManyToMany, tab Table) string {
	return fmt.Sprintf(ManyToManyDeleteStmtTemp, tab.ModelName, mtm.DstTab.ModelName)
}

func GenStmtVar(tab Table) string {
	list := make([]string, 0, 8)
	list = append(list, GenInsertStmt(tab))
	list = append(list, GenUpdateStmt(tab))
	list = append(list, GenDeleteStmt(tab))
	list = append(list, GenInsertOrUpdateStmt(tab))
	for _, mtm := range tab.ManyToManys {
		list = append(list, GenInsertMidStmt(tab, mtm))
		list = append(list, GenManyToManyDeleteStmt(mtm, tab))
	}
	return strings.Join(list, "\n")
}

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

func GenDeleteStmtInitBlock(tab Table, db Database) string {
	return fmt.Sprintf(DeleteStmtInitTemp,
		tab.ModelName,
		db.ObjName,
		BackQuote(tab.TableName),
		BackQuote(tab.AutoIncrement.ColumnName))
}

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

func GenInsertMidStmtInitBlock(mtm ManyToMany, tab Table, db Database) string {
	return fmt.Sprintf(InsertMidStmtInitTemp,
		tab.ModelName+"To"+mtm.DstTab.ModelName,
		db.ObjName,
		BackQuote(mtm.MidTab.TableName),
		BackQuote(mtm.MidLeftCol.ColumnName),
		BackQuote(mtm.MidRightCol.ColumnName))
}

func GenManyToManyDeleteStmtInitBlock(mtm ManyToMany, tab Table, db Database) string {
	return fmt.Sprintf(ManyToManyDeleteStmtInitTemp,
		tab.ModelName,
		mtm.DstTab.ModelName,
		db.ObjName,
		BackQuote(mtm.MidTab.TableName),
		BackQuote(mtm.MidLeftCol.ColumnName),
		BackQuote(mtm.MidRightCol.ColumnName))
}

func GenStmtInitBlock(tab Table, db Database) string {
	list := make([]string, 0, 8)
	list = append(list, GenInsertStmtInitBlock(tab, db))
	list = append(list, GenUpdateStmtInitBlock(tab, db))
	list = append(list, GenDeleteStmtInitBlock(tab, db))
	list = append(list, GenInsertOrUpdateStmtInitBlock(tab, db))
	for _, mtm := range tab.ManyToManys {
		list = append(list, GenInsertMidStmtInitBlock(mtm, tab, db))
		list = append(list, GenManyToManyDeleteStmtInitBlock(mtm, tab, db))
	}
	return strings.Join(list, "\n")
}

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

func GenStmtArgNilToDefaultBlock(tab Table) string {
	list := make([]string, 0, len(tab.Columns))
	for _, col := range tab.Columns {
		if !col.AutoIncrement && col.Default != "CURRENT_TIMESTAMP" && col.On != "UPDATE CURRENT_TIMESTAMP" {
			if !col.Nullable && col.Default != "" {
				// list = append(list, fmt.Sprintf(StmtArgNilToDefaultTemp, col.FieldName))
				list = append(list, GenStmtArgNilToDefault(col))
			} else {
				list = append(list, fmt.Sprintf(StmtArgTemp, col.FieldName))
			}
		}
	}
	return fmt.Sprintf(StmtArgNilToDefaultBlockTemp, len(list), strings.Join(list, "\n"))
}

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
		buf.WriteString(GenFromRowsFunc(tab) + "\n\n")
		buf.WriteString(GenFromRowFunc(tab) + "\n\n")
		buf.WriteString(GenModelExistsMethod(tab, db) + "\n\n")
		buf.WriteString(GenModelCheckMethod(tab) + "\n\n")
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
