package acl

import (
	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppote-core/uhppote"
	"reflect"
	"strings"
	"testing"
)

const tsv = `Card Number	From	To	Workshop	Side Door	Front Door	Garage
65537	2020-01-02	2020-10-31	N	N	Y	N
65538	2020-02-03	2020-11-30	Y	N	Y	N
65539	2020-03-04	2020-12-31	N	N	N	N
`

const tsv2 = `Card Number	From	To	D1	D2	D3	D4	Workshop	Side Door	Front Door	Garage
65537	2020-01-02	2020-10-31	Y	Y	N	Y	N	N	Y	N
65538	2020-02-03	2020-11-30	Y	N	Y	Y	Y	N	Y	N
65539	2020-03-04	2020-12-31	N	Y	Y	Y	N	N	N	N
`

func TestParseTSV(t *testing.T) {
	expected := []types.Card{
		types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: []bool{true, false, false, false}},
		types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-11-30"), Doors: []bool{true, false, false, true}},
		types.Card{CardNumber: 65539, From: date("2020-03-04"), To: date("2020-12-31"), Doors: []bool{false, false, false, false}},
	}

	d := uhppote.Device{
		DeviceID: 12345,
		Doors:    []string{"Front Door", "Side Door", "Garage", "Workshop"},
	}

	devices := []*uhppote.Device{&d}
	r := strings.NewReader(tsv)

	m, err := ParseTSV(r, devices)
	if err != nil {
		t.Fatalf("Unexpected error parsing TSV: %v", err)
	}

	if len(m) != len(devices) {
		t.Fatalf("ParseTSV returned invalid ACL (%v)", m)
	}

	for _, device := range devices {
		if l := m[device.DeviceID]; l == nil {
			t.Errorf("ACL missing access list for device ID %v", device.DeviceID)
		} else {
			if len(l) != len(expected) {
				t.Errorf("device %v: record counts do not match - expected %d, got %d", device.DeviceID, len(expected), len(l))
			}

			for _, card := range expected {
				if c, ok := l[card.CardNumber]; !ok {
					t.Errorf("device %v: missing record for card %v", device.DeviceID, card.CardNumber)
				} else if !reflect.DeepEqual(c, card) {
					t.Errorf("device %v: invalid record for card %v\n  expected: %v\n  got:      %v", device.DeviceID, card.CardNumber, card, c)
				}
			}
		}
	}
}

func TestParseTSVWithMultipleDevices(t *testing.T) {
	expected := map[uint32][]types.Card{
		12345: []types.Card{
			types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: []bool{true, false, false, false}},
			types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-11-30"), Doors: []bool{true, false, false, true}},
			types.Card{CardNumber: 65539, From: date("2020-03-04"), To: date("2020-12-31"), Doors: []bool{false, false, false, false}},
		},

		54321: []types.Card{
			types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: []bool{true, true, false, true}},
			types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-11-30"), Doors: []bool{true, false, true, true}},
			types.Card{CardNumber: 65539, From: date("2020-03-04"), To: date("2020-12-31"), Doors: []bool{false, true, true, true}},
		},
	}

	d1 := uhppote.Device{
		DeviceID: 12345,
		Doors:    []string{"Front Door", "Side Door", "Garage", "Workshop"},
	}

	d2 := uhppote.Device{
		DeviceID: 54321,
		Doors:    []string{"D1", "D2", "D3", "D4"},
	}

	devices := []*uhppote.Device{&d1, &d2}
	r := strings.NewReader(tsv2)

	m, err := ParseTSV(r, devices)
	if err != nil {
		t.Fatalf("Unexpected error parsing TSV: %v", err)
	}

	if len(m) != len(devices) {
		t.Fatalf("ParseTSV returned invalid ACL (%v)", m)
	}

	for _, device := range devices {
		e := expected[device.DeviceID]

		if l := m[device.DeviceID]; l == nil {
			t.Errorf("ACL missing access list for device ID %v", device.DeviceID)
		} else {
			if len(l) != len(e) {
				t.Errorf("device %v: record counts do not match - expected %d, got %d", device.DeviceID, len(e), len(l))
			}

			for _, card := range e {
				if c, ok := l[card.CardNumber]; !ok {
					t.Errorf("device %v: missing record for card %v", device.DeviceID, card.CardNumber)
				} else if !reflect.DeepEqual(c, card) {
					t.Errorf("device %v: invalid record for card %v\n  expected: %v\n  got:      %v", device.DeviceID, card.CardNumber, card, c)
				}
			}
		}
	}
}

func TestMakeTSVWithMissingACL(t *testing.T) {
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

	err := MakeTSV(acl, devices, &w)
	if err == nil {
		t.Fatalf("Expected error creating TSV")
	}
}

func TestMakeTSV(t *testing.T) {
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

	expected := `Card Number	From	To	Front Door	Side Door	Garage	Workshop	D1	D2	D3	D4
65536	2020-01-01	2020-12-31	Y	N	Y	N	N	N	N	N
65537	2020-01-01	2020-12-31	Y	N	N	N	Y	Y	N	Y
65538	2020-02-03	2020-11-30	Y	N	N	Y	Y	N	Y	Y
65539	2020-01-03	2020-12-31	N	N	N	N	N	Y	Y	Y
65540	2019-01-01	2021-12-31	N	N	N	N	N	Y	N	Y
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

	err := MakeTSV(acl, devices, &w)
	if err != nil {
		t.Fatalf("Unexpected error creating TSV: %v", err)
	}

	s := w.String()
	if s != expected {
		t.Errorf("Returned incorrect TSV - expected:\n%v\ngot:\n%v\n", expected, s)
	}
}
