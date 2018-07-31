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
	return fmt.Sprintf(MapElemTemp, "@"+col.FieldName, col.ColumnName)
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

func GenModelRelation(srcTab, dstTab Table) string {
	return fmt.Sprintf(ModelRelationTemp, srcTab.ModelName, dstTab.ModelName, dstTab.ModelName, dstTab.ModelName)
}

func GenNewFuncArg(col Column) string {
	return fmt.Sprintf(FuncArgTemp, col.ArgName, col.FieldType)
}

func GenNewFunc(tab Table) string {
	argList := make([]string, len(tab.Columns))
	argNameList := make([]string, len(tab.Columns))
	for i, col := range tab.Columns {
		argList[i] = GenNewFuncArg(col)
		argNameList[i] = col.ArgName
	}
	return fmt.Sprintf(NewModelFuncTemp, tab.ModelName, strings.Join(argList, ", "), tab.ModelName, tab.ArgName, tab.ModelName,
		strings.Join(argNameList, ", "), tab.ArgName)
}

func GenAllFunc(tab Table, db Database) string {
	return fmt.Sprintf(AllModelFuncTemp, tab.ModelName, tab.ModelName, db.ObjName, BackQuote(tab.TableName), tab.ModelName, tab.ModelName)
}

func GenQueryFunc(tab Table, db Database) string {
	return fmt.Sprintf(QueryModelFuncTemp, tab.ModelName, tab.ModelName, tab.ModelName, db.ObjName, BackQuote(tab.TableName), tab.ModelName, tab.ModelName)
}

func GenManyToManyFunc(srcTab, dstTab, midTab Table, srcCol, dstCol, midLeftCol, midRightCol Column, db Database) string {
	allSql := fmt.Sprintf(ManyToManyAllSQLTemp,
		BackQuote(dstTab.TableName),
		BackQuote(srcTab.TableName),
		BackQuote(midTab.TableName),
		BackQuote(srcTab.TableName),
		BackQuote(srcCol.ColumnName),
		BackQuote(midTab.TableName),
		BackQuote(midLeftCol.ColumnName),
		BackQuote(dstTab.TableName),
		BackQuote(midTab.TableName),
		BackQuote(midRightCol.ColumnName),
		BackQuote(dstTab.TableName),
		BackQuote(dstCol.ColumnName),
		BackQuote(srcTab.TableName),
		BackQuote(srcCol.ColumnName))
	filterSql := allSql + " AND ?"
	return fmt.Sprintf(ModelRelationFuncTemp, srcTab.ModelName, dstTab.ModelName, srcCol.FieldName, srcTab.ModelName, dstTab.ModelName, srcTab.ModelName,
		dstTab.ModelName, dstTab.ModelName, db.ObjName, allSql, srcCol.FieldName, dstTab.ModelName, dstTab.ModelName, dstTab.ModelName,
		dstTab.ModelName, db.ObjName, filterSql, srcCol.FieldName, dstTab.ModelName, dstTab.ModelName)
}

func GenForeignKeyFunc(srcTab, dstTab Table, srcCol, dstCol Column, db Database) string {
	allSql := fmt.Sprintf(ForeignKeyAllSQLTemp, BackQuote(dstTab.TableName), BackQuote(srcTab.TableName), BackQuote(dstTab.TableName),
		BackQuote(srcTab.TableName), BackQuote(srcCol.ColumnName), BackQuote(dstTab.TableName), BackQuote(dstCol.ColumnName), BackQuote(srcTab.TableName),
		BackQuote(srcCol.ColumnName))
	filterSql := allSql + " AND ?"
	return fmt.Sprintf(ModelRelationFuncTemp, srcTab.ModelName, dstTab.ModelName, srcCol.FieldName, srcTab.ModelName, dstTab.ModelName, srcTab.ModelName,
		dstTab.ModelName, dstTab.ModelName, db.ObjName, allSql, srcCol.FieldName, dstTab.ModelName, dstTab.ModelName, dstTab.ModelName,
		dstTab.ModelName, db.ObjName, filterSql, srcCol.FieldName, dstTab.ModelName, dstTab.ModelName)
}

func GenCheckStringBlock(col Column) string {
	return fmt.Sprintf(ModelCheckStringBlockTemp, col.FieldName, col.ColumnName, col.FieldName)
}

func GenCheckIntBlock(col Column) string {
	return fmt.Sprintf(ModelCheckIntBlockTemp, col.FieldName, col.ColumnName, col.FieldName)
}

func GenCheckFloatBlock(col Column) string {
	return fmt.Sprintf(ModelCheckFloatBlockTemp, col.FieldName, col.ColumnName, col.FieldName)
}

