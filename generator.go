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

func genPackage(db Database) (string, error) {
	return nbfmt.Fmt(packageTemp, map[string]interface{}{"DB": db})
}

func genDB(db Database) (string, error) {
	return nbfmt.Fmt(dbTemp, map[string]interface{}{"DB": db})
}

func genInitFunc(db Database) (string, error) {
	return nbfmt.Fmt(initFuncTemp, map[string]interface{}{"DB": db})
}

func genWhereMap(db Database) (string, error) {
	return nbfmt.Fmt(whereMapTemp, map[string]interface{}{"DB": db})
}

func genModel(db Database) (string, error) {
	return nbfmt.Fmt(modelTypeTemp, map[string]interface{}{"DB": db})
}

func genNewModelFunc(db Database) (string, error) {
	return nbfmt.Fmt(newModelFuncTemp, map[string]interface{}{"DB": db})
}

func genAllFunc(db Database) (string, error) {
	return nbfmt.Fmt(allModelFuncTemp, map[string]interface{}{"DB": db})
}

func genQueryFunc(db Database) (string, error) {
	return nbfmt.Fmt(queryFuncTemp, map[string]interface{}{"DB": db})
}

func genQueryOneFunc(db Database) (string, error) {
	return nbfmt.Fmt(queryOneFuncTemp, map[string]interface{}{"DB": db})
}

func genForeignKeyMethod(db Database) (string, error) {
	return nbfmt.Fmt(foreignKeyMethodTemp, map[string]interface{}{"DB": db})
}

func genReverseForeignKeyMethod(db Database) (string, error) {
	return nbfmt.Fmt(reverseForeignKeyMethodTemp, map[string]interface{}{"DB": db})
}

func genManyToManyMethod(db Database) (string, error) {
	return nbfmt.Fmt(manyToManyMethodTemp, map[string]interface{}{"DB": db})
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

func genModelInsertMethod(db Database) (string, error) {
	return nbfmt.Fmt(modelInsertMethodTemp, map[string]interface{}{"DB": db})
}

func genModelUpdateMethod(db Database) (string, error) {
	return nbfmt.Fmt(modelUpdateMethodTemp, map[string]interface{}{"DB": db})
}

func genModelInsertOrUpdateMethod(db Database) (string, error) {
	return nbfmt.Fmt(modelInsertOrUpdateMethodTemp, map[string]interface{}{"DB": db})
}

//genModelDeleteMethod generate XXX.Delete() method
func genModelDeleteMethod(db Database) (string, error) {
	return nbfmt.Fmt(modelDeleteMethodTemp, map[string]interface{}{"DB": db})
}

//genNewMiddleTypeBlock generate new middle type consponsed to field type block
// func genNewMiddleTypeBlock(col Column) string {
// 	s, err := nbfmt.Fmt(newMiddleTypeTemp, map[string]interface{}{"Column": col})
// 	if err != nil {
// 		fmt.Println("error: genNewMiddleTypeBlock() failed")
// 		log.Fatal(err)
// 	}
// 	return s
// }

func genFromRowsFunc(db Database) (string, error) {
	return nbfmt.Fmt(fromRowsFuncTemp, map[string]interface{}{"DB": db})
}

func genFromRowFunc(db Database) (string, error) {
	return nbfmt.Fmt(fromRowFuncTemp, map[string]interface{}{"DB": db})
}

func genModelCheckMethod(db Database) (string, error) {
	return nbfmt.Fmt(modelCheckMethodTemp, map[string]interface{}{"DB": db})
}

func genModelListType(db Database) (string, error) {
	return nbfmt.Fmt(modelListTypeTemp, map[string]interface{}{"DB": db})
}

func genCountFunc(db Database) (string, error) {
	return nbfmt.Fmt(countFuncTemp, map[string]interface{}{"DB": db})
}

func genFieldType(db Database) (string, error) {
	return nbfmt.Fmt(fieldTypeTemp, map[string]interface{}{"DB": db})
}

func genModelInsertFunc(db Database) (string, error) {
	return nbfmt.Fmt(modelInsertFuncTemp, map[string]interface{}{"DB": db})
}

func genUpdateFunc(db Database) (string, error) {
	return nbfmt.Fmt(updateFuncTemp, map[string]interface{}{"DB": db})
}

func genModelInvalidateMethod(db Database) (string, error) {
	return nbfmt.Fmt(modelInvalidateMethodTemp, map[string]interface{}{"DB": db})
}

func genDeleteFunc(db Database) (string, error) {
	return nbfmt.Fmt(deleteFuncTemp, map[string]interface{}{"DB": db})
}

//Gen generate database definition
func Gen(db Database, outName string) error {
	buf := bytes.NewBuffer([]byte{})
	s, err := genPackage(db)
	if err != nil {
		return err
	}
	buf.WriteString(s)
	buf.WriteString(importTemp + "\n\n")
	s, err = genDB(db)
	if err != nil {
		return err
	}
	buf.WriteString(s)
	s, err = genInitFunc(db)
	if err != nil {
		return err
	}
	buf.WriteString(s)
	s, err = genWhereMap(db)
	if err != nil {
		return err
	}
	buf.WriteString(s)
	s, err = genFieldType(db)
	if err != nil {
		return err
	}
	buf.WriteString(s)
	s, err = genModelInsertFunc(db)
	if err != nil {
		return err
	}
	buf.WriteString(s)
	s, err = genModel(db)
	if err != nil {
		return err
	}
	buf.WriteString(s)
	s, err = genNewModelFunc(db)
	if err != nil {
		return err
	}
	buf.WriteString(s)
	s, err = genModelInsertMethod(db)
	if err != nil {
		return err
	}
	buf.WriteString(s)
	s, err = genModelUpdateMethod(db)
	if err != nil {
		return err
	}
	buf.WriteString(s)
	s, err = genUpdateFunc(db)
	if err != nil {
		return err
	}
	buf.WriteString(s)
	s, err = genModelInvalidateMethod(db)
	if err != nil {
		return err
	}
	buf.WriteString(s)
	s, err = genModelDeleteMethod(db)
	if err != nil {
		return err
	}
	buf.WriteString(s)
	s, err = genDeleteFunc(db)
	if err != nil {
		return err
	}
	buf.WriteString(s)
	s, err = genFromRowsFunc(db)
	if err != nil {
		return err
	}
	buf.WriteString(s)
	s, err = genFromRowFunc(db)
	if err != nil {
		return err
	}
	buf.WriteString(s)
	s, err = genAllFunc(db)
	if err != nil {
		return err
	}
	buf.WriteString(s)
	s, err = genQueryFunc(db)
	if err != nil {
		return err
	}
	buf.WriteString(s)
	s, err = genQueryOneFunc(db)
	if err != nil {
		return err
	}
	buf.WriteString(s)
	s, err = genCountFunc(db)
	if err != nil {
		return err
	}
	buf.WriteString(s)
	s, err = genModelInsertOrUpdateMethod(db)
	if err != nil {
		return err
	}
	buf.WriteString(s)
	s, err = genModelCheckMethod(db)
	if err != nil {
		return err
	}
	buf.WriteString(s)
	s, err = genModelListType(db)
	if err != nil {
		return err
	}
	buf.WriteString(s)
	s, err = genForeignKeyMethod(db)
	if err != nil {
		return err
	}
	buf.WriteString(s)
	s, err = genReverseForeignKeyMethod(db)
	if err != nil {
		return err
	}
	buf.WriteString(s)
	s, err = genManyToManyMethod(db)
	if err != nil {
		return err
	}
	buf.WriteString(s)
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
