package nbmysql

import "errors"

var ErrRecordNotExists = errors.New("record not exists")
var ErrDupKey = errors.New("duplicate key")