func GenCheckTimeBlock(col Column) string {
	return fmt.Sprintf(ModelCheckTimeBlockTemp, col.FieldName, col.ColumnName, col.FieldName)
}

func GenCheckBoolBlock(col Column) string {
	return fmt.Sprintf(ModelCheckBoolBlockTemp, col.FieldName, col.ColumnName, col.FieldName)
}

func GenInsertFunc(tab Table, db Database) string {
	list := make([]string, len(tab.Columns))
	for i, col := range tab.Columns {
		switch col.FieldType {
		case "string":
			list[i] = GenCheckStringBlock(col)
		case "int64":
			list[i] = GenCheckIntBlock(col)
		case "float64":
			list[i] = GenCheckFloatBlock(col)
		case "time.Time":
			list[i] = GenCheckTimeBlock(col)
		}
	}
	return fmt.Sprintf(ModelInsertMethodTemp, tab.ModelName, strings.Join(list, "\n"), db.ObjName, BackQuote(tab.TableName), db.ObjName, tab.PrimaryKey.FieldName)
}

func GenUpdateFunc(tab Table, db Database) string {
	list := make([]string, len(tab.Columns))
	for i, col := range tab.Columns {
		switch col.FieldType {
		case "string":
			list[i] = GenCheckStringBlock(col)
		case "int64":
			list[i] = GenCheckIntBlock(col)
		case "float64":
			list[i] = GenCheckFloatBlock(col)
		case "time.Time":
			list[i] = GenCheckTimeBlock(col)

		}
	}
	return fmt.Sprintf(ModelUpdateMethodTemp, tab.ModelName, strings.Join(list, "\n"), db.ObjName, BackQuote(tab.TableName), tab.PrimaryKey.FieldName)
}

func GenDeleteFunc(tab Table, db Database) string {
	return fmt.Sprintf(ModelDeleteMethodTemp, tab.ModelName, db.ObjName, BackQuote(tab.TableName), tab.PrimaryKey.FieldName)
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
	varList := make([]string, len(tab.Columns))
	nameList := make([]string, len(tab.Columns))
	checkList := make([]string, len(tab.Columns))
	for i, col := range tab.Columns {
		midList[i] = GenNewMiddleTypeBlock(col)
		midNameList[i] = "_" + col.ArgName
		varList[i] = fmt.Sprintf("%s *%s", col.ArgName, col.FieldType)
		nameList[i] = col.ArgName
		checkList[i] = GenFromRowsCheckBlock(col)
	}
	return fmt.Sprintf(ModelFromRowsFuncTemp, tab.ModelName, tab.ModelName, strings.Join(midList, "\n"), strings.Join(midNameList, ", "),
		strings.Join(varList, "\n"), strings.Join(checkList, "\n"), tab.ModelName, strings.Join(nameList, ", "))
}

func Gen(db Database, outName string) error {
	buf := bytes.NewBuffer([]byte{})
	buf.WriteString(GenPackage(db.Package) + "\n")
	buf.WriteString(ImportTemp + "\n")
	buf.WriteString(GenDb(db) + "\n")
	buf.WriteString(GenInitFunc(db) + "\n")
	for _, tab := range db.Tables {
		buf.WriteString(GenQueryFieldMapBlock(tab) + "\n")
		buf.WriteString(GenModel(tab) + "\n")
		for _, fk := range tab.ForeignKeys {
			buf.WriteString(GenModelRelation(tab, *fk.DstTab) + "\n")
			buf.WriteString(GenForeignKeyFunc(tab, *fk.DstTab, *fk.SrcCol, *fk.DstCol, db) + "\n")
		}
		for _, mtm := range tab.ManyToManys {
			buf.WriteString(GenModelRelation(tab, *mtm.DstTab) + "\n")
			buf.WriteString(GenManyToManyFunc(tab, *mtm.DstTab, *mtm.MidTab, *mtm.SrcCol, *mtm.DstCol, *mtm.MidLeftCol, *mtm.MidRightCol, db) + "\n")
		}
		buf.WriteString(GenNewFunc(tab) + "\n")
		buf.WriteString(GenAllFunc(tab, db) + "\n")
		buf.WriteString(GenQueryFunc(tab, db) + "\n")
		buf.WriteString(GenInsertFunc(tab, db) + "\n")
		buf.WriteString(GenUpdateFunc(tab, db) + "\n")
		buf.WriteString(GenDeleteFunc(tab, db) + "\n")
		buf.WriteString(GenFromRowsFunc(tab) + "\n")
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
