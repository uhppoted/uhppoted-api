package acl

import (
	"github.com/uhppoted/uhppote-core/types"
	"reflect"
	"testing"
)

func TestCompare(t *testing.T) {
	src := ACL{
		12345: {
			923321456: types.Card{CardNumber: 923321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, true, false}},
			823321456: types.Card{CardNumber: 823321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, true, false}},
			623321456: types.Card{CardNumber: 623321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, true, false}},
		},
	}

	dest := ACL{
		12345: {
			923321456: types.Card{CardNumber: 923321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, true, false}},
			823321456: types.Card{CardNumber: 823321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, false, false}},
			723321456: types.Card{CardNumber: 723321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, true, false}},
		},
	}

	expected := map[uint32]Diff{
		12345: Diff{
			Unchanged: []types.Card{
				types.Card{CardNumber: 923321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, true, false}},
			},
			Updated: []types.Card{
				types.Card{CardNumber: 823321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, false, false}},
			},
			Added: []types.Card{
				types.Card{CardNumber: 723321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, true, false}},
			},
			Deleted: []types.Card{
				types.Card{CardNumber: 623321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, true, false}},
			},
		},
	}

	diff, err := Compare(src, dest)
	if err != nil {
		t.Fatalf("Unexpected error comparing ACL: %v", err)
	}

	if diff == nil {
		t.Fatalf("Compare(..) returned 'nil'")
	}

	if !reflect.DeepEqual(diff, expected) {
		t.Fatalf("Compare(..) returned invalid 'diff':\n   expected: %+v\n   got:      %+v", expected, diff)
	}
}

func TestCompareWithMultipleDevices(t *testing.T) {
	src := ACL{
		12345: {
			923321456: types.Card{CardNumber: 923321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, true, false}},
			823321456: types.Card{CardNumber: 823321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, true, false}},
			623321456: types.Card{CardNumber: 623321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, true, false}},
		},
		54321: {
			923321456: types.Card{CardNumber: 923321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, true, false}},
			823321456: types.Card{CardNumber: 823321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, true, false}},
			723321456: types.Card{CardNumber: 723321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, true, false}},
		},
	}

	dest := ACL{
		12345: {
			923321456: types.Card{CardNumber: 923321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, true, false}},
			823321456: types.Card{CardNumber: 823321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, false, false}},
			723321456: types.Card{CardNumber: 723321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, true, false}},
		},
		54321: {
			923321456: types.Card{CardNumber: 923321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, true, false}},
			823321456: types.Card{CardNumber: 823321456, From: date("2020-01-01"), To: date("2020-12-31"), Doors: []bool{false, false, true, false}},
			623321456: types.Card{CardNumber: 623321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, true, false}},
		},
	}

	expected := map[uint32]Diff{
		12345: Diff{
			Unchanged: []types.Card{
				types.Card{CardNumber: 923321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, true, false}},
			},
			Updated: []types.Card{
				types.Card{CardNumber: 823321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, false, false}},
			},
			Added: []types.Card{
				types.Card{CardNumber: 723321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, true, false}},
			},
			Deleted: []types.Card{
				types.Card{CardNumber: 623321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, true, false}},
			},
		},
		54321: Diff{
			Unchanged: []types.Card{
				types.Card{CardNumber: 923321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, true, false}},
			},
			Updated: []types.Card{
				types.Card{CardNumber: 823321456, From: date("2020-01-01"), To: date("2020-12-31"), Doors: []bool{false, false, true, false}},
			},
			Added: []types.Card{
				types.Card{CardNumber: 623321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, true, false}},
			},
			Deleted: []types.Card{
				types.Card{CardNumber: 723321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, true, false}},
			},
		},
	}

	diff, err := Compare(src, dest)
	if err != nil {
		t.Fatalf("Unexpected error comparing ACL: %v", err)
	}

	if diff == nil {
		t.Fatalf("Compare(..) returned 'nil'")
	}

	if !reflect.DeepEqual(diff, expected) {
		t.Fatalf("Compare(..) returned invalid 'diff':\n   expected: %+v\n   got:      %+v", expected, diff)
	}
}
