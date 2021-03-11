package acl

import (
	"sync"

	"github.com/uhppoted/uhppote-core/device"
	"github.com/uhppoted/uhppote-core/types"
)

func PutACL(u device.IDevice, acl ACL, dryrun bool) (map[uint32]Report, []error) {
	errors := []error{}
	report := sync.Map{}

	for id, _ := range acl {
		report.Store(id, Report{
			Unchanged: []uint32{},
			Updated:   []uint32{},
			Added:     []uint32{},
			Deleted:   []uint32{},
			Failed:    []uint32{},
			Errored:   []uint32{},
			Errors:    []error{},
		})
	}

	var wg sync.WaitGroup

	for k, v := range acl {
		id := k
		cards := v

		wg.Add(1)
		go func() {
			var rpt *Report
			var err error

			if dryrun {
				rpt, err = fakePutACL(u, id, cards)
			} else {
				rpt, err = putACL(u, id, cards)
			}

			if rpt != nil {
				report.Store(id, *rpt)
			}

			if err != nil {
				errors = append(errors, err)
			}

			wg.Done()
		}()
	}

	wg.Wait()

	r := map[uint32]Report{}
	report.Range(func(k, v interface{}) bool {
		r[k.(uint32)] = v.(Report)
		return true
	})

	return r, errors
}

func putACL(u device.IDevice, deviceID uint32, cards map[uint32]types.Card) (*Report, error) {
	current, err := getACL(u, deviceID)
	if err != nil {
		return nil, err
	}

	diff := compare(current, cards)

	report := Report{
		Unchanged: []uint32{},
		Updated:   []uint32{},
		Added:     []uint32{},
		Deleted:   []uint32{},
		Failed:    []uint32{},
		Errored:   []uint32{},
		Errors:    []error{},
	}

	for _, card := range diff.Unchanged {
		report.Unchanged = append(report.Unchanged, card.CardNumber)
	}

	for _, card := range diff.Updated {
		if ok, err := u.PutCard(deviceID, card); err != nil {
			report.Errored = append(report.Errored, card.CardNumber)
			report.Errors = append(report.Errors, err)
		} else if !ok {
			report.Failed = append(report.Failed, card.CardNumber)
		} else {
			report.Updated = append(report.Updated, card.CardNumber)
		}
	}

	for _, card := range diff.Added {
		if ok, err := u.PutCard(deviceID, card); err != nil {
			report.Errored = append(report.Errored, card.CardNumber)
			report.Errors = append(report.Errors, err)
		} else if !ok {
			report.Failed = append(report.Failed, card.CardNumber)
		} else {
			report.Added = append(report.Added, card.CardNumber)
		}
	}

	for _, card := range diff.Deleted {
		if ok, err := u.DeleteCard(deviceID, card.CardNumber); err != nil {
			report.Errored = append(report.Errored, card.CardNumber)
			report.Errors = append(report.Errors, err)
		} else if !ok {
			report.Failed = append(report.Failed, card.CardNumber)
		} else {
			report.Deleted = append(report.Deleted, card.CardNumber)
		}
	}

	return &report, nil
}

func fakePutACL(u device.IDevice, deviceID uint32, cards map[uint32]types.Card) (*Report, error) {
	current, err := getACL(u, deviceID)
	if err != nil {
		return nil, err
	}

	diff := compare(current, cards)

	report := Report{
		Unchanged: []uint32{},
		Updated:   []uint32{},
		Added:     []uint32{},
		Deleted:   []uint32{},
		Failed:    []uint32{},
		Errored:   []uint32{},
		Errors:    []error{},
	}

	for _, card := range diff.Unchanged {
		report.Unchanged = append(report.Unchanged, card.CardNumber)
	}

	for _, card := range diff.Updated {
		report.Updated = append(report.Updated, card.CardNumber)
	}

	for _, card := range diff.Added {
		report.Added = append(report.Added, card.CardNumber)
	}

	for _, card := range diff.Deleted {
		report.Deleted = append(report.Deleted, card.CardNumber)
	}

	return &report, nil
}
