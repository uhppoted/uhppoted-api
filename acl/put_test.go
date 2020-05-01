package acl

import (
	"github.com/uhppoted/uhppote-core/types"
	"reflect"
	"testing"
)

func TestPutACL(t *testing.T) {
	acl := ACL{
		12345: map[uint32]types.Card{
			65536: types.Card{CardNumber: 65536, From: date("2020-03-04"), To: date("2020-12-31"), Doors: []bool{true, false, true, false}},
			65537: types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: []bool{true, false, false, false}},
			65538: types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-12-31"), Doors: []bool{false, false, false, false}},
		},
	}

	expected := []types.Card{
		types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: []bool{true, false, false, false}},
		types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-12-31"), Doors: []bool{false, false, false, false}},
		types.Card{CardNumber: 65536, From: date("2020-03-04"), To: date("2020-12-31"), Doors: []bool{true, false, true, false}},
	}

	report := map[uint32]Report{
		12345: Report{Unchanged: 1, Updated: 1, Added: 1, Deleted: 1, Failed: 0},
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

	rpt, err := PutACL(&u, acl)
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
			65536: types.Card{CardNumber: 65536, From: date("2020-01-01"), To: date("2020-12-31"), Doors: []bool{true, false, true, false}},
			65537: types.Card{CardNumber: 65537, From: date("2020-01-01"), To: date("2020-12-31"), Doors: []bool{true, false, false, false}},
			65538: types.Card{CardNumber: 65538, From: date("2020-01-01"), To: date("2020-12-31"), Doors: []bool{false, false, false, true}},
		},

		54321: map[uint32]types.Card{
			65536: types.Card{CardNumber: 65536, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, true, false}},
			65537: types.Card{CardNumber: 65537, From: date("2020-03-04"), To: date("2020-11-30"), Doors: []bool{false, true, false, false}},
			65538: types.Card{CardNumber: 65538, From: date("2020-05-06"), To: date("2020-10-29"), Doors: []bool{true, false, false, false}},
		},
	}

	expected := map[uint32][]types.Card{
		12345: []types.Card{
			types.Card{CardNumber: 65537, From: date("2020-01-01"), To: date("2020-12-31"), Doors: []bool{true, false, false, false}},
			types.Card{CardNumber: 65538, From: date("2020-01-01"), To: date("2020-12-31"), Doors: []bool{false, false, false, true}},
			types.Card{CardNumber: 65536, From: date("2020-01-01"), To: date("2020-12-31"), Doors: []bool{true, false, true, false}},
		},

		54321: []types.Card{
			types.Card{CardNumber: 65536, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, true, false}},
			types.Card{CardNumber: 65537, From: date("2020-03-04"), To: date("2020-11-30"), Doors: []bool{false, true, false, false}},
			types.Card{CardNumber: 65538, From: date("2020-05-06"), To: date("2020-10-29"), Doors: []bool{true, false, false, false}},
		},
	}

	report := map[uint32]Report{
		12345: Report{Unchanged: 1, Updated: 1, Added: 1, Deleted: 1, Failed: 0},
		54321: Report{Unchanged: 1, Updated: 2, Added: 0, Deleted: 1, Failed: 0},
	}

	cards := map[uint32][]types.Card{
		12345: []types.Card{
			types.Card{CardNumber: 65537, From: date("2020-01-01"), To: date("2020-12-31"), Doors: []bool{true, false, false, false}},
			types.Card{CardNumber: 65538, From: date("2020-01-01"), To: date("2020-11-30"), Doors: []bool{true, true, true, true}},
			types.Card{CardNumber: 65539, From: date("2020-01-01"), To: date("2020-10-31"), Doors: []bool{true, true, true, true}},
		},

		54321: []types.Card{
			types.Card{CardNumber: 65536, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, true, false}},
			types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: []bool{false, false, true, false}},
			types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-11-30"), Doors: []bool{false, true, false, false}},
			types.Card{CardNumber: 65539, From: date("2020-03-04"), To: date("2020-12-31"), Doors: []bool{true, false, false, false}},
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

	rpt, err := PutACL(&u, acl)
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
