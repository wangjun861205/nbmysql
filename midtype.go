package nbmysql

import (
	"strconv"
	"time"
)

//Int middle type for scan int64 field
type Int struct {
	Value  int64
	IsNull bool
}

//NewInt create new *Int
func NewInt(value int64, isNull bool) *Int {
	return &Int{value, isNull}
}

//ToGo convert *int to *int64
func (i *Int) ToGo() *int64 {
	if i.IsNull {
		return nil
	}
	return &i.Value
}

//Scan implement sql.Scanner interface for Int
func (i *Int) Scan(v interface{}) error {
	if v != nil {
		i.IsNull = false
		switch val := v.(type) {
		case int64:
			i.Value = val
		case []byte:
			i64, err := strconv.ParseInt(string(val), 10, 64)
			if err != nil {
				return err
			}
			i.Value = i64
		}
		return nil
	}
	i.IsNull = true
	return nil
}

//Float middle type for scan float64 field
type Float struct {
	Value  float64
	IsNull bool
}

//NewFloat create new *Float
func NewFloat(value float64, isNull bool) *Float {
	return &Float{value, isNull}
}

//ToGo convert *Float to *float64
func (f *Float) ToGo() *float64 {
	if f.IsNull {
		return nil
	}
	return &f.Value
}

//Scan implement sql.Scanner interface for Float
func (f *Float) Scan(v interface{}) error {
	if v != nil {
		f.IsNull = false
		switch val := v.(type) {
		case float64:
			f.Value = val
		case []byte:
			f64, err := strconv.ParseFloat(string(val), 64)
			if err != nil {
				return err
			}
			f.Value = f64
		}
		return nil
	}
	f.IsNull = true
	return nil

}

//Bool middle type for scan bool field
type Bool struct {
	Value  bool
	IsNull bool
}

//NewBool create new Bool
func NewBool(value bool, isNull bool) *Bool {
	return &Bool{value, isNull}
}

//ToGo convert *Bool to *bool
func (b *Bool) ToGo() *bool {
	if b.IsNull {
		return nil
	}
	return &b.Value
}

//Scan implement sql.Scanner interface for Bool
func (b *Bool) Scan(v interface{}) error {
	if v != nil {
		b.IsNull = false
		switch val := v.(type) {
		case bool:
			b.Value = val
		case int64:
			bl, err := strconv.ParseBool(strconv.FormatInt(val, 64))
			if err != nil {
				return err
			}
			b.Value = bl
		}
		return nil
	}
	b.IsNull = true
	return nil
}

//String middle type for scan string field
type String struct {
	Value  string
	IsNull bool
}

//NewString create new *String
func NewString(value string, isNull bool) *String {
	return &String{value, isNull}
}

//ToGo convert *String to *string
func (s *String) ToGo() *string {
	if s.IsNull {
		return nil
	}
	return &s.Value
}

//Scan implement sql.Scanner interface for String
func (s *String) Scan(v interface{}) error {
	if v != nil {
		s.IsNull = false
		s.Value = string(v.([]byte))
		return nil
	}
	s.IsNull = true
	return nil
}

//Time middle type for scan time.Time field
type Time struct {
	Value  time.Time
	IsNull bool
}

//NewTime create new *time.Time
func NewTime(value time.Time, isNull bool) *Time {
	return &Time{value, isNull}
}

//ToGo convert *Time to *time
func (t *Time) ToGo() *time.Time {
	if t.IsNull {
		return nil
	}
	return &t.Value
}

//Scan implement sql.Scanner for Time
func (t *Time) Scan(v interface{}) error {
	if v != nil {
		t.IsNull = false
		tv, err := time.Parse("2006-01-02 15:04:05", string(v.([]byte)))
		if err != nil {
			tv, err = time.Parse("2006-01-02", string(v.([]byte)))
			if err != nil {
				return err
			}
		}
		t.Value = tv
		return nil
	}
	t.IsNull = true
	return nil
}
