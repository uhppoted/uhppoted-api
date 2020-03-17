package acl

import (
	"fmt"
	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppote-core/uhppote"
	"reflect"
	"testing"
	"time"
)

type mock struct {
	getCards       func(uint32) (uint32, error)
	getCardByIndex func(uint32, uint32) (*types.Card, error)
	putCard        func(uint32, types.Card) (bool, error)
	deleteCard     func(uint32, types.Card) (bool, error)
}

func (m *mock) GetCardsN(deviceID uint32) (uint32, error) {
	return m.getCards(deviceID)
}

func (m *mock) GetCardByIndexN(deviceID, index uint32) (*types.Card, error) {
	return m.getCardByIndex(deviceID, index)
}

func (m *mock) PutCardN(deviceID uint32, card types.Card) (bool, error) {
	return m.putCard(deviceID, card)
}

func (m *mock) DeleteCardN(deviceID uint32, card types.Card) (bool, error) {
	return m.deleteCard(deviceID, card)
}

var date = func(s string) types.Date {
	d, _ := time.ParseInLocation("2006-01-02", s, time.Local)
	return types.Date(d)
}

var deviceA = uhppote.Device{
	DeviceID: 12345,
	Doors:    []string{"Front Door", "Side Door", "Garage", "Workshop"},
}

var deviceB = uhppote.Device{
	DeviceID: 54321,
	Doors:    []string{"D1", "D2", "D3", "D4"},
}

var aclA = ACL{
	12345: map[uint32]types.Card{
		65537: types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: []bool{true, false, false, false}},
		65538: types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-11-30"), Doors: []bool{true, false, false, true}},
		65539: types.Card{CardNumber: 65539, From: date("2020-03-04"), To: date("2020-12-31"), Doors: []bool{false, false, false, false}},
	},
}

var cardsA = []types.Card{
	types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: []bool{true, false, false, false}},
	types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-11-30"), Doors: []bool{true, false, false, true}},
	types.Card{CardNumber: 65539, From: date("2020-03-04"), To: date("2020-12-31"), Doors: []bool{false, false, false, false}},
}

func TestGetACL(t *testing.T) {
	expected := aclA
	cards := cardsA

	u := mock{
		getCards: func(deviceID uint32) (uint32, error) {
			return uint32(len(cards)), nil
		},
		getCardByIndex: func(deviceID, index uint32) (*types.Card, error) {
			if int(index) < 0 || int(index) > len(cards) {
				return nil, nil
			}
			return &cards[index-1], nil
		},
	}

	devices := []*uhppote.Device{&deviceA}

	acl, err := GetACL(&u, devices)
	if err != nil {
		t.Fatalf("Unexpected error getting ACL: %v", err)
	}

	if len(acl) != len(devices) {
		t.Errorf("Incorrect ACL record count: expected %v, got %v", len(devices), len(acl))
	}

	if l, ok := acl[deviceA.DeviceID]; !ok {
		t.Errorf("Missing access list for device ID %v", deviceA.DeviceID)
	} else {
		e := expected[deviceA.DeviceID]
		if len(l) != len(e) {
			t.Errorf("device %v: record counts do not match - expected %d, got %d", deviceA.DeviceID, len(e), len(l))
		}

		for _, card := range e {
			if c, ok := l[card.CardNumber]; !ok {
				t.Errorf("device %v: missing record for card %v", deviceA.DeviceID, card.CardNumber)
			} else if !reflect.DeepEqual(c, card) {
				t.Errorf("device %v: invalid record for card %v\n  expected: %v\n  got:      %v", deviceA.DeviceID, card.CardNumber, card, c)
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

	devices := []*uhppote.Device{&deviceA, &deviceB}

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

func TestPutACL(t *testing.T) {
	acl := ACL{
		12345: map[uint32]types.Card{
			65536: types.Card{CardNumber: 65536, From: date("2020-01-01"), To: date("2020-12-31"), Doors: []bool{true, false, true, false}},
			65537: types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: []bool{true, false, false, false}},
			65538: types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-12-31"), Doors: []bool{false, false, false, false}},
		},
	}

	expected := []types.Card{
		types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: []bool{true, false, false, false}},
		types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-12-31"), Doors: []bool{false, false, false, false}},
		types.Card{CardNumber: 65536, From: date("2020-01-01"), To: date("2020-12-31"), Doors: []bool{true, false, true, false}},
	}

	cards := []types.Card{
		types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: []bool{true, false, false, false}},
		types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-11-30"), Doors: []bool{true, false, false, true}},
		types.Card{CardNumber: 65539, From: date("2020-03-04"), To: date("2020-12-31"), Doors: []bool{false, false, false, false}},
	}

	u := mock{
		getCards: func(deviceID uint32) (uint32, error) {
			return uint32(len(cards)), nil
		},
		getCardByIndex: func(deviceID, index uint32) (*types.Card, error) {
			if int(index) < 0 || int(index) > len(cards) {
				return nil, nil
			}
			return &cards[index-1], nil
		},
		putCard: func(deviceID uint32, card types.Card) (bool, error) {
			for ix, c := range cards {
				if c.CardNumber == card.CardNumber {
					cards[ix] = card
					return true, nil
				}
			}

			cards = append(cards, card)

			return true, nil
		},
		deleteCard: func(deviceID uint32, card types.Card) (bool, error) {
			for ix, c := range cards {
				if c.CardNumber == card.CardNumber {
					cards = append(cards[:ix], cards[ix+1:]...)
					return true, nil
				}
			}

			return false, nil
		},
	}

	err := PutACL(&u, acl)
	if err != nil {
		t.Fatalf("Unexpected error putting ACL: %v", err)
	}

	if !reflect.DeepEqual(cards, expected) {
		t.Errorf("Device internal card list not updated correctly:\n    expected:%+v\n    got:     %+v", expected, cards)
	}
}

func TestCompare(t *testing.T) {
	src := ACL{
		12345: {
			923321456: types.Card{CardNumber: 923321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, true, false}},
			823321456: types.Card{CardNumber: 823321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, true, false}},
			623321456: types.Card{CardNumber: 623321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, true, false}},
		},
	}

	dest := ACL{
		12345: {
			923321456: types.Card{CardNumber: 923321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, true, false}},
			823321456: types.Card{CardNumber: 823321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, false, false}},
			723321456: types.Card{CardNumber: 723321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, true, false}},
		},
	}

	expected := map[uint32]Diff{
		12345: Diff{
			Unchanged: []types.Card{
				types.Card{CardNumber: 923321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, true, false}},
			},
			Updated: []types.Card{
				types.Card{CardNumber: 823321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, false, false}},
			},
			Added: []types.Card{
				types.Card{CardNumber: 723321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, true, false}},
			},
			Deleted: []types.Card{
				types.Card{CardNumber: 623321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, true, false}},
			},
		},
	}

	diff, err := Compare(src, dest)
	if err != nil {
		t.Fatalf("Unexpected error comparing ACL: %v", err)
	}

	if diff == nil {
		t.Fatalf("Compare(..) returned 'nil'")
	}

	if !reflect.DeepEqual(diff, expected) {
		t.Fatalf("Compare(..) returned invalid 'diff':\n   expected: %+v\n   got:      %+v", expected, diff)
	}
}

