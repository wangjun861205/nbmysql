package nbmysql

import (
	"strconv"
	"time"
)

type Int struct {
	Value  int64
	IsNull bool
}

func NewInt(value int64, isNull bool) *Int {
	return &Int{value, isNull}
}

func (i *Int) ToGo() *int64 {
	if i.IsNull {
		return nil
	}
	return &i.Value
}

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

type Float struct {
	Value  float64
	IsNull bool
}

func NewFloat(value float64, isNull bool) *Float {
	return &Float{value, isNull}
}

func (f *Float) ToGo() *float64 {
	if f.IsNull {
		return nil
	}
	return &f.Value
}

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

type Bool struct {
	Value  bool
	IsNull bool
}

func NewBool(value bool, isNull bool) *Bool {
	return &Bool{value, isNull}
}

func (b *Bool) ToGo() *bool {
	if b.IsNull {
		return nil
	}
	return &b.Value
}

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

type String struct {
	Value  string
	IsNull bool
}

func NewString(value string, isNull bool) *String {
	return &String{value, isNull}
}

func (s *String) ToGo() *string {
	if s.IsNull {
		return nil
	}
	return &s.Value
}

func (s *String) Scan(v interface{}) error {
	if v != nil {
		s.IsNull = false
		s.Value = string(v.([]byte))
		return nil
	}
	s.IsNull = true
	return nil
}

type Time struct {
	Value  time.Time
	IsNull bool
}

func NewTime(value time.Time, isNull bool) *Time {
	return &Time{value, isNull}
}

func (t *Time) ToGo() *time.Time {
	if t.IsNull {
		return nil
	}
	return &t.Value
}

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
