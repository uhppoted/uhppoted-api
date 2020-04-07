package acl

import (
	"fmt"
	"github.com/uhppoted/uhppote-core/uhppote"
	"io"
)

func MakeFlatFile(acl ACL, devices []*uhppote.Device, f io.Writer) error {
	t, err := MakeTable(acl, devices)
	if err != nil {
		return err
	}

	formats := make([]string, len(t.header))

	for i, h := range t.header {
		width := len(h)
		for _, r := range t.records {
			if len(r[i]) > width {
				width = len(r[i])
			}
		}

		formats[i] = fmt.Sprintf("%%-%ds", width)
	}

	separator := ""
	for i, h := range t.header {
		fmt.Fprintf(f, "%s", separator)
		fmt.Fprintf(f, formats[i], h)
		separator = "  "
	}
	fmt.Fprintln(f)

	for _, r := range t.records {
		separator := ""
		for i, v := range r {
			fmt.Fprintf(f, "%s", separator)
			fmt.Fprintf(f, formats[i], v)
			separator = "  "
		}
		fmt.Fprintln(f)
	}

	return nil
}
