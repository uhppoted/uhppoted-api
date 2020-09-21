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
			Unchanged: []types.CardX{
				types.CardX{CardNumber: 923321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: true, 4: false}},
				types.CardX{CardNumber: 233214569, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: true, 4: false}},
			},
			Updated: []types.CardX{
				types.CardX{CardNumber: 823321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: false, 4: false}},
				types.CardX{CardNumber: 233214568, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: false, 4: false}},
			},
			Added: []types.CardX{
				types.CardX{CardNumber: 723321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: true, 4: false}},
				types.CardX{CardNumber: 233214567, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: true, 4: false}},
			},
			Deleted: []types.CardX{
				types.CardX{CardNumber: 623321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: true, 4: false}},
				types.CardX{CardNumber: 233214566, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: true, 4: false}},
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
				types.CardX{CardNumber: 233214569, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: true, 4: false}},
			},
			Updated: []types.CardX{
				types.CardX{CardNumber: 233214568, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: false, 4: false}},
			},
			Added: []types.CardX{
				types.CardX{CardNumber: 233214567, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: true, 4: false}},
			},
			Deleted: []types.CardX{
				types.CardX{CardNumber: 233214566, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: true, 4: false}},
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

func TestConsolidateDiffWithAddAndUpdateSameCard(t *testing.T) {
	expected := ConsolidatedDiff{
		Unchanged: []uint32{233214569, 923321456},
		Updated:   []uint32{233214568, 823321456},
		Added:     []uint32{233214567},
		Deleted:   []uint32{233214566, 623321456},
	}

	diff := SystemDiff{
		12345: Diff{
			Unchanged: []types.CardX{
				types.CardX{CardNumber: 923321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: true, 4: false}},
				types.CardX{CardNumber: 233214569, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: true, 4: false}},
			},
			Updated: []types.CardX{
				types.CardX{CardNumber: 823321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: false, 4: false}},
				types.CardX{CardNumber: 233214568, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: false, 4: false}},
			},
			Added: []types.CardX{
				types.CardX{CardNumber: 823321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: true, 4: false}},
				types.CardX{CardNumber: 233214567, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: true, 4: false}},
			},
			Deleted: []types.CardX{
				types.CardX{CardNumber: 623321456, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: true, 4: false}},
				types.CardX{CardNumber: 233214566, From: date("2020-01-02"), To: date("2020-12-31"), Doors: map[uint8]bool{1: false, 2: false, 3: true, 4: false}},
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
