package nbmysql

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func GenPackage(pkg string) string {
	return fmt.Sprintf(PackageTemp, pkg)
}

func GenDb(db Database) string {
	return fmt.Sprintf(DbTemp, db.ObjName)
}

func GenInitFunc(db Database) string {
	return fmt.Sprintf(InitFuncTemp, db.Username, db.Password, db.Address, db.DatabaseName, db.ObjName)
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

// func GenNewFuncArg(col Column, tab Table) string {
// 	var argType string
// 	switch col.FieldType {
// 	case "string":
// 		argType = "nbmysql.String"
// 	case "int64":
// 		argType = "nbmysql.Int"
// 	case "float64":
// 		argType = "nbmysql.Float"
// 	case "bool":
// 		argType = "nbmysql.Bool"
// 	case "time.Time":
// 		argType = "nbmysql.Time"
// 	}
// 	return fmt.Sprintf(FuncArgTemp, tab.ArgName, col.FieldName, argType)
// }

func GenNewFuncArgList(tab Table) string {
	list := make([]string, len(tab.Columns))
	for i, col := range tab.Columns {
		var argType string
		switch col.FieldType {
		case "string":
			argType = "[]byte"
		case "int64":
			argType = "complex128"
		case "float64":
			argType = "complex128"
		case "bool":
			argType = "int"
		case "time.Time":
			argType = "*time.Time"
		}
		list[i] = fmt.Sprintf(FuncArgTemp, col.ArgName, argType)
	}
	return strings.Join(list, ", ")
}

// func GenNewFuncArgName(col Column) string {
// 	return fmt.Sprintf(FuncArgNameTemp, "_"+col.ArgName)
// }

func GenNewFuncArgNameList(tab Table) string {
	list := make([]string, len(tab.Columns))
	for i, col := range tab.Columns {
		list[i] = "_" + col.ArgName
	}
	return strings.Join(list, ", ")
}

func GenNewFuncVarBlock(tab Table) string {
	list := make([]string, len(tab.Columns))
	for i, col := range tab.Columns {
		list[i] = fmt.Sprintf(NewModelFuncVarTemp, "_"+col.ArgName, col.FieldType)
	}
	return strings.Join(list, "\n")
}

func GenNewFuncArgCheckBlock(tab Table) string {
	list := make([]string, len(tab.Columns))
	for i, col := range tab.Columns {
		switch col.FieldType {
		case "int64":
			list[i] = fmt.Sprintf(CheckIntArgTemp, col.ArgName, "_"+col.ArgName, col.ArgName, "_"+col.ArgName)
		case "float64":
			list[i] = fmt.Sprintf(CheckFloatArgTemp, col.ArgName, "_"+col.ArgName, col.ArgName, "_"+col.ArgName)
		case "string":
			list[i] = fmt.Sprintf(CheckStringArgTemp, col.ArgName, "_"+col.ArgName, col.ArgName, "_"+col.ArgName)
		case "bool":
			list[i] = fmt.Sprintf(CheckBoolArgTemp, col.ArgName, "_"+col.ArgName, col.ArgName, "_"+col.ArgName, "_"+col.ArgName)
		case "time.Time":
			list[i] = fmt.Sprintf(CheckTimeArgTemp, "_"+col.ArgName, col.ArgName)
		}
	}
	return strings.Join(list, "\n")
}

// func GenMiddleTypeToGo(tab Table) string {
// 	list := make([]string, len(tab.Columns))
// 	for i, col := range tab.Columns {
// 		list[i] = fmt.Sprintf(MiddleTypeToGoTemp, "_"+col.ArgName, tab.ArgName+col.FieldName)
// 	}
// 	return strings.Join(list, "\n")
// }

// func GenNewFunc(tab Table) string {
// 	argList := make([]string, len(tab.Columns))
// 	argNameList := make([]string, len(tab.Columns))
// 	checkList := make([]string, len(tab.Columns))
// 	for i, col := range tab.Columns {
// 		argList[i] = GenNewFuncArg(col)
// 		argNameList[i] = GenNewFuncArgName(col)
// 	}
// 	return fmt.Sprintf(NewModelFuncTemp, tab.ModelName, strings.Join(argList, ", "), tab.ModelName, GenMiddleTypeToGo(tab), tab.ArgName, tab.ModelName,
// 		strings.Join(argNameList, ", "), tab.ArgName)
// }

func GenNewFunc(tab Table) string {
	return fmt.Sprintf(NewModelFuncTemp,
		tab.ModelName,
		GenNewFuncArgList(tab),
		tab.ModelName,
		GenNewFuncVarBlock(tab),
		GenNewFuncArgCheckBlock(tab),
		"__"+tab.ArgName,
		tab.ModelName,
		GenNewFuncArgNameList(tab),
		"__"+tab.ArgName)
}

func GenAllFunc(tab Table, db Database) string {
	return fmt.Sprintf(AllModelFuncTemp, tab.ModelName, tab.ModelName, db.ObjName, BackQuote(tab.TableName), tab.ModelName, tab.ModelName)
}

func GenQueryFunc(tab Table, db Database) string {
	return fmt.Sprintf(QueryModelFuncTemp, tab.ModelName, tab.ModelName, tab.ModelName, db.ObjName, BackQuote(tab.TableName), tab.ModelName, tab.ModelName)
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

func GenManyToManyAddMethod(mtm ManyToMany, srcTab Table, db Database) string {
	return fmt.Sprintf(ManyToManyAddMethodTemp,
		mtm.DstTab.ArgName,
		mtm.DstTab.ModelName,
		mtm.DstTab.ArgName,
		mtm.DstTab.ModelName,
		db.ObjName,
		fmt.Sprintf(ManyToManyAddSQLTemp, BackQuote(mtm.MidTab.TableName), BackQuote(mtm.MidLeftCol.ColumnName), BackQuote(mtm.MidRightCol.ColumnName)),
		mtm.SrcCol.FieldName,
		mtm.DstTab.ArgName,
		mtm.DstCol.FieldName)
}

func GenManyToManyRemoveMethod(mtm ManyToMany, srcTab Table, db Database) string {
	return fmt.Sprintf(ManyToManyRemoveMethodTemp,
		mtm.DstTab.ArgName,
		mtm.DstTab.ModelName,
		mtm.DstTab.ArgName,
		mtm.DstTab.ModelName,
		db.ObjName,
		fmt.Sprintf(ManyToManyRemoveSQLTemp, BackQuote(mtm.MidTab.TableName), BackQuote(mtm.MidLeftCol.ColumnName), BackQuote(mtm.MidRightCol.ColumnName)),
		mtm.SrcCol.FieldName,
		mtm.DstTab.ArgName,
		mtm.DstCol.FieldName)
}

func GenManyToManyMethod(mtm ManyToMany, srcTab Table, db Database) string {
	allMethod := GenManyToManyAllMethod(mtm, srcTab, db)
	queryMethod := GenManyToManyQueryMethod(mtm, srcTab, db)
	addMethod := GenManyToManyAddMethod(mtm, srcTab, db)
	removeMethod := GenManyToManyRemoveMethod(mtm, srcTab, db)
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

func GenInsertFunc(tab Table, db Database) string {
	list := make([]string, len(tab.Columns))
	for i, col := range tab.Columns {
		switch col.FieldType {
		case "string":
			list[i] = GenCheckStringBlock(col, "m")
		case "int64":
			list[i] = GenCheckIntBlock(col, "m")
		case "float64":
			list[i] = GenCheckFloatBlock(col, "m")
		case "time.Time":
			list[i] = GenCheckTimeBlock(col, "m")
		case "bool":
			list[i] = GenCheckBoolBlock(col, "m")
		}
	}
	return fmt.Sprintf(ModelInsertMethodTemp, tab.ModelName, strings.Join(list, "\n"), db.ObjName, GenInsertSQL(tab), tab.AutoIncrement.FieldName)
}

func GenUpdateFunc(tab Table, db Database) string {
	list := make([]string, len(tab.Columns))
	for i, col := range tab.Columns {
		if !col.AutoIncrement && col.ColumnName != tab.PrimaryKey.ColumnName {
			switch col.FieldType {
			case "string":
				list[i] = GenCheckStringBlock(col, "m")
			case "int64":
				list[i] = GenCheckIntBlock(col, "m")
			case "float64":
				list[i] = GenCheckFloatBlock(col, "m")
			case "time.Time":
				list[i] = GenCheckTimeBlock(col, "m")
			case "bool":
				list[i] = GenCheckBoolBlock(col, "m")

			}
		}
	}
	return fmt.Sprintf(ModelUpdateMethodTemp, tab.ModelName, strings.Join(list, "\n"), db.ObjName, BackQuote(tab.TableName),
		BackQuote(tab.PrimaryKey.ColumnName), tab.PrimaryKey.FieldName)
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

func GenInsertOrUpdateMethod(tab Table) string {
	return fmt.Sprintf(ModelInsertOrUpdateMethodTemp, tab.ModelName)
}

func GenDeleteSQL(tab Table) string {
	return fmt.Sprintf(DeleteSQLTemp, BackQuote(tab.TableName), BackQuote(tab.PrimaryKey.ColumnName))
}

func GenManyToManyDeleteSQL(info ManyToMany) string {
	return fmt.Sprintf(DeleteSQLTemp, BackQuote(info.MidTab.TableName), BackQuote(info.MidLeftCol.ColumnName))
}

func GenManyToManyDeleteBlock(info ManyToMany) string {
	return fmt.Sprintf(ManyToManyDeleteBlockTemp, GenManyToManyDeleteSQL(info), info.SrcCol.FieldName)
}

func GenCasecadeDeleteLoop(rfk ReverseForeignKey, db Database) string {
	return fmt.Sprintf(CascadeDeleteLoopTemp,
		fmt.Sprintf(CascadeDeleteSQLTemp, BackQuote(rfk.DstTab.TableName), BackQuote(rfk.DstCol.ColumnName)),
		rfk.SrcCol.FieldName)
}

func GenDeleteFunc(tab Table, db Database) string {
	list := make([]string, len(tab.ManyToManys))
	for i, mtm := range tab.ManyToManys {
		list[i] = GenManyToManyDeleteBlock(mtm)
	}
	rfkList := make([]string, len(tab.ReverseForeignKeys))
	for i, rfk := range tab.ReverseForeignKeys {
		rfkList[i] = GenCasecadeDeleteLoop(rfk, db)
	}
	var cascadeDelete string
	if len(rfkList) > 0 {
		cascadeDelete = fmt.Sprintf(CascadeDeleteBlockTemp, strings.Join(rfkList, "\n"))
	}
	return fmt.Sprintf(ModelDeleteMethodTemp, tab.ModelName, db.ObjName, strings.Join(list, "\n"), cascadeDelete, GenDeleteSQL(tab), tab.PrimaryKey.FieldName)
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
		if !col.Nullable && !col.AutoIncrement {
			list = append(list, fmt.Sprintf(FieldCheckNullTemp, col.FieldName, tab.ModelName, col.FieldName))
		}
	}
	return fmt.Sprintf(ModelCheckMethodTemp, tab.ModelName, strings.Join(list, "\n"))
}

func Gen(db Database, outName string) error {
	buf := bytes.NewBuffer([]byte{})
	buf.WriteString(GenPackage(db.Package) + "\n\n")
	buf.WriteString(ImportTemp + "\n\n")
	buf.WriteString(GenDb(db) + "\n\n")
	buf.WriteString(GenInitFunc(db) + "\n\n")
	for _, tab := range db.Tables {
		buf.WriteString(GenQueryFieldMapBlock(tab) + "\n\n")
		buf.WriteString(GenModel(tab) + "\n\n")
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
		buf.WriteString(GenInsertFunc(tab, db) + "\n\n")
		buf.WriteString(GenInsertOrUpdateMethod(tab) + "\n\n")
		buf.WriteString(GenUpdateFunc(tab, db) + "\n\n")
		buf.WriteString(GenDeleteFunc(tab, db) + "\n\n")
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
