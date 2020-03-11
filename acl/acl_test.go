package acl

import (
	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppote-core/uhppote"
	"reflect"
	"strings"
	"testing"
	"time"
)

const tsv = `Card Number	From	To	Workshop	Side Door	Front Door	Garage
65537	2020-01-02	2020-10-31	N	N	Y	N
65538	2020-02-03	2020-11-30	Y	N	Y	N
65539	2020-03-04	2020-12-31	N	N	N	N
`

func TestParseTSV(t *testing.T) {
	date := func(s string) types.Date {
		d, _ := time.ParseInLocation("2006-01-02", s, time.Local)
		return types.Date(d)
	}

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
			for _, card := range expected {
				if c := l[card.CardNumber]; c == nil {
					t.Errorf("device %v: missing record for card %v", device.DeviceID, card.CardNumber)
				} else if !reflect.DeepEqual(*c, card) {
					t.Errorf("device %v: invalid record for card %v\n  expected: %v\n  got:      %v", device.DeviceID, card.CardNumber, card, c)
				}
			}
		}
	}
}
