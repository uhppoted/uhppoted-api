package acl

import (
	"sort"
)

type Report struct {
	Unchanged []uint32
	Updated   []uint32
	Added     []uint32
	Deleted   []uint32
	Failed    []uint32
	Errored   []uint32
	Errors    []error
}

type ReportSummary map[uint32]struct {
	Unchanged int `json:"unchanged"`
	Updated   int `json:"updated"`
	Added     int `json:"added"`
	Deleted   int `json:"deleted"`
	Failed    int `json:"failed"`
	Errored   int `json:"errored"`
}

type ConsolidatedReport struct {
	Unchanged []uint32 `json:"unchanged"`
	Updated   []uint32 `json:"updated"`
	Added     []uint32 `json:"added"`
	Deleted   []uint32 `json:"deleted"`
	Failed    []uint32 `json:"failed"`
	Errored   []uint32 `json:"errored"`
}

func Summarize(report map[uint32]Report) ReportSummary {
	summary := ReportSummary{}

	for k, v := range report {
		summary[k] = struct {
			Unchanged int `json:"unchanged"`
			Updated   int `json:"updated"`
			Added     int `json:"added"`
			Deleted   int `json:"deleted"`
			Failed    int `json:"failed"`
			Errored   int `json:"errored"`
		}{
			Unchanged: len(v.Unchanged),
			Updated:   len(v.Updated),
			Added:     len(v.Added),
			Deleted:   len(v.Deleted),
			Failed:    len(v.Failed),
			Errored:   len(v.Errored),
		}
	}

	return summary
}

func Consolidate(report map[uint32]Report) ConsolidatedReport {

	// ... consolidate report

	consolidated := map[uint32]*struct {
		updated bool
		added   bool
		deleted bool
		failed  bool
		errored bool
	}{}

	for _, r := range report {
		lists := [][]uint32{r.Updated, r.Added, r.Deleted, r.Failed}
		for _, l := range lists {
			for _, card := range l {
				consolidated[card] = &struct {
					updated bool
					added   bool
					deleted bool
					failed  bool
					errored bool
				}{}
			}
		}
	}

	for _, r := range report {
		for _, card := range r.Updated {
			consolidated[card].updated = true
		}

		for _, card := range r.Added {
			consolidated[card].added = true
		}

		for _, card := range r.Deleted {
			consolidated[card].deleted = true
		}

		for _, card := range r.Failed {
			consolidated[card].failed = true
		}

		for _, card := range r.Errored {
			consolidated[card].errored = true
		}
	}

	updated := []uint32{}
	added := []uint32{}
	deleted := []uint32{}
	failed := []uint32{}
	errored := []uint32{}

	for card, s := range consolidated {
		if s.updated {
			updated = append(updated, card)
		}

		if s.added {
			added = append(added, card)
		}

		if s.deleted {
			deleted = append(deleted, card)
		}

		if s.failed {
			failed = append(failed, card)
		}

		if s.errored {
			errored = append(errored, card)
		}
	}

	f := func(a []uint32) {
		sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
	}

	f(updated)
	f(added)
	f(deleted)
	f(failed)
	f(errored)

	return ConsolidatedReport{
		Updated: updated,
		Added:   added,
		Deleted: deleted,
		Failed:  failed,
		Errored: errored,
	}
}
