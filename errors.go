package nbmysql

import "errors"

var ErrRecordNotExists = errors.New("record not exists")
var ErrDupKey = errors.New("duplicate key")
var ErrFkConstraint = errors.New("foreign key constraint fails(MySQL Error Number: 1452)")
