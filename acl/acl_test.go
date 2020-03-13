package acl

import (
	"fmt"
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

const tsv2 = `Card Number	From	To	D1	D2	D3	D4	Workshop	Side Door	Front Door	Garage
65537	2020-01-02	2020-10-31	Y	Y	N	Y	N	N	Y	N
65538	2020-02-03	2020-11-30	Y	N	Y	Y	Y	N	Y	N
65539	2020-03-04	2020-12-31	N	Y	Y	Y	N	N	N	N
`

var date = func(s string) types.Date {
	d, _ := time.ParseInLocation("2006-01-02", s, time.Local)
	return types.Date(d)
}

type mock struct {
	getCards       func(uint32) (uint32, error)
	getCardByIndex func(uint32, uint32) (*types.Card, error)
}

func (m *mock) GetCardsN(deviceID uint32) (uint32, error) {
	return m.getCards(deviceID)
}

func (m *mock) GetCardByIndexN(deviceID, index uint32) (*types.Card, error) {
	return m.getCardByIndex(deviceID, index)
}

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

func TestGetACL(t *testing.T) {
	expected := ACL{
		12345: map[uint32]types.Card{
			65537: types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: []bool{true, false, false, false}},
			65538: types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-11-30"), Doors: []bool{true, false, false, true}},
			65539: types.Card{CardNumber: 65539, From: date("2020-03-04"), To: date("2020-12-31"), Doors: []bool{false, false, false, false}},
		},
	}

	cards := []types.Card{
		types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: []bool{true, false, false, false}},
		types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-11-30"), Doors: []bool{true, false, false, true}},
		types.Card{CardNumber: 65539, From: date("2020-03-04"), To: date("2020-12-31"), Doors: []bool{false, false, false, false}},
	}

	u := mock{
		getCards: func(deviceID uint32) (uint32, error) {
			return 3, nil
		},
		getCardByIndex: func(deviceID, index uint32) (*types.Card, error) {
			if int(index) < 0 || int(index) > len(cards) {
				return nil, nil
			}
			return &cards[index-1], nil
		},
	}

	d := uhppote.Device{
		DeviceID: 12345,
		Doors:    []string{"Front Door", "Side Door", "Garage", "Workshop"},
	}

	devices := []*uhppote.Device{&d}

	acl, err := GetACL(&u, devices)
	if err != nil {
		t.Fatalf("Unexpected error getting ACL: %v", err)
	}

	if len(acl) != len(devices) {
		t.Errorf("Incorrect ACL record count: expected %v, got %v", len(devices), len(acl))
	}

	if l, ok := acl[d.DeviceID]; !ok {
		t.Errorf("Missing access list for device ID %v", d.DeviceID)
	} else {
		e := expected[d.DeviceID]
		if len(l) != len(e) {
			t.Errorf("device %v: record counts do not match - expected %d, got %d", d.DeviceID, len(e), len(l))
		}

		for _, card := range e {
			if c, ok := l[card.CardNumber]; !ok {
				t.Errorf("device %v: missing record for card %v", d.DeviceID, card.CardNumber)
			} else if !reflect.DeepEqual(c, card) {
				t.Errorf("device %v: invalid record for card %v\n  expected: %v\n  got:      %v", d.DeviceID, card.CardNumber, card, c)
			}
		}
	}
}

func TestGetACLWithMultipleDevices(t *testing.T) {
	expected := ACL{
		12345: map[uint32]types.Card{
			65537: types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: []bool{true, false, false, false}},
			65538: types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-11-30"), Doors: []bool{true, false, false, true}},
			65539: types.Card{CardNumber: 65539, From: date("2020-03-04"), To: date("2020-12-31"), Doors: []bool{false, false, false, false}},
		},

		54321: map[uint32]types.Card{
			65536: types.Card{CardNumber: 65536, From: date("2020-01-02"), To: date("2020-10-31"), Doors: []bool{false, false, false, true}},
			65537: types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: []bool{false, false, true, false}},
			65538: types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-11-30"), Doors: []bool{false, true, false, false}},
			65539: types.Card{CardNumber: 65539, From: date("2020-03-04"), To: date("2020-12-31"), Doors: []bool{true, false, false, false}},
		},
	}

	cards := map[uint32][]types.Card{
		12345: []types.Card{
			types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: []bool{true, false, false, false}},
			types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-11-30"), Doors: []bool{true, false, false, true}},
			types.Card{CardNumber: 65539, From: date("2020-03-04"), To: date("2020-12-31"), Doors: []bool{false, false, false, false}},
		},

		54321: []types.Card{
			types.Card{CardNumber: 65536, From: date("2020-01-02"), To: date("2020-10-31"), Doors: []bool{false, false, false, true}},
			types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: []bool{false, false, true, false}},
			types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-11-30"), Doors: []bool{false, true, false, false}},
			types.Card{CardNumber: 65539, From: date("2020-03-04"), To: date("2020-12-31"), Doors: []bool{true, false, false, false}},
		},
	}

	u := mock{
		getCards: func(deviceID uint32) (uint32, error) {
			list, ok := cards[deviceID]
			if !ok {
				return 0, fmt.Errorf("Unexpected device: %v", deviceID)
			}

			return uint32(len(list)), nil
		},

		getCardByIndex: func(deviceID, index uint32) (*types.Card, error) {
			list, ok := cards[deviceID]
			if !ok {
				return nil, fmt.Errorf("Unexpected device: %v", deviceID)
			}

			if int(index) < 0 || int(index) > len(list) {
				return nil, nil
			}
			return &list[index-1], nil
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

	acl, err := GetACL(&u, devices)
	if err != nil {
		t.Fatalf("Unexpected error getting ACL: %v", err)
	}

	if len(acl) != len(devices) {
		t.Errorf("Incorrect ACL record count: expected %v, got %v", len(devices), len(acl))
	}

	for _, d := range devices {
		if l, ok := acl[d.DeviceID]; !ok {
			t.Errorf("Missing access list for device ID %v", d.DeviceID)
		} else {
			e := expected[d.DeviceID]
			if len(l) != len(e) {
				t.Errorf("device %v: record counts do not match - expected %d, got %d", d.DeviceID, len(e), len(l))
			}

			for _, card := range e {
				if c, ok := l[card.CardNumber]; !ok {
					t.Errorf("device %v: missing record for card %v", d.DeviceID, card.CardNumber)
				} else if !reflect.DeepEqual(c, card) {
					t.Errorf("device %v: invalid record for card %v\n  expected: %v\n  got:      %v", d.DeviceID, card.CardNumber, card, c)
				}
			}
		}
	}
}
