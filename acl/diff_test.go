package acl

import (
	"github.com/uhppoted/uhppote-core/types"
	"reflect"
	"testing"
)

func TestConsolidateDiff(t *testing.T) {
	expected := ConsolidatedDiff{
		Unchanged: []uint32{233214569, 923321456},
		Updated:   []uint32{233214568, 823321456},
		Added:     []uint32{233214567, 723321456},
		Deleted:   []uint32{233214566, 623321456},
	}

	diff := SystemDiff{
		12345: Diff{
			Unchanged: []types.Card{
				types.Card{CardNumber: 923321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, true, false}},
				types.Card{CardNumber: 233214569, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, true, false}},
			},
			Updated: []types.Card{
				types.Card{CardNumber: 823321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, false, false}},
				types.Card{CardNumber: 233214568, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, false, false}},
			},
			Added: []types.Card{
				types.Card{CardNumber: 723321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, true, false}},
				types.Card{CardNumber: 233214567, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, true, false}},
			},
			Deleted: []types.Card{
				types.Card{CardNumber: 623321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, true, false}},
				types.Card{CardNumber: 233214566, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, true, false}},
			},
		},
	}

	consolidated := diff.Consolidate()
	if consolidated == nil {
		t.Fatalf("ConsolidateDiff(..) returned 'nil'")
	}

	if !reflect.DeepEqual(consolidated, &expected) {
		t.Fatalf("Compare(..) returned invalid consolidated 'diff':\n   expected: %+v\n   got:      %+v", expected, *consolidated)
	}
}

func TestConsolidateDiffWithMultipleDevices(t *testing.T) {
	expected := ConsolidatedDiff{
		Unchanged: []uint32{233214569, 923321456},
		Updated:   []uint32{233214568, 823321456},
		Added:     []uint32{233214567, 723321456},
		Deleted:   []uint32{233214566, 623321456},
	}

	diff := SystemDiff{
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
				types.Card{CardNumber: 233214569, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, true, false}},
			},
			Updated: []types.Card{
				types.Card{CardNumber: 233214568, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, false, false}},
			},
			Added: []types.Card{
				types.Card{CardNumber: 233214567, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, true, false}},
			},
			Deleted: []types.Card{
				types.Card{CardNumber: 233214566, From: date("2020-01-02"), To: date("2020-12-31"), Doors: []bool{false, false, true, false}},
			},
		},
	}

	consolidated := diff.Consolidate()
	if consolidated == nil {
		t.Fatalf("ConsolidateDiff(..) returned 'nil'")
	}

	if !reflect.DeepEqual(consolidated, &expected) {
		t.Fatalf("Compare(..) returned invalid consolidated 'diff':\n   expected: %+v\n   got:      %+v", expected, *consolidated)
	}
}
