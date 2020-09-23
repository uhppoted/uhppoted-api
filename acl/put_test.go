package acl

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/uhppoted/uhppote-core/types"
)

func TestPutACL(t *testing.T) {
	acl := ACL{
		12345: map[uint32]types.Card{
			65536: types.Card{CardNumber: 65536, From: date("2020-03-04"), To: date("2020-12-31"), Doors: map[uint8]bool{1: true, 2: false, 3: true, 4: false}},
			65537: types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]bool{1: true, 2: false, 3: false, 4: false}},
			65538: types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: false, 4: false}},
		},
	}

	expected := []types.Card{
		types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]bool{1: true, 2: false, 3: false, 4: false}},
		types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: false, 4: false}},
		types.Card{CardNumber: 65536, From: date("2020-03-04"), To: date("2020-12-31"), Doors: map[uint8]bool{1: true, 2: false, 3: true, 4: false}},
	}

	report := map[uint32]Report{
		12345: Report{
			Unchanged: []uint32{65537},
			Updated:   []uint32{65538},
			Added:     []uint32{65536},
			Deleted:   []uint32{65539},
			Failed:    []uint32{},
			Errored:   []uint32{},
			Errors:    []error{},
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

	rpt, err := PutACL(&u, acl, false)
	if err != nil {
		t.Fatalf("Unexpected error putting ACL: %v", err)
	}

	if !reflect.DeepEqual(cards, expected) {
		t.Errorf("Device internal card list not updated correctly:\n    expected:%+v\n    got:     %+v", expected, cards)
	}

	if !reflect.DeepEqual(rpt, report) {
		t.Errorf("Returned report does not match expected:\n    expected:%+v\n    got:     %+v", report, rpt)
	}
}

func TestPutACLDryRun(t *testing.T) {
	acl := ACL{
		12345: map[uint32]types.Card{
			65536: types.Card{CardNumber: 65536, From: date("2020-03-04"), To: date("2020-12-31"), Doors: map[uint8]bool{1: true, 2: false, 3: true, 4: false}},
			65537: types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]bool{1: true, 2: false, 3: false, 4: false}},
			65538: types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: false, 4: false}},
		},
	}

	expected := []types.Card{
		types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]bool{1: true, 2: false, 3: false, 4: false}},
		types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-11-30"), Doors: map[uint8]bool{1: true, 2: false, 3: false, 4: true}},
		types.Card{CardNumber: 65539, From: date("2020-03-04"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: false, 4: false}},
	}

	report := map[uint32]Report{
		12345: Report{
			Unchanged: []uint32{65537},
			Updated:   []uint32{65538},
			Added:     []uint32{65536},
			Deleted:   []uint32{65539},
			Failed:    []uint32{},
			Errored:   []uint32{},
			Errors:    []error{},
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

	rpt, err := PutACL(&u, acl, true)
	if err != nil {
		t.Fatalf("Unexpected error putting ACL: %v", err)
	}

	if !reflect.DeepEqual(cards, expected) {
		t.Errorf("Device internal card list not updated correctly:\n    expected:%+v\n    got:     %+v", expected, cards)
	}

	if !reflect.DeepEqual(rpt, report) {
		t.Errorf("Returned report does not match expected:\n    expected:%+v\n    got:     %+v", report, rpt)
	}
}

func TestPutACLWithMultipleDevices(t *testing.T) {
	acl := ACL{
		12345: map[uint32]types.Card{
			65536: types.Card{CardNumber: 65536, From: date("2020-01-01"), To: date("2020-12-31"), Doors: map[uint8]bool{1: true, 2: false, 3: true, 4: false}},
			65537: types.Card{CardNumber: 65537, From: date("2020-01-01"), To: date("2020-12-31"), Doors: map[uint8]bool{1: true, 2: false, 3: false, 4: false}},
			65538: types.Card{CardNumber: 65538, From: date("2020-01-01"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: false, 4: true}},
		},

		54321: map[uint32]types.Card{
			65536: types.Card{CardNumber: 65536, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: true, 4: false}},
			65537: types.Card{CardNumber: 65537, From: date("2020-03-04"), To: date("2020-11-30"), Doors: map[uint8]bool{1: false, 2: true, 3: false, 4: false}},
			65538: types.Card{CardNumber: 65538, From: date("2020-05-06"), To: date("2020-10-29"), Doors: map[uint8]bool{1: true, 2: false, 3: false, 4: false}},
		},
	}

	expected := map[uint32][]types.Card{
		12345: []types.Card{
			types.Card{CardNumber: 65537, From: date("2020-01-01"), To: date("2020-12-31"), Doors: map[uint8]bool{1: true, 2: false, 3: false, 4: false}},
			types.Card{CardNumber: 65538, From: date("2020-01-01"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: false, 4: true}},
			types.Card{CardNumber: 65536, From: date("2020-01-01"), To: date("2020-12-31"), Doors: map[uint8]bool{1: true, 2: false, 3: true, 4: false}},
		},

		54321: []types.Card{
			types.Card{CardNumber: 65536, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: true, 4: false}},
			types.Card{CardNumber: 65537, From: date("2020-03-04"), To: date("2020-11-30"), Doors: map[uint8]bool{1: false, 2: true, 3: false, 4: false}},
			types.Card{CardNumber: 65538, From: date("2020-05-06"), To: date("2020-10-29"), Doors: map[uint8]bool{1: true, 2: false, 3: false, 4: false}},
		},
	}

	report := map[uint32]Report{
		12345: Report{
			Unchanged: []uint32{65537},
			Updated:   []uint32{65538},
			Added:     []uint32{65536},
			Deleted:   []uint32{65539},
			Failed:    []uint32{},
			Errored:   []uint32{},
			Errors:    []error{},
		},

		54321: Report{
			Unchanged: []uint32{65536},
			Updated:   []uint32{65537, 65538},
			Added:     []uint32{},
			Deleted:   []uint32{65539},
			Failed:    []uint32{},
			Errored:   []uint32{},
			Errors:    []error{},
		},
	}

	cards := map[uint32][]types.Card{
		12345: []types.Card{
			types.Card{CardNumber: 65537, From: date("2020-01-01"), To: date("2020-12-31"), Doors: map[uint8]bool{1: true, 2: false, 3: false, 4: false}},
			types.Card{CardNumber: 65538, From: date("2020-01-01"), To: date("2020-11-30"), Doors: map[uint8]bool{1: true, 2: true, 3: true, 4: true}},
			types.Card{CardNumber: 65539, From: date("2020-01-01"), To: date("2020-10-31"), Doors: map[uint8]bool{1: true, 2: true, 3: true, 4: true}},
		},

		54321: []types.Card{
			types.Card{CardNumber: 65536, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: true, 4: false}},
			types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]bool{1: false, 2: false, 3: true, 4: false}},
			types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-11-30"), Doors: map[uint8]bool{1: false, 2: true, 3: false, 4: false}},
			types.Card{CardNumber: 65539, From: date("2020-03-04"), To: date("2020-12-31"), Doors: map[uint8]bool{1: true, 2: false, 3: false, 4: false}},
		},
	}

	u := mock{
		getCards: func(deviceID uint32) (uint32, error) {
			return uint32(len(cards[deviceID])), nil
		},
		getCardByIndex: func(deviceID, index uint32) (*types.Card, error) {
			if int(index) < 0 || int(index) > len(cards[deviceID]) {
				return nil, nil
			}
			return &cards[deviceID][index-1], nil
		},

		putCard: func(deviceID uint32, card types.Card) (bool, error) {
			for ix, c := range cards[deviceID] {
				if c.CardNumber == card.CardNumber {
					cards[deviceID][ix] = card
					return true, nil
				}
			}

			cards[deviceID] = append(cards[deviceID], card)

			return true, nil
		},

		deleteCard: func(deviceID uint32, card types.Card) (bool, error) {
			for ix, c := range cards[deviceID] {
				if c.CardNumber == card.CardNumber {
					cards[deviceID] = append(cards[deviceID][:ix], cards[deviceID][ix+1:]...)
					return true, nil
				}
			}

			return false, nil
		},
	}

	rpt, err := PutACL(&u, acl, false)
	if err != nil {
		t.Fatalf("Unexpected error putting ACL: %v", err)
	}

	if !reflect.DeepEqual(cards, expected) {
		t.Errorf("Device internal card list not updated correctly:\n    expected:%+v\n    got:     %+v", expected, cards)
	}

	if !reflect.DeepEqual(rpt, report) {
		t.Errorf("Returned report does not match expected:\n    expected:%+v\n    got:     %+v", report, rpt)
	}
}

