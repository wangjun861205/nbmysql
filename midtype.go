package nbmysql

import (
	"strconv"
	"time"
)

type Int struct {
	Value  int64
	IsNull bool
}

func (i *Int) Scan(v interface{}) error {
	if v != nil {
		i.IsNull = false
		i64, err := strconv.ParseInt(string(v.([]byte)), 10, 64)
		if err != nil {
			return err
		}
		i.Value = i64
		return nil
	}
	i.IsNull = true
	return nil
}

type Float struct {
	Value  float64
	IsNull bool
}

func (f *Float) Scan(v interface{}) error {
	if v != nil {
		f.IsNull = false
		f64, err := strconv.ParseFloat(string(v.([]byte)), 64)
		if err != nil {
			return err
		}
		f.Value = f64
		return nil
	}
	f.IsNull = true
	return nil

}

type Bool struct {
	Value  bool
	IsNull bool
}

func (b *Bool) Scan(v interface{}) error {
	if v != nil {
		b.IsNull = false
		bl, err := strconv.ParseBool(string(v.([]byte)))
		if err != nil {
			return err
		}
		b.Value = bl
		return nil
	}
	b.IsNull = true
	return nil
}

type String struct {
	Value  string
	IsNull bool
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

func (t *Time) Scan(v interface{}) error {
	if v != nil {
		t.IsNull = false
		tv, err := time.Parse("2006-01-02 15:04:05", string(v.([]byte)))
		if err != nil {
			return err
		}
		t.Value = tv
		return nil
	}
	t.IsNull = true
	return nil
}
