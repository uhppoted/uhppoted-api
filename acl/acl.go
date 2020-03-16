package acl

import (
	"github.com/uhppoted/uhppote-core/device"
	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppote-core/uhppote"
	"reflect"
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
		N, err := u.GetCardsN(device.DeviceID)
		if err != nil {
			return acl, err
		}

		for index := uint32(0); index < N; index++ {
			card, err := u.GetCardByIndexN(device.DeviceID, index+1)
			if err != nil {
				return nil, err
			}

			if card != nil {
				acl[device.DeviceID][card.CardNumber] = card.Clone()
			}
		}
	}

	return acl, nil
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

	return diff
}
