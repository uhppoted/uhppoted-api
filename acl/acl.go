package acl

import (
	"github.com/uhppoted/uhppote-core/device"
	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppote-core/uhppote"
	"reflect"
	"sort"
)

type ACL map[uint32]map[uint32]types.Card

type Diff struct {
	Unchanged []types.Card
	Updated   []types.Card
	Added     []types.Card
	Deleted   []types.Card
}

type index struct {
	cardnumber int
	from       int
	to         int
	doors      map[uint32][]int
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

func PutACL(u device.IDevice, acl ACL) error {
	for id, cards := range acl {
		err := putACL(u, id, cards)
		if err != nil {
			return err
		}
	}

	return nil
}

func putACL(u device.IDevice, deviceID uint32, cards map[uint32]types.Card) error {
	current, err := getACL(u, deviceID)
	if err != nil {
		return err
	}

	diff := compare(current, cards)
	updated := []types.Card{}
	added := []types.Card{}
	deleted := []types.Card{}
	failed := []types.Card{}

	for _, card := range diff.Updated {
		if ok, err := u.PutCardN(deviceID, card); err != nil {
			return err
		} else if !ok {
			failed = append(failed, card)
		} else {
			updated = append(updated, card)
		}
	}

	for _, card := range diff.Added {
		if ok, err := u.PutCardN(deviceID, card); err != nil {
			return err
		} else if !ok {
			failed = append(failed, card)
		} else {
			added = append(added, card)
		}
	}

	for _, card := range diff.Deleted {
		if ok, err := u.DeleteCardN(deviceID, card); err != nil {
			return err
		} else if !ok {
			failed = append(failed, card)
		} else {
			deleted = append(deleted, card)
		}
	}

	return nil
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
