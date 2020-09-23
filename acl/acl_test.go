package acl

import (
	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppote-core/uhppote"
	"time"
)

type mock struct {
	getCards       func(uint32) (uint32, error)
	getCardByIndex func(uint32, uint32) (*types.Card, error)
	getCardByID    func(uint32, uint32) (*types.Card, error)
	putCard        func(uint32, types.Card) (bool, error)
	deleteCard     func(uint32, types.Card) (bool, error)
}

func (m *mock) GetCardsN(deviceID uint32) (uint32, error) {
	return m.getCards(deviceID)
}

func (m *mock) GetCardByIndex(deviceID, index uint32) (*types.Card, error) {
	return m.getCardByIndex(deviceID, index)
}

func (m *mock) GetCardByIdN(deviceID, cardID uint32) (*types.Card, error) {
	return m.getCardByID(deviceID, cardID)
}

func (m *mock) PutCardN(deviceID uint32, card types.Card) (bool, error) {
	return m.putCard(deviceID, card)
}

func (m *mock) DeleteCardN(deviceID uint32, card types.Card) (bool, error) {
	return m.deleteCard(deviceID, card)
}

var date = func(s string) *types.Date {
	d, _ := time.ParseInLocation("2006-01-02", s, time.Local)
	p := types.Date(d)
	return &p
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
		65537: types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]bool{1: true, 2: false, 3: false, 4: false}},
		65538: types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-11-30"), Doors: map[uint8]bool{1: true, 2: false, 3: false, 4: true}},
		65539: types.Card{CardNumber: 65539, From: date("2020-03-04"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: false, 4: false}},
	},
}

var cardsA = []types.Card{
	types.Card{CardNumber: 65537, From: date("2020-01-02"), To: date("2020-10-31"), Doors: map[uint8]bool{1: true, 2: false, 3: false, 4: false}},
	types.Card{CardNumber: 65538, From: date("2020-02-03"), To: date("2020-11-30"), Doors: map[uint8]bool{1: true, 2: false, 3: false, 4: true}},
	types.Card{CardNumber: 65539, From: date("2020-03-04"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: false, 4: false}},
}
