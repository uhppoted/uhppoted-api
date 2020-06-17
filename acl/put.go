package acl

import (
	"github.com/uhppoted/uhppote-core/device"
	"github.com/uhppoted/uhppote-core/types"
)

func PutACL(u device.IDevice, acl ACL, dryrun bool) (map[uint32]Report, error) {
	report := map[uint32]Report{}
	for id, _ := range acl {
		report[id] = Report{}
	}

	for id, cards := range acl {
		var rpt Report
		var err error

		if dryrun {
			rpt, err = fakePutACL(u, id, cards)
		} else {
			rpt, err = putACL(u, id, cards)
		}

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

func fakePutACL(u device.IDevice, deviceID uint32, cards map[uint32]types.Card) (Report, error) {
	report := Report{}
	current, err := getACL(u, deviceID)
	if err != nil {
		return report, err
	}

	diff := compare(current, cards)

	report.Unchanged = len(diff.Unchanged)

	for _, _ = range diff.Updated {
		report.Updated++
	}

	for _, _ = range diff.Added {
		report.Added++
	}

	for _, _ = range diff.Deleted {
		report.Deleted++
	}

	return report, nil
}
