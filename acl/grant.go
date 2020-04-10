package acl

import (
	"fmt"
	"github.com/uhppoted/uhppote-core/device"
	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppote-core/uhppote"
	"strings"
)

func Grant(u device.IDevice, devices []*uhppote.Device, cardID uint32, from, to types.Date, doors []string) error {
	m := map[string]struct {
		deviceID uint32
		door     uint8
	}{}

	for _, d := range devices {
		for i, dd := range d.Doors {
			door := strings.ToLower(strings.ReplaceAll(dd, " ", ""))
			if e, ok := m[door]; ok {
				return fmt.Errorf("Ambiguous reference to door '%s': defined for both devices %v and %v", dd, e.deviceID, d.DeviceID)
			}

			m[door] = struct {
				deviceID uint32
				door     uint8
			}{
				deviceID: d.DeviceID,
				door:     uint8(i + 1),
			}
		}
	}

	for _, d := range devices {
		list := []uint8{}

		for _, dd := range doors {
			door := strings.ToLower(strings.ReplaceAll(dd, " ", ""))
			if e, ok := m[door]; ok && e.deviceID == d.DeviceID {
				list = append(list, e.door)
			}
		}

		grant(u, d.DeviceID, cardID, from, to, list)
	}

	return nil
}

func grant(u device.IDevice, deviceID uint32, cardID uint32, from, to types.Date, doors []uint8) error {
	if len(doors) == 0 {
		return nil
	}

	card, err := u.GetCardByIdN(deviceID, cardID)
	if err != nil {
		return err
	} else if card == nil {
		card = &types.Card{
			CardNumber: cardID,
			From:       from,
			To:         to,
			Doors:      []bool{false, false, false, false},
		}
	}

	if card.From.After(from) {
		card.From = from
	}

	if card.To.Before(to) {
		card.To = to
	}

	for _, d := range doors {
		card.Doors[d-1] = true
	}

	if ok, err := u.PutCardN(deviceID, *card); err != nil {
		return err
	} else if !ok {
		return fmt.Errorf("Failed to update access rights for card '%v' on device '%v'", cardID, deviceID)
	}

	return nil
}
