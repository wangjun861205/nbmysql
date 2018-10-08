package nbmysql

import (
	"bytes"
	"os"
	"os/exec"

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

// func genWhereMap(db Database) (string, error) {
// 	return nbfmt.Fmt(whereMapTemp, map[string]interface{}{"DB": db})
// }

func genModelInstanceType(db Database) (string, error) {
	return nbfmt.Fmt(modelInstanceTypeTemp, map[string]interface{}{"DB": db})
}

// func genNewModelFunc(db Database) (string, error) {
// 	return nbfmt.Fmt(newModelFuncTemp, map[string]interface{}{"DB": db})
// }

// func genAllFunc(db Database) (string, error) {
// 	return nbfmt.Fmt(allModelFuncTemp, map[string]interface{}{"DB": db})
// }

// func genQueryFunc(db Database) (string, error) {
// 	return nbfmt.Fmt(queryFuncTemp, map[string]interface{}{"DB": db})
// }

// func genQueryOneFunc(db Database) (string, error) {
// 	return nbfmt.Fmt(queryOneFuncTemp, map[string]interface{}{"DB": db})
// }

func genGetFieldMethod(db Database) (string, error) {
	return nbfmt.Fmt(getFieldMethodTemp, map[string]interface{}{"DB": db})
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

// func genFieldType(db Database) (string, error) {
// 	return nbfmt.Fmt(fieldTypeTemp, map[string]interface{}{"DB": db})
// }

// func genModelInsertFunc(db Database) (string, error) {
// 	return nbfmt.Fmt(modelInsertFuncTemp, map[string]interface{}{"DB": db})
// }

// func genUpdateFunc(db Database) (string, error) {
// 	return nbfmt.Fmt(updateFuncTemp, map[string]interface{}{"DB": db})
// }

func genModelInvalidateMethod(db Database) (string, error) {
	return nbfmt.Fmt(modelInvalidateMethodTemp, map[string]interface{}{"DB": db})
}

// func genDeleteFunc(db Database) (string, error) {
// 	return nbfmt.Fmt(deleteFuncTemp, map[string]interface{}{"DB": db})
// }

func genSetLastInsertIDMethod(db Database) (string, error) {
	return nbfmt.Fmt(setLastInsertIDMethodTemp, map[string]interface{}{"DB": db})
}

// func genStmtType(db Database) (string, error) {
// 	return nbfmt.Fmt(stmtTypeTemp, map[string]interface{}{"DB": db})
// }

func genInsertStmtMethod(db Database) (string, error) {
	return nbfmt.Fmt(insertStmtMethodTemp, map[string]interface{}{"DB": db})
}

// func genInsertStmtFunc(db Database) (string, error) {
// 	return nbfmt.Fmt(insertStmtFuncTemp, map[string]interface{}{"DB": db})
// }

func genUpdateStmtMethod(db Database) (string, error) {
	return nbfmt.Fmt(updateStmtMethodTemp, map[string]interface{}{"DB": db})
}

// func genUpdateStmtFunc(db Database) (string, error) {
// 	return nbfmt.Fmt(updateStmtFuncTemp, map[string]interface{}{"DB": db})
// }

func genDeleteStmtMethod(db Database) (string, error) {
	return nbfmt.Fmt(deleteStmtMethodTemp, map[string]interface{}{"DB": db})
}

// func genDeleteStmtFunc(db Database) (string, error) {
// 	return nbfmt.Fmt(deleteStmtFuncTemp, map[string]interface{}{"DB": db})
// }

// func genModelDistinctMethod(db Database) (string, error) {
// 	return nbfmt.Fmt(modelDistinctMethodTemp, map[string]interface{}{"DB": db})
// }

func genModelListDistinctMethod(db Database) (string, error) {
	return nbfmt.Fmt(modelListDistinctMethodTemp, map[string]interface{}{"DB": db})
}

// func genModelListSortByMethod(db Database) (string, error) {
// 	return nbfmt.Fmt(modelListSortByMethodTemp, map[string]interface{}{"DB": db})
// }

func genModelListSortMethod(db Database) (string, error) {
	return nbfmt.Fmt(modelListSortMethodTemp, map[string]interface{}{"DB": db})
}

func genModelExistsMethod(db Database) (string, error) {
	return nbfmt.Fmt(modelExistsMethodTemp, map[string]interface{}{"DB": db})
}

// func genExistsFunc(db Database) (string, error) {
// 	return nbfmt.Fmt(existsFuncTemp, map[string]interface{}{"DB": db})
// }

func genModelClassType(db Database) (string, error) {
	return nbfmt.Fmt(modelClassTypeTemp, map[string]interface{}{"DB": db})
}

func genModelClassInit(db Database) (string, error) {
	return nbfmt.Fmt(modelClassInitTemp, map[string]interface{}{"DB": db})
}

func genClassesFileImport(db Database) (string, error) {
	return classesFileImportTemp, nil
}

func genFieldsMethod(db Database) (string, error) {
	return nbfmt.Fmt(fieldsMethodTemp, map[string]interface{}{"DB": db})
}

// func genGetInstanceClassMethod(db Database) (string, error) {
// 	return nbfmt.Fmt(getInstanceClassMethodTemp, map[string]interface{}{"DB": db})
// }

//Gen generate database definition
func Gen(db Database, outName string) error {
	err := generateClassesFile(db)
	if err != nil {
		return err
	}
	err = generateInstancesFile(db)
	if err != nil {
		return err
	}
	cmd := exec.Command("go", "fmt")
	return cmd.Run()
}

func generateClassesFile(db Database) error {
	buf := bytes.NewBuffer([]byte{})
	s, err := genPackage(db)
	if err != nil {
		return err
	}
	buf.WriteString(s)
	s, err = genClassesFileImport(db)
	if err != nil {
		return err
	}
	buf.WriteString(s)
	s, err = genModelClassType(db)
	if err != nil {
		return err
	}
	buf.WriteString(s)
	s, err = genModelClassInit(db)
	if err != nil {
		return err
	}
	buf.WriteString(s)
	f, err := os.OpenFile("classes.go", os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	defer f.Close()
	n, err := f.Write(buf.Bytes())
	if err != nil {
		return err
	}
	return f.Truncate(int64(n))
}

func generateInstancesFile(db Database) error {
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
	s, err = genModelInstanceType(db)
	if err != nil {
		return err
	}
	buf.WriteString(s)
	s, err = genFieldsMethod(db)
	if err != nil {
		return err
	}
	buf.WriteString(s)
	s, err = genGetFieldMethod(db)
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
	s, err = genSetLastInsertIDMethod(db)
	if err != nil {
		return err
	}
	buf.WriteString(s)
	s, err = genInsertStmtMethod(db)
	if err != nil {
		return err
	}
	buf.WriteString(s)
	s, err = genUpdateStmtMethod(db)
	if err != nil {
		return err
	}
	buf.WriteString(s)
	s, err = genDeleteStmtMethod(db)
	if err != nil {
		return err
	}
	buf.WriteString(s)
	s, err = genModelListDistinctMethod(db)
	if err != nil {
		return err
	}
	buf.WriteString(s)
	s, err = genModelListSortMethod(db)
	if err != nil {
		return err
	}
	buf.WriteString(s)
	s, err = genModelExistsMethod(db)
	if err != nil {
		return err
	}
	buf.WriteString(s)
	f, err := os.OpenFile("instances.go", os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	defer f.Close()
	n, err := f.Write(buf.Bytes())
	if err != nil {
		return err
	}
	return f.Truncate(int64(n))
}
