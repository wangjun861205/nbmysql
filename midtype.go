package nbmysql

import "time"

type Int struct {
	Value  int64
	IsNull bool
}

func (i *Int) Scan(v interface{}) error {
	if v != nil {
		i.IsNull = false
		i.Value = v.(int64)
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
		f.Value = v.(float64)
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
		b.Value = v.(bool)
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
		t.Value = v.(time.Time)
		return nil
	}
	t.IsNull = true
	return nil
}
