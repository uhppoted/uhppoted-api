package acl

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppote-core/uhppote"
)

type DuplicateCardError struct {
	CardNumber uint32
}

func (e *DuplicateCardError) Error() string {
	return fmt.Sprintf("%-10v Duplicate card number", e.CardNumber)
}

func ParseTable(table *Table, devices []*uhppote.Device, strict bool) (*ACL, []error, error) {
	acl := make(ACL)
	for _, device := range devices {
		acl[device.DeviceID] = make(map[uint32]types.Card)
	}

	index, err := parseHeader(table.Header, devices)
	if err != nil {
		return nil, nil, err
	}

	list := []map[uint32]types.Card{}
	for row, record := range table.Records {
		cards, err := parseRecord(record, index)
		if err != nil {
			return nil, nil, fmt.Errorf("Error parsing table - row %d: %w", row+1, err)
		}

		list = append(list, cards)
	}

	duplicates := map[uint32]int{}
	for _, cards := range list {
		for _, card := range cards {
			count, _ := duplicates[card.CardNumber]
			duplicates[card.CardNumber] = count + 1
			break
		}
	}

	warnings := []error{}
	for _, cards := range list {
	loop:
		for id, card := range cards {
			if acl[id] != nil {
				if count, _ := duplicates[card.CardNumber]; count > 1 {
					if strict {
						return nil, nil, fmt.Errorf("Duplicate card number (%v)", card.CardNumber)
					} else {
						warning := &DuplicateCardError{card.CardNumber}
						for i := range warnings {
							if reflect.DeepEqual(warnings[i], warning) {
								continue loop
							}
						}

						warnings = append(warnings, warning)
					}

					continue
				}

				acl[id][card.CardNumber] = card
			}
		}
	}

	return &acl, warnings, nil
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
			if clean(door) != "" {
				jndex[i] = index[clean(door)]
			}
		}

		for cardno, c := range v {
			fmt.Printf(">>>> CARD: %#v\n", c)
			record, ok := cards[cardno]
			if !ok {
				record = card{
					cardnumber: c.CardNumber,
					from:       *c.From,
					to:         *c.To,
					doors:      make([]types.Permission, len(index)),
				}
			}

			if c.From.Before(record.from) {
				record.from = *c.From
			}

			if c.To.After(record.to) {
				record.to = *c.To
			}

			for i := uint8(1); i <= 4; i++ {
				ix := jndex[i-1]

				if ix == 0 && clean(d.Doors[i-1]) != "" {
					return nil, fmt.Errorf("Missing door ID for device %v, door:%v", d.DeviceID, i)
				}

				if ix != 0 {
					switch v := c.Doors[i].(type) {
					case bool:
						record.doors[ix-1] = v

					case int:
						switch {
						case v == 0:
							record.doors[ix-1] = false

						case v == 1:
							record.doors[ix-1] = true

						case v > 1 && v < 255:
							record.doors[ix-1] = types.Permission(v)

						default:
							record.doors[ix-1] = false
						}

					case uint:
						switch {
						case v == 0:
							record.doors[ix-1] = false

						case v == 1:
							record.doors[ix-1] = true

						case v > 1 && v < 255:
							record.doors[ix-1] = types.Permission(v)

						default:
							record.doors[ix-1] = false
						}

					case types.Permission:
						vv, _ := v.(uint8)
						switch {
						case vv == 0:
							record.doors[ix-1] = false

						case vv == 1:
							record.doors[ix-1] = true

						case vv > 1 && vv < 255:
							record.doors[ix-1] = v

						default:
							record.doors[ix-1] = false
						}

					default:
						record.doors[ix-1] = false
					}
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
			switch v := d.(type) {
			case bool:
				if v {
					record = append(record, "Y")
				} else {
					record = append(record, "N")
				}

			case int:
				switch {
				case v == 0:
					record = append(record, "N")
				case v == 1:
					record = append(record, "Y")
				case v > 1 && v < 255:
					record = append(record, fmt.Sprintf("%v", v))
				default:
					record = append(record, "N")
				}

			case uint:
				switch {
				case v == 0:
					record = append(record, "N")
				case v == 1:
					record = append(record, "Y")
				case v > 1 && v < 255:
					record = append(record, fmt.Sprintf("%v", v))
				default:
					record = append(record, "N")
				}

			case types.Permission:
				vv, _ := v.(uint8)
				switch {
				case vv == 0:
					record = append(record, "N")
				case vv == 1:
					record = append(record, "Y")
				case vv > 1 && vv < 255:
					record = append(record, fmt.Sprintf("%v", v))
				default:
					record = append(record, "N")
				}

			default:
				if d == nil {
					record = append(record, "N")
				} else {
					record = append(record, "?")
				}
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
				for _, door := range d.Doors {
					if clean(door) != "" {
						header = append(header, strings.TrimSpace(door))
					}
				}
			}
		}
	}

	return header, nil
}
