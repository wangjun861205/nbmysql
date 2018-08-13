package nbmysql

import "errors"

//ErrRecordNotExists Error when query return empty result
var ErrRecordNotExists = errors.New("record not exists")

//ErrDupKey Error when insert on duplicate key
var ErrDupKey = errors.New("duplicate key")

//ErrFkConstraint Error when insert or update statement not compliant to existed foreign key constraint
var ErrFkConstraint = errors.New("foreign key constraint fails(MySQL Error Number: 1452)")

//ErrModelNotStoredInDB Error when model object has not linked to a database row
var ErrModelNotStoredInDB = errors.New("model is not strored in database")