func TestCompareWithMultipleDevices(t *testing.T) {
	src := ACL{
		12345: {
			923321456: types.Card{CardNumber: 923321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, true, false}},
			823321456: types.Card{CardNumber: 823321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, true, false}},
			623321456: types.Card{CardNumber: 623321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, true, false}},
		},
		54321: {
			923321456: types.Card{CardNumber: 923321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, true, false}},
			823321456: types.Card{CardNumber: 823321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, true, false}},
			723321456: types.Card{CardNumber: 723321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, true, false}},
		},
	}

	dest := ACL{
		12345: {
			923321456: types.Card{CardNumber: 923321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, true, false}},
			823321456: types.Card{CardNumber: 823321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, false, false}},
			723321456: types.Card{CardNumber: 723321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, true, false}},
		},
		54321: {
			923321456: types.Card{CardNumber: 923321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, true, false}},
			823321456: types.Card{CardNumber: 823321456, From: date("2020-01-01"), To: date("2020-12-31"), Doors: []bool{false, false, true, false}},
			623321456: types.Card{CardNumber: 623321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, true, false}},
		},
	}

	expected := map[uint32]Diff{
		12345: Diff{
			Unchanged: []types.Card{
				types.Card{CardNumber: 923321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, true, false}},
			},
			Updated: []types.Card{
				types.Card{CardNumber: 823321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, false, false}},
			},
			Added: []types.Card{
				types.Card{CardNumber: 723321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, true, false}},
			},
			Deleted: []types.Card{
				types.Card{CardNumber: 623321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, true, false}},
			},
		},
		54321: Diff{
			Unchanged: []types.Card{
				types.Card{CardNumber: 923321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, true, false}},
			},
			Updated: []types.Card{
				types.Card{CardNumber: 823321456, From: date("2020-01-01"), To: date("2020-12-31"), Doors: []bool{false, false, true, false}},
			},
			Added: []types.Card{
				types.Card{CardNumber: 623321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, true, false}},
			},
			Deleted: []types.Card{
				types.Card{CardNumber: 723321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, true, false}},
			},
		},
	}

	diff, err := Compare(src, dest)
	if err != nil {
		t.Fatalf("Unexpected error comparing ACL: %v", err)
	}

	if diff == nil {
		t.Fatalf("Compare(..) returned 'nil'")
	}

	if !reflect.DeepEqual(diff, expected) {
		t.Fatalf("Compare(..) returned invalid 'diff':\n   expected: %+v\n   got:      %+v", expected, diff)
	}
}
