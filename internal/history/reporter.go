package history

import (
	"fmt"
	"io"
	"text/tabwriter"
	"time"
)

// PrintReport writes a human-readable summary of entries to w.
func PrintReport(w io.Writer, entries []Entry) {
	if len(entries) == 0 {
		fmt.Fprintln(w, "No history recorded.")
		return
	}
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "TIMESTAMP\tOPENED\tCLOSED")
	for _, e := range entries {
		fmt.Fprintf(tw, "%s\t%s\t%s\n",
			e.Timestamp.Format(time.RFC3339),
			joinOrDash(e.Opened),
			joinOrDash(e.Closed),
		)
	}
	_ = tw.Flush()
}

func joinOrDash(vals []string) string {
	if len(vals) == 0 {
		return "-"
	}
	out := ""
	for i, v := range vals {
		if i > 0 {
			out += ","
		}
		out += v
	}
	return out
}
