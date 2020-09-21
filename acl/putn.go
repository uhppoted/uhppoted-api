package acl

import (
//	"github.com/uhppoted/uhppote-core/device"
//	"github.com/uhppoted/uhppote-core/types"
)

// func PutACLN(u device.IDevice, acl ACLN, dryrun bool) (map[uint32]Report, error) {
// 	report := map[uint32]Report{}
// 	for id, _ := range acl {
// 		report[id] = Report{}
// 	}
//
// 	for id, cards := range acl {
// 		var rpt *Report
// 		var err error
//
// 		if dryrun {
// 			rpt, err = fakePutACLN(u, id, cards)
// 		} else {
// 			rpt, err = putACLN(u, id, cards)
// 		}
//
// 		if rpt != nil {
// 			report[id] = *rpt
// 		}
//
// 		if err != nil {
// 			return report, err
// 		}
// 	}
//
// 	return report, nil
// }
//
// func putACLN(u device.IDevice, deviceID uint32, cards map[uint32]types.CardN) (*Report, error) {
// 	current, err := getACLN(u, deviceID)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	diff := compareN(current, cards)
//
// 	report := Report{
// 		Unchanged: []uint32{},
// 		Updated:   []uint32{},
// 		Added:     []uint32{},
// 		Deleted:   []uint32{},
// 		Failed:    []uint32{},
// 		Errored:   []uint32{},
// 		Errors:    []error{},
// 	}
//
// 	for _, card := range diff.Unchanged {
// 		report.Unchanged = append(report.Unchanged, card.CardNumber)
// 	}
//
// 	for _, card := range diff.Updated {
// 		if ok, err := u.PutCardX(deviceID, card); err != nil {
// 			report.Errored = append(report.Errored, card.CardNumber)
// 			report.Errors = append(report.Errors, err)
// 		} else if !ok {
// 			report.Failed = append(report.Failed, card.CardNumber)
// 		} else {
// 			report.Updated = append(report.Updated, card.CardNumber)
// 		}
// 	}
//
// 	for _, card := range diff.Added {
// 		if ok, err := u.PutCardX(deviceID, card); err != nil {
// 			report.Errored = append(report.Errored, card.CardNumber)
// 			report.Errors = append(report.Errors, err)
// 		} else if !ok {
// 			report.Failed = append(report.Failed, card.CardNumber)
// 		} else {
// 			report.Added = append(report.Added, card.CardNumber)
// 		}
// 	}
//
// 	for _, card := range diff.Deleted {
// 		if ok, err := u.DeleteCardX(deviceID, card); err != nil {
// 			report.Errored = append(report.Errored, card.CardNumber)
// 			report.Errors = append(report.Errors, err)
// 		} else if !ok {
// 			report.Failed = append(report.Failed, card.CardNumber)
// 		} else {
// 			report.Deleted = append(report.Deleted, card.CardNumber)
// 		}
// 	}
//
// 	return &report, nil
// }
//
// func fakePutACLN(u device.IDevice, deviceID uint32, cards map[uint32]types.CardN) (*Report, error) {
// 	current, err := getACLN(u, deviceID)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	diff := compareN(current, cards)
//
// 	report := Report{
// 		Unchanged: []uint32{},
// 		Updated:   []uint32{},
// 		Added:     []uint32{},
// 		Deleted:   []uint32{},
// 		Failed:    []uint32{},
// 		Errored:   []uint32{},
// 		Errors:    []error{},
// 	}
//
// 	for _, card := range diff.Unchanged {
// 		report.Unchanged = append(report.Unchanged, card.CardNumber)
// 	}
//
// 	for _, card := range diff.Updated {
// 		report.Updated = append(report.Updated, card.CardNumber)
// 	}
//
// 	for _, card := range diff.Added {
// 		report.Added = append(report.Added, card.CardNumber)
// 	}
//
// 	for _, card := range diff.Deleted {
// 		report.Deleted = append(report.Deleted, card.CardNumber)
// 	}
//
// 	return &report, nil
// }
