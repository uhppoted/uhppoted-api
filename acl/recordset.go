package acl

import (
	"fmt"
	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppote-core/uhppote"
	"sort"
)

type Table struct {
	Header  []string
	Records [][]string
}

func ParseTable(table Table, devices []*uhppote.Device) (*ACL, error) {
	acl := make(ACL)
	for _, device := range devices {
		acl[device.DeviceID] = make(map[uint32]types.Card)
	}

	index, err := parseHeader(table.Header, devices)
	if err != nil {
		return nil, err
	}

	for i, record := range table.Records {
		row := i + 1
		cards, err := parseRecord(record, index)
		if err != nil {
			return nil, fmt.Errorf("Error parsing table - row %d: %w\n", row, err)
		}

		for id, card := range cards {
			if acl[id] != nil {
				if _, ok := acl[id][card.CardNumber]; ok {
					return nil, fmt.Errorf("Duplicate card number (%v)\n", card.CardNumber)
				}

				acl[id][card.CardNumber] = card
			}
		}
	}

	return &acl, nil
}

func MakeTable(acl ACL, devices []*uhppote.Device) (*Table, error) {
	header, err := makeHeader(devices)
	if err != nil {
		return nil, err
	}

	index := map[string]int{}
	for i, h := range header {
		if i > 2 {
			index[clean(h)] = i - 2
		}
	}

	cards := map[uint32]card{}
	for _, d := range devices {
		v, ok := acl[d.DeviceID]
		if !ok {
			return nil, fmt.Errorf("ACL missing for device %v", d.DeviceID)
		}

		jndex := []int{0, 0, 0, 0}
		for i, door := range d.Doors {
			jndex[i] = index[clean(door)]
		}

		for cardno, c := range v {
			record, ok := cards[cardno]
			if !ok {
				record = card{
					cardnumber: c.CardNumber,
					from:       *c.From,
					to:         *c.To,
					doors:      make([]bool, len(index)),
				}
			}

			if c.From.Before(record.from) {
				record.from = *c.From
			}

			if c.To.After(record.to) {
				record.to = *c.To
			}

			for i, door := range c.Doors {
				if ix := jndex[i]; ix == 0 {
					return nil, fmt.Errorf("Missing door ID for device %v, door:%v", d.DeviceID, i+1)
				} else {
					record.doors[ix-1] = door
				}
			}

			cards[cardno] = record
		}
	}

	keys := []uint32{}
	for k, _ := range cards {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

	records := [][]string{}
	for _, k := range keys {
		c := cards[k]
		record := []string{
			fmt.Sprintf("%v", c.cardnumber),
			fmt.Sprintf("%s", c.from),
			fmt.Sprintf("%s", c.to),
		}

		for _, d := range c.doors {
			if d {
				record = append(record, "Y")
			} else {
				record = append(record, "N")
			}
		}

		records = append(records, record)
	}

	rs := Table{
		Header:  header,
		Records: records,
	}

	return &rs, nil
}

func makeHeader(devices []*uhppote.Device) ([]string, error) {
	keys := []uint32{}
	for _, d := range devices {
		keys = append(keys, d.DeviceID)
	}

	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

	header := []string{
		"Card Number",
		"From",
		"To",
	}

	for _, id := range keys {
		for _, d := range devices {
			if d.DeviceID == id {
				header = append(header, d.Doors...)
			}
		}
	}

	return header, nil
}
