package acl

import (
	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppote-core/uhppote"
	"strings"
	"testing"
)

func TestMakeFlatFileWithMissingACL(t *testing.T) {
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

	var w strings.Builder

	err := MakeFlatFile(acl, devices, &w)
	if err == nil {
		t.Fatalf("Expected error creating flat file")
	}
}

func TestMakeFlatFile(t *testing.T) {
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

	expected := `Card Number  From        To          Front Door  Side Door  Garage  Workshop  D1  D2  D3  D4
65536        2020-01-01  2020-12-31  Y           N          Y       N         N   N   N   N 
65537        2020-01-01  2020-12-31  Y           N          N       N         Y   Y   N   Y 
65538        2020-02-03  2020-11-30  Y           N          N       Y         Y   N   Y   Y 
65539        2020-01-03  2020-12-31  N           N          N       N         N   Y   Y   Y 
65540        2019-01-01  2021-12-31  N           N          N       N         N   Y   N   Y 
`
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

	var w strings.Builder

	err := MakeFlatFile(acl, devices, &w)
	if err != nil {
		t.Fatalf("Unexpected error creating flat file: %v", err)
	}

	s := w.String()
	if s != expected {
		t.Errorf("Returned incorrect flat file - expected:\n%v\ngot:\n%v\n", expected, s)
	}
}