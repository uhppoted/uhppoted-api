package acl

import (
	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppote-core/uhppote"
	"reflect"
	"testing"
)

func TestMakeTable(t *testing.T) {
	acl := ACL{
		12345: map[uint32]types.Card{
			65537: types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: []bool{true, false, false, false}},
			65538: types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-11-30"), Doors: []bool{true, false, false, true}},
			65539: types.Card{CardNumber: 65539, From: date("2020-03-04"), To: date("2020-12-31"), Doors: []bool{false, false, false, false}},
		},
	}

	expected := table{
		header: []string{"Card Number", "From", "To", "Front Door", "Side Door", "Garage", "Workshop"},
		records: [][]string{
			[]string{"65537", "2020-01-02", "2020-10-31", "Y", "N", "N", "N"},
			[]string{"65538", "2020-02-03", "2020-11-30", "Y", "N", "N", "Y"},
			[]string{"65539", "2020-03-04", "2020-12-31", "N", "N", "N", "N"},
		},
	}

	d := uhppote.Device{
		DeviceID: 12345,
		Doors:    []string{"Front Door", "Side Door", "Garage", "Workshop"},
	}

	devices := []*uhppote.Device{&d}

	rs, err := MakeTable(acl, devices)
	if err != nil {
		t.Fatalf("Unexpected error creating table: %v", err)
	}

	if rs == nil {
		t.Fatalf("MakeTable returned invalid result: %v", rs)
	}

	if !reflect.DeepEqual(*rs, expected) {
		t.Errorf("Returned incorrect table - expected:\n%+v\ngot:\n%+v\n", expected, *rs)
	}
}

func TestMakeTableWithMultipleDevices(t *testing.T) {
	acl := ACL{
		12345: map[uint32]types.Card{
			65537: types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: []bool{true, false, false, false}},
			65538: types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-11-30"), Doors: []bool{true, false, false, true}},
			65539: types.Card{CardNumber: 65539, From: date("2020-03-04"), To: date("2020-12-31"), Doors: []bool{false, false, false, false}},
		},
		54321: map[uint32]types.Card{
			65537: types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: []bool{true, true, false, true}},
			65538: types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-11-30"), Doors: []bool{true, false, true, true}},
			65539: types.Card{CardNumber: 65539, From: date("2020-03-04"), To: date("2020-12-31"), Doors: []bool{false, true, true, true}},
		},
	}

	expected := table{
		header: []string{"Card Number", "From", "To", "Front Door", "Side Door", "Garage", "Workshop", "D1", "D2", "D3", "D4"},
		records: [][]string{
			[]string{"65537", "2020-01-02", "2020-10-31", "Y", "N", "N", "N", "Y", "Y", "N", "Y"},
			[]string{"65538", "2020-02-03", "2020-11-30", "Y", "N", "N", "Y", "Y", "N", "Y", "Y"},
			[]string{"65539", "2020-03-04", "2020-12-31", "N", "N", "N", "N", "N", "Y", "Y", "Y"},
		},
	}

	devices := []*uhppote.Device{
		&uhppote.Device{
			DeviceID: 12345,
			Doors:    []string{"Front Door", "Side Door", "Garage", "Workshop"},
		},
		&uhppote.Device{
			DeviceID: 54321,
			Doors:    []string{"D1", "D2", "D3", "D4"},
		},
	}

	rs, err := MakeTable(acl, devices)
	if err != nil {
		t.Fatalf("Unexpected error creating table: %v", err)
	}

	if rs == nil {
		t.Fatalf("MakeTable returned invalid result: %v", rs)
	}

	if !reflect.DeepEqual(*rs, expected) {
		t.Errorf("Returned incorrect table - expected:\n%+v\ngot:\n%+v\n", expected, *rs)
	}
}

func TestMakeTableWithMissingACL(t *testing.T) {
	acl := ACL{
		12345: map[uint32]types.Card{
			65537: types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: []bool{true, false, false, false}},
			65538: types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-11-30"), Doors: []bool{true, false, false, true}},
			65539: types.Card{CardNumber: 65539, From: date("2020-03-04"), To: date("2020-12-31"), Doors: []bool{false, false, false, false}},
		},
	}

	devices := []*uhppote.Device{
		&uhppote.Device{
			DeviceID: 12345,
			Doors:    []string{"Front Door", "Side Door", "Garage", "Workshop"},
		},
		&uhppote.Device{
			DeviceID: 54321,
			Doors:    []string{"D1", "D2", "D3", "D4"},
		},
	}

	_, err := MakeTable(acl, devices)
	if err == nil {
		t.Fatalf("Expected error creating table")
	}
}

