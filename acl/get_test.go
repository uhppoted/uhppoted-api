package acl

import (
	"fmt"
	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppote-core/uhppote"
	"reflect"
	"testing"
)

func TestGetACL(t *testing.T) {
	expected := ACL{
		12345: map[uint32]types.Card{
			65537: types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]bool{1: true, 2: false, 3: false, 4: false}},
			65538: types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-11-30"), Doors: map[uint8]bool{1: true, 2: false, 3: false, 4: true}},
			65539: types.Card{CardNumber: 65539, From: date("2020-03-04"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: false, 4: false}},
		},
	}

	cards := []types.Card{
		types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]bool{1: true, 2: false, 3: false, 4: false}},
		types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-11-30"), Doors: map[uint8]bool{1: true, 2: false, 3: false, 4: true}},
		types.Card{CardNumber: 65539, From: date("2020-03-04"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: false, 4: false}},
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
	}

	devices := []*uhppote.Device{
		&uhppote.Device{
			DeviceID: 12345,
			Doors:    []string{"Front Door", "Side Door", "Garage", "Workshop"},
		},
	}

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
			65537: types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]bool{1: true, 2: false, 3: false, 4: false}},
			65538: types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-11-30"), Doors: map[uint8]bool{1: true, 2: false, 3: false, 4: true}},
			65539: types.Card{CardNumber: 65539, From: date("2020-03-04"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: false, 4: false}},
		},

		54321: map[uint32]types.Card{
			65536: types.Card{CardNumber: 65536, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]bool{1: false, 2: false, 3: false, 4: true}},
			65537: types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]bool{1: false, 2: false, 3: true, 4: false}},
			65538: types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-11-30"), Doors: map[uint8]bool{1: false, 2: true, 3: false, 4: false}},
			65539: types.Card{CardNumber: 65539, From: date("2020-03-04"), To: date("2020-12-31"), Doors: map[uint8]bool{1: true, 2: false, 3: false, 4: false}},
		},
	}

	cards := map[uint32][]types.Card{
		12345: []types.Card{
			types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]bool{1: true, 2: false, 3: false, 4: false}},
			types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-11-30"), Doors: map[uint8]bool{1: true, 2: false, 3: false, 4: true}},
			types.Card{CardNumber: 65539, From: date("2020-03-04"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: false, 4: false}},
		},

		54321: []types.Card{
			types.Card{CardNumber: 65536, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]bool{1: false, 2: false, 3: false, 4: true}},
			types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]bool{1: false, 2: false, 3: true, 4: false}},
			types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-11-30"), Doors: map[uint8]bool{1: false, 2: true, 3: false, 4: false}},
			types.Card{CardNumber: 65539, From: date("2020-03-04"), To: date("2020-12-31"), Doors: map[uint8]bool{1: true, 2: false, 3: false, 4: false}},
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
