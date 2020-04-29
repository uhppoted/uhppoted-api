package acl

import (
	"encoding/csv"
	"fmt"
	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppote-core/uhppote"
	"io"
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

	if err := w.Write(t.Header); err != nil {
		return err
	}

	for _, r := range t.Records {
		if err := w.Write(r); err != nil {
			return err
		}
	}

	w.Flush()

	return nil
}