func TestPutACLWithFailures(t *testing.T) {
	acl := ACL{
		12345: map[uint32]types.Card{
			65536: types.Card{CardNumber: 65536, From: date("2020-03-04"), To: date("2020-12-31"), Doors: map[uint8]bool{1: true, 2: false, 3: true, 4: false}},
			65537: types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]bool{1: true, 2: false, 3: false, 4: false}},
			65538: types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: false, 4: false}},
		},
	}

	expected := []types.Card{
		types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]bool{1: true, 2: false, 3: false, 4: false}},
		types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-11-30"), Doors: map[uint8]bool{1: true, 2: false, 3: false, 4: true}},
		types.Card{CardNumber: 65536, From: date("2020-03-04"), To: date("2020-12-31"), Doors: map[uint8]bool{1: true, 2: false, 3: true, 4: false}},
	}

	report := map[uint32]Report{
		12345: Report{
			Unchanged: []uint32{65537},
			Updated:   []uint32{},
			Added:     []uint32{65536},
			Deleted:   []uint32{65539},
			Failed:    []uint32{65538},
			Errored:   []uint32{},
			Errors:    []error{},
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
		putCard: func(deviceID uint32, card types.Card) (bool, error) {
			if card.CardNumber == 65538 {
				return false, nil
			}

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

	rpt, err := PutACL(&u, acl, false)
	if err != nil {
		t.Fatalf("Unexpected error putting ACL: %v", err)
	}

	if !reflect.DeepEqual(cards, expected) {
		t.Errorf("Device internal card list not updated correctly:\n    expected:%+v\n    got:     %+v", expected, cards)
	}

	if !reflect.DeepEqual(rpt, report) {
		t.Errorf("Returned report does not match expected:\n    expected:%+v\n    got:     %+v", report, rpt)
	}
}

func TestPutACLWithErrors(t *testing.T) {
	acl := ACL{
		12345: map[uint32]types.Card{
			65536: types.Card{CardNumber: 65536, From: date("2020-03-04"), To: date("2020-12-31"), Doors: map[uint8]bool{1: true, 2: false, 3: true, 4: false}},
			65537: types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]bool{1: true, 2: false, 3: false, 4: false}},
			65538: types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: false, 4: false}},
		},
	}

	expected := []types.Card{
		types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]bool{1: true, 2: false, 3: false, 4: false}},
		types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-11-30"), Doors: map[uint8]bool{1: true, 2: false, 3: false, 4: true}},
		types.Card{CardNumber: 65536, From: date("2020-03-04"), To: date("2020-12-31"), Doors: map[uint8]bool{1: true, 2: false, 3: true, 4: false}},
	}

	report := map[uint32]Report{
		12345: Report{
			Unchanged: []uint32{65537},
			Updated:   []uint32{},
			Added:     []uint32{65536},
			Deleted:   []uint32{65539},
			Failed:    []uint32{},
			Errored:   []uint32{65538},
			Errors:    []error{fmt.Errorf("Mysterious error updating card %v", 65538)},
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
		putCard: func(deviceID uint32, card types.Card) (bool, error) {
			if card.CardNumber == 65538 {
				return false, fmt.Errorf("Mysterious error updating card %v", card.CardNumber)
			}

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

	rpt, err := PutACL(&u, acl, false)
	if err != nil {
		t.Fatalf("Unexpected error putting ACL: %v", err)
	}

	if !reflect.DeepEqual(cards, expected) {
		t.Errorf("Device internal card list not updated correctly:\n    expected:%+v\n    got:     %+v", expected, cards)
	}

	if !reflect.DeepEqual(rpt, report) {
		t.Errorf("Returned report does not match expected:\n    expected:%+v\n    got:     %+v", report, rpt)
	}
}

func TestPutACLWithNoCurrentPermissions(t *testing.T) {
	acl := ACL{
		12345: map[uint32]types.Card{
			65537: types.Card{CardNumber: 65537, From: date("2020-03-04"), To: date("2020-10-31"), Doors: map[uint8]bool{1: false, 2: false, 3: true, 4: false}},
		},
	}

	expected := []types.Card{
		types.Card{CardNumber: 65537, From: date("2020-03-04"), To: date("2020-10-31"), Doors: map[uint8]bool{1: false, 2: false, 3: true, 4: false}},
	}

	report := map[uint32]Report{
		12345: Report{
			Unchanged: []uint32{},
			Updated:   []uint32{65537},
			Added:     []uint32{},
			Deleted:   []uint32{},
			Failed:    []uint32{},
			Errored:   []uint32{},
			Errors:    []error{},
		},
	}

	cards := []types.Card{
		types.Card{CardNumber: 65537, From: date("2020-01-01"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: false, 4: false}},
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

	rpt, err := PutACL(&u, acl, false)
	if err != nil {
		t.Fatalf("Unexpected error putting ACL: %v", err)
	}

	if !reflect.DeepEqual(cards, expected) {
		t.Errorf("Device internal card list not updated correctly:\n    expected:%+v\n    got:     %+v", expected, cards)
	}

	if !reflect.DeepEqual(rpt, report) {
		t.Errorf("Returned report does not match expected:\n    expected:%+v\n    got:     %+v", report, rpt)
	}
}
