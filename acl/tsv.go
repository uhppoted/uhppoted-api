package acl

import (
	"encoding/csv"
	"fmt"
	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppote-core/uhppote"
	"io"
	"strconv"
	"strings"
	"time"
)

func ParseTSV(f io.Reader, devices []*uhppote.Device) (ACL, error) {
	acl := make(ACL)
	for _, device := range devices {
		acl[device.DeviceID] = make(map[uint32]types.Card)
	}

	r := csv.NewReader(f)
	r.Comma = '\t'

	header, err := r.Read()
	if err != nil {
		return nil, err
	}

	index, err := parseHeader(header, devices)
	if err != nil {
		return nil, err
	}

	line := 0
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		line += 1
		cards, err := parseRecord(record, index)
		if err != nil {
			return nil, fmt.Errorf("Error parsing TSV - line %d: %w\n", line, err)
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

	return acl, nil
}

func MakeTSV(acl ACL, devices []*uhppote.Device, f io.Writer) error {
	t, err := MakeTable(acl, devices)
	if err != nil {
		return err
	}

	w := csv.NewWriter(f)
	w.Comma = '\t'

	if err := w.Write(t.header); err != nil {
		return err
	}

	for _, r := range t.records {
		if err := w.Write(r); err != nil {
			return err
		}
	}

	w.Flush()

	return nil
}

func parseHeader(header []string, devices []*uhppote.Device) (*index, error) {
	columns := make(map[string]int)

	index := index{
		cardnumber: 0,
		from:       0,
		to:         0,
		doors:      make(map[uint32][]int),
	}

	for _, d := range devices {
		index.doors[d.DeviceID] = make([]int, 4)
	}

	for c, field := range header {
		key := clean(field)
		ix := c + 1

		if columns[key] != 0 {
			return nil, fmt.Errorf("Duplicate column name '%s'", field)
		}

		columns[key] = ix
	}

	index.cardnumber = columns["cardnumber"]
	index.from = columns["from"]
	index.to = columns["to"]

	for _, device := range devices {
		for i, door := range device.Doors {
			if d := clean(door); d != "" {
				index.doors[device.DeviceID][i] = columns[d]
			}
		}
	}

	if index.cardnumber == 0 {
		return nil, fmt.Errorf("Missing 'Card Number' column")
	}

	if index.from == 0 {
		return nil, fmt.Errorf("Missing 'From' column")
	}

	if index.to == 0 {
		return nil, fmt.Errorf("Missing 'To' column")
	}

	for _, device := range devices {
		for i, door := range device.Doors {
			if d := clean(door); d != "" {
				if index.doors[device.DeviceID][i] == 0 {
					return nil, fmt.Errorf("Missing column for door '%s'", door)
				}
			}
		}
	}

	return &index, nil
}

func parseRecord(record []string, index *index) (map[uint32]types.Card, error) {
	cards := make(map[uint32]types.Card, 0)

	for k, v := range index.doors {
		cardno, err := getCardNumber(record, index)
		if err != nil {
			return nil, err
		}

		from, err := getFromDate(record, index)
		if err != nil {
			return nil, err
		}

		to, err := getToDate(record, index)
		if err != nil {
			return nil, err
		}

		doors, err := getDoors(record, v)
		if err != nil {
			return nil, err
		}

		cards[k] = types.Card{
			CardNumber: cardno,
			From:       *from,
			To:         *to,
			Doors:      doors,
		}
	}

	return cards, nil
}

func getCardNumber(record []string, index *index) (uint32, error) {
	f := field(record, index.cardnumber)
	cardnumber, err := strconv.ParseUint(f, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("Invalid card number '%s' (%w)", f, err)
	}

	return uint32(cardnumber), nil
}

func getFromDate(record []string, index *index) (*types.Date, error) {
	f := field(record, index.from)
	date, err := time.ParseInLocation("2006-01-02", f, time.Local)
	if err != nil {
		return nil, fmt.Errorf("Invalid 'from' date '%s' (%w)", f, err)
	}

	from := types.Date(date)

	return &from, nil
}

func getToDate(record []string, index *index) (*types.Date, error) {
	f := field(record, index.to)
	date, err := time.ParseInLocation("2006-01-02", f, time.Local)
	if err != nil {
		return nil, fmt.Errorf("Invalid 'to' date '%s' (%w)", f, err)
	}

	to := types.Date(date)

	return &to, nil
}

func getDoors(record []string, v []int) ([]bool, error) {
	doors := make([]bool, 4)

	for i, d := range v {
		if d == 0 {
			doors[i] = false
			continue
		}

		switch field(record, d) {
		case "Y":
			doors[i] = true
		case "N":
			doors[i] = false
		default:
			return doors, fmt.Errorf("Expected 'Y/N' for door: '%s'", record[d])
		}
	}

	return doors, nil
}

func field(record []string, ix int) string {
	return strings.TrimSpace(record[ix-1])
}