func TestMakeRecordsetWithMismatchedDates(t *testing.T) {
	acl := ACL{
		12345: map[uint32]types.Card{
			65537: types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: []bool{true, false, false, false}},
			65538: types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-11-30"), Doors: []bool{true, false, false, true}},
			65539: types.Card{CardNumber: 65539, From: date("2020-03-04"), To: date("2020-12-31"), Doors: []bool{false, false, false, false}},
		},
		54321: map[uint32]types.Card{
			65537: types.Card{CardNumber: 65537, From: date("2020-01-01"), To: date("2020-12-31"), Doors: []bool{true, true, false, true}},
			65538: types.Card{CardNumber: 65538, From: date("2020-03-01"), To: date("2020-10-31"), Doors: []bool{true, false, true, true}},
			65539: types.Card{CardNumber: 65539, From: date("2020-01-03"), To: date("2020-11-30"), Doors: []bool{false, true, true, true}},
		},
	}

	expected := table{
		header: []string{"Card Number", "From", "To", "Front Door", "Side Door", "Garage", "Workshop", "D1", "D2", "D3", "D4"},
		records: [][]string{
			[]string{"65537", "2020-01-01", "2020-12-31", "Y", "N", "N", "N", "Y", "Y", "N", "Y"},
			[]string{"65538", "2020-02-03", "2020-11-30", "Y", "N", "N", "Y", "Y", "N", "Y", "Y"},
			[]string{"65539", "2020-01-03", "2020-12-31", "N", "N", "N", "N", "N", "Y", "Y", "Y"},
		},
	}
	devices := []*uhppote.Device{
		&uhppote.Device{
			DeviceID: 12345,
			Doors:    []string{"Front Door", "Side Door", "Garage", "Workshop"},
		},
		&uhppote.Device{
			DeviceID: 54321,
			Doors:    []string{"D1", "D2", "D3", "D4"},
		},
	}

	rs, err := MakeTable(acl, devices)
	if err != nil {
		t.Fatalf("Unexpected error creating table: %v", err)
	}

	if rs == nil {
		t.Fatalf("MakeTable returned invalid result: %v", rs)
	}

	if !reflect.DeepEqual(*rs, expected) {
		t.Errorf("Returned incorrect table - expected:\n%+v\ngot:\n%+v\n", expected, *rs)
	}
}

func TestMakeTableWithMismatchedCards(t *testing.T) {
	acl := ACL{
		12345: map[uint32]types.Card{
			65536: types.Card{CardNumber: 65536, From: date("2020-01-01"), To: date("2020-12-31"), Doors: []bool{true, false, true, false}},
			65537: types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: []bool{true, false, false, false}},
			65538: types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-11-30"), Doors: []bool{true, false, false, true}},
			65539: types.Card{CardNumber: 65539, From: date("2020-03-04"), To: date("2020-12-31"), Doors: []bool{false, false, false, false}},
		},
		54321: map[uint32]types.Card{
			65537: types.Card{CardNumber: 65537, From: date("2020-01-01"), To: date("2020-12-31"), Doors: []bool{true, true, false, true}},
			65538: types.Card{CardNumber: 65538, From: date("2020-03-01"), To: date("2020-10-31"), Doors: []bool{true, false, true, true}},
			65539: types.Card{CardNumber: 65539, From: date("2020-01-03"), To: date("2020-11-30"), Doors: []bool{false, true, true, true}},
			65540: types.Card{CardNumber: 65540, From: date("2019-01-01"), To: date("2021-12-31"), Doors: []bool{false, true, false, true}},
		},
	}

	expected := table{
		header: []string{"Card Number", "From", "To", "Front Door", "Side Door", "Garage", "Workshop", "D1", "D2", "D3", "D4"},
		records: [][]string{
			[]string{"65536", "2020-01-01", "2020-12-31", "Y", "N", "Y", "N", "N", "N", "N", "N"},
			[]string{"65537", "2020-01-01", "2020-12-31", "Y", "N", "N", "N", "Y", "Y", "N", "Y"},
			[]string{"65538", "2020-02-03", "2020-11-30", "Y", "N", "N", "Y", "Y", "N", "Y", "Y"},
			[]string{"65539", "2020-01-03", "2020-12-31", "N", "N", "N", "N", "N", "Y", "Y", "Y"},
			[]string{"65540", "2019-01-01", "2021-12-31", "N", "N", "N", "N", "N", "Y", "N", "Y"},
		},
	}

	devices := []*uhppote.Device{
		&uhppote.Device{
			DeviceID: 12345,
			Doors:    []string{"Front Door", "Side Door", "Garage", "Workshop"},
		},
		&uhppote.Device{
			DeviceID: 54321,
			Doors:    []string{"D1", "D2", "D3", "D4"},
		},
	}

	rs, err := MakeTable(acl, devices)
	if err != nil {
		t.Fatalf("Unexpected error creating table: %v", err)
	}

	if rs == nil {
		t.Fatalf("MakeTable returned invalid result: %v", rs)
	}

	if !reflect.DeepEqual(*rs, expected) {
		t.Errorf("Returned incorrect table - expected:\n%+v\ngot:\n%+v\n", expected, *rs)
	}
}