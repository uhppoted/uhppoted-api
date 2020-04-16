package acl

import (
	"fmt"
	"github.com/uhppoted/uhppote-core/device"
	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppote-core/uhppote"
	"reflect"
	"sort"
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

func PutACL(u device.IDevice, acl ACL) (map[uint32]Report, error) {
	report := map[uint32]Report{}
	for id, _ := range acl {
		report[id] = Report{}
	}

	for id, cards := range acl {
		rpt, err := putACL(u, id, cards)
		report[id] = rpt
		if err != nil {
			return report, err
		}
	}

	return report, nil
}

func putACL(u device.IDevice, deviceID uint32, cards map[uint32]types.Card) (Report, error) {
	report := Report{}
	current, err := getACL(u, deviceID)
	if err != nil {
		return report, err
	}

	diff := compare(current, cards)

	report.Unchanged = len(diff.Unchanged)

	for _, card := range diff.Updated {
		if ok, err := u.PutCardN(deviceID, card); err != nil {
			return report, err
		} else if !ok {
			report.Failed++
		} else {
			report.Updated++
		}
	}

	for _, card := range diff.Added {
		if ok, err := u.PutCardN(deviceID, card); err != nil {
			return report, err
		} else if !ok {
			report.Failed++
		} else {
			report.Added++
		}
	}

	for _, card := range diff.Deleted {
		if ok, err := u.DeleteCardN(deviceID, card); err != nil {
			return report, err
		} else if !ok {
			report.Failed++
		} else {
			report.Deleted++
		}
	}

	return report, nil
}

func Compare(src, dst ACL) (map[uint32]Diff, error) {
	m := map[uint32]Diff{}

	for k, _ := range src {
		m[k] = Diff{}
	}

	for k, _ := range dst {
		m[k] = Diff{}
	}

	for k, _ := range m {
		p := src[k]
		q := dst[k]
		m[k] = compare(p, q)
	}

	return m, nil
}

func compare(p, q map[uint32]types.Card) Diff {
	cards := map[uint32]struct{}{}

	for k, _ := range p {
		cards[k] = struct{}{}
	}

	for k, _ := range q {
		cards[k] = struct{}{}
	}

	diff := Diff{
		Unchanged: []types.Card{},
		Added:     []types.Card{},
		Updated:   []types.Card{},
		Deleted:   []types.Card{},
	}

	for k, _ := range cards {
		u, hasu := p[k]
		v, hasv := q[k]

		if hasu && hasv {
			if reflect.DeepEqual(u, v) {
				diff.Unchanged = append(diff.Unchanged, u)
			} else {
				diff.Updated = append(diff.Updated, v)
			}
		} else if !hasu && hasv {
			diff.Added = append(diff.Added, v)
		} else if hasu && !hasv {
			diff.Deleted = append(diff.Deleted, u)
		}
	}

	for _, list := range [][]types.Card{
		diff.Unchanged,
		diff.Added,
		diff.Updated,
		diff.Deleted} {
		sort.Slice(list, func(i, j int) bool { return list[i].CardNumber < list[j].CardNumber })
	}

	return diff
}

func GetCard(u device.IDevice, devices []*uhppote.Device, cardID uint32) (map[string]DateRange, error) {
	acl := map[string]DateRange{}
	lookup, err := mapDeviceDoors(devices)
	if err != nil {
		return acl, err
	}

	for _, device := range devices {
		card, err := u.GetCardByIdN(device.DeviceID, cardID)
		if err != nil {
			return acl, err
		}

		if card != nil {
			for i, d := range card.Doors {
				ix := uint8(i + 1)
				if d {
					for _, v := range lookup {
						if v.deviceID == device.DeviceID && v.door == ix {
							acl[v.name] = DateRange{
								From: card.From,
								To:   card.To,
							}
						}
					}
				}
			}
		}
	}

	return acl, nil
}

func clean(s string) string {
	return strings.ReplaceAll(strings.ToLower(s), " ", "")
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
