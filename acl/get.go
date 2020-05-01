package acl

import (
	"github.com/uhppoted/uhppote-core/device"
	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppote-core/uhppote"
)

func GetACL(u device.IDevice, devices []*uhppote.Device) (ACL, error) {
	acl := make(ACL)
	for _, device := range devices {
		acl[device.DeviceID] = make(map[uint32]types.Card)
	}

	for _, device := range devices {
		cards, err := getACL(u, device.DeviceID)
		if err != nil {
			return acl, err
		}

		acl[device.DeviceID] = cards
	}

	return acl, nil
}

func getACL(u device.IDevice, deviceID uint32) (map[uint32]types.Card, error) {
	cards := map[uint32]types.Card{}

	N, err := u.GetCardsN(deviceID)
	if err != nil {
		return cards, err
	}

	for index := uint32(0); index < N; index++ {
		card, err := u.GetCardByIndexN(deviceID, index+1)
		if err != nil {
			return nil, err
		}

		if card != nil {
			cards[card.CardNumber] = card.Clone()
		}
	}

	return cards, nil
}
