package acl

import (
	"github.com/uhppoted/uhppote-core/device"
	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppote-core/uhppote"
)

func GetACLN(u device.IDevice, devices []*uhppote.Device) (ACLN, error) {
	acl := make(ACLN)
	for _, device := range devices {
		acl[device.DeviceID] = make(map[uint32]types.CardN)
	}

	for _, device := range devices {
		cards, err := getACLN(u, device.DeviceID)
		if err != nil {
			return acl, err
		}

		acl[device.DeviceID] = cards
	}

	return acl, nil
}

func getACLN(u device.IDevice, deviceID uint32) (map[uint32]types.CardN, error) {
	cards := map[uint32]types.CardN{}

	N, err := u.GetCardsN(deviceID)
	if err != nil {
		return cards, err
	}

	var index uint32 = 1
	for count := 0; count < int(N); {
		card, err := u.GetCardByIndex(deviceID, index)
		if err != nil {
			return nil, err
		}

		if card != nil {
			cards[card.CardNumber] = card.CloneN()
			count++
		}

		index++
	}

	return cards, nil
}
