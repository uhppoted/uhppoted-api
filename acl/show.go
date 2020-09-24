package acl

import (
	"github.com/uhppoted/uhppote-core/device"
	"github.com/uhppoted/uhppote-core/uhppote"
)

func GetCard(u device.IDevice, devices []*uhppote.Device, cardID uint32) (map[string]DateRange, error) {
	acl := map[string]DateRange{}
	lookup, err := mapDeviceDoors(devices)
	if err != nil {
		return acl, err
	}

	for _, device := range devices {
		card, err := u.GetCardByID(device.DeviceID, cardID)
		if err != nil {
			return acl, err
		}

		if card != nil {
			for ix, d := range card.Doors {
				if d {
					for _, v := range lookup {
						if v.deviceID == device.DeviceID && v.door == ix {
							acl[v.name] = DateRange{
								From: *card.From,
								To:   *card.To,
							}
						}
					}
				}
			}
		}
	}

	return acl, nil
}
