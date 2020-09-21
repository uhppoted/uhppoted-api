package acl

import (
	"github.com/uhppoted/uhppote-core/types"
	"reflect"
	"testing"
)

func TestCompare(t *testing.T) {
	src := ACL{
		12345: {
			923321456: types.CardX{CardNumber: 923321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: true, 4: false}},
			823321456: types.CardX{CardNumber: 823321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: true, 4: false}},
			623321456: types.CardX{CardNumber: 623321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: true, 4: false}},
		},
	}

	dest := ACL{
		12345: {
			923321456: types.CardX{CardNumber: 923321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: true, 4: false}},
			823321456: types.CardX{CardNumber: 823321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: false, 4: false}},
			723321456: types.CardX{CardNumber: 723321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: true, 4: false}},
		},
	}

	expected := map[uint32]Diff{
		12345: Diff{
			Unchanged: []types.CardX{
				types.CardX{CardNumber: 923321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: true, 4: false}},
			},
			Updated: []types.CardX{
				types.CardX{CardNumber: 823321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: false, 4: false}},
			},
			Added: []types.CardX{
				types.CardX{CardNumber: 723321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: true, 4: false}},
			},
			Deleted: []types.CardX{
				types.CardX{CardNumber: 623321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: true, 4: false}},
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
			923321456: types.CardX{CardNumber: 923321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: true, 4: false}},
			823321456: types.CardX{CardNumber: 823321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: true, 4: false}},
			623321456: types.CardX{CardNumber: 623321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: true, 4: false}},
		},
		54321: {
			923321456: types.CardX{CardNumber: 923321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: true, 4: false}},
			823321456: types.CardX{CardNumber: 823321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: true, 4: false}},
			723321456: types.CardX{CardNumber: 723321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: true, 4: false}},
		},
	}

	dest := ACL{
		12345: {
			923321456: types.CardX{CardNumber: 923321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: true, 4: false}},
			823321456: types.CardX{CardNumber: 823321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: false, 4: false}},
			723321456: types.CardX{CardNumber: 723321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: true, 4: false}},
		},
		54321: {
			923321456: types.CardX{CardNumber: 923321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: true, 4: false}},
			823321456: types.CardX{CardNumber: 823321456, From: date("2020-01-01"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: true, 4: false}},
			623321456: types.CardX{CardNumber: 623321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: true, 4: false}},
		},
	}

	expected := map[uint32]Diff{
		12345: Diff{
			Unchanged: []types.CardX{
				types.CardX{CardNumber: 923321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: true, 4: false}},
			},
			Updated: []types.CardX{
				types.CardX{CardNumber: 823321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: false, 4: false}},
			},
			Added: []types.CardX{
				types.CardX{CardNumber: 723321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: true, 4: false}},
			},
			Deleted: []types.CardX{
				types.CardX{CardNumber: 623321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: true, 4: false}},
			},
		},
		54321: Diff{
			Unchanged: []types.CardX{
				types.CardX{CardNumber: 923321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: true, 4: false}},
			},
			Updated: []types.CardX{
				types.CardX{CardNumber: 823321456, From: date("2020-01-01"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: true, 4: false}},
			},
			Added: []types.CardX{
				types.CardX{CardNumber: 623321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: true, 4: false}},
			},
			Deleted: []types.CardX{
				types.CardX{CardNumber: 723321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: true, 4: false}},
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
