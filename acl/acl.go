package acl

import (
	"fmt"
	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppote-core/uhppote"
	"strings"
)

type ACL map[uint32]map[uint32]types.Card

type Diff struct {
	Unchanged []types.Card
	Updated   []types.Card
	Added     []types.Card
	Deleted   []types.Card
}

type Report struct {
	Unchanged int
	Updated   int
	Added     int
	Deleted   int
	Failed    int
}

type DateRange struct {
	From types.Date
	To   types.Date
}

type index struct {
	cardnumber int
	from       int
	to         int
	doors      map[uint32][]int
}

type doormap map[string]struct {
	deviceID uint32
	door     uint8
	name     string
}

type card struct {
	cardnumber uint32
	from       types.Date
	to         types.Date
	doors      []bool
}

func clean(s string) string {
	return strings.ReplaceAll(strings.ToLower(s), " ", "")
}

func noAccess(c types.Card) bool {
	for _, d := range c.Doors {
		if d {
			return false
		}
	}

	return true
}

func mapDeviceDoors(devices []*uhppote.Device) (doormap, error) {
	m := doormap{}

	for _, d := range devices {
		for i, dd := range d.Doors {
			door := strings.ToLower(strings.ReplaceAll(dd, " ", ""))
			if e, ok := m[door]; ok {
				return m, fmt.Errorf("Ambiguous reference to door '%s': defined for both devices %v and %v", dd, e.deviceID, d.DeviceID)
			}

			m[door] = struct {
				deviceID uint32
				door     uint8
				name     string
			}{
				deviceID: d.DeviceID,
				door:     uint8(i + 1),
				name:     strings.TrimSpace(dd),
			}
		}
	}

	return m, nil
}
