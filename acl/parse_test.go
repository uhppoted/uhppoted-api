package acl

import (
	"github.com/uhppoted/uhppote-core/uhppote"
	"reflect"
	"testing"
)

func TestParseHeader(t *testing.T) {
	expected := index{
		cardnumber: 1,
		from:       2,
		to:         3,
		doors: map[uint32][]int{
			12345: []int{6, 5, 7, 4},
		},
	}

	header := []string{"Card Number", "From", "To", "Workshop", "Side Door", "Front Door", "Garage"}

	devices := []*uhppote.Device{
		&uhppote.Device{
			DeviceID: 12345,
			Doors:    []string{"Front Door", "Side Door", "Garage", "Workshop"},
		},
	}

	ix, err := parseHeader(header, devices)
	if err != nil {
		t.Fatalf("Unexpected error parsing header: %v", err)
	} else if ix == nil {
		t.Fatalf("parseHeader returned 'nil'")
	}

	if !reflect.DeepEqual(*ix, expected) {
		t.Errorf("Invalid index\n   expected: %+v\n   got:      %+v", expected, *ix)
	}
}

func TestParseHeaderWithMultipleDevices(t *testing.T) {
	expected := index{
		cardnumber: 1,
		from:       2,
		to:         3,
		doors: map[uint32][]int{
			12345: []int{6, 5, 7, 4},
			54321: []int{8, 9, 10, 11},
		},
	}

	header := []string{"Card Number", "From", "To", "Workshop", "Side Door", "Front Door", "Garage", "D1", "D2", "D3", "D4"}

	devices := []*uhppote.Device{
		&uhppote.Device{
			DeviceID: 12345,
			Doors:    []string{"Front Door", "Side Door", "Garage", "Workshop"},
		},
		&uhppote.Device{
			DeviceID: 54321,
			Doors:    []string{"D1", "D2", "D3", "D4"},
		},
	}

	ix, err := parseHeader(header, devices)
	if err != nil {
		t.Fatalf("Unexpected error parsing header: %v", err)
	} else if ix == nil {
		t.Fatalf("parseHeader returned 'nil'")
	}

	if !reflect.DeepEqual(*ix, expected) {
		t.Errorf("Invalid index\n   expected: %+v\n   got:      %+v", expected, *ix)
	}
}

func TestParseHeaderWithMissingColumn(t *testing.T) {
	expected := index{
		cardnumber: 1,
		from:       2,
		to:         3,
		doors: map[uint32][]int{
			12345: []int{6, 5, 7, 4},
			54321: []int{8, 9, 0, 10},
		},
	}

	header := []string{"Card Number", "From", "To", "Workshop", "Side Door", "Front Door", "Garage", "D1", "D2", "D4"}

	devices := []*uhppote.Device{
		&uhppote.Device{
			DeviceID: 12345,
			Doors:    []string{"Front Door", "Side Door", "Garage", "Workshop"},
		},
		&uhppote.Device{
			DeviceID: 54321,
			Doors:    []string{"D1", "D2", "D3", "D4"},
		},
	}

	ix, err := parseHeader(header, devices)
	if err != nil {
		t.Fatalf("Unexpected error parsing header: %v", err)
	} else if ix == nil {
		t.Fatalf("parseHeader returned 'nil'")
	}

	if !reflect.DeepEqual(*ix, expected) {
		t.Errorf("Invalid index\n   expected: %+v\n   got:      %+v", expected, *ix)
	}
}

func TestParseHeaderWithInvalidColumn(t *testing.T) {
	header := []string{"Card Number", "From", "To", "Workshop", "Side Door", "Front Door", "Garage", "D1", "D2", "D3X", "D4"}

	devices := []*uhppote.Device{
		&uhppote.Device{
			DeviceID: 12345,
			Doors:    []string{"Front Door", "Side Door", "Garage", "Workshop"},
		},
		&uhppote.Device{
			DeviceID: 54321,
			Doors:    []string{"D1", "D2", "D3", "D4"},
		},
	}

	ix, err := parseHeader(header, devices)
	if err == nil {
		t.Fatalf("Expected error parsing header with invalid column: %+v", *ix)
	}
}
