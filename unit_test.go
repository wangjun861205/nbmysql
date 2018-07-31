package nbmysql

import "testing"

func TestUnit(t *testing.T) {
	db, err := ParseDatabase("example.nbdb")
	if err != nil {
		t.Fatal(err)
	}
	err = Gen(db, "dalianModel.go")
	if err != nil {
		t.Fatal(err)
	}
	for _, tab := range db.Tables {
		err := db.CreateTableIfNotExists(tab)
		if err != nil {
			t.Fatal(err)
		}
	}
}
