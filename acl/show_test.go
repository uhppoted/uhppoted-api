package acl

import (
	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppote-core/uhppote"
	"reflect"
	"testing"
)

func TestGetCard(t *testing.T) {
	expected := map[string]DateRange{
		"Front Door": DateRange{From: *date("2020-02-03"), To: *date("2020-11-30")},
		"Workshop":   DateRange{From: *date("2020-02-03"), To: *date("2020-11-30")},
	}

	cards := []types.CardX{
		types.CardX{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]bool{1: true, 2: false, 3: false, 4: false}},
		types.CardX{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-11-30"), Doors: map[uint8]bool{1: true, 2: false, 3: false, 4: true}},
		types.CardX{CardNumber: 65539, From: date("2020-03-04"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: false, 4: false}},
	}

	u := mock{
		getCardByID: func(deviceID, cardID uint32) (*types.CardX, error) {
			for _, c := range cards {
				if c.CardNumber == cardID {
					return &c, nil
				}
			}
			return nil, nil
		},
	}

	devices := []*uhppote.Device{
		&uhppote.Device{
			DeviceID: 12345,
			Doors:    []string{"Front Door", "Side Door", "Garage", "Workshop"},
		},
	}

	doors, err := GetCard(&u, devices, 65538)
	if err != nil {
		t.Fatalf("Unexpected error getting card ACL: %v", err)
	}

	if !reflect.DeepEqual(doors, expected) {
		t.Errorf("invalid ACL for card %v\n  expected: %v\n  got:      %v", 65538, expected, doors)
	}
}

func TestGetCardWithUnknownCard(t *testing.T) {
	expected := map[string]DateRange{}

	cards := []types.CardX{
		types.CardX{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]bool{1: true, 2: false, 3: false, 4: false}},
		types.CardX{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-11-30"), Doors: map[uint8]bool{1: true, 2: false, 3: false, 4: true}},
		types.CardX{CardNumber: 65539, From: date("2020-03-04"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: false, 4: false}},
	}

	u := mock{
		getCardByID: func(deviceID, cardID uint32) (*types.CardX, error) {
			for _, c := range cards {
				if c.CardNumber == cardID {
					return &c, nil
				}
			}
			return nil, nil
		},
	}

	devices := []*uhppote.Device{
		&uhppote.Device{
			DeviceID: 12345,
			Doors:    []string{"Front Door", "Side Door", "Garage", "Workshop"},
		},
	}

	doors, err := GetCard(&u, devices, 65536)
	if err != nil {
		t.Fatalf("Unexpected error getting card ACL: %v", err)
	}

	if !reflect.DeepEqual(doors, expected) {
		t.Errorf("invalid ACL for card %v\n  expected: %v\n  got:      %v", 65538, expected, doors)
	}
}

func TestGetCardWithMultipleDevices(t *testing.T) {
	expected := map[string]DateRange{
		"Front Door": DateRange{From: *date("2020-02-03"), To: *date("2020-11-30")},
		"Workshop":   DateRange{From: *date("2020-02-03"), To: *date("2020-11-30")},
		"D2":         DateRange{From: *date("2020-01-01"), To: *date("2020-12-31")},
	}

	cards := map[uint32][]types.CardX{
		12345: []types.CardX{
			types.CardX{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]bool{1: true, 2: false, 3: false, 4: false}},
			types.CardX{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-11-30"), Doors: map[uint8]bool{1: true, 2: false, 3: false, 4: true}},
			types.CardX{CardNumber: 65539, From: date("2020-03-04"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: false, 4: false}},
		},

		54321: []types.CardX{
			types.CardX{CardNumber: 65536, From: date("2020-01-01"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: false, 4: true}},
			types.CardX{CardNumber: 65537, From: date("2020-01-01"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: true, 4: false}},
			types.CardX{CardNumber: 65538, From: date("2020-01-01"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: true, 3: false, 4: false}},
			types.CardX{CardNumber: 65539, From: date("2020-01-01"), To: date("2020-12-31"), Doors: map[uint8]bool{1: true, 2: false, 3: false, 4: false}},
		},
	}

	u := mock{
		getCardByID: func(deviceID, cardID uint32) (*types.CardX, error) {
			l := cards[deviceID]
			for _, c := range l {
				if c.CardNumber == cardID {
					return &c, nil
				}
			}
			return nil, nil
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

	doors, err := GetCard(&u, devices, 65538)
	if err != nil {
		t.Fatalf("Unexpected error getting card ACL: %v", err)
	}

	if !reflect.DeepEqual(doors, expected) {
		t.Errorf("invalid ACL for card %v\n  expected: %v\n  got:      %v", 65538, expected, doors)
	}
}
